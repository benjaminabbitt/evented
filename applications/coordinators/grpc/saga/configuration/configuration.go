package configuration

import (
	"fmt"
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) OtherCommandHandlerURL() string {
	//TODO: map of other command handler
	return viper.GetString("otherCommandHandler.url")
}

func (o *Configuration) SagaURL() string {
	return viper.GetString("saga.url")
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

func (o *Configuration) Port() uint {
	return uint(viper.GetInt64("port"))
}

func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("queryHandler.url")
}
