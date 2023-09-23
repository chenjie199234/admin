package xgrpc

import (
	"crypto/tls"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/service"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cgrpc/mids"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s *cgrpc.CGrpcServer

// StartCGrpcServer -
func StartCGrpcServer() {
	c := config.GetCGrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				log.Error(nil, "[xgrpc] load cert failed", log.String("cert", cert), log.String("key", key), log.CError(e))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	var e error
	if s, e = cgrpc.NewCGrpcServer(c.ServerConfig, model.Project, model.Group, model.Name, tlsc); e != nil {
		log.Error(nil, "[xgrpc] new server failed", log.CError(e))
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCGrpcServer(s, service.SvcStatus, mids.AllMids())
	api.RegisterAppCGrpcServer(s, service.SvcApp, mids.AllMids())
	api.RegisterUserCGrpcServer(s, service.SvcUser, mids.AllMids())
	api.RegisterPermissionCGrpcServer(s, service.SvcPermission, mids.AllMids())
	api.RegisterInitializeCGrpcServer(s, service.SvcInitialize, mids.AllMids())
	//example
	//api.RegisterExampleCGrpcServer(s, service.SvcExample, mids.AllMids())

	if e = s.StartCGrpcServer(":10000"); e != nil && e != cgrpc.ErrServerClosed {
		log.Error(nil, "[xgrpc] start server failed", log.CError(e))
		return
	}
	log.Info(nil, "[xgrpc] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	if s != nil {
		s.UpdateHandlerTimeout(timeout)
	}
}

// StopCGrpcServer -
func StopCGrpcServer() {
	if s != nil {
		s.StopCGrpcServer(false)
	}
}
