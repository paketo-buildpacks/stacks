package stack_test

import (
	"testing"

	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testMixins(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		mixins stack.Mixins
	)

	it("can get the mixins", func() {
		buildPackages := []string{"testlib1", "testlib2", "testlib3", "testlib4", "testlib5", "testlib6"}
		runPackages := []string{"testlib1", "testlib5", "testlib6", "testlib7"}

		buildMixins, runMixins := mixins.GetMixins(buildPackages, runPackages)
		Expect(buildMixins).To(Equal([]string{"testlib1", "testlib5", "testlib6", "build:testlib2", "build:testlib3", "build:testlib4"}))
		Expect(runMixins).To(Equal([]string{"testlib1", "testlib5", "testlib6", "run:testlib7"}))
	})
}
