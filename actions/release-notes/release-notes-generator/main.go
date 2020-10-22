package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type USN struct {
	Title            string   `json:"title"`
	Link             string   `json:"link"`
	CveArray         []CVE    `json:"cves"`
	AffectedPackages []string `json:"affected_packages"`
}

type CVE struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

type RecordedUSN struct {
	Title   string `json:"title"`
	Release string `json:"release"`
}

func main() {
	var (
		buildBaseImage    string
		buildBaseImageTag string
		buildCNBImage     string
		runBaseImage      string
		runCNBImage       string
		buildReceiptDiff  string
		runReceiptDiff    string
		relevantUSNs      string
		allUSNs           string
		releaseVersion    string
		stack             string
	)

	flag.StringVar(&buildBaseImage, "build-base-image", "", "Fully qualified build base image")
	flag.StringVar(&buildBaseImageTag, "build-base-image-tag", "", "Build base image tag")
	flag.StringVar(&buildCNBImage, "build-cnb-image", "", "Fully qualified build CNB image")
	flag.StringVar(&runBaseImage, "run-base-image", "", "Fully qualified run base image")
	flag.StringVar(&runCNBImage, "run-cnb-image", "", "Fully qualified run CNB image")
	flag.StringVar(&buildReceiptDiff, "build-receipt-diff", "", "Build receipt diff")
	flag.StringVar(&runReceiptDiff, "run-receipt-diff", "", "Run receipt diff")
	flag.StringVar(&relevantUSNs, "relevant-usns", "", "Path to relevant USNs")
	flag.StringVar(&allUSNs, "all-usns", "", "Path to all USNs")
	flag.StringVar(&releaseVersion, "release-version", "", "Release version")
	flag.StringVar(&stack, "stack", "", "Stack")

	flag.Parse()

	if buildCNBImage == "" || runBaseImage == "" || runCNBImage == "" || relevantUSNs == "" || allUSNs == "" ||
		releaseVersion == "" || stack == "" {
		flag.Usage()
		os.Exit(1)
	}

	digestNotes := documentDigests(runBaseImage, runCNBImage, buildBaseImage, buildCNBImage, releaseVersion, stack, buildBaseImageTag)
	receiptNotes := documentReceiptDiffs(buildReceiptDiff, runReceiptDiff)
	usnNotes, err := documentUSNs(relevantUSNs, allUSNs, buildReceiptDiff, runReceiptDiff, releaseVersion)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error documenting USNs: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(digestNotes + usnNotes + receiptNotes)
}

func documentDigests(runBaseImage, runCNBImage, buildBaseImage, buildCNBImage, releaseVersion, stack, buildBaseImageTag string) string {
	runCNBImageTag := fmt.Sprintf("%s:%s-%s-cnb", strings.Split(runCNBImage, "@")[0], releaseVersion, stack)
	runBaseImageTag := fmt.Sprintf("%s:%s-%s", strings.Split(runBaseImage, "@")[0], releaseVersion, stack)
	buildCNBImageTag := fmt.Sprintf("%s:%s-%s-cnb", strings.Split(buildCNBImage, "@")[0], releaseVersion, stack)

	if buildBaseImageTag == "" {
		buildBaseImageTag = fmt.Sprintf("%s:%s-%s", strings.Split(buildBaseImage, "@")[0], releaseVersion, stack)
	}

	digestNotes := "## Release Images\n\n" +
		"### Runtime Base Images\n\n" +
		"#### For CNB Builds:\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", runCNBImageTag) +
		fmt.Sprintf("- Digest: `%s`\n\n", runCNBImage) +
		"#### Source Image (e.g., for Dockerfile builds):\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", runBaseImageTag) +
		fmt.Sprintf("- Digest: `%s`\n\n", runBaseImage) +
		"### Build-time Base Images\n\n" +
		"#### For CNB Builds:\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", buildCNBImageTag) +
		fmt.Sprintf("- Digest: `%s`\n\n", buildCNBImage) +
		"#### Source Image (e.g., for Dockerfile builds):\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", buildBaseImageTag) +
		fmt.Sprintf("- Digest: `%s`", buildBaseImage)

	return digestNotes
}

func documentReceiptDiffs(buildReceiptDiff, runReceiptDiff string) string {

	var receiptNotes string

	if buildReceiptDiff != "" {
		receiptNotes = receiptNotes + fmt.Sprintf("\n## Build Receipt Diff\n```\n%s\n```", formatReceiptDiff(buildReceiptDiff))
	}

	if runReceiptDiff != "" {
		receiptNotes = receiptNotes + fmt.Sprintf("\n## Run Receipt Diff\n```\n%s\n```", formatReceiptDiff(runReceiptDiff))
	}

	return receiptNotes
}

func getLongestItemLength(receiptDiff string, columnNumber int) int {
	longestItemLength := 0
	for _, line := range strings.Split(receiptDiff, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		packageName := fields[columnNumber]

		if len(packageName) > longestItemLength {
			longestItemLength = len(packageName)
		}
	}

	return longestItemLength
}

func formatReceiptDiff(receiptDiff string) string {
	longestPackageNameLength := getLongestItemLength(receiptDiff, 1)
	longestVersionLength := getLongestItemLength(receiptDiff, 2)
	longestArchLength := getLongestItemLength(receiptDiff, 3)

	var output string
	for _, line := range strings.Split(receiptDiff, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		descriptionPadding := "  "
		if len(fields) < 5 {
			descriptionPadding = ""
		}

		formatString := fmt.Sprintf("%%s  %%-%ds  %%-%ds  %%-%ds%s%%s\n", longestPackageNameLength, longestVersionLength, longestArchLength, descriptionPadding)

		output += fmt.Sprintf(formatString, fields[0], fields[1], fields[2], fields[3], strings.Join(fields[4:], " "))
	}

	return strings.TrimSpace(output)
}

func documentUSNs(relevantUSNsPath, allUSNsPath, buildReceiptDiff, runReceiptDiff, releaseVersion string) (string, error) {
	relevantUSNBytes, err := ioutil.ReadFile(relevantUSNsPath)
	if err != nil {
		return "", fmt.Errorf("error reading relevant USNs file: %w", err)
	}

	var relevantUSNs []RecordedUSN
	err = json.Unmarshal(relevantUSNBytes, &relevantUSNs)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling relevant USNs: %w", err)
	}

	allUSNBytes, err := ioutil.ReadFile(allUSNsPath)
	if err != nil {
		return "", fmt.Errorf("error reading all USNs file: %w", err)
	}

	var allUSNs []USN
	err = json.Unmarshal(allUSNBytes, &allUSNs)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling all USNs: %w", err)
	}

	allPkgs := append(getNewPackages(buildReceiptDiff), getNewPackages(runReceiptDiff)...)

	var patchedUSNs []USN
	for i, usn := range relevantUSNs {
		if usn.Release == "unreleased" {
			fullUSN := patchedUSN(usn, allUSNs, allPkgs)
			if fullUSN.Title != "" {
				patchedUSNs = append(patchedUSNs, fullUSN)
				relevantUSNs[i].Release = releaseVersion
			}
		}
	}

	err = updateReleasedUSNs(relevantUSNs, relevantUSNsPath)
	if err != nil {
		return "", fmt.Errorf("error updating relevant USNs with release: %w", err)
	}

	return generateUSNNotes(patchedUSNs), nil
}

func updateReleasedUSNs(usns []RecordedUSN, usnPath string) error {
	usnBytes, err := json.Marshal(usns)
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

func patchedUSN(recordedUSN RecordedUSN, allUSNs []USN, pkgs []string) USN {
	for _, usn := range allUSNs {
		if usn.Title == recordedUSN.Title {
			if patchedPkg(usn.AffectedPackages, pkgs) {
				return usn
			}
		}
	}

	return USN{}
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

func generateUSNNotes(usns []USN) string {
	if len(usns) == 0 {
		return ""
	}

	notes := "\n## Patched USNs"
	for _, usn := range usns {
		notes = notes + "\n" + formatUSN(usn) + "\n"
	}

	return notes
}

func formatUSN(usn USN) string {
	usnTitle := strings.Split(usn.Title, ":")
	usnNumber := usnTitle[0]
	usnDescription := usnTitle[1]
	notes := fmt.Sprintf("[%s](%s): %s", usnNumber, usn.Link, usnDescription)

	for _, cve := range usn.CveArray {
		notes = notes + fmt.Sprintf("\n* [%s](%s): %s", cve.Title, cve.Link, cve.Description)
	}

	return notes
}
