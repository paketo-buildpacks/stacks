package packages_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestPackages(t *testing.T) {
	suite := spec.New("packages", spec.Report(report.Terminal{}))
	suite("Bionic", testBionic)
	suite("Tiny", testTiny)
	suite.Run(t)
}
