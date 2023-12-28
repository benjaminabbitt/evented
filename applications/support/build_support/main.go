package main

import (
	_ "github.com/benjaminabbitt/evented/applications/support/build_support/actions/build_time"
	"github.com/benjaminabbitt/evented/applications/support/build_support/actions/root"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	err := root.Execute()
	if err != nil {
		return
	}
}
