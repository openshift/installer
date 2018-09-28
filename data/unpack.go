package data

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// Unpack unpacks the assets from this package into a target directory.
func Unpack(base string, uri string) (err error) {
	file, err := Assets.Open(uri)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		os.Mkdir(base, 0777)
		children, err := file.Readdir(0)
		if err != nil {
			return err
		}
		file.Close()

		for _, childInfo := range children {
			name := childInfo.Name()
			err = Unpack(filepath.Join(base, name), path.Join(uri, name))
			if err != nil {
				return err
			}
		}
		return nil
	}

	out, err := os.OpenFile(base, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}
