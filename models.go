package resource

type Source struct {
	Src         string `json:"source_url"`
	Target      string `json:"target_url"`
	Artifact    string `json:"artifact"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PostExecute string `json:"post_execute,omitempty"`
	Verbose     bool   `json:"verbose,omitempty"`
}

type Version struct {
	Version string `json:"version,omitempty"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
