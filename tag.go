package di

import (
	"github.com/ljun20160606/go-lib/reflectl"
	"reflect"
	"strings"
)

type dropInfo struct {
	auto      bool
	dependent bool
	name      string
	prefix    string
	path      string
}

// 正则表达式说明
// 1. di:"[^\.]*" 由di控制依赖注入
// 2. di:"[^\.]*\..*" 由di plugin控制依赖注入
func newDropInfo(tagDi string, class reflect.StructField) *dropInfo {
	si := &dropInfo{}
	dotIndex := strings.Index(tagDi, ".")
	switch {
	case dotIndex == 0:
		panic(ErrorTagDotIndex.Panic(class.Name, class.Type, class.Tag, ". 不能放在首位"))
	case dotIndex == len(tagDi)-1:
		panic(ErrorTagDotIndex.Panic(class.Name, class.Type, class.Tag, ". 不能放在末尾"))
	case dotIndex > 0: // 2.
		si.prefix = tagDi[:dotIndex]
		si.path = tagDi[dotIndex+1:]
		si.name = reflectl.GetTypeDefaultName(class.Type)
		return si
	}
	// 1.
	if tagDi == "*" {
		si.auto = true
		si.dependent = true
		si.name = reflectl.GetTypeDefaultName(class.Type)
		return si
	}
	si.dependent = true
	si.name = tagDi
	return si
}
