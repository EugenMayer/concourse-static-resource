package shared

import (
	"fmt"

	"github.com/eugenmayer/concourse-static-resource/model"
)

func Curlopt(source model.Source) string {
	curlOpts := "--location --retry 3 --fail --show-error --silent"

	if source.Authentication.User != "" {
		curlOpts = fmt.Sprintf("%s -u '%s:%s'", curlOpts, source.Authentication.User, source.Authentication.Password)
	}

	if source.SkipSslVaidation {
		curlOpts = curlOpts + " -k"
	}
	return curlOpts
}
