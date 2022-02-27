package support

import (
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

type ConfigInit struct {
	consulHost       string
	consulKey        string
	name             string
	configMgmt       string
	consulConfigType string
}

func (o *ConfigInit) AppName() string {
	return o.name
}

func (o *ConfigInit) ConsulHost() string {
	return o.consulHost
}

func (o *ConfigInit) ConsulKey() string {
	return o.consulKey
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

func (o *ConfigInit) Initialize(log *zap.SugaredLogger, config interface{}) interface{} {
	viper.AutomaticEnv()
	viper.SetDefault(ConfigMgmtType, Consul)
	viper.SetDefault(ConsulConfigType, Yaml)
	viper.SetDefault(ConsulHostType, LocalConsul)
	viper.SetDefault(AppNameType, AppName)
	viper.SetConfigType(Yaml)

	o.configMgmt = viper.GetString(ConfigMgmtType)
	o.name = viper.GetString(AppNameType)
	log.Infow("Configuring.", ConfigMgmtType, o.configMgmt)
	if o.configMgmt == Consul {
		o.consulHost = viper.GetString(ConsulHostType)
		o.consulKey = viper.GetString(ConsulKeyType)
		o.consulConfigType = viper.GetString(ConsulConfigType)
		log.Infow("Attempting to fetch configuration", "provider", o.configMgmt, "host", o.consulHost, "key", o.consulKey)
		err := viper.AddRemoteProvider("consul", o.consulHost, o.consulKey)
		if err != nil {
			log.Fatal(err)
		}
		err = viper.ReadRemoteConfig()
		if err != nil {
			log.Fatalw(err.Error(), "key", o.consulKey)
		}
		err = viper.Unmarshal(config)
		log.Infow("Configuration set", "configuration", fmt.Sprintf("%+v", config))
		if err != nil {
			log.Fatalw(err.Error(), "key", o.consulKey)
		}

		log.Infow("Read consul.", "Proof", viper.GetString("proof"))
		return config
	}
	return nil
}
