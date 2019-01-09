package out

import (
	"github.com/lorands/maven-stage-resource"
)

type Request struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	From   string `json:"from"`
	FromRe string `json:"from_re_filter,omitempty"`
	To     string `json:"to,omitempty"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
