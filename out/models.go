package out

import (
	"github.com/lorands/maven-stage-resource"
)

type Request struct {
	Source resource.Source `json:"source"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
