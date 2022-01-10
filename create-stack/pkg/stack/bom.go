package stack

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/format"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

type BOM struct{}

func NewBOM() BOM {
	return BOM{}
}

type InputOutputMapping struct {
	FileName string
	File     io.ReadCloser
	DstPath  string
	Size     int64
}

// Generate takes in an existent image tag, and generates two Bill of Materials
// using the anchore/syft library, and then returns the filepaths to both created files.
// The function generates two files. One in the Syft JSON format,
// and the other in CycloneDX 1.3 JSON.
func (b BOM) Generate(imageTag string) ([]string, error) {
	src, cleanup, err := source.New(imageTag, &image.RegistryOptions{}, []string{})
	if err != nil {
		return []string{}, fmt.Errorf("syft failed to source image: %w", err)
	}
	defer cleanup()

	cfg := cataloger.Config{
		Search: cataloger.SearchConfig{
			Scope: source.SquashedScope,
		},
	}

	catalog, _, distro, err := syft.CatalogPackages(src, cfg)
	if err != nil {
		return []string{}, fmt.Errorf("syft failed to catalog packages from image: %w", err)
	}

	sbomResult := sbom.SBOM{
		Artifacts: sbom.Artifacts{
			PackageCatalog: catalog,
			Distro:         distro,
		},
		Source: src.Metadata,
	}

	syftSBOM, err := encodeBOM(sbomResult, format.JSONOption, "syft.json")
	if err != nil {
		return []string{}, fmt.Errorf("failed to encode Syft BOM: %w", err)
	}
	if err := syftSBOM.Close(); err != nil {
		return []string{}, fmt.Errorf("failed to close Syft BOM file: %w", err)
	}

	cyclonedxSBOM, err := encodeBOM(sbomResult, format.CycloneDxJSONOption, "cdx.json")
	if err != nil {
		return []string{}, fmt.Errorf("failed to encode CycloneDX BOM: %w", err)
	}

	return []string{syftSBOM.Name(), cyclonedxSBOM.Name()}, nil
}

// Attach takes in the image tag of an image on the local Docker daemon,
// and files to add to that image. It expects a files to have `syft.json` in the name to denote a
// Syft-type BOM file, and it denotes all other given files as `bom.cdx.json` to
// denote CycloneDX.
// Files are added to `/cnb/sbom/`, and are turned into a layer and added to the image.
// This function also adds a layer to the image under
// `io.buildpacks.base.sbom`, with the value set to the diffID of the newly
// added layer.
func (b BOM) Attach(imageTag string, files []string) error {
	ref, err := name.ParseReference(imageTag)
	if err != nil {
		return fmt.Errorf("bad imageRefString: %w", err)
	}

	img, err := daemon.Image(ref)
	if err != nil {
		return fmt.Errorf("failed to retrieve image: %w", err)
	}

	// From open discussion on the upstream RFC:
	// https://github.com/buildpacks/rfcs/pull/186#discussion_r744368384
	digest, err := img.Digest()
	if err != nil {
		return fmt.Errorf("failed to retrieve image digest: %w", err)
	}
	// For use in:
	//DstPath:  path.Join("/cnb/sbom", digest.Hex[:8]+"."+strings.TrimPrefix(ext, ".")),

	dstPaths := []InputOutputMapping{}
	for _, file := range files {
		ext := "syft.json"
		if !strings.Contains(file, "syft") {
			ext = "cdx.json"
		}

		dstPaths = append(dstPaths, InputOutputMapping{
			FileName: file,
			// DstPath:  path.Join("/cnb/sbom", "sbom."+ext),
			DstPath: path.Join("/cnb/sbom", digest.Hex[:8]+"."+strings.TrimPrefix(ext, ".")),
		})
	}

	layer, err := b.CreateLayer(dstPaths)
	if err != nil {
		return fmt.Errorf("failed to create BOM layer: %w", err)
	}

	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return fmt.Errorf("failed to append BOM layer to image: %w", err)
	}
	diffID, err := layer.DiffID()
	if err != nil {
		return fmt.Errorf("failed to retrieve layer diff: %w", err)
	}
	cf, err := img.ConfigFile()
	if err != nil {
		return fmt.Errorf("failed to retrieve image config: %w", err)
	}
	if cf.Config.Labels == nil {
		cf.Config.Labels = make(map[string]string, 1)
	}
	cf.Config.Labels["io.buildpacks.base.sbom"] = diffID.String()
	img, err = mutate.ConfigFile(img, cf)
	if err != nil {
		return fmt.Errorf("failed to modify image config: %w", err)
	}

	tag, err := name.NewTag(imageTag)
	if err != nil {
		return fmt.Errorf("failed to create tag with %s: %w", imageTag, err)
	}

	_, err = daemon.Write(tag, img)
	if err != nil {
		return fmt.Errorf("failed to write %s to daemon: %w", tag, err)
	}

	// Clean up the create BOM files
	for _, path := range files {
		defer os.Remove(path)
	}

	return err
}

func (b BOM) CreateLayer(mappings []InputOutputMapping) (v1.Layer, error) {
	if len(mappings) == 0 {
		return nil, errors.New("no BOM files provided")
	}
	return tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		for i, entry := range mappings {
			f, err := os.Open(entry.FileName)
			if err != nil {
				return nil, err
			}
			fd, err := f.Stat()
			if err != nil {
				return nil, err
			}
			entry.Size = fd.Size()
			entry.File = f
			mappings[i] = entry
		}
		return tarFile(mappings), nil
	})
}

func tarFile(mappings []InputOutputMapping) io.ReadCloser {
	out, w := io.Pipe()
	go func() {
		tw := tar.NewWriter(w)

		for _, entry := range mappings {
			header := &tar.Header{
				Name: entry.DstPath,
				Size: entry.Size,
				Mode: 0600,
			}

			if err := tw.WriteHeader(header); err != nil {
				w.CloseWithError(err)
				return
			}

			if _, err := io.Copy(tw, entry.File); err != nil {
				w.CloseWithError(err)
				return
			}

			defer entry.File.Close()
			defer w.Close() // always nil + never overrides
		}
		if err := tw.Close(); err != nil {
			w.CloseWithError(err)
			return
		}
	}()
	return out
}

func encodeBOM(sbomStruct sbom.SBOM, format format.Option, filename string) (*os.File, error) {
	bomBytes, err := syft.Encode(sbomStruct, format)
	if err != nil {
		return nil, fmt.Errorf("could not encode data into %s format: %w", format, err)
	}

	bomFile, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file %s: %w", filename, err)
	}

	if _, err = bomFile.Write(bomBytes); err != nil {
		return nil, fmt.Errorf("could not write BOM to %s: %w", bomFile.Name(), err)
	}

	return bomFile, nil
}
