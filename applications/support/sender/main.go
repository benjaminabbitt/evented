package main

import (
	"github.com/benjaminabbitt/evented/applications/support/sender/actions/root"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	root.Execute()
}
