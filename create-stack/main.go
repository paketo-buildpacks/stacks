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
		RunDestination   string `long:"run-destination"   description:"Destination to tag and publish run image to"  required:"true"`
		Version          string `long:"version"           description:"Version to include in image tags"             required:"true"`
		StackName        string `long:"stack"             description:"Stack name (base, full, tiny)"                required:"true"`
		StacksDir        string `long:"stacks-dir"        description:"Stacks Base Directory"                        required:"true"`
		Publish          bool   `long:"publish"           description:"Push to docker registry"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	buildBaseTag := fmt.Sprintf("%s:%s-%s", opts.BuildDestination, opts.Version, opts.StackName)
	runBaseTag := fmt.Sprintf("%s:%s-%s", opts.RunDestination, opts.Version, opts.StackName)

	var definition stackPkg.Definition
	var packageFinder stackPkg.PackageFinder

	switch opts.StackName {
	case "full":
		packageFinder = packages.Bionic{}
		definition, err = stackPkg.NewFullStack(buildBaseTag, runBaseTag, opts.StacksDir, opts.Publish)
		if err != nil {
			log.Fatal(err)
		}

	case "base":
		packageFinder = packages.Bionic{}
		definition, err = stackPkg.NewBaseStack(buildBaseTag, runBaseTag, opts.StacksDir, opts.Publish)
		if err != nil {
			log.Fatal(err)
		}

	case "tiny":
		packageFinder = packages.Tiny{BuildPkgs: packages.Bionic{}}
		definition, err = stackPkg.NewTinyStack(buildBaseTag, runBaseTag, opts.StacksDir, opts.Publish)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Only full, base, and tiny stacks supported")
	}

	_, attachBOM := os.LookupEnv("EXPERIMENTAL_ATTACH_RUN_IMAGE_SBOM")
	creator := stackPkg.Creator{
		PackageFinder:   packageFinder,
		MixinsGenerator: stackPkg.Mixins{},
		ImageClient:     image.Client{},
		BOMGenerator:    stackPkg.BOM{},
		AttachBOM:       attachBOM,
	}

	err = creator.Execute(definition)
	if err != nil {
		log.Fatal(err)
	}
}
