package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eugenmayer/concourse-static-resource/model"
	"github.com/eugenmayer/concourse-static-resource/shared"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	if err := os.MkdirAll(destination, 0755); err != nil {
		log.Fatal("creating destination", err)
	}

	var request model.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatal("reading request", err)
	}

	// that is the version handed over from check, we do NOT use the Source.version_static argument here
	// eventhough that would be the same when having a check to in handover
	// though when having a out to in handover, it would be the version out had, maybe something else VersionFromFile
	var version string = request.Version.Ref

	var SourceURL string = shared.InjectVersionIntoPath(request.Source.URI, version, "<version>")
	URI, err := url.Parse(SourceURL)
	if err != nil {
		log.Fatal("parsing uri", err)
	}

	var curlOpts string = shared.Curlopt(request.Source)

	// placeholder for the curlPipe source arg
	curlOpts = curlOpts

	var command string
	if request.Source.Extract == true {
		// $2 is the placeholder for the curlPipe destination arg
		command = fmt.Sprintf("curl %s '%s' | tar --warning=no-unknown-keyword -C '%s' -zxf -", curlOpts, URI.String(), destination)
	} else {
		// $2 is the placeholder for the curlPipe destination arg
		command = fmt.Sprintf("cd '%s'; curl %s '%s' -O", destination, curlOpts, URI.String())
	}

	curlPipe := exec.Command(
		"sh",
		"-ec",
		command,
		"sh",
	)

	curlPipe.Stderr = os.Stderr

	err = curlPipe.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Url: "+URI.String())
		os.Exit(1)
	}

	metavalue := []model.MetaDataPair{
		model.MetaDataPair{
			Name: "filename",
			// we expect the filename to be tha last path snippet
			Value: filepath.Base(URI.String()),
		},
	}
	json.NewEncoder(os.Stdout).Encode(model.InResponse{
		Version:  model.Version{version},
		MetaData: metavalue,
	})
}
