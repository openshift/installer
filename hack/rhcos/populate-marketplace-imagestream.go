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
	x86   = "x86_64"
	arm64 = "aarch64"
)

type streamConfig struct {
	name       string
	inputFile  string
	outputFile string
}

var (
	streamRHEL9 = streamConfig{
		name:       "rhel-9",
		inputFile:  "data/data/coreos/coreos-rhel-9.json",
		outputFile: "data/data/coreos/marketplace/coreos-rhel-9.json",
	}
	streamRHEL10 = streamConfig{
		name:       "rhel-10",
		inputFile:  "data/data/coreos/coreos-rhel-10.json",
		outputFile: "data/data/coreos/marketplace/coreos-rhel-10.json",
	}
)

// arch -> marketplace
type marketplaceStream map[string]*rhcos.Marketplace

func main() {
	ctx := context.Background()

	rhel9 := marketplaceStream{}
	if err := rhel9.populate(ctx, streamRHEL9); err != nil {
		log.Fatalln("Failed to populate RHEL 9 marketplace stream:", err)
	}
	if err := rhel9.write(streamRHEL9); err != nil {
		log.Fatalln("Failed to write RHEL 9 marketplace stream:", err)
	}
	log.Printf("Successfully wrote marketplace stream to %s", streamRHEL9.outputFile)

	rhel10 := marketplaceStream{}
	if err := rhel10.populateWithFallback(ctx, streamRHEL10, rhel9); err != nil {
		log.Fatalln("Failed to populate RHEL 10 marketplace stream:", err)
	}
	if err := rhel10.write(streamRHEL10); err != nil {
		log.Fatalln("Failed to write RHEL 10 marketplace stream:", err)
	}
	log.Printf("Successfully wrote marketplace stream to %s", streamRHEL10.outputFile)
}

// populate gathers the marketplace images for each cloud
// and adds them to the marketplace stream data structure.
func (s marketplaceStream) populate(ctx context.Context, cfg streamConfig) error {
	clouds := []func(ctx context.Context, arch string, cfg streamConfig) error{
		s.azure,
	}

	for _, supportedArch := range []string{arm64, x86} {
		s[supportedArch] = &rhcos.Marketplace{}
		for _, populateCloud := range clouds {
			if err := populateCloud(ctx, supportedArch, cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

// populateWithFallback attempts to populate marketplace data for each architecture.
// Any individual image types not found are filled in from the fallback stream.
func (s marketplaceStream) populateWithFallback(ctx context.Context, cfg streamConfig, fallback marketplaceStream) error {
	for _, supportedArch := range []string{arm64, x86} {
		s[supportedArch] = &rhcos.Marketplace{}

		if err := s.azure(ctx, supportedArch, cfg); err != nil {
			return err
		}

		if fb, ok := fallback[supportedArch]; ok && fb != nil {
			azure.FillMissing(s[supportedArch].Azure, fb.Azure)
		}
	}
	return nil
}

// write serializes the marketplace stream to disk.
func (s marketplaceStream) write(cfg streamConfig) error {
	contents, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling stream: %w", err)
	}

	contents = append(contents, []byte("\n")...)

	if err := os.WriteFile(cfg.outputFile, contents, 0644); err != nil {
		return fmt.Errorf("error writing stream: %w", err)
	}
	return nil
}

func (s marketplaceStream) azure(ctx context.Context, arch string, cfg streamConfig) error {
	var err error

	s[arch].Azure = &rhcos.AzureMarketplace{}

	rel, err := getReleaseFromStream(cfg)
	if err != nil {
		return fmt.Errorf("failed to get release from %s rhcos stream: %w", cfg.name, err)
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
func getReleaseFromStream(cfg streamConfig) (string, error) {
	if rel, ok := os.LookupEnv("STREAM_RELEASE_OVERRIDE"); ok {
		log.Printf("Found STREAM_RELEASE_OVERRIDE %s", rel)
		return rel, nil
	}
	fileContents, err := os.ReadFile(cfg.inputFile)
	if err != nil {
		return "", err
	}

	st := &stream.Stream{}
	if err := json.Unmarshal(fileContents, st); err != nil {
		return "", fmt.Errorf("failed to unmarshal RHCOS stream: %w", err)
	}

	return strings.Split(st.Stream, "-")[1], nil
}
