package main

import (
	"encoding/json"
	"os"

	"github.com/eugenmayer/concourse-static-resource/model"
	"github.com/eugenmayer/concourse-static-resource/log"
	"errors"
)

func main() {
	var request model.CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatal("reading request", err)
	}

	if request.Source.VersionStatic == "" {
		log.Fatal("Accessing version_static from source", errors.New("please provide a source.version_static value for the version"))
	}

	json.NewEncoder(os.Stdout).Encode( []model.Version{ model.Version{ Ref: request.Source.VersionStatic}} )
}
