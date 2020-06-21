package consul

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type EventedConsul struct {
	Log        *zap.SugaredLogger
	ConsulHost string
}

func (o *EventedConsul) Register(name string, id string, port uint) error {
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
		Port:    int(port),
		Connect: &api.AgentServiceConnect{},
	})
	return err
}
