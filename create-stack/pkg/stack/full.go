package stack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func NewFullStack(buildTag, runTag, stackDir string, publish bool) (Definition, error) {
	sources, err := ioutil.ReadFile(filepath.Join(stackDir, "arch", arch, "sources.list"))
	if err != nil {
		return Definition{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "build"))
	if err != nil {
		return Definition{}, fmt.Errorf("failed to read build packages list file: %w", err)
	}

	runPackages, err := ioutil.ReadFile(filepath.Join(stackDir, "packages", "full", "run"))
	if err != nil {
		return Definition{}, fmt.Errorf("failed to read run packages list file: %w", err)
	}

	useBuildKit := false

	return Definition{
		BuildBase: Image{
			UseBuildKit: useBuildKit,
			Publish:     publish,
			Tag:         buildTag,
			Dockerfile:  fmt.Sprintf("%s/bionic/dockerfile/build", stackDir),
			Args: []string{
				fmt.Sprintf("sources=%s", sources),
				fmt.Sprintf("packages=%s", buildPackages),
			},
		},
		BuildCNB: Image{
			Publish:     publish,
			Tag:         fmt.Sprintf("%s-cnb", buildTag),
			Dockerfile:  fmt.Sprintf("%s/bionic/cnb/build", stackDir),
			Description: "ubuntu:bionic + many common C libraries and utilities",
		},
		RunBase: Image{
			UseBuildKit: useBuildKit,
			Publish:     publish,
			Tag:         runTag,
			Dockerfile:  fmt.Sprintf("%s/bionic/dockerfile/run", stackDir),
			Args: []string{
				fmt.Sprintf("sources=%s", sources),
				fmt.Sprintf("packages=%s", runPackages),
			},
		},
		RunCNB: Image{
			Publish:     publish,
			Tag:         fmt.Sprintf("%s-cnb", runTag),
			Dockerfile:  fmt.Sprintf("%s/bionic/cnb/run", stackDir),
			Description: "ubuntu:bionic + many common C libraries and utilities",
		},
	}, nil
}
