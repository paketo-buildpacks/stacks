package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
)

func TestCreateStackBionic(t *testing.T) {
	spec.Run(t, "CreateStackBionic", testCreateStackBionic, spec.Report(report.Terminal{}))
}

func testCreateStackBionic(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath string
		require = requirepkg.New(t)
		assert  = assertpkg.New(t)
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

		assertCommonLabels(t, BionicStackID, buildImageConfig)

		buildDescription := "ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
		assert.Equal(buildDescription, buildImageConfig.StackLabels.Description)

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

		assertCommonLabels(t, BionicStackID, runImageConfig)

		runDescription := "ubuntu:bionic + openssl + CA certs"
		assert.Equal(runDescription, runImageConfig.StackLabels.Description)

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

		assertCommonLabels(t, BionicStackID, buildImageConfig)

		buildDescription := "ubuntu:bionic + many common C libraries and utilities"
		assert.Equal(buildDescription, buildImageConfig.StackLabels.Description)

		assert.JSONEq(`{}`, buildImageConfig.StackLabels.Metadata)

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

		assertCommonLabels(t, BionicStackID, runImageConfig)

		runDescription := "ubuntu:bionic + many common C libraries and utilities"
		assert.Equal(runDescription, runImageConfig.StackLabels.Description)

		assert.JSONEq(`{}`, runImageConfig.StackLabels.Metadata)

		assert.Contains(runImageConfig.StackLabels.Mixins, `"ca-certificates"`)
		assert.NotContains(runImageConfig.StackLabels.Mixins, "build:")

		runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
		require.NoError(err)
		assert.WithinDuration(time.Now(), runReleaseDate, time.Minute*10)
		assert.Equal(buildReleaseDate, runReleaseDate)

		assert.Contains(runImageConfig.StackLabels.Packages, `"ca-certificates"`)
	})
}
