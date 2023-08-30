package ecode

import (
	"net/http"

	"github.com/chenjie199234/Corelib/cerror"
)

var (
	ErrUnknown    = cerror.ErrUnknown    //10000 // http code 500
	ErrReq        = cerror.ErrReq        //10001 // http code 400
	ErrResp       = cerror.ErrResp       //10002 // http code 500
	ErrSystem     = cerror.ErrSystem     //10003 // http code 500
	ErrToken      = cerror.ErrToken      //10004 // http code 401
	ErrSession    = cerror.ErrSession    //10005 // http code 401
	ErrKey        = cerror.ErrKey        //10006 // http code 401
	ErrSign       = cerror.ErrSign       //10007 // http code 401
	ErrPermission = cerror.ErrPermission //10008 // http code 403
	ErrTooFast    = cerror.ErrTooFast    //10009 // http code 403
	ErrBan        = cerror.ErrBan        //10010 // http code 403
	ErrBusy       = cerror.ErrBusy       //10011 // http code 503
	ErrNotExist   = cerror.ErrNotExist   //10012 // http code 404

	// ErrInitConflict        = cerror.MakeError(20000, http.StatusBadRequest, "init conflict")
	ErrNotInited           = cerror.MakeError(20001, http.StatusBadRequest, "not inited")
	ErrAlreadyInited       = cerror.MakeError(20002, http.StatusBadRequest, "already inited")
	ErrPasswordLength      = cerror.MakeError(20003, http.StatusBadRequest, "password length must less then 32")
	ErrPasswordWrong       = cerror.MakeError(20004, http.StatusBadRequest, "password wrong")
	ErrProjectNotExist     = cerror.MakeError(20005, http.StatusBadRequest, "project doesn't exist")
	ErrProjectAlreadyExist = cerror.MakeError(20006, http.StatusBadRequest, "project already exist")

	ErrAppNotExist           = cerror.MakeError(20010, http.StatusBadRequest, "app doesn't exist")
	ErrAppPermissionMissing  = cerror.MakeError(20011, http.StatusBadRequest, "app's permission node id missing")
	ErrKeyNotExist           = cerror.MakeError(20012, http.StatusBadRequest, "key doesn't exist")
	ErrAppAlreadyExist       = cerror.MakeError(20013, http.StatusBadRequest, "app already exist")
	ErrKeyAlreadyExist       = cerror.MakeError(20014, http.StatusBadRequest, "key already exist")
	ErrIndexNotExist         = cerror.MakeError(20015, http.StatusBadRequest, "config index doesn't exist")
	ErrWrongSecret           = cerror.MakeError(20016, http.StatusBadRequest, "wrong secret")
	ErrSecretLength          = cerror.MakeError(20017, http.StatusBadRequest, "secret length must less then 32")
	ErrDataBroken            = cerror.MakeError(20018, http.StatusBadRequest, "data broken")
	ErrSignCheckFailed       = cerror.MakeError(20019, http.StatusBadRequest, "sign check failed")
	ErrProxyPathNotExist     = cerror.MakeError(20020, http.StatusBadRequest, "proxy path doesn't exist")
	ErrProxyPathAlreadyExist = cerror.MakeError(20021, http.StatusBadRequest, "proxy path already exist")

	ErrNodeNotExist       = cerror.MakeError(20100, http.StatusBadRequest, "node not exist")
	ErrPNodeNotExist      = cerror.MakeError(20101, http.StatusBadRequest, "parent node not exist")
	ErrRoleNotExist       = cerror.MakeError(20102, http.StatusBadRequest, "role doesn't exist")
	ErrRoleAlreadyExist   = cerror.MakeError(20103, http.StatusBadRequest, "role already exist")
	ErrUserNotExist       = cerror.MakeError(20104, http.StatusBadRequest, "user not exist")
	ErrUserAlreadyInvited = cerror.MakeError(20105, http.StatusBadRequest, "user already invited")
	ErrUserNotInProject   = cerror.MakeError(20106, http.StatusBadRequest, "user not in project")

	ErrPageOverflow = cerror.MakeError(30001, http.StatusBadRequest, "page overflow")
)

func ReturnEcode(originerror error, defaulterror *cerror.Error) error {
	if _, ok := originerror.(*cerror.Error); ok {
		return originerror
	}
	return defaulterror
}
