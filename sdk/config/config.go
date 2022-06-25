package config

import (
	"context"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/web"
)

type Sdk struct {
	pingclient   api.StatusWebClient
	configclient api.ConfigWebClient
	wait         chan *struct{}
	lker         sync.Mutex
	keys         map[string]*api.WatchData
	keysnotice   map[string]NoticeHandler
	ctx          context.Context
	cancel       context.CancelFunc
}

//keyvalue: map's key is the key name,map's value is the key's data
//keytype: map's key is the key name,map's value is the type of the key's data
type NoticeHandler func(key, keyvalue, keytype string)

func NewConfigSdk(selfgroup, selfname, servergroup, serverhost string) (*Sdk, error) {
	tmpclient, _ := web.NewWebClient(&web.ClientConfig{}, selfgroup, selfname, servergroup, "admin", serverhost)
	instance := &Sdk{
		pingclient:   api.NewStatusWebClient(tmpclient),
		configclient: api.NewConfigWebClient(tmpclient),
		wait:         make(chan *struct{}, 1),
		keys:         make(map[string]*api.WatchData),
		keysnotice:   make(map[string]NoticeHandler),
	}
	if _, e := instance.pingclient.Ping(context.Background(), &api.Pingreq{Timestamp: time.Now().UnixNano()}, nil); e != nil {
		log.Error(nil, "[config.sdk.init] ping server:", serverhost, "error:", e)
		return nil, e
	}
	go instance.watch(selfgroup, selfname)
	return instance, nil
}
func (instance *Sdk) watch(selfgroup, selfname string) {
	for {
		instance.lker.Lock()
		keys := make(map[string]int32)
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
		resp, e := instance.configclient.Watch(instance.ctx, &api.WatchReq{Groupname: selfgroup, Appname: selfname, Keys: keys}, nil)
		if e != nil && !cerror.Equal(e, cerror.ErrCanceled) && e != context.Canceled {
			log.Error(nil, "[config.sdk.watch] keys:", keys, "error:", e)
			time.Sleep(time.Millisecond * 100)
		} else if e == nil {
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
				instance.keys[key] = data
				notice, ok := instance.keysnotice[key]
				if !ok || notice == nil {
					continue
				}
				notice(key, data.Value, data.ValueType)
			}
			instance.lker.Unlock()
		}
		instance.cancel()
	}
}

//watch the same key will overwrite the old one's notice function
//but the old's cancel function can still work
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
		Version:   -1,
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
