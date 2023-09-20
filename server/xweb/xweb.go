package xweb

import (
	"crypto/tls"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/service"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/ctime"
	"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/web/mids"
)

var s *web.WebServer

// StartWebServer -
func StartWebServer() {
	c := config.GetWebServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				log.Error(nil, "[xweb] load cert failed:", map[string]interface{}{"cert": cert, "key": key, "error": e})
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	var e error
	if s, e = web.NewWebServer(c.ServerConfig, model.Project, model.Group, model.Name, tlsc); e != nil {
		log.Error(nil, "[xweb] new server failed", map[string]interface{}{"error": e})
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)
	UpdateWebPathRewrite(config.AC.WebPathRewrite)

	r := s.NewRouter()

	//this place can register global midwares
	//r.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusWebServer(r, service.SvcStatus, mids.AllMids())
	api.RegisterAppWebServer(r, service.SvcApp, mids.AllMids())
	api.RegisterUserWebServer(r, service.SvcUser, mids.AllMids())
	api.RegisterPermissionWebServer(r, service.SvcPermission, mids.AllMids())
	api.RegisterInitializeWebServer(r, service.SvcInitialize, mids.AllMids())
	//example
	//api.RegisterExampleWebServer(r, service.SvcExample, mids.AllMids())

	if e = s.StartWebServer(":8000"); e != nil && e != web.ErrServerClosed {
		log.Error(nil, "[xweb] start server failed", map[string]interface{}{"error": e})
		return
	}
	log.Info(nil, "[xweb] server closed", nil)
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	if s != nil {
		s.UpdateHandlerTimeout(timeout)
	}
}

// UpdateWebPathRewrite -
// key origin url,value rewrite url
func UpdateWebPathRewrite(rewrite map[string]map[string]string) {
	if s != nil {
		s.UpdateHandlerRewrite(rewrite)
	}
}

// StopWebServer -
func StopWebServer() {
	if s != nil {
		s.StopWebServer(false)
	}
}
