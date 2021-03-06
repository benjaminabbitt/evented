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

func (o *Configuration) ProjectorURL() string {
	return viper.GetString("projector.url")
}

func (o *Configuration) Name() string {
	return viper.GetString("name")
}

func (o *Configuration) Port() uint {
	return uint(viper.GetInt64("port"))
}

func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("queryHandler.url")
}

func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
