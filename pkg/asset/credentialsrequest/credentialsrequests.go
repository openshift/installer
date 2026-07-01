package credentialsrequest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	ccov1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// providerSpecDecoderFunc decodes a raw provider spec into a
// platform-specific type. Each platform registers its decoder via
// registerProviderSpecDecoder in an init() function.
type providerSpecDecoderFunc func(raw *runtime.RawExtension) (interface{}, error)

var providerSpecDecoders = map[string]providerSpecDecoderFunc{}

// registerProviderSpecDecoder registers a provider spec decoder for
// the given cloud platform. Called from platform-specific init() functions.
func registerProviderSpecDecoder(cloud string, fn providerSpecDecoderFunc) {
	providerSpecDecoders[cloud] = fn
}

// CredentialRequest holds a parsed CredentialsRequest with platform-agnostic
// common fields and a platform-specific ProviderSpec.
type CredentialRequest struct {
	Name                string
	Namespace           string
	SecretRefName       string
	SecretRefNamespace  string
	ServiceAccountNames []string
	ProviderSpec        interface{}
}

const (
	credentialsRequestFileName        = "99_credentials-request_%s.yaml"
	credentialsRequestFileNamePattern = "99_credentials-request_*.yaml"
	credentialsRequestDir             = "manifests"
)

// CredentialsRequests is an asset that extracts and parses
// CredentialsRequests from the release image for the target platform.
// It is consumed by infrastructure provisioning (e.g. IAM roles) and
// manifest generation (credential secrets).
type CredentialsRequests struct {
	FileList []*asset.File
	Requests []CredentialRequest
}

var _ asset.WritableAsset = (*CredentialsRequests)(nil)

// Name returns a human-friendly name for the asset.
func (*CredentialsRequests) Name() string {
	return "Credentials Requests"
}

// Dependencies returns the assets required to extract credentials requests.
func (*CredentialsRequests) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&releaseimage.Image{},
	}
}

// Generate extracts CredentialsRequests from the release image for the
// target platform and parses them into platform-specific structs.
func (cr *CredentialsRequests) Generate(ctx context.Context, dependencies asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	ri := &releaseimage.Image{}
	dependencies.Get(ic, ri)

	cloud := ic.Config.Platform.Name()
	switch cloud {
	case awstypes.Name:
		if !ic.Config.Platform.AWS.IsSTSManaged() {
			return nil
		}
	default:
		logrus.Debugf("Skipped generating Credentials Requests for unsupported platform %s", cloud)
		return nil
	}

	pullSecret := ic.Config.PullSecret
	mirrorConfig := types.BuildMirrorConfig(ic.Config)

	icDir, err := os.MkdirTemp("", "install-config")
	if err != nil {
		return fmt.Errorf("failed to create temp directory for install-config: %w", err)
	}
	defer os.RemoveAll(icDir)

	icData, err := yaml.Marshal(ic.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal install-config: %w", err)
	}
	icPath := filepath.Join(icDir, "install-config.yaml")
	if err := os.WriteFile(icPath, icData, 0o600); err != nil {
		return fmt.Errorf("failed to write install-config: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "credentials-requests")
	if err != nil {
		return fmt.Errorf("failed to create temp directory for credentials requests: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := extractCredentialsRequests(ctx, cloud, pullSecret, ri.PullSpec, tmpDir, mirrorConfig, icPath); err != nil {
		return fmt.Errorf("failed to extract credentials requests from release image: %w", err)
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read credentials requests directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !(strings.HasSuffix(entry.Name(), ".yaml") || strings.HasSuffix(entry.Name(), ".yml")) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(tmpDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}
		req, err := parseCredentialRequestBytes(data, entry.Name())
		if err != nil {
			return fmt.Errorf("failed to parse credentials request %s: %w", entry.Name(), err)
		}
		if req == nil {
			continue
		}
		cr.Requests = append(cr.Requests, *req)
		cr.FileList = append(cr.FileList, &asset.File{
			Filename: filepath.Join(credentialsRequestDir, fmt.Sprintf(credentialsRequestFileName, req.Namespace+"-"+req.Name)),
			Data:     data,
		})
	}
	if len(cr.Requests) == 0 {
		return fmt.Errorf("no credentials requests found in release image")
	}
	logrus.Infof("Extracted %d credentials requests from release image", len(cr.Requests))
	return nil
}

// Files returns the files generated by the asset.
func (cr *CredentialsRequests) Files() []*asset.File {
	return cr.FileList
}

// Load reads previously persisted credentials requests from disk.
func (cr *CredentialsRequests) Load(f asset.FileFetcher) (bool, error) {
	fileList, err := f.FetchByPattern(filepath.Join(credentialsRequestDir, credentialsRequestFileNamePattern))
	if err != nil {
		return false, err
	}
	if len(fileList) == 0 {
		return false, nil
	}

	cr.FileList = fileList
	cr.Requests = make([]CredentialRequest, 0, len(fileList))
	for _, file := range fileList {
		req, err := parseCredentialRequestBytes(file.Data, file.Filename)
		if err != nil {
			return false, fmt.Errorf("failed to parse %s: %w", file.Filename, err)
		}
		if req != nil {
			cr.Requests = append(cr.Requests, *req)
		}
	}
	return true, nil
}

// extractCredentialsRequests calls oc to extract credentials requests
// for the specified cloud platform from the release image.
func extractCredentialsRequests(ctx context.Context, cloud, pullSecret, releaseImage, destDir string, mirrorConfig types.MirrorConfig, icPath string) error {
	if _, err := exec.LookPath("oc"); err != nil {
		return fmt.Errorf("oc command not found in PATH: %w", err)
	}

	cmd := []string{
		"oc", "adm", "release", "extract",
		"--credentials-requests",
		"--cloud=" + cloud,
		"--included=true",
		"--install-config=" + icPath,
		"--to=" + destDir,
	}

	if mirrorConfig.HasMirrors() {
		mirrorArg, cleanup, err := getMirrorArg(mirrorConfig)
		if err != nil {
			return err
		}
		if mirrorArg != "" {
			defer cleanup()
			cmd = append(cmd, mirrorArg)
		}
	}

	cmd = append(cmd, releaseImage)

	logrus.Debugf("Extracting credentials requests: %s", strings.Join(cmd, " "))
	_, err := executeOC(ctx, pullSecret, cmd)
	return err
}

// parseCredentialRequestBytes parses a single CredentialsRequest YAML
// document, trying all registered provider spec decoders.
func parseCredentialRequestBytes(data []byte, filename string) (*CredentialRequest, error) {
	cr := &ccov1.CredentialsRequest{}
	if err := yaml.Unmarshal(data, cr); err != nil {
		return nil, fmt.Errorf("failed to parse %s as CredentialsRequest: %w", filename, err)
	}
	if cr.Spec.ProviderSpec == nil {
		return nil, fmt.Errorf("%s has no provider spec", filename)
	}
	var providerSpec interface{}
	for _, decoder := range providerSpecDecoders {
		spec, err := decoder(cr.Spec.ProviderSpec)
		if err == nil {
			providerSpec = spec
			break
		}
	}
	if providerSpec == nil {
		return nil, fmt.Errorf("no registered decoder matched provider spec in %s", filename)
	}
	return &CredentialRequest{
		Name:                cr.Name,
		Namespace:           cr.Namespace,
		SecretRefName:       cr.Spec.SecretRef.Name,
		SecretRefNamespace:  cr.Spec.SecretRef.Namespace,
		ServiceAccountNames: cr.Spec.ServiceAccountNames,
		ProviderSpec:        providerSpec,
	}, nil
}

// executeOC runs an oc command with the pull secret written to a temp file
// for registry authentication. This is a local copy of the pattern from
// pkg/asset/agent/oc.go to avoid cross-package dependencies.
func executeOC(ctx context.Context, pullSecret string, command []string) (string, error) {
	ps, err := os.CreateTemp("", "registry-config")
	if err != nil {
		return "", err
	}
	defer func() {
		ps.Close()
		os.Remove(ps.Name()) //nolint:gosec
	}()
	if _, err := ps.Write([]byte(pullSecret)); err != nil {
		return "", err
	}
	ps.Close()

	registryConfig := "--registry-config=" + ps.Name()
	command = append(command, registryConfig)

	var stdoutBytes, stderrBytes bytes.Buffer
	cmd := exec.CommandContext(ctx, command[0], command[1:]...) //nolint:gosec
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s",
				strings.Join(command, " "), exitErr.ExitCode(), stdoutBytes.String(), stderrBytes.String())
		}
		return "", fmt.Errorf("command '%s' failed: %w", strings.Join(command, " "), err)
	}
	return strings.TrimSpace(stdoutBytes.String()), nil
}

// getMirrorArg creates a temporary ICSP file from mirror config for
// use with oc commands. This is a local copy of the pattern from
// pkg/asset/rhcos/releaseextract.go.
func getMirrorArg(mirrorConfig types.MirrorConfig) (string, func(), error) {
	if !mirrorConfig.HasMirrors() {
		return "", nil, nil
	}

	contents, err := mirrorConfig.GetICSPContents()
	if err != nil {
		return "", nil, err
	}

	icspFile, err := os.CreateTemp("", "icsp-file")
	if err != nil {
		return "", nil, err
	}

	if _, err := icspFile.Write(contents); err != nil {
		icspFile.Close()
		os.Remove(icspFile.Name()) //nolint:gosec
		return "", nil, err
	}
	icspFile.Close()

	remove := func() {
		os.Remove(icspFile.Name()) //nolint:gosec
	}

	return "--icsp-file=" + icspFile.Name(), remove, nil
}
