package curlopts

import (
	"github.com/eugenmayer/concourse-static-resource/model"
	"fmt"
)

func Curlopt(source model.Source) string {
	curlOpts := "--location --retry 3 --fail"

	if (source.Authentication.User != "") {
		curlOpts = fmt.Sprintf("%s -u '%s:%s'", curlOpts, source.Authentication.User, source.Authentication.Password)
	}

	if (source.SkipSslVaidation) {
		curlOpts = curlOpts + " -k"
	}
	return curlOpts
}
