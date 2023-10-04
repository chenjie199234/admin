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

	//https://login.dingtalk.com/oauth2/auth?redirect_uri={REDIRECT_URI}&response_type=code&client_id={APPKEY}&scope=openid&state=DingTalk&prompt=consent
	DingTalkOauth2    string `json:"dingtalk_oauth2"`
	DingTalkAppKey    string `json:"dingtalk_app_key"`
	DingTalkAppSecret string `json:"dingtalk_app_secret"`

	WeComOauth2 string `json:"wecom_oauth2"`

	//https://open.feishu.cn/open-apis/authen/v1/authorize?redirect_uri={REDIRECT_URI}&app_id={APPID}&state=FeiShu&scope=contact:user.employee_id:readonly%20contact:user.phone:readonly
	FeiShuOauth2    string `json:"feishu_oauth2"`
	FeiShuAppID     string `json:"feishu_app_id"`
	FeiShuAppSecret string `json:"feishu_app_secret"`
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
	if ac.Service.FeiShuAppSecret != "" {
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
	if ac.Service.FeiShuOauth2 != "" && (ac.Service.FeiShuAppID == "" || ac.Service.FeiShuAppSecret == "") {
		log.Error(nil, "[config.validateAppConfig] missing feishu_app_id or feishu_app_secret setting")
		Close()
		os.Exit(1)
	}
}

// AC -
var AC *AppConfig
