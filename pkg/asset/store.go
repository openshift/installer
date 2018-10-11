package asset

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	stateFileName = ".openshift_install_state.json"
)

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given asset, generating it and its
	// dependencies if necessary.
	Fetch(Asset) error
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	assets          map[reflect.Type]Asset
	stateFileAssets map[string]json.RawMessage
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) error {
	return s.fetch(asset, "")
}

// Load retrieves the state from the state file present in the given directory
// and returns the assets map
func (s *StoreImpl) Load(dir string) error {
	path := filepath.Join(dir, stateFileName)
	assets := make(map[string]json.RawMessage)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, &assets)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal state file %s", path)
	}
	s.stateFileAssets = assets
	return nil
}

// GetStateAsset renders the asset object arguments from the state file contents
// also returns a boolean indicating whether the object was found in the state file or not
func (s *StoreImpl) GetStateAsset(asset Asset) (bool, error) {
	bytes, ok := s.stateFileAssets[reflect.TypeOf(asset).String()]
	if !ok {
		return false, nil
	}
	err := json.Unmarshal(bytes, asset)
	return true, err
}

// Save dumps the entire state map into a file
func (s *StoreImpl) Save(dir string) error {
	assetMap := make(map[string]Asset)
	for k, v := range s.assets {
		assetMap[k.String()] = v
	}
	data, err := json.MarshalIndent(&assetMap, "", "    ")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, stateFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
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

	// Before generating the asset, look if we have it all ready in the state file
	// if yes, then use it instead
	logrus.Debugf("%sLooking up asset from state file: %s", indent, reflect.TypeOf(asset).String())
	ok, err := s.GetStateAsset(asset)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal asset '%s' from state file '%s'", asset.Name(), stateFileName)
	}
	if ok {
		logrus.Debugf("%sAsset found in state file", indent)
	} else {
		logrus.Debugf("%sAsset not found in state file. Generating %s...", indent, asset.Name())
		err := asset.Generate(parents)
		if err != nil {
			return errors.Wrapf(err, "failed to generate asset %s", asset.Name())
		}
	}
	if s.assets == nil {
		s.assets = make(map[reflect.Type]Asset)
	}
	s.assets[reflect.TypeOf(asset)] = asset
	return nil
}
