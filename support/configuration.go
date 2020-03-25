package support

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func SetupConfig(name *string, configpath *string, cmdLine *flag.FlagSet)(err error) {
	err = viper.BindPFlags(cmdLine)
	if err != nil {
		return err
	}

	viper.SetConfigName(*name)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(*configpath)

	viper.SetEnvPrefix(*name)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			return err
		}
	}
	return nil
}