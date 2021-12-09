package acceptance_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testTiny(cliPath string) func(*testing.T, spec.G, spec.S) {
	return func(t *testing.T, when spec.G, it spec.S) {
		var (
			Expect = NewWithT(t).Expect

			containerID string
			settings    struct {
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

			settings.Build.Destination = "tiny-test/build"
			settings.Build.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Build.Destination, settings.Version, "tiny")
			settings.Build.CNBRef = settings.Build.BaseRef + "-cnb"

			settings.Run.Destination = "tiny-test/run"
			settings.Run.BaseRef = fmt.Sprintf("%s:%s-%s", settings.Run.Destination, settings.Version, "tiny")
			settings.Run.CNBRef = settings.Run.BaseRef + "-cnb"
		})

		it.After(func() {
			commands := [][]string{
				{"image", "rm", settings.Run.CNBRef, "--force"},
				{"image", "rm", settings.Run.BaseRef, "--force"},
				{"image", "rm", settings.Build.CNBRef, "--force"},
				{"image", "rm", settings.Build.BaseRef, "--force"},
			}

			if containerID != "" {
				commands = append([][]string{{"container", "rm", "--force", containerID}}, commands...)
			}

			for _, command := range commands {
				cmd := exec.Command("docker", command...)
				output, err := cmd.CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(output))
			}
		})

		it("builds tiny stack", func() {
			stacksDir, err := getStacksDirectory()
			Expect(err).NotTo(HaveOccurred())

			cmd := exec.Command(cliPath,
				"--build-destination", settings.Build.Destination,
				"--run-destination", settings.Run.Destination,
				"--version", settings.Version,
				"--stack", "tiny",
				"--stacks-dir", stacksDir,
			)
			cmd.Env = append(os.Environ(), "EXPERIMENTAL_ATTACH_RUN_IMAGE_SBOM=true")
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

			assertCommonLabels(t, TinyStackID, buildImageConfig)

			Expect(buildImageConfig.StackLabels.Description).To(Equal("ubuntu:bionic + openssl + CA certs + compilers + shell utilities"))
			Expect(buildImageConfig.StackLabels.Metadata).To(MatchJSON("{}"))
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

			assertCommonLabels(t, TinyStackID, runImageConfig)
			assertSBOMAttached(t, settings.Run.CNBRef, runImageConfig.StackLabels)

			Expect(runImageConfig.StackLabels.Description).To(Equal("distroless-like bionic + glibc + openssl + CA certs"))
			Expect(runImageConfig.StackLabels.Metadata).To(MatchJSON("{}"))
			Expect(runImageConfig.StackLabels.Mixins).To(ContainSubstring(`"ca-certificates"`))
			Expect(runImageConfig.StackLabels.Mixins).NotTo(ContainSubstring("build:"))
			Expect(runImageConfig.StackLabels.Packages).To(ContainSubstring(`"ca-certificates"`))
			// BOM label should contain the SHA of the last added layer on the image
			layers := runImageMetadata.RootFS.Layers
			Expect(runImageConfig.StackLabels.SBOM).To(Equal(layers[len(layers)-1]))

			runReleaseDate, err := time.Parse(time.RFC3339, runImageConfig.StackLabels.Released)
			Expect(err).NotTo(HaveOccurred(), string(output))
			Expect(runReleaseDate).To(BeTemporally("~", time.Now(), 10*time.Minute))

			Expect(runReleaseDate).To(Equal(buildReleaseDate))

			cmd = exec.Command(
				"docker", "create",
				settings.Run.BaseRef,
				"/tmp/app",
			)
			output, err = cmd.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to create container: %s", output))

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
			Expect(err).NotTo(HaveOccurred())
			Expect(passwdContent).To(Equal(`root:x:0:0:root:/root:/sbin/nologin
nobody:x:65534:65534:nobody:/nonexistent:/sbin/nologin
nonroot:x:65532:65532:nonroot:/home/nonroot:/sbin/nologin
`))

			groupContent, err := getContainerFile(containerID, "/etc/group")
			Expect(err).NotTo(HaveOccurred())
			Expect(groupContent).To(Equal(`root:x:0:
nobody:x:65534:
tty:x:5:
staff:x:50:
nonroot:x:65532:
`))

			osReleaseContent, err := getContainerFile(containerID, "/etc/os-release")
			Expect(err).NotTo(HaveOccurred())
			Expect(osReleaseContent).To(ContainSubstring(`PRETTY_NAME="Cloud Foundry Tiny"`))
			Expect(osReleaseContent).To(ContainSubstring(`HOME_URL="https://github.com/cloudfoundry/stacks"`))
			Expect(osReleaseContent).To(ContainSubstring(`SUPPORT_URL="https://github.com/cloudfoundry/stacks/blob/master/README.md"`))
			Expect(osReleaseContent).To(ContainSubstring(`BUG_REPORT_URL="https://github.com/cloudfoundry/stacks/issues/new"`))

			originalOSReleaseContent, err := getOriginalOSRelease(containerID)
			Expect(err).NotTo(HaveOccurred())

			assertOSReleaseEqual(t, originalOSReleaseContent, osReleaseContent)

			testAppFile, err := os.CreateTemp("", "create-stack-tiny-test-app")
			Expect(err).NotTo(HaveOccurred())

			testAppPath := testAppFile.Name()

			err = testAppFile.Close()
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(testAppPath)

			goBuild := exec.Command("go", "build", "-o", testAppPath, ".")
			goBuild.Env = append(os.Environ(), "GOOS=linux")
			goBuild.Dir = filepath.Join(stacksDir, "create-stack", "test", "acceptance", "fixtures", "tiny")
			output, err = goBuild.CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to build test app: %s", output))

			output, err = exec.Command("docker", "cp", testAppPath, containerID+":/tmp/app").CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to copy test app into container: %s", output))

			output, err = exec.Command("docker", "start", "-a", containerID).CombinedOutput()
			Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed to run test app in container: %s", output))
		})
	}
}
