package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type BaseStack struct {
	sources                 string
	buildPackages           string
	runPackages             string
	baseBuildDockerfilePath string
	baseRunDockerfilePath   string
	cnbBuildDockerfilePath  string
	cnbRunDockerfilePath    string
	architecture            string
}

func (bs BaseStack) WithBuildKit() bool {
	return false
}

func (bs BaseStack) GetSecretArgs() map[string]string {
	return nil
}

func (bs BaseStack) GetArchitecture() string {
	return bs.architecture
}

func (bs BaseStack) GetBaseBuildArgs() []string {
	return []string{
		fmt.Sprintf("ubuntu_image=%s", GetUbuntuImage(bs.GetArchitecture())),
		fmt.Sprintf("sources=%s", bs.sources),
		fmt.Sprintf("packages=%s", bs.buildPackages),
	}
}

func (bs BaseStack) GetBaseRunArgs() []string {
	return []string{
		fmt.Sprintf("ubuntu_image=%s", GetUbuntuImage(bs.GetArchitecture())),
		fmt.Sprintf("sources=%s", bs.sources),
		fmt.Sprintf("packages=%s", bs.runPackages),
	}
}

func (bs BaseStack) GetCNBBuildArgs() []string {
	return []string{}
}

func (bs BaseStack) GetCNBRunArgs() []string {
	return []string{}
}

func (bs BaseStack) GetName() string {
	return "base"
}

func (bs BaseStack) GetBaseBuildDockerfilePath() string {
	return bs.baseBuildDockerfilePath
}

func (bs BaseStack) GetBaseRunDockerfilePath() string {
	return bs.baseRunDockerfilePath
}

func (bs BaseStack) GetCNBBuildDockerfilePath() string {
	return bs.cnbBuildDockerfilePath
}

func (bs BaseStack) GetCNBRunDockerfilePath() string {
	return bs.cnbRunDockerfilePath
}

func (bs BaseStack) GetBuildDescription() string {
	return "ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
}

func (bs BaseStack) GetRunDescription() string {
	return "ubuntu:bionic + openssl + CA certs"
}

func NewBaseStack(stackDir string, architecture string) (BaseStack, error) {

	sources, err := ioutil.ReadFile(filepath.Join(stackDir, "arch", architecture, "sources.list"))
	if err != nil {
		return BaseStack{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "base", "build"))
	if err != nil {
		return BaseStack{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	runPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "base", "run"))
	if err != nil {
		return BaseStack{}, fmt.Errorf("failed to read run packages list file: %w", err)
	}

	baseBuildDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/build", stackDir)
	baseRunDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/run", stackDir)
	cnbBuildDockerfilePath := fmt.Sprintf("%s/bionic/cnb/build", stackDir)
	cnbRunDockerfilePath := fmt.Sprintf("%s/bionic/cnb/run", stackDir)

	return BaseStack{
		sources:                 string(sources),
		buildPackages:           string(buildPackages),
		runPackages:             string(runPackages),
		baseBuildDockerfilePath: baseBuildDockerfilePath,
		baseRunDockerfilePath:   baseRunDockerfilePath,
		cnbBuildDockerfilePath:  cnbBuildDockerfilePath,
		cnbRunDockerfilePath:    cnbRunDockerfilePath,
		architecture:            architecture,
	}, nil
}
