package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
)

const (
	DistroName    = "Ubuntu"
	DistroVersion = "18.04"
	Homepage      = "https://github.com/paketo-buildpacks/stacks"
	BionicStackID = "io.buildpacks.stacks.bionic"
	Maintainer    = "Paketo Buildpacks"
)

func TestEntrypoint(t *testing.T) {
	spec.Run(t, "Entrypoint", testEntrypoint, spec.Report(report.Terminal{}))
}

func testEntrypoint(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath string
		require = requirepkg.New(t)
		assert  = assertpkg.New(t)
	)

	it.Before(func() {
		tempFile, err := ioutil.TempFile("", "entrypoint")
		require.NoError(err)

		cliPath = tempFile.Name()
		require.NoError(tempFile.Close())

		goBuild := exec.Command("go", "build", "-o", cliPath, ".")
		output, err := goBuild.CombinedOutput()
		require.NoError(err, "failed to build CLI: %s", string(output))
	})

	it.After(func() {
		_ = os.Remove(cliPath)
	})

	it("builds base bionic-stack", func() {
		buildRepo := "paketotesting/test-create-stack-base-build"
		runRepo := "paketotesting/test-create-stack-base-run"
		version := fmt.Sprintf("dev-%d", time.Now().UnixNano())

		stacksDir, err := getStacksDirectory()
		require.NoError(err)

		cmd := exec.Command(cliPath,
			"--build-destination", buildRepo,
			"--run-destination", runRepo,
			"--version", version,
			"--stack", "base",
			"--stacks-dir", stacksDir,
			"--publish",
		)
		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		buildBaseImageRef := fmt.Sprintf("%s:%s-%s", buildRepo, version, "base")
		buildCNBImageRef := buildBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", buildCNBImageRef, "--format={{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var buildImageConfig ImageConfig
		err = json.Unmarshal(output, &buildImageConfig)
		require.NoError(err)

		assertCommonLabels(t, buildImageConfig)

		buildDescription := "ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
		assert.Equal(buildDescription, buildImageConfig.StackLabels.Description)

		output, err = exec.Command("docker", "pull", buildBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", buildBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		assert.JSONEq(fmt.Sprintf(`{"base-image": "%s"}`, strings.TrimSpace(string(output))), buildImageConfig.StackLabels.Metadata)

		assert.Contains(buildImageConfig.StackLabels.Mixins, `"build:make"`)
		assert.Contains(buildImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(buildImageConfig.StackLabels.Mixins, `"run:"`)

		buildReleaseDate, err := time.Parse(time.RFC3339, buildImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), buildReleaseDate, time.Minute*10)

		assert.Contains(buildImageConfig.StackLabels.Packages, `"ca-certificates"`)

		runBaseImageRef := fmt.Sprintf("%s:%s-%s", runRepo, version, "base")
		runCNBImageRef := runBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", runCNBImageRef, "--format", "{{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var runImageConfig ImageConfig
		err = json.Unmarshal(output, &runImageConfig)
		require.NoError(err)

		assertCommonLabels(t, runImageConfig)

		runDescription := "ubuntu:bionic + openssl + CA certs"
		assert.Equal(runDescription, runImageConfig.StackLabels.Description)

		output, err = exec.Command("docker", "pull", runBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", runBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		assert.JSONEq(fmt.Sprintf(`{"base-image": "%s"}`, strings.TrimSpace(string(output))), runImageConfig.StackLabels.Metadata)

		assert.Contains(runImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(runImageConfig.StackLabels.Mixins, "build:")

		runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), runReleaseDate, time.Minute*10)
		assert.Equal(buildReleaseDate, runReleaseDate)

		assert.Contains(runImageConfig.StackLabels.Packages, `"ca-certificates"`)
	})

	it("builds full bionic-stack", func() {
		buildRepo := "paketotesting/test-create-stack-full-build"
		runRepo := "paketotesting/test-create-stack-full-run"
		version := fmt.Sprintf("dev-%d", time.Now().UnixNano())

		stacksDir, err := getStacksDirectory()
		require.NoError(err)

		cmd := exec.Command(cliPath,
			"--build-destination", buildRepo,
			"--run-destination", runRepo,
			"--version", version,
			"--stack", "full",
			"--stacks-dir", stacksDir,
			"--publish",
		)
		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		buildBaseImageRef := fmt.Sprintf("%s:%s-%s", buildRepo, version, "full")
		buildCNBImageRef := buildBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", buildCNBImageRef, "--format={{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var buildImageConfig ImageConfig
		err = json.Unmarshal(output, &buildImageConfig)
		require.NoError(err)

		assertCommonLabels(t, buildImageConfig)

		buildDescription := "ubuntu:bionic + many common C libraries and utilities"
		assert.Equal(buildDescription, buildImageConfig.StackLabels.Description)

		output, err = exec.Command("docker", "pull", buildBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", buildBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		assert.JSONEq(fmt.Sprintf(`{"base-image": "%s"}`, strings.TrimSpace(string(output))), buildImageConfig.StackLabels.Metadata)

		assert.Contains(buildImageConfig.StackLabels.Mixins, `"build:cmake"`)
		assert.Contains(buildImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(buildImageConfig.StackLabels.Mixins, `"run:"`)

		buildReleaseDate, err := time.Parse(time.RFC3339, buildImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), buildReleaseDate, time.Minute*10)

		assert.Contains(buildImageConfig.StackLabels.Packages, `"ca-certificates"`)

		runBaseImageRef := fmt.Sprintf("%s:%s-%s", runRepo, version, "full")
		runCNBImageRef := runBaseImageRef + "-cnb"
		output, err = exec.Command("docker", "inspect", runCNBImageRef, "--format", "{{json .Config}}").CombinedOutput()
		require.NoError(err, string(output))

		var runImageConfig ImageConfig
		err = json.Unmarshal(output, &runImageConfig)
		require.NoError(err)

		assertCommonLabels(t, runImageConfig)

		runDescription := "ubuntu:bionic + many common C libraries and utilities"
		assert.Equal(runDescription, runImageConfig.StackLabels.Description)

		output, err = exec.Command("docker", "pull", runBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", runBaseImageRef).CombinedOutput()
		require.NoError(err, string(output))

		assert.JSONEq(fmt.Sprintf(`{"base-image": "%s"}`, strings.TrimSpace(string(output))), runImageConfig.StackLabels.Metadata)

		assert.Contains(runImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(runImageConfig.StackLabels.Mixins, "build:")

		runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), runReleaseDate, time.Minute*10)
		assert.Equal(buildReleaseDate, runReleaseDate)

		assert.Contains(runImageConfig.StackLabels.Packages, `"ca-certificates"`)
	})
}

func assertCommonLabels(t *testing.T, imageConfig ImageConfig) {
	assertpkg.Equal(t, DistroName, imageConfig.StackLabels.DistroName)
	assertpkg.Equal(t, DistroVersion, imageConfig.StackLabels.DistroVersion)
	assertpkg.Equal(t, Homepage, imageConfig.StackLabels.Homepage)
	assertpkg.Equal(t, BionicStackID, imageConfig.StackLabels.ID)
	assertpkg.Equal(t, Maintainer, imageConfig.StackLabels.Maintainer)
}

func getStacksDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to obtain directory")
	}
	return strings.TrimSuffix(filename, "create-stack/main_test.go"), nil
}

type ImageConfig struct {
	User        string      `json:"User"`
	Env         []string    `json:"Env"`
	StackLabels StackLabels `json:"Labels"`
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
	Released      string `json:"io.buildpacks.stack.released"`
	Packages      string `json:"io.paketo.stack.packages"`
}
