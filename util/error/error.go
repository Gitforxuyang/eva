package error

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UNKNOW_ERROR = 1001
)

type EvaError struct {
	AppId   string     //错误发生的服务
	Code    int32      //错误码 业务的code
	Message string     //错误消息
	Detail  string     //更详细的错误消息 不对外展示的
	Status  codes.Code //grpc的错误码
}

func (e *EvaError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(appId, message, detail string, code int32, status codes.Code) error {
	return &EvaError{
		AppId:   appId,
		Code:    code,
		Detail:  detail,
		Message: message,
		Status:  status,
	}
}

func Parse(err string) *EvaError {
	e := new(EvaError)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
		e.Code = UNKNOW_ERROR
		e.Message = "未知错误"
		e.Status = codes.Unknown
	}
	return e
}

func FromError(err error) *EvaError {
	if verr, ok := err.(*EvaError); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

func EncodeStatus(e *EvaError) *status.Status {
	status := status.New(e.Status, e.Error())
	return status
}

func DecodeStatus(e error) *EvaError {
	status, ok := status.FromError(e)
	if !ok {
		return Parse(e.Error())
	} else {
		return Parse(status.Message())
	}
}