package asset

import (
	"context"
)

// Generator is used to generate assets.
type Generator interface {
	// Generate generates this asset given
	// the states of its parent assets.
	GenerateWithContext(context.Context, Parents) error
}

// generatorAdapter wraps an asset to provide the
// Generate with context function.
type generatorAdapter struct {
	a Asset
}

// NewDefaultGenerator creates a new adapter to generate
// an asset with a context.
func NewDefaultGenerator(a Asset) Generator {
	return &generatorAdapter{a: a}
}

// Generate calls Generate on an asset, dropping the context
// to maintain compatibility with assets that do not implement
// generate with context.
func (a *generatorAdapter) GenerateWithContext(_ context.Context, p Parents) error {
	return a.a.Generate(p)
}
