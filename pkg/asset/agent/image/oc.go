package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openshift/assisted-service/models"
	"github.com/openshift/assisted-service/pkg/executer"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
)

const (
	MachineOsImageName = "machine-os-images"
	CoreOsFileName     = "/coreos/coreos-x86_64.iso"
	OcDefaultTries     = 5
	OcDefaltRetryDelay = time.Second * 5
)

type Config struct {
	MaxTries   uint
	RetryDelay time.Duration
}

type Release interface {
	GetBaseIso(log logrus.FieldLogger, releaseImage string, pullSecret string, platformType models.PlatformType) (string, error)
}

type release struct {
	executer executer.Executer
	config   Config
}

func NewRelease(executer executer.Executer, config Config) Release {
	return &release{executer: executer, config: config}
}

const (
	templateGetImage     = "oc adm release info --image-for=%s --insecure=%t %s"
	templateImageExtract = "oc image extract --file=%s --path /:/%s --confirm %s"
)

// Get the CoreOS ISO from the releaseImage
func (r *release) GetBaseIso(log logrus.FieldLogger, releaseImage string, pullSecret string, platformType models.PlatformType) (string, error) {

	// Get the machine-os-images pullspec from the release and use that to get the CoreOS ISO
	image, err := r.getImageFromRelease(log, MachineOsImageName, releaseImage, pullSecret, true)
	if err != nil {
		return "", err
	}

	cacheDir, err := GetCacheDir(imageDataType)
	if err != nil {
		return "", err
	}

	// Check if file is already cached
	filePath, err := GetFileFromCache(CoreOsFileName, cacheDir)
	if err != nil {
		return "", err
	}
	if filePath != "" {
		// Found cached file
		return filePath, nil
	}

	path, err := r.extractFileFromImage(log, image, cacheDir, pullSecret, true, platformType)
	if err != nil {
		return "", err
	}
	return path, err
}

func (r *release) getImageFromRelease(log logrus.FieldLogger, imageName, releaseImage, pullSecret string, insecure bool) (string, error) {
	cmd := fmt.Sprintf(templateGetImage, imageName, insecure, releaseImage)

	log.Infof("Fetching image from OCP release (%s)", cmd)
	image, err := execute(log, r.executer, pullSecret, cmd)
	if err != nil {
		return "", err
	}

	return image, nil
}

func (r *release) extractFileFromImage(log logrus.FieldLogger, image, cacheDir, pullSecret string, insecure bool, platformType models.PlatformType) (string, error) {

	file := CoreOsFileName
	cmd := fmt.Sprintf(templateImageExtract, file, cacheDir, image)
	log.Infof("extracting %s to %s, %s", file, cacheDir, cmd)
	_, err := retry.Do(r.config.MaxTries, r.config.RetryDelay, execute, log, r.executer, pullSecret, cmd)
	if err != nil {
		return "", err
	}
	// set path
	path := filepath.Join(cacheDir, file)
	log.Infof("Successfully extracted %s binary from the release to: %s", file, path)
	return path, nil
}

func execute(log logrus.FieldLogger, executer executer.Executer, pullSecret string, command string) (string, error) {

	// TODO add in icsp file
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
	} else {
		err = fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s", executeCommand, exitCode, stdout, stderr)
		log.Error(err)
		return "", err
	}
}
