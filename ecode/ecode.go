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
	ErrAuth       = cerror.ErrAuth       //10004 // http code 401
	ErrPermission = cerror.ErrPermission //10005 // http code 403
	ErrTooFast    = cerror.ErrTooFast    //10006 // http code 403
	ErrBan        = cerror.ErrBan        //10007 // http code 403
	ErrBusy       = cerror.ErrBusy       //10008 // http code 503
	ErrNotExist   = cerror.ErrNotExist   //10009 // http code 404

	ErrNotInited        = cerror.MakeError(20001, http.StatusBadRequest, "not inited")
	ErrAlreadyInited    = cerror.MakeError(20002, http.StatusBadRequest, "already inited")
	ErrAppNotExist      = cerror.MakeError(20003, http.StatusBadRequest, "app doesn't exist")
	ErrKeyNotExist      = cerror.MakeError(20004, http.StatusBadRequest, "key doesn't exist")
	ErrAppAlreadyExist  = cerror.MakeError(20005, http.StatusBadRequest, "app already exist")
	ErrKeyAlreadyExist  = cerror.MakeError(20006, http.StatusBadRequest, "key already exist")
	ErrIndexNotExist    = cerror.MakeError(20007, http.StatusBadRequest, "config index doesn't exist")
	ErrWrongCipher      = cerror.MakeError(20008, http.StatusBadRequest, "wrong cipher")
	ErrCipherLength     = cerror.MakeError(20009, http.StatusBadRequest, "cipher must be empty or 32 byte length")
	ErrRoleNotExist     = cerror.MakeError(20010, http.StatusBadRequest, "role doesn't exist")
	ErrRoleAlreadyExist = cerror.MakeError(20011, http.StatusBadRequest, "role already exist")
	ErrUserNotExist     = cerror.MakeError(20012, http.StatusBadRequest, "user not exist")
	ErrNodeNotExist     = cerror.MakeError(20013, http.StatusBadRequest, "node not exist")
	ErrPNodeNotExist    = cerror.MakeError(20014, http.StatusBadRequest, "parent node not exist")

	ErrPageOverflow = cerror.MakeError(30001, http.StatusBadRequest, "page overflow")
)

func ReturnEcode(originerror error, defaulterror *cerror.Error) error {
	if _, ok := originerror.(*cerror.Error); ok {
		return originerror
	}
	return defaulterror
}
