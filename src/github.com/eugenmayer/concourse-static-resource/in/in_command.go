package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/eugenmayer/concourse-static-resource/curlopts"
	"github.com/eugenmayer/concourse-static-resource/log"
	"github.com/eugenmayer/concourse-static-resource/model"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		log.Fatal("creating destination", err)
	}

	var request model.InRequest

	err = json.NewDecoder(os.Stdin).Decode(&request)

	if err != nil {
		log.Fatal("reading request", err)
	}

	sourceURL, err := url.Parse(request.Source.URI)

	if err != nil {
		log.Fatal("parsing uri", err)
	}

	var curlOpts string = curlopts.Curlopt(request.Source)

	// placeholder for the curlPipe source arg
	curlOpts = curlOpts + " \"$1\""

	var command string
	if (request.Source.Extract == true) {
		// $2 is the placeholder for the curlPipe destination arg
		command = fmt.Sprintf("curl %s %s", curlOpts, "| tar --warning=no-unknown-keyword -C \"$2\" -zxf -")
	} else {
		// $2 is the placeholder for the curlPipe destination arg
		command = fmt.Sprintf("cd \"$2\"; curl %s %s", curlOpts, "-O")
	}

	curlPipe := exec.Command(
		"sh",
		"-ec",
		command,
		"sh", sourceURL.String(), destination,
	)

	curlPipe.Stdout = os.Stderr
	curlPipe.Stderr = os.Stderr

	err = curlPipe.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(model.InResponse{})
}
