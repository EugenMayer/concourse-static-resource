package model

type CheckRequest struct {
	Source Source `json:"source"`
}
type CheckResponse struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
	Params InParams `json:"params"`

}

type InResponse struct {
	Version  Version        `json:"version"`
	MetaData []MetaDataPair `json:"metadata"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
	Version Version `json:"version"`
}

type OutResponse struct {
	Version  Version        `json:"version"`
	MetaData []MetaDataPair `json:"metadata"`
}


type OutParams struct {
	SourceFilepathGlob string `json:"source_filepath"`
	VersionFilepath  string   `json:"version_filepath"`
}


type InParams struct {
	VersionFilepath  string   `json:"version_filepath"`
}

type Version struct {
	Ref string `json:"ref"`
}

type MetaDataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}


type Source struct {
	URI              string   `json:"uri"`
	VersionStatic    string   `json:"version_static"`
	Authentication   AuthPair `json:"authentication"`
	SkipSslVaidation bool     `json:"skip_ssl_validation"`
	Extract          bool     `json:"extract"`
}

type AuthPair struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type PseudoVersion struct {
	Name      string `json:"name,omitempty"`
	VersionID string `json:"version_id,omitempty"`
}
