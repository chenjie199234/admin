package main

import (
	_ "embed"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/server/xcrpc"
	"github.com/chenjie199234/admin/server/xgrpc"
	"github.com/chenjie199234/admin/server/xweb"
	"github.com/chenjie199234/admin/service"

	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	_ "github.com/chenjie199234/Corelib/monitor"
)

// this is used for DirectSDK
//
//go:embed AppConfig.json
var AppConfigTemplate []byte

// this is used for DirectSDK
//
//go:embed SourceConfig.json
var SourceConfigTemplate []byte

func main() {
	config.Init(func(ac *config.AppConfig) {
		//this is a notice callback every time appconfig changes
		//this function works in sync mode
		//don't write block logic inside this
		dao.UpdateAPI(ac)
		xcrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xgrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateWebPathRewrite(ac.WebPathRewrite)
		publicmids.UpdateRateConfig(ac.HandlerRate)
		publicmids.UpdateTokenConfig(ac.TokenSecret, ac.SessionTokenExpire.StdDuration())
		publicmids.UpdateSessionConfig(ac.SessionTokenExpire.StdDuration())
		publicmids.UpdateAccessConfig(ac.Accesses)
	}, AppConfigTemplate, SourceConfigTemplate)
	defer config.Close()
	if rateredis := config.GetRedis("rate_redis"); rateredis != nil {
		publicmids.UpdateRateRedisInstance(rateredis)
	} else {
		log.Warning(nil, "[main] rate redis missing,all rate check will be failed", nil)
	}
	if sessionredis := config.GetRedis("session_redis"); sessionredis != nil {
		publicmids.UpdateSessionRedisInstance(sessionredis)
	} else {
		log.Warning(nil, "[main] session redis missing,all session event will be failed", nil)
	}
	//start the whole business service
	if e := service.StartService(); e != nil {
		log.Error(nil, "[main] start service failed", map[string]interface{}{"error": e})
		return
	}
	//start low level net service
	ch := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		xcrpc.StartCrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StartWebServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StartCGrpcServer()
		select {
		case ch <- syscall.SIGTERM:
		default:
		}
		wg.Done()
	}()
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	//stop the whole business service
	service.StopService()
	//stop low level net service
	wg.Add(1)
	go func() {
		xcrpc.StopCrpcServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xweb.StopWebServer()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		xgrpc.StopCGrpcServer()
		wg.Done()
	}()
	wg.Wait()
}
