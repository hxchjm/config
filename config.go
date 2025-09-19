package config

import (
	//"examples/config/config/file"
	"flag"
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
		if defaultConfig, err = NewFile(_path); err != nil {
			return err
		}
	} else {
		if defaultConfig, err = NewNacos(); err != nil {
			return err
		}
	}
	return nil
}

func Bind(key string, value interface{}) error {
	return defaultConfig.Bind(key, value)
}
