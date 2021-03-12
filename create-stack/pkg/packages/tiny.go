package packages

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Tiny struct {
	BuildPkgs Bionic
}

func (t Tiny) GetBuildPackagesList(imageName string) ([]string, error) {
	return t.BuildPkgs.getPackagesList(imageName)
}

func (t Tiny) GetRunPackagesList(imageName string) ([]string, error) {
	tmpDir, err := t.getRunPackagesDir(imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get run packages for image %s: %w", imageName, err)
	}

	defer os.RemoveAll(tmpDir)

	var packages []string
	err = filepath.WalkDir(filepath.Join(tmpDir, "status.d"), func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		for _, line := range strings.Split(string(contents), "\n") {
			if strings.HasPrefix(line, "Package:") {
				packages = append(packages, strings.Split(line, " ")[1])
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read tiny packages: %w", err)
	}

	return packages, nil
}

func (t Tiny) GetBuildPackageMetadata(imageName string) (string, error) {
	return t.BuildPkgs.getPackageMetadata(imageName)
}

func (t Tiny) GetRunPackageMetadata(imageName string) (string, error) {
	tmpDir, err := t.getRunPackagesDir(imageName)
	if err != nil {
		return "", fmt.Errorf("failed to get run packages for image %s: %w", imageName, err)
	}

	defer os.RemoveAll(tmpDir)

	var packageMetadata []PackageMetadata

	err = filepath.WalkDir(filepath.Join(tmpDir, "status.d"), func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var name, version, arch, sourcePackage, sourceVersion, sourceUpstreamVersion, summary string
		for _, line := range strings.Split(string(contents), "\n") {
			switch {
			case strings.HasPrefix(line, "Package:"):
				name = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Version:"):
				version = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Architecture:"):
				arch = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Source-Package:"):
				sourcePackage = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Source-Version:"):
				sourceVersion = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Source-Upstream-Version:"):
				sourceUpstreamVersion = strings.Split(line, " ")[1]
			case strings.HasPrefix(line, "Description:"):
				summary = strings.TrimSpace(strings.TrimPrefix(line, "Description: "))
			}
		}

		packageMetadata = append(packageMetadata, PackageMetadata{
			Name:    name,
			Version: version,
			Arch:    arch,
			Source: SourcePackage{
				Name:            sourcePackage,
				Version:         sourceVersion,
				UpstreamVersion: sourceUpstreamVersion,
			},
			Summary: summary,
		})
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to read tiny packages: %w", err)
	}

	jsonMetadata, err := json.Marshal(packageMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal package metadata: %w", err)
	}

	return string(jsonMetadata), nil
}

func (t Tiny) getRunPackagesDir(imageName string) (string, error) {
	output, err := exec.Command("docker", "create", imageName, "sleep").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create run image %s: %w\n%s", imageName, err, string(output))
	}

	containerID := strings.TrimSpace(string(output))
	tmpDir, err := ioutil.TempDir("", "tiny")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp file: %w", err)
	}

	output, err = exec.Command("docker", "cp", fmt.Sprintf("%s:/var/lib/dpkg/status.d", containerID),
		tmpDir).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to copy package info from %s: %w\n%s", imageName, err, string(output))
	}

	output, err = exec.Command("docker", "rm", containerID).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to remove image %s: %w\n%s", containerID, err, string(output))
	}

	return tmpDir, nil
}
