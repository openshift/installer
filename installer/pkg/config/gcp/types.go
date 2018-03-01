package gcp

type component struct {
	DiskSize string `yaml:"diskSize,omitempty"`
	DiskType string `yaml:"diskType,omitempty"`
	GCEType  string `yaml:"gceType,omitempty"`
}
