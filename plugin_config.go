package di

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/ljun20160606/gox/fs"
	"github.com/ljun20160606/gox/reflectx"
	"github.com/ljun20160606/jsonpath"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
)

var (
	JSON Serializer = new(jsonSerializer)
	TOML Serializer = new(tomlSerializer)
	YAML Serializer = new(yamlSerializer)
)

type Serializer interface {
	Unmarshal(data []byte, v interface{}) error
	Marshal(v interface{}) ([]byte, error)
}

type jsonSerializer struct {
}

func (*jsonSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (*jsonSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type tomlSerializer struct {
}

func (*tomlSerializer) Unmarshal(data []byte, v interface{}) error {
	_, err := toml.Decode(string(data), v)
	return err
}

func (*tomlSerializer) Marshal(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := toml.NewEncoder(buffer).Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

type yamlSerializer struct {
}

func (*yamlSerializer) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (*yamlSerializer) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

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
	value      interface{}
	serializer Serializer
}

func (d *diConfig) Load(path string, ice Ice) {
	// prefix 模式
	value := ice.Value()
	if strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}") {
		defaultName := reflectx.GetTypeDefaultName(ice.Type())
		water := GetWithName(defaultName)
		if water != nil {
			if value.Kind() == reflect.Ptr {
				value.Set(reflect.ValueOf(water))
			} else {
				value.Set(reflect.ValueOf(water).Elem())
			}
			return
		}
		prefix := path[1 : len(path)-1]
		var b []byte
		// 如果prefix则取本身
		if prefix == "" {
			b, _ = d.serializer.Marshal(d.value)
		} else {
			v, err := jsonpath.JsonPathLookup(d.value, "$."+prefix)
			if err != nil {
				panic(errors.Wrap(err, ErrorMissing.Panic("#."+path)))
			}
			b, _ = d.serializer.Marshal(v)
		}
		var vv interface{}
		if value.Kind() == reflect.Ptr {
			vv = value.Interface()
		} else {
			vv = value.Addr().Interface()
		}
		err := d.serializer.Unmarshal(b, vv)
		if err != nil {
			panic(errors.Wrap(err, ErrorUnserialize.String()))
		}
		// 放入容器
		Put(vv)
		return
	}
	// value 模式
	v, err := jsonpath.JsonPathLookup(d.value, "$."+path)
	if err != nil {
		panic(errors.Wrap(err, ErrorMissing.Panic("#."+path)))
	}
	value.Set(reflect.ValueOf(v))
}

func (d *diConfig) Prefix() string {
	return "#"
}

func (d *diConfig) Priority() int {
	return 0
}

func ConfigLoadFile(path string, ct Serializer) error {
	return fs.ReadFile(path, func(reader io.Reader) error {
		return ConfigLoadReader(reader, ct)
	})
}

func ConfigLoadReader(reader io.Reader, ct Serializer) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return ConfigLoad(string(b), ct)
}

func ConfigLoad(content string, serializer Serializer) error {
	loadConfig()

	err := serializer.Unmarshal([]byte(content), &(configPlugin.value))
	if err != nil {
		return err
	}
	configPlugin.serializer = serializer
	return nil
}
