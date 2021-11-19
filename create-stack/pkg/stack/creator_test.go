package stack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

type imageBuildInvocation struct {
	Tag            string
	DockerfilePath string
	WithBuildKit   bool
	Secrets        map[string]string
	BuildArgs      []string
}

func testCreator(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		fakePackageFinder     *fakes.PackageFinder
		fakeMixinsGenerator   *fakes.MixinsGenerator
		fakeImageClient       *fakes.ImageClient
		fakeBOMGenerator      *fakes.BOMGenerator
		imageBuildInvocations []imageBuildInvocation
		creator               stack.Creator
		dir                   string
		bomPath               string
		err                   error
	)

	it.Before(func() {
		fakePackageFinder = &fakes.PackageFinder{}
		fakeMixinsGenerator = &fakes.MixinsGenerator{}
		fakeBOMGenerator = &fakes.BOMGenerator{}

		imageBuildInvocations = []imageBuildInvocation{}
		fakeImageClient = &fakes.ImageClient{}
		fakeImageClient.BuildCall.Stub = func(tag, dockerfilePath string, withBuildKit bool, secrets map[string]string, buildArgs ...string) error {
			imageBuildInvocations = append(imageBuildInvocations, imageBuildInvocation{
				Tag:            tag,
				DockerfilePath: dockerfilePath,
				WithBuildKit:   withBuildKit,
				Secrets:        secrets,
				BuildArgs:      buildArgs,
			})
			return nil
		}

		dir, err = os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())
		bomPath = filepath.Join(dir, "syft.json")

		creator = stack.Creator{
			PackageFinder:   fakePackageFinder,
			MixinsGenerator: fakeMixinsGenerator,
			ImageClient:     fakeImageClient,
			BOMGenerator:    fakeBOMGenerator,
		}
	})

	it.After(func() {
		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	it("create the bionic stack", func() {
		fakeMixinsGenerator.GetMixinsCall.Returns.BuildMixins = []string{"test1", "test2", "build:test3"}
		fakeMixinsGenerator.GetMixinsCall.Returns.RunMixins = []string{"test1", "test2", "run:test4"}
		fakeBOMGenerator.GenerateCall.Returns.OutputPaths = []string{bomPath}

		err := creator.Execute(stack.Definition{
			BuildBase: stack.Image{
				Tag:        "test-build-base-tag",
				Dockerfile: "test-base-build-dockerfile-path",
				Args:       []string{"sources=test-sources", "packages=test-build-packages"},
			},
			BuildCNB: stack.Image{
				Tag:         "test-build-base-tag-cnb",
				Dockerfile:  "test-cnb-build-dockerfile-path",
				Description: "test-build-description",
			},
			RunBase: stack.Image{
				Tag:        "test-run-base-tag",
				Dockerfile: "test-base-run-dockerfile-path",
				Args:       []string{"sources=test-sources", "packages=test-run-packages"},
			},
			RunCNB: stack.Image{
				Tag:         "test-run-base-tag-cnb",
				Dockerfile:  "test-cnb-run-dockerfile-path",
				Description: "test-run-description",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeImageClient.PullCall.Receives.Tag).To(Equal("ubuntu:bionic"))
		Expect(fakeImageClient.PullCall.Receives.Keychain).To(Equal(authn.DefaultKeychain))

		Expect(fakeImageClient.BuildCall.CallCount).To(Equal(4))

		Expect(imageBuildInvocations[0].Tag).To(Equal("test-build-base-tag"))
		Expect(imageBuildInvocations[0].DockerfilePath).To(Equal("test-base-build-dockerfile-path"))
		Expect(imageBuildInvocations[0].BuildArgs).To(Equal([]string{"sources=test-sources", "packages=test-build-packages"}))

		Expect(imageBuildInvocations[1].Tag).To(Equal("test-run-base-tag"))
		Expect(imageBuildInvocations[1].DockerfilePath).To(Equal("test-base-run-dockerfile-path"))
		Expect(imageBuildInvocations[1].BuildArgs).To(Equal([]string{"sources=test-sources", "packages=test-run-packages"}))

		Expect(imageBuildInvocations[2].Tag).To(Equal("test-build-base-tag-cnb"))
		Expect(imageBuildInvocations[2].DockerfilePath).To(Equal("test-cnb-build-dockerfile-path"))
		Expect(imageBuildInvocations[2].BuildArgs[0]).To(Equal("base_image=test-build-base-tag"))
		Expect(imageBuildInvocations[2].BuildArgs[1]).To(Equal("description=test-build-description"))
		Expect(imageBuildInvocations[2].BuildArgs[2]).To(Equal("mixins=[\"test1\",\"test2\",\"build:test3\"]"))
		Expect(imageBuildInvocations[2].BuildArgs[4]).To(Equal("metadata={}"))

		Expect(imageBuildInvocations[3].Tag).To(Equal("test-run-base-tag-cnb"))
		Expect(imageBuildInvocations[3].DockerfilePath).To(Equal("test-cnb-run-dockerfile-path"))
		Expect(imageBuildInvocations[3].BuildArgs[0]).To(Equal("base_image=test-run-base-tag"))
		Expect(imageBuildInvocations[3].BuildArgs[1]).To(Equal("description=test-run-description"))
		Expect(imageBuildInvocations[3].BuildArgs[2]).To(Equal("mixins=[\"test1\",\"test2\",\"run:test4\"]"))
		Expect(imageBuildInvocations[3].BuildArgs[4]).To(Equal("metadata={}"))

		buildReleaseDate := imageBuildInvocations[2].BuildArgs[3]
		runReleaseDate := imageBuildInvocations[3].BuildArgs[3]
		Expect(runReleaseDate).To(Equal(buildReleaseDate))

		Expect(fakeImageClient.PushCall.CallCount).To(Equal(0))
		Expect(fakeBOMGenerator.GenerateCall.Receives.ImageTag).To(Equal("test-run-base-tag"))
		Expect(fakeBOMGenerator.AttachCall.Receives.CnbImageTag).To(Equal("test-run-base-tag-cnb"))
		Expect(fakeBOMGenerator.AttachCall.Receives.Files).To(Equal([]string{bomPath}))
	})

	it("create the bionic stack and publish", func() {
		fakeImageClient.PushCall.Stub = func(tag string) (string, error) {
			refs := []string{
				"test-build-base-tag@sha256:124124214hjshfskahfkjh12312",
				"test-run-base-tag@sha256:8098908s908f90asf8980989898",
			}
			if fakeImageClient.PushCall.CallCount <= len(refs) {
				return refs[fakeImageClient.PushCall.CallCount-1], nil
			}

			return "", nil
		}

		fakeMixinsGenerator.GetMixinsCall.Returns.BuildMixins = []string{"test1", "test2", "build:test3"}
		fakeMixinsGenerator.GetMixinsCall.Returns.RunMixins = []string{"test1", "test2", "run:test4"}
		fakeBOMGenerator.GenerateCall.Returns.OutputPaths = []string{bomPath}

		err := creator.Execute(stack.Definition{
			BuildBase: stack.Image{
				Publish:    true,
				Tag:        "test-build-base-tag:latest-test-stack",
				Dockerfile: "test-base-build-dockerfile-path",
				Args:       []string{"sources=test-sources", "packages=test-build-packages"},
			},
			BuildCNB: stack.Image{
				Publish:     true,
				Tag:         "test-build-base-tag:latest-test-stack-cnb",
				Dockerfile:  "test-cnb-build-dockerfile-path",
				Description: "test-build-description",
			},
			RunBase: stack.Image{
				Publish:    true,
				Tag:        "test-run-base-tag:latest-test-stack",
				Dockerfile: "test-base-run-dockerfile-path",
				Args:       []string{"sources=test-sources", "packages=test-run-packages"},
			},
			RunCNB: stack.Image{
				Publish:     true,
				Tag:         "test-run-base-tag:latest-test-stack-cnb",
				Dockerfile:  "test-cnb-run-dockerfile-path",
				Description: "test-run-description",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeImageClient.PullCall.Receives.Tag).To(Equal("ubuntu:bionic"))
		Expect(fakeImageClient.PullCall.Receives.Keychain).To(Equal(authn.DefaultKeychain))

		Expect(fakeImageClient.BuildCall.CallCount).To(Equal(4))

		Expect(imageBuildInvocations[0].Tag).To(Equal("test-build-base-tag:latest-test-stack"))
		Expect(imageBuildInvocations[0].DockerfilePath).To(Equal("test-base-build-dockerfile-path"))
		Expect(imageBuildInvocations[0].BuildArgs).To(Equal([]string{"sources=test-sources", "packages=test-build-packages"}))

		Expect(imageBuildInvocations[1].Tag).To(Equal("test-run-base-tag:latest-test-stack"))
		Expect(imageBuildInvocations[1].DockerfilePath).To(Equal("test-base-run-dockerfile-path"))
		Expect(imageBuildInvocations[1].BuildArgs).To(Equal([]string{"sources=test-sources", "packages=test-run-packages"}))

		Expect(imageBuildInvocations[2].Tag).To(Equal("test-build-base-tag:latest-test-stack-cnb"))
		Expect(imageBuildInvocations[2].DockerfilePath).To(Equal("test-cnb-build-dockerfile-path"))
		Expect(imageBuildInvocations[2].BuildArgs[0]).To(Equal("base_image=test-build-base-tag:latest-test-stack"))
		Expect(imageBuildInvocations[2].BuildArgs[1]).To(Equal("description=test-build-description"))
		Expect(imageBuildInvocations[2].BuildArgs[2]).To(Equal("mixins=[\"test1\",\"test2\",\"build:test3\"]"))
		Expect(imageBuildInvocations[2].BuildArgs[4]).To(Equal("metadata={\"base-image\":\"test-build-base-tag@sha256:124124214hjshfskahfkjh12312\"}"))

		Expect(imageBuildInvocations[3].Tag).To(Equal("test-run-base-tag:latest-test-stack-cnb"))
		Expect(imageBuildInvocations[3].DockerfilePath).To(Equal("test-cnb-run-dockerfile-path"))
		Expect(imageBuildInvocations[3].BuildArgs[0]).To(Equal("base_image=test-run-base-tag:latest-test-stack"))
		Expect(imageBuildInvocations[3].BuildArgs[1]).To(Equal("description=test-run-description"))
		Expect(imageBuildInvocations[3].BuildArgs[2]).To(Equal("mixins=[\"test1\",\"test2\",\"run:test4\"]"))
		Expect(imageBuildInvocations[3].BuildArgs[4]).To(Equal("metadata={\"base-image\":\"test-run-base-tag@sha256:8098908s908f90asf8980989898\"}"))

		buildReleaseDate := imageBuildInvocations[2].BuildArgs[3]
		runReleaseDate := imageBuildInvocations[3].BuildArgs[3]
		Expect(runReleaseDate).To(Equal(buildReleaseDate))
		Expect(fakeBOMGenerator.GenerateCall.Receives.ImageTag).To(Equal("test-run-base-tag:latest-test-stack"))
		Expect(fakeBOMGenerator.AttachCall.Receives.CnbImageTag).To(Equal("test-run-base-tag:latest-test-stack-cnb"))
		Expect(fakeBOMGenerator.AttachCall.Receives.Files).To(Equal([]string{bomPath}))
	})

	it("passes additional args when building CNB image", func() {
		err := creator.Execute(stack.Definition{
			BuildBase: stack.Image{
				Publish: true,
				Tag:     "test-build-base-tag:latest-test-stack",
			},
			BuildCNB: stack.Image{
				Publish: true,
				Tag:     "test-build-base-tag:latest-test-stack-cnb",
				Args:    []string{"stack_id=some.stack.id"},
			},
			RunBase: stack.Image{
				Publish: true,
				Tag:     "test-run-base-tag:latest-test-stack",
			},
			RunCNB: stack.Image{
				Publish: true,
				Tag:     "test-run-base-tag:latest-test-stack-cnb",
				Args:    []string{"stack_id=some.stack.id"},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(imageBuildInvocations[2].BuildArgs).To(ContainElement("stack_id=some.stack.id"))
		Expect(imageBuildInvocations[3].BuildArgs).To(ContainElement("stack_id=some.stack.id"))
	})
}
