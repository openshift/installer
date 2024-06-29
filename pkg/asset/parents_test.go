package asset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type parentsAsset struct {
	x int
}

func (a *parentsAsset) Name() string {
	return "parents-asset"
}

func (a *parentsAsset) Dependencies() []Asset {
	return []Asset{}
}

func (a *parentsAsset) Generate(context.Context, Parents) error {
	return nil
}

func TestParentsGetPointer(t *testing.T) {
	origAsset := &parentsAsset{x: 1}
	parents := Parents{}
	parents.Add(origAsset)

	retrievedAsset := &parentsAsset{}
	parents.Get(retrievedAsset)
	assert.Equal(t, 1, retrievedAsset.x)
}
