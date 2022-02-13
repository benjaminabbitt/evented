package configuration

import (
	"fmt"
	"github.com/benjaminabbitt/evented/support"
	"github.com/spf13/viper"
)

type Configuration struct {
	support.ConfigInit
}

func (o *Configuration) DatabaseURL() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.url", o.DatabaseType()))
}

func (o *Configuration) DatabaseType() string {
	return viper.GetString("eventStore.type")
}

func (o *Configuration) DatabaseName() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.database", o.DatabaseType()))
}

func (o *Configuration) DatabaseCollection() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.collection", o.DatabaseType()))
}

func (o *Configuration) Port() uint {
	return uint(viper.GetInt64("port"))
}

func (o *Configuration) EventBookTargetSize() uint {
	return uint(viper.GetInt64("targetSize"))
}
