package cucumber

import (
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func RunTestsWithCucumber(m *testing.M, suite godog.TestSuite, opts godog.Options) {
	status := suite.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func GetOptions(format string) godog.Options {
	opts := godog.Options{
		Format:    format,
		Paths:     []string{"."},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}
	return opts
}

func GetFormat() string {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}
	return format
}

func GetFlagOptions() godog.Options {
	var opts = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "progress", // can define default values
	}
	return opts
}

func Init(opts godog.Options) {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

// assertExpectedAndActual is a helper function to allow the step function to call
// assertion functions where you want to compare an expected and an actual value.
func AssertExpectedAndActual(a ExpectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t Asserter
	a(&t, expected, actual, msgAndArgs...)
	return t.err
}

type ExpectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

// assertActual is a helper function to allow the step function to call
// assertion functions where you want to compare an actual value to a
// predined state like nil, empty or true/false.
func AssertActual(a ActualAssertion, actual interface{}, msgAndArgs ...interface{}) error {
	var t Asserter
	a(&t, actual, msgAndArgs...)
	return t.err
}

type ActualAssertion func(t assert.TestingT, actual interface{}, msgAndArgs ...interface{}) bool

// asserter is used to be able to retrieve the error reported by the called assertion
type Asserter struct {
	err error
}

// Errorf is used by the called assertion to report an error
func (a *Asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

func Uint64ToUint32WithErrorPassthrough(value uint64, err error) (uint32, error) {
	return uint32(value), err
}
