// Package asset abstracts generated asset representations.
package asset

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path"
)

// Asset is a named byte slice.
type Asset interface {
	Name() string
	Data() []byte
}

// An asset is an in-memory Asset.
type asset struct {
	name string
	data []byte
}

// New returns a new Asset.
func New(name string, data []byte) Asset {
	return asset{
		name: name,
		data: data,
	}
}

// Name returns the name of the asset.
func (a asset) Name() string {
	return a.name
}

// Data returns the byte contents of the asset.
func (a asset) Data() []byte {
	return a.data
}

// Find returns the Asset with the given name.
func Find(assets []Asset, name string) (Asset, error) {
	for _, a := range assets {
		if a.Name() == name {
			return a, nil
		}
	}
	return asset{}, fmt.Errorf("Asset %q was not found", name)
}

// Replace inserts the given Asset or replaces an Asset that has the same name.
func Replace(assets []Asset, asset Asset) ([]Asset, error) {
	for i, a := range assets {
		if a.Name() == asset.Name() {
			assets[i] = asset
			return assets, nil
		}
	}
	return assets, fmt.Errorf("Asset %q was not found", asset.Name())
}

// ZipAssets zips and provides the result back as a []byte
func ZipAssets(assets []Asset) ([]byte, error) {
	// Create a buffer to write the zip archive to
	buf := new(bytes.Buffer)
	// Create a zip archive Writer
	zw := zip.NewWriter(buf)
	// Add asset files to the archive
	for _, asset := range assets {
		f, err := zw.Create(path.Join("assets", asset.Name()))
		if err != nil {
			return nil, fmt.Errorf("Cannot add file to zip: %v", err)
		}
		_, err = f.Write(asset.Data())
		if err != nil {
			return nil, fmt.Errorf("Cannot write data to file: %v", err)
		}
	}
	// Close the archive
	err := zw.Close()
	if err != nil {
		return nil, fmt.Errorf("Cannot close zip.Writer: %v", err)
	}
	return buf.Bytes(), nil
}
