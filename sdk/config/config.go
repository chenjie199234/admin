package config

import (
	"context"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/web"
)

type Sdk struct {
	client     api.ConfigWebClient
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

func NewConfigSdk(selfgroup, selfname, servergroup, serverhost, secret string) (*Sdk, error) {
	tmpclient, e := web.NewWebClient(&web.ClientConfig{}, selfgroup, selfname, servergroup, "admin", serverhost)
	if e != nil {
		return nil, e
	}
	instance := &Sdk{
		client:     api.NewConfigWebClient(tmpclient),
		secret:     secret,
		wait:       make(chan *struct{}, 1),
		keys:       make(map[string]*api.WatchData),
		keysnotice: make(map[string]NoticeHandler),
	}
	go instance.watch(selfgroup, selfname)
	return instance, nil
}
func (instance *Sdk) watch(selfgroup, selfname string) {
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
		resp, e := instance.client.Watch(instance.ctx, &api.WatchReq{Groupname: selfgroup, Appname: selfname, Keys: keys}, nil)
		if e != nil {
			if !cerror.Equal(e, cerror.ErrCanceled) {
				log.Error(nil, "[config.sdk.watch] keys:", keys, e)
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
				log.Error(nil, "[config.sdk.watch] key:", data.Key, "return version 0")
				continue
			}
			if instance.secret != "" {
				plaintext, e := util.Decrypt(instance.secret, data.Value)
				if e != nil {
					broken = true
					log.Error(nil, "[config.sdk.watch] decrypt key:", data.Key, e)
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
