package in

import (
	"github.com/lorands/maven-stage-resource"
)

type Request struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
	Params  Params           `json:"params,omitempty"`
}

type Params struct {
	Version      string `json:"version,omitempty"`
	DownloadOnly bool   `json:"download_only,omitempty"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
