package image_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/image"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testImageClient(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		imageClient image.Client
		tag         = "stack-test-image"
	)

	it.Before(func() {
		localTag, err := name.NewTag(tag)
		Expect(err).NotTo(HaveOccurred())

		image, err := random.Image(1, 1)
		Expect(err).NotTo(HaveOccurred())

		_, err = daemon.Write(localTag, image)
		Expect(err).NotTo(HaveOccurred())
	})

	it("can set labels", func() {
		err := imageClient.SetLabel(tag, "some-key", "some-value")
		Expect(err).NotTo(HaveOccurred())

		labels := getLabels(tag, t)

		Expect(labels["some-key"]).To(Equal("some-value"))
	})

	it("can build images", func() {
		dir, err := ioutil.TempDir("", "dockerfile-test")
		Expect(err).NotTo(HaveOccurred())

		file, err := os.Create(fmt.Sprintf("%s/%s", dir, "Dockerfile"))
		Expect(err).NotTo(HaveOccurred())

		_, err = file.WriteString(`FROM alpine
ARG test_build_arg
LABEL testing.key=some-value
LABEL testing.build.arg.key=$test_build_arg`)
		Expect(err).NotTo(HaveOccurred())

		defer os.RemoveAll(dir)

		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		err = imageClient.Build(tag, dir, false, nil, "test_build_arg=1")
		Expect(err).NotTo(HaveOccurred())

		labels := getLabels(tag, t)
		Expect(labels["testing.key"]).To(Equal("some-value"))
		Expect(labels["testing.build.arg.key"]).To(Equal("1"))
	})

	it("can build with docker buildkit", func() {
		dir, err := ioutil.TempDir("", "dockerfile-test")
		Expect(err).NotTo(HaveOccurred())

		file, err := os.Create(fmt.Sprintf("%s/%s", dir, "Dockerfile"))
		Expect(err).NotTo(HaveOccurred())

		_, err = file.WriteString(`FROM alpine`)
		Expect(err).NotTo(HaveOccurred())

		defer os.RemoveAll(dir)

		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		err = imageClient.Build(tag, dir, true, nil, "test_build_arg=1")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.Getenv("DOCKER_BUILDKIT")).To(Equal("1"))
	})

	it("can pass secrets to docker build command", func() {
		dir, err := ioutil.TempDir("", "dockerfile-test")
		Expect(err).NotTo(HaveOccurred())

		file, err := os.Create(fmt.Sprintf("%s/%s", dir, "Dockerfile"))
		Expect(err).NotTo(HaveOccurred())

		_, err = file.WriteString(`# syntax=docker/dockerfile:experimental
FROM alpine
RUN --mount=type=secret,id=test-secret,dst=/temp cat /temp > /secret`)
		Expect(err).NotTo(HaveOccurred())

		defer os.RemoveAll(dir)

		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		err = imageClient.Build(tag, dir, true, map[string]string{"test-secret": "some-secret"})
		Expect(err).NotTo(HaveOccurred())

		contents, err := exec.Command("docker", "run", tag, "cat", "/secret").CombinedOutput()
		Expect(err).NotTo(HaveOccurred())

		Expect(string(contents)).To(Equal("some-secret"))
	})
}
