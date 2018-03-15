package di

import (
	"reflect"
)

type (
	// break if return true
	CupFunc func(name string, cup *Cup) bool
	CupSet  map[*Cup]struct{}

	Cup struct {
		Water        Water
		Class        reflect.Type  // classInfo.kind() Ptr
		Value        reflect.Value // Lookup.kind() Ptr
		Container    Container
		Dependencies []*Cup
	}
)

func newCup(water Water, typee reflect.Type, value reflect.Value, container *container) *Cup {
	return &Cup{
		Water:        water,
		Class:        typee,
		Value:        value,
		Container:    container,
		Dependencies: []*Cup{},
	}
}

func (c *Cup) injectDependency() {
	classElem := c.Class.Elem()
	valueElem := c.Value.Elem()
	// mustBe(struct)
	if classElem.Kind() != reflect.Struct {
		return
	}
	for num := valueElem.NumField() - 1; num >= 0; num-- {
		classInfo := classElem.Field(num)
		value := valueElem.Field(num)
		tag, ok := classInfo.Tag.Lookup(DI)
		if !ok || tag == Ignore {
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
		c.Dependencies = append(c.Dependencies, cup)
		value.Set(cup.Value)
	}
}

// 传入接口类型，并根据接口调用
func (c *Cup) lifeCycle(set CupSet, t reflect.Type) {
	if _, has := set[c]; has {
		return
	}
	set[c] = struct{}{}
	for _, dependency := range c.Dependencies {
		dependency.lifeCycle(set, t)
	}
	if c.Class.Implements(t) {
		c.Value.MethodByName(t.Method(0).Name).Call(nil)
	}
}
