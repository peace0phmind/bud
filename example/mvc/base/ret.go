package base

import (
	"fmt"
	"gorm.io/gorm"
)

type RetCode int

const (
	RetCodeRetOk              RetCode = 0 // OK
	RetCodeParamsError        RetCode = 1
	RetCodeRecordNotExists    RetCode = 2
	RetCodeSysError           RetCode = 3
	RetCodeAdminInitedError   RetCode = 4
	RetCodeMethodNotSupported RetCode = 5

	RetCodeCreateErr RetCode = 10
	RetCodeUpdateErr RetCode = 11
	RetCodeDeleteErr RetCode = 12
	RetCodeImportErr RetCode = 13
	RetCodeExportErr RetCode = 14

	RetCodeQueryErr RetCode = 20

	RetCodeUnauthorized RetCode = 401
	RetCodeForbidden    RetCode = 403
)

var (
	RetCodeMessage = map[RetCode]string{
		RetCodeRetOk:              "Success",
		RetCodeParamsError:        "Params error: %s",
		RetCodeRecordNotExists:    "Record not exists.",
		RetCodeSysError:           "System error: %s",
		RetCodeAdminInitedError:   "System init error: admin is %s",
		RetCodeMethodNotSupported: "Method %s not supported",

		RetCodeCreateErr: "Create err: %s",
		RetCodeUpdateErr: "Update err: %s",
		RetCodeDeleteErr: "Delete err: %s",
		RetCodeImportErr: "Import err: %s",
		RetCodeExportErr: "Export err: %s",

		RetCodeQueryErr: "Query err: %s",

		RetCodeUnauthorized: "Login error, username or password error.",
		RetCodeForbidden:    "forbidden: %s",
	}

	//Ret_Code_Message_CN = map[RetCode]string{
	//	Ret_Code_RET_OK:      "成功。",
	//	Ret_Code_Unauthorized: "登录失败，用户名或密码错误。",
	//}
)

type Ret[T any] struct {
	Code    RetCode `json:"code"`
	Message string  `json:"message"`
	Data    T       `json:"data,omitempty"`
}

type PageRet[T any] struct {
	Items []T         `json:"items"`
	Page  *PageParams `json:"page"`
}

func NewRet[T any](code RetCode, data T, messageParams ...interface{}) *Ret[T] {

	if msg, ok := RetCodeMessage[code]; !ok {
		panic(fmt.Sprintf("Message code not exists: %d", code))
	} else {
		if len(messageParams) > 0 {
			msg = fmt.Sprintf(msg, messageParams)
		}

		return &Ret[T]{
			Code:    code,
			Message: msg,
			Data:    data,
		}
	}
}

func RetErr(code RetCode, messageParams ...interface{}) *Ret[any] {
	if msg, ok := RetCodeMessage[code]; !ok {
		panic(fmt.Sprintf("Message code not exists: %d", code))
	} else {
		if len(messageParams) > 0 {
			msg = fmt.Sprintf(msg, messageParams)
		}

		return &Ret[any]{
			Code:    code,
			Message: msg,
			Data:    nil,
		}
	}
}

func RetOk[T any](data T) *Ret[T] {
	return NewRet(RetCodeRetOk, data)
}

func RetOkPage[T any](data []T, page *PageParams) *Ret[*PageRet[T]] {
	return RetOk(&PageRet[T]{
		Items: data,
		Page:  page,
	})
}

type QueryCondition interface {
	QueryCondition(db *gorm.DB) *gorm.DB
}
