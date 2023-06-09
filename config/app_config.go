package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/fsnotify/fsnotify"
)

// AppConfig can hot update
// this is the config used for this app
type AppConfig struct {
	HandlerTimeout     map[string]map[string]ctime.Duration      `json:"handler_timeout"`      //first key handler path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite     map[string]map[string]string              `json:"web_path_rewrite"`     //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate        map[string][]*publicmids.PathRateConfig   `json:"handler_rate"`         //key path
	Accesses           map[string][]*publicmids.PathAccessConfig `json:"accesses"`             //key path
	TokenSecret        string                                    `json:"token_secret"`         //if don't need token check,this can be ingored
	SessionTokenExpire ctime.Duration                            `json:"session_token_expire"` //if don't need session and token check,this can be ignored
	Service            *ServiceConfig                            `json:"service"`
}
type ServiceConfig struct {
	//add your config here
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
}

// AC -
var AC *AppConfig

var watcher *fsnotify.Watcher

func initlocalapp(notice func(*AppConfig)) {
	data, e := os.ReadFile("./AppConfig.json")
	if e != nil {
		log.Error(nil, "[config.local.app] read config file error:", e)
		Close()
		os.Exit(1)
	}
	AC = &AppConfig{}
	if e = json.Unmarshal(data, AC); e != nil {
		log.Error(nil, "[config.local.app] config file format error:", e)
		Close()
		os.Exit(1)
	}
	validateAppConfig(AC)
	log.Info(nil, "[config.remote.app] new config:", AC)
	if notice != nil {
		notice(AC)
	}
	watcher, e = fsnotify.NewWatcher()
	if e != nil {
		log.Error(nil, "[config.local.app] create watcher for hot update error:", e)
		Close()
		os.Exit(1)
	}
	if e = watcher.Add("./"); e != nil {
		log.Error(nil, "[config.local.app] create watcher for hot update error:", e)
		Close()
		os.Exit(1)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) != "AppConfig.json" || (!event.Has(fsnotify.Create) && !event.Has(fsnotify.Write)) {
					continue
				}
				data, e := os.ReadFile("./AppConfig.json")
				if e != nil {
					log.Error(nil, "[config.local.app] hot update read config file error:", e)
					continue
				}
				c := &AppConfig{}
				if e = json.Unmarshal(data, c); e != nil {
					log.Error(nil, "[config.local.app] hot update config file format error:", e)
					continue
				}
				validateAppConfig(c)
				log.Info(nil, "[config.local.app] new config:", c)
				if notice != nil {
					notice(c)
				}
				AC = c
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(nil, "[config.local.app] hot update watcher error:", err)
			}
		}
	}()
}
func initremoteapp(notice func(*AppConfig), wait chan *struct{}) (stopwatch func()) {
	return RemoteConfigSdk.Watch("AppConfig", func(key, keyvalue, keytype string) {
		//only support json
		if keytype != "json" {
			log.Error(nil, "[config.remote.app] config data can only support json format")
			return
		}
		c := &AppConfig{}
		if e := json.Unmarshal(common.Str2byte(keyvalue), c); e != nil {
			log.Error(nil, "[config.remote.app] config data format error:", e)
			return
		}
		validateAppConfig(c)
		log.Info(nil, "[config.remote.app] new config:", c)
		if notice != nil {
			notice(c)
		}
		AC = c
		select {
		case wait <- nil:
		default:
		}
	})
}
