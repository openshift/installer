package asset

import (
	"context"
)

// AssetGenerator is used to generate assets.
type AssetGenerator interface {
	// Generate generates this asset given
	// the states of its parent assets.
	GenerateWithContext(context.Context, Parents) error
}

// assetGeneratorAdapter wraps an asset to provide the
// Generate with context function.
type assetGeneratorAdapter struct {
	a Asset
}

// NewDefaultAssetGenerator creates a new adapter to generate
// an asset with a context.
func NewDefaultAssetGenerator(a Asset) AssetGenerator {
	return &assetGeneratorAdapter{a: a}
}

// Generate calls Generate on an asset, dropping the context
// to maintain compatibility with assets that do not implement
// generate with context.
func (a *assetGeneratorAdapter) GenerateWithContext(_ context.Context, p Parents) error {
	return a.a.Generate(p)
}
