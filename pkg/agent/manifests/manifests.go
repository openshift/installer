package manifests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// getFileData reads a YAML file and unmarshals the contents
func getFileData(fileName string, output interface{}) error {

	path := filepath.Join("./manifests/", fileName)

	contents, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("error reading file %s: %w", path, err)
	} else if err = yaml.Unmarshal(contents, output); err != nil {
		err = fmt.Errorf("error unmarshalling contents of %s: %w", path, err)
	}

	return err
}

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

// GetPullSecret gets the pull secret from the manifest file
func GetPullSecret() (string, error) {
	var secret corev1.Secret
	if err := getFileData("pull-secret.yaml", &secret); err != nil {
		return "", err
	}

	pullSecret := secret.StringData[".dockerconfigjson"]
	return pullSecret, nil
}

// GetAgentClusterInstall gets the AgentClusterInstall resource from the
// manifest file
func GetAgentClusterInstall() (hiveext.AgentClusterInstall, error) {
	var aci hiveext.AgentClusterInstall
	if err := getFileData("agent-cluster-install.yaml", &aci); err != nil {
		return aci, err
	}

	return aci, nil
}

// GetInfraEnv gets the InfraEnv resource from the manifest file
func GetInfraEnv() (aiv1beta1.InfraEnv, error) {
	var infraEnv aiv1beta1.InfraEnv
	if err := getFileData("infraenv.yaml", &infraEnv); err != nil {
		return infraEnv, err
	}

	return infraEnv, nil
}
