package main

import (
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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/redis/go-redis/v9"
	_ "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	config.InitInternal()
	config.Init(func(ac *config.AppConfig) {
		//this is a notice callback every time appconfig changes
		//this function works in sync mode
		//don't write block logic inside this
		dao.UpdateAppConfig(ac)
		xcrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xgrpc.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateHandlerTimeout(ac.HandlerTimeout)
		xweb.UpdateWebPathRewrite(ac.WebPathRewrite)
		publicmids.UpdateRateConfig(ac.HandlerRate)
		publicmids.UpdateTokenConfig(ac.TokenSecret, ac.SessionTokenExpire.StdDuration())
		publicmids.UpdateSessionConfig(ac.SessionTokenExpire.StdDuration())
		publicmids.UpdateAccessConfig(ac.Accesses)
	})
	defer config.Close()
	if rateredis := config.GetRedis("rate_redis"); rateredis != nil {
		publicmids.UpdateRateRedisInstance(rateredis)
	} else {
		log.Warn(nil, "[main] rate redis missing,all rate check will be failed")
	}
	if sessionredis := config.GetRedis("session_redis"); sessionredis != nil {
		publicmids.UpdateSessionRedisInstance(sessionredis)
	} else {
		log.Warn(nil, "[main] session redis missing,all session event will be failed")
	}
	//start the whole business service
	if e := service.StartService(); e != nil {
		log.Error(nil, "[main] start service failed", log.CError(e))
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
	config.StopInternal()
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
