package main

import (
	"flag"
	"fmt"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/image"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/packages"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
	"log"
	"os"
)

func main() {
	var (
		buildDestination string
		runDestination   string
		version          string
		stackName        string
		stacksDir        string
		publish          bool
	)

	flag.StringVar(&buildDestination, "build-destination", "", "Destination to tag and publish base image to")
	flag.StringVar(&runDestination, "run-destination", "", "Destination to tag and publish run image to")
	flag.StringVar(&version, "version", "", "Version to include in image tags")
	flag.StringVar(&stackName, "stack", "", "Stack name (base, full)")
	flag.StringVar(&stacksDir, "stacks-dir", "", "Stacks Base Directory")
	flag.BoolVar(&publish, "publish", false, "Push to docker registry")

	flag.Parse()

	if buildDestination == "" || runDestination == "" || version == "" || stackName == "" || stacksDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	buildBaseTag := fmt.Sprintf("%s:%s-%s", buildDestination, version, stackName)
	runBaseTag := fmt.Sprintf("%s:%s-%s", runDestination, version, stackName)

	var bionicStack stack.Stack
	var err error

	if stackName == "full" {
		bionicStack, err = stack.NewFullStack(stacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if stackName == "base" {
		bionicStack, err = stack.NewBaseStack(stacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Only full/base stacks supported")
	}

	creator := stack.Creator{
		PackageFinder:   packages.Bionic{},
		MixinsGenerator: stack.Mixins{},
		ImageClient:     image.Client{},
	}

	err = creator.CreateBionicStack(bionicStack, buildBaseTag, runBaseTag, publish)
	if err != nil {
		log.Fatal(err)
	}
}
