// We want to keep the stream metadata stored "directly"
// in git so it's easy to read and validate.  This build
// script is invoked as part of the container build to
// inject the data into a ConfigMap that will be installed
// via CVO manifests into the target cluster, and maintained
// across upgrades.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/coreos/stream-metadata-go/stream/rhcos"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	streamRHCOSJSON            = "data/data/coreos/rhcos.json"
	streamSCOSJSON             = "data/data/coreos/scos.json"
	streamMarketplaceRHCOSJSON = "data/data/coreos/marketplace-rhcos.json"
	scosTAG                    = "scos"
	dest                       = "bin/manifests/coreos-bootimages.yaml"
)

func run() error {
	bootimages, err := getBootImages()
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

const rhcosStreamName = "rhcos-4.21"

func getBootImages() ([]byte, error) {
	var okd bool
	var streamJSONFile string
	var streamName string
	tags, _ := os.LookupEnv("TAGS")
	switch {
	case strings.Contains(tags, scosTAG):
		streamJSONFile = streamSCOSJSON
		okd = true
	default:
		streamJSONFile = streamRHCOSJSON
		streamName = rhcosStreamName
	}

	fileData, err := os.ReadFile(streamJSONFile)
	if err != nil {
		return nil, err
	}

	if okd {
		// okd does not yet have marketplace images, so we are done
		return fileData, nil
	}

	// Parse the multi-stream map and extract the target stream
	var store map[string]json.RawMessage
	if err := json.Unmarshal(fileData, &store); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stream store: %w", err)
	}
	streamJSON, ok := store[streamName]
	if !ok {
		return nil, fmt.Errorf("stream %q not found in stream store", streamName)
	}

	return mergeMarketplaceStream(streamJSON)
}

type marketplaceStream map[string]*rhcos.Marketplace

func mergeMarketplaceStream(streamJSON []byte) ([]byte, error) {
	mktStream := marketplaceStream{}
	mktJSON, err := os.ReadFile(streamMarketplaceRHCOSJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to open marketplace file: %w", err)
	}
	if err := json.Unmarshal(mktJSON, &mktStream); err != nil {
		return nil, fmt.Errorf("failed to unmarshal market stream: %w", err)
	}

	stream := stream.Stream{}
	if err := json.Unmarshal(streamJSON, &stream); err != nil {
		return nil, fmt.Errorf("failed to unmarshal boot image stream: %w", err)
	}

	for name, arch := range stream.Architectures {
		if mkt, ok := mktStream[name]; ok {
			if arch.RHELCoreOSExtensions == nil {
				arch.RHELCoreOSExtensions = &rhcos.Extensions{}
			}
			arch.RHELCoreOSExtensions.Marketplace = mkt
		}
	}

	bootImgs, err := json.Marshal(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged boot image stream: %w", err)
	}
	return bootImgs, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
