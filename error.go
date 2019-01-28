package di

import "fmt"

type errorDi int

const (
	ErrorType errorDi = iota
	ErrorMissing
	ErrorUnexported
	ErrorTagDotIndex
	ErrorStopIterator
	ErrorUnserialize
)

var errorString = [...]string{
	ErrorType:         "请传入Ptr \n当前basic type %v",
	ErrorMissing:      "找不到 %v",
	ErrorUnexported:   "%v: %v %v %v 需要大写变量首字母 mustBeExported",
	ErrorTagDotIndex:  "错误field: %v %v `%v`  %v",
	ErrorStopIterator: "停止循环",
	ErrorUnserialize:  "反序列化失败，类型可能有误",
}

func (e errorDi) String() string {
	return errorString[e]
}

func (e errorDi) Error() string {
	return errorString[e]
}

func (e errorDi) Panic(argv ...interface{}) string {
	return fmt.Sprintf(e.String(), argv...)
}
