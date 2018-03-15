package di

import "reflect"

type Lifecycle int

const (
	BeforeInit Lifecycle = iota
	TheInit
	BeforeReady
	TheReady
)

var lifecycleTypes = []reflect.Type{
	BeforeInit:  nil,
	TheInit:     reflect.TypeOf((*Init)(nil)).Elem(),
	BeforeReady: nil,
	TheReady:    reflect.TypeOf((*Ready)(nil)).Elem(),
}

func (l Lifecycle) Type() reflect.Type {
	return lifecycleTypes[l]
}
