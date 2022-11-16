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

	ErrNotInited        = cerror.MakeError(20001, http.StatusBadRequest, "not inited")
	ErrAlreadyInited    = cerror.MakeError(20002, http.StatusBadRequest, "already inited")
	ErrPasswordLength   = cerror.MakeError(20003, http.StatusBadRequest, "password length must less then 32")
	ErrPasswordWrong    = cerror.MakeError(20004, http.StatusBadRequest, "password wrong")
	ErrOldPasswordWrong = cerror.MakeError(20005, http.StatusBadRequest, "old password wrong")

	ErrAppNotExist      = cerror.MakeError(20010, http.StatusBadRequest, "app doesn't exist")
	ErrKeyNotExist      = cerror.MakeError(20011, http.StatusBadRequest, "key doesn't exist")
	ErrAppAlreadyExist  = cerror.MakeError(20012, http.StatusBadRequest, "app already exist")
	ErrKeyAlreadyExist  = cerror.MakeError(20013, http.StatusBadRequest, "key already exist")
	ErrIndexNotExist    = cerror.MakeError(20014, http.StatusBadRequest, "config index doesn't exist")
	ErrWrongSecret      = cerror.MakeError(20015, http.StatusBadRequest, "wrong secret")
	ErrSecretLength     = cerror.MakeError(20016, http.StatusBadRequest, "secret length must less then 32")
	ErrConfigDataBroken = cerror.MakeError(20017, http.StatusBadRequest, "config data broken")

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
