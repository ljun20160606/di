package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func reset() {
	c := defaultContainer.(*container)
	c.cupMap = make(map[string][]*Cup)
	c.iceMap = make(map[string][]Ice)
}

func TestPut(t *testing.T) {
	reset()

	type Duck struct {
		Name *string `di:"*"`
	}

	ast := assert.New(t)

	ast.Panics(func() {
		Put("duck")
	})

	duck := Duck{}

	Put(&duck)

	ast.NotPanics(func() {
		duckName := "duck"
		Put(&duckName)

		Start()

		ast.Equal(&duckName, duck.Name)
	})
}

func TestPutWithName(t *testing.T) {
	reset()

	type Duck struct {
		Name *string `di:"duckName"`
	}

	ast := assert.New(t)

	duck := Duck{}

	Put(&duck)

	ast.NotPanics(func() {
		duckName := "duck"
		PutWithName(&duckName, "duckName")
		Start()

		ast.Equal(&duckName, duck.Name)
	})
}

func TestGetWithName(t *testing.T) {
	reset()

	ast := assert.New(t)

	name := "duck"

	Put(&name)

	withName := GetWithName("string")

	ast.Equal(&name, withName)
}

func TestTomlLoad(t *testing.T) {
	reset()

	type Duck struct {
		Name string `di:"#.name"`
	}

	ast := assert.New(t)

	ast.NoError(TomlLoad(`name = "duck"`))

	duck := Duck{}

	Put(&duck)

	Start()

	ast.Equal("duck", duck.Name)
}

func TestYamlLoad(t *testing.T) {
	reset()

	type Duck struct {
		Name string `di:"#.name"`
	}

	ast := assert.New(t)

	ast.NoError(ConfigLoad(`name: duck`, YAML))

	duck := Duck{}

	Put(&duck)

	Start()

	ast.Equal("duck", duck.Name)
}

type Sex struct {
}

func (s *Sex) BeforeInit() interface{} {
	return "male"
}

func TestBeforeInitPlugin(t *testing.T) {
	reset()

	type Duck struct {
		Sex    string `di:"@.*"`
		SexToo string `di:"@.di.Sex"`
	}

	ast := assert.New(t)

	duck := Duck{}
	sex := Sex{}

	Put(&duck)
	Put(&sex)

	Start()

	ast.Equal(sex.BeforeInit(), duck.SexToo)
	ast.Equal(duck.Sex, duck.SexToo)
}
