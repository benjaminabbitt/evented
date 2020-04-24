package support

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ConfigInit struct{}

func (o *ConfigInit) Name() string {
	return *flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
}

func (o *ConfigInit) ConfigPath() string {
	return *flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
}

func (o *ConfigInit) Initialize(log *zap.SugaredLogger) {
	flag.Parse()

	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Error(err)
	}

	viper.SetConfigName(o.Name())
	viper.SetConfigType("yaml")

	viper.AddConfigPath(o.ConfigPath())

	viper.SetEnvPrefix(o.Name())
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			log.Error(err)
		}
	}
}
