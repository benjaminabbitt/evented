package memory

import (
	"github.com/benjaminabbitt/evented/support/cucumber"
	"github.com/cucumber/godog"
	"os"
	"reflect"
	"testing"
)

func init() {
	cucumber.Init(cucumber.GetFlagOptions())
}

type Empty struct{}

func TestMain(m *testing.M) {
	os.Exit(executeCucumber())
}

func Test_Cucumber(t *testing.T) {
	result := executeCucumber()
	t.Logf("Cucumber tests executed with status %d", result)
	if result != 0 {
		t.Fail()
	}
}

func executeCucumber() int {
	format := cucumber.GetFormat()
	opts := cucumber.GetOptions(format)
	opts.Paths = []string{"../"}
	suite := godog.TestSuite{
		Name:                 reflect.TypeOf(Empty{}).PkgPath(),
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}
	return cucumber.RunTestsWithCucumber(suite, opts)
}
