package acceptance_test

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/ulikunitz/xz"
)

func TestAcceptance(t *testing.T) {
	var Expect = NewWithT(t).Expect

	tempFile, err := os.CreateTemp("", "create-stack")
	Expect(err).NotTo(HaveOccurred())

	cliPath := tempFile.Name()
	Expect(tempFile.Close()).To(Succeed())

	goBuild := exec.Command("go", "build", "-o", cliPath, ".")

	stacksDir, err := getStacksDirectory()
	Expect(err).NotTo(HaveOccurred())

	goBuild.Dir = stacksDir + "/create-stack"

	output, err := goBuild.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to build CLI: %s", output))

	suite := spec.New("acceptance", spec.Report(report.Terminal{}))
	suite("Base", testBase(cliPath))
	suite("Full", testFull(cliPath))
	suite("Publish", testPublish(cliPath))
	suite("Tiny", testTiny(cliPath))
	suite.Run(t)

	Expect(os.Remove(cliPath)).To(Succeed())
}

const (
	DistroName    = "Ubuntu"
	DistroVersion = "18.04"
	Homepage      = "https://github.com/paketo-buildpacks/stacks"
	BionicStackID = "io.buildpacks.stacks.bionic"
	TinyStackID   = "io.paketo.stacks.tiny"
	Maintainer    = "Paketo Buildpacks"
)

func assertCommonLabels(t *testing.T, stackID string, imageConfig ImageConfig) {
	t.Helper()
	var Expect = NewWithT(t).Expect

	Expect(imageConfig.StackLabels.DistroName).To(Equal(DistroName))
	Expect(imageConfig.StackLabels.DistroVersion).To(Equal(DistroVersion))
	Expect(imageConfig.StackLabels.Homepage).To(Equal(Homepage))
	Expect(imageConfig.StackLabels.ID).To(Equal(stackID))
	Expect(imageConfig.StackLabels.Maintainer).To(Equal(Maintainer))
}

func getStacksDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to obtain directory")
	}
	filename = filepath.Dir(filename)
	return strings.TrimSuffix(filename, "create-stack/test/acceptance"), nil
}

type ImageConfig struct {
	User        string      `json:"User"`
	Env         []string    `json:"Env"`
	StackLabels StackLabels `json:"Labels"`
}

type ImageMetadata struct {
	ImageConfig ImageConfig `json:"Config"`
	RootFS      struct {
		Layers []string `json:"Layers"`
	} `json:"RootFS"`
}

type StackLabels struct {
	Description   string `json:"io.buildpacks.stack.description"`
	DistroName    string `json:"io.buildpacks.stack.distro.name"`
	DistroVersion string `json:"io.buildpacks.stack.distro.version"`
	Homepage      string `json:"io.buildpacks.stack.homepage"`
	ID            string `json:"io.buildpacks.stack.id"`
	Maintainer    string `json:"io.buildpacks.stack.maintainer"`
	Metadata      string `json:"io.buildpacks.stack.metadata"`
	Mixins        string `json:"io.buildpacks.stack.mixins"`
	Packages      string `json:"io.paketo.stack.packages"`
	Released      string `json:"io.buildpacks.stack.released"`
	SBOM          string `json:"io.buildpacks.base.sbom"`
}

func assertCorrectBaseImage(t *testing.T, baseImageRef string) {
	t.Helper()

	var Expect = NewWithT(t).Expect

	cnbImageRef := baseImageRef + "-cnb"
	output, err := exec.Command("docker", "inspect", cnbImageRef, "--format", "{{json .Config}}").CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	var runImageConfig ImageConfig
	err = json.Unmarshal(output, &runImageConfig)
	Expect(err).NotTo(HaveOccurred())

	output, err = exec.Command("docker", "pull", baseImageRef).CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", baseImageRef).CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), string(output))

	Expect(runImageConfig.StackLabels.Metadata).To(MatchJSON(fmt.Sprintf(`{"base-image": %q}`, strings.TrimSpace(string(output)))))
}

func assertContainerFileExists(t *testing.T, containerID string, path string) {
	t.Helper()
	var Expect = NewWithT(t).Expect

	output, err := exec.Command(
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", containerID, path),
		"-",
	).CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("expected %s to exist:\n%s", path, output))
}

func assertContainerFileDoesNotExist(t *testing.T, containerID string, path string) {
	t.Helper()
	var Expect = NewWithT(t).Expect

	output, err := exec.Command(
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", containerID, path),
		"-",
	).CombinedOutput()
	Expect(err).To(HaveOccurred(), fmt.Sprintf("expected %s not to exist:\n%s", path, output))
}

func assertSBOMAttached(t *testing.T, imageRef string, labels StackLabels) {
	t.Helper()
	var Expect = NewWithT(t).Expect

	ref, err := name.ParseReference(imageRef)
	Expect(err).NotTo(HaveOccurred())
	img, err := daemon.Image(ref)
	Expect(err).NotTo(HaveOccurred())

	containsBOMs, err := hasBOMs(labels.SBOM, img)
	Expect(err).NotTo(HaveOccurred())
	Expect(containsBOMs).To(BeTrue())
}

func assertOSReleaseEqual(t *testing.T, expectedOSReleaseContent, actualOSReleaseContent string) {
	t.Helper()
	var Expect = NewWithT(t).Expect

	var expectedKVPairs []string
	for _, kvPair := range strings.Split(strings.TrimSpace(expectedOSReleaseContent), "\n") {
		if shouldIgnoreOSReleaseKey(kvPair) {
			continue
		}

		expectedKVPairs = append(expectedKVPairs, kvPair)
	}

	var actualKVPairs []string
	for _, kvPair := range strings.Split(strings.TrimSpace(actualOSReleaseContent), "\n") {
		if shouldIgnoreOSReleaseKey(kvPair) {
			continue
		}

		actualKVPairs = append(actualKVPairs, kvPair)
	}

	sort.Strings(expectedKVPairs)
	sort.Strings(actualKVPairs)

	Expect(expectedKVPairs).To(Equal(actualKVPairs))
}

func getContainerFile(containerID string, path string) (string, error) {
	tempFile, err := os.CreateTemp("", "create-stack-tiny-acceptance")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)

	err = tempFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	output, err := exec.Command(
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", containerID, path),
		tempFilePath,
	).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to copy file from container: %w\n%s", err, string(output))
	}

	contents, err := os.ReadFile(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return string(contents), nil
}

func getOriginalOSRelease(containerID string) (string, error) {
	baseFilesPackageMetadata, err := getContainerFile(containerID, "/var/lib/dpkg/status.d/base-files")
	if err != nil {
		return "", fmt.Errorf("failed to get base-files metadata: %w", err)
	}

	re := regexp.MustCompile("(?m)^Version: (.*)")
	matches := re.FindStringSubmatch(baseFilesPackageMetadata)
	if len(matches) != 2 {
		return "", fmt.Errorf("expected 2 matches but got %d: %v", len(matches), matches)
	}

	baseFilesVersion := matches[1]

	packageURL := fmt.Sprintf("http://archive.ubuntu.com/ubuntu/pool/main/b/base-files/base-files_%s.tar.xz", baseFilesVersion)
	resp, err := http.Get(packageURL)
	if err != nil {
		return "", fmt.Errorf("failed to get base-files package: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("got unsuccessful status code: %s", resp.Status)
	}

	xzReader, err := xz.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to create xz reader: %w", err)
	}

	tarReader := tar.NewReader(xzReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to get next item from tar reader: %w", err)
		}

		if strings.HasSuffix(hdr.Name, "etc/os-release") {
			contents, err := io.ReadAll(tarReader)
			if err != nil {
				return "", fmt.Errorf("failed to read os-release from tar file: %w", err)
			}

			return string(contents), nil
		}
	}
	return "", fmt.Errorf("could not find os-release file in base-files package")
}

func shouldIgnoreOSReleaseKey(kvPair string) bool {
	return strings.HasPrefix(kvPair, "PRETTY_NAME") ||
		strings.HasPrefix(kvPair, "HOME_URL") ||
		strings.HasPrefix(kvPair, "SUPPORT_URL") ||
		strings.HasPrefix(kvPair, "BUG_REPORT_URL")
}

func hasBOMs(layerDiffID string, img v1.Image) (bool, error) {
	var seenSyft bool
	var seenCyclonedx bool

	diffID, err := v1.NewHash(layerDiffID)
	if err != nil {
		return false, err
	}

	bomLayer, err := img.LayerByDiffID(diffID)
	if err != nil {
		return false, err
	}

	layerReader, err := bomLayer.Uncompressed()
	if err != nil {
		return false, err
	}
	tr := tar.NewReader(layerReader)
	for {
		header, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return false, err
			}
			break
		}

		if header.Typeflag != tar.TypeReg || !strings.HasPrefix(header.Name, "/cnb/sbom") {
			continue
		}

		switch strings.TrimPrefix(header.Name, "/cnb/sbom/") {
		case "bom.syft.json":
			seenSyft = true
		case "bom.cdx.json":
			seenCyclonedx = true
		}
	}
	return seenSyft && seenCyclonedx, nil
}
