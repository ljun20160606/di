package di

import (
	"reflect"
)

type (
	// break if return true
	CupFunc func(name string, cup *Cup) bool
	CupSet  map[*Cup]struct{}

	Cup struct {
		Water      Water
		Class      reflect.Type  // classInfo.kind() Ptr
		Value      reflect.Value // Lookup.kind() Ptr
		Container  Container
		Dependents []*Cup
	}
)

func newCup(water Water, typee reflect.Type, value reflect.Value, container *container) *Cup {
	return &Cup{
		Water:      water,
		Class:      typee,
		Value:      value,
		Container:  container,
		Dependents: []*Cup{},
	}
}

func (c *Cup) injectDependency() {
	classElem := c.Class.Elem()
	valueElem := c.Value.Elem()
	// mustBe(Struct)
	if classElem.Kind() != reflect.Struct {
		return
	}
	for num := valueElem.NumField() - 1; num >= 0; num-- {
		classInfo := classElem.Field(num)
		value := valueElem.Field(num)
		tag, ok := classInfo.Tag.Lookup("di")
		if !ok || tag == "-" {
			continue
		}
		if !value.CanSet() {
			panic(ErrorUnexported.Panic(classInfo.PkgPath, classInfo.Name, classInfo.Type, classInfo.Tag))
		}
		dropInfo := newDropInfo(tag, classInfo)
		if !dropInfo.dependent {
			c.Container.PutIce(&ice{
				classInfo: classInfo,
				value:     value,
				dropInfo:  dropInfo,
				parentCup: c,
			})
			continue
		}
		cup := c.Container.GetCupWithName(dropInfo.name, value.Type())
		if cup == nil {
			otherCup := c.Container.GetCup(value.Type(), dropInfo.name)
			if otherCup == nil {
				panic(ErrorMissing.Panic(value.Type()))
			}
			cup = otherCup
		}
		c.Dependents = append(c.Dependents, cup)
		value.Set(cup.Value)
	}
}

func (c *Cup) init(set CupSet) {
	if _, has := set[c]; has {
		return
	}
	set[c] = struct{}{}
	for _, dependence := range c.Dependents {
		dependence.init(set)
	}
	if init, ok := c.Water.(Init); ok {
		init.Init()
	}
}

func (c *Cup) ready(set CupSet) {
	if _, has := set[c]; has {
		return
	}
	set[c] = struct{}{}
	for _, dependence := range c.Dependents {
		dependence.ready(set)
	}
	if init, ok := c.Water.(Ready); ok {
		init.Ready()
	}
}
