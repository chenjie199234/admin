package xcrpc

import (
	"crypto/tls"
	"log/slog"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/service"

	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/crpc/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

var s *crpc.CrpcServer

// StartCrpcServer -
func StartCrpcServer() {
	c := config.GetCrpcServerConfig()
	var tlsc *tls.Config
	if len(c.Certs) > 0 {
		certificates := make([]tls.Certificate, 0, len(c.Certs))
		for cert, key := range c.Certs {
			temp, e := tls.LoadX509KeyPair(cert, key)
			if e != nil {
				slog.ErrorContext(nil, "[xcrpc] load cert failed", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	var e error
	if s, e = crpc.NewCrpcServer(c.ServerConfig, model.Project, model.Group, model.Name, tlsc); e != nil {
		slog.ErrorContext(nil, "[xcrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//s.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCrpcServer(s, service.SvcStatus, mids.AllMids())
	// api.RegisterAppCrpcServer(s, service.SvcApp, mids.AllMids())
	// api.RegisterUserCrpcServer(s, service.SvcUser, mids.AllMids())
	// api.RegisterPermissionCrpcServer(s, service.SvcPermission, mids.AllMids())
	// api.RegisterInitializeCrpcServer(s, service.SvcInitialize, mids.AllMids())
	//example
	//api.RegisterExampleCrpcServer(s, service.SvcExample,mids.AllMids())

	if e = s.StartCrpcServer(":9000"); e != nil && e != crpc.ErrServerClosed {
		slog.ErrorContext(nil, "[xcrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xcrpc] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	if s != nil {
		s.UpdateHandlerTimeout(timeout)
	}
}

// StopCrpcServer -
func StopCrpcServer() {
	if s != nil {
		s.StopCrpcServer(false)
	}
}
