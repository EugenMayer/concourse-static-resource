package log

import "os"

func Fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
