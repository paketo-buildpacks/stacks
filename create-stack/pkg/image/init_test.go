package image_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

func TestImage(t *testing.T) {
	suite := spec.New("image", spec.Report(report.Terminal{}))
	suite("ImageClient", testImageClient)
	suite.Run(t)
}

func getLabels(tag string, t *testing.T) map[string]string {
	t.Helper()
	var Expect = NewWithT(t).Expect

	ref, err := name.ParseReference(tag)
	Expect(err).NotTo(HaveOccurred())

	image, err := daemon.Image(ref)
	Expect(err).NotTo(HaveOccurred())

	configFile, err := image.ConfigFile()
	Expect(err).NotTo(HaveOccurred())

	return configFile.Config.Labels
}
