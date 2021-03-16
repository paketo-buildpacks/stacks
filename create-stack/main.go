package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paketo-buildpacks/stacks/create-stack/pkg/image"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/packages"
	stackPkg "github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
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
	flag.StringVar(&stackName, "stack", "", "Stack name (base, full, tiny)")
	flag.StringVar(&stacksDir, "stacks-dir", "", "Stacks Base Directory")
	flag.BoolVar(&publish, "publish", false, "Push to docker registry")

	flag.Parse()

	if buildDestination == "" || runDestination == "" || version == "" || stackName == "" || stacksDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	buildBaseTag := fmt.Sprintf("%s:%s-%s", buildDestination, version, stackName)
	runBaseTag := fmt.Sprintf("%s:%s-%s", runDestination, version, stackName)

	var stack stackPkg.Stack
	var packageFinder stackPkg.PackageFinder
	var err error

	if stackName == "full" {
		packageFinder = packages.Bionic{}
		stack, err = stackPkg.NewFullStack(stacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if stackName == "base" {
		packageFinder = packages.Bionic{}
		stack, err = stackPkg.NewBaseStack(stacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if stackName == "tiny" {
		packageFinder = packages.Tiny{BuildPkgs: packages.Bionic{}}
		stack, err = stackPkg.NewTinyStack(stacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Only full, base, and tiny stacks supported")
	}

	creator := stackPkg.Creator{
		PackageFinder:   packageFinder,
		MixinsGenerator: stackPkg.Mixins{},
		ImageClient:     image.Client{},
	}

	err = creator.CreateStack(stack, buildBaseTag, runBaseTag, publish)
	if err != nil {
		log.Fatal(err)
	}
}
