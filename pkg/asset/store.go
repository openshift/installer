package asset

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given asset, generating it and its
	// dependencies if necessary.
	Fetch(Asset) error
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	assets map[reflect.Type]Asset
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) error {
	return s.fetch(asset, "")
}

func (s *StoreImpl) fetch(asset Asset, indent string) error {
	logrus.Debugf("%sFetching %s...", indent, asset.Name())
	storedAsset, ok := s.assets[reflect.TypeOf(asset)]
	if ok {
		logrus.Debugf("%sFound %s...", indent, asset.Name())
		reflect.ValueOf(asset).Elem().Set(reflect.ValueOf(storedAsset).Elem())
		return nil
	}

	dependencies := asset.Dependencies()
	parents := make(Parents, len(dependencies))
	if len(dependencies) > 0 {
		logrus.Debugf("%sGenerating dependencies of %s...", indent, asset.Name())
	}
	for _, d := range dependencies {
		err := s.fetch(d, indent+"  ")
		if err != nil {
			return errors.Wrapf(err, "failed to fetch dependency for %s", asset.Name())
		}
		parents.Add(d)
	}

	logrus.Debugf("%sGenerating %s...", indent, asset.Name())
	err := asset.Generate(parents)
	if err != nil {
		return errors.Wrapf(err, "failed to generate asset %s", asset.Name())
	}
	if s.assets == nil {
		s.assets = make(map[reflect.Type]Asset)
	}
	s.assets[reflect.TypeOf(asset)] = asset
	return nil
}
