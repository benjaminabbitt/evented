package consul

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func NewEventedConsul(host string, port uint) (consul *EventedConsul) {
	consul = &EventedConsul{
		ConsulHost: host,
		ConsulPort: port,
	}
	return consul
}

type EventedConsul struct {
	Log        *zap.SugaredLogger
	ConsulHost string
	ConsulPort uint
}

func (o *EventedConsul) Register(name string, id string) error {
	client, err := api.NewClient(&api.Config{
		Address: o.ConsulHost,
	})
	if err != nil {
		o.Log.Error(err)
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      name + "-" + id,
		Name:    name,
		Tags:    []string{"evented", name},
		Address: o.ConsulHost,
		Port:    int(o.ConsulPort),
		Connect: &api.AgentServiceConnect{},
	})
	return err
}
