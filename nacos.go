package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	_nacosHost        string
	_nacosNamespaceId string //默认为空
	_nacosGroup       string //默认 DEFAULT_GROUP
)

func init() {
	flag.StringVar(&_nacosHost, "nacos.host", os.Getenv("NACOS_HOST"), `config api host.`)
	flag.StringVar(&_nacosNamespaceId, "nacos.namespaceid", os.Getenv("NACOS_NAMESPACEID"), `nacos namespace id`)
	flag.StringVar(&_nacosGroup, "nacos.group", os.Getenv("NACOS_GROUP"), `nacos group`)
}

type nacos struct {
	client config_client.IConfigClient
}

func NewNacos() (Config, error) {
	var (
		err  error
		ip   string
		port uint64 = 8849
	)

	as := strings.Split(_nacosHost, ":")
	if len(as) == 2 {
		port, _ = strconv.ParseUint(as[1], 10, 64)
	}
	ip = as[0]

	if _nacosGroup == "" {
		_nacosGroup = "DEFAULT_GROUP"
	}
	//create ServerConfig  "118.195.198.233:8849"
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(_nacosNamespaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(filepath.Join(os.TempDir(), "nacos", "log")),
		constant.WithCacheDir(filepath.Join(os.TempDir(), "nacos", "cache")),
		//constant.WithLogLevel("debug"),
	)
	n := &nacos{}
	// create config client
	if n.client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	); err != nil {
		return nil, errors.Wrap(err, "newNacos faild")
	}
	return n, nil
}

func (n *nacos) Get(key string, value interface{}) error {
	v, err := n.client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  _nacosGroup,
	})
	if err != nil {
		fmt.Printf("get key err is (%+v)\n", err)
		return err
	}

	if v == "" {
		value = ""
		return nil
	}
	return n.unmarshal([]byte(v), value)
}

func (n *nacos) unmarshal(v []byte, out interface{}) error {
	var err error
	if err = json.Unmarshal(v, &out); err == nil {
		return nil
	}
	if err = yaml.Unmarshal(v, &out); err == nil {
		return nil
	}
	return err
}
