package manifests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type decodeFormat interface {
	newDecodedYaml(decoder *yaml.YAMLToJSONDecoder) (interface{}, error)
}

// Read a YAML file containing multiple YAML definitions of the same format
// Each specific format must be of type decodeFormat
func getFileMultipleYamls(filename string, decoder decodeFormat) ([]interface{}, error) {

	// TODO - this location must be defined by user input via cobra flags
	path := filepath.Join("./manifests", filename)

	contents, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("error reading file %s: %w", path, err)
		return nil, err
	}

	r := bytes.NewReader(contents)
	dec := yaml.NewYAMLToJSONDecoder(r)

	var outputList []interface{}
	for {
		decodedData, err := decoder.newDecodedYaml(dec)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			err = fmt.Errorf("error reading multiple yamls in file %s: %w", path, err)
			return nil, err
		}

		outputList = append(outputList, decodedData)
	}

	return outputList, nil
}
