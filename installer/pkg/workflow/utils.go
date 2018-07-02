package workflow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
	configgenerator "github.com/coreos/tectonic-installer/installer/pkg/config-generator"
)

const (
	assetsStep                 = "assets"
	binaryPrefix               = "installer"
	bootstrapOff               = "-var=tectonic_bootstrap=false"
	bootstrapOn                = "-var=tectonic_bootstrap=true"
	configFileName             = "config.yaml"
	etcdStep                   = "etcd"
	generatedPath              = "generated"
	internalFileName           = "internal.yaml"
	joinWorkersStep            = "joining_workers"
	kcoConfigFileName          = "kco-config.yaml"
	kubeSystemFileName         = "cluster-config.yaml"
	kubeSystemPath             = "generated/manifests"
	mastersStep                = "masters"
	newTLSPath                 = "generated/newTLS"
	newTLSStep                 = "newtls"
	stepsBaseDir               = "steps"
	tectonicSystemFileName     = "cluster-config.yaml"
	tectonicSystemPath         = "generated/tectonic"
	terraformVariablesFileName = "terraform.tfvars"
	tlsStep                    = "tls"
	tncDNSStep                 = "tnc_dns"
	tncoConfigFileName         = "tnco-config.yaml"
	topologyStep               = "topology"
)

func buildInternalConfig(clusterDir string) error {
	if clusterDir == "" {
		return errors.New("no cluster dir given for building internal config")
	}

	// fill the internal struct
	clusterID, err := configgenerator.GenerateClusterID(16)
	if err != nil {
		return err
	}
	internalCfg := config.Internal{
		ClusterID: clusterID,
	}

	// store the content
	yamlContent, err := yaml.Marshal(internalCfg)
	internalFileContent := []byte("# Do not touch, auto-generated\n")
	internalFileContent = append(internalFileContent, yamlContent...)
	if err != nil {
		return err
	}
	return writeFile(filepath.Join(clusterDir, internalFileName), string(internalFileContent))
}

func copyFile(fromFilePath, toFilePath string) error {
	from, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

// returns the directory containing templates for a given step. If platform is
// specified, it looks for a subdirectory with platform first, falling back if
// there are no platform-specific templates for that step
func findStepTemplates(stepName string, platform config.Platform) (string, error) {
	base, err := baseLocation()
	if err != nil {
		return "", fmt.Errorf("error looking up step %s templates: %v", stepName, err)
	}
	for _, path := range []string{
		filepath.Join(base, stepsBaseDir, stepName, platformPath(platform)),
		filepath.Join(base, stepsBaseDir, stepName)} {

		stat, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("invalid path for '%s' templates: %s", base, err)
		}
		if !stat.IsDir() {
			return "", fmt.Errorf("invalid path for '%s' templates", base)
		}
		return path, nil
	}
	return "", os.ErrNotExist
}

func platformPath(platform config.Platform) string {
	switch platform {
	case config.PlatformLibvirt:
		return "libvirt"
	case config.PlatformAWS:
		return "aws"
	}
	panic("invalid platform")
}

func generateClusterConfigMaps(m *metadata) error {
	clusterGeneratedPath := filepath.Join(m.clusterDir, generatedPath)
	if err := os.MkdirAll(clusterGeneratedPath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Failed to create cluster generated directory at %s", clusterGeneratedPath)
	}

	configGenerator := configgenerator.New(m.cluster)

	kcoConfig, err := configGenerator.CoreConfig()
	if err != nil {
		return err
	}

	kcoConfigFilePath := filepath.Join(clusterGeneratedPath, kcoConfigFileName)
	if err := writeFile(kcoConfigFilePath, kcoConfig); err != nil {
		return err
	}

	tncoConfig, err := configGenerator.TncoConfig()
	if err != nil {
		return err
	}

	tncoConfigFilePath := filepath.Join(clusterGeneratedPath, tncoConfigFileName)
	if err := writeFile(tncoConfigFilePath, tncoConfig); err != nil {
		return err
	}

	kubeSystem, err := configGenerator.KubeSystem()
	if err != nil {
		return err
	}

	kubePath := filepath.Join(m.clusterDir, kubeSystemPath)
	if err := os.MkdirAll(kubePath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Failed to create manifests directory at %s", kubePath)
	}

	kubeSystemConfigFilePath := filepath.Join(kubePath, kubeSystemFileName)
	if err := writeFile(kubeSystemConfigFilePath, kubeSystem); err != nil {
		return err
	}

	tectonicSystem, err := configGenerator.TectonicSystem()
	if err != nil {
		return err
	}

	tectonicPath := filepath.Join(m.clusterDir, tectonicSystemPath)
	if err := os.MkdirAll(tectonicPath, os.ModeDir|0755); err != nil {
		return fmt.Errorf("Failed to create tectonic directory at %s", tectonicPath)
	}

	tectonicSystemConfigFilePath := filepath.Join(tectonicPath, tectonicSystemFileName)
	return writeFile(tectonicSystemConfigFilePath, tectonicSystem)
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

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := fmt.Fprintln(w, content); err != nil {
		return err
	}
	w.Flush()

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

func clusterIsBootstrapped(stateDir string) bool {
	return hasStateFile(stateDir, topologyStep) &&
		hasStateFile(stateDir, mastersStep) &&
		hasStateFile(stateDir, tncDNSStep)
}

func createTNCCNAME(m *metadata) error {
	return runInstallStep(m, tncDNSStep, []string{bootstrapOn}...)
}

func createTNCARecord(m *metadata) error {
	return runInstallStep(m, tncDNSStep, []string{bootstrapOff}...)
}

func destroyTNCDNS(m *metadata) error {
	return runDestroyStep(m, tncDNSStep, []string{bootstrapOff}...)
}

func generateTerraformVariablesStep(m *metadata) error {
	vars, err := m.cluster.TFVars()
	if err != nil {
		return err
	}

	terraformVariablesFilePath := filepath.Join(m.clusterDir, terraformVariablesFileName)
	return writeFile(terraformVariablesFilePath, vars)
}
