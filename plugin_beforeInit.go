package di

import (
	"github.com/ljun20160606/gox/reflectx"
	"reflect"
)

func init() {
	RegisterPlugin(BeforeInit, new(BeforeInitPlugin))
}

var (
	beforeInitType = reflect.TypeOf((*BeforeInitType)(nil)).Elem()
)

type (
	BeforeInitType interface {
		BeforeInit() interface{}
	}
	BeforeInitPlugin struct {
	}
)

func (b *BeforeInitPlugin) Lookup(path string, ice Ice) (v interface{}) {
	if path == "*" {
		ice.Container().EachCup(func(name string, cup *Cup) bool {
			if beforeInitType, ok := cup.Water.(BeforeInitType); ok {
				vv := beforeInitType.BeforeInit()
				if reflectx.TypeEqual(ice.Type(), reflect.TypeOf(vv)) {
					v = vv
					return true
				}
			}
			return false
		})
		if v == nil {
			panic(ErrorMissing.Panic(b.Prefix() + "." + path))
		}
		return v
	}
	cup := ice.Container().GetCupWithName(path, beforeInitType)
	if cup == nil {
		panic(ErrorMissing.Panic(b.Prefix() + "." + path))
	}
	return cup.Water.(BeforeInitType).BeforeInit()
}

func (b *BeforeInitPlugin) Prefix() string {
	return "@"
}

func (b *BeforeInitPlugin) Priority() int {
	return 1
}
