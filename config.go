package config

import (
	//"examples/config/config/file"
	"flag"
	"fmt"
	"os"
)

type Config interface {
	Bind(string, interface{}) error
}

var (
	_path string

	defaultConfig Config
)

func init() {
	flag.StringVar(&_path, "conf", os.Getenv("CONF_PATH"), `config file path.`)
}

func Init() error {
	var err error
	if _path != "" {
		if defaultConfig, err = NewFile(_path); err == nil {
			return nil
		}
	} else if _nacosHost != "" {
		if defaultConfig, err = NewNacos(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("init config err")
}

func Bind(key string, value interface{}) error {
	return defaultConfig.Bind(key, value)
}
