package gcp

type component struct {
	DiskSize string `yaml:"DiskSize,omitempty"`
	DiskType string `yaml:"DiskType,omitempty"`
	GCEType  string `yaml:"GCEType,omitempty"`
}
