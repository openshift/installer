package image

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/stream-metadata-go/arch"
	"github.com/coreos/stream-metadata-go/stream"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/rhcos/cache"
)

const (
	machineOsImageName   = "machine-os-images"
	coreOsFileName       = "/coreos/coreos-%s.iso"
	coreOsSha256FileName = "/coreos/coreos-%s.iso.sha256"
	coreOsStreamFileName = "/coreos/coreos-stream.json"
	//OcDefaultTries is the number of times to execute the oc command on failues
	OcDefaultTries = 5
	// OcDefaultRetryDelay is the time between retries
	OcDefaultRetryDelay = time.Second * 5
)

// Config is used to set up the retries for extracting the base ISO
type Config struct {
	MaxTries   uint
	RetryDelay time.Duration
}

// Release is the interface to use the oc command to the get image info
type Release interface {
	GetBaseIso(architecture string) (string, error)
	GetBaseIsoVersion(architecture string) (string, error)
	ExtractFile(image string, filename string, architecture string) ([]string, error)
}

type release struct {
	config       Config
	releaseImage string
	pullSecret   string
	mirrorConfig []mirror.RegistriesConfig
}

// NewRelease is used to set up the executor to run oc commands
func NewRelease(config Config, releaseImage string, pullSecret string, mirrorConfig []mirror.RegistriesConfig) Release {
	return &release{
		config:       config,
		releaseImage: releaseImage,
		pullSecret:   pullSecret,
		mirrorConfig: mirrorConfig,
	}
}

const (
	templateGetImage             = "oc adm release info --image-for=%s --filter-by-os=linux/%s --insecure=%t %s"
	templateGetImageWithIcsp     = "oc adm release info --image-for=%s --filter-by-os=linux/%s --insecure=%t --icsp-file=%s %s"
	templateImageExtract         = "oc image extract --path %s:%s --filter-by-os=linux/%s --confirm %s"
	templateImageExtractWithIcsp = "oc image extract --path %s:%s --filter-by-os=linux/%s --confirm --icsp-file=%s %s"
)

// ExtractFile extracts the specified file from the given image name, and store it in the cache dir.
func (r *release) ExtractFile(image string, filename string, architecture string) ([]string, error) {
	imagePullSpec, err := r.getImageFromRelease(image, architecture)
	if err != nil {
		return nil, err
	}

	cacheDir, err := cache.GetCacheDir(cache.FilesDataType, cache.AgentApplicationName)
	if err != nil {
		return nil, err
	}

	path, err := r.extractFileFromImage(imagePullSpec, filename, cacheDir, architecture)
	if err != nil {
		return nil, err
	}
	return path, err
}

// Get the CoreOS ISO from the releaseImage.
func (r *release) GetBaseIso(architecture string) (string, error) {
	// Get the machine-os-images pullspec from the release and use that to get the CoreOS ISO
	image, err := r.getImageFromRelease(machineOsImageName, architecture)
	if err != nil {
		return "", err
	}

	cacheDir, err := cache.GetCacheDir(cache.ImageDataType, cache.AgentApplicationName)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf(coreOsFileName, architecture)
	// Check if file is already cached
	cachedFile, err := cache.GetFileFromCache(path.Base(filename), cacheDir)
	if err != nil {
		return "", err
	}
	if cachedFile != "" {
		logrus.Info("Verifying cached file")
		valid, err := r.verifyCacheFile(image, cachedFile, architecture)
		if err != nil {
			return "", err
		}
		if valid {
			logrus.Infof("Using cached Base ISO %s", cachedFile)
			return cachedFile, nil
		}
	}

	// Get the base ISO from the payload
	path, err := r.extractFileFromImage(image, filename, cacheDir, architecture)
	if err != nil {
		return "", err
	}
	logrus.Infof("Base ISO obtained from release and cached at %s", path)
	return path[0], err
}

func (r *release) GetBaseIsoVersion(architecture string) (string, error) {
	files, err := r.ExtractFile(machineOsImageName, coreOsStreamFileName, architecture)
	if err != nil {
		return "", err
	}

	if len(files) > 1 {
		return "", fmt.Errorf("too many files found for %s", coreOsStreamFileName)
	}

	rawData, err := os.ReadFile(files[0])
	if err != nil {
		return "", err
	}

	var st stream.Stream
	if err := json.Unmarshal(rawData, &st); err != nil {
		return "", errors.Wrap(err, "failed to parse CoreOS stream metadata")
	}

	streamArch, err := st.GetArchitecture(architecture)
	if err != nil {
		return "", err
	}

	if metal, ok := streamArch.Artifacts["metal"]; ok {
		return metal.Release, nil
	}

	return "", errors.New("unable to determine CoreOS release version")
}

func (r *release) getImageFromRelease(imageName string, architecture string) (string, error) {
	// This requires the 'oc' command so make sure its available
	_, err := exec.LookPath("oc")
	var cmd string
	if err != nil {
		if len(r.mirrorConfig) > 0 {
			logrus.Warning("Unable to validate mirror config because \"oc\" command is not available")
		} else {
			logrus.Debug("Skipping ISO extraction; \"oc\" command is not available")
		}
		return "", err
	}

	archName := arch.GoArch(architecture)

	if len(r.mirrorConfig) > 0 {
		logrus.Debugf("Using mirror configuration")
		icspFile, err := getIcspFileFromRegistriesConfig(r.mirrorConfig)
		if err != nil {
			return "", err
		}
		defer removeIcspFile(icspFile)
		cmd = fmt.Sprintf(templateGetImageWithIcsp, imageName, archName, true, icspFile, r.releaseImage)
	} else {
		cmd = fmt.Sprintf(templateGetImage, imageName, archName, true, r.releaseImage)
	}

	logrus.Debugf("Fetching image from OCP release (%s)", cmd)
	image, err := execute(r.pullSecret, cmd)
	if err != nil {
		if strings.Contains(err.Error(), "unknown flag: --icsp-file") {
			logrus.Warning("Using older version of \"oc\" that does not support mirroring")
		}
		return "", err
	}

	return image, nil
}

func (r *release) extractFileFromImage(image, file, cacheDir string, architecture string) ([]string, error) {
	var cmd string
	archName := arch.GoArch(architecture)
	if len(r.mirrorConfig) > 0 {
		icspFile, err := getIcspFileFromRegistriesConfig(r.mirrorConfig)
		if err != nil {
			return nil, err
		}
		defer removeIcspFile(icspFile)
		cmd = fmt.Sprintf(templateImageExtractWithIcsp, file, cacheDir, archName, icspFile, image)
	} else {
		cmd = fmt.Sprintf(templateImageExtract, file, cacheDir, archName, image)
	}
	path := filepath.Join(cacheDir, path.Base(file))
	// Remove file if it exists
	if err := removeCacheFile(path); err != nil {
		return nil, err
	}

	logrus.Debugf("extracting %s to %s, %s", file, cacheDir, cmd)
	_, err := retry.Do(r.config.MaxTries, r.config.RetryDelay, execute, r.pullSecret, cmd)
	if err != nil {
		return nil, err
	}

	// Make sure file(s) exist after extraction
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	if matches == nil {
		return nil, fmt.Errorf("file %s was not found", file)
	}

	return matches, nil
}

// Get hash from rhcos.json
func getHashFromInstaller(architecture string) (bool, string) {
	// Get hash from metadata in the installer
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	st, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return false, ""
	}

	streamArch, err := st.GetArchitecture(architecture)
	if err != nil {
		return false, ""
	}
	if artifacts, ok := streamArch.Artifacts["metal"]; ok {
		if format, ok := artifacts.Formats["iso"]; ok {
			return true, format.Disk.Sha256
		}
	}

	return false, ""
}

func matchingHash(imageSha []byte, sha string) bool {
	decoded, err := hex.DecodeString(sha)
	if err == nil && bytes.Equal(imageSha, decoded) {
		return true
	}

	return false
}

// Check if there is a different base ISO in the release payload
func (r *release) verifyCacheFile(image, file, architecture string) (bool, error) {
	// Get hash of cached file
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}
	fileSha := h.Sum(nil)

	// Check if the hash of cached file matches hash in rhcos.json
	found, rhcosSha := getHashFromInstaller(architecture)
	if found && matchingHash(fileSha, rhcosSha) {
		logrus.Debug("Found matching hash in installer metadata")
		return true, nil
	}

	// If no match, get the file containing the coreos sha256 and compare that
	tempDir, err := os.MkdirTemp("", "cache")
	if err != nil {
		return false, err
	}

	defer os.RemoveAll(tempDir)

	shaFilename := fmt.Sprintf(coreOsSha256FileName, architecture)
	shaFile, err := r.extractFileFromImage(image, shaFilename, tempDir, architecture)
	if err != nil {
		logrus.Debug("Could not get SHA from payload for cache comparison")
		return false, nil
	}

	payloadSha, err := os.ReadFile(shaFile[0])
	if err != nil {
		return false, err
	}
	if matchingHash(fileSha, string(payloadSha)) {
		logrus.Debugf("Found matching hash in %s", shaFilename)
		return true, nil
	}

	logrus.Debugf("Cached file %s is not most recent", file)
	return false, nil
}

// Remove any existing files in the cache.
func removeCacheFile(path string) error {
	matches, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, file := range matches {
		if err = os.Remove(file); err != nil {
			return err
		}
		logrus.Debugf("Removed file %s", file)
	}
	return nil
}

func execute(pullSecret, command string) (string, error) {
	ps, err := os.CreateTemp("", "registry-config")
	if err != nil {
		return "", err
	}
	defer func() {
		ps.Close()
		os.Remove(ps.Name())
	}()
	_, err = ps.Write([]byte(pullSecret))
	if err != nil {
		return "", err
	}
	// flush the buffer to ensure the file can be read
	ps.Close()
	executeCommand := command[:] + " --registry-config=" + ps.Name()
	args := strings.Split(executeCommand, " ")

	var stdoutBytes, stderrBytes bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err = cmd.Run()
	if err == nil {
		return strings.TrimSpace(stdoutBytes.String()), nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		err = fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s", executeCommand, exitErr.ExitCode(), stdoutBytes.String(), stderrBytes.String())
	} else {
		err = fmt.Errorf("command '%s' failed: %w", executeCommand, err)
	}
	return "", err
}

// Create a temporary file containing the ImageContentPolicySources
func getIcspFileFromRegistriesConfig(mirrorConfig []mirror.RegistriesConfig) (string, error) {
	contents, err := getIcspContents(mirrorConfig)
	if err != nil {
		return "", err
	}
	if contents == nil {
		logrus.Debugf("No registry entries to build ICSP file")
		return "", nil
	}

	icspFile, err := os.CreateTemp("", "icsp-file")
	if err != nil {
		return "", err
	}

	if _, err := icspFile.Write(contents); err != nil {
		icspFile.Close()
		os.Remove(icspFile.Name())
		return "", err
	}
	icspFile.Close()

	return icspFile.Name(), nil
}

// Convert the data in registries.conf into ICSP format
func getIcspContents(mirrorConfig []mirror.RegistriesConfig) ([]byte, error) {
	icsp := operatorv1alpha1.ImageContentSourcePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1alpha1.SchemeGroupVersion.String(),
			Kind:       "ImageContentSourcePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "image-policy",
			// not namespaced
		},
	}

	icsp.Spec.RepositoryDigestMirrors = make([]operatorv1alpha1.RepositoryDigestMirrors, len(mirrorConfig))
	for i, mirrorRegistries := range mirrorConfig {
		icsp.Spec.RepositoryDigestMirrors[i] = operatorv1alpha1.RepositoryDigestMirrors{Source: mirrorRegistries.Location, Mirrors: []string{mirrorRegistries.Mirror}}
	}

	// Convert to json first so json tags are handled
	jsonData, err := json.Marshal(&icsp)
	if err != nil {
		return nil, err
	}
	contents, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func removeIcspFile(filename string) {
	if filename != "" {
		os.Remove(filename)
	}
}
