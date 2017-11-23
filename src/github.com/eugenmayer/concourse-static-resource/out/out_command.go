package main

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"

	"github.com/eugenmayer/concourse-static-resource/curlopts"
	"github.com/eugenmayer/concourse-static-resource/log"
	"github.com/eugenmayer/concourse-static-resource/model"
)

func main() {
	var request model.OutRequest
	var destUrl *url.URL
	var err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatal("reading request", err)
	}

	destUrl, err = url.Parse(request.Source.URI)

	if err != nil {
		log.Fatal("parsing uri", err)
	}

	var curlOpts string = curlopts.Curlopt(request.Source)

	// placeholder for the curlPipe dest arg $1 and upload-destination $2
	var command string = "curl " + curlOpts + " --upload-file \"$1\" \"$2\""

	curlPipe := exec.Command(
		"sh",
		"-exc",
		command,
		"sh", request.Params.Filepath, destUrl.String(),
	)


	curlPipe.Stdout = os.Stderr
	curlPipe.Stderr = os.Stderr

	err = curlPipe.Run()
	if err != nil {
		log.Fatal("uploading file", err)
	}
	// no-op check
	json.NewEncoder(os.Stdout).Encode([]interface{}{})
}
