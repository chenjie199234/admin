package config

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/web"
)

type Sdk struct {
	client     api.AppWebClient
	secret     string
	wait       chan *struct{}
	lker       sync.Mutex
	keys       map[string]*api.WatchData
	keysnotice map[string]NoticeHandler
	ctx        context.Context
	cancel     context.CancelFunc
}

// keyvalue: map's key is the key name,map's value is the key's data
// keytype: map's key is the key name,map's value is the type of the key's data
type NoticeHandler func(key, keyvalue, keytype string)

var (
	ErrMissingEnvPROJECT = errors.New("missing env REMOTE_CONFIG_SERVICE_PROJECT")
	ErrMissingEnvGroup   = errors.New("missing env REMOTE_CONFIG_SERVICE_GROUP")
	ErrMissingEnvHost    = errors.New("missing env REMOTE_CONFIG_SERVICE_WEB_HOST")
	ErrWrongEnvPort      = errors.New("env REMOTE_CONFIG_SERVICE_WEB_PORT must be number <= 65535")
	ErrWrongEnvSecret    = errors.New("env REMOTE_CONFIG_SECRET too long")
)

// if tlsc is not nil,the tls will be actived
// must set below env:
// REMOTE_CONFIG_SERVICE_PROJECT
// REMOTE_CONFIG_SERVICE_GROUP
// REMOTE_CONFIG_SERVICE_WEB_HOST
// REMOTE_CONFIG_SERVICE_WEB_PORT
// REMOTE_CONFIG_SECRET
func NewConfigSdk(selfprojectname, selfappgroup, selfappname string, tlsc *tls.Config) (*Sdk, error) {
	serverprojectname, group, host, port, secret, e := env()
	if e != nil {
		return nil, e
	}
	di, e := discover.NewStaticDiscover(serverprojectname, group, "admin", []string{host}, 0, 0, port)
	if e != nil {
		return nil, e
	}
	tmpclient, e := web.NewWebClient(nil, di, selfprojectname, selfappgroup, selfappname, serverprojectname, group, "admin", tlsc)
	if e != nil {
		return nil, e
	}
	instance := &Sdk{
		client:     api.NewAppWebClient(tmpclient),
		secret:     secret,
		wait:       make(chan *struct{}, 1),
		keys:       make(map[string]*api.WatchData),
		keysnotice: make(map[string]NoticeHandler),
	}
	go instance.watch(selfprojectname, selfappgroup, selfappname)
	return instance, nil
}
func env() (projectname string, group string, host string, port int, secret string, e error) {
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SERVICE_PROJECT"); ok && str != "<REMOTE_CONFIG_SERVICE_PROJECT>" && str != "" {
		projectname = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvPROJECT
	}
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SERVICE_GROUP"); ok && str != "<REMOTE_CONFIG_SERVICE_GROUP>" && str != "" {
		group = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvGroup
	}
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SERVICE_WEB_HOST"); ok && str != "<REMOTE_CONFIG_SERVICE_WEB_HOST>" && str != "" {
		host = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvHost
	}
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SERVICE_WEB_PORT"); ok && str != "<REMOTE_CONFIG_SERVICE_WEB_PORT>" && str != "" {
		var e error
		port, e = strconv.Atoi(str)
		if e != nil || port < 0 || port > 65535 {
			return "", "", "", 0, "", ErrWrongEnvPort
		}
	}
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SECRET"); ok && str != "<REMOTE_CONFIG_SECRET>" && str != "" {
		secret = str
	}
	if len(secret) >= 32 {
		return "", "", "", 0, "", ErrWrongEnvSecret
	}
	return
}
func (instance *Sdk) watch(selfprojectname, selfappgroup, selfappname string) {
	for {
		instance.lker.Lock()
		keys := make(map[string]uint32)
		for k, v := range instance.keys {
			keys[k] = v.Version
		}
		if len(keys) == 0 {
			instance.lker.Unlock()
			<-instance.wait
			continue
		}
		instance.ctx, instance.cancel = context.WithCancel(context.Background())
		instance.lker.Unlock()
		resp, e := instance.client.Watch(instance.ctx, &api.WatchReq{ProjectName: selfprojectname, GName: selfappgroup, AName: selfappname, Keys: keys}, nil)
		if e != nil {
			if !cerror.Equal(e, cerror.ErrCanceled) {
				log.Error(nil, "[config.sdk.watch] failed", map[string]interface{}{"error": e, "watch_keys": keys})
				time.Sleep(time.Millisecond * 100)
				instance.cancel()
			}
			continue
		}
		broken := false
		instance.lker.Lock()
		for key, data := range resp.Datas {
			if keys[key] == data.Version {
				//didn't changed
				continue
			}
			_, ok := instance.keys[key]
			if !ok {
				//already deleted
				continue
			}
			if data.Version == 0 {
				log.Error(nil, "[config.sdk.watch] key's value's version == 0", map[string]interface{}{"key": data.Key})
				continue
			}
			if instance.secret != "" {
				plaintext, e := secure.AesDecrypt(instance.secret, data.Value)
				if e != nil {
					broken = true
					log.Error(nil, "[config.sdk.watch] decrypt keys's value failed", map[string]interface{}{"key": data.Key, "error": e})
					continue
				}
				data.Value = common.Byte2str(plaintext)
			}
			instance.keys[key] = data
			notice, ok := instance.keysnotice[key]
			if !ok || notice == nil {
				continue
			}
			notice(key, data.Value, data.ValueType)
		}
		instance.lker.Unlock()
		if broken {
			time.Sleep(time.Millisecond * 100)
		}
		instance.cancel()
	}
}

// Warning!!!Don't block in notice function
// watch the same key will overwrite the old one's notice function
// but the old's cancel function can still work
func (instance *Sdk) Watch(key string, notice NoticeHandler) (cancel func()) {
	instance.lker.Lock()
	defer instance.lker.Unlock()
	if _, ok := instance.keys[key]; ok {
		instance.keysnotice[key] = notice
		return func() {
			instance.lker.Lock()
			if _, ok := instance.keys[key]; ok {
				delete(instance.keys, key)
				delete(instance.keysnotice, key)
				if instance.cancel != nil {
					instance.cancel()
				}
			}
			instance.lker.Unlock()
		}
	}
	instance.keys[key] = &api.WatchData{
		Key:       key,
		Value:     "",
		ValueType: "raw",
		Version:   0,
	}
	instance.keysnotice[key] = notice
	if instance.cancel != nil {
		instance.cancel()
	}
	select {
	case instance.wait <- nil:
	default:
	}
	return func() {
		instance.lker.Lock()
		if _, ok := instance.keys[key]; ok {
			delete(instance.keys, key)
			delete(instance.keysnotice, key)
			if instance.cancel != nil {
				instance.cancel()
			}
		}
		instance.lker.Unlock()
	}
}
