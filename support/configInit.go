package support

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

type ConfigInit interface {
	SetConsulHost(host string)
	SetConsulKey(key string)
	SetName(name string)
	SetConfigMgmt(mgmt string)
	SetConsulConfigType(configType string)
}

type BasicConfigInit struct {
	ConsulHost       string
	ConsulKey        string
	Name             string
	ConfigMgmt       string
	ConsulConfigType string
}

func (c *BasicConfigInit) SetConsulKey(key string) {
	c.ConsulKey = key
}

func (c *BasicConfigInit) SetName(name string) {
	c.Name = name
}

func (c *BasicConfigInit) SetConfigMgmt(mgmt string) {
	c.ConfigMgmt = mgmt
}

func (c *BasicConfigInit) SetConsulConfigType(configType string) {
	c.ConsulConfigType = configType
}

func (c *BasicConfigInit) SetConsulHost(host string) {
	c.ConsulHost = host
}

const ConfigMgmtType = "CONFIG_MGMT_TYPE"
const ConsulConfigType = "CONSUL_CONFIG_TYPE"
const ConsulHostType = "CONSUL_HOST"
const ConsulKeyType = "CONSUL_KEY"
const LocalConsul = "localhost:8500"
const Consul = "consul"
const Yaml = "yaml"
const AppNameType = "APP_NAME"
const AppName = "UnnamedEventedApplication"

func Initialize(log *zap.SugaredLogger, config ConfigInit) ConfigInit {
	viper.AutomaticEnv()
	viper.SetDefault(ConfigMgmtType, Consul)
	viper.SetDefault(ConsulConfigType, Yaml)
	viper.SetDefault(ConsulHostType, LocalConsul)
	viper.SetDefault(AppNameType, AppName)
	viper.SetConfigType(Yaml)

	//Viper doesn't handle environment variable mapping properly, so workaround by setting manually
	config.SetConfigMgmt(viper.GetString(ConfigMgmtType))
	config.SetConsulHost(viper.GetString(ConsulHostType))
	config.SetConsulConfigType(viper.GetString(ConsulConfigType))
	config.SetName(viper.GetString(AppNameType))
	config.SetConsulKey(viper.GetString(ConsulKeyType))

	configMgmt := viper.GetString(ConfigMgmtType)
	log.Infow("Configuring.", ConfigMgmtType, configMgmt)
	if configMgmt == Consul {
		consulHost := viper.GetString(ConsulHostType)
		consulKey := viper.GetString(ConsulKeyType)
		log.Infow("Attempting to fetch configuration", "provider", configMgmt, "host", consulHost, "key", consulKey)
		err := viper.AddRemoteProvider("consul", consulHost, consulKey)
		if err != nil {
			log.Fatal(err)
		}
		err = viper.ReadRemoteConfig()
		if err != nil {
			log.Fatalw(err.Error(), "key", consulKey)
		}
		err = viper.Unmarshal(config)
		log.Infow("Configuration set", "configuration", fmt.Sprintf("%+v", config))
		if err != nil {
			log.Fatalw(err.Error(), "key", consulKey)
		}

		log.Infow("Read consul.", "Proof", viper.GetString("proof"))
		return config
	}
	return nil
}
