package support

import (
	"fmt"
	"github.com/dsnet/try"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"os"
)

type ConfigInit interface {
	SetConsulHost(host string)
	SetConsulKey(key string)
	SetName(name string)
	GetName() string
	SetConfigMgmt(mgmt string)
	SetConfigFormat(configType string)
}

type BasicConfigInit struct {
	ConsulHost string
	ConsulKey  string
	Name       string
	ConfigMgmt string
	ConfigType string
}

func (c *BasicConfigInit) SetConsulKey(key string) {
	c.ConsulKey = key
}

func (c *BasicConfigInit) SetName(name string) {
	c.Name = name
}

func (c *BasicConfigInit) GetName() string {
	return c.Name
}

func (c *BasicConfigInit) SetConfigMgmt(mgmt string) {
	c.ConfigMgmt = mgmt
}

func (c *BasicConfigInit) SetConfigFormat(configType string) {
	c.ConfigType = configType
}

func (c *BasicConfigInit) SetConsulHost(host string) {
	c.ConsulHost = host
}

const ConfigMgmtType = "CONFIG_MGMT_TYPE"
const ConfigFormat = "CONFIG_FORMAT"
const ConsulHostType = "CONSUL_HOST"
const ConsulKeyType = "CONSUL_KEY"
const LocalConsul = "localhost:8500"

type ConfigType string

const (
	Consul ConfigType = "consul"
	File   ConfigType = "file"
)

type ConfigFormatType string

const (
	yaml ConfigFormatType = "yaml"
	json ConfigFormatType = "json"
)
const AppNameType = "APP_NAME"
const AppName = "UnnamedEventedApplication"

func Initialize(log *zap.SugaredLogger, config ConfigInit) (any, error) {
	viper.AutomaticEnv()
	viper.SetDefault(ConfigMgmtType, Consul)
	viper.SetDefault(ConfigFormat, yaml)
	viper.SetDefault(ConsulHostType, LocalConsul)
	viper.SetDefault(AppNameType, AppName)
	viper.SetConfigType(string(yaml))

	configMgmt := viper.GetString(ConfigMgmtType)
	log.Infow("Configuring.", ConfigMgmtType, configMgmt)
	switch ConfigType(configMgmt) {
	case File:
		viper.AddConfigPath(".")
		viper.AddConfigPath(fmt.Sprintf("%s/.evented/%s/", try.E1(os.UserHomeDir()), config.GetName()))
		viper.SetConfigType(viper.GetString(ConfigFormat))
		viper.SetConfigName("serve")
		viper.AutomaticEnv()
		try.E(viper.ReadInConfig())
		try.E(viper.Unmarshal(config))
		return config, nil
	case Consul:
		consulHost := viper.GetString(ConsulHostType)
		consulKey := viper.GetString(ConsulKeyType)
		log.Infow("Attempting to fetch configuration", "provider", configMgmt, "host", consulHost, "key", consulKey)
		try.E(viper.AddRemoteProvider("consul", consulHost, consulKey))
		try.E(viper.ReadRemoteConfig())
		try.E(viper.Unmarshal(config))
		log.Infow("Configuration set", "configuration", fmt.Sprintf("%+v", config))

		log.Infow("Read consul.", "Proof", viper.GetString("proof"))
		return config, nil
	default:
		return nil, fmt.Errorf("configuration provider %s not matched with valid types", configMgmt)
	}
}
