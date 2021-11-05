package acceptance_test

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testPublish(cliPath string) func(*testing.T, spec.G, spec.S) {
	return func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			settings struct {
				Version string
				Build   struct {
					Destination string
					BaseRef     string
					CNBRef      string
				}
				Run struct {
					Destination string
					BaseRef     string
					CNBRef      string
				}
			}
		)

		it.Before(func() {
			settings.Version = fmt.Sprintf("dev-%d", time.Now().UnixNano())

			settings.Build.Destination = "paketotesting/publish-test-build"
			settings.Build.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Build.Destination, settings.Version, "full")
			settings.Build.CNBRef = settings.Build.BaseRef + "-cnb"

			settings.Run.Destination = "paketotesting/publish-test-run"
			settings.Run.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Run.Destination, settings.Version, "full")
			settings.Run.CNBRef = settings.Run.BaseRef + "-cnb"
		})

		it.After(func() {
			for _, command := range [][]string{
				{"image", "rm", settings.Run.CNBRef, "--force"},
				{"image", "rm", settings.Run.BaseRef, "--force"},
				{"image", "rm", settings.Build.CNBRef, "--force"},
				{"image", "rm", settings.Build.BaseRef, "--force"},
			} {
				cmd := exec.Command("docker", command...)
				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(output))
			}
		})

		it("builds and publishes full bionic-stack", func() {
			stacksDir, err := getStacksDirectory()
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(cliPath,
				"--build-destination", settings.Build.Destination,
				"--run-destination", settings.Run.Destination,
				"--version", settings.Version,
				"--stack", "full",
				"--stacks-dir", stacksDir,
				"--publish",
			)
			output, err := cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), string(output))

			assertCorrectBaseImage(t, settings.Build.BaseRef)
			assertCorrectBaseImage(t, settings.Run.BaseRef)
		})
	}
}
