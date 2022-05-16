package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
)

func writeToTemp(w io.Writer) {
	fmt.Printf("Try writing to file... ")
	content := []byte("temporary file's content")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		fmt.Printf("FAIL\n")
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		fmt.Printf("FAIL\n")
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Printf("FAIL\n")
		log.Fatal(err)
	}
	fmt.Fprintln(w, tmpfile.Name())
	fmt.Printf("PASS\n")
}

func httpsClient() {
	fmt.Printf("Try connecting via HTTPS... ")
	_, err := http.Get("https://example.com/")
	if err != nil {
		fmt.Printf("FAIL\n")
		log.Fatal(err)
	}
	fmt.Printf("PASS\n")

	fmt.Printf("Try connecting to an untrusted HTTPS... ")
	_, err = http.Get("https://untrusted-root.badssl.com/")
	if err == nil {
		fmt.Printf("FAIL\n")
		log.Fatalf("Expected error when connecting to https://untrusted-root.badssl.com/")
	}
	fmt.Printf("PASS\n")
}

func checkUsers() {
	fmt.Printf("Try checking users... ")
	for _, username := range []string{"root","nonroot","nobody"} {
		_, err := user.Lookup(username)
		if err != nil {
			fmt.Printf("FAIL\n")
			log.Fatal(err)
		}
	}
	fmt.Printf("PASS\n")
}

func checkGroups() {
	fmt.Printf("Try checking groups... ")
	for _, groupname := range []string{"root","nonroot","nobody","tty","staff"} {
		_, err := user.LookupGroup(groupname)
		if err != nil {
			fmt.Printf("FAIL\n")
			log.Fatal(err)
		}
	}
	fmt.Printf("PASS\n")
}

func checkServices() {
	fmt.Printf("Try checking services... ")
	_, err := net.LookupPort("tcp", "ldap")
	if err != nil {
		fmt.Printf("FAIL\n")
		log.Fatal(err)
	}
	fmt.Printf("PASS\n")
}

func main() {
	writeToTemp(bytes.NewBuffer([]byte("foo")))
	httpsClient()
	checkUsers()
	checkGroups()
	checkServices()
	fmt.Println("All tests PASSED")
}
