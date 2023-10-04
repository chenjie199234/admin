package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"os"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/initinternal"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/mongo"
	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/ctime"
)

// EnvConfig can't hot update,all these data is from system env setting
// nil field means that system env not exist
type EnvConfig struct {
	RunEnv    *string
	DeployEnv *string
}

// EC -
var EC *EnvConfig

var Sdk *initinternal.InternalSdk

// notice is a sync function
// don't write block logic inside it
func Init(notice func(c *AppConfig)) {
	EC = &EnvConfig{}
	if str, ok := os.LookupEnv("RUN_ENV"); ok && str != "<RUN_ENV>" && str != "" {
		EC.RunEnv = &str
	} else {
		log.Warn(nil, "[config.Init] missing env RUN_ENV")
	}
	if str, ok := os.LookupEnv("DEPLOY_ENV"); ok && str != "<DEPLOY_ENV>" && str != "" {
		EC.DeployEnv = &str
	} else {
		log.Warn(nil, "[config.Init] missing env DEPLOY_ENV")
	}
	var secret string
	if str, ok := os.LookupEnv("CONFIG_SECRET"); ok && str != "<CONFIG_SECRET>" && str != "" {
		secret = str
	}
	tmer := time.NewTimer(time.Second)
	appch := make(chan *struct{}, 1)
	sourcech := make(chan *struct{}, 1)
	go func() {
		var appversion, sourceversion uint32
		ch, cancel, e := Sdk.GetNoticeByProjectID(model.AdminProjectID, model.Group, model.Name)
		if e != nil {
			log.Error(nil, "[config.Init] get notice failed", log.CError(e))
			Close()
			os.Exit(1)
		}
		defer cancel()
		for {
			<-ch
			app, e := Sdk.GetAppConfigByProjectID(model.AdminProjectID, model.Group, model.Name)
			if e != nil {
				if e == ecode.ErrServerClosing {
					return
				}
				log.Error(nil, "[config.Init] get app config failed", log.CError(e))
				continue
			}
			appkey, ok := app.Keys["AppConfig"]
			if !ok {
				log.Error(nil, "[config.Init] key: AppConfig missing")
				continue
			}
			sourcekey, ok := app.Keys["SourceConfig"]
			if !ok {
				log.Error(nil, "[config.Init] key: SourceConfig missing")
				continue
			}
			if appkey.CurVersion == appversion && sourcekey.CurVersion == sourceversion {
				continue
			}
			if appkey.CurValueType != "json" || sourcekey.CurValueType != "json" {
				log.Error(nil, "[config.Init] config data can only support json format")
				continue
			}
			if appkey.CurVersion != appversion {
				var plaintxt []byte
				if secret != "" {
					plaintxt, e = secure.AesDecrypt(secret, appkey.CurValue)
					if e != nil {
						log.Error(nil, "[config.Init] decrypt failed", log.CError(e))
						continue
					}
				} else {
					plaintxt = common.Str2byte(appkey.CurValue)
				}
				c := &AppConfig{}
				if e := json.Unmarshal(plaintxt, c); e != nil {
					log.Error(nil, "[config.Init] key: AppConfig data format wrong", log.CError(e))
					continue
				}
				validateAppConfig(c)
				AC = c
				log.Info(nil, "[config.Init] update app config success", log.Any("config", AC))
				if notice != nil {
					notice(AC)
				}
				appversion = appkey.CurVersion
				select {
				case appch <- nil:
				default:
				}
			}
			if sourcekey.CurVersion != sourceversion && sc == nil {
				//source config can't hot update,can only init once
				var plaintxt []byte
				if secret != "" {
					plaintxt, e = secure.AesDecrypt(secret, sourcekey.CurValue)
					if e != nil {
						log.Error(nil, "[config.Init] decrypt failed", log.CError(e))
						continue
					}
				} else {
					plaintxt = common.Str2byte(sourcekey.CurValue)
				}
				c := &sourceConfig{}
				if e := json.Unmarshal(plaintxt, c); e != nil {
					log.Error(nil, "[config.Init] key: SourceConfig data format wrong", log.CError(e))
					continue
				}
				log.Info(nil, "[config.remote.source] update source config success", log.Any("config", c))
				sc = c
				sourceversion = sourcekey.CurVersion
				initsource()
				select {
				case sourcech <- nil:
				default:
				}
			}
		}
	}()
	for {
		select {
		case <-appch:
		case <-sourcech:
		case <-tmer.C:
			log.Error(nil, "[config.Init] timeout")
			Close()
			os.Exit(1)
		}
		if AC != nil && sc != nil {
			break
		}
	}
}

// Close -
func Close() {
	log.Close()
}

func InitInternal() {
	var secret string
	if str, ok := os.LookupEnv("CONFIG_SECRET"); ok && str != "<CONFIG_SECRET>" && str != "" {
		secret = str
	}
	if len(secret) >= 32 {
		log.Error(nil, "[config.InitInternal] env CONFIG_SECRET length must < 32")
		Close()
		os.Exit(1)
	}
	sctemplate, e := os.ReadFile("./SourceConfig.json")
	if e != nil {
		log.Error(nil, "[config.InitInternal] read ./SourceConfig.json failed", log.CError(e))
		Close()
		os.Exit(1)
	}
	tmpsc := &sourceConfig{}
	if e := json.Unmarshal(sctemplate, tmpsc); e != nil {
		log.Error(nil, "[config.InitInternal] ./SourceConfig.json format wrong", log.CError(e))
		Close()
		os.Exit(1)
	}
	mongoc, ok := tmpsc.Mongo["admin_mongo"]
	if !ok {
		log.Error(nil, "[config.InitInternal] ./SourceConfig.json missing mongo config for 'admin_mongo'")
		Close()
		os.Exit(1)
	}
	mongoc.MongoName = "admin_mongo"
	if len(mongoc.Addrs) == 0 {
		mongoc.Addrs = []string{"127.0.0.1:27017"}
	}
	if mongoc.MaxConnIdletime <= 0 {
		mongoc.MaxConnIdletime = ctime.Duration(time.Minute * 5)
	}
	if mongoc.IOTimeout <= 0 {
		mongoc.IOTimeout = ctime.Duration(time.Millisecond * 500)
	}
	if mongoc.DialTimeout <= 0 {
		mongoc.DialTimeout = ctime.Duration(time.Millisecond * 250)
	}
	var tlsc *tls.Config
	if mongoc.TLS {
		tlsc = &tls.Config{}
		if len(mongoc.SpecificCAPaths) > 0 {
			tlsc.RootCAs = x509.NewCertPool()
			for _, certpath := range mongoc.SpecificCAPaths {
				cert, e := os.ReadFile(certpath)
				if e != nil {
					log.Error(nil, "[config.InitInternal] read specific cert failed",
						log.String("mongo", "admin_mongo"),
						log.String("cert_path", certpath),
						log.CError(e))
					Close()
					os.Exit(1)
				}
				if ok := tlsc.RootCAs.AppendCertsFromPEM(cert); !ok {
					log.Error(nil, "[config.InitInternal] specific cert load failed",
						log.String("mongo", "admin_mongo"),
						log.String("cert_path", certpath),
						log.CError(e))
					Close()
					os.Exit(1)
				}
			}
		}
	}
	db, e := mongo.NewMongo(mongoc.Config, tlsc)
	if e != nil {
		log.Error(nil, "[config.InitInternal] new mongo failed", log.String("mongo", "admin_mongo"), log.CError(e))
		Close()
		os.Exit(1)
	}
	if initinternal.InitDatabase(secret, db.Client) != nil {
		Close()
		os.Exit(1)
	}
	if Sdk, e = initinternal.InitWatch(secret, db.Client); e != nil {
		Close()
		os.Exit(1)
	}
	return
}
func StopInternal() {
	Sdk.Stop()
}
