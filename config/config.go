package config

import (
	"os"
	"strconv"
	"time"

	"github.com/chenjie199234/admin/config/internal"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/log"
)

// EnvConfig can't hot update,all these data is from system env setting
// nil field means that system env not exist
type EnvConfig struct {
	ConfigType *int
	RunEnv     *string
	DeployEnv  *string
}

// EC -
var EC *EnvConfig

// RemoteConfigSdk -
var RemoteConfigSdk *internal.Sdk

// notice is a sync function
// don't write block logic inside it
func Init(notice func(c *AppConfig), AppConfigTemplate, SourceConfigTemplate []byte) {
	initenv(AppConfigTemplate, SourceConfigTemplate)
	if EC.ConfigType != nil && *EC.ConfigType == 1 {
		tmer := time.NewTimer(time.Second * 2)
		waitapp := make(chan *struct{})
		waitsource := make(chan *struct{})
		initremoteapp(notice, waitapp)
		initremotesource(waitsource)
		appinit := false
		sourceinit := false
		for {
			select {
			case <-waitapp:
				appinit = true
			case <-waitsource:
				sourceinit = true
			case <-tmer.C:
				log.Error(nil, "[config.initremote] timeout")
				Close()
				os.Exit(1)
			}
			if appinit && sourceinit {
				break
			}
		}
	} else {
		initlocalapp(notice)
		initlocalsource()
	}
}

// Close -
func Close() {
	log.Close()
}

func initenv(AppConfigTemplate, SourceConfigTemplate []byte) {
	EC = &EnvConfig{}
	if str, ok := os.LookupEnv("CONFIG_TYPE"); ok && str != "<CONFIG_TYPE>" && str != "" {
		configtype, e := strconv.Atoi(str)
		if e != nil || (configtype != 0 && configtype != 1 && configtype != 2) {
			log.Error(nil, "[config.initenv] env CONFIG_TYPE must be number in [0,1,2]")
			Close()
			os.Exit(1)
		}
		EC.ConfigType = &configtype
	} else {
		log.Warning(nil, "[config.initenv] missing env CONFIG_TYPE")
	}
	if EC.ConfigType != nil && *EC.ConfigType == 1 {
		var mongourl string
		if str, ok := os.LookupEnv("REMOTE_CONFIG_MONGO_URL"); ok && str != "<REMOTE_CONFIG_MONGO_URL>" && str != "" {
			mongourl = str
		} else {
			log.Error(nil, "[config.initenv] missing env REMOTE_CONFIG_MONGO_URL")
			Close()
			os.Exit(1)
		}
		var secret string
		if str, ok := os.LookupEnv("REMOTE_CONFIG_SECRET"); ok && str != "<REMOTE_CONFIG_SECRET>" && str != "" {
			secret = str
		}
		if len(secret) >= 32 {
			log.Error(nil, "[config.initenv] REMOTE_CONFIG_SECRET length too long")
			Close()
			os.Exit(1)
		}
		var e error
		if RemoteConfigSdk, e = internal.NewDirectSdk(model.Group, model.Name, mongourl, secret, AppConfigTemplate, SourceConfigTemplate); e != nil {
			log.Error(nil, "[config.initenv] new remote config sdk:", e)
			Close()
			os.Exit(1)
		}
	}
	if str, ok := os.LookupEnv("RUN_ENV"); ok && str != "<RUN_ENV>" && str != "" {
		EC.RunEnv = &str
	} else {
		log.Warning(nil, "[config.initenv] missing env RUN_ENV")
	}
	if str, ok := os.LookupEnv("DEPLOY_ENV"); ok && str != "<DEPLOY_ENV>" && str != "" {
		EC.DeployEnv = &str
	} else {
		log.Warning(nil, "[config.initenv] missing env DEPLOY_ENV")
	}
}
