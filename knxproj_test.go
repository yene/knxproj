package knxproj_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/yene/knxproj"
)

const tmpFolder = "temp/"

var update = flag.Bool("update", false, "update golden file")

func TestAll(t *testing.T) {
	knxprojFile := "test-fixtures/DemoCase.knxproj"

	err := knxproj.Unzip(knxprojFile, tmpFolder)
	if err != nil {
		t.Error("unzip error", err)
	}

	knxproj.Language = "en-US"
	knxproj.ParseManufacturers(tmpFolder)
	knxproj.ReadManufacturerData(tmpFolder)
	knxproj.ReadProjects(tmpFolder)
	os.RemoveAll(tmpFolder)
	rankingsJSON, _ := json.MarshalIndent(knxproj.MainProject, "", "  ")

	if *update {
		ioutil.WriteFile("test-fixtures/DemoCase.json.golden", rankingsJSON, 0644)
	}

	expected, _ := ioutil.ReadFile("test-fixtures/DemoCase.json.golden")
	if !bytes.Equal(rankingsJSON, expected) {
		t.Error("Generated JSON is different from DemoCase.json.golden")
	}
}
