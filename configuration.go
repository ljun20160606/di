package di

import "log"

type Configurator func(Container)

func WithLoggerVerbose(verbose bool) Configurator {
	return func(c Container) {
		c.(*container).verbose = verbose
	}
}

func WithLogger(logger *log.Logger) Configurator {
	return func(c Container) {
		c.(*container).logger = logger
	}
}
