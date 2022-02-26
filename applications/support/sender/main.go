package main

import (
	"github.com/benjaminabbitt/evented/applications/support/sender/cmd"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	cmd.Execute()
}
