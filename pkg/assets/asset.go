package assets

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-log/log"
	"github.com/pkg/errors"
)

var filenameRegexp = regexp.MustCompile("[^A-Za-z0-9.-]+")
var subDir = ".state"

// GetByString retrieves an asset from a store.  Returning an Asset
// (instead of an *Asset) creates a copy to ensure the caller cannot
// adjust the immutable *Asset held in the store.  It returns
// os.ErrNotExist if the asset is not found.
type GetByString func(ctx context.Context, id string) (asset Asset, err error)

// GetByBytes retrieves an asset from a store.  Returning an Asset
// (instead of an *Asset) creates a copy to ensure the caller cannot
// adjust the immutable *Asset held in the store.  It returns
// os.ErrNotExist if the asset is not found.
type GetByBytes func(ctx context.Context, id []byte) (asset Asset, err error)

// Rebuild rebuilds an asset after a parent changes.
type Rebuild func(ctx context.Context, getByName GetByString) (*Asset, error)

// Asset is a node in the asset graph.
type Asset struct {
	// Parents holds hashes for our parent assets.
	Parents []Reference `json:"parents,omitempty"`

	// Name summarizes the semantic meaning of this asset in a form
	// which could be used as a filename.
	Name string `json:"name,omitempty"`

	// Data holds the asset payload (e.g. the password for a password
	// asset).
	Data []byte `json:"data,omitempty"`

	// Frozen marks assets that shoud no longer be regenerated
	// (e.g. because a user has overridden their values).
	Frozen bool `json:"frozen,omitempty"`

	// RebuildHelper rebuilds this asset based on updated dependencies.
	RebuildHelper Rebuild `json:"-"`
}

// Hash returns a cryptographic hash for this asset.
func (asset *Asset) Hash() (hash []byte, err error) {
	sort.Sort(referenceSlice(asset.Parents))
	data, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	hsh := sha1.Sum(data)
	return hsh[:], nil
}

// GetParents retrieves parents by name and adds them to Asset.Parent.
func (asset *Asset) GetParents(ctx context.Context, getByName GetByString, names ...string) (parents map[string]*Asset, err error) {
	parents = make(map[string]*Asset)
	for _, name := range names {
		prnt, err := getByName(ctx, name)
		if err != nil {
			return nil, err
		}
		parents[name] = &prnt

		hash, err := parents[name].Hash()
		if err != nil {
			return nil, errors.Wrapf(err, "hash %q parent %q", asset.Name, name)
		}

		asset.Parents = append(asset.Parents, Reference{
			Name: name,
			Hash: hash,
		})
	}

	return parents, nil
}

// path returns a slugged version of the asset name.
func (asset *Asset) path() (subPath string) {
	segments := []string{}
	remaining := asset.Name
	var file string
	for remaining != "." {
		remaining, file = path.Split(remaining)
		segments = append([]string{filenameRegexp.ReplaceAllString(file, "-")}, segments...)
		next := path.Dir(remaining) // drop any trailing slash
		if next == remaining {
			break
		}
		remaining = next
	}
	return filepath.Join(segments...)
}

// Write writes the asset to the given directory, using a slugged
// version of the name as the filename.  If getByHash is non-nil,
// Write will recurse through parent assets.  If written is non-nil,
// every filename written will be added to the map with a true value.
func (asset *Asset) Write(ctx context.Context, directory string, getByHash GetByBytes, written map[string]bool) (err error) {
	data, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	subPath := asset.path()
	path := filepath.Join(directory, subDir, subPath)
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}

	path = filepath.Join(directory, subPath)
	dir = filepath.Dir(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, asset.Data, 0666)
	if err != nil {
		return err
	}

	if written != nil {
		written[subPath] = true
		written[filepath.Join(subDir, subPath)] = true
	}

	if getByHash != nil {
		for _, hash := range asset.Parents {
			parent, err := getByHash(ctx, hash.Hash)
			if err != nil {
				return errors.Errorf("failed to retrieve %x by hash", hash)
			}

			err = (&parent).Write(ctx, directory, getByHash, written)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Read reads the asset from the given directory, using a slugged
// version of the name as the filename.
func (asset *Asset) Read(ctx context.Context, directory string, logger log.Logger) (err error) {
	name := asset.Name
	subPath := asset.path()
	path := filepath.Join(directory, subDir, subPath)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &asset)
	if err != nil {
		return err
	}
	if asset.Name != name {
		return errors.Errorf("name %q read from %q does not match the expected %q", asset.Name, path, name)
	}

	path = filepath.Join(directory, subPath)
	data, err = ioutil.ReadFile(path)
	if err == nil && !bytes.Equal(data, asset.Data) {
		if asset.Parents == nil {
			logger.Logf("%q data changed via %q", name, path)
		} else {
			parents := make([]string, 0, len(asset.Parents))
			for _, parent := range asset.Parents {
				parents = append(parents, parent.Name)
			}
			logger.Logf("%q data changed via %q, future changes to the previous parents (%s) will no longer update this asset", name, path, strings.Join(parents, ", "))
		}
		asset.Data = data
		asset.Parents = nil
		asset.Frozen = true
		asset.RebuildHelper = ConstantDataRebuilder(ctx, name, data, true)
	} else if asset.Frozen {
		logger.Logf("%q is frozen due to a previous user change", name)
		asset.RebuildHelper = ConstantDataRebuilder(ctx, name, asset.Data, true)
	}

	return nil
}
