package server

import (
	"fmt"
	"strings"
)

// KnownPlatforms is the list of supported platforms.
var KnownPlatforms = []string{
	"aws-tf",
	"bare-metal-tf",
	"aws",
	"bare-metal",
	"azure",
	"openstack",
}

// PlatformsValue is a flag.Value/flag.Getter compatible type for reading platform arguments
type PlatformsValue struct {
	Names []string
}

// String formats the platform list in a command-line-acceptable way
func (p *PlatformsValue) String() string {
	return strings.Join(p.Names, ",")
}

// Set parses a command line value into Names, or returns an error.
func (p *PlatformsValue) Set(s string) error {
	p.Names = strings.Split(s, ",")
	for _, s := range p.Names {
		found := false
		for _, known := range KnownPlatforms {
			if s == known {
				found = true
				break
			}
		}

		if !found {
			platforms := strings.Join(KnownPlatforms, ",")
			return fmt.Errorf("unrecognized platform \"%s\": known platforms are %s", s, platforms)
		}
	}

	return nil
}
