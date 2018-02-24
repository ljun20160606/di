package di

import (
	"reflect"
)

type (
	Ice interface {
		// 获取 容器
		Container() Container
		// 获取 父级cup
		ParentCup() *Cup
		Value() reflect.Value
		Type() reflect.Type
		LoadPlugin(p Plugin)
		Prefix() string
	}

	ice struct {
		classInfo reflect.StructField
		parentCup *Cup
		value     reflect.Value
		dropInfo  *dropInfo
	}
)

func (s *ice) LoadPlugin(p Plugin) {
	s.value.Set(reflect.ValueOf(p.Lookup(s.dropInfo.path, s)))
}

func (s *ice) Container() Container {
	return s.parentCup.Container
}

func (s *ice) ParentCup() *Cup {
	return s.parentCup
}

func (s *ice) Value() reflect.Value {
	return s.value
}

func (s *ice) Type() reflect.Type {
	return s.value.Type()
}

func (s *ice) Prefix() string {
	return s.dropInfo.prefix
}
