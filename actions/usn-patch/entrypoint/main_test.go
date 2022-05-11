package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/paketo-buildpacks/stacks/actions/usn-patch/entrypoint"
)

func TestEntrypoint(t *testing.T) {
	spec.Run(t, "Entrypoint", testEntrypoint, spec.Report(report.Terminal{}))
}

func testEntrypoint(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath              string
		relevantUSNs         *os.File
		allUSNs              *os.File
		relevantUsnArrayJson []byte
		require              = require.New(t)
		assert               = assert.New(t)
	)

	const (
		releaseVersion = "1.0"
	)

	it.Before(func() {
		tempFile, err := os.CreateTemp("", "entrypoint")
		require.NoError(err)

		cliPath = tempFile.Name()
		require.NoError(tempFile.Close())

		relevantUSNs, err = os.CreateTemp("", "relevant-usns")
		require.NoError(err)

		relevantUSNArray := []RecordedUSN{
			{
				Title:   "USN-4498-1: Loofah vulnerability",
				Release: "unreleased",
			},
			{
				Title:   "USN-4593-1: FreeType vulnerability",
				Release: "unreleased",
			},
		}
		relevantUsnArrayJson, err = json.MarshalIndent(relevantUSNArray, "", "    ")
		require.NoError(err)

		allUSNs, err = os.CreateTemp("", "all-usns")
		require.NoError(err)

		allUSNArray := []USN{
			{
				Title:            "USN-4593-1: FreeType vulnerability",
				AffectedPackages: []string{"libfreetype6"},
			},
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
		jsonUSNArray, err := json.Marshal(allUSNArray)
		_, err = allUSNs.Write(jsonUSNArray)
		require.NoError(err)

		goBuild := exec.Command("go", "build", "-o", cliPath, ".")
		output, err := goBuild.CombinedOutput()
		require.NoError(err, "failed to build CLI: %s", string(output))
	})

	it.After(func() {
		_ = os.Remove(cliPath)
		_ = os.Remove(relevantUSNs.Name())
		_ = os.Remove(allUSNs.Name())
	})

	when("there is a patched usn", func() {
		it("updates the usn release", func() {
			_, err := relevantUSNs.Write(relevantUsnArrayJson)
			require.NoError(err)

			buildReceiptDiff := `-ii  ruby-loofah          1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah       1.6.12ubuntu0.1      all    some longer description
-ii  ruby-loofaher     1.6.0   amd64    some description
+ii  ruby-loofaher     1.6.12Trusty0.1.23   amd64   some description
-ii  libfreetype6:amd64      2.8.1-2ubuntu2      amd64  some other description
+ii  libfreetype6:amd64      2.8.1-2ubuntu2.1    amd64  some other description`

			runReceiptDiff := `-ii  ruby-loofah      1.6.10ubuntu0.1  amd64
+ii  ruby-loofah      1.6.12ubuntu0.1       amd64
-ii  ruby-boofah      1.6.10ubuntu0.1         amd64
+ii  ruby-boofah      1.6.12ubuntu0.1      amd64`

			cmd := exec.Command(cliPath,
				"--build-receipt-diff="+buildReceiptDiff,
				"--run-receipt-diff="+runReceiptDiff,
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			relevantUSNsContent, err := os.ReadFile(relevantUSNs.Name())
			require.NoError(err)

			var updatedUSNs []RecordedUSN
			err = json.Unmarshal(relevantUSNsContent, &updatedUSNs)
			require.NoError(err)

			assert.Len(updatedUSNs, 2)
			assert.Equal(releaseVersion, updatedUSNs[0].Release)
			assert.Equal(releaseVersion, updatedUSNs[1].Release)
		})
	})

	when("there is no run patched usn", func() {
		it("does nothing", func() {
			_, err := relevantUSNs.Write(relevantUsnArrayJson)
			require.NoError(err)

			cmd := exec.Command(cliPath,
				"--build-receipt-diff=",
				"--run-receipt-diff=",
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			relevantUSNsContent, err := os.ReadFile(relevantUSNs.Name())
			require.NoError(err)

			fmt.Println(string(relevantUSNsContent))
			fmt.Println("*********************")
			fmt.Println(string(relevantUsnArrayJson))

			assert.Equal(relevantUsnArrayJson, relevantUSNsContent)
		})
	})
}
