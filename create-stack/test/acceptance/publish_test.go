package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	assertpkg "github.com/stretchr/testify/assert"
	requirepkg "github.com/stretchr/testify/require"
)

func TestCreateStackBionicPublish(t *testing.T) {
	spec.Run(t, "CreateStackBionicPublish", testCreateStackBionicPublish, spec.Report(report.Terminal{}))
}

func testCreateStackBionicPublish(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath string
		require = requirepkg.New(t)
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

	it("builds and publishes full bionic-stack", func() {
		buildRepo := "paketotesting/test-create-stack-base-build"
		runRepo := "paketotesting/test-create-stack-base-run"
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
		assertCorrectBaseImage(t, buildBaseImageRef)

		runBaseImageRef := fmt.Sprintf("%s:%s-%s", runRepo, version, "full")
		assertCorrectBaseImage(t, runBaseImageRef)
	})
}

func assertCorrectBaseImage(t *testing.T, baseImageRef string) {
	cnbImageRef := baseImageRef + "-cnb"
	output, err := exec.Command("docker", "inspect", cnbImageRef, "--format", "{{json .Config}}").CombinedOutput()
	requirepkg.NoError(t, err, string(output))

	var runImageConfig ImageConfig
	err = json.Unmarshal(output, &runImageConfig)
	requirepkg.NoError(t, err)

	output, err = exec.Command("docker", "pull", baseImageRef).CombinedOutput()
	requirepkg.NoError(t, err, string(output))

	output, err = exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", baseImageRef).CombinedOutput()
	requirepkg.NoError(t, err, string(output))

	assertpkg.JSONEq(t, fmt.Sprintf(`{"base-image": "%s"}`, strings.TrimSpace(string(output))), runImageConfig.StackLabels.Metadata)
}
