package main

import (
	"encoding/json"
	"errors"
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
		log.Fatal("accessing first argument", errors.New("usage: %s <sources directory>"))
	}
	var sourceDir = os.Args[1]

	var request model.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatal("reading request", err)
	}

	var curlOpts string = shared.Curlopt(request.Source)

	// read the version if the path is actually provided
	var version string

	// if no version should be read from a file, use the version handed over by the ref
	if request.Params.VersionFilepath != "" {
		version = shared.GetVersionFromFile(request.Params.VersionFilepath, sourceDir)
	} else {
		version = request.Version.Ref
	}

	// depending if destFilenamePattern has a placeholder, us version to replace it and set our destFilename
	var destURL string = shared.InjectVersionIntoPath(request.Source.URI, version, "<version>")
	URI, err := url.Parse(destURL)
	if err != nil {
		log.Fatal("parsing uri", err)
	}

	// resolve our glob so we have the source file
	var sourceFile string = shared.GetSourceFile(request.Params.SourceFilepathGlob, sourceDir)

	// placeholder for the curlPipe dest arg $1 and upload-destination $2
	var command = fmt.Sprintf("curl %s --upload-file '%s' '%s'", curlOpts, sourceFile, URI.String())

	curlPipe := exec.Command(
		"sh",
		"-ec",
		command,
		"sh",
	)

	//curlPipe.Stdout = os.Stderr
	curlPipe.Stderr = os.Stderr

	if err := curlPipe.Run(); err != nil {
		if request.Source.Debug {
			fmt.Fprintln(os.Stderr, command)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Url: "+URI.String())
		log.Fatal("uploading file", err)
	}

	metavalue := []model.MetaDataPair{
		model.MetaDataPair{
			Name: "filename",
			// we expect the filename to be tha last path snippet
			Value: filepath.Base(URI.String()),
		},
	}
	json.NewEncoder(os.Stdout).Encode(model.OutResponse{
		Version:  model.Version{version},
		MetaData: metavalue,
	})
}
