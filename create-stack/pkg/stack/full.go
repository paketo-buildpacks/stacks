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
}

func (bs FullStack) GetSources() string {
	return bs.sources
}

func (bs FullStack) GetBuildPackages() string {
	return bs.buildPackages
}

func (bs FullStack) GetRunPackages() string {
	return bs.runPackages
}

func (bs FullStack) GetName() string {
	return "full"
}

func (bs FullStack) GetBaseBuildDockerfilePath() string {
	return bs.baseBuildDockerfilePath
}

func (bs FullStack) GetBaseRunDockerfilePath() string {
	return bs.baseRunDockerfilePath
}

func (bs FullStack) GetCNBBuildDockerfilePath() string {
	return bs.cnbBuildDockerfilePath
}

func (bs FullStack) GetCNBRunDockerfilePath() string {
	return bs.cnbRunDockerfilePath
}

func (bs FullStack) GetBuildDescription() string {
	return "ubuntu:bionic + many common C libraries and utilities"
}

func (bs FullStack) GetRunDescription() string {
	return "ubuntu:bionic + many common C libraries and utilities"
}

func NewFullStack(stackDir string) (FullStack, error) {

	sources, err := ioutil.ReadFile(filepath.Join(stackDir, "arch", arch, "sources.list"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "build"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	runPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "run"))
	if err != nil {
		return FullStack{}, fmt.Errorf("failed to read run packages list file: %w", err)
	}

	baseBuildDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/build", stackDir)
	baseRunDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/run", stackDir)
	cnbBuildDockerfilePath := fmt.Sprintf("%s/bionic/cnb/build", stackDir)
	cnbRunDockerfilePath := fmt.Sprintf("%s/bionic/cnb/run", stackDir)

	return FullStack{
		sources:                 string(sources),
		buildPackages:           string(buildPackages),
		runPackages:             string(runPackages),
		baseBuildDockerfilePath: baseBuildDockerfilePath,
		baseRunDockerfilePath:   baseRunDockerfilePath,
		cnbBuildDockerfilePath:  cnbBuildDockerfilePath,
		cnbRunDockerfilePath:    cnbRunDockerfilePath,
	}, nil
}
