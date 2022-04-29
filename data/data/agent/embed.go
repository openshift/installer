package agent

import "embed"

// IgnitionData contains the source data for building the ignition file
//go:embed *
var IgnitionData embed.FS
