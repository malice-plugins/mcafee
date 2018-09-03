package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TestParseResult tests the __ function.
func TestParseMalwareResultXML(t *testing.T) {
	xmlFile, err := os.Open("tests/av_malware.xml")
	if err != nil {
		fmt.Print(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var r McAfeeResults
	err = xml.Unmarshal(byteValue, &r)

	if strings.EqualFold(r.File.Status, "infected") {
		if true {
			t.Log("Infected: ", r.File.Status)
			t.Log("Result: ", strings.TrimSpace(r.File.VirusName))
		}
	}
}

func TestParseCleanResultXML(t *testing.T) {
	xmlFile, err := os.Open("tests/av_clean.xml")
	if err != nil {
		fmt.Print(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var r McAfeeResults
	err = xml.Unmarshal(byteValue, &r)

	if strings.EqualFold(r.File.Status, "") {
		r.File.Status = "clean"
		if true {
			t.Log("SHIZ IS CLEAN YO!")
			t.Log("Infected: ", r.File.Status)
			t.Log("Result: ", strings.TrimSpace(r.File.VirusName))
		}
	}
}
