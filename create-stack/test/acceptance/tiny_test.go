package acceptance_test

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
)

func TestCreateStackTiny(t *testing.T) {
	spec.Run(t, "CreateStackTiny", testCreateStackTiny, spec.Report(report.Terminal{}))
}

func testCreateStackTiny(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath     string
		containerID string
		require     = requirepkg.New(t)
		assert      = assertpkg.New(t)
	)

	it.Before(func() {
		tempFile, err := ioutil.TempFile("", "create-stack")
		require.NoError(err)

		cliPath = tempFile.Name()
		require.NoError(tempFile.Close())

		goBuild := exec.Command("go", "build", "-o", cliPath, ".")

		stacksDir, err := getStacksDirectory()
		require.NoError(err)

		goBuild.Dir = stacksDir + "/create-stack"

		output, err := goBuild.CombinedOutput()
		require.NoError(err, "failed to build CLI: %s", string(output))
	})

	it.After(func() {
		_ = os.Remove(cliPath)

		if containerID != "" {
			output, err := exec.Command("docker", "rm", "-f", containerID).CombinedOutput()
			require.NoError(err, "failed to remove container %s: %s", containerID, string(output))
		}
	})

	it("builds tiny stack", func() {
		buildRepo := "paketotesting/test-create-stack-tiny-build"
		runRepo := "paketotesting/test-create-stack-tiny-run"
		version := fmt.Sprintf("dev-%d", time.Now().UnixNano())

		stacksDir, err := getStacksDirectory()
		require.NoError(err)

		cmd := exec.Command(cliPath,
			"--build-destination", buildRepo,
			"--run-destination", runRepo,
			"--version", version,
			"--stack", "tiny",
			"--stacks-dir", stacksDir,
		)
		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		buildBaseImageRef := fmt.Sprintf("%s:%s-%s", buildRepo, version, "tiny")
		buildCNBImageRef := buildBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", buildCNBImageRef, "--format={{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var buildImageConfig ImageConfig
		err = json.Unmarshal(output, &buildImageConfig)
		require.NoError(err)

		assertCommonLabels(t, BionicStackID, buildImageConfig)

		buildDescription := "ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
		assert.Equal(buildDescription, buildImageConfig.StackLabels.Description)

		assert.JSONEq(`{}`, buildImageConfig.StackLabels.Metadata)

		assert.Contains(buildImageConfig.StackLabels.Mixins, `"build:make"`)
		assert.Contains(buildImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(buildImageConfig.StackLabels.Mixins, `"run:"`)

		buildReleaseDate, err := time.Parse(time.RFC3339, buildImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), buildReleaseDate, time.Minute*10)

		assert.Contains(buildImageConfig.StackLabels.Packages, `"ca-certificates"`)

		runBaseImageRef := fmt.Sprintf("%s:%s-%s", runRepo, version, "tiny")
		runCNBImageRef := runBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", runCNBImageRef, "--format", "{{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var runImageConfig ImageConfig
		err = json.Unmarshal(output, &runImageConfig)
		require.NoError(err)

		assertCommonLabels(t, TinyStackID, runImageConfig)

		runDescription := "distroless-like bionic + glibc + openssl + CA certs"
		assert.Equal(runDescription, runImageConfig.StackLabels.Description)

		assert.JSONEq(`{}`, runImageConfig.StackLabels.Metadata)

		assert.Contains(runImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.Contains(runImageConfig.StackLabels.Mixins, `"run:netbase"`)
		assert.NotContains(runImageConfig.StackLabels.Mixins, "build:")

		runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), runReleaseDate, time.Minute*10)
		assert.Equal(buildReleaseDate, runReleaseDate)

		assert.Contains(runImageConfig.StackLabels.Packages, `"ca-certificates"`)

		output, err = exec.Command("docker", "create", runBaseImageRef, "/tmp/app").CombinedOutput()
		require.NoError(err, "failed to create container: %s", string(output))
		containerID = strings.TrimSpace(string(output))

		assertContainerFileExists(t, containerID, "/usr/share/doc/ca-certificates/copyright")
		assertContainerFileExists(t, containerID, "/etc/ssl/certs/ca-certificates.crt")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/base-files")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/ca-certificates")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/libc6")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/libssl1.1")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/netbase")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/openssl")
		assertContainerFileExists(t, containerID, "/var/lib/dpkg/status.d/tzdata")
		assertContainerFileExists(t, containerID, "/root")
		assertContainerFileExists(t, containerID, "/home/nonroot")
		assertContainerFileExists(t, containerID, "/tmp")
		assertContainerFileExists(t, containerID, "/etc/services")
		assertContainerFileExists(t, containerID, "/etc/nsswitch.conf")
		assertContainerFileExists(t, containerID, "/etc/passwd")
		assertContainerFileExists(t, containerID, "/etc/os-release")
		assertContainerFileExists(t, containerID, "/etc/group")
		assertContainerFileDoesNotExist(t, containerID, "/usr/share/ca-certificates")

		passwdContent, err := getContainerFile(containerID, "/etc/passwd")
		require.NoError(err)
		assert.Equal(
			`root:x:0:0:root:/root:/sbin/nologin
nobody:x:65534:65534:nobody:/nonexistent:/sbin/nologin
nonroot:x:65532:65532:nonroot:/home/nonroot:/sbin/nologin
`, passwdContent)

		groupContent, err := getContainerFile(containerID, "/etc/group")
		require.NoError(err)
		assert.Equal(`root:x:0:
nobody:x:65534:
tty:x:5:
staff:x:50:
nonroot:x:65532:
`, groupContent)

		osReleaseContent, err := getContainerFile(containerID, "/etc/os-release")
		require.NoError(err)
		assert.Contains(osReleaseContent, `PRETTY_NAME="Cloud Foundry Tiny"`)
		assert.Contains(osReleaseContent, `HOME_URL="https://github.com/cloudfoundry/stacks"`)
		assert.Contains(osReleaseContent, `SUPPORT_URL="https://github.com/cloudfoundry/stacks/blob/master/README.md"`)
		assert.Contains(osReleaseContent, `BUG_REPORT_URL="https://github.com/cloudfoundry/stacks/issues/new"`)

		originalOSReleaseContent, err := getOriginalOSRelease(containerID)
		require.NoError(err)

		assertOSReleaseEqual(t, originalOSReleaseContent, osReleaseContent)

		testAppFile, err := ioutil.TempFile("", "create-stack-tiny-test-app")
		require.NoError(err)

		testAppPath := testAppFile.Name()

		err = testAppFile.Close()
		require.NoError(err)
		defer os.Remove(testAppPath)

		goBuild := exec.Command("go", "build", "-o", testAppPath, ".")
		goBuild.Env = append(os.Environ(), "GOOS=linux")
		goBuild.Dir = filepath.Join(stacksDir, "create-stack", "test", "acceptance", "fixtures", "tiny")
		output, err = goBuild.CombinedOutput()
		require.NoError(err, "failed to build test app: %s", string(output))

		output, err = exec.Command("docker", "cp", testAppPath, containerID+":/tmp/app").CombinedOutput()
		require.NoError(err, "failed to copy test app into container: %s", string(output))

		output, err = exec.Command("docker", "start", "-a", containerID).CombinedOutput()
		require.NoError(err, "failed to run test app in container: %s", string(output))
	})
}

func assertContainerFileExists(t *testing.T, containerID string, path string) {
	output, err := exec.Command(
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", containerID, path),
		"-",
	).CombinedOutput()
	requirepkg.NoError(t, err, "expected %s to exist:\n%s", path, string(output))
}

func assertContainerFileDoesNotExist(t *testing.T, containerID string, path string) {
	output, err := exec.Command(
		"docker",
		"cp",
		fmt.Sprintf("%s:%s", containerID, path),
		"-",
	).CombinedOutput()
	requirepkg.Error(t, err, "expected %s not to exist:\n%s", path, string(output))
}

func assertOSReleaseEqual(t *testing.T, expectedOSReleaseContent, actualOSReleaseContent string) {
	var expectedKVPairs, actualKVPairs []string

	for _, kvPair := range strings.Split(strings.TrimSpace(expectedOSReleaseContent), "\n") {
		if shouldIgnoreOSReleaseKey(kvPair) {
			continue
		}

		expectedKVPairs = append(expectedKVPairs, kvPair)
	}

	for _, kvPair := range strings.Split(strings.TrimSpace(actualOSReleaseContent), "\n") {
		if shouldIgnoreOSReleaseKey(kvPair) {
			continue
		}

		actualKVPairs = append(actualKVPairs, kvPair)
	}

	sort.Strings(expectedKVPairs)
	sort.Strings(actualKVPairs)

	assertpkg.Equal(t, expectedKVPairs, actualKVPairs)
}

func getContainerFile(containerID string, path string) (string, error) {
	tempFile, err := ioutil.TempFile("", "create-stack-tiny-acceptance")
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

	contents, err := ioutil.ReadFile(tempFilePath)
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
			contents, err := ioutil.ReadAll(tarReader)
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
