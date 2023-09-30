package discover

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chenjie199234/admin/api"

	"github.com/chenjie199234/Corelib/cerror"
	cdiscover "github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/chenjie199234/Corelib/util/name"
	"github.com/chenjie199234/Corelib/web"
)

var (
	ErrMissingEnvProject = errors.New("missing env ADMIN_SERVICE_PROJECT")
	ErrMissingEnvGroup   = errors.New("missing env ADMIN_SERVICE_GROUP")
	ErrMissingEnvHost    = errors.New("missing env ADMIN_SERVICE_WEB_HOST")
	ErrWrongEnvPort      = errors.New("env ADMIN_SERVICE_WEB_PORT must be number <= 65535")
)

type DiscoverSdk struct {
	target string
	client api.AppWebClient
	status int32 //0 idle,1 discovering
	triger chan *struct{}

	ctx    context.Context
	cancel context.CancelFunc

	lker         *sync.RWMutex
	notices      map[chan *struct{}]*struct{}
	version      int64
	discovermode string
	dnshost      string
	dnsinterval  uint32 //uint seconds
	addrs        []string
	crpcport     uint32
	cgrpcport    uint32
	webport      uint32
	lasterror    error
}

// if tlsc is not nil,the tls will be actived
// required env:
// ADMIN_SERVICE_PROJECT
// ADMIN_SERVICE_GROUP
// ADMIN_SERVICE_WEB_HOST
// ADMIN_SERVICE_WEB_PORT
func NewAdminDiscover(selfproject, selfgroup, selfapp string, targetproject, targetgroup, targetapp string, tlsc *tls.Config) (cdiscover.DI, error) {
	targetfullname, e := name.MakeFullName(targetproject, targetgroup, targetapp)
	if e != nil {
		return nil, e
	}
	project, group, host, port, e := env()
	if e != nil {
		return nil, e
	}
	di, e := cdiscover.NewStaticDiscover(project, group, "admin", []string{host}, 0, 0, port)
	if e != nil {
		return nil, e
	}
	tmpclient, e := web.NewWebClient(nil, di, selfproject, selfgroup, selfapp, project, group, "admin", tlsc)
	if e != nil {
		return nil, e
	}
	sdk := &DiscoverSdk{
		target:  targetfullname,
		client:  api.NewAppWebClient(tmpclient),
		status:  1,
		triger:  make(chan *struct{}, 1),
		lker:    &sync.RWMutex{},
		notices: make(map[chan *struct{}]*struct{}, 10),
	}
	sdk.ctx, sdk.cancel = context.WithCancel(context.Background())
	go sdk.watch(targetproject, targetgroup, targetapp)
	go sdk.run()
	return sdk, nil
}

func (s *DiscoverSdk) Now() {
	if !atomic.CompareAndSwapInt32(&s.status, 0, 1) {
		return
	}
	select {
	case s.triger <- nil:
	default:
	}
}
func (s *DiscoverSdk) Stop() {
	s.cancel()
}
func env() (projectname, group string, host string, port int, e error) {
	if str, ok := os.LookupEnv("ADMIN_SERVICE_PROJECT"); ok && str != "<ADMIN_SERVICE_PROJECT>" && str != "" {
		projectname = str
	} else {
		return "", "", "", 0, ErrMissingEnvProject
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_GROUP"); ok && str != "<ADMIN_SERVICE_GROUP>" && str != "" {
		group = str
	} else {
		return "", "", "", 0, ErrMissingEnvGroup
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_HOST"); ok && str != "<ADMIN_SERVICE_WEB_HOST>" && str != "" {
		host = str
	} else {
		return "", "", "", 0, ErrMissingEnvHost
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_PORT"); ok && str != "<ADMIN_SERVICE_WEB_PORT>" && str != "" {
		var e error
		port, e = strconv.Atoi(str)
		if e != nil || port < 0 || port > 65535 {
			return "", "", "", 0, ErrWrongEnvPort
		}
	}
	return
}
func (s *DiscoverSdk) watch(project, group, app string) {
	for {
		resp, e := s.client.WatchDiscover(s.ctx, &api.WatchDiscoverReq{
			ProjectName:     project,
			GName:           group,
			AName:           app,
			CurDiscoverMode: s.discovermode,
			CurDnsHost:      s.dnshost,
			CurDnsInterval:  s.dnsinterval,
			CurAddrs:        s.addrs,
			CrpcPort:        s.crpcport,
			CgrpcPort:       s.cgrpcport,
			WebPort:         s.webport,
		}, nil)
		if e != nil {
			if cerror.Equal(e, cerror.ErrCanceled) {
				return
			}
			log.Error(nil, "[discover.admin] watch failed", log.String("target", s.target), log.CError(e))
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if resp.DiscoverMode == "dns" && (resp.DnsHost == "" || resp.DnsInterval == 0) {
			log.Error(nil, "[discover.admin] dns setting broken", log.String("target", s.target))
			time.Sleep(time.Millisecond * 100)
			continue
		}
		s.lker.Lock()
		s.discovermode = resp.DiscoverMode
		s.dnshost = resp.DnsHost
		s.dnsinterval = resp.DnsInterval
		s.crpcport = resp.CrpcPort
		s.cgrpcport = resp.CgrpcPort
		s.webport = resp.WebPort
		if s.discovermode != "dns" {
			s.addrs = resp.Addrs
			s.version = time.Now().UnixNano()
			s.lasterror = nil
		}
		s.lker.Unlock()
		select {
		case s.triger <- nil:
		default:
		}
	}
}
func (s *DiscoverSdk) run() {
	tmer := time.NewTimer(0)
	<-tmer.C
	defer func() {
		s.lker.Lock()
		for notice := range s.notices {
			delete(s.notices, notice)
			close(notice)
		}
		s.lker.Unlock()
	}()
	for {
		var dnshost string
		var dnsinterval time.Duration
		var version int64
		select {
		case <-s.ctx.Done():
			log.Info(nil, "[discover.admin] discover stopped", log.String("target", s.target))
			s.lasterror = cerror.ErrDiscoverStopped
			return
		case <-tmer.C:
		case <-s.triger:
		}
		s.status = 1
		if !tmer.Stop() {
			for len(tmer.C) > 0 {
				<-tmer.C
			}
		}
		s.lker.RLock()
		if s.discovermode == "dns" {
			//copy the current dns discover setting
			//the watch may change the current setting when we do dns look up
			dnshost = s.dnshost
			dnsinterval = time.Duration(s.dnsinterval) * time.Second
			version = s.version
		} else {
			for notice := range s.notices {
				select {
				case notice <- nil:
				default:
				}
			}
			s.status = 0
		}
		s.lker.RUnlock()
		if dnshost == "" {
			continue
		}
		tmer.Reset(dnsinterval)
		addrs, e := net.DefaultResolver.LookupHost(s.ctx, dnshost)
		if e != nil && cerror.Equal(errors.Unwrap(e), cerror.ErrCanceled) {
			log.Info(nil, "[discover.admin] discover stopped", log.String("target", s.target))
			s.lasterror = cerror.ErrDiscoverStopped
			return
		}
		if e != nil {
			log.Error(nil, "[discover.admin] dns look up failed",
				log.String("target", s.target),
				log.String("host", dnshost),
				log.CDuration("interval", ctime.Duration(dnsinterval)),
				log.CError(e))
			s.lker.Lock()
			s.lasterror = e
			for notice := range s.notices {
				select {
				case notice <- nil:
				default:
				}
			}
			s.status = 0
			s.lker.Unlock()
			continue
		}
		s.lker.Lock()
		if version != s.version {
			s.lker.Unlock()
			//watch changed the version
			s.status = 0
			continue
		}
		different := len(addrs) != len(s.addrs)
		if !different {
			for _, newaddr := range addrs {
				find := false
				for _, oldaddr := range s.addrs {
					if newaddr == oldaddr {
						find = true
						break
					}
				}
				if !find {
					different = true
					break
				}
			}
		}
		if different {
			s.addrs = addrs
			s.version = time.Now().UnixNano()
		}
		s.lasterror = nil
		for notice := range s.notices {
			select {
			case notice <- nil:
			default:
			}
		}
		s.status = 0
		s.lker.Unlock()
	}
}

// don't close the returned channel,it will be closed in cases:
// 1.the cancel function be called
// 2.this discover stopped
func (s *DiscoverSdk) GetNotice() (notice <-chan *struct{}, cancel func()) {
	ch := make(chan *struct{}, 1)
	s.lker.Lock()
	if s.status == 0 {
		ch <- nil
		s.notices[ch] = nil
	} else {
		select {
		case <-s.ctx.Done():
			close(ch)
		default:
			s.notices[ch] = nil
		}
	}
	s.lker.Unlock()
	return ch, func() {
		s.lker.Lock()
		if _, ok := s.notices[ch]; ok {
			delete(s.notices, ch)
			close(ch)
		}
		s.lker.Unlock()
	}
}

func (s *DiscoverSdk) GetAddrs(pt cdiscover.PortType) (addrs map[string]*cdiscover.RegisterData, version cdiscover.Version, lasterror error) {
	s.lker.RLock()
	defer s.lker.RUnlock()
	r := make(map[string]*cdiscover.RegisterData)
	reg := &cdiscover.RegisterData{
		DServers: map[string]*struct{}{"admin": nil},
	}
	for _, addr := range s.addrs {
		switch pt {
		case cdiscover.NotNeed:
		case cdiscover.Crpc:
			if s.crpcport > 0 {
				if strings.Contains(addr, ":") {
					//ipv6
					addr = "[" + addr + "]:" + strconv.FormatUint(uint64(s.crpcport), 10)
				} else {
					//ipv4
					addr = addr + ":" + strconv.FormatUint(uint64(s.crpcport), 10)
				}
			}
		case cdiscover.Cgrpc:
			if s.cgrpcport > 0 {
				if strings.Contains(addr, ":") {
					//ipv6
					addr = "[" + addr + "]:" + strconv.FormatUint(uint64(s.cgrpcport), 10)
				} else {
					//ipv4
					addr = addr + ":" + strconv.FormatUint(uint64(s.cgrpcport), 10)
				}
			}
		case cdiscover.Web:
			if s.webport > 0 {
				if strings.Contains(addr, ":") {
					//ipv6
					addr = "[" + addr + "]:" + strconv.FormatUint(uint64(s.webport), 10)
				} else {
					//ipv4
					addr = addr + ":" + strconv.FormatUint(uint64(s.webport), 10)
				}
			}
		}
		r[addr] = reg
	}
	return r, s.version, s.lasterror
}
func (s *DiscoverSdk) CheckTarget(target string) bool {
	return s.target == target
}
