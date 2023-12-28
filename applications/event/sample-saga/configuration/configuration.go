package configuration

import (
	"github.com/benjaminabbitt/evented/support"
)

type Configuration struct {
	support.BasicConfigInit
	Port uint
}
