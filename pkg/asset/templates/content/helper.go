package content

import (
	"io/ioutil"
	"path"

	"github.com/openshift/installer/data"
)

const (
	// TemplateDir is the target directory for all template assets' files
	TemplateDir     = "templates"
	bootkubeDataDir = "manifests/bootkube/"
	tectonicDataDir = "manifests/tectonic/"
)

// GetBootkubeTemplate returns the contents of the file in bootkube data dir
func GetBootkubeTemplate(uri string) ([]byte, error) {
	return getFileContents(path.Join(bootkubeDataDir, uri))
}

// GetTectonicTemplate returns the contents of the file in tectonic data dir
func GetTectonicTemplate(uri string) ([]byte, error) {
	return getFileContents(path.Join(tectonicDataDir, uri))
}

// getFileContents the content of the given URI, assuming that it's a file
func getFileContents(uri string) ([]byte, error) {
	file, err := data.Assets.Open(uri)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
