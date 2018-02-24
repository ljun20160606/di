package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Bar interface {
	Values() string
}

type (
	BarImpl struct {
		Value string
	}

	Foo struct {
		BarAny *BarImpl `di:"*"`
		Bar    Bar      `di:"*"`
		Sex    *string  `di:"*"`
	}
)

func (f *BarImpl) Values() string {
	return f.Value
}

func TestPut(t *testing.T) {
	ast := assert.New(t)

	foo := new(Foo)
	bar := &BarImpl{Value: "dc"}
	sex := "sex"

	Put(foo)
	Put(bar)
	Put(&sex)

	Start()

	ast.Equal(bar, foo.BarAny)
	ast.Equal(bar, foo.Bar)
	ast.Equal(&sex, foo.Sex)
	ast.Equal(bar.Value, foo.Bar.Values())
}
