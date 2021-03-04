package stack

import (
	"fmt"
)

type Mixins struct{}

func (m Mixins) GetMixins(buildPackages, runPackages []string) ([]string, []string) {
	buildPackagesDiff := m.getDifference(buildPackages, runPackages)
	runPackagesDiff := m.getDifference(runPackages, buildPackages)

	buildMixins := m.tagMixins(buildPackagesDiff, "build")
	runMixins := m.tagMixins(runPackagesDiff, "run")

	sharedPackages := m.getIntersection(buildPackages, runPackages)
	sharedPackages2 := m.getIntersection(buildPackages, runPackages)

	buildMixins = append(sharedPackages, buildMixins...)
	runMixins = append(sharedPackages2, runMixins...)

	return buildMixins, runMixins
}

func (m Mixins) getDifference(slice1 []string, slice2 []string) []string {
	var diff []string

	for _, el := range slice1 {
		exist := false
		for _, el2 := range slice2 {
			if el == el2 {
				exist = true
				break
			}
		}

		if !exist {
			diff = append(diff, el)
		}
	}

	return diff
}

func (m Mixins) getIntersection(slice1 []string, slice2 []string) []string {
	var intersection []string

	for _, el := range slice1 {
		for _, el2 := range slice2 {
			if el == el2 {
				intersection = append(intersection, el)
			}
		}
	}

	return intersection
}

func (m Mixins) tagMixins(packages []string, phase string) []string {
	var mixins []string

	for _, el := range packages {
		mixins = append(mixins, fmt.Sprintf("%s:%s", phase, el))
	}

	return mixins
}
