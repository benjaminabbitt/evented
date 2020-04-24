package main

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) Port() uint {
	return viper.GetUint("port")
}
func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}
