package main

import (
	"sort"
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"github.com/docopt/docopt-go"
	"github.com/Masterminds/semver"
)

func main() {

	usage := `
Usage:
  chart_checker <chart> <index>
  chart_checker (-h | --help)
  chart_checker --version

Options:
  -h --help     Show this screen.
  --version     Show version.

Error Code:

  1				Validation failure
  126			Format error (the yaml or version can't be parsed)
`

	arguments, _ := docopt.ParseArgs(usage, os.Args[1:], "1.0.0")
	chart := arguments["<chart>"].(string)
	index := arguments["<index>"].(string)

	parsedChart, err := readChart(chart)
	if err != nil {
		log.Print("Error parsing chart")
		log.Print(err)
		// we took the exit code forme https://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(126)
	}

	appName := parsedChart.Name
	chartVersion, err := semver.NewVersion(parsedChart.Version)
	if err != nil {
		log.Print("Error parsing index")
		log.Print(err)
		// we took the exit code forme https://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(126)
	}

	latestInIndex, err := getLatest(appName, index)
	if err != nil {
		log.Print(err)
		// we took the exit code forme https://tldp.org/LDP/abs/html/exitcodes.html
		os.Exit(126)
	}

	if latestInIndex.Compare(chartVersion) >= 0 { // Chart greatre than latest
		msg := fmt.Sprintf("The chart version (%#v) is less than or equals the latest on index (%#v)", chartVersion, latestInIndex)
		log.Print(msg)
		os.Exit(1)
	}

}

type Chart struct {
	Name string `yaml:name`
	Version string `yaml:version`
}

func readChart(chartPath string) (*Chart, error){
	buf, err := ioutil.ReadFile(chartPath)
	if err != nil {
		return nil, err
	}

	c := &Chart{}
	err = yaml.Unmarshal(buf, c)

	return c, err
}


type IndexItem struct {
	Version string `yaml:version`
}

type Index struct {
	Entries map[string][]IndexItem
}

func getLatest(appName string, index string) (*semver.Version, error){
	indexParsed, err := readIndex(index)
	if err != nil {
		return nil, err
	}

	vs := make([]*semver.Version, len(indexParsed.Entries[appName]))
	for i, r := range indexParsed.Entries[appName] {
		v, err := semver.NewVersion(r.Version)
		if err != nil {
			return nil, err
		}
		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	return vs[len(vs)-1], err
}

func readIndex(indexPath string) (*Index, error) {
	buf, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	i := &Index{}
	err = yaml.Unmarshal(buf, i)

	return i, err
}
