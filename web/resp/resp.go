package resp

import (
	"fmt"
	"net/http"
)

type DataResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Data(data interface{}) (int, DataResp) {
	resp := DataResp{
		Code: 200,
		Msg:  "",
		Data: data,
	}
	return http.StatusOK, resp
}

func UnexpectedError(desc ...interface{}) (int, interface{}) {
	msg := fmt.Sprint(desc...)
	resp := DataResp{
		Code: 500,
		Msg:  msg,
		Data: nil,
	}
	return http.StatusInternalServerError, resp
}

func UnexpectedErrorf(format string, a ...interface{}) (int, interface{}) {
	msg := fmt.Sprintf(format, a...)
	resp := DataResp{
		Code: 500,
		Msg:  msg,
		Data: nil,
	}
	return http.StatusInternalServerError, resp
}
