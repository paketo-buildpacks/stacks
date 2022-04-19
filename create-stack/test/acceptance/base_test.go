package acceptance_test

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBase(cliPath string) func(*testing.T, spec.G, spec.S) {
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

			settings.Build.Destination = "base-test/build"
			settings.Build.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Build.Destination, settings.Version, "base")
			settings.Build.CNBRef = settings.Build.BaseRef + "-cnb"

			settings.Run.Destination = "base-test/run"
			settings.Run.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Run.Destination, settings.Version, "base")
			settings.Run.CNBRef = settings.Run.BaseRef + "-cnb"
		})

		it.After(func() {
			commands := [][]string{
				{"image", "rm", settings.Run.CNBRef, "--force"},
				{"image", "rm", settings.Run.BaseRef, "--force"},
				{"image", "rm", settings.Build.CNBRef, "--force"},
				{"image", "rm", settings.Build.BaseRef, "--force"},
			}

			for _, command := range commands {
				cmd := exec.Command("docker", command...)
				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(output))
			}
		})

		it("builds base bionic-stack", func() {
			stacksDir, err := getStacksDirectory()
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(cliPath,
				"--build-destination", settings.Build.Destination,
				"--run-destination", settings.Run.Destination,
				"--version", settings.Version,
				"--stack", "base",
				"--stacks-dir", stacksDir,
			)
			output, err := cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), string(output))

			cmd = exec.Command(
				"docker", "inspect",
				settings.Build.CNBRef,
				"--format={{json .Config}}",
			)
			output, err = cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), string(output))

			var buildImageConfig ImageConfig
			err = json.Unmarshal(output, &buildImageConfig)
			Expect(err).NotTo(HaveOccurred())

			assertCommonLabels(t, BionicStackID, buildImageConfig)

			Expect(buildImageConfig.StackLabels.Description).To(Equal("ubuntu:bionic + openssl + CA certs + compilers + shell utilities"))
			Expect(buildImageConfig.StackLabels.Mixins).To(ContainSubstring(`"build:make"`))
			Expect(buildImageConfig.StackLabels.Mixins).To(ContainSubstring(`"ca-certificates"`))
			Expect(buildImageConfig.StackLabels.Mixins).NotTo(ContainSubstring("run:"))
			Expect(buildImageConfig.StackLabels.Packages).To(ContainSubstring(`"ca-certificates"`))

			buildReleaseDate, err := time.Parse(time.RFC3339, buildImageConfig.StackLabels.Released)
			Expect(err).NotTo(HaveOccurred())
			Expect(buildReleaseDate).To(BeTemporally("~", time.Now(), 10*time.Minute))

			cmd = exec.Command(
				"docker", "inspect",
				settings.Run.CNBRef,
				"--format", "{{json .}}",
			)
			output, err = cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), string(output))

			var runImageMetadata ImageMetadata
			var runImageConfig ImageConfig
			err = json.Unmarshal(output, &runImageMetadata)
			Expect(err).NotTo(HaveOccurred(), string(output))

			runImageConfig = runImageMetadata.ImageConfig

			assertCommonLabels(t, BionicStackID, runImageConfig)

			Expect(runImageConfig.StackLabels.Description).To(Equal("ubuntu:bionic + openssl + CA certs"))
			Expect(runImageConfig.StackLabels.Mixins).To(ContainSubstring(`"ca-certificates"`))
			Expect(runImageConfig.StackLabels.Mixins).NotTo(ContainSubstring("build:"))
			Expect(runImageConfig.StackLabels.Packages).To(ContainSubstring(`"ca-certificates"`))
			Expect(runImageConfig.StackLabels.SBOM).To(BeEmpty())

			runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
			Expect(err).NotTo(HaveOccurred())
			Expect(runReleaseDate).To(BeTemporally("~", time.Now(), 10*time.Minute))
			Expect(runReleaseDate).To(Equal(buildReleaseDate))

			cmd = exec.Command(
				"docker", "run", "--rm",
				settings.Build.BaseRef,
				"git", "config", "--system", "--list",
			)
			output, err = cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to run container: %s", output))

			Expect(string(output)).To(ContainSubstring("safe.directory=/workspace"))
			Expect(string(output)).To(ContainSubstring("safe.directory=/workspace/source-ws"))
			Expect(string(output)).To(ContainSubstring("safe.directory=/workspace/source"))
		})
	}
}
