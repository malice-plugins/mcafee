package main

import "encoding/xml"

const tpl = `#### McAfee
{{- with .Results }}
| Infected      | Result      | Engine      | Updated      |
|:-------------:|:-----------:|:-----------:|:------------:|
| {{.Infected}} | {{.Result}} | {{.Engine}} | {{.Updated}} |
{{ end -}}
`

type productName struct {
	Value string `xml:"value,attr"`
}

type mVersion struct {
	Value string `xml:"value,attr"`
}

type licenseInfo struct {
	Value string `xml:"value,attr"`
}

type aVEngineVersion struct {
	Value string `xml:"value,attr"`
}

type datSetVersion struct {
	Value string `xml:"value,attr"`
}

type preamble struct {
	XMLName         xml.Name        `xml:"Preamble"`
	ProductName     productName     `xml:"Product_name"`
	Version         mVersion        `xml:"Version"`
	LicenseInfo     licenseInfo     `xml:"License_info"`
	AVEngineVersion aVEngineVersion `xml:"AV_Engine_version"`
	DatSetVersion   datSetVersion   `xml:"Dat_set_version"`
}

type scanDateTime struct {
	Value string `xml:"value,attr"`
}

type scanOptions struct {
	Value string `xml:"value,attr"`
}

type fileResults struct {
	Name          string `xml:"name,attr"`
	Status        string `xml:"status,attr"`
	VirusName     string `xml:"virus-name,attr"`
	DetectionType string `xml:"detection-type,attr"`
}

type timeToScan struct {
	Value string `xml:"value,attr"`
}

// McAfeeResults is the xml data struct
type McAfeeResults struct {
	XMLName  xml.Name     `xml:"Uvscan"`
	Preamble preamble     `xml:"Preamble"`
	DateTime scanDateTime `xml:"Date_Time"`
	Options  scanOptions  `xml:"Options"`
	File     fileResults  `xml:"File"`
	Time     timeToScan   `xml:"Time"`
}
