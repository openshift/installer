package asset

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given asset, generating it and its
	// dependencies if necessary.
	Fetch(Asset) (*State, error)
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	assets map[Asset]*State
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) (*State, error) {
	return s.fetch(asset, "")
}

func (s *StoreImpl) fetch(asset Asset, indent string) (*State, error) {
	logrus.Debugf("%sFetching %s...", indent, asset.Name())
	state, ok := s.assets[asset]
	if ok {
		logrus.Debugf("%sFound %s...", indent, asset.Name())
		return state, nil
	}

	dependencies := asset.Dependencies()
	dependenciesStates := make(map[Asset]*State, len(dependencies))
	if len(dependencies) > 0 {
		logrus.Debugf("%sGenerating dependencies of %s...", indent, asset.Name())
	}
	for _, d := range dependencies {
		ds, err := s.fetch(d, indent+"  ")
		if err != nil {
			return nil, err
		}
		dependenciesStates[d] = ds
	}

	logrus.Debugf("%sGenerating %s...", indent, asset.Name())
	state, err := asset.Generate(dependenciesStates)
	if err != nil {
		return nil, fmt.Errorf("failed to generate asset %q: %v", asset.Name(), err)
	}
	if s.assets == nil {
		s.assets = make(map[Asset]*State)
	}
	s.assets[asset] = state
	return state, nil
}
