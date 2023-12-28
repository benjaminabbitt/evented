package configuration

import (
	"github.com/benjaminabbitt/evented/support"
	_ "github.com/spf13/viper/remote"
)

type Configuration struct {
	support.BasicConfigInit
	Business struct {
		Url string
	}
	Port   uint
	Domain string
	Sync   struct {
		Sagas []struct {
			Name string
			Url  string
		}
		Projectors []struct {
			Name string
			Url  string
		}
	}
	Snapshots struct {
		Kind    string
		Mongodb SnapshotStore
	}
	Transport struct {
		Kind     string
		Rabbitmq struct {
			Url      string
			Exchange string
		}
	}
	Events struct {
		Kind    string
		Mongodb struct {
			Url        string
			Name       string
			Collection string
		}
	}
}

type SnapshotStore struct {
	Url        string
	Name       string
	Collection string
}

func (o Configuration) SnapshotStore() SnapshotStore {
	return o.Snapshots.Mongodb
}
