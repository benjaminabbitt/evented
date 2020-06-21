package support

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ConfigInit struct {
	consulHost string
	consulKey  string
	name       string
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

func (o *ConfigInit) Initialize(log *zap.SugaredLogger) {
	viper.AutomaticEnv()
	o.consulHost = viper.GetString("CONSUL_HOST")
	o.consulKey = viper.GetString("CONSUL_KEY")
	o.name = viper.GetString("APP_NAME")

	log.Infow("Attempting to reach Consul K/V", "host", o.consulHost, "key", o.consulKey)

	var err error
	err = viper.AddRemoteProvider("consul", o.consulHost, o.consulKey)
	if err != nil {
		log.Error(err)
	}

	viper.SetConfigType("yaml")

	err = viper.ReadRemoteConfig()
	if err != nil {
		log.Error(err)
	}
}
