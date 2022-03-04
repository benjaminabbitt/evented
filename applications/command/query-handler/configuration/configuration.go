package configuration

import (
	"github.com/benjaminabbitt/evented/support"
)

type Configuration struct {
	support.BasicConfigInit
	EventStore struct {
		Kind    string
		Mongodb struct {
			Url        string
			Name       string
			Collection string
		}
	}
	Port       uint
	TargetSize uint
	Proof      string
}
