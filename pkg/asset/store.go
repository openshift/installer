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

// assetSource indicates from where the asset was fetched
type assetSource int

const (
	// unsourced indicates that the asset has not been fetched
	unfetched assetSource = iota
	// generatedSource indicates that the asset was generated
	generatedSource
	// onDiskSource indicates that the asset was fetched from disk
	onDiskSource
	// stateFileSource indicates that the asset was fetched from the state file
	stateFileSource
)

type assetState struct {
	// asset is the asset.
	// If the asset has not been fetched, then this will be nil.
	asset Asset
	// source is the source from which the asset was fetched
	source assetSource
	// anyParentsDirty is true if any of the parents of the asset are dirty
	anyParentsDirty bool
	// presentOnDisk is true if the asset in on-disk. This is set whether the
	// asset is sourced from on-disk or not. It is used in purging consumed assets.
	presentOnDisk bool
}

// StoreImpl is the implementation of Store.
type StoreImpl struct {
	directory       string
	assets          map[reflect.Type]*assetState
	stateFileAssets map[string]json.RawMessage
	fileFetcher     FileFetcher
}

// NewStore returns an asset store that implements the Store interface.
func NewStore(dir string) (Store, error) {
	store := &StoreImpl{
		directory:   dir,
		fileFetcher: &fileFetcher{directory: dir},
		assets:      map[reflect.Type]*assetState{},
	}

	if err := store.loadStateFile(); err != nil {
		return nil, err
	}
	return store, nil
}

// Fetch retrieves the state of the given asset, generating it and its
// dependencies if necessary.
func (s *StoreImpl) Fetch(asset Asset) error {
	if err := s.fetch(asset, ""); err != nil {
		return err
	}
	if err := s.saveStateFile(); err != nil {
		return errors.Wrapf(err, "failed to save state")
	}
	if wa, ok := asset.(WritableAsset); ok {
		return errors.Wrapf(s.purge(wa), "failed to purge asset")
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
	return s.saveStateFile()
}

// loadStateFile retrieves the state from the state file present in the given directory
// and returns the assets map
func (s *StoreImpl) loadStateFile() error {
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

// saveStateFile dumps the entire state map into a file
func (s *StoreImpl) saveStateFile() error {
	if s.stateFileAssets == nil {
		s.stateFileAssets = map[string]json.RawMessage{}
	}
	for k, v := range s.assets {
		if v.source == unfetched {
			continue
		}
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
func (s *StoreImpl) fetch(asset Asset, indent string) error {
	logrus.Debugf("%sFetching %q...", indent, asset.Name())

	assetState, ok := s.assets[reflect.TypeOf(asset)]
	if !ok {
		if _, err := s.load(asset, ""); err != nil {
			return err
		}
		assetState = s.assets[reflect.TypeOf(asset)]
	}

	// Return immediately if the asset has been fetched before,
	// this is because we are doing a depth-first-search, it's guaranteed
	// that we always fetch the parent before children, so we don't need
	// to worry about invalidating anything in the cache.
	if assetState.source != unfetched {
		logrus.Debugf("%sReusing previously-fetched %q", indent, asset.Name())
		reflect.ValueOf(asset).Elem().Set(reflect.ValueOf(assetState.asset).Elem())
		return nil
	}

	// Re-generate the asset
	dependencies := asset.Dependencies()
	parents := make(Parents, len(dependencies))
	for _, d := range dependencies {
		if err := s.fetch(d, increaseIndent(indent)); err != nil {
			return errors.Wrapf(err, "failed to fetch dependency of %q", asset.Name())
		}
		parents.Add(d)
	}
	logrus.Debugf("%sGenerating %q...", indent, asset.Name())
	if err := asset.Generate(parents); err != nil {
		return errors.Wrapf(err, "failed to generate asset %q", asset.Name())
	}
	assetState.asset = asset
	assetState.source = generatedSource
	return nil
}

// load loads the asset and all of its ancestors from on-disk and the state file.
func (s *StoreImpl) load(asset Asset, indent string) (*assetState, error) {
	logrus.Debugf("%sLoading %q...", indent, asset.Name())

	// Stop descent if the asset has already been loaded.
	if state, ok := s.assets[reflect.TypeOf(asset)]; ok {
		return state, nil
	}

	// Load dependencies from on-disk.
	anyParentsDirty := false
	for _, d := range asset.Dependencies() {
		state, err := s.load(d, increaseIndent(indent))
		if err != nil {
			return nil, err
		}
		if state.anyParentsDirty || state.source == onDiskSource {
			anyParentsDirty = true
		}
	}

	// Try to load from on-disk.
	var (
		onDiskAsset WritableAsset
		foundOnDisk bool
	)
	if _, isWritable := asset.(WritableAsset); isWritable {
		onDiskAsset = reflect.New(reflect.TypeOf(asset).Elem()).Interface().(WritableAsset)
		var err error
		foundOnDisk, err = onDiskAsset.Load(s.fileFetcher)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load asset %q", asset.Name())
		}
	}

	// Try to load from state file.
	var (
		stateFileAsset         Asset
		foundInStateFile       bool
		onDiskMatchesStateFile bool
	)
	// Do not need to bother with loading from state file if any of the parents
	// are dirty because the asset must be re-generated in this case.
	if !anyParentsDirty {
		foundInStateFile = s.isAssetInState(asset)
		if foundInStateFile {
			stateFileAsset = reflect.New(reflect.TypeOf(asset).Elem()).Interface().(Asset)
			if err := s.loadAssetFromState(stateFileAsset); err != nil {
				return nil, errors.Wrapf(err, "failed to load asset %q from state file", asset.Name())
			}
		}

		if foundOnDisk && foundInStateFile {
			logrus.Debugf("%sLoading %q from both state file and target directory", indent, asset.Name())

			// If the on-disk asset is the same as the one in the state file, there
			// is no need to consider the one on disk and to mark the asset dirty.
			onDiskMatchesStateFile = reflect.DeepEqual(onDiskAsset, stateFileAsset)
			if onDiskMatchesStateFile {
				logrus.Debugf("%sOn-disk %q matches asset in state file", indent, asset.Name())
			}
		}
	}

	var (
		assetToStore Asset
		source       assetSource
	)
	switch {
	// A parent is dirty. The asset must be re-generated.
	case anyParentsDirty:
		if foundOnDisk {
			logrus.Warningf("%sDiscarding the %q that was provided in the target directory because its dependencies are dirty and it needs to be regenerated", indent, asset.Name())
		}
		source = unfetched
	// The asset is on disk and that differs from what is in the source file.
	// The asset is sourced from on disk.
	case foundOnDisk && !onDiskMatchesStateFile:
		logrus.Debugf("%sUsing %q loaded from target directory", indent, asset.Name())
		assetToStore = onDiskAsset
		source = onDiskSource
	// The asset is in the state file. The asset is sourced from state file.
	case foundInStateFile:
		logrus.Debugf("%sUsing %q loaded from state file", indent, asset.Name())
		assetToStore = stateFileAsset
		source = stateFileSource
	// There is no existing source for the asset. The asset will be generated.
	default:
		source = unfetched
	}

	state := &assetState{
		asset:           assetToStore,
		source:          source,
		anyParentsDirty: anyParentsDirty,
		presentOnDisk:   foundOnDisk,
	}
	s.assets[reflect.TypeOf(asset)] = state
	return state, nil
}

// purge deletes the on-disk assets that are consumed already.
// E.g., install-config.yaml will be deleted after fetching 'manifests'.
// The target asset is excluded.
func (s *StoreImpl) purge(excluded WritableAsset) error {
	for _, assetState := range s.assets {
		if !assetState.presentOnDisk {
			continue
		}
		if reflect.TypeOf(assetState.asset) == reflect.TypeOf(excluded) {
			continue
		}
		logrus.Infof("Consuming %q from target directory", assetState.asset.Name())
		if err := deleteAssetFromDisk(assetState.asset.(WritableAsset), s.directory); err != nil {
			return err
		}
		assetState.presentOnDisk = false
	}
	return nil
}

func increaseIndent(indent string) string {
	return indent + "  "
}
