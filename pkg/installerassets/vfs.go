package installerassets

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/openshift/installer/data"
)

func addAssetDefaults(defaults map[string]Defaulter, base string, rel string) error {
	uri := path.Join(base, rel)
	file, err := data.Assets.Open(uri)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		children, err := file.Readdir(0)
		if err != nil {
			return err
		}
		file.Close()

		for _, childInfo := range children {
			name := childInfo.Name()
			err = addAssetDefaults(defaults, base, path.Join(rel, name))
			if err != nil {
				return err
			}
		}

		return nil
	}

	if path.Base(rel) == "OWNERS" {
		return nil
	}

	defaults[rel] = func(ctx context.Context) ([]byte, error) {
		file, err := data.Assets.Open(uri)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		return ioutil.ReadAll(file)
	}

	return nil
}

func init() {
	err := addAssetDefaults(Defaults, "bootstrap", "")
	if err != nil {
		panic(err)
	}
}
