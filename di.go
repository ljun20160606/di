package di

import (
	"log"
	"os"
)

var (
	logger           = log.New(os.Stderr, "[di] ", log.Lshortfile)
	defaultContainer = NewContainer()
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

func Start() {
	defaultContainer.Start()
}
