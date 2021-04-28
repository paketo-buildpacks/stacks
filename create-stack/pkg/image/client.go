package image

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Client struct{}

func (c Client) Build(tag, dockerfilePath string, withBuildKit bool, secrets map[string]string, buildArgs ...string) error {
	if withBuildKit {
		err := os.Setenv("DOCKER_BUILDKIT", "1")
		if err != nil {
			return fmt.Errorf("failed to set DOCKER_BUILDKIT environment variable: %w", err)
		}
	}
	finalArgs := []string{"build", "-t", tag, "--no-cache"}

	if len(secrets) != 0 {
		dir, err := ioutil.TempDir("", "docker-secrets")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}

		for id, secret := range secrets {
			file, err := os.Create(fmt.Sprintf("%s/secret-%s", dir, id))
			if err != nil {
				return fmt.Errorf("failed to create secret file: %w", err)
			}

			_, err = file.WriteString(secret)
			if err != nil {
				return fmt.Errorf("failed to write secret to file: %w", err)
			}

			finalArgs = append(finalArgs, "--secret", fmt.Sprintf(`id=%s,src=%s`, id, file.Name()))
		}
	}

	for _, elem := range buildArgs {
		finalArgs = append(finalArgs, "--build-arg", elem)
	}

	finalArgs = append(finalArgs, dockerfilePath)

	output, err := exec.Command("docker", finalArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build cnb image: %w\n%s", err, string(output))
	}

	return nil
}

func (c Client) Push(tag string) (string, error) {
	ref, err := name.ParseReference(tag)
	if err != nil {
		return "", fmt.Errorf("could not parse reference: %w", err)
	}

	image, err := daemon.Image(ref)
	if err != nil {
		return "", fmt.Errorf("could not get image from daemon: %w", err)
	}

	err = remote.Write(ref, image, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", fmt.Errorf("could not write image manifest: %w", err)
	}

	digest, err := image.Digest()
	if err != nil {
		return "", fmt.Errorf("could not get image digest: %w", err)
	}

	repo := strings.Split(tag, ":")[0]

	return fmt.Sprintf("%s@%s", repo, digest.String()), nil
}

func (c Client) SetLabel(tag, key, value string) error {
	ref, err := name.ParseReference(tag)
	if err != nil {
		return fmt.Errorf("could not parse reference: %w", err)
	}

	image, err := daemon.Image(ref)
	if err != nil {
		return fmt.Errorf("could not get image from daemon: %w", err)
	}

	updatedImage, err := setLabelOnImage(image, key, value)
	if err != nil {
		return fmt.Errorf("could not set label: %w", err)
	}

	localTag, err := name.NewTag(tag)
	if err != nil {
		return fmt.Errorf("could not create local tag: %w", err)
	}

	_, err = daemon.Write(localTag, updatedImage)
	if err != nil {
		return fmt.Errorf("could not write image to docker daemon: %w", err)
	}

	return nil
}

func (c Client) Pull(tag string, keychain authn.Keychain) (v1.Image, error) {

	ref, err := name.ParseReference(tag)
	if err != nil {
		return nil, fmt.Errorf("could not pull image manifest: %w", err)
	}

	image, err := remote.Image(ref, remote.WithAuthFromKeychain(keychain))
	if err != nil {
		return nil, fmt.Errorf("could not pull image manifest: %w", err)
	}

	localTag, err := name.NewTag(tag)
	if err != nil {
		return nil, fmt.Errorf("could not create local tag: %w", err)
	}

	_, err = daemon.Write(localTag, image)
	if err != nil {
		return nil, fmt.Errorf("could not write image to docker daemon: %w", err)
	}

	return image, nil
}

func setLabelOnImage(img v1.Image, key, value string) (v1.Image, error) {
	configFile, err := img.ConfigFile()
	if err != nil {
		return img, fmt.Errorf("could not get config file: %w", err)
	}

	if configFile.Config.Labels == nil {
		configFile.Config.Labels = make(map[string]string)
	}
	configFile.Config.Labels[key] = value

	newImage, err := mutate.ConfigFile(img, configFile)
	if err != nil {
		return img, fmt.Errorf("could not mutate config file: %w", err)
	}

	return newImage, nil
}
