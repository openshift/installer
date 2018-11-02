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

	// Destroy removes the asset from all its internal state and also from
	// disk if possible.
	Destroy(Asset) error
}

// assetState includes an asset and a boolean that indicates
// whether it's dirty or not.
type assetState struct {
	asset Asset
	dirty bool
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	directory       string
	assets          map[reflect.Type]assetState
	stateFileAssets map[string]json.RawMessage
	fileFetcher     *fileFetcher

	markedForPurge []WritableAsset // This records the on-disk assets that are loaded already, which will be cleaned up in the end.
}

// NewStore returns an asset store that implements the Store interface.
func NewStore(dir string) (Store, error) {
	store := &StoreImpl{
		directory:   dir,
		fileFetcher: &fileFetcher{directory: dir},
		assets:      make(map[reflect.Type]assetState),
	}

	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) error {
	if _, err := s.fetch(asset, ""); err != nil {
		return err
	}
	if err := s.save(); err != nil {
		return errors.Wrapf(err, "failed to save state")
	}
	if wa, ok := asset.(WritableAsset); ok {
		return errors.Wrapf(s.purge([]WritableAsset{wa}), "failed to purge asset")
	}
	return nil
}

// Destroy removes the asset from all its internal state and also from
// disk if possible.
func (s *StoreImpl) Destroy(asset Asset) error {
	if sa, ok := s.assets[reflect.TypeOf(asset)]; ok {
		reflect.ValueOf(asset).Elem().Set(reflect.ValueOf(sa.asset).Elem())
	} else if s.isAssetInState(asset) {
		if err := s.loadAssetFromState(asset); err != nil {
			return err
		}
	} else {
		// nothing to do
		return nil
	}

	if wa, ok := asset.(WritableAsset); ok {
		if err := deleteAssetFromDisk(wa, s.directory); err != nil {
			return err
		}
	}

	delete(s.assets, reflect.TypeOf(asset))
	delete(s.stateFileAssets, reflect.TypeOf(asset).String())
	return s.save()
}

// load retrieves the state from the state file present in the given directory
// and returns the assets map
func (s *StoreImpl) load() error {
	path := filepath.Join(s.directory, stateFileName)
	assets := map[string]json.RawMessage{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, &assets)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal state file %q", path)
	}
	s.stateFileAssets = assets
	return nil
}

// loadAssetFromState renders the asset object arguments from the state file contents.
func (s *StoreImpl) loadAssetFromState(asset Asset) error {
	bytes, ok := s.stateFileAssets[reflect.TypeOf(asset).String()]
	if !ok {
		return errors.Errorf("asset %q is not found in the state file", asset.Name())
	}
	return json.Unmarshal(bytes, asset)
}

// isAssetInState tests whether the asset is in the state file.
func (s *StoreImpl) isAssetInState(asset Asset) bool {
	_, ok := s.stateFileAssets[reflect.TypeOf(asset).String()]
	return ok
}

// save dumps the entire state map into a file
func (s *StoreImpl) save() error {
	if s.stateFileAssets == nil {
		s.stateFileAssets = map[string]json.RawMessage{}
	}
	for k, v := range s.assets {
		data, err := json.MarshalIndent(v.asset, "", "    ")
		if err != nil {
			return err
		}
		s.stateFileAssets[k.String()] = json.RawMessage(data)
	}
	data, err := json.MarshalIndent(s.stateFileAssets, "", "    ")
	if err != nil {
		return err
	}

	path := filepath.Join(s.directory, stateFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

// fetch populates the given asset, generating it and its dependencies if
// necessary, and returns whether or not the asset had to be regenerated and
// any errors.
func (s *StoreImpl) fetch(asset Asset, indent string) (bool, error) {
	logrus.Debugf("%sFetching %q...", indent, asset.Name())

	// Return immediately if the asset is found in the cache,
	// this is because we are doing a depth-first-search, it's guaranteed
	// that we always fetch the parent before children, so we don't need
	// to worry about invalidating anything in the cache.
	storedAsset, ok := s.assets[reflect.TypeOf(asset)]
	if ok {
		logrus.Debugf("%sReusing previously-fetched %q", indent, asset.Name())
		reflect.ValueOf(asset).Elem().Set(reflect.ValueOf(storedAsset.asset).Elem())
		return storedAsset.dirty, nil
	}

	dependencies := asset.Dependencies()
	parents := make(Parents, len(dependencies))
	if len(dependencies) > 0 {
		logrus.Debugf("%sFetching dependencies of %q...", indent, asset.Name())
	}

	var anyParentsDirty bool
	for _, d := range dependencies {
		dirty, err := s.fetch(d, indent+"  ")
		if err != nil {
			return false, errors.Wrapf(err, "failed to fetch dependency of %q", asset.Name())
		}
		if dirty {
			anyParentsDirty = true
		}
		parents.Add(d)
	}

	// Try to find the asset from the state file.
	foundInStateFile := s.isAssetInState(asset)
	if foundInStateFile {
		logrus.Debugf("%sFound %q in state file", indent, asset.Name())
	}

	// Try to load from the provided files.
	var foundOnDisk bool
	if as, ok := asset.(WritableAsset); ok {
		var err error
		foundOnDisk, err = as.Load(s.fileFetcher)
		if err != nil {
			return false, errors.Wrapf(err, "failed to load asset %q", asset.Name())
		}
		if foundOnDisk {
			logrus.Infof("Consuming %q from target directory", asset.Name())
			s.markedForPurge = append(s.markedForPurge, as)
		}
	}

	if anyParentsDirty && foundOnDisk {
		logrus.Warningf("%sDiscarding the %q that was provided in the target directory because its dependencies are dirty and it needs to be regenerated", indent, asset.Name())
	}

	if anyParentsDirty || (!foundOnDisk && !foundInStateFile) {
		logrus.Debugf("%sGenerating %q...", indent, asset.Name())
		if err := asset.Generate(parents); err != nil {
			return false, errors.Wrapf(err, "failed to generate asset %q", asset.Name())
		}
	} else if foundInStateFile && foundOnDisk {
		logrus.Debugf("%sLoading %q from both state file and target directory", indent, asset.Name())

		stateAsset := reflect.New(reflect.TypeOf(asset).Elem()).Interface().(Asset)
		if err := s.loadAssetFromState(stateAsset); err != nil {
			return false, errors.Wrapf(err, "failed to load asset %q from state file", asset.Name())
		}

		// If the on-disk asset is the same as the one in the state file, there
		// is no need to consider the one on disk and to mark the asset dirty.
		if reflect.DeepEqual(stateAsset, asset) {
			foundOnDisk = false
		}
	} else if foundInStateFile {
		logrus.Debugf("%sLoading %q from state file", indent, asset.Name())
		if err := s.loadAssetFromState(asset); err != nil {
			return false, errors.Wrapf(err, "failed to load asset %q from state file", asset.Name())
		}
	} else if foundOnDisk {
		logrus.Debugf("%sLoading %q from target directory", indent, asset.Name())
	}

	dirty := anyParentsDirty || foundOnDisk
	s.assets[reflect.TypeOf(asset)] = assetState{asset: asset, dirty: dirty}
	return dirty, nil
}

// purge deletes the on-disk assets that are consumed already.
// E.g., install-config.yml will be deleted after fetching 'manifests'.
// The target assets are excluded.
func (s *StoreImpl) purge(excluded []WritableAsset) error {
	var toPurge []WritableAsset
	for _, asset := range s.markedForPurge {
		var found bool
		for _, as := range excluded {
			if reflect.TypeOf(as) == reflect.TypeOf(asset) {
				found = true
				break
			}
		}
		if !found {
			toPurge = append(toPurge, asset)
		}
	}

	for _, asset := range toPurge {
		if err := deleteAssetFromDisk(asset, s.directory); err != nil {
			return err
		}
	}
	s.markedForPurge = []WritableAsset{}
	return nil
}
