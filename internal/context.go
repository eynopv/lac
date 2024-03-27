package internal

import "sync"

type context struct {
	Config *Config
}

var instance *context
var once sync.Once

func GetContext() *context {
	once.Do(func() {
		instance = &context{Config: &Config{}}
	})

	return instance
}
