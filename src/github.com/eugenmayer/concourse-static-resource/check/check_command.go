package main

import (
	"encoding/json"
	"os"

	"github.com/eugenmayer/concourse-static-resource/model"
)

func main() {
	// no-op check
	version := model.PseudoVersion{}
	version.Name = "Pseudo Version"
	version.VersionID = "123456"
	json.NewEncoder(os.Stdout).Encode([]model.PseudoVersion{{
		Name:      "Pseudo Version",
		VersionID: "123456",
	}})
}
