package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type FullStack struct {
	sources                 string
	buildPackages           string
	runPackages             string
	baseBuildDockerfilePath string
	baseRunDockerfilePath   string
	cnbBuildDockerfilePath  string
	cnbRunDockerfilePath    string
	architecture            string
}

func (fs FullStack) GetArchitecture() string {
	return fs.architecture
}

func (fs FullStack) WithBuildKit() bool {
	return false
}

func (fs FullStack) GetSecretArgs() map[string]string {
	return nil
}

func (fs FullStack) GetBaseBuildArgs() []string {
	return []string{
		fmt.Sprintf("ubuntu_image=%s", GetUbuntuImage(fs.GetArchitecture())),
		fmt.Sprintf("sources=%s", fs.sources),
		fmt.Sprintf("packages=%s", fs.buildPackages),
	}
}

func (fs FullStack) GetBaseRunArgs() []string {
	return []string{
		fmt.Sprintf("ubuntu_image=%s", GetUbuntuImage(fs.GetArchitecture())),
		fmt.Sprintf("sources=%s", fs.sources),
		fmt.Sprintf("packages=%s", fs.runPackages),
	}
}

func (fs FullStack) GetCNBBuildArgs() []string {
	return []string{}
}

func (fs FullStack) GetCNBRunArgs() []string {
	return []string{}
}

func (fs FullStack) GetName() string {
	return "full"
}

func (fs FullStack) GetBaseBuildDockerfilePath() string {
	return fs.baseBuildDockerfilePath
}

func (fs FullStack) GetBaseRunDockerfilePath() string {
	return fs.baseRunDockerfilePath
}

func (fs FullStack) GetCNBBuildDockerfilePath() string {
	return fs.cnbBuildDockerfilePath
}

func (fs FullStack) GetCNBRunDockerfilePath() string {
	return fs.cnbRunDockerfilePath
}

func (fs FullStack) GetBuildDescription() string {
	return "ubuntu:bionic + many common C libraries and utilities"
}

func (fs FullStack) GetRunDescription() string {
	return "ubuntu:bionic + many common C libraries and utilities"
}

func NewFullStack(stackDir string, architecture string) (FullStack, error) {

	sources, err := ioutil.ReadFile(filepath.Join(stackDir, "arch", architecture, "sources.list"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "build", "_common"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	buildPackagesArch, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "build", architecture))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	runPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "run", "_common"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read run packages list file: %w", err)
	}

	runPackagesArch, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "run", architecture))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read run packages list file: %w", err)
	}

	baseBuildDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/build", stackDir)
	baseRunDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/run", stackDir)
	cnbBuildDockerfilePath := fmt.Sprintf("%s/bionic/cnb/build", stackDir)
	cnbRunDockerfilePath := fmt.Sprintf("%s/bionic/cnb/run", stackDir)

	return FullStack{
		sources:                 string(sources),
		buildPackages:           string(buildPackages) + string(buildPackagesArch),
		runPackages:             string(runPackages) + string(runPackagesArch),
		baseBuildDockerfilePath: baseBuildDockerfilePath,
		baseRunDockerfilePath:   baseRunDockerfilePath,
		cnbBuildDockerfilePath:  cnbBuildDockerfilePath,
		cnbRunDockerfilePath:    cnbRunDockerfilePath,
		architecture:            architecture,
	}, nil
}
