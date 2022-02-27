package configuration

import "github.com/benjaminabbitt/evented/support"

type Configuration struct {
	support.ConfigInitS
	Name  string
	Port  uint
	Proof string
}
