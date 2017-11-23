package model

type InRequest struct {
	Source Source `json:"source"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type OutParams struct {
	Filepath string `json:"filepath"`
}

type Source struct {
	URI              string   `json:"uri"`
	Authentication   AuthPair `json:"authentication"`
	SkipSslVaidation bool     `json:"skip_ssl_validation"`
	Extract          bool     `json:"extract"`
}

type AuthPair struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type InResponse struct{}

type CheckRequest struct{}

type PseudoVersion struct {
	Name      string `json:"name,omitempty"`
	VersionID string `json:"version_id,omitempty"`
}
