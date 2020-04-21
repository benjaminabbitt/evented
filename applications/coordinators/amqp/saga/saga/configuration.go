package saga

import "github.com/spf13/viper"

type Configuration struct {
}

func (o *Configuration) CommandHandlerURL() string {
	return viper.GetString("commandHandler.url")
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
