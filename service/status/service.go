package status

import (
	"context"
	// "log/slog"
	"time"

	//"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/api"
	statusdao "github.com/chenjie199234/admin/dao/status"
	//"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/cotel"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/host"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for status business
type Service struct {
	stop *graceful.Graceful

	statusDao *statusdao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		//statusDao: statusdao.NewDao(config.GetSql("status_sql"), config.GetRedis("status_redis"), config.GetMongo("status_mongo")),
		statusDao: statusdao.NewDao(nil, nil, nil),
	}, nil
}

// Ping -
func (s *Service) Ping(ctx context.Context, in *api.Pingreq) (*api.Pingresp, error) {
	//if _, ok := ctx.(*crpc.NoStreamServerContext); ok {
	//        slog.InfoContext("this is a crpc call")
	//}
	//if _, ok := ctx.(*cgrpc.NoStreamServerContext); ok {
	//        slog.InfoContext("this is a cgrpc call")
	//}
	//if _, ok := ctx.(*web.Context); ok {
	//        Slog.InfoContext("this is a web call")
	//}
	totalmem, lastmem, maxmem := cotel.GetMEM()
	lastcpu, maxcpu, avgcpu := cotel.GetCPU()
	return &api.Pingresp{
		ClientTimestamp: in.Timestamp,
		ServerTimestamp: time.Now().UnixNano(),
		TotalMem:        totalmem,
		CurMemUsage:     lastmem,
		MaxMemUsage:     maxmem,
		CpuNum:          cotel.CPUNum,
		CurCpuUsage:     lastcpu,
		AvgCpuUsage:     avgcpu,
		MaxCpuUsage:     maxcpu,
		Host:            host.Hostname,
		Ip:              host.Hostip,
	}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
