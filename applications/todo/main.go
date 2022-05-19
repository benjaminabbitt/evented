package main

import (
	"github.com/benjaminabbitt/evented/applications/todo/commands/root"
	"github.com/dsnet/try"
)

func main() {
	try.E(root.Execute())
}
