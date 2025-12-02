package reflect

// Options provides configuration settings for how the reflection behavior
// works, letting callers tweak different behaviors based on their needs.
type Options struct {
	// UnhandledNullAsEmpty controls whether null values should be
	// translated into empty values without provider interaction, or if
	// they must be explicitly handled.
	UnhandledNullAsEmpty bool

	// UnhandledUnknownAsEmpty controls whether null values should be
	// translated into empty values without provider interaction, or if
	// they must be explicitly handled.
	UnhandledUnknownAsEmpty bool

	// AllowRoundingNumbers silently rounds numbers that don't fit
	// perfectly in the types they're being stored in, rather than
	// returning errors. Numbers will always be rounded towards 0.
	AllowRoundingNumbers bool
}
