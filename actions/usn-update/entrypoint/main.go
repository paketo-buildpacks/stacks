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
	CVEs             []CVE    `json:"cves"`
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
		buildReceiptPath    string
		runReceiptPath      string
		fullUSNListPath     string
		relevantUSNListPath string
	)

	flag.StringVar(&buildReceiptPath, "build-receipt", "", "Path to build receipt")
	flag.StringVar(&runReceiptPath, "run-receipt", "", "Path to run receipt")
	flag.StringVar(&fullUSNListPath, "full-usn-list", "", "Path to full USN list")
	flag.StringVar(&relevantUSNListPath, "relevant-usn-list", "", "Path to relevant USN list")

	flag.Parse()

	if buildReceiptPath == "" || runReceiptPath == "" || fullUSNListPath == "" || relevantUSNListPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := recordRelevantUSNs(buildReceiptPath, runReceiptPath, fullUSNListPath, relevantUSNListPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error checking USNs: %s\n", err.Error())
		os.Exit(1)
	}
}

func recordRelevantUSNs(buildReceiptPath, runReceiptPath, fullUSNListPath, relevantUSNListPath string) error {
	allUSNs, err := getUSNs(fullUSNListPath)
	if err != nil {
		return fmt.Errorf("failed to get full USN list: %w", err)
	}

	allPkgs, err := getAllPkgs(buildReceiptPath, runReceiptPath)
	if err != nil {
		return fmt.Errorf("failed to get packages: %w", err)
	}

	recordedUSNs, err := getRecordedUSNs(relevantUSNListPath)
	if err != nil {
		return fmt.Errorf("failed to get recorded USNs: %w", err)
	}

	var relevantUSNs []USN
	for _, usn := range allUSNs {
		if isRelevantUSN(usn, allPkgs, recordedUSNs) {
			relevantUSNs = append(relevantUSNs, usn)
		}
	}

	err = recordUSNs(relevantUSNs, relevantUSNListPath)
	if err != nil {
		return fmt.Errorf("failed to record USNs: %w", err)
	}

	return nil
}

func getAllPkgs(buildReceiptPath, runReceiptPath string) ([]string, error) {
	buildPkgList, err := getPkgList(buildReceiptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get build package list: %w", err)
	}

	runPkgList, err := getPkgList(runReceiptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get run package list: %w", err)
	}

	fullPkgList := append(buildPkgList, runPkgList...)

	return fullPkgList, nil
}

func getUSNs(usnPath string) ([]USN, error) {
	contents, err := ioutil.ReadFile(usnPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read USN file: %w", err)
	}

	var usns []USN
	err = json.Unmarshal(contents, &usns)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal USNs: %w", err)
	}

	return usns, nil
}

func getRecordedUSNs(recordedUSNPath string) ([]RecordedUSN, error) {
	contents, err := ioutil.ReadFile(recordedUSNPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read USN file: %w", err)
	}

	var usns []RecordedUSN
	err = json.Unmarshal(contents, &usns)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal USNs: %w", err)
	}

	return usns, nil
}

func isRelevantUSN(usn USN, allPkgs []string, recordedUSNs []RecordedUSN) bool {
	for _, recordedUSN := range recordedUSNs {
		if recordedUSN.Title == usn.Title {
			return false
		}
	}

	for _, affectedPkg := range usn.AffectedPackages {
		for _, pkg := range allPkgs {
			if affectedPkg == pkg {
				return true
			}
		}
	}
	return false
}

func getPkgList(receiptPath string) ([]string, error) {
	contents, err := ioutil.ReadFile(receiptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read receipt file: %w", err)
	}
	formattedContents := strings.ReplaceAll(string(contents), "\n", " ")

	re := regexp.MustCompile(`(ii).*?(\S+)`)
	pkgList := re.FindAllStringSubmatch(formattedContents, -1)

	var finalPkgList []string
	for _, p := range pkgList {
		finalPkgList = append(finalPkgList, strings.Split(p[2], ":")[0])
	}

	return finalPkgList, nil
}

func recordUSNs(usns []USN, usnListPath string) error {
	contents, err := ioutil.ReadFile(usnListPath)
	if err != nil {
		return fmt.Errorf("failed to read USN list file: %w", err)
	}

	var recordedUSNs []RecordedUSN
	err = json.Unmarshal(contents, &recordedUSNs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal usns: %w", err)
	}

	for i := len(usns) - 1; i >= 0; i-- {
		usn := RecordedUSN{
			Title:   usns[i].Title,
			Release: "unreleased",
		}
		recordedUSNs = append([]RecordedUSN{usn}, recordedUSNs...)
	}

	usnBytes, err := json.Marshal(recordedUSNs)
	if err != nil {
		return fmt.Errorf("failed to marshal usns: %w", err)
	}

	newUSNFile, err := os.Create(usnListPath)
	if err != nil {
		return fmt.Errorf("failed to create USN list file: %w", err)
	}
	defer newUSNFile.Close()

	_, err = newUSNFile.Write(usnBytes)
	if err != nil {
		return fmt.Errorf("failed to write USN list file: %w", err)
	}
	return nil
}
