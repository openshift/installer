// Package aws contains AWS-specific structures for installer
// configuration and management.
package aws

// CustomEndpoint store the configuration of a custom url to
// override existing defaults of AWS Services.
// Currently Supports - EC2, IAM, ELB, S3 and Route53.
type CustomEndpoint struct {
	Service string `json:"service"`
	URL     string `json:"url"`
}
