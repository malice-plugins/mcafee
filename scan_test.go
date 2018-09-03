package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

// TestParseResult tests the __ function.
func TestParseResult(t *testing.T) {
	r, err := ioutil.ReadFile("tests/av_scan.out")
	if err != nil {
		fmt.Print(err)
	}

	result := strings.Split(string(r), "\t")

	if !strings.Contains(string(r), "[OK]") {
		if true {
			t.Log("Infected: ", true)
			t.Log("Result: ", strings.TrimSpace(result[1]))
		}
	}

}

// TestParseVersion tests the __ function.
func TestParseVersion(t *testing.T) {
	v, err := ioutil.ReadFile("tests/av_version.out")
	if err != nil {
		fmt.Print(err)
	}

	d, err := ioutil.ReadFile("tests/av_vps.out")
	if err != nil {
		fmt.Print(err)
	}

	version := strings.TrimSpace(string(v))
	database := strings.TrimSpace(string(d))

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
