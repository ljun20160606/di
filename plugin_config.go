package di

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/ljun20160606/gox/fs"
	"github.com/ljun20160606/jsonpath"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"reflect"
	"sync"
)

type Unmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
}

type UnmarshalFunc func(data []byte, v interface{}) error

func (u UnmarshalFunc) Unmarshal(data []byte, v interface{}) error {
	return u(data, v)
}

var (
	JSON Unmarshaler = UnmarshalFunc(json.Unmarshal)
	TOML Unmarshaler = UnmarshalFunc(func(data []byte, v interface{}) error {
		_, err := toml.Decode(string(data), v)
		return err
	})
	YAML Unmarshaler = UnmarshalFunc(yaml.Unmarshal)
)

var (
	loadPluginConfigOnce = sync.Once{}
	configPlugin         = new(diConfig)
)

func loadConfig() {
	loadPluginConfigOnce.Do(func() {
		RegisterPlugin(BeforeInit, configPlugin)
	})
}

type diConfig struct {
	value interface{}
}

func (d *diConfig) Load(path string, ice Ice) {
	v, err := jsonpath.JsonPathLookup(d.value, "$."+path)
	if v == nil {
		panic(errors.Wrap(err, ErrorMissing.Panic("#."+path)))
	}
	ice.Value().Set(reflect.ValueOf(v))
}

func (d *diConfig) Prefix() string {
	return "#"
}

func (d *diConfig) Priority() int {
	return 0
}

func ConfigLoadFile(path string, ct Unmarshaler) error {
	return fs.ReadFile(path, func(reader io.Reader) error {
		return ConfigLoadReader(reader, ct)
	})
}

func ConfigLoadReader(reader io.Reader, ct Unmarshaler) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ConfigLoad(string(bytes), ct)
}

func ConfigLoad(content string, unmarshaler Unmarshaler) error {
	loadConfig()

	err := unmarshaler.Unmarshal([]byte(content), &(configPlugin.value))
	if err != nil {
		return err
	}
	return nil
}
