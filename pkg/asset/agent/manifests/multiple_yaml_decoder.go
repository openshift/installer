package manifests

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type DecodeFormat interface {
	NewDecodedYaml(decoder *yaml.YAMLToJSONDecoder) (interface{}, error)
}

// Read a YAML file containing multiple YAML definitions of the same format
// Each specific format must be of type DecodeFormat
func GetMultipleYamls(contents []byte, decoder DecodeFormat) ([]interface{}, error) {

	r := bytes.NewReader(contents)
	dec := yaml.NewYAMLToJSONDecoder(r)

	var outputList []interface{}
	for {
		decodedData, err := decoder.NewDecodedYaml(dec)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading multiple YAMLs")
		}

		outputList = append(outputList, decodedData)
	}

	return outputList, nil
}
