package main

import (
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"bufio"
	"path/filepath"

	"github.com/eugenmayer/concourse-static-resource/curlopts"
	"github.com/eugenmayer/concourse-static-resource/log"
	"github.com/eugenmayer/concourse-static-resource/model"
	"errors"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Accessing first argument", errors.New("usage: %s <sources directory>\n"))
		os.Exit(1)
	}
	var sourceDir string = os.Args[1]

	var request model.OutRequest
	var destBaseUrl *url.URL
	var err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatal("reading request", err)
	}

	destBaseUrl, err = url.Parse(request.Source.URI)
	if err != nil {
		log.Fatal("parsing uri", err)
	}

	// read the version if the path is actually provided
	var version string = getVersionFromFile(request.Params.VersionFilepath, sourceDir)
	if version != "" {
		fmt.Println("Loaded version: " + version)
	}

	var curlOpts string = curlopts.Curlopt(request.Source)

	// depending if destFilenamePattern has a placeholder, us version to replace it and set our destFilename
	var destFilenamePattern string = request.Params.DestFilenamePattern
	var destFilename string = ""

	if version == "" && strings.Contains(destFilenamePattern, "<version>") {
		log.Fatal("Inject version", errors.New("You have provided pattern with a <version> placeholder but your provided version is empty - cannot replace"))
	} else if version != "" && !strings.Contains(destFilenamePattern, "<version>") {
		log.Fatal("Inject version", errors.New("You have provided a version but your pattern does miss the <version> placeholder"))
	} else if version != "" && strings.Contains(destFilenamePattern, "<version>") {
		destFilename = strings.Replace(destFilenamePattern, "<version>", version, -1)
	} else {
		destFilename = destFilenamePattern
	}

	var sourceFile string = getSourceFile(request.Params.SourceFilepathGlob)

	// placeholder for the curlPipe dest arg $1 and upload-destination $2
	// the dest URL looks like <URI>/<destFilename>
	var fullDestUrl string = destBaseUrl.String() + "/" + destFilename
	var command string = "curl " + curlOpts + " --upload-file \"$1\" " + fullDestUrl

	curlPipe := exec.Command(
		"sh",
		"-exc",
		command,
		"sh", sourceFile, destBaseUrl.String(), destFilename,
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

func getVersionFromFile(versionFilepath string, sourceDir string) string {
	if versionFilepath != "" {
		var realpath string = filepath.Join(sourceDir, versionFilepath)
		file, err := os.Open(realpath)
		if err != nil {
			log.Fatal("could not find version file at:"+realpath, err)
		}
		defer file.Close()

		var scanner = bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		scanner.Scan()
		var version = scanner.Text()
		if version == "" {
			log.Fatal("reading version from version file", errors.New("Your version file seems to be empty"))
		}
		// probably validate further
		return version
	}
	// else
	return ""
}

func getSourceFile(sourceFileGlob string) string {
	matches, err := filepath.Glob(sourceFileGlob)
	fmt.Println(sourceFileGlob)

	if err != nil {
		log.Fatal("using source glob did not match a file", err)
	}

	if matches == nil {
		log.Fatal("using source glob did not match a file", errors.New(""))
	}
	fmt.Println(matches[0])


	return matches[0]
}
