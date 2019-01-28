package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Duck struct {
	Name string `di:"#.name"`
}

func TestConfigLoadJson(t *testing.T) {
	reset()
	ast := assert.New(t)
	_ = ConfigLoad(`{"name":"duck"}`, JSON)
	//di.ConfigLoadFile("path", JSON)
	//di.ConfigLoadReader(reader, JSON)
	duck := Duck{}
	Put(&duck)
	Start()

	ast.Equal("duck", duck.Name)
}

func TestConfigLoadToml(t *testing.T) {
	reset()
	ast := assert.New(t)
	_ = ConfigLoad(`name = "duck"`, TOML)
	//di.ConfigLoadFile("path", TOML)
	//di.ConfigLoadReader(reader, TOML)
	duck := Duck{}
	Put(&duck)
	Start()

	ast.Equal("duck", duck.Name)
}

func TestConfigLoadYaml(t *testing.T) {
	reset()
	ast := assert.New(t)
	_ = ConfigLoad(`name: duck`, YAML)
	//di.ConfigLoadFile("path", YAML)
	//di.ConfigLoadReader(reader, YAML)
	duck := Duck{}
	Put(&duck)
	Start()

	ast.Equal("duck", duck.Name)
}

type Prefix struct {
	Properties PrefixProperties `di:"#.{test.properties}"`
}

type AnotherPrefix struct {
	Properties *PrefixProperties `di:"#.{test.properties}"`
}

type PrefixProperties struct {
	Type        string
	Value       string
	SnakeFormat string `yaml:"snake_format"`
}

func TestPrefixYaml(t *testing.T) {
	reset()
	ast := assert.New(t)
	_ = ConfigLoad(`
test:
  properties:
    type: prefix
    value: test
    snake_format: snake`, YAML)

	prefix := Prefix{}
	Put(&prefix)
	anotherPrefix := AnotherPrefix{}
	Put(&anotherPrefix)
	Start()

	ast.Equal("prefix", prefix.Properties.Type)
	ast.Equal("test", prefix.Properties.Value)
	ast.Equal("snake", prefix.Properties.SnakeFormat)
	ast.Equal(prefix.Properties, *anotherPrefix.Properties)
}
