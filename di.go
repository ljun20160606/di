package di

import (
	"log"
	"os"
)

const (
	DI     = "di"
	Ignore = "-"
)

var (
	Verbose          = true
	logger           = log.New(os.Stderr, "["+DI+"] ", log.Lshortfile)
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
