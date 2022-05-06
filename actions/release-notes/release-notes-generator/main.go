package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
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
	var opts struct {
		BuildBaseImage   string `long:"build-base-image" description:"Fully qualified build base image"`
		BuildCNBImage    string `long:"build-cnb-image" description:"Fully qualified build CNB image" required:"true"`
		RunBaseImage     string `long:"run-base-image" description:"Fully qualified run base image"`
		RunCNBImage      string `long:"run-cnb-image" description:"Fully qualified run CNB image" required:"true"`
		BuildReceiptDiff string `long:"build-receipt-diff" description:"Build receipt diff"`
		RunReceiptDiff   string `long:"run-receipt-diff" description:"Run receipt diff"`
		RelevantUSNs     string `long:"relevant-usns" description:"Path to relevant USNs" required:"true"`
		AllUSNs          string `long:"all-usns" description:"Path to all USNs" required:"true"`
		ReleaseVersion   string `long:"release-version" description:"Release version" required:"true"`
		Stack            string `long:"stack" description:"Stack" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	digestNotes := documentDigests(opts.RunBaseImage, opts.RunCNBImage, opts.BuildBaseImage, opts.BuildCNBImage, opts.ReleaseVersion, opts.Stack)
	receiptNotes := documentReceiptDiffs(opts.BuildReceiptDiff, opts.RunReceiptDiff)
	usnNotes, err := documentUSNs(opts.RelevantUSNs, opts.AllUSNs, opts.ReleaseVersion)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error documenting USNs: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(digestNotes + usnNotes + receiptNotes)
}

func documentDigests(runBaseImage, runCNBImage, buildBaseImage, buildCNBImage, releaseVersion, stack string) string {
	runCNBImageTag := fmt.Sprintf("%s:%s-%s-cnb", strings.Split(runCNBImage, "@")[0], releaseVersion, stack)
	runBaseImageTag := fmt.Sprintf("%s:%s-%s", strings.Split(runBaseImage, "@")[0], releaseVersion, stack)
	buildCNBImageTag := fmt.Sprintf("%s:%s-%s-cnb", strings.Split(buildCNBImage, "@")[0], releaseVersion, stack)
	buildBaseImageTag := fmt.Sprintf("%s:%s-%s", strings.Split(buildBaseImage, "@")[0], releaseVersion, stack)

	digestNotes := "## Release Images\n\n" +
		"### Runtime Base Images\n\n" +
		"#### For CNB Builds:\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", runCNBImageTag) +
		fmt.Sprintf("- Digest: `%s`\n\n", runCNBImage)

	if runBaseImage != "" {
		digestNotes += "#### Source Image (e.g., for Dockerfile builds):\n" +
			fmt.Sprintf("- Tag: **`%s`**\n", runBaseImageTag) +
			fmt.Sprintf("- Digest: `%s`\n\n", runBaseImage)
	}

	digestNotes += "### Build-time Base Images\n\n" +
		"#### For CNB Builds:\n" +
		fmt.Sprintf("- Tag: **`%s`**\n", buildCNBImageTag) +
		fmt.Sprintf("- Digest: `%s`", buildCNBImage)

	if buildBaseImage != "" {
		digestNotes += "\n\n#### Source Image (e.g., for Dockerfile builds):\n" +
			fmt.Sprintf("- Tag: **`%s`**\n", buildBaseImageTag) +
			fmt.Sprintf("- Digest: `%s`", buildBaseImage)
	}

	return digestNotes
}

func documentReceiptDiffs(buildReceiptDiff, runReceiptDiff string) string {
	var receiptNotes string

	if buildReceiptDiff != "" {
		receiptNotes = receiptNotes + fmt.Sprintf("\n\n## Build Receipt Diff\n```\n%s\n```", formatReceiptDiff(buildReceiptDiff))
	}

	if runReceiptDiff != "" {
		receiptNotes = receiptNotes + fmt.Sprintf("\n\n## Run Receipt Diff\n```\n%s\n```", formatReceiptDiff(runReceiptDiff))
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

func documentUSNs(relevantUSNsPath, allUSNsPath, releaseVersion string) (string, error) {
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

	var patchedUSNs []USN
	for _, usn := range relevantUSNs {
		if usn.Release == releaseVersion {
			fullUSN := patchedUSN(usn, allUSNs)
			if fullUSN.Title != "" {
				patchedUSNs = append(patchedUSNs, fullUSN)
			}
		}
	}

	return generateUSNNotes(patchedUSNs), nil
}

func patchedUSN(recordedUSN RecordedUSN, allUSNs []USN) USN {
	for _, usn := range allUSNs {
		if usn.Title == recordedUSN.Title {
			return usn
		}
	}

	return USN{}
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
