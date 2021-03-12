package packages

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PackageMetadata struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Arch    string        `json:"arch"`
	Source  SourcePackage `json:"source"`
	Summary string        `json:"summary"`
}

type SourcePackage struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	UpstreamVersion string `json:"upstreamVersion"`
}

type Bionic struct{}

func (t Bionic) GetBuildPackagesList(imageName string) ([]string, error) {
	return t.getPackagesList(imageName)
}

func (t Bionic) GetRunPackagesList(imageName string) ([]string, error) {
	return t.getPackagesList(imageName)
}

func (t Bionic) getPackagesList(imageName string) ([]string, error) {
	output, err := exec.Command("docker", "run", "--rm", imageName, "dpkg-query", "-f",
		"${Package}\\n", "-W",
	).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s packages list: %w\n%s", imageName, err, string(output))
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func (t Bionic) GetBuildPackageMetadata(imageName string) (string, error) {
	return t.getPackageMetadata(imageName)
}

func (t Bionic) GetRunPackageMetadata(imageName string) (string, error) {
	return t.getPackageMetadata(imageName)
}

func (t Bionic) getPackageMetadata(imageName string) (string, error) {
	output, err := exec.Command("docker", "run", "--rm", imageName, "dpkg-query", "-W", "-f",
		"${binary:Package};${Version};${Architecture};${binary:Summary};${source:Package};${source:Version};${source:Upstream-Version}\\n",
	).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get package metadata for %s: %w\n%s", imageName, err, string(output))
	}
	var packageMetadata []PackageMetadata

	for _, metadata := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		metadataFields := strings.Split(metadata, ";")
		if len(metadataFields) < 7 {
			return "", fmt.Errorf("not enough fields present in metadata %s", metadata)
		}
		packageMetadata = append(packageMetadata, PackageMetadata{
			Name:    metadataFields[0],
			Version: metadataFields[1],
			Arch:    metadataFields[2],
			Summary: strings.TrimSpace(metadataFields[3]),
			Source: SourcePackage{
				Name:            metadataFields[4],
				Version:         metadataFields[5],
				UpstreamVersion: metadataFields[6],
			},
		})
	}

	jsonMetadata, err := json.Marshal(packageMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal package metadata: %w", err)
	}

	return string(jsonMetadata), nil
}
