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
	//TODO: this will need to be an array of other command handlers, by domain
	return viper.GetString("otherCommandHandler.url")
}

func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("query-handler.url")
}

func (o *Configuration) AMQPURL() string {
	return viper.GetString("transport.url")
}

func (o *Configuration) AMQPExchange() string {
	return viper.GetString("transport.exchange")
}

func (o *Configuration) AMQPQueue() string {
	return viper.GetString("transport.queue")
}

func (o *Configuration) BusinessURL() string {
	return viper.GetString("placeholder-saga.url")
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
func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
