package service

import (
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/service/config"
	"github.com/chenjie199234/admin/service/initialize"
	"github.com/chenjie199234/admin/service/permission"
	"github.com/chenjie199234/admin/service/status"
	"github.com/chenjie199234/admin/service/user"
)

// SvcStatus one specify sub service
var SvcStatus *status.Service
var SvcConfig *config.Service
var SvcUser *user.Service
var SvcPermission *permission.Service
var SvcInitialize *initialize.Service

// StartService start the whole service
func StartService() error {
	if e := dao.NewApi(); e != nil {
		return e
	}
	//start sub service
	SvcStatus = status.Start()
	SvcConfig = config.Start()
	SvcUser = user.Start()
	SvcPermission = permission.Start()
	SvcInitialize = initialize.Start()
	return nil
}

// StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
	SvcConfig.Stop()
	SvcUser.Stop()
	SvcPermission.Stop()
	SvcInitialize.Stop()
}
