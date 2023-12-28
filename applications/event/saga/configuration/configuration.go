package configuration

import "github.com/benjaminabbitt/evented/support"

type Config struct {
	support.BasicConfigInit
	Proof        string
	QueryHandler struct {
		Url string
	}
	Transport struct {
		Kind     string
		Rabbitmq struct {
			Url      string
			Exchange string
			Queue    string
		}
	}
	Evented struct {
		Url string
	}
	Database struct {
		Kind    string
		Mongodb struct {
			Url        string
			Name       string
			Collection string
		}
	}
	Domain string
	Saga   struct {
		Url string
	}
	Port                 uint
	OtherCommandHandlers []struct {
		Domain string
		Url    string
	}
}
