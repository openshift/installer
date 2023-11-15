package store

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
)

const (
	stateFileName = ".openshift_install_state.json"
)

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
	asset asset.Asset
	// source is the source from which the asset was fetched
	source assetSource
	// anyParentsDirty is true if any of the parents of the asset are dirty
	anyParentsDirty bool
	// presentOnDisk is true if the asset in on-disk. This is set whether the
	// asset is sourced from on-disk or not. It is used in purging consumed assets.
	presentOnDisk bool
}

// storeImpl is the implementation of Store.
type storeImpl struct {
	directory       string
	assets          map[reflect.Type]*assetState
	stateFileAssets map[string]json.RawMessage
	fileFetcher     asset.FileFetcher
}

// NewStore returns an asset store that implements the asset.Store interface.
func NewStore(dir string) (asset.Store, error) {
	return newStore(dir)
}

func newStore(dir string) (*storeImpl, error) {
	store := &storeImpl{
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
// dependencies if necessary. When purging consumed assets, none of the
// assets in preserved will be purged.
func (s *storeImpl) Fetch(ctx context.Context, a asset.Asset, preserved ...asset.WritableAsset) error {
	if err := s.fetch(ctx, a, ""); err != nil {
		return err
	}
	if err := s.saveStateFile(); err != nil {
		return errors.Wrap(err, "failed to save state")
	}
	if wa, ok := a.(asset.WritableAsset); ok {
		return errors.Wrap(s.purge(append(preserved, wa)), "failed to purge asset")
	}
	return nil
}

// Destroy removes the asset from all its internal state and also from
// disk if possible.
func (s *storeImpl) Destroy(a asset.Asset) error {
	if sa, ok := s.assets[reflect.TypeOf(a)]; ok {
		reflect.ValueOf(a).Elem().Set(reflect.ValueOf(sa.asset).Elem())
	} else if s.isAssetInState(a) {
		if err := s.loadAssetFromState(a); err != nil {
			return err
		}
	} else {
		// nothing to do
		return nil
	}

	if wa, ok := a.(asset.WritableAsset); ok {
		if err := asset.DeleteAssetFromDisk(wa, s.directory); err != nil {
			return err
		}
	}

	delete(s.assets, reflect.TypeOf(a))
	delete(s.stateFileAssets, reflect.TypeOf(a).String())
	return s.saveStateFile()
}

// DestroyState removes the state file from disk
func (s *storeImpl) DestroyState() error {
	s.stateFileAssets = nil
	path := filepath.Join(s.directory, stateFileName)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// loadStateFile retrieves the state from the state file present in the given directory
// and returns the assets map
func (s *storeImpl) loadStateFile() error {
	path := filepath.Join(s.directory, stateFileName)
	assets := map[string]json.RawMessage{}
	data, err := os.ReadFile(path)
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
func (s *storeImpl) loadAssetFromState(a asset.Asset) error {
	bytes, ok := s.stateFileAssets[reflect.TypeOf(a).String()]
	if !ok {
		return errors.Errorf("asset %q is not found in the state file", a.Name())
	}
	return json.Unmarshal(bytes, a)
}

// isAssetInState tests whether the asset is in the state file.
func (s *storeImpl) isAssetInState(a asset.Asset) bool {
	_, ok := s.stateFileAssets[reflect.TypeOf(a).String()]
	return ok
}

// saveStateFile dumps the entire state map into a file
func (s *storeImpl) saveStateFile() error {
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
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o640); err != nil { //nolint:gosec // no sensitive info
		return err
	}
	return nil
}

// fetch populates the given asset, generating it and its dependencies if
// necessary, and returns whether or not the asset had to be regenerated and
// any errors.
func (s *storeImpl) fetch(ctx context.Context, a asset.Asset, indent string) error {
	logrus.Debugf("%sFetching %s...", indent, a.Name())

	assetState, ok := s.assets[reflect.TypeOf(a)]
	if !ok {
		if _, err := s.load(a, ""); err != nil {
			return err
		}
		assetState = s.assets[reflect.TypeOf(a)]
	}

	// Return immediately if the asset has been fetched before,
	// this is because we are doing a depth-first-search, it's guaranteed
	// that we always fetch the parent before children, so we don't need
	// to worry about invalidating anything in the cache.
	if assetState.source != unfetched {
		logrus.Debugf("%sReusing previously-fetched %s", indent, a.Name())
		reflect.ValueOf(a).Elem().Set(reflect.ValueOf(assetState.asset).Elem())
		return nil
	}

	// Re-generate the asset
	dependencies := a.Dependencies()
	parents := make(asset.Parents, len(dependencies))
	for _, d := range dependencies {
		if err := s.fetch(ctx, d, increaseIndent(indent)); err != nil {
			return errors.Wrapf(err, "failed to fetch dependency of %q", a.Name())
		}
		parents.Add(d)
	}
	logrus.Debugf("%sGenerating %s...", indent, a.Name())
	if err := a.Generate(ctx, parents); err != nil {
		return errors.Wrapf(err, "failed to generate asset %q", a.Name())
	}
	assetState.asset = a
	assetState.source = generatedSource
	return nil
}

// load loads the asset and all of its ancestors from on-disk and the state file.
func (s *storeImpl) load(a asset.Asset, indent string) (*assetState, error) {
	logrus.Debugf("%sLoading %s...", indent, a.Name())

	// Stop descent if the asset has already been loaded.
	if state, ok := s.assets[reflect.TypeOf(a)]; ok {
		return state, nil
	}

	// Load dependencies from on-disk.
	anyParentsDirty := false
	for _, d := range a.Dependencies() {
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
		onDiskAsset asset.WritableAsset
		foundOnDisk bool
	)
	if _, isWritable := a.(asset.WritableAsset); isWritable {
		onDiskAsset = reflect.New(reflect.TypeOf(a).Elem()).Interface().(asset.WritableAsset)
		var err error
		foundOnDisk, err = onDiskAsset.Load(s.fileFetcher)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load asset %q", a.Name())
		}
	}

	// Try to load from state file.
	var (
		stateFileAsset         asset.Asset
		foundInStateFile       bool
		onDiskMatchesStateFile bool
	)
	// Do not need to bother with loading from state file if any of the parents
	// are dirty because the asset must be re-generated in this case.
	if !anyParentsDirty {
		foundInStateFile = s.isAssetInState(a)
		if foundInStateFile {
			stateFileAsset = reflect.New(reflect.TypeOf(a).Elem()).Interface().(asset.Asset)
			if err := s.loadAssetFromState(stateFileAsset); err != nil {
				return nil, errors.Wrapf(err, "failed to load asset %q from state file", a.Name())
			}
		}

		if foundOnDisk && foundInStateFile {
			logrus.Debugf("%sLoading %s from both state file and target directory", indent, a.Name())

			// If the on-disk asset is the same as the one in the state file, there
			// is no need to consider the one on disk and to mark the asset dirty.
			onDiskMatchesStateFile = reflect.DeepEqual(onDiskAsset, stateFileAsset)
			if onDiskMatchesStateFile {
				logrus.Debugf("%sOn-disk %s matches asset in state file", indent, a.Name())
			}
		}
	}

	var (
		assetToStore asset.Asset
		source       assetSource
	)
	switch {
	// A parent is dirty. The asset must be re-generated.
	case anyParentsDirty:
		if foundOnDisk {
			logrus.Warningf("%sDiscarding the %s that was provided in the target directory because its dependencies are dirty and it needs to be regenerated", indent, a.Name())
		}
		source = unfetched
	// The asset is on disk and that differs from what is in the source file.
	// The asset is sourced from on disk.
	case foundOnDisk && !onDiskMatchesStateFile:
		logrus.Debugf("%sUsing %s loaded from target directory", indent, a.Name())
		assetToStore = onDiskAsset
		source = onDiskSource
	// The asset is in the state file. The asset is sourced from state file.
	case foundInStateFile:
		logrus.Debugf("%sUsing %s loaded from state file", indent, a.Name())
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
	s.assets[reflect.TypeOf(a)] = state
	return state, nil
}

// purge deletes the on-disk assets that are consumed already.
// E.g., install-config.yaml will be deleted after fetching 'manifests'.
// The target asset is excluded.
func (s *storeImpl) purge(excluded []asset.WritableAsset) error {
	excl := make(map[reflect.Type]bool, len(excluded))
	for _, a := range excluded {
		excl[reflect.TypeOf(a)] = true
	}
	for _, assetState := range s.assets {
		if !assetState.presentOnDisk || excl[reflect.TypeOf(assetState.asset)] {
			continue
		}
		logrus.Infof("Consuming %s from target directory", assetState.asset.Name())
		if err := asset.DeleteAssetFromDisk(assetState.asset.(asset.WritableAsset), s.directory); err != nil {
			return err
		}
		assetState.presentOnDisk = false
	}
	return nil
}

func increaseIndent(indent string) string {
	return indent + "  "
}

// Load retrieves the given asset if it is present in the store and does not generate the asset
// if it does not exist and will return nil.
func (s *storeImpl) Load(a asset.Asset) (asset.Asset, error) {
	foundOnDisk, err := s.load(a, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load asset")
	}

	if foundOnDisk.source == unfetched {
		return nil, nil
	}

	return s.assets[reflect.TypeOf(a)].asset, nil
}

func asFileWriter(a asset.WritableAsset) asset.FileWriter {
	switch v := a.(type) {
	case asset.FileWriter:
		return v
	default:
		return asset.NewDefaultFileWriter(a)
	}
}

func (s *storeImpl) PersistToFile(wa asset.WritableAsset) error {
	return asFileWriter(wa).PersistToFile(s.directory)
}
