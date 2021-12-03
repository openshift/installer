// Package plugins is collection of all the terraform plugins that are used/required by installer.
package plugins

import (
	"archive/tar"
	"bufio"
	"compress/bzip2"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const TerraformResourceFile = "bin.tbz2"

var (
	// KnownPlugins is a map of all the known plugin names to their exec functions.
	KnownPlugins = map[string]string{}

	//go:embed bin.tbz2
	TerraformResourceHandle embed.FS
)

// XXX: find better names for these

func GetPluginPath() string {
	userCacheDir, _ := os.UserCacheDir()
	return filepath.Join(userCacheDir, "openshift-installer", "terraform")
}

func GetPluginBinPath() string {
	return filepath.Join(GetPluginPath(), "bin")
}

// XXX: This should probably go in the terraform package. Keep here for now.
func GetTerraformPath() (string, error) {
	terraformPath := filepath.Join(GetPluginBinPath(), "terraform")
	_, err := os.Stat(terraformPath)
	if err != nil {
		return "", err
	}

	return terraformPath, nil
}

// Write out the terraform binary and providers. These are symlinked when
// terraform exec runs.
// XXX: Move this out of init() function.
func init() {
	f, err := TerraformResourceHandle.Open(TerraformResourceFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s", TerraformResourceFile))
	}

	br := bufio.NewReader(f)
	cr := bzip2.NewReader(br)
	tbz := tar.NewReader(cr)

	pluginDir := GetPluginPath()
	for {
		hdr, err := tbz.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err.Error())
		}

		filePath := filepath.Join(pluginDir, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeReg:
			fp, err := os.Create(filePath)
			if err != nil {
				panic(fmt.Sprintf("Failed to create file %s: %s", filePath, err))
			}

			_, err = io.Copy(fp, tbz)
			if err != nil {
				panic(fmt.Sprintf("Failed to copy to file %s: %s", filePath, err))
			}

			err = os.Chmod(filePath, os.FileMode(hdr.Mode))
			if err != nil {
				panic(fmt.Sprintf("Failed to set permissions on file %s: %s", filePath, err))
			}

			basePath := filepath.Base(filePath)
			if basePath != "terraform" {
				KnownPlugins[filepath.Base(filePath)] = filePath
			}
			fp.Close()
		case tar.TypeDir:
			err = os.MkdirAll(filePath, os.FileMode(hdr.Mode))
			if err != nil {
				panic(fmt.Sprintf("Failed to create directory: %s", err))
			}
		default:
			panic(fmt.Sprintf("Failed to parse tar type: %c in file: %s", hdr.Typeflag, hdr.Name))
		}
	}
}
