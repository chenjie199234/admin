package discover

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
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
	"github.com/chenjie199234/Corelib/util/name"
	"github.com/chenjie199234/Corelib/web"
)

var (
	ErrMissingEnvProject   = errors.New("missing env ADMIN_SERVICE_PROJECT")
	ErrMissingEnvGroup     = errors.New("missing env ADMIN_SERVICE_GROUP")
	ErrMissingEnvHost      = errors.New("missing env ADMIN_SERVICE_WEB_HOST")
	ErrMissingEnvAccessKey = errors.New("missing env ADMIN_SERVICE_DISCOVER_ACCESS_KEY")
	ErrWrongEnvPort        = errors.New("env ADMIN_SERVICE_WEB_PORT must be number <= 65535")
)

type DiscoverSdk struct {
	target string
	status int32 //0 idle,1 discovering

	lker        *sync.RWMutex
	notices     map[chan *struct{}]*struct{}
	di          cdiscover.DI
	noportaddrs map[string]*cdiscover.RegisterData
	version     cdiscover.Version
	lasterror   error

	ctx    context.Context
	cancel context.CancelFunc

	client       api.AppWebClient
	accesskey    string
	discovermode string
	dnshost      string
	dnsinterval  uint32 //uint seconds
	staticaddrs  []string
	kubernetesns string
	kubernetesls string
	kubernetesfs string
	crpcport     uint32
	cgrpcport    uint32
	webport      uint32
}

// if tlsc is not nil,the tls will be actived
// required env:
// ADMIN_SERVICE_PROJECT
// ADMIN_SERVICE_GROUP
// ADMIN_SERVICE_WEB_HOST
// ADMIN_SERVICE_WEB_PORT
// ADMIN_SERVICE_DISCOVER_ACCESS_KEY
func NewAdminDiscover(selfproject, selfgroup, selfapp string, targetproject, targetgroup, targetapp string, tlsc *tls.Config) (cdiscover.DI, error) {
	targetfullname, e := name.MakeFullName(targetproject, targetgroup, targetapp)
	if e != nil {
		return nil, e
	}
	project, group, host, port, accesskey, e := env()
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
		target: targetfullname,
		status: 1,

		lker:    &sync.RWMutex{},
		notices: make(map[chan *struct{}]*struct{}, 10),

		client:    api.NewAppWebClient(tmpclient),
		accesskey: accesskey,
	}
	sdk.ctx, sdk.cancel = context.WithCancel(context.Background())

	once := make(chan *struct{}, 1)
	go sdk.watch(targetproject, targetgroup, targetapp, once)
	go sdk.run(targetproject, targetgroup, targetapp, once)
	return sdk, nil
}

func (s *DiscoverSdk) Now() {
	if !atomic.CompareAndSwapInt32(&s.status, 0, 1) {
		return
	}
	curdi := s.di
	if curdi != nil {
		curdi.Now()
	}
}
func (s *DiscoverSdk) Stop() {
	s.cancel()
}
func env() (projectname, group string, host string, port int, accesskey string, e error) {
	if str, ok := os.LookupEnv("ADMIN_SERVICE_PROJECT"); ok && str != "<ADMIN_SERVICE_PROJECT>" && str != "" {
		projectname = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvProject
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_GROUP"); ok && str != "<ADMIN_SERVICE_GROUP>" && str != "" {
		group = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvGroup
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_HOST"); ok && str != "<ADMIN_SERVICE_WEB_HOST>" && str != "" {
		host = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvHost
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_PORT"); ok && str != "<ADMIN_SERVICE_WEB_PORT>" && str != "" {
		var e error
		port, e = strconv.Atoi(str)
		if e != nil || port < 0 || port > 65535 {
			return "", "", "", 0, "", ErrWrongEnvPort
		}
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_DISCOVER_ACCESS_KEY"); ok && str != "<ADMIN_SERVICE_DISCOVER_ACCESS_KEY>" && str != "" {
		accesskey = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvAccessKey
	}
	return
}
func (s *DiscoverSdk) watch(project, group, app string, once chan *struct{}) {
	defer func() {
		s.discovermode = ""
		if s.di != nil {
			olddi := s.di
			s.di = nil
			olddi.Stop()
		}
		close(once)
	}()
	for {
		header := make(http.Header)
		header.Set("Access-Key", s.accesskey)
		resp, e := s.client.WatchDiscover(s.ctx, &api.WatchDiscoverReq{
			ProjectName:                project,
			GName:                      group,
			AName:                      app,
			CurDiscoverMode:            s.discovermode,
			CurDnsHost:                 s.dnshost,
			CurDnsInterval:             s.dnsinterval,
			CurStaticAddrs:             s.staticaddrs,
			CurKubernetesNamespace:     s.kubernetesns,
			CurKubernetesLabelselector: s.kubernetesls,
			CurKubernetesFieldselector: s.kubernetesfs,
			CurCrpcPort:                s.crpcport,
			CurCgrpcPort:               s.cgrpcport,
			CurWebPort:                 s.webport,
		}, header)
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
		if resp.DiscoverMode == "Static" && len(resp.StaticAddrs) == 0 {
			log.Error(nil, "[discover.admin] static setting broken", log.String("target", s.target))
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if resp.DiscoverMode == "kubernetes" && (resp.KubernetesNamespace == "" || (resp.KubernetesFieldselector == "" && resp.KubernetesLabelselector == "")) {
			log.Error(nil, "[discover.admin] kubernetes setting broken", log.String("target", s.target))
			time.Sleep(time.Millisecond * 100)
			continue
		}
		s.lker.Lock()
		s.discovermode = resp.DiscoverMode
		s.dnshost = resp.DnsHost
		s.dnsinterval = resp.DnsInterval
		s.staticaddrs = resp.StaticAddrs
		s.kubernetesns = resp.KubernetesNamespace
		s.kubernetesls = resp.KubernetesLabelselector
		s.kubernetesfs = resp.KubernetesFieldselector
		s.crpcport = resp.CrpcPort
		s.cgrpcport = resp.CgrpcPort
		s.webport = resp.WebPort
		if s.di != nil {
			olddi := s.di
			s.di = nil
			go olddi.Stop()
		}
		s.lker.Unlock()
		select {
		case once <- nil:
		default:
		}
	}
}
func (s *DiscoverSdk) run(project, group, app string, once chan *struct{}) {
	for {
		s.lker.Lock()
		if s.di == nil && s.discovermode != "" {
			var e error
			switch s.discovermode {
			case "dns":
				s.di, e = cdiscover.NewDNSDiscover(project, group, app, s.dnshost, time.Duration(s.dnsinterval)*time.Second, int(s.crpcport), int(s.cgrpcport), int(s.webport))
			case "static":
				s.di, e = cdiscover.NewStaticDiscover(project, group, app, s.staticaddrs, int(s.crpcport), int(s.cgrpcport), int(s.webport))
			case "kubernetes":
				s.di, e = cdiscover.NewKubernetesDiscover(project, group, app, s.kubernetesns, s.kubernetesfs, s.kubernetesls, int(s.crpcport), int(s.cgrpcport), int(s.webport))
			default:
				log.Error(nil, "[discover.admin] unknown discover type", log.String("target", s.target))
				time.Sleep(time.Millisecond * 100)
				s.lker.Unlock()
				continue
			}
			if e != nil {
				log.Error(nil, "[discover.admin] create discover failed", log.String("target", s.target))
				time.Sleep(time.Millisecond * 100)
				s.lker.Unlock()
				continue
			}
		}
		if s.di == nil {
			s.lker.Unlock()
			_, ok := <-once
			if !ok {
				return
			}
			continue
		}
		curdi := s.di
		s.lker.Unlock()

		ch, cancel := curdi.GetNotice()
		for {
			if _, ok := <-ch; !ok {
				break
			}
			s.lker.Lock()
			s.noportaddrs, s.version, s.lasterror = curdi.GetAddrs(cdiscover.NotNeed)
			s.lker.Unlock()
			atomic.StoreInt32(&s.status, 0)
			for notice := range s.notices {
				select {
				case notice <- nil:
				default:
				}
			}
		}
		cancel()
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
	for addr, reg := range s.noportaddrs {
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
