package shared

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func InjectVersionIntoPath(path string, version string, pattern string) string {
	if pattern == "" {
		pattern = "<version>"
	}

	if version == "" && strings.Contains(path, pattern) {
		log.Fatal("Inject version", errors.New("You have provided pattern with a <version> placeholder but your provided version is empty - cannot replace"))
	} else if version != "" && !strings.Contains(path, pattern) {
		// log.Fatal("Inject version", errors.New("You have provided a version but your pattern does miss the <version> placeholder"))
	} else if version != "" && strings.Contains(path, pattern) {
		path = strings.Replace(path, pattern, version, -1)
	}

	return path
}

func GetVersionFromFile(versionFilepath string, sourceDir string) string {
	if versionFilepath != "" {
		var realpath = filepath.Join(sourceDir, versionFilepath)
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
			log.Fatal("reading version from version file", errors.New("your version file seems to be empty"))
		}
		// probably validate further
		return version
	}
	// else
	return ""
}

func GetSourceFile(sourceFileGlob string, sourceDir string) string {
	var realpath = filepath.Join(sourceDir, sourceFileGlob)
	matches, err := filepath.Glob(realpath)

	if err != nil {
		log.Fatal("using source glob did not match a file", err)
	}

	if matches == nil {
		log.Fatal("using source glob did not match a file", errors.New(""))
	}

	return matches[0]
}
