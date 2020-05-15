package support

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ConfigInit struct {
	AppName      *string
	configPath   *string
	configSource *string
	consulHost   *string
}

func (o *ConfigInit) Name() string {
	return *o.AppName
}

func (o *ConfigInit) ConfigPath() string {
	return *o.configPath
}

func (o *ConfigInit) Initialize(name string, log *zap.SugaredLogger) {
	o.AppName = flag.String("appName", name, "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	o.configPath = flag.String("configFilePath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	o.configSource = flag.String("config", "file", "The configuration source of the application.  Valid values are \"file\" and \"consul\".")
	o.consulHost = flag.String("consulHost", "localhost:8500", "The consul URL to bootstrap against.")

	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Error(err)
	}
	viper.SetConfigName(o.Name())

	if *o.configSource == "consul" {
		err = viper.AddRemoteProvider("consul", *o.consulHost, *o.AppName)
		if err != nil {
			log.Error(err)
		}

		viper.SetConfigType("json") // Need to explicitly set this to json

		err = viper.ReadRemoteConfig()
		if err != nil {
			log.Error(err)
		}
		log.Info("test consul: " + viper.GetString(fmt.Sprintf("%s.port", *o.AppName)))
	} else if *o.configSource == "file" {
		viper.SetConfigType("yaml")

		viper.AddConfigPath(o.ConfigPath())

		viper.SetEnvPrefix(o.Name())
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Errorw("Config file not found.", "path", o.ConfigPath(), "name", o.Name())
				// Config file not found; ignore error if desired
			} else {
				log.Error(err)
			}
		}
	}

}
