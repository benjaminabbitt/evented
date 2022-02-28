package configuration

import "github.com/benjaminabbitt/evented/support"

type Configuration struct {
	support.DefaultConfigInit
	Name  string
	Port  uint
	Proof string
}
