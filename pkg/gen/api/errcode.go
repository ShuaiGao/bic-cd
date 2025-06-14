package api

import (
	"errors"
)

type ErrCode interface {
	Code() int
	String() string
}

type Err struct {
	code   int
	msg    string
	detail error
}

func (e Err) Code() int {
	return e.code
}

func (e Err) String() string {
	//if !gin.IsDebugging() {
	//	return e.msg
	//}
	if e.detail != nil {
		return e.msg + ", detail: " + e.detail.Error()
	}
	return e.msg
}

func (e Err) Wrap(err interface{}) ErrCode {
	switch err.(type) {
	case error:
		e.detail = err.(error)
	case string:
		e.detail = errors.New(err.(string))
	}
	return e
}

func NewErr(e Err, err interface{}) Err {
	switch err.(type) {
	case error:
		e.detail = err.(error)
	case string:
		e.detail = errors.New(err.(string))
	}
	return e
}

func Equal(a, b ErrCode) bool {
	return a.Code() == b.Code()
}

var (
	ECSuccess               = Err{code: 0, msg: "成功"}
	ECAuthError             = Err{code: 100001, msg: "授权登录失败"}
	ECParamError            = Err{code: 100002, msg: "参数错误"}
	ECServerError           = Err{code: 100003, msg: "服务端错误"}
	ECNoSupport             = Err{code: 100004, msg: "暂不支持该功能"}
	ECCannotOpSelf          = Err{code: 100005, msg: "不能给自己分配权限"}
	ECDbError               = Err{code: 102000, msg: "数据库查询错误"}
	ECDbNotfound            = Err{code: 102001, msg: "id对应数据不存在"}
	ECDbCreateError         = Err{code: 102002, msg: "数据库创建数据错误"}
	ECDbFindError           = Err{code: 102003, msg: "数据库查找出错"}
	ECDbDeleteError         = Err{code: 102004, msg: "删除id对应数据不存在"}
	ECDbTransactionError    = Err{code: 102005, msg: "数据库提交事务出错"}
	ECDbNotfound_permission = Err{code: 102006, msg: "权限项ID不存在"}
	ECDbSaveError           = Err{code: 102007, msg: "数据库保存数据错误"}
	ECDbCountError          = Err{code: 102008, msg: "数据库count错误"}
	ECDbExecError           = Err{code: 102009, msg: "sql执行错误"}
	ECParamValidatorError   = Err{code: 102200, msg: "参数校验错误"}
	ECRepeatedVersion       = Err{code: 102201, msg: "该服务版本已存在"}
)
