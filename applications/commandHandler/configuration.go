package main

import (
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) BusinessURL() string {
	return viper.GetString("business.url")
}

func (o *Configuration) Port() uint {
	return viper.GetUint("port")
}

func (o *Configuration) Domain() string {
	return viper.GetString("domain")
}

func (o *Configuration) SagaURLs() (urls []string) {
	sagaConfig := viper.GetStringMap("sync.sagas")
	for name, _ := range sagaConfig {
		urls = append(urls, viper.GetString("sync.sagas."+name+".url"))
	}
	return urls
}

func (o *Configuration) ProjectorURLs() (urls []string) {
	projectorConfig := viper.GetStringSlice("sync.projectors")
	for _, url := range projectorConfig {
		urls = append(urls, url)
	}
	return urls
}
