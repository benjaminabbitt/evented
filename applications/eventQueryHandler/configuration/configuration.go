package configuration

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) DatabaseURL() string {
	return viper.GetString("database.url")
}

func (o *Configuration) DatabaseName() string {
	return viper.GetString("database.name")
}

func (o *Configuration) DatabaseCollection() string {
	return viper.GetString("database.collection")
}

func (o *Configuration) Port() uint {
	return viper.GetUint("port")
}

func (o *Configuration) EventBookTargetSize() uint {
	return viper.GetUint("targetSize")
}
