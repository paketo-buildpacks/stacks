package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type TinyStack struct {
	sources                 string
	buildPackages           string
	runPackages             string
	baseBuildDockerfilePath string
	baseRunDockerfilePath   string
	cnbBuildDockerfilePath  string
	cnbRunDockerfilePath    string
}

func (ts TinyStack) GetBaseBuildArgs() []string {
	return []string{
		fmt.Sprintf("sources=%s", ts.sources),
		fmt.Sprintf("packages=%s", ts.buildPackages),
	}
}

func (ts TinyStack) GetBaseRunArgs() []string {
	return []string{}
}

func (ts TinyStack) GetName() string {
	return "tiny"
}

func (ts TinyStack) GetBaseBuildDockerfilePath() string {
	return ts.baseBuildDockerfilePath
}

func (ts TinyStack) GetBaseRunDockerfilePath() string {
	return ts.baseRunDockerfilePath
}

func (ts TinyStack) GetCNBBuildDockerfilePath() string {
	return ts.cnbBuildDockerfilePath
}

func (ts TinyStack) GetCNBRunDockerfilePath() string {
	return ts.cnbRunDockerfilePath
}

func (ts TinyStack) GetBuildDescription() string {
	return "ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
}

func (ts TinyStack) GetRunDescription() string {
	return "distroless-like bionic + glibc + openssl + CA certs"
}

func NewTinyStack(stackDir string) (TinyStack, error) {

	sources, err := ioutil.ReadFile(filepath.Join(stackDir, "arch", arch, "sources.list"))
	if err != nil {
		return TinyStack{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "base", "build"))
	if err != nil {
		return TinyStack{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	baseBuildDockerfilePath := fmt.Sprintf("%s/bionic/dockerfile/build", stackDir)
	baseRunDockerfilePath := fmt.Sprintf("%s/tiny/dockerfile/run", stackDir)
	cnbBuildDockerfilePath := fmt.Sprintf("%s/bionic/cnb/build", stackDir)
	cnbRunDockerfilePath := fmt.Sprintf("%s/tiny/cnb/run", stackDir)

	return TinyStack{
		sources:                 string(sources),
		buildPackages:           string(buildPackages),
		baseBuildDockerfilePath: baseBuildDockerfilePath,
		baseRunDockerfilePath:   baseRunDockerfilePath,
		cnbBuildDockerfilePath:  cnbBuildDockerfilePath,
		cnbRunDockerfilePath:    cnbRunDockerfilePath,
	}, nil
}
