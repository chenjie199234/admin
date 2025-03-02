package xgrpc

import (
	"crypto/tls"
	"log/slog"
	"sync/atomic"
	"unsafe"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/service"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cgrpc/mids"
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
				slog.ErrorContext(nil, "[xgrpc] load cert failed:", slog.String("cert", cert), slog.String("key", key), slog.String("error", e.Error()))
				return
			}
			certificates = append(certificates, temp)
		}
		tlsc = &tls.Config{Certificates: certificates}
	}
	server, e := cgrpc.NewCGrpcServer(c.ServerConfig, tlsc)
	if e != nil {
		slog.ErrorContext(nil, "[xgrpc] new server failed", slog.String("error", e.Error()))
		return
	}
	//avoid race when build/run in -race mode
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s)), unsafe.Pointer(server))
	UpdateHandlerTimeout(config.AC.HandlerTimeout)

	//this place can register global midwares
	//server.Use(globalmidwares)

	//you just need to register your service here
	api.RegisterStatusCGrpcServer(server, service.SvcStatus, mids.AllMids())
	// api.RegisterAppCGrpcServer(s, service.SvcApp, mids.AllMids())
	// api.RegisterUserCGrpcServer(s, service.SvcUser, mids.AllMids())
	// api.RegisterPermissionCGrpcServer(s, service.SvcPermission, mids.AllMids())
	// api.RegisterInitializeCGrpcServer(s, service.SvcInitialize, mids.AllMids())
	//example
	//api.RegisterExampleCGrpcServer(server, service.SvcExample, mids.AllMids())

	if e = server.StartCGrpcServer(":10000"); e != nil && e != cgrpc.ErrServerClosed {
		slog.ErrorContext(nil, "[xgrpc] start server failed", slog.String("error", e.Error()))
		return
	}
	slog.InfoContext(nil, "[xgrpc] server closed")
}

// UpdateHandlerTimeout -
// first key path,second key method,value timeout duration
func UpdateHandlerTimeout(timeout map[string]map[string]ctime.Duration) {
	//avoid race when build/run in -race mode
	tmps := (*cgrpc.CGrpcServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.UpdateHandlerTimeout(timeout)
	}
}

// StopCGrpcServer force - false(graceful),true(not graceful)
func StopCGrpcServer(force bool) {
	//avoid race when build/run in -race mode
	tmps := (*cgrpc.CGrpcServer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s))))
	if tmps != nil {
		tmps.StopCGrpcServer(force)
	}
}
