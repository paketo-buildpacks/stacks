package stack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

//go:generate faux --interface PackageFinder --output fakes/package_finder.go
type PackageFinder interface {
	GetBuildPackagesList(image string) (list []string, err error)
	GetRunPackagesList(image string) (list []string, err error)
	GetBuildPackageMetadata(image string) (metadata string, err error)
	GetRunPackageMetadata(image string) (metadata string, err error)
}

//go:generate faux --interface MixinsGenerator --output fakes/mixins_generator.go
type MixinsGenerator interface {
	GetMixins(buildPackages, runPackages []string) (buildMixins []string, runMixins []string)
}

//go:generate faux --interface BOMGenerator --output fakes/bom_generator.go
type BOMGenerator interface {
	Generate(imageTag string) (outputPaths []string, err error)
	Attach(cnbImageTag string, files []string) (err error)
}

//go:generate faux --interface ImageClient --output fakes/image_client.go
type ImageClient interface {
	Build(tag, dockerfilePath string, withBuildKit bool, secrets map[string]string, buildArgs ...string) error
	Push(tag string) (string, error)
	Pull(tag string, keychain authn.Keychain) (v1.Image, error)
	SetLabel(tag, key, value string) error
}

type Stack interface {
	GetName() string
	WithBuildKit() bool
	GetSecretArgs() map[string]string
	GetBaseBuildArgs() []string
	GetBaseRunArgs() []string
	GetCNBBuildArgs() []string
	GetCNBRunArgs() []string
	GetBaseBuildDockerfilePath() string
	GetBaseRunDockerfilePath() string
	GetCNBBuildDockerfilePath() string
	GetCNBRunDockerfilePath() string
	GetBuildDescription() string
	GetRunDescription() string
}

type Creator struct {
	PackageFinder   PackageFinder
	MixinsGenerator MixinsGenerator
	ImageClient     ImageClient
	BOMGenerator    BOMGenerator
	AttachBOM       bool
}

type Definition struct {
	BuildBase Image
	BuildCNB  Image
	RunBase   Image
	RunCNB    Image
}

type Image struct {
	UseBuildKit bool
	Publish     bool
	Tag         string
	Dockerfile  string
	Description string
	Args        []string
	Secrets     map[string]string
}

// Deprecated: use Execute instead
func (c Creator) CreateStack(stackable Stack, buildBaseTag, runBaseTag string, publish bool) error {
	return c.Execute(Definition{
		BuildBase: Image{
			UseBuildKit: stackable.WithBuildKit(),
			Publish:     publish,
			Tag:         buildBaseTag,
			Dockerfile:  stackable.GetBaseBuildDockerfilePath(),
			Args:        stackable.GetBaseBuildArgs(),
			Secrets:     stackable.GetSecretArgs(),
		},
		BuildCNB: Image{
			Publish:     publish,
			Tag:         fmt.Sprintf("%s-cnb", buildBaseTag),
			Dockerfile:  stackable.GetCNBBuildDockerfilePath(),
			Description: stackable.GetBuildDescription(),
			Args:        stackable.GetCNBBuildArgs(),
		},
		RunBase: Image{
			UseBuildKit: stackable.WithBuildKit(),
			Publish:     publish,
			Tag:         runBaseTag,
			Dockerfile:  stackable.GetBaseRunDockerfilePath(),
			Args:        stackable.GetBaseRunArgs(),
			Secrets:     stackable.GetSecretArgs(),
		},
		RunCNB: Image{
			Publish:     publish,
			Tag:         fmt.Sprintf("%s-cnb", runBaseTag),
			Dockerfile:  stackable.GetCNBRunDockerfilePath(),
			Description: stackable.GetRunDescription(),
			Args:        stackable.GetCNBRunArgs(),
		},
	})
}

func (c Creator) Execute(def Definition) error {
	_, err := c.ImageClient.Pull("ubuntu:bionic", authn.DefaultKeychain)
	if err != nil {
		return fmt.Errorf("error pulling bionic image: %w", err)
	}

	buildBaseRef, err := c.buildBaseImage(def.BuildBase)
	if err != nil {
		return fmt.Errorf("error building base build image: %w", err)
	}

	buildBasePackageList, err := c.PackageFinder.GetBuildPackagesList(def.BuildBase.Tag)
	if err != nil {
		return fmt.Errorf("error getting packages from base build image: %w", err)
	}

	buildPackageMetadata, err := c.PackageFinder.GetBuildPackageMetadata(def.BuildBase.Tag)
	if err != nil {
		return fmt.Errorf("failed to get build package metadata: %w", err)
	}

	runBaseRef, err := c.buildBaseImage(def.RunBase)
	if err != nil {
		return fmt.Errorf("error building base run image: %w", err)
	}

	runBasePackageList, err := c.PackageFinder.GetRunPackagesList(def.RunBase.Tag)
	if err != nil {
		return fmt.Errorf("error getting packages from base run image: %w", err)
	}

	runPackageMetadata, err := c.PackageFinder.GetRunPackageMetadata(def.RunBase.Tag)
	if err != nil {
		return fmt.Errorf("failed to get run package metadata: %w", err)
	}

	buildMixins, runMixins := c.MixinsGenerator.GetMixins(buildBasePackageList, runBasePackageList)
	releaseDate := time.Now()

	err = c.buildCNBImage(def.BuildBase, def.BuildCNB, buildBaseRef, buildPackageMetadata, releaseDate, buildMixins, false)
	if err != nil {
		return fmt.Errorf("error building cnb build image: %w", err)
	}

	err = c.buildCNBImage(def.RunBase, def.RunCNB, runBaseRef, runPackageMetadata, releaseDate, runMixins, c.AttachBOM)
	if err != nil {
		return fmt.Errorf("error building cnb run image: %w", err)
	}

	return nil
}

func (c Creator) buildBaseImage(image Image) (string, error) {
	err := c.ImageClient.Build(image.Tag, image.Dockerfile, image.UseBuildKit, image.Secrets, image.Args...)
	if err != nil {
		return "", fmt.Errorf("failed to build base image: %w", err)
	}

	if image.Publish {
		ref, err := c.ImageClient.Push(image.Tag)
		if err != nil {
			return "", fmt.Errorf("failed to push tag %s: %w", image.Tag, err)
		}

		return ref, nil
	}

	return "", nil
}

func (c Creator) buildCNBImage(base, image Image, baseRef, packageMetadata string, releaseDate time.Time, mixinsList []string, bom bool) error {
	mixinsJSON, err := json.Marshal(mixinsList)
	if err != nil {
		return fmt.Errorf("failed to marshal mixin array: %w", err)
	}

	type Metadata struct {
		BaseImage string `json:"base-image,omitempty"`
	}

	var metadata Metadata
	if image.Publish {
		metadata = Metadata{BaseImage: baseRef}
	}

	metadataMap, err := json.Marshal(metadata)

	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	buildArgs := []string{
		fmt.Sprintf(`base_image=%s`, base.Tag),
		fmt.Sprintf("description=%s", image.Description),
		fmt.Sprintf("mixins=%s", string(mixinsJSON)),
		fmt.Sprintf("released=%s", releaseDate.Format(time.RFC3339)),
		fmt.Sprintf(`metadata=%s`, string(metadataMap)),
	}
	buildArgs = append(buildArgs, image.Args...)

	err = c.ImageClient.Build(image.Tag, image.Dockerfile, false, nil, buildArgs...)

	if err != nil {
		return fmt.Errorf("failed to build cnb image: %w", err)
	}

	if bom {
		// Generate 2 BOMs
		runBaseBOMPaths, err := c.BOMGenerator.Generate(base.Tag)
		if err != nil {
			return fmt.Errorf("error generating BOM: %w", err)
		}

		// Add BOMs to Layer
		err = c.BOMGenerator.Attach(image.Tag, runBaseBOMPaths)
		if err != nil {
			return fmt.Errorf("error attaching bom: %w", err)
		}
	}

	err = c.ImageClient.SetLabel(image.Tag, "io.paketo.stack.packages", packageMetadata)
	if err != nil {
		return fmt.Errorf("failed to set label: %w", err)
	}

	if image.Publish {
		_, err := c.ImageClient.Push(image.Tag)
		if err != nil {
			return fmt.Errorf("failed to push tag %s: %w", image.Tag, err)
		}
		return nil
	}

	return nil
}
