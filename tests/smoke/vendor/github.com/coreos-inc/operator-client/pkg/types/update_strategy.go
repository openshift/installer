package types

// UpdateStrategy describes the type of update to perform.
type UpdateStrategy string

const (
	// StrategyImage will update only the image field for a component.
	StrategyImage UpdateStrategy = "image"
	// StrategyReplace will replace the entire manifest on the server
	// with one we store locally.
	StrategyReplace UpdateStrategy = "replace"
)
