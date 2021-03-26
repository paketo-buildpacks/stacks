package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/image"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/packages"
	stackPkg "github.com/paketo-buildpacks/stacks/create-stack/pkg/stack"
)

func main() {
	var opts struct {
		BuildDestination string `long:"build-destination" description:"Destination to tag and publish base image to" required:"true"`
		RunDestination   string `long:"run-destination" description:"Destination to tag and publish run image to" required:"true"`
		Version          string `long:"version" description:"Version to include in image tags" required:"true"`
		StackName        string `long:"stack" description:"Stack name (base, full, tiny)" required:"true"`
		StacksDir        string `long:"stacks-dir" description:"Stacks Base Directory" required:"true"`
		Publish          bool   `long:"publish" description:"Push to docker registry"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	buildBaseTag := fmt.Sprintf("%s:%s-%s", opts.BuildDestination, opts.Version, opts.StackName)
	runBaseTag := fmt.Sprintf("%s:%s-%s", opts.RunDestination, opts.Version, opts.StackName)

	var stack stackPkg.Stack
	var packageFinder stackPkg.PackageFinder

	if opts.StackName == "full" {
		packageFinder = packages.Bionic{}
		stack, err = stackPkg.NewFullStack(opts.StacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if opts.StackName == "base" {
		packageFinder = packages.Bionic{}
		stack, err = stackPkg.NewBaseStack(opts.StacksDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if opts.StackName == "tiny" {
		packageFinder = packages.Tiny{BuildPkgs: packages.Bionic{}}
		stack, err = stackPkg.NewTinyStack(opts.StacksDir)
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

	err = creator.CreateStack(stack, buildBaseTag, runBaseTag, opts.Publish)
	if err != nil {
		log.Fatal(err)
	}
}
