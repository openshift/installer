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
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/coreos/stream-metadata-go/stream/rhcos"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	instTypes "github.com/openshift/installer/pkg/types"
)

const (
	scosTAG           = "scos"
	dest              = "bin/manifests/coreos-bootimages.yaml"
	defaultStreamName = instTypes.OSImageStreamRHCOS9
)

func run() error {
	streamsFiles, err := getOSStreamPaths()
	if err != nil {
		return err
	}

	var defaultStreamFiles *streamPath
	if len(streamsFiles) == 0 {
		return fmt.Errorf("no OS image streams found")
	} else if len(streamsFiles) == 1 {
		defaultFile := streamsFiles[slices.Collect(maps.Keys(streamsFiles))[0]]
		defaultStreamFiles = &defaultFile
	} else {

		defaultPaths, ok := streamsFiles[string(defaultStreamName)]
		if !ok {
			return fmt.Errorf("no %v image streams found. %v is considered the default stream", defaultStreamName, defaultStreamName)
		}
		defaultStreamFiles = &defaultPaths
	}

	defaultStreamBytes, err := defaultStreamFiles.getBytes()
	if err != nil {
		return fmt.Errorf("could not read default stream bytes: %w", err)
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
			"stream":         string(defaultStreamBytes),
		},
	}

	streamsData := make(map[string]json.RawMessage)

	for name, files := range streamsFiles {
		streamBytes, err := files.getBytes()
		if err != nil {
			return fmt.Errorf("could not read %s stream bytes: %w", name, err)
		}
		streamsData[name] = streamBytes
	}
	streamsJsonData, err := json.MarshalIndent(streamsData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal streams data: %w", err)
	}
	cm.Data["streams"] = string(streamsJsonData)

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

type streamPath struct {
	streamFile      string
	marketplaceFile string
	name            string
}

func (s *streamPath) getBytes() ([]byte, error) {
	streamJson, err := os.ReadFile(s.streamFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream file: %w", err)
	}

	bootImageStream := stream.Stream{}
	if err := json.Unmarshal(streamJson, &bootImageStream); err != nil {
		return nil, fmt.Errorf("failed to unmarshal boot image stream: %w", err)
	}

	if s.marketplaceFile != "" {
		mktStream := marketplaceStream{}
		mktJSON, err := os.ReadFile(s.marketplaceFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open marketplace file: %w", err)
		}
		if err := json.Unmarshal(mktJSON, &mktStream); err != nil {
			return nil, fmt.Errorf("failed to unmarshal market stream: %w", err)
		}
		for name, arch := range bootImageStream.Architectures {
			if mkt, ok := mktStream[name]; ok {
				if arch.RHELCoreOSExtensions == nil {
					arch.RHELCoreOSExtensions = &rhcos.Extensions{}
				}
				arch.RHELCoreOSExtensions.Marketplace = mkt
				bootImageStream.Architectures[name] = arch
			}
		}
	}
	bootImgs, err := json.MarshalIndent(bootImageStream, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged boot image stream: %w", err)
	}
	return bootImgs, nil
}

func getOSStreamPaths() (map[string]streamPath, error) {
	name := "coreos"
	if isOKD() {
		name = "scos"
	}
	matches, err := filepath.Glob(fmt.Sprintf("data/data/coreos/%s*.json", name))
	if err != nil {
		return nil, fmt.Errorf("failed to find image streams: %v", err)
	}
	streams := make(map[string]streamPath)
	for _, match := range matches {
		filename := filepath.Base(match)
		if strings.HasPrefix(filename, name+"-") {
			// It's a stream

			streamName := strings.TrimPrefix(strings.TrimSuffix(filename, ".json"), name+"-")
			streamFile := streamPath{
				streamFile: match,
				name:       streamName,
			}

			marketPlaceFile := filepath.Join(filepath.Dir(match), "marketplace", filename)
			if _, err := os.Stat(marketPlaceFile); err == nil {
				streamFile.marketplaceFile = marketPlaceFile
			}
			streams[streamName] = streamFile
		} else if !strings.Contains(filename, "-") {
			// It's a plain no OSImageStream, ie SCOS
			streams[""] = streamPath{
				streamFile: match,
			}
		}
	}
	return streams, nil
}

func isOKD() bool {
	tags, _ := os.LookupEnv("TAGS")
	return strings.Contains(tags, scosTAG)
}

type marketplaceStream map[string]*rhcos.Marketplace

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
