package config

import (
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
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
}

// AC -
var AC *AppConfig
