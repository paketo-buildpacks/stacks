package main_test

import (
	"encoding/json"
	. "github.com/paketo-buildpacks/stack-usns/actions/usn-monitor/entrypoint"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEntrypoint(t *testing.T) {
	spec.Run(t, "Entrypoint", testEntrypoint, spec.Report(report.Terminal{}))
}

func testEntrypoint(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath             string
		fullUSNListPath     string
		relevantUSNListPath string
		buildReceiptPath    string
		runReceiptPath      string
		tmpDirPath          string
		require             = require.New(t)
		assert              = assert.New(t)
	)

	it.Before(func() {
		var err error
		tmpDirPath, err = ioutil.TempDir("", "usn-update")

		fullUSNListPath = filepath.Join(tmpDirPath, "fullUSNList")

		relevantUSNListPath = filepath.Join(tmpDirPath, "relevantUSNList")

		buildReceiptPath = filepath.Join(tmpDirPath, "buildReceipt")

		runReceiptPath = filepath.Join(tmpDirPath, "runReceipt")

		cliPath = filepath.Join(tmpDirPath, "entrypoint")

		goBuild := exec.Command("go", "build", "-o", cliPath, ".")
		output, err := goBuild.CombinedOutput()
		require.NoError(err, "failed to build CLI: %s", string(output))
	})

	it.After(func() {
		_ = os.RemoveAll(tmpDirPath)
	})

	it("successfully finds relevant USNs", func() {
		allUSNs := []USN{
			{
				Title:            "USN-4504-1: OpenSSL vulnerabilities",
				AffectedPackages: []string{"openssl", "openssl1.0"},
			},
			{
				Title:            "USN-4498-1: Loofah vulnerability",
				AffectedPackages: []string{"ruby-loofah"},
			},
			{
				Title:            "USN-4499-1: MilkyTracker vulnerabilities",
				AffectedPackages: []string{"milkytracker"},
			},
		}

		jsonUSNArray, err := json.Marshal(allUSNs)
		require.NoError(err)
		err = ioutil.WriteFile(fullUSNListPath, jsonUSNArray, 0644)
		require.NoError(err)

		buildReceiptContent := `ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some-description`
		err = ioutil.WriteFile(buildReceiptPath, []byte(buildReceiptContent), 0644)
		require.NoError(err)

		runReceiptContents := `ii  milkytracker  1.6.10ubuntu0.1  amd64  some-description`
		err = ioutil.WriteFile(runReceiptPath, []byte(runReceiptContents), 0644)
		require.NoError(err)

		err = ioutil.WriteFile(relevantUSNListPath, []byte("[]"), 0644)
		require.NoError(err)

		cmd := exec.Command(cliPath,
			"--build-receipt", buildReceiptPath,
			"--run-receipt", runReceiptPath,
			"--full-usn-list", fullUSNListPath,
			"--relevant-usn-list", relevantUSNListPath)
		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		relevantBuildUSN := RecordedUSN{
			Title:   "USN-4498-1: Loofah vulnerability",
			Release: "unreleased",
		}
		relevantRunUSN := RecordedUSN{
			Title:   "USN-4499-1: MilkyTracker vulnerabilities",
			Release: "unreleased",
		}

		content, err := ioutil.ReadFile(relevantUSNListPath)
		require.NoError(err)

		var actualUSNArray []RecordedUSN
		err = json.Unmarshal(content, &actualUSNArray)
		assert.NoError(err)

		assert.Len(actualUSNArray, 2)
		assert.Contains(actualUSNArray, relevantBuildUSN)
		assert.Contains(actualUSNArray, relevantRunUSN)
	})
	it("omits duplicate relevant USNs", func() {
		allUSNs := []USN{
			{
				Title:            "USN-4499-1: MilkyTracker vulnerabilities",
				AffectedPackages: []string{"milkytracker"},
			},
		}

		jsonUSNArray, err := json.Marshal(allUSNs)
		require.NoError(err)
		err = ioutil.WriteFile(fullUSNListPath, jsonUSNArray, 0644)
		require.NoError(err)

		buildReceiptContent := `ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some-description`
		err = ioutil.WriteFile(buildReceiptPath, []byte(buildReceiptContent), 0644)
		require.NoError(err)

		runReceiptContents := `ii  milkytracker  1.6.10ubuntu0.1  amd64  some-description`
		err = ioutil.WriteFile(runReceiptPath, []byte(runReceiptContents), 0644)
		require.NoError(err)

		relevantRunUSNs := []RecordedUSN{{
			Title:   "USN-4499-1: MilkyTracker vulnerabilities",
			Release: "unreleased",
		}}
		relevantRunUSNsJson, err := json.Marshal(relevantRunUSNs)
		require.NoError(err)

		err = ioutil.WriteFile(relevantUSNListPath, relevantRunUSNsJson, 0644)
		require.NoError(err)

		cmd := exec.Command(cliPath,
			"--build-receipt", buildReceiptPath,
			"--run-receipt", runReceiptPath,
			"--full-usn-list", fullUSNListPath,
			"--relevant-usn-list", relevantUSNListPath)
		output, err := cmd.CombinedOutput()
		require.NoError(err, string(output))

		content, err := ioutil.ReadFile(relevantUSNListPath)
		require.NoError(err)

		var actualUSNArray []RecordedUSN
		err = json.Unmarshal(content, &actualUSNArray)
		assert.NoError(err)

		assert.Equal(actualUSNArray, relevantRunUSNs)
	})
}
