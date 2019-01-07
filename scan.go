package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	name     = "mcafee"
	category = "av"
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string

	path string
	hash string
	// es is the elasticsearch database object
	es elasticsearch.Database
)

type pluginResults struct {
	ID   string      `json:"id" structs:"id,omitempty"`
	Data ResultsData `json:"mcafee" structs:"mcafee"`
}

// McAfee json object
type McAfee struct {
	Results ResultsData `json:"mcafee"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" structs:"infected"`
	Result   string `json:"result" structs:"result"`
	Engine   string `json:"engine" structs:"engine"`
	Database string `json:"database" structs:"database"`
	Updated  string `json:"updated" structs:"updated"`
	MarkDown string `json:"markdown,omitempty" structs:"markdown,omitempty"`
}

func assert(err error) {
	if err != nil {
		// skip exit code 13 (which means a virus was found)
		if err.Error() != "exit status 13" {
			log.WithFields(log.Fields{
				"plugin":   name,
				"category": category,
				"path":     path,
			}).Fatal(err)
		}
	}
}

// AvScan performs antivirus scan
func AvScan(timeout int) McAfee {

	defer os.Remove("/tmp/" + hash + ".xml")

	var results ResultsData

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, err := utils.RunCommand(ctx, "/usr/local/uvscan/uvscan_secure", path, "--xmlpath=/tmp/"+hash+".xml")
	assert(err)
	results, err = ParseMcAfeeOutput(output)

	if err != nil {
		// If fails try a second time
		output, err := utils.RunCommand(ctx, "/usr/local/uvscan/uvscan_secure", path, "--xmlpath=/tmp/"+hash+".xml")
		assert(err)
		results, err = ParseMcAfeeOutput(output)
		assert(err)
	}

	return McAfee{
		Results: results,
	}
}

// ParseMcAfeeOutput convert mcafee output into ResultsData struct
func ParseMcAfeeOutput(mcafeeout string) (ResultsData, error) {

	log.WithFields(log.Fields{
		"plugin":   name,
		"category": category,
		"path":     path,
	}).Debug("McAfee Output: ", mcafeeout)

	xmlFile, err := os.Open("/tmp/" + hash + ".xml")
	if err != nil {
		fmt.Print(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var results McAfeeResults
	err = xml.Unmarshal(byteValue, &results)

	mcafee := ResultsData{
		Infected: strings.EqualFold(results.File.Status, "infected"),
		Engine:   results.Preamble.AVEngineVersion.Value,
		Database: results.Preamble.DatSetVersion.Value,
		Updated:  getUpdatedDate(),
		Result:   results.File.VirusName,
	}

	return mcafee, nil
}

func parseUpdatedDate(date string) string {
	layout := "Mon, 02 Jan 2006 15:04:05 +0000"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return string(updated)
}

func updateAV(ctx context.Context) error {
	fmt.Println("Updating McAfee...")
	fmt.Println(utils.RunCommand(ctx, "/usr/local/uvscan/update"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func didLicenseExpire() bool {
	log.Error("could not find expiration date in license file")
	return false
}

func generateMarkDownTable(a McAfee) string {
	var tplOut bytes.Buffer

	t := template.Must(template.New("mcafee").Parse(tpl))

	err := t.Execute(&tplOut, a)
	if err != nil {
		log.Println("executing template:", err)
	}

	return tplOut.String()
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func webService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scan", webAvScan).Methods("POST")
	log.WithFields(log.Fields{
		"plugin":   name,
		"category": category,
	}).Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}

func webAvScan(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("malware")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Please supply a valid file to scan.")
		log.WithFields(log.Fields{
			"plugin":   name,
			"category": category,
		}).Error(err)
	}
	defer file.Close()

	log.WithFields(log.Fields{
		"plugin":   name,
		"category": category,
	}).Debug("Uploaded fileName: ", header.Filename)

	tmpfile, err := ioutil.TempFile("/malware", "web_")
	assert(err)
	defer os.Remove(tmpfile.Name()) // clean up

	data, err := ioutil.ReadAll(file)
	assert(err)

	if _, err = tmpfile.Write(data); err != nil {
		assert(err)
	}
	if err = tmpfile.Close(); err != nil {
		assert(err)
	}

	// Do AV scan
	path = tmpfile.Name()
	mcafee := AvScan(60)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mcafee); err != nil {
		assert(err)
	}
}

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()

	app.Name = "mcafee"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice McAfee AntiVirus Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "",
			Usage:       "elasticsearch url for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH_URL",
			Destination: &es.URL,
		},
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "callback, c",
			Usage:  "POST results back to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  120,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update virus definitions",
			Action: func(c *cli.Context) error {
				return updateAV(nil)
			},
		},
		{
			Name:  "web",
			Usage: "Create a McAfee scan web service",
			Action: func(c *cli.Context) error {
				webService()
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		var err error

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if c.Args().Present() {
			path, err = filepath.Abs(c.Args().First())
			assert(err)

			if _, err = os.Stat(path); os.IsNotExist(err) {
				assert(err)
			}

			// if didLicenseExpire() {
			// 	log.Errorln("mcafee license has expired")
			// 	log.Errorln("please get a new one here: http://www.mcafee.com/ca/about/contact-us.aspx")
			// }
			hash = utils.GetSHA256(path)

			mcafee := AvScan(c.Int("timeout"))
			mcafee.Results.MarkDown = generateMarkDownTable(mcafee)
			// upsert into Database
			if len(c.String("elasticsearch")) > 0 {
				err := es.Init()
				if err != nil {
					return errors.Wrap(err, "failed to initalize elasticsearch")
				}
				err = es.StorePluginResults(database.PluginResults{
					ID:       utils.Getopt("MALICE_SCANID", hash),
					Name:     name,
					Category: category,
					Data:     structs.Map(mcafee.Results),
				})
				if err != nil {
					return errors.Wrapf(err, "failed to index malice/%s results", name)
				}
			}

			if c.Bool("table") {
				fmt.Printf(mcafee.Results.MarkDown)
			} else {
				mcafee.Results.MarkDown = ""
				mcafeeJSON, err := json.Marshal(mcafee)
				assert(err)
				if c.Bool("callback") {
					request := gorequest.New()
					if c.Bool("proxy") {
						request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
					}
					request.Post(os.Getenv("MALICE_ENDPOINT")).
						Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", utils.GetSHA256(path))).
						Send(string(mcafeeJSON)).
						End(printStatus)

					return nil
				}
				fmt.Println(string(mcafeeJSON))
			}
		} else {
			log.WithFields(log.Fields{
				"plugin":   name,
				"category": category,
			}).Fatal(fmt.Errorf("Please supply a file to scan with malice/%s", name))
		}
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
