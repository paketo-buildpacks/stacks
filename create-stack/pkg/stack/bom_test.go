package stack_test

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cyclonedx "github.com/CycloneDX/cyclonedx-go"
	"github.com/anchore/syft/syft"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"github.com/sclevine/spec"

	gocontext "context"

	. "github.com/onsi/gomega"
)

func testBOM(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect        = NewWithT(t).Expect
		err           error
		bom           stack.BOM
		inputSyftFile string
		inputCdxFile  string
		dir           string
	)

	it.Before(func() {
		dir, err = os.MkdirTemp("", "")
		Expect(err).ToNot(HaveOccurred())

		inputSyftFile = filepath.Join(dir, "bom.syft.json")
		err = os.WriteFile(inputSyftFile, []byte("syft.json file contents"), 0600)
		Expect(err).ToNot(HaveOccurred())

		inputCdxFile = filepath.Join(dir, "bom.cdx.json")
		err = os.WriteFile(inputCdxFile, []byte("cdx.json file contents"), 0600)
		Expect(err).ToNot(HaveOccurred())
	})

	it.After(func() {
		os.RemoveAll(dir)
	})

	context("Generate", func() {
		var (
			dockerClient *client.Client
		)
		it.Before(func() {
			var err error
			dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			Expect(err).NotTo(HaveOccurred())

			ctx := gocontext.Background()

			logs, err := dockerClient.ImagePull(ctx, "alpine:latest", types.ImagePullOptions{})
			Expect(err).NotTo(HaveOccurred())
			defer logs.Close()

			// NOTE: required to force a wait for the image to pull
			_, err = io.Copy(io.Discard, logs)
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			_, err := dockerClient.ImageRemove(gocontext.Background(), "alpine:latest", types.ImageRemoveOptions{Force: true})
			Expect(err).NotTo(HaveOccurred())

			// Clean up untagged alpine image
			filters := filters.NewArgs()
			filters.Add("reference", "alpine")
			leftoverAlpineImages, err := dockerClient.ImageList(gocontext.Background(), types.ImageListOptions{Filters: filters})
			Expect(err).NotTo(HaveOccurred())

			for _, leftoverAlpineImage := range leftoverAlpineImages {
				for _, digest := range leftoverAlpineImage.RepoDigests {
					_, err = dockerClient.ImageRemove(gocontext.Background(), digest, types.ImageRemoveOptions{Force: true})
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		it("generates a syft and cyclonedx BOM", func() {
			outputPaths, err := bom.Generate("alpine:latest")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(outputPaths)).To(Equal(2))

			for _, outputPath := range outputPaths {
				bomContent, err := os.Open(outputPath)
				Expect(err).NotTo(HaveOccurred())

				if strings.Contains(outputPath, "syft") {
					sbomStruct, _, err := syft.Decode(bomContent)
					Expect(err).NotTo(HaveOccurred())
					Expect(sbomStruct.Artifacts.PackageCatalog.PackageCount()).Should(BeNumerically(">", 0))
					Expect(sbomStruct.Artifacts.Distro.Name()).To(Equal("alpine"))
				} else {
					var bom cyclonedx.BOM
					decoder := cyclonedx.NewBOMDecoder(bomContent, cyclonedx.BOMFileFormatJSON)
					Expect(decoder.Decode(&bom)).To(Succeed())
					Expect(len(*bom.Components)).Should(BeNumerically(">", 0))
					Expect(bom.Metadata.Component.Name).To(Equal("alpine:latest"))
				}
				Expect(os.Remove(outputPath)).To(Succeed())
			}
		})
		context("failure cases", func() {
			context("cannot source image to generate BOM for", func() {
				it("returns an error", func() {
					_, err = bom.Generate("nonexistent-image")
					Expect(err).To(MatchError(ContainSubstring("syft failed to source image:")))
				})
			})
		})
	})

	context("Attach", func() {
		var (
			dockerClient *client.Client
			// ogImgDigest  v1.Hash
		)
		it.Before(func() {
			var err error
			dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			Expect(err).NotTo(HaveOccurred())

			ctx := gocontext.Background()

			logs, err := dockerClient.ImagePull(ctx, "alpine:latest", types.ImagePullOptions{})
			Expect(err).NotTo(HaveOccurred())
			defer logs.Close()

			// NOTE: required to force a wait for the image to pull
			_, err = io.Copy(io.Discard, logs)
			Expect(err).NotTo(HaveOccurred())

			// If we end up following: https://github.com/buildpacks/rfcs/pull/186#discussion_r744368384
			// ogRef, err := name.ParseReference("alpine:latest")
			// Expect(err).NotTo(HaveOccurred())
			// ogImg, err := daemon.Image(ogRef)
			// Expect(err).NotTo(HaveOccurred())
			// ogImgDigest, err = ogImg.Digest()
			// Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			_, err := dockerClient.ImageRemove(gocontext.Background(), "alpine:latest", types.ImageRemoveOptions{Force: true})
			Expect(err).NotTo(HaveOccurred())

			// Clean up untagged alpine image
			filters := filters.NewArgs()
			filters.Add("reference", "alpine")
			leftoverAlpineImages, err := dockerClient.ImageList(gocontext.Background(), types.ImageListOptions{Filters: filters})
			Expect(err).NotTo(HaveOccurred())

			for _, leftoverAlpineImage := range leftoverAlpineImages {
				for _, digest := range leftoverAlpineImage.RepoDigests {
					_, err = dockerClient.ImageRemove(gocontext.Background(), digest, types.ImageRemoveOptions{Force: true})
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		it("attaches the BOM layer and corresponding label to the image", func() {
			err := bom.Attach("alpine:latest", []string{inputSyftFile, inputCdxFile})
			Expect(err).NotTo(HaveOccurred())

			ref, err := name.ParseReference("alpine:latest")
			Expect(err).NotTo(HaveOccurred())
			img, err := daemon.Image(ref)
			Expect(err).NotTo(HaveOccurred())

			// Check that the SBOM label exists
			sbomLayerSHA, labelExists, err := getLabel(img, "io.buildpacks.base.sbom")
			Expect(err).NotTo(HaveOccurred())
			Expect(labelExists).To(BeTrue())

			// Check SBOM layer from label contains SBOM files
			retrievedBOMs, err := retrieveBOMFromLayer(sbomLayerSHA, dir, img)
			Expect(err).NotTo(HaveOccurred())

			for _, bomPath := range retrievedBOMs {
				contents, err := os.ReadFile(bomPath)
				Expect(err).ToNot(HaveOccurred())
				fileBase := strings.TrimPrefix(filepath.Base(bomPath), "bom.")
				Expect(string(contents)).To(Equal(fmt.Sprintf("%s file contents", fileBase)))
			}
		})

		context("failure cases", func() {
			context("given an invalid image reference name", func() {
				it("returns an error", func() {
					err := bom.Attach("", []string{inputSyftFile})
					Expect(err).To(MatchError(ContainSubstring("bad imageRefString:")))
				})
			})

			context("given an invalid image to attach to", func() {
				it("returns an error", func() {
					err := bom.Attach("nonexistent-image", []string{inputSyftFile})
					Expect(err).To(MatchError(ContainSubstring("failed to retrieve image config:")))
				})
			})

			context("given no BOM files", func() {
				it("returns an error", func() {
					err := bom.Attach("alpine:latest", []string{})
					Expect(err).To(MatchError(ContainSubstring("failed to create BOM layer: no BOM files provided")))
				})
			})
		})
	})

	context("CreateLayer", func() {
		var (
			outputSyftFile string
			outputCdxFile  string
		)

		it.Before(func() {
			outputSyftFile = filepath.Join(dir, "new-syft.json")
			outputCdxFile = filepath.Join(dir, "new-cdx.json")
		})

		it("creates a layer that contains the given files when uncompressed", func() {
			mapping := []stack.InputOutputMapping{
				{
					FileName: inputSyftFile,
					DstPath:  outputSyftFile,
				},
				{
					FileName: inputCdxFile,
					DstPath:  outputCdxFile,
				},
			}

			layer, err := bom.CreateLayer(mapping)
			Expect(err).ToNot(HaveOccurred())

			// uncompress the layer and check out the files
			uncompressedLayer, err := layer.Uncompressed()
			Expect(err).ToNot(HaveOccurred())

			defer uncompressedLayer.Close()
			tr := tar.NewReader(uncompressedLayer)

			entryIndex := 0
			for {
				header, err := tr.Next()
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Fatal(err)
				}
				Expect(header.Name).To(Equal(mapping[entryIndex].DstPath))
				entryIndex++
			}
		})
		context("failure cases", func() {
			context("cannot open given input file", func() {
				it.Before(func() {
					Expect(os.Chmod(inputSyftFile, 0000)).To(Succeed())
				})
				it("returns an error", func() {
					_, err := bom.CreateLayer([]stack.InputOutputMapping{
						{
							FileName: inputSyftFile,
							DstPath:  outputSyftFile,
						},
					})
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}

func getLabel(img v1.Image, label string) (sha string, exists bool, err error) {
	cf, err := img.ConfigFile()
	if err != nil {
		return "", false, err
	}
	sha, ok := cf.Config.Labels[label]
	return sha, ok, nil
}

func retrieveBOMFromLayer(layerDiff, outputDir string, img v1.Image) (outputFiles []string, err error) {
	outputFiles = []string{}

	diffID, err := v1.NewHash(layerDiff)
	if err != nil {
		return outputFiles, err
	}

	bomLayer, err := img.LayerByDiffID(diffID)
	if err != nil {
		return outputFiles, err
	}

	layerReader, err := bomLayer.Uncompressed()
	if err != nil {
		return outputFiles, err
	}

	tr := tar.NewReader(layerReader)
	for {
		header, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return outputFiles, err
			}
			break
		}

		if header.Typeflag != tar.TypeReg ||
			!strings.HasPrefix(header.Name, "/cnb/sbom") {
			continue
		}
		s := strings.TrimPrefix(strings.TrimPrefix(header.Name, "/cnb/sbom"), "/")
		name := strings.Join(strings.Split(s, "/"), ".")
		size, err := writeFile(filepath.Join(outputDir, name), tr)
		if err != nil {
			return outputFiles, err
		}
		if size != header.Size {
			return outputFiles, errors.New("invalid tar: size mismatch")
		}
		outputFiles = append(outputFiles, filepath.Join(outputDir, name))
	}
	return outputFiles, nil
}

func writeFile(name string, r io.Reader) (n int64, err error) {
	f, err := os.Create(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}
