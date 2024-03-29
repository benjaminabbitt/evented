package configuration

import "github.com/benjaminabbitt/evented/support"

type Configuration struct {
	support.BasicConfigInit
	QueryHandler struct {
		Url  string
		Name string
	}
	Transport struct {
		Kind string
		AMQP struct {
			Url      string
			Exchange string
			Queue    string
		}
	}
	Projector struct {
		Url  string
		Name string
	}
	Database struct {
		Kind    string
		Mongodb struct {
			Url        string
			Name       string
			Collection string
		}
	}
	Port   uint
	Domain string
}
