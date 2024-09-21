package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/initinternal"
	"github.com/chenjie199234/admin/model"

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
		slog.WarnContext(nil, "[config.Init] missing env RUN_ENV")
	}
	if str, ok := os.LookupEnv("DEPLOY_ENV"); ok && str != "<DEPLOY_ENV>" && str != "" {
		EC.DeployEnv = &str
	} else {
		slog.WarnContext(nil, "[config.Init] missing env DEPLOY_ENV")
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
			slog.ErrorContext(nil, "[config.Init] get notice failed", slog.String("error", e.Error()))
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
				slog.ErrorContext(nil, "[config.Init] get app config failed", slog.String("error", e.Error()))
				continue
			}
			appkey, ok := app.Keys["AppConfig"]
			if !ok {
				slog.ErrorContext(nil, "[config.Init] key: AppConfig missing")
				continue
			}
			sourcekey, ok := app.Keys["SourceConfig"]
			if !ok {
				slog.ErrorContext(nil, "[config.Init] key: SourceConfig missing")
				continue
			}
			if appkey.CurVersion == appversion && sourcekey.CurVersion == sourceversion {
				continue
			}
			if appkey.CurValueType != "json" || sourcekey.CurValueType != "json" {
				slog.ErrorContext(nil, "[config.Init] config data can only support json format")
				continue
			}
			if appkey.CurVersion != appversion {
				var plaintxt []byte
				if secret != "" {
					plaintxt, e = secure.AesDecrypt(secret, appkey.CurValue)
					if e != nil {
						slog.ErrorContext(nil, "[config.Init] decrypt failed", slog.String("error", e.Error()))
						continue
					}
				} else {
					plaintxt = common.STB(appkey.CurValue)
				}
				c := &AppConfig{}
				if e := json.Unmarshal(plaintxt, c); e != nil {
					slog.ErrorContext(nil, "[config.Init] key: AppConfig data format wrong", slog.String("error", e.Error()))
					continue
				}
				validateAppConfig(c)
				AC = c
				slog.InfoContext(nil, "[config.Init] update app config success", slog.Any("config", AC))
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
						slog.ErrorContext(nil, "[config.Init] decrypt failed", slog.String("error", e.Error()))
						continue
					}
				} else {
					plaintxt = common.STB(sourcekey.CurValue)
				}
				c := &sourceConfig{}
				if e := json.Unmarshal(plaintxt, c); e != nil {
					slog.ErrorContext(nil, "[config.Init] key: SourceConfig data format wrong", slog.String("error", e.Error()))
					continue
				}
				slog.InfoContext(nil, "[config.remote.source] update source config success", slog.Any("config", c))
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
			slog.ErrorContext(nil, "[config.Init] timeout")
			os.Exit(1)
		}
		if AC != nil && sc != nil {
			break
		}
	}
}

func InitInternal() {
	var secret string
	if str, ok := os.LookupEnv("CONFIG_SECRET"); ok && str != "<CONFIG_SECRET>" && str != "" {
		secret = str
	}
	if len(secret) >= 32 {
		slog.ErrorContext(nil, "[config.InitInternal] env CONFIG_SECRET length must < 32")
		os.Exit(1)
	}
	sctemplate, e := os.ReadFile("./SourceConfig.json")
	if e != nil {
		slog.ErrorContext(nil, "[config.InitInternal] read ./SourceConfig.json failed", slog.String("error", e.Error()))
		os.Exit(1)
	}
	tmpsc := &sourceConfig{}
	if e := json.Unmarshal(sctemplate, tmpsc); e != nil {
		slog.ErrorContext(nil, "[config.InitInternal] ./SourceConfig.json format wrong", slog.String("error", e.Error()))
		os.Exit(1)
	}
	mongoc, ok := tmpsc.Mongo["admin_mongo"]
	if !ok {
		slog.ErrorContext(nil, "[config.InitInternal] ./SourceConfig.json missing mongo config for 'admin_mongo'")
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
					slog.ErrorContext(nil, "[config.InitInternal] read specific cert failed",
						slog.String("mongo", "admin_mongo"),
						slog.String("cert_path", certpath),
						slog.String("error", e.Error()))
					os.Exit(1)
				}
				if ok := tlsc.RootCAs.AppendCertsFromPEM(cert); !ok {
					slog.ErrorContext(nil, "[config.InitInternal] specific cert load failed",
						slog.String("mongo", "admin_mongo"),
						slog.String("cert_path", certpath),
						slog.String("error", e.Error()))
					os.Exit(1)
				}
			}
		}
	}
	db, e := mongo.NewMongo(mongoc.Config, tlsc)
	if e != nil {
		slog.ErrorContext(nil, "[config.InitInternal] new mongo failed", slog.String("mongo", "admin_mongo"), slog.String("error", e.Error()))
		os.Exit(1)
	}
	if initinternal.InitDatabase(secret, db.Client) != nil {
		os.Exit(1)
	}
	if Sdk, e = initinternal.InitWatch(secret, db.Client); e != nil {
		os.Exit(1)
	}
	return
}
func StopInternal() {
	Sdk.Stop()
}
