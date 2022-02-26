package configuration

import (
	"fmt"
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) DatabaseType() string {
	return viper.GetString("database.type")
}

func (o *Configuration) DatabaseURL() string {
	return viper.GetString(fmt.Sprintf("database.%s.url", o.DatabaseType()))
}

func (o *Configuration) DatabaseName() string {
	return viper.GetString(fmt.Sprintf("database.%s.name", o.DatabaseType()))
}

func (o *Configuration) CollectionName() string {
	return viper.GetString(fmt.Sprintf("database.%s.collectionName", o.DatabaseType()))
}

func (o *Configuration) ProjectorURL() string {
	return viper.GetString("evented.url")
}

func (o *Configuration) Name() string {
	return viper.GetString("name")
}

func (o *Configuration) Port() uint {
	return uint(viper.GetInt64("port"))
}

func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("query-handler.url")
}

func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
