package configuration

import "github.com/benjaminabbitt/evented/support"

type Configuration struct {
	support.BasicConfigInit
	Name  string
	Port  uint
	Proof string
}
