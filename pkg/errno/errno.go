package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode                 = 0
	ServiceErrCode              = 10001
	ParamErrCode                = 10002
	UserAlreadyExistErrCode     = 10003
	UserNotExistErrCode         = 10004
	AuthorizationFailedErrCode  = 10005
	OSSUploadFailedErrCode      = 10006
	DBOperationFailedErrCode    = 10007
	CommentExistErrCode         = 10008
	CommentNotExistErrCode      = 10009
	FavoriteExistErrCode        = 10010
	FavoriteNotExistErrCode     = 10011
	RedisOperationFailedErrCode = 10012
)

type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success                    = NewErrNo(SuccessCode, "Success")
	ServiceErr                 = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	ParamErr                   = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	UserAlreadyExistErr        = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	UserNotExistErr            = NewErrNo(UserNotExistErrCode, "User not exists")
	AuthorizationFailedErr     = NewErrNo(AuthorizationFailedErrCode, "Authorization failed")
	OSSUploadFailedErr         = NewErrNo(OSSUploadFailedErrCode, "Upload file to oss failed")
	DatabaseOperationFailedErr = NewErrNo(DBOperationFailedErrCode, "Database operation error")
	CommentExistErr            = NewErrNo(CommentExistErrCode, "Comment already exists")
	CommentNotExistErr         = NewErrNo(CommentNotExistErrCode, "Comment not exists")
	FavoriteExistErr           = NewErrNo(FavoriteExistErrCode, "Favorite already exists")
	FavoriteNotExistErr        = NewErrNo(FavoriteNotExistErrCode, "Favorite not exists")
)

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
