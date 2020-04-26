package configuration

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) OtherCommandHandlerURL() string {
	return viper.GetString("otherCommandHandlerURL")
}

func (o *Configuration) DatabaseURL() string {
	return viper.GetString("database.url")
}

func (o *Configuration) TargetURL() string {
	return viper.GetString("target.url")
}

func (o *Configuration) DatabaseName() string {
	return viper.GetString("database.name")
}

func (o *Configuration) Port() uint {
	return viper.GetUint("port")
}

func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("queryHandler.url")
}
