package configuration

import (
	"fmt"
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
	return viper.GetStringSlice("sync.projectors")
}

func (o *Configuration) SnapshotStoreType() string {
	return viper.GetString("snapshotStore.type")
}

func (o *Configuration) SnapshotStoreURL() string {
	return viper.GetString(fmt.Sprintf("snapshotStore.%s.url", o.SnapshotStoreType()))
}

func (o *Configuration) SnapshotStoreDatabaseName() string {
	return viper.GetString(fmt.Sprintf("snapshotStore.%s.database", o.SnapshotStoreType()))
}

func (o *Configuration) TransportType() string {
	return viper.GetString("transport.type")
}

func (o *Configuration) TransportURL() string {
	return viper.GetString(fmt.Sprintf("transport.%s.url", o.TransportType()))
}

func (o *Configuration) TransportExchange() string {
	return viper.GetString(fmt.Sprintf("transport.%s.exchange", o.TransportType()))
}

func (o *Configuration) EventRepoType() string {
	return viper.GetString("eventStore.type")
}

func (o *Configuration) EventStoreURL() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.url", o.EventRepoType()))
}

func (o *Configuration) EventStoreDatabaseName() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.database", o.EventRepoType()))
}

func (o *Configuration) EventStoreCollectionName() string {
	return viper.GetString(fmt.Sprintf("eventStore.%s.collection", o.EventRepoType()))
}
