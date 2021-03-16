package stack_test

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack/stackfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
	"testing"
)

func TestCreator(t *testing.T) {
	spec.Run(t, "Creator", testCreator, spec.Report(report.Terminal{}))
}

func testCreator(t *testing.T, when spec.G, it spec.S) {

	var (
		fakePackageFinder   = &stackfakes.FakePackageFinder{}
		fakeMixinsGenerator = &stackfakes.FakeMixinsGenerator{}
		fakeImageClient     = &stackfakes.FakeImageClient{}
		fakeStack           = &stackfakes.FakeStack{}
		assert              = assertpkg.New(t)
		require             = requirepkg.New(t)
		creator             stack.Creator
	)

	it.Before(func() {
		creator = stack.Creator{
			PackageFinder:   fakePackageFinder,
			MixinsGenerator: fakeMixinsGenerator,
			ImageClient:     fakeImageClient,
		}
	})

	it("create the bionic stack", func() {
		fakeStack.GetNameReturns("test-stack")
		fakeStack.GetBaseBuildArgsReturns([]string{"sources=test-sources", "packages=test-build-packages"})
		fakeStack.GetBaseRunArgsReturns([]string{"sources=test-sources", "packages=test-run-packages"})
		fakeStack.GetBaseBuildDockerfilePathReturns("test-base-build-dockerfile-path")
		fakeStack.GetBaseRunDockerfilePathReturns("test-base-run-dockerfile-path")
		fakeStack.GetCNBBuildDockerfilePathReturns("test-cnb-build-dockerfile-path")
		fakeStack.GetCNBRunDockerfilePathReturns("test-cnb-run-dockerfile-path")
		fakeStack.GetBuildDescriptionReturns("test-build-description")
		fakeStack.GetRunDescriptionReturns("test-run-description")

		fakeMixinsGenerator.GetMixinsReturns([]string{"test1", "test2", "build:test3"}, []string{"test1", "test2", "run:test4"})

		err := creator.CreateStack(fakeStack, "test-build-base-tag", "test-run-base-tag", false)
		require.NoError(err)

		imagePullTag, pullAuth := fakeImageClient.PullArgsForCall(0)
		assert.Equal("ubuntu:bionic", imagePullTag)
		assert.Equal(authn.DefaultKeychain, pullAuth)

		assert.Equal(4, fakeImageClient.BuildCallCount())

		expectedBaseBuildTag := "test-build-base-tag"
		expectedBaseBuildDockerfile := "test-base-build-dockerfile-path"
		expectedBaseBuildArgs := []string{"sources=test-sources", "packages=test-build-packages"}

		actualBaseBuildTag, actualBaseBuildDockerfile, actualBaseBuildArgs := fakeImageClient.BuildArgsForCall(0)
		assert.Equal(expectedBaseBuildTag, actualBaseBuildTag)
		assert.Equal(actualBaseBuildDockerfile, expectedBaseBuildDockerfile)
		assert.Equal(actualBaseBuildArgs, expectedBaseBuildArgs)

		expectedBaseRunTag := "test-run-base-tag"
		expectedBaseRunDockerfile := "test-base-run-dockerfile-path"
		expectedBaseRunArgs := []string{"sources=test-sources", "packages=test-run-packages"}

		actualBaseRunTag, actualBaseRunDockerfile, actualBaseRunArgs := fakeImageClient.BuildArgsForCall(1)
		assert.Equal(expectedBaseRunTag, actualBaseRunTag)
		assert.Equal(expectedBaseRunDockerfile, actualBaseRunDockerfile)
		assert.Equal(expectedBaseRunArgs, actualBaseRunArgs)

		expectedCNBBuildTag := "test-build-base-tag-cnb"
		expectedCNBBuildDockerfile := "test-cnb-build-dockerfile-path"

		actualCNBBuildTag, actualCNBBuildDockerfile, actualCNBBuildArgs := fakeImageClient.BuildArgsForCall(2)
		assert.Equal(expectedCNBBuildTag, actualCNBBuildTag)
		assert.Equal(actualCNBBuildDockerfile, expectedCNBBuildDockerfile)
		assert.Equal("base_image=test-build-base-tag", actualCNBBuildArgs[0])
		assert.Equal("description=test-build-description", actualCNBBuildArgs[1])
		assert.Equal("mixins=[\"test1\",\"test2\",\"build:test3\"]", actualCNBBuildArgs[2])
		assert.Equal("metadata={}", actualCNBBuildArgs[4])

		buildReleaseDate := actualCNBBuildArgs[3]

		expectedCNBRunTag := "test-run-base-tag-cnb"
		expectedCNBRunDockerfile := "test-cnb-run-dockerfile-path"

		actualCNBRunTag, actualCNBRunDockerfile, actualCNBRunArgs := fakeImageClient.BuildArgsForCall(3)
		assert.Equal(expectedCNBRunTag, actualCNBRunTag)
		assert.Equal(actualCNBRunDockerfile, expectedCNBRunDockerfile)
		assert.Equal("base_image=test-run-base-tag", actualCNBRunArgs[0])
		assert.Equal("description=test-run-description", actualCNBRunArgs[1])
		assert.Equal("mixins=[\"test1\",\"test2\",\"run:test4\"]", actualCNBRunArgs[2])
		assert.Equal("metadata={}", actualCNBRunArgs[4])

		runReleaseDate := actualCNBRunArgs[3]

		assert.Equal(buildReleaseDate, runReleaseDate)

		assert.Equal(0, fakeImageClient.PushCallCount())
	})

	it("create the bionic stack and publish", func() {
		fakeStack.GetNameReturns("test-stack")
		fakeStack.GetBaseBuildArgsReturns([]string{"sources=test-sources", "packages=test-build-packages"})
		fakeStack.GetBaseRunArgsReturns([]string{"sources=test-sources", "packages=test-run-packages"})
		fakeStack.GetBaseBuildDockerfilePathReturns("test-base-build-dockerfile-path")
		fakeStack.GetBaseRunDockerfilePathReturns("test-base-run-dockerfile-path")
		fakeStack.GetCNBBuildDockerfilePathReturns("test-cnb-build-dockerfile-path")
		fakeStack.GetCNBRunDockerfilePathReturns("test-cnb-run-dockerfile-path")
		fakeStack.GetBuildDescriptionReturns("test-build-description")
		fakeStack.GetRunDescriptionReturns("test-run-description")

		fakeImageClient.PushReturnsOnCall(0, "test-build-base-tag@sha256:124124214hjshfskahfkjh12312", nil)
		fakeImageClient.PushReturnsOnCall(1, "test-run-base-tag@sha256:8098908s908f90asf8980989898", nil)

		fakeMixinsGenerator.GetMixinsReturns([]string{"test1", "test2", "build:test3"}, []string{"test1", "test2", "run:test4"})

		err := creator.CreateStack(fakeStack, "test-build-base-tag:latest-test-stack", "test-run-base-tag:latest-test-stack", true)
		require.NoError(err)

		imagePullTag, pullAuth := fakeImageClient.PullArgsForCall(0)
		assert.Equal("ubuntu:bionic", imagePullTag)
		assert.Equal(authn.DefaultKeychain, pullAuth)

		assert.Equal(4, fakeImageClient.BuildCallCount())

		expectedBaseBuildTag := "test-build-base-tag:latest-test-stack"
		expectedBaseBuildDockerfile := "test-base-build-dockerfile-path"
		expectedBaseBuildArgs := []string{"sources=test-sources", "packages=test-build-packages"}

		actualBaseBuildTag, actualBaseBuildDockerfile, actualBaseBuildArgs := fakeImageClient.BuildArgsForCall(0)
		assert.Equal(expectedBaseBuildTag, actualBaseBuildTag)
		assert.Equal(actualBaseBuildDockerfile, expectedBaseBuildDockerfile)
		assert.Equal(actualBaseBuildArgs, expectedBaseBuildArgs)

		expectedBaseRunTag := "test-run-base-tag:latest-test-stack"
		expectedBaseRunDockerfile := "test-base-run-dockerfile-path"
		expectedBaseRunArgs := []string{"sources=test-sources", "packages=test-run-packages"}

		actualBaseRunTag, actualBaseRunDockerfile, actualBaseRunArgs := fakeImageClient.BuildArgsForCall(1)
		assert.Equal(expectedBaseRunTag, actualBaseRunTag)
		assert.Equal(expectedBaseRunDockerfile, actualBaseRunDockerfile)
		assert.Equal(expectedBaseRunArgs, actualBaseRunArgs)

		expectedCNBBuildTag := "test-build-base-tag:latest-test-stack-cnb"
		expectedCNBBuildDockerfile := "test-cnb-build-dockerfile-path"

		actualCNBBuildTag, actualCNBBuildDockerfile, actualCNBBuildArgs := fakeImageClient.BuildArgsForCall(2)
		assert.Equal(expectedCNBBuildTag, actualCNBBuildTag)
		assert.Equal(actualCNBBuildDockerfile, expectedCNBBuildDockerfile)
		assert.Equal("base_image=test-build-base-tag:latest-test-stack", actualCNBBuildArgs[0])
		assert.Equal("description=test-build-description", actualCNBBuildArgs[1])
		assert.Equal("mixins=[\"test1\",\"test2\",\"build:test3\"]", actualCNBBuildArgs[2])
		assert.Equal("metadata={\"base-image\":\"test-build-base-tag@sha256:124124214hjshfskahfkjh12312\"}", actualCNBBuildArgs[4])

		buildReleaseDate := actualCNBBuildArgs[3]

		expectedCNBRunTag := "test-run-base-tag:latest-test-stack-cnb"
		expectedCNBRunDockerfile := "test-cnb-run-dockerfile-path"

		actualCNBRunTag, actualCNBRunDockerfile, actualCNBRunArgs := fakeImageClient.BuildArgsForCall(3)
		assert.Equal(expectedCNBRunTag, actualCNBRunTag)
		assert.Equal(actualCNBRunDockerfile, expectedCNBRunDockerfile)
		assert.Equal("base_image=test-run-base-tag:latest-test-stack", actualCNBRunArgs[0])
		assert.Equal("description=test-run-description", actualCNBRunArgs[1])
		assert.Equal("mixins=[\"test1\",\"test2\",\"run:test4\"]", actualCNBRunArgs[2])
		assert.Equal("metadata={\"base-image\":\"test-run-base-tag@sha256:8098908s908f90asf8980989898\"}", actualCNBRunArgs[4])

		runReleaseDate := actualCNBRunArgs[3]

		assert.Equal(buildReleaseDate, runReleaseDate)
	})
}
