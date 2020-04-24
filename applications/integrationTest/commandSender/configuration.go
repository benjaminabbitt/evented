package main

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) CommandHandlerURL() string {
	return viper.GetString("commandHandlerURL")
}
func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
