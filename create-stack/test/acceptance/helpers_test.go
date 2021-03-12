package acceptance_test

import (
	"fmt"
	assertpkg "github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	DistroName    = "Ubuntu"
	DistroVersion = "18.04"
	Homepage      = "https://github.com/paketo-buildpacks/stacks"
	BionicStackID = "io.buildpacks.stacks.bionic"
	TinyStackID   = "io.paketo.stacks.tiny"
	Maintainer    = "Paketo Buildpacks"
)

func assertCommonLabels(t *testing.T, stackID string, imageConfig ImageConfig) {
	assertpkg.Equal(t, DistroName, imageConfig.StackLabels.DistroName)
	assertpkg.Equal(t, DistroVersion, imageConfig.StackLabels.DistroVersion)
	assertpkg.Equal(t, Homepage, imageConfig.StackLabels.Homepage)
	assertpkg.Equal(t, stackID, imageConfig.StackLabels.ID)
	assertpkg.Equal(t, Maintainer, imageConfig.StackLabels.Maintainer)
}

func getStacksDirectory() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to obtain directory")
	}
	filename = filepath.Dir(filename)
	return strings.TrimSuffix(filename, "create-stack/test/acceptance"), nil
}

type ImageConfig struct {
	User        string      `json:"User"`
	Env         []string    `json:"Env"`
	StackLabels StackLabels `json:"Labels"`
}

type StackLabels struct {
	Description   string `json:"io.buildpacks.stack.description"`
	DistroName    string `json:"io.buildpacks.stack.distro.name"`
	DistroVersion string `json:"io.buildpacks.stack.distro.version"`
	Homepage      string `json:"io.buildpacks.stack.homepage"`
	ID            string `json:"io.buildpacks.stack.id"`
	Maintainer    string `json:"io.buildpacks.stack.maintainer"`
	Metadata      string `json:"io.buildpacks.stack.metadata"`
	Mixins        string `json:"io.buildpacks.stack.mixins"`
	Released      string `json:"io.buildpacks.stack.released"`
	Packages      string `json:"io.paketo.stack.packages"`
}
