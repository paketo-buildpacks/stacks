package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type USN struct {
	Title            string   `json:"title"`
	AffectedPackages []string `json:"affected_packages"`
}

type RecordedUSN struct {
	Title   string `json:"title"`
	Release string `json:"release"`
}

func main() {
	var opts struct {
		BuildReceiptDiff string `long:"build-receipt-diff" description:"Build receipt diff"`
		RunReceiptDiff   string `long:"run-receipt-diff" description:"Run receipt diff"`
		RelevantUSNs     string `long:"relevant-usns" description:"Path to relevant USNs" required:"true"`
		AllUSNs          string `long:"all-usns" description:"Path to all USNs" required:"true"`
		ReleaseVersion   string `long:"release-version" description:"Release version" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	err = updateUSNs(opts.RelevantUSNs, opts.AllUSNs, opts.BuildReceiptDiff, opts.RunReceiptDiff, opts.ReleaseVersion)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error updating USNs: %s\n", err.Error())
		os.Exit(1)
	}
}

func updateUSNs(relevantUSNsPath, allUSNsPath, buildReceiptDiff, runReceiptDiff, releaseVersion string) error {
	relevantUSNBytes, err := ioutil.ReadFile(relevantUSNsPath)
	if err != nil {
		return fmt.Errorf("error reading relevant USNs file: %w", err)
	}

	var relevantUSNs []RecordedUSN
	err = json.Unmarshal(relevantUSNBytes, &relevantUSNs)
	if err != nil {
		return fmt.Errorf("error unmarshalling relevant USNs: %w", err)
	}

	allUSNBytes, err := ioutil.ReadFile(allUSNsPath)
	if err != nil {
		return fmt.Errorf("error reading all USNs file: %w", err)
	}

	var allUSNs []USN
	err = json.Unmarshal(allUSNBytes, &allUSNs)
	if err != nil {
		return fmt.Errorf("error unmarshalling all USNs: %w", err)
	}

	allPkgs := append(getNewPackages(buildReceiptDiff), getNewPackages(runReceiptDiff)...)

	for i, usn := range relevantUSNs {
		if usn.Release == "unreleased" {
			if patchedUSN(usn, allUSNs, allPkgs) {
				relevantUSNs[i].Release = releaseVersion
			}
		}
	}

	err = updateReleasedUSNs(relevantUSNs, relevantUSNsPath)
	if err != nil {
		return fmt.Errorf("error updating relevant USNs with release: %w", err)
	}

	return nil
}

func updateReleasedUSNs(usns []RecordedUSN, usnPath string) error {
	usnBytes, err := json.MarshalIndent(usns, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling relevant USNs: %w", err)
	}

	usnFile, err := os.Create(usnPath)
	if err != nil {
		return fmt.Errorf("error creating new relevant USN file: %w", err)
	}
	defer usnFile.Close()

	_, err = usnFile.Write(usnBytes)
	if err != nil {
		return fmt.Errorf("error writing to relevant USN file: %w", err)
	}
	return nil
}

func getNewPackages(receiptDiff string) []string {
	re := regexp.MustCompile(`(?m:^\+ii\s*?)(\S+)`)
	pkgMatches := re.FindAllStringSubmatch(receiptDiff, -1)

	var pkgs []string
	for _, match := range pkgMatches {
		pkgs = append(pkgs, strings.Split(match[1], ":")[0])
	}

	return pkgs
}

func patchedUSN(recordedUSN RecordedUSN, allUSNs []USN, pkgs []string) bool {
	for _, usn := range allUSNs {
		if usn.Title == recordedUSN.Title {
			if patchedPkg(usn.AffectedPackages, pkgs) {
				return true
			}
		}
	}

	return false
}

func patchedPkg(affectedPkgs, allPkgs []string) bool {
	for _, affectedPkg := range affectedPkgs {
		for _, pkg := range allPkgs {
			if strings.TrimSpace(affectedPkg) == strings.TrimSpace(pkg) {
				return true
			}
		}
	}

	return false
}
