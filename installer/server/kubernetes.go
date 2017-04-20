package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"k8s.io/kubernetes/pkg/apis/abac"
	"k8s.io/kubernetes/pkg/apis/abac/v1beta1"
	"k8s.io/kubernetes/pkg/runtime/serializer/json"

	"github.com/coreos/tectonic-installer/installer/server/asset"
)

type secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Type       string            `yaml:"type"`
	Data       map[string]string `yaml:"data"`
}

// secretFromAssets writes named Assets into a Kubernetes Secret.
func secretFromAssets(name, namespace string, assetNames []string, assets []asset.Asset) ([]byte, error) {
	data := make(map[string]string)
	for _, an := range assetNames {
		a, err := asset.Find(assets, an)
		if err != nil {
			return []byte{}, err
		}
		data[filepath.Base(a.Name())] = base64.StdEncoding.EncodeToString(a.Data())
	}
	return yaml.Marshal(secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
		Data: data,
	})
}

// secretFromData writes key-value secrets into a Kubernetes Secret.
func secretFromData(name, namespace string, secrets map[string][]byte) ([]byte, error) {
	data := make(map[string]string)
	for key, value := range secrets {
		data[key] = base64.StdEncoding.EncodeToString(value)
	}
	return yaml.Marshal(secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
		Data: data,
	})
}

type configMap struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

// configMapFromData writes key-value data into a Kubernetes ConfigMap.
func configMapFromData(name, namespace string, data map[string]string) ([]byte, error) {
	return yaml.Marshal(configMap{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
		Data: data,
	})
}

// abacPolicyToJSONL encodes ABAC policies as a line delimited JSON file, the
// format the API server consumes.
func abacPolicyToJSONL(policies []abac.Policy) ([]byte, error) {
	// Use the Kubernetes serializers to add the appropriate "kind" and "metadata"
	// fields.
	encoder := abac.Codecs.EncoderForVersion(
		json.NewSerializer(
			json.DefaultMetaFactory,
			abac.Scheme, // For some reason the ABAC types are registered on abac.Scheme instead of api.Scheme.
			abac.Scheme,
			false,
		),
		v1beta1.GroupVersion,
	)
	p := new(bytes.Buffer)
	for _, policy := range policies {
		// Encode automatically adds a newline.
		if err := encoder.Encode(&policy, p); err != nil {
			return nil, fmt.Errorf("failed to encode policy: %v", err)
		}
	}
	return p.Bytes(), nil
}
