package send

import (
	"github.com/benjaminabbitt/evented/applications/todo/commands/send/commands/create"
	"github.com/cucumber/godog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	opts := godog.Options{
		Format: format,
		Paths:  []string{"features"},
	}

	status := godog.TestSuite{
		Name:                "create",
		ScenarioInitializer: create.InitializeScenario,
		Options:             &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
