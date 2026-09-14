package test

import (
	"os"
	"testing"

	"github.com/RogerGCruz/pos-goexpert-desafio5/test/steps"

	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                "loadtest",
		ScenarioInitializer: steps.InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
		},
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
