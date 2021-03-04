package stack_test

import (
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	"testing"
)

func TestMixins(t *testing.T) {
	spec.Run(t, "Mixins", testMixins, spec.Report(report.Terminal{}))
}

func testMixins(t *testing.T, when spec.G, it spec.S) {
	var (
		assert = assertpkg.New(t)
		mixins stack.Mixins
	)

	it.Before(func() {
		mixins = stack.Mixins{}
	})

	it("can get the mixins", func() {
		buildPackages := []string{"testlib1", "testlib2", "testlib3", "testlib4", "testlib5", "testlib6"}
		runPackages := []string{"testlib1", "testlib5", "testlib6", "testlib7"}

		expectedBuildMixins := []string{"testlib1", "testlib5", "testlib6", "build:testlib2", "build:testlib3", "build:testlib4"}
		expectedRunMixins := []string{"testlib1", "testlib5", "testlib6", "run:testlib7"}

		actualBuildMixins, actualRunMixins := mixins.GetMixins(buildPackages, runPackages)
		assert.Equal(expectedBuildMixins, actualBuildMixins)
		assert.Equal(expectedRunMixins, actualRunMixins)
	})
}
