package ecode

import (
	"net/http"

	"github.com/chenjie199234/Corelib/cerror"
)

var (
	ErrServerClosing     = cerror.ErrServerClosing     //1000  // http code 449 Warning!! Client will retry on this error,be careful to use this error
	ErrDataConflict      = cerror.ErrDataConflict      //9001  // http code 500
	ErrDataBroken        = cerror.ErrDataBroken        //9002  // http code 500
	ErrDBDataConflict    = cerror.ErrDBDataConflict    //9101  // http code 500
	ErrDBDataBroken      = cerror.ErrDBDataBroken      //9102  // http code 500
	ErrCacheDataConflict = cerror.ErrCacheDataConflict //9201  // http code 500
	ErrCacheDataBroken   = cerror.ErrCacheDataBroken   //9202  // http code 500
	ErrMQDataBroken      = cerror.ErrMQDataBroken      //9301  // http code 500
	ErrUnknown           = cerror.ErrUnknown           //10000 // http code 500
	ErrReq               = cerror.ErrReq               //10001 // http code 400
	ErrResp              = cerror.ErrResp              //10002 // http code 500
	ErrSystem            = cerror.ErrSystem            //10003 // http code 500
	ErrToken             = cerror.ErrToken             //10004 // http code 401
	ErrSession           = cerror.ErrSession           //10005 // http code 401
	ErrAccessKey         = cerror.ErrAccessKey         //10006 // http code 401
	ErrAccessSign        = cerror.ErrAccessSign        //10007 // http code 401
	ErrPermission        = cerror.ErrPermission        //10008 // http code 403
	ErrTooFast           = cerror.ErrTooFast           //10009 // http code 403
	ErrBan               = cerror.ErrBan               //10010 // http code 403
	ErrBusy              = cerror.ErrBusy              //10011 // http code 503
	ErrNotExist          = cerror.ErrNotExist          //10012 // http code 404
	ErrAlreadyExist      = cerror.ErrAlreadyExist      //10013 // http code 400
	ErrPasswordWrong     = cerror.ErrPasswordWrong     //10014 // http code 400
	ErrPasswordLength    = cerror.ErrPasswordLength    //10015 // http code 400

	ErrNotInited           = cerror.MakeCError(20001, http.StatusBadRequest, "not inited")
	ErrAlreadyInited       = cerror.MakeCError(20002, http.StatusBadRequest, "already inited")
	ErrProjectNotExist     = cerror.MakeCError(20003, http.StatusBadRequest, "project doesn't exist")
	ErrProjectAlreadyExist = cerror.MakeCError(20004, http.StatusBadRequest, "project already exist")

	ErrAppNotExist           = cerror.MakeCError(20010, http.StatusBadRequest, "app doesn't exist")
	ErrAppPermissionMissing  = cerror.MakeCError(20011, http.StatusBadRequest, "app's permission node id missing")
	ErrKeyNotExist           = cerror.MakeCError(20012, http.StatusBadRequest, "key doesn't exist")
	ErrAppAlreadyExist       = cerror.MakeCError(20013, http.StatusBadRequest, "app already exist")
	ErrKeyAlreadyExist       = cerror.MakeCError(20014, http.StatusBadRequest, "key already exist")
	ErrIndexNotExist         = cerror.MakeCError(20015, http.StatusBadRequest, "config index doesn't exist")
	ErrProxyPathNotExist     = cerror.MakeCError(20016, http.StatusBadRequest, "proxy path doesn't exist")
	ErrProxyPathAlreadyExist = cerror.MakeCError(20017, http.StatusBadRequest, "proxy path already exist")

	ErrNodeNotExist       = cerror.MakeCError(20100, http.StatusBadRequest, "node not exist")
	ErrPNodeNotExist      = cerror.MakeCError(20101, http.StatusBadRequest, "parent node not exist")
	ErrRoleNotExist       = cerror.MakeCError(20102, http.StatusBadRequest, "role doesn't exist")
	ErrRoleAlreadyExist   = cerror.MakeCError(20103, http.StatusBadRequest, "role already exist")
	ErrUserNotExist       = cerror.MakeCError(20104, http.StatusBadRequest, "user not exist")
	ErrUserAlreadyInvited = cerror.MakeCError(20105, http.StatusBadRequest, "user already invited")
	ErrUserNotInProject   = cerror.MakeCError(20106, http.StatusBadRequest, "user not in project")
)

func ReturnEcode(originerror error, defaulterror *cerror.Error) error {
	if _, ok := originerror.(*cerror.Error); ok {
		return originerror
	}
	return defaulterror
}
