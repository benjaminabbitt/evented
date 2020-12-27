package memory

import (
	"github.com/benjaminabbitt/evented/support/cucumber"
	"github.com/cucumber/godog"
	"reflect"
	"testing"
)

func init() {
	cucumber.Init(cucumber.GetFlagOptions())
}

type Empty struct{}

func TestMain(m *testing.M) {
	format := cucumber.GetFormat()
	opts := cucumber.GetOptions(format)
	opts.Paths = []string{"../"}
	suite := godog.TestSuite{
		Name:                 reflect.TypeOf(Empty{}).PkgPath(),
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}
	cucumber.RunTestsWithCucumber(m, suite, opts)
}
