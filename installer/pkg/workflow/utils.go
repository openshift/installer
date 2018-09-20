package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	configgenerator "github.com/openshift/installer/installer/pkg/config-generator"
	"github.com/openshift/installer/pkg/types/config"
)

const (
	binaryPrefix     = "installer"
	configFileName   = "config.yaml"
	internalFileName = "internal.yaml"
)

func generateClusterConfigMaps(m *metadata) error {
	clusterGeneratedPath := filepath.Join(m.clusterDir, generatedPath)
	if err := os.MkdirAll(clusterGeneratedPath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create cluster generated directory at %s", clusterGeneratedPath)
	}

	configGenerator := configgenerator.New(m.cluster)

	kcoConfig, err := configGenerator.CoreConfig()
	if err != nil {
		return err
	}

	kcoConfigFilePath := filepath.Join(clusterGeneratedPath, kcoConfigFileName)
	if err := ioutil.WriteFile(kcoConfigFilePath, []byte(kcoConfig), 0666); err != nil {
		return err
	}

	kubeSystem, err := configGenerator.KubeSystem(m.clusterDir)
	if err != nil {
		return err
	}

	kubePath := filepath.Join(m.clusterDir, kubeSystemPath)
	if err := os.MkdirAll(kubePath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create manifests directory at %s", kubePath)
	}

	kubeSystemConfigFilePath := filepath.Join(kubePath, kubeSystemFileName)
	if err := ioutil.WriteFile(kubeSystemConfigFilePath, []byte(kubeSystem), 0666); err != nil {
		return err
	}

	tectonicSystem, err := configGenerator.TectonicSystem()
	if err != nil {
		return err
	}

	tectonicPath := filepath.Join(m.clusterDir, tectonicSystemPath)
	if err := os.MkdirAll(tectonicPath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("failed to create tectonic directory at %s", tectonicPath)
	}

	tectonicSystemConfigFilePath := filepath.Join(tectonicPath, tectonicSystemFileName)
	return ioutil.WriteFile(tectonicSystemConfigFilePath, []byte(tectonicSystem), 0666)
}

func readClusterConfig(configFilePath string, internalFilePath string) (*config.Cluster, error) {
	cfg, err := config.ParseConfigFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid config file: %s", configFilePath, err)
	}

	if internalFilePath != "" {
		internal, err := config.ParseInternalFile(internalFilePath)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid internal file: %s", internalFilePath, err)
		}
		cfg.Internal = *internal
	}

	return cfg, nil
}

func readClusterConfigStep(m *metadata) error {
	if m.clusterDir == "" {
		return errors.New("no cluster dir given for reading config")
	}
	configFilePath := filepath.Join(m.clusterDir, configFileName)
	internalFilePath := filepath.Join(m.clusterDir, internalFileName)

	cluster, err := readClusterConfig(configFilePath, internalFilePath)
	if err != nil {
		return err
	}

	if err := cluster.ValidateAndLog(); err != nil {
		return err
	}

	m.cluster = *cluster

	return nil
}

func baseLocation() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("undetermined location of own executable: %s", err)
	}
	ex = path.Dir(ex)
	if path.Base(ex) != binaryPrefix {
		return "", fmt.Errorf("%s executable in unknown location: %s", path.Base(ex), err)
	}
	return path.Dir(ex), nil
}
