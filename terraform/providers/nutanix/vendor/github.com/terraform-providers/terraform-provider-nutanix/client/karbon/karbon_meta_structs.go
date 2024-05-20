package karbon

type MetaVersionResponse struct {
	BuildDate *string `json:"build_date" mapstructure:"build_date, omitempty"`
	GitCommit *string `json:"git_commit" mapstructure:"git_commit, omitempty"`
	Version   *string `json:"version" mapstructure:"version, omitempty"`
}

type MetaSemanticVersionResponse struct {
	MajorVersion    int64 `json:"major_version" mapstructure:"major_version, omitempty"`
	MinorVersion    int64 `json:"minor_version" mapstructure:"minor_version, omitempty"`
	RevisionVersion int64 `json:"revision_version" mapstructure:"revision_version, omitempty"`
}
