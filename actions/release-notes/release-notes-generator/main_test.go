package main_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	. "github.com/paketo-buildpacks/stacks/actions/release-notes/release-notes-generator"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReleaseNotesGenerator(t *testing.T) {
	spec.Run(t, "ReleaseNotesGenerator", testReleaseNotesGenerator, spec.Report(report.Terminal{}))
}

func testReleaseNotesGenerator(t *testing.T, when spec.G, it spec.S) {
	var (
		cliPath                       string
		relevantUSNs                  *os.File
		allUSNs                       *os.File
		fullReleaseNotes              string
		releaseNotesWithoutRunDiff    string
		releaseNotesWithoutUSNs       string
		releaseNotesWithoutBaseImages string
		relevantUsnArrayJson          []byte
		require                       = require.New(t)
		assert                        = assert.New(t)
	)

	const (
		fullReleaseNotesFilePath              = "testdata/full_release_notes.md"
		releaseNotesWithoutRunDiffFilePath    = "testdata/release_notes_without_run_diff.md"
		releaseNotesWithoutUSNsFilePath       = "testdata/release_notes_without_usns.md"
		releaseNotesWithoutBaseImagesFilePath = "testdata/release_notes_without_base_images.md"
		buildBaseImage                        = "some-registry/build@sha256:some-base-sha"
		buildCNBImage                         = "some-registry/build@sha256:some-cnb-sha"
		runBaseImage                          = "some-registry/run@sha256:some-base-sha"
		runCNBImage                           = "some-registry/run@sha256:some-cnb-sha"
		releaseVersion                        = "1.0"
		stack                                 = "tiny"
	)

	it.Before(func() {
		tempFile, err := os.CreateTemp("", "release-notes-generator")
		require.NoError(err)

		cliPath = tempFile.Name()
		require.NoError(tempFile.Close())

		relevantUSNs, err = os.CreateTemp("", "relevant-usns")
		require.NoError(err)

		relevantUSNArray := []RecordedUSN{
			{
				Title:   "USN-4498-1: Loofah vulnerability",
				Release: "1.0",
			},
			{
				Title:   "USN-4593-1: FreeType vulnerability",
				Release: "1.0",
			},
		}
		relevantUsnArrayJson, err = json.Marshal(relevantUSNArray)
		require.NoError(err)

		allUSNs, err = os.CreateTemp("", "all-usns")
		require.NoError(err)

		allUSNArray := []USN{
			{
				Title: "USN-4593-1: FreeType vulnerability",
				Link:  "https://ubuntu.com/security/notices/USN-4593-1",
				CveArray: []CVE{
					{
						Title:       "CVE-2020-15999",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2020-15999",
						Description: "A buffer overflow was discovered in Load_SBit_Png.",
					},
				},
				AffectedPackages: []string{"libfreetype6"},
			},
			{
				Title: "USN-4504-1: OpenSSL vulnerabilities",
				Link:  "https://ubuntu.com/security/notices/USN-4504-1",
				CveArray: []CVE{
					{
						Title:       "CVE-2019-1547",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-1547",
						Description: "Normally in OpenSSL EC groups always have a co-factor present and this is used in side channel resistant code paths. However, in some cases, it is possible to construct a group using explicit parameters (instead of using a named curve). In those cases it is possible that such a group does not have the cofactor present. This can occur even where all the parameters match a known named curve. If such a curve is used then OpenSSL falls back to non-side channel resistant code paths which may result in full key recovery during an ECDSA signature operation. In order to be vulnerable an attacker would have to have the ability to time the creation of a large number of signatures where explicit parameters with no co-factor present are in use by an application using libcrypto. For the avoidance of doubt libssl is not vulnerable because explicit parameters are never used. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
					},
					{
						Title:       "CVE-2019-1551",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-1551",
						Description: "There is an overflow bug in the x64_64 Montgomery squaring procedure used in exponentiation with 512-bit moduli. No EC algorithms are affected. Analysis suggests that attacks against 2-prime RSA1024, 3-prime RSA1536, and DSA1024 as a result of this defect would be very difficult to perform and are not believed likely. Attacks against DH512 are considered just feasible. However, for an attack the target would have to re-use the DH512 private key, which is not recommended anyway. Also applications directly using the low level API BN_mod_exp may be affected if they use BN_FLG_CONSTTIME. Fixed in OpenSSL 1.1.1e (Affected 1.1.1-1.1.1d). Fixed in OpenSSL 1.0.2u (Affected 1.0.2-1.0.2t).",
					},
					{
						Title:       "CVE-2019-1563",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-1563",
						Description: "In situations where an attacker receives automated notification of the success or failure of a decryption attempt an attacker, after sending a very large number of messages to be decrypted, can recover a CMS/PKCS7 transported encryption key or decrypt any RSA encrypted message that was encrypted with the public RSA key, using a Bleichenbacher padding oracle attack. Applications are not affected if they use a certificate together with the private RSA key to the CMS_decrypt or PKCS7_decrypt functions to select the correct recipient info to decrypt. Fixed in OpenSSL 1.1.1d (Affected 1.1.1-1.1.1c). Fixed in OpenSSL 1.1.0l (Affected 1.1.0-1.1.0k). Fixed in OpenSSL 1.0.2t (Affected 1.0.2-1.0.2s).",
					},
					{
						Title:       "CVE-2020-1968",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2020-1968",
						Description: "Raccoon Attack",
					},
				},
				AffectedPackages: []string{"openssl", "openssl1.0"},
			},
			{
				Title: "USN-4498-1: Loofah vulnerability",
				Link:  "https://ubuntu.com/security/notices/USN-4498-1",
				CveArray: []CVE{
					{
						Title:       "CVE-2019-15587",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-15587",
						Description: "In the Loofah gem for Ruby through v2.3.0 unsanitized JavaScript may occur in sanitized output when a crafted SVG element is republished.",
					},
				},
				AffectedPackages: []string{"ruby-loofah"},
			},
			{
				Title: "USN-4499-1: MilkyTracker vulnerabilities",
				Link:  "https://ubuntu.com/security/notices/USN-4499-1",
				CveArray: []CVE{
					{
						Title:       "CVE-2019-14464",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-14464",
						Description: "XMFile::read in XMFile.cpp in milkyplay in MilkyTracker 1.02.00 has a heap-based buffer overflow.",
					},
					{
						Title:       "CVE-2019-14496",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-14496",
						Description: "LoaderXM::load in LoaderXM.cpp in milkyplay in MilkyTracker 1.02.00 has a stack-based buffer overflow.",
					},
					{
						Title:       "CVE-2019-14497",
						Link:        "https://people.canonical.com/~ubuntu-security/cve/CVE-2019-14497",
						Description: "ModuleEditor::convertInstrument in tracker/ModuleEditor.cpp in MilkyTracker 1.02.00 has a heap-based buffer overflow.",
					},
				},
				AffectedPackages: []string{"milkytracker"},
			},
		}
		jsonUSNArray, err := json.Marshal(allUSNArray)
		_, err = allUSNs.Write(jsonUSNArray)
		require.NoError(err)

		fullReleaseNotesBytes, err := os.ReadFile(fullReleaseNotesFilePath)
		require.NoError(err)

		fullReleaseNotes = string(fullReleaseNotesBytes)

		releaseNotesWithoutRunDiffBytes, err := os.ReadFile(releaseNotesWithoutRunDiffFilePath)
		require.NoError(err)

		releaseNotesWithoutRunDiff = string(releaseNotesWithoutRunDiffBytes)

		releaseNotesWithoutUSNsBytes, err := os.ReadFile(releaseNotesWithoutUSNsFilePath)
		require.NoError(err)

		releaseNotesWithoutUSNs = string(releaseNotesWithoutUSNsBytes)

		releaseNotesWithoutBaseImagesBytes, err := os.ReadFile(releaseNotesWithoutBaseImagesFilePath)
		require.NoError(err)

		releaseNotesWithoutBaseImages = string(releaseNotesWithoutBaseImagesBytes)

		goBuild := exec.Command("go", "build", "-o", cliPath, ".")
		output, err := goBuild.CombinedOutput()
		require.NoError(err, "failed to build CLI: %s", string(output))
	})

	it.After(func() {
		_ = os.Remove(cliPath)
		_ = os.Remove(relevantUSNs.Name())
		_ = os.Remove(allUSNs.Name())
	})

	when("there is a run diff, build diff, and usns", func() {
		it("generates properly formatted release notes", func() {
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
				"--build-base-image="+buildBaseImage,
				"--build-cnb-image="+buildCNBImage,
				"--run-base-image="+runBaseImage,
				"--run-cnb-image="+runCNBImage,
				"--build-receipt-diff="+buildReceiptDiff,
				"--run-receipt-diff="+runReceiptDiff,
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
				"--stack="+stack,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			assert.Equal(fullReleaseNotes, string(output))
		})
	})

	when("there is no run receipt diff", func() {
		it("generates properly formatted release notes without a run receipt diff", func() {
			_, err := relevantUSNs.Write(relevantUsnArrayJson)
			require.NoError(err)

			buildReceiptDiff := `-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description
-ii  libfreetype6:amd64      2.8.1-2ubuntu2      amd64  some other description
+ii  libfreetype6:amd64      2.8.1-2ubuntu2.1    amd64  some other description`

			cmd := exec.Command(cliPath,
				"--build-base-image="+buildBaseImage,
				"--build-cnb-image="+buildCNBImage,
				"--run-base-image="+runBaseImage,
				"--run-cnb-image="+runCNBImage,
				"--build-receipt-diff="+buildReceiptDiff,
				"--run-receipt-diff=",
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
				"--stack="+stack,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			assert.Equal(releaseNotesWithoutRunDiff, string(output))
		})
	})

	when("there is no usn patch", func() {
		it("generates properly formatted release notes without a usn patch", func() {
			_, err := relevantUSNs.Write([]byte("[]"))
			require.NoError(err)

			buildReceiptDiff := `-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description`

			runReceiptDiff := `-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description`

			cmd := exec.Command(cliPath,
				"--build-base-image="+buildBaseImage,
				"--build-cnb-image="+buildCNBImage,
				"--run-base-image="+runBaseImage,
				"--run-cnb-image="+runCNBImage,
				"--build-receipt-diff="+buildReceiptDiff,
				"--run-receipt-diff="+runReceiptDiff,
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
				"--stack="+stack,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			assert.Equal(releaseNotesWithoutUSNs, string(output))
		})
	})

	when("there are no base images", func() {
		it("generates release notes without base image refs", func() {
			_, err := relevantUSNs.Write(relevantUsnArrayJson)
			require.NoError(err)

			buildReceiptDiff := `-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description
-ii  libfreetype6:amd64      2.8.1-2ubuntu2      amd64  some other description
+ii  libfreetype6:amd64      2.8.1-2ubuntu2.1    amd64  some other description`

			cmd := exec.Command(cliPath,
				"--build-base-image=",
				"--build-cnb-image="+buildCNBImage,
				"--run-base-image=",
				"--run-cnb-image="+runCNBImage,
				"--build-receipt-diff="+buildReceiptDiff,
				"--run-receipt-diff=",
				"--relevant-usns="+relevantUSNs.Name(),
				"--all-usns="+allUSNs.Name(),
				"--release-version="+releaseVersion,
				"--stack="+stack,
			)
			output, err := cmd.CombinedOutput()
			require.NoError(err, string(output))

			assert.Equal(releaseNotesWithoutBaseImages, string(output))
		})
	})
}
