package di

import (
	"github.com/ljun20160606/gox/reflectx"
	"log"
	"os"
)

const (
	DI     = "di"
	Ignore = "-"
)

var (
	loggerConfiguration        = WithLogger(log.New(os.Stderr, "["+DI+"] ", log.Lshortfile))
	loggerVerboseConfiguration = WithLoggerVerbose(false)
	defaultContainer           = NewContainer(loggerConfiguration, loggerVerboseConfiguration)
)

func Put(water Water) {
	defaultContainer.Put(water)
}

func PutWithName(water Water, name string) {
	defaultContainer.PutWithName(water, name)
}

func GetWithName(name string) Water {
	return defaultContainer.GetWaterWithName(name)
}

func RegisterPlugin(lifecycle Lifecycle, p Plugin) {
	defaultContainer.RegisterPlugin(lifecycle, p)
}

func Configure(configurators ...Configurator) Container {
	return defaultContainer.Configure(configurators...)
}

func Start() {
	defaultContainer.Start()
}

func GetDefaultName(water Water) string {
	return reflectx.GetInterfaceDefaultName(water)
}
