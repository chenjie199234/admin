package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/common"
	ctime "github.com/chenjie199234/Corelib/util/time"
	"github.com/fsnotify/fsnotify"
)

//AppConfig can hot update
//this is the config used for this app
type AppConfig struct {
	HandlerTimeout map[string]map[string]ctime.Duration `json:"handler_timeout"` //first key handler path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	HandlerRate    []*publicmids.RateConfig             `json:"handler_rate"`
	WhiteIP        []string                             `json:"white_ip"`
	BlackIP        []string                             `json:"black_ip"`
	WebPathRewrite map[string]map[string]string         `json:"web_path_rewrite"` //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	AccessKeys     map[string][]string                  `json:"access_keys"`      //key-specific path,value specific seckey,key-"default",value default seckey
	TokenSecret    string                               `json:"token_secret"`
	TokenExpire    ctime.Duration                       `json:"token_expire"`
	Service        *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
}

//every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
	os.Setenv("TOKEN_SECRET", ac.TokenSecret)
}

//AC -
var AC *AppConfig

var watcher *fsnotify.Watcher

func initlocalapp(notice func(*AppConfig)) {
	data, e := os.ReadFile("./AppConfig.json")
	if e != nil {
		log.Error(nil, "[config.initlocalapp] read config file error:", e)
		Close()
		os.Exit(1)
	}
	AC = &AppConfig{}
	if e = json.Unmarshal(data, AC); e != nil {
		log.Error(nil, "[config.initlocalapp] config file format error:", e)
		Close()
		os.Exit(1)
	}
	validateAppConfig(AC)
	if notice != nil {
		notice(AC)
	}
	watcher, e = fsnotify.NewWatcher()
	if e != nil {
		log.Error(nil, "[config.initlocalapp] create watcher for hot update error:", e)
		Close()
		os.Exit(1)
	}
	if e = watcher.Add("./"); e != nil {
		log.Error(nil, "[config.initlocalapp] create watcher for hot update error:", e)
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
				if filepath.Base(event.Name) != "AppConfig.json" || (event.Op&fsnotify.Create == 0 && event.Op&fsnotify.Write == 0) {
					continue
				}
				data, e := os.ReadFile("./AppConfig.json")
				if e != nil {
					log.Error(nil, "[config.initlocalapp] hot update read config file error:", e)
					continue
				}
				c := &AppConfig{}
				if e = json.Unmarshal(data, c); e != nil {
					log.Error(nil, "[config.initlocalapp] hot update config file format error:", e)
					continue
				}
				validateAppConfig(c)
				if notice != nil {
					notice(c)
				}
				AC = c
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(nil, "[config.initlocalapp] hot update watcher error:", err)
			}
		}
	}()
}
func initremoteapp(notice func(*AppConfig), wait chan *struct{}) (stopwatch func()) {
	return RemoteConfigSdk.Watch("AppConfig", func(key, keyvalue, keytype string) {
		//only support json now,so keytype will be ignore
		c := &AppConfig{}
		if e := json.Unmarshal(common.Str2byte(keyvalue), c); e != nil {
			log.Error(nil, "[config.initremoteapp] config data format error:", e)
			return
		}
		validateAppConfig(c)
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
