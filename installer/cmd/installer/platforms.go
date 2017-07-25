package main

import (
	"fmt"
	"strings"
)

// knownPlatforms is the list of supported platforms.
var knownPlatforms = []string{
	"aws",
	"aws-tf",
	"bare-metal",
	"bare-metal-tf",
	"azure",
	"openstack",
}

// platformsValue is a flag.Value/flag.Getter compatible type for reading platform arguments
type platformsValue struct {
	names []string
}

// String formats the platform list in a command-line-acceptable way
func (p *platformsValue) String() string {
	return strings.Join(p.names, ",")
}

// Set parses a command line value into Names, or returns an error.
func (p *platformsValue) Set(s string) error {
	p.names = strings.Split(s, ",")
	for _, s := range p.names {
		found := false
		for _, known := range knownPlatforms {
			if s == known {
				found = true
				break
			}
		}

		if !found {
			platforms := strings.Join(knownPlatforms, ",")
			return fmt.Errorf("unrecognized platform \"%s\": known platforms are %s", s, platforms)
		}
	}

	return nil
}
