// Package assets defines a generic Merkle DAG for assets.
package assets

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-log/log"
	"github.com/pkg/errors"
)

// GetData retrieves injected data by asset name during generation.
// It returns os.ErrNotExist if the asset is not found.
type GetData func(ctx context.Context, name string) (data []byte, err error)

// Put adds an asset to a store.
type Put func(asset Asset) (hash []byte, err error)

// Assets holds a directed, acyclic graph of assets, which can be used
// for building, and rebuilding, before installation.
type Assets struct {
	// assetsByHash stores all the assets by hash.
	assetsByHash map[string]*Asset

	// assetsByName stores the most-recently-put asset for each name.
	assetsByName map[string]*Asset

	// Root references the root asset.
	Root Reference

	// Rebuilders registers asset rebuilders by name.
	Rebuilders map[string]Rebuild
}

// Prune removes any assets from the store which are not accessible
// from Root.
func (assets *Assets) Prune() (err error) {
	currentHashes := map[string]bool{}
	currentNames := map[string]bool{}
	stack := [][]byte{assets.Root.Hash}
	for len(stack) > 0 {
		hash := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if currentHashes[string(hash)] {
			continue
		}
		currentHashes[string(hash)] = true
		asset, ok := assets.assetsByHash[string(hash)]
		if !ok {
			continue // Already over-pruned.
		}
		currentNames[asset.Name] = true
		for _, reference := range asset.Parents {
			stack = append(stack, reference.Hash)
		}
	}

	for hash := range assets.assetsByHash {
		if !currentHashes[string(hash)] {
			delete(assets.assetsByHash, hash)
		}
	}

	for name := range assets.assetsByName {
		if !currentNames[name] {
			delete(assets.assetsByName, name)
		}
	}

	return nil
}

// rebuildAsset is a per-asset helper for Assets.Rebuild.  It rebuilds
// the asset and returns the new asset and its hash on success.
// Ancestors are rebuilt inside getByName as they are retrieved.
func (assets *Assets) rebuildAsset(ctx context.Context, asset *Asset, getByName GetByString, put Put, rebuilt map[string]Reference, logger log.Logger) (newAsset *Asset, hash []byte, err error) {
	logger.Logf("rebuilding %q", asset.Name)

	oldReference := Reference{Name: asset.Name}
	oldReference.Hash, err = asset.Hash()
	if err != nil {
		return nil, nil, err
	}

	if asset.RebuildHelper == nil {
		return nil, nil, errors.Errorf("cannot rebuild %s without a rebuilder", oldReference.String())
	}

	newAsset, err = asset.RebuildHelper(ctx, getByName)
	if err != nil {
		return nil, nil, err
	}
	if newAsset == nil {
		return nil, nil, errors.Errorf("RebuildHelper returned nil for %q", asset.Name)
	}

	newReference := Reference{Name: newAsset.Name}
	newReference.Hash, err = put(*newAsset)
	if err != nil {
		return newAsset, nil, err
	}

	if bytes.Equal(newReference.Hash, oldReference.Hash) {
		logger.Logf("%q is still fresh (%x)", asset.Name, string(oldReference.Hash))
		rebuilt[oldReference.String()] = oldReference
		return newAsset, oldReference.Hash, nil
	}

	logger.Logf("rebuilt %q (%x -> %x)", asset.Name, oldReference.Hash, newReference.Hash)
	rebuilt[oldReference.String()] = newReference
	return newAsset, newReference.Hash, nil
}

// getByName is a helper for Assets.Rebuild.  It rebuilds checks for
// injections and rebuilds any requested assets that aren't already in
// the store.
func (assets *Assets) getByName(ctx context.Context, name string, getByName GetByString, getInjection GetByString, rebuilt map[string]Reference, logger log.Logger) (asset *Asset, err error) {
	assetValue, err := assets.GetByName(ctx, name)
	asset = &assetValue
	if err == nil || !os.IsNotExist(err) {
		return asset, err
	}

	hash, err := asset.Hash()
	if err != nil {
		return asset, err
	}

	oldReference := Reference{
		Name: name,
		Hash: hash,
	}

	if getInjection != nil {
		assetValue, err = getInjection(ctx, name)
		if err == nil {
			asset = &assetValue
		} else if !os.IsNotExist(err) {
			return &assetValue, errors.Wrapf(err, "inject content for %q", name)
		}
	}

	asset, newHash, err := assets.rebuildAsset(ctx, asset, getByName, assets.Put, rebuilt, logger)
	if err != nil {
		return asset, errors.Wrapf(err, "retrieve %q by name", name)
	}

	if bytes.Equal(newHash, oldReference.Hash) {
		logger.Logf("rebuilt %q (%x -> %x)", name, oldReference.Hash, newHash)
	}

	return asset, nil
}

// Rebuild rebuilds the asset store, pulling in any injected data.
func (assets *Assets) Rebuild(ctx context.Context, getInjection GetByString, logger log.Logger) (err error) {
	rebuilt := map[string]Reference{}

	var getByName GetByString
	getByName = func(ctx context.Context, name string) (Asset, error) {
		asset, err := assets.getByName(ctx, name, getByName, getInjection, rebuilt, logger)
		if asset == nil {
			asset = &Asset{Name: name}
		}
		return *asset, err
	}

	var asset Asset
	var newHash []byte
	if assets.Root.Hash != nil {
		asset, err = assets.GetByHash(ctx, assets.Root.Hash)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to retrieve root %x by hash", assets.Root.Hash)
		}

		var newAsset *Asset
		newAsset, newHash, err = assets.rebuildAsset(ctx, &asset, getByName, assets.Put, rebuilt, logger)
		if err != nil {
			return err
		}
		asset = *newAsset
	} else if assets.Root.Name != "" {
		asset, err = getByName(ctx, assets.Root.Name)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, "retrieve root")
		}

		newHash, err = asset.Hash()
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	} else {
		return nil
	}

	assets.Root.Hash = newHash
	assets.Root.Name = asset.Name
	return nil
}

// GetByHash retrieves an asset from the store by hash.
func (assets *Assets) GetByHash(ctx context.Context, hash []byte) (asset Asset, err error) {
	pointer, ok := assets.assetsByHash[string(hash)]
	if ok {
		return *pointer, nil
	}
	return asset, os.ErrNotExist
}

// GetByName retrieves an asset from the store by name.
func (assets *Assets) GetByName(ctx context.Context, name string) (asset Asset, err error) {
	pointer, ok := assets.assetsByName[name]
	if ok {
		return *pointer, nil
	}
	return asset, os.ErrNotExist
}

// Put adds an asset to the store.
func (assets *Assets) Put(asset Asset) (hash []byte, err error) {
	hash, err = (&asset).Hash()
	if err != nil {
		return hash, err
	}

	if assets.assetsByHash == nil {
		assets.assetsByHash = map[string]*Asset{}
	}
	assets.assetsByHash[string(hash)] = &asset

	if assets.assetsByName == nil {
		assets.assetsByName = map[string]*Asset{}
	}
	assets.assetsByName[asset.Name] = &asset
	return hash, nil
}

// Write writes assets to the target directory.  If prune is true, it
// also removes any files from that directory which it did not write
// (for example, leftovers from previous invocations).
func (assets *Assets) Write(ctx context.Context, directory string, prune bool) (err error) {
	err = os.MkdirAll(filepath.Join(directory, subDir), 0777)
	if err != nil {
		return err
	}

	written := map[string]bool{}

	if assets.Root.Hash != nil {
		asset, ok := assets.assetsByHash[string(assets.Root.Hash)]
		if !ok {
			return errors.Errorf("failed to retrieve root %x by hash", assets.Root.Hash)
		}

		err = asset.Write(ctx, directory, assets.GetByHash, written)
		if err != nil {
			return err
		}
	}

	for path := range written {
		remaining := path
		for remaining != "." {
			written[remaining] = true
			remaining = filepath.Dir(remaining)
		}
	}

	if prune {
		filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == directory {
				return nil
			}

			rel, err := filepath.Rel(directory, path)
			if err != nil {
				return err
			}

			if written[rel] {
				return nil
			}
			fmt.Printf("XXX prune %q\n", path)

			if info.IsDir() {
				err = os.RemoveAll(path)
				if err != nil {
					return err
				}
				return filepath.SkipDir
			}

			return os.Remove(path)
		})
	}

	return nil
}

// Read rebuilds the asset store, pulling in data from a previous
// Assets.Write call and falling back to getDefault for requested
// assets that aren't in that directory.
func (assets *Assets) Read(ctx context.Context, directory string, getDefault GetData, logger log.Logger) (err error) {
	return assets.Rebuild(ctx, func(ctx context.Context, name string) (Asset, error) {
		logger.Logf("checking injection for %q", name)
		asset := &Asset{Name: name}
		loaded := false
		err := asset.Read(ctx, directory, logger)
		if err == nil {
			loaded = true
			logger.Logf("loaded %q from %q", name, asset.path())
		} else if !os.IsNotExist(err) {
			return *asset, err
		}

		if !loaded && getDefault != nil {
			var data []byte
			data, err = getDefault(ctx, name)
			if err == nil {
				asset.RebuildHelper = ConstantDataRebuilder(ctx, name, data, false)
				if len(data) > 10 {
					logger.Logf("default %q to \"%s...\"", name, string(data[:10]))
				} else {
					logger.Logf("default %q to %q", name, string(data))
				}
			} else if !os.IsNotExist(err) {
				return *asset, errors.Wrapf(err, "defaulting %q", name)
			}
		}

		if asset.RebuildHelper == nil {
			var ok bool
			asset.RebuildHelper, ok = assets.Rebuilders[name]
			if !ok {
				if loaded {
					asset.RebuildHelper = ConstantDataRebuilder(ctx, name, asset.Data, false)
				} else {
					return *asset, errors.Wrapf(err, "cannot inject %q without a file, default, or rebuilder", name)
				}
			}
		}

		return *asset, nil
	}, logger)
}
