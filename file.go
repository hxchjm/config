package config

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"os"
)

var (
	ERRKeyNotFound = errors.New("Key not found")
)

type file struct {
	path    string
	content map[string]interface{}
}

func NewFile(path string) (Config, error) {
	f := &file{
		path:    path,
		content: make(map[string]interface{}),
	}
	if err := f.Load(); err != nil {
		return nil, errors.Wrap(err, "NewFile faild")
	}
	return f, nil
}

func (f *file) Get(key string, value interface{}) error {
	if v, ok := f.content[key]; ok {
		return f.unmarshal(v, value)
	}
	return  errors.Wrap(ERRKeyNotFound, fmt.Sprintf("(%s)", key))
}

func (f *file) Load() error {
	bs, err := os.ReadFile(f.path)
	if err != nil {
		return errors.Wrapf(err, "read file failed")
	}
	if err = json.Unmarshal(bs, &f.content); err == nil {
		return nil
	}
	if err = yaml.Unmarshal(bs, &f.content); err == nil {
		return nil
	}
	return errors.Wrap(err, "file Load err,not json or yaml")
}

func (f *file) unmarshal(v interface{}, out interface{}) error {
	var err error
	if bs, err := json.Marshal(v); err == nil {
		if err := json.Unmarshal(bs, &out); err == nil {
			return nil
		}
	}
	if bs, err := yaml.Marshal(v); err == nil {
		if err = yaml.Unmarshal(bs, &out); err == nil {
			return nil
		}
	}
	return err
}
