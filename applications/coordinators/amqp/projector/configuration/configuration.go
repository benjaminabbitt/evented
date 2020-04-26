package configuration

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) ProjectorURL() string {
	return viper.GetString("projector.url")
}

func (o *Configuration) QueryHandlerURL() string {
	return viper.GetString("queryHandler.url")
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
	return viper.GetString("saga.url")
}

func (o *Configuration) DatabaseURL() string {
	return viper.GetString("database.url")
}

func (o *Configuration) DatabaseName() string {
	return viper.GetString("database.name")
}
func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
