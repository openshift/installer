//go:build generate
// +build generate

package main

import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

//go:generate go run scripts/get_test_image/get_test_image.go
