package service

import (
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/service/app"
	"github.com/chenjie199234/admin/service/initialize"
	"github.com/chenjie199234/admin/service/permission"
	"github.com/chenjie199234/admin/service/raw"
	"github.com/chenjie199234/admin/service/status"
	"github.com/chenjie199234/admin/service/user"
)

// SvcStatus one specify sub service
var SvcStatus *status.Service

// SvcRaw one specify sub service
var SvcRaw *raw.Service

var SvcApp *app.Service
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
	SvcRaw = raw.Start()
	SvcApp = app.Start()
	SvcUser = user.Start()
	SvcPermission = permission.Start()
	SvcInitialize = initialize.Start()
	return nil
}

// StopService stop the whole service
func StopService() {
	//stop sub service
	SvcStatus.Stop()
	SvcRaw.Stop()
	SvcApp.Stop()
	SvcUser.Stop()
	SvcPermission.Stop()
	SvcInitialize.Stop()
}
