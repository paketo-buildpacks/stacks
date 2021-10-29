package stack

import "fmt"

const version = "bionic"

func GetUbuntuImage(architecture string) string {
	var image string
	label := version // can be overwitten below if it defers for a given architecture
	switch architecture {
	case "arm64":
		image = "arm64v8/ubuntu"
	default:
		image = "ubuntu"
	}
	return fmt.Sprintf("%s:%s", image, label)
}
