package support

import (
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

func (o *ConfigInit) Initialize(log *zap.SugaredLogger) {
	viper.AutomaticEnv()
	viper.SetDefault(ConfigMgmtType, "consul")
	viper.SetDefault(ConsulConfigType, "yaml")

	o.configMgmt = viper.GetString(ConfigMgmtType)
	o.name = viper.GetString("APP_NAME")
	log.Infow("Configuring.", ConfigMgmtType, o.configMgmt)
	if o.configMgmt == "consul" {
		o.consulHost = viper.GetString("CONSUL_HOST")
		o.consulKey = viper.GetString("CONSUL_KEY")
		o.consulConfigType = viper.GetString(ConsulConfigType)
		err := viper.AddRemoteProvider("consul", o.consulHost, o.consulKey)
		if err != nil {
			log.Fatal(err)
		}
		viper.SetConfigType("yaml")
		err = viper.ReadRemoteConfig()
		if err != nil {
			log.Fatal(err)
		}

		log.Infow("Read consul.", "Proof", viper.GetString("proof"))
	}

}
