package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/coreos/stream-metadata-go/stream/rhcos"

	"github.com/openshift/installer/pkg/rhcos/marketplace/azure"
)

const (
	streamRHCOSJSON            = "data/data/coreos/rhcos.json"
	streamMarketplaceRHCOSJSON = "data/data/coreos/marketplace-rhcos.json"

	x86   = "x86_64"
	arm64 = "aarch64"
)

// arch -> marketplace
type marketplaceStream map[string]*rhcos.Marketplace

func main() {
	ctx := context.Background()
	stream := marketplaceStream{}

	if err := stream.populate(ctx); err != nil {
		log.Fatalln("Failed to populate marketplace stream:", err)
	}

	if err := stream.write(); err != nil {
		log.Fatalln("Failed to write marketplace stream:", err)
	}
	log.Printf("Successfully wrote marketplace stream to %s", streamMarketplaceRHCOSJSON)
}

// populate gathers the marketplace images for each cloud
// and adds them to the marketplace stream data structure.
func (s marketplaceStream) populate(ctx context.Context) error {
	clouds := []func(ctx context.Context, arch string) error{
		s.azure,
	}

	for _, supportedArch := range []string{arm64, x86} {
		s[supportedArch] = &rhcos.Marketplace{}
		for _, populateCloud := range clouds {
			if err := populateCloud(ctx, supportedArch); err != nil {
				return err
			}
		}
	}
	return nil
}

// write serializes the marketplace stream to disk.
func (s marketplaceStream) write() error {
	contents, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling stream: %w", err)
	}

	// TODO(padillon): dumb question time, git is still complaining \ No newline at end of file
	// what am I doing wrong?
	contents = append(contents, []byte("\n")...)

	if err := os.WriteFile(streamMarketplaceRHCOSJSON, contents, 0644); err != nil {
		return fmt.Errorf("error writing stream: %w", err)
	}
	return nil
}

func (s marketplaceStream) azure(ctx context.Context, arch string) error {
	var err error

	s[arch].Azure = &rhcos.AzureMarketplace{}

	rel, err := getReleaseFromStream()
	if err != nil {
		return fmt.Errorf("failed to get release from rhcos stream: %w", err)
	}

	azClient, err := azure.NewStreamClient()
	if err != nil {
		return fmt.Errorf("failed to initialize azure client: %w", err)
	}

	if s[arch].Azure, err = azClient.Populate(ctx, arch, rel); err != nil {
		return err
	}

	return nil
}

// getXYFromStream obtains the X.Y version from rhcos.json.
func getReleaseFromStream() (string, error) {
	if rel, ok := os.LookupEnv("STREAM_RELEASE_OVERRIDE"); ok {
		log.Printf("Found STREAM_RELEASE_OVERRIDE %s", rel)
		return rel, nil
	}
	fileContents, err := os.ReadFile(streamRHCOSJSON)
	if err != nil {
		return "", err
	}

	st := &stream.Stream{}
	if err := json.Unmarshal(fileContents, st); err != nil {
		return "", fmt.Errorf("failed to unmarshal RHCOS stream: %w", err)
	}

	return strings.Split(st.Stream, "-")[1], nil
}
