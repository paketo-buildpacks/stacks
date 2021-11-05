package stack_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestStack(t *testing.T) {
	suite := spec.New("stack", spec.Report(report.Terminal{}))
	suite("Creator", testCreator)
	suite("Mixins", testMixins)
	suite.Run(t)
}
