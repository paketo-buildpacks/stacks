package stack

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewTinyStack(buildTag, runTag, stackDir string, publish bool) (Definition, error) {
	sources, err := os.ReadFile(filepath.Join(stackDir, "arch", arch, "sources.list"))
	if err != nil {
		return Definition{}, fmt.Errorf("failed to read sources list file: %w", err)
	}

	buildPackages, err := os.ReadFile(filepath.Join(stackDir, "packages", "base", "build"))
	if err != nil {
		return Definition{}, fmt.Errorf("failed to read build packages list file: %w", err)
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
			Description: "ubuntu:bionic + openssl + CA certs + compilers + shell utilities",
			Args: []string{
				"stack_id=io.paketo.stacks.tiny",
			},
		},
		RunBase: Image{
			UseBuildKit: useBuildKit,
			Publish:     publish,
			Tag:         runTag,
			Dockerfile:  fmt.Sprintf("%s/tiny/dockerfile/run", stackDir),
		},
		RunCNB: Image{
			Publish:     publish,
			Tag:         fmt.Sprintf("%s-cnb", runTag),
			Dockerfile:  fmt.Sprintf("%s/tiny/cnb/run", stackDir),
			Description: "distroless-like bionic + glibc + openssl + CA certs",
		},
	}, nil
}
