package config

import (
	"os"

	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

// AppConfig can hot update
// this is the config used for this app
type AppConfig struct {
	HandlerTimeout     map[string]map[string]ctime.Duration `json:"handler_timeout"`      //first key handler path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite     map[string]map[string]string         `json:"web_path_rewrite"`     //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate        publicmids.MultiPathRateConfigs      `json:"handler_rate"`         //key:path
	Accesses           publicmids.MultiPathAccessConfigs    `json:"accesses"`             //key:path
	TokenSecret        string                               `json:"token_secret"`         //if don't need token check,this can be ingored
	SessionTokenExpire ctime.Duration                       `json:"session_token_expire"` //if don't need session and token check,this can be ignored
	Service            *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
	DingTalkOauth2    string `json:"dingtalk_oauth2"`
	DingTalkAppKey    string `json:"dingtalk_app_key"`
	DingTalkAppSecret string `json:"dingtalk_app_secret"`
	WeComOauth2       string `json:"wecom_oauth2"`
	LarkOauth2        string `json:"lark_oauth2"`
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
	oauth2count := 0
	if ac.Service.DingTalkOauth2 != "" {
		oauth2count++
	}
	if ac.Service.WeComOauth2 != "" {
		oauth2count++
	}
	if ac.Service.LarkOauth2 != "" {
		oauth2count++
	}
	if oauth2count == 0 {
		log.Warn(nil, "[config.validateAppConfig] no oauth2 service,only root account can login by password")
	} else if oauth2count > 1 {
		log.Error(nil, "[config.validateAppConfig] too many oauth2 service")
		Close()
		os.Exit(1)
	}
	if ac.Service.DingTalkOauth2 != "" && (ac.Service.DingTalkAppKey == "" || ac.Service.DingTalkAppSecret == "") {
		log.Error(nil, "[config.validateAppConfig] missing dingtalk_app_key or dingtalk_app_secret setting")
		Close()
		os.Exit(1)
	}
}

// AC -
var AC *AppConfig
