package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/assisted-service/pkg/executer"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	machineOsImageName = "machine-os-images"
	coreOsFileName     = "/coreos/coreos-%s.iso"
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
	GetBaseIso(log logrus.FieldLogger, releaseImage string, pullSecret string, mirrorConfig []mirror.RegistriesConfig, architecture string) (string, error)
}

type release struct {
	executer executer.Executer
	config   Config
}

// NewRelease is used to set up the executor to run oc commands
func NewRelease(executer executer.Executer, config Config) Release {
	return &release{executer: executer, config: config}
}

const (
	templateGetImage             = "oc adm release info --image-for=%s --insecure=%t %s"
	templateImageExtract         = "oc image extract --file=%s --path /:/%s --confirm %s"
	templateImageExtractWithIcsp = "oc image extract --file=%s --path /:/%s --confirm --icsp-file=%s %s"
)

// Get the CoreOS ISO from the releaseImage
func (r *release) GetBaseIso(log logrus.FieldLogger, releaseImage string, pullSecret string, mirrorConfig []mirror.RegistriesConfig, architecture string) (string, error) {

	// Get the machine-os-images pullspec from the release and use that to get the CoreOS ISO
	image, err := r.getImageFromRelease(log, machineOsImageName, releaseImage, pullSecret, len(mirrorConfig) > 0)
	if err != nil {
		return "", err
	}

	cacheDir, err := GetCacheDir(imageDataType)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf(coreOsFileName, architecture)
	// Check if file is already cached
	filePath, err := GetFileFromCache(filename, cacheDir)
	if err != nil {
		return "", err
	}
	if filePath != "" {
		// Found cached file
		return filePath, nil
	}

	path, err := r.extractFileFromImage(log, image, filename, cacheDir, pullSecret, mirrorConfig)
	if err != nil {
		return "", err
	}
	return path, err
}

func (r *release) getImageFromRelease(log logrus.FieldLogger, imageName, releaseImage, pullSecret string, haveMirror bool) (string, error) {
	// This requires the 'oc' command so make sure its available
	_, err := exec.LookPath("oc")
	if err != nil {
		if haveMirror {
			log.Warning("Unable to validate mirror config because \"oc\" command is not available")
		} else {
			log.Debug("Skipping ISO extraction; \"oc\" command is not available")
		}
		return "", err
	}

	cmd := fmt.Sprintf(templateGetImage, imageName, true, releaseImage)

	log.Debugf("Fetching image from OCP release (%s)", cmd)
	image, err := execute(log, r.executer, pullSecret, cmd)
	if err != nil {
		return "", err
	}

	return image, nil
}

func (r *release) extractFileFromImage(log logrus.FieldLogger, image, file, cacheDir, pullSecret string, mirrorConfig []mirror.RegistriesConfig) (string, error) {

	var cmd string
	if len(mirrorConfig) > 0 {
		log.Debugf("Using mirror configuration")
		icspFile, err := getIcspFileFromRegistriesConfig(log, mirrorConfig)
		if err != nil {
			return "", err
		}
		defer removeIcspFile(icspFile)
		cmd = fmt.Sprintf(templateImageExtractWithIcsp, file, cacheDir, icspFile, image)
	} else {
		cmd = fmt.Sprintf(templateImageExtract, file, cacheDir, image)
	}

	log.Debugf("extracting %s to %s, %s", file, cacheDir, cmd)
	_, err := retry.Do(r.config.MaxTries, r.config.RetryDelay, execute, log, r.executer, pullSecret, cmd)
	if err != nil {
		return "", err
	}
	// set path
	path := filepath.Join(cacheDir, file)
	log.Info("Successfully extracted base ISO from the release")
	log.Debugf("Base ISO %s cached at %s", file, path)
	return path, nil
}

func execute(log logrus.FieldLogger, executer executer.Executer, pullSecret string, command string) (string, error) {

	ps, err := executer.TempFile("", "registry-config")
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

	stdout, stderr, exitCode := executer.Execute(args[0], args[1:]...)

	if exitCode == 0 {
		return strings.TrimSpace(stdout), nil
	}

	err = fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s", executeCommand, exitCode, stdout, stderr)
	log.Error(err)
	return "", err
}

// Create a temporary file containing the ImageContentPolicySources
func getIcspFileFromRegistriesConfig(log logrus.FieldLogger, mirrorConfig []mirror.RegistriesConfig) (string, error) {

	contents, err := getIcspContents(mirrorConfig)
	if err != nil {
		return "", err
	}
	if contents == nil {
		log.Debugf("No registry entries to build ICSP file")
		return "", nil
	}

	icspFile, err := ioutil.TempFile("", "icsp-file")
	if err != nil {
		return "", err
	}
	log.Debugf("Building ICSP file from registries.conf with contents %s", contents)
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
