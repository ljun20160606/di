package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Bar1 struct {
}

func (b *Bar1) BeforeInit() interface{} {
	return "2"
}

type Foo1 struct {
	A string `di:"#.ab"`
	B string `di:"@.*"`
	C string `di:"@.di.Bar1"`
}

func TestTomlLoad(t *testing.T) {
	ast := assert.New(t)

	ast.NoError(TomlLoad(`ab="1"`))

	bar1 := new(Bar1)
	foo1 := new(Foo1)

	Put(bar1)
	Put(foo1)
	PutWithName(foo1, "Foo1")

	Start()

	ast.Equal(bar1, GetWithName("di.Bar1"))
	ast.Equal(foo1, GetWithName("di.Foo1"))
	ast.Equal(foo1, GetWithName("Foo1"))

	ast.Equal("1", foo1.A)
	ast.Equal("2", foo1.B)
	ast.Equal("2", foo1.C)
}
