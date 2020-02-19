package provider

import (
	"time"
)

// DefaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
const DefaultAPITimeout = time.Minute * 5
