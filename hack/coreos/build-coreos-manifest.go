// We want to keep the stream metadata stored "directly"
// in git so it's easy to read and validate.  This build
// script is invoked as part of the container build to
// inject the data into a ConfigMap that will be installed
// via CVO manifests into the target cluster, and maintained
// across upgrades.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	streamRHCOSJSON = "data/data/coreos/rhcos.json"
	streamFCOSJSON  = "data/data/coreos/fcos.json"
	fcosTAG         = "okd"
	dest            = "bin/manifests/coreos-bootimages.yaml"
)

func run() error {
	streamJSON := streamRHCOSJSON
	if tags, _ := os.LookupEnv("TAGS"); strings.Contains(tags, fcosTAG) {
		streamJSON = streamFCOSJSON
	}
	bootimages, err := os.ReadFile(streamJSON)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-config-operator",
			Name:      "coreos-bootimages",
			Annotations: map[string]string{
				"include.release.openshift.io/ibm-cloud-managed":              "true",
				"include.release.openshift.io/self-managed-high-availability": "true",
				"include.release.openshift.io/single-node-developer":          "true",
			},
		},
		Data: map[string]string{
			"releaseVersion": "0.0.1-snapshot",
			"stream":         string(bootimages),
		},
	}

	b, err := yaml.Marshal(cm)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	err = os.WriteFile(dest, b, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
