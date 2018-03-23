package workflow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
	"github.com/coreos/tectonic-installer/installer/pkg/config-generator"

	log "github.com/Sirupsen/logrus"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	stepsBaseDir     = "steps"
	assetsStep       = "assets"
	bootstrapStep    = "bootstrap"
	joinStep         = "joining"
	configFileName   = "config.yaml"
	internalFileName = "internal.yaml"
	kubeConfigPath   = "generated/auth/kubeconfig"
	binaryPrefix     = "tectonic-installer"
	tncDaemonSet     = "tectonic-node-controller"
)

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

func destroyCNAME(clusterDir string) error {
	templatesPath, err := findTemplates(bootstrapStep)
	if err != nil {
		return err
	}
	return terraformExec(clusterDir, "destroy", "-force", fmt.Sprintf("-state=%s.tfstate", bootstrapStep), "-target=aws_route53_record.tectonic_tnc", templatesPath)
}

func findTemplates(relativePath string) (string, error) {
	base, err := baseLocation()
	if err != nil {
		return "", fmt.Errorf("unknown base path for '%s' templates: %s", relativePath, err)
	}
	base = filepath.Join(base, stepsBaseDir, relativePath)
	stat, err := os.Stat(base)
	if err != nil {
		return "", fmt.Errorf("invalid path for '%s' templates: %s", base, err)
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("invalid path for '%s' templates", base)
	}
	return base, nil
}

func generateClusterConfigStep(m *metadata) error {
	configGenerator := configgenerator.New(m.cluster)

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

func importAutoScalingGroup(m *metadata) error {
	templatesPath, err := findTemplates(joinStep)
	if err != nil {
		return err
	}
	err = terraformExec(
		m.clusterDir,
		"import",
		fmt.Sprintf("-state=%s.tfstate", joinStep),
		fmt.Sprintf("-config=%s", templatesPath),
		"aws_autoscaling_group.masters",
		fmt.Sprintf("%s-masters", m.cluster.Name))
	if err != nil {
		return err
	}
	err = terraformExec(
		m.clusterDir,
		"import",
		fmt.Sprintf("-state=%s.tfstate", joinStep),
		fmt.Sprintf("-config=%s", templatesPath),
		"aws_autoscaling_group.workers",
		fmt.Sprintf("%s-workers", m.cluster.Name))
	return err
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
	var configFilePath string
	var internalFilePath string

	if m.configFilePath != "" {
		configFilePath = m.configFilePath
	} else {
		configFilePath = filepath.Join(m.clusterDir, configFileName)
		internalFilePath = filepath.Join(m.clusterDir, internalFileName)
	}

	cluster, err := readClusterConfig(configFilePath, internalFilePath)
	if err != nil {
		return err
	}

	if errs := cluster.Validate(); len(errs) != 0 {
		log.Errorf("Found %d errors in the cluster definition:", len(errs))
		for i, err := range errs {
			log.Errorf("error %d: %v", i+1, err)
		}
		return fmt.Errorf("found %d cluster definition errors", len(errs))
	}

	m.cluster = *cluster

	return nil
}

func waitForTNC(m *metadata) error {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(m.clusterDir, kubeConfigPath))
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	retries := 180
	wait := 10
	for retries > 0 {
		// client will error until api sever is up
		ds, _ := client.AppsV1().DaemonSets("kube-system").Get(tncDaemonSet, meta.GetOptions{})
		log.Info("Waiting for TNC to be running, this might take a while...")
		if ds.Status.NumberReady >= 1 {
			return nil
		}
		time.Sleep(time.Second * time.Duration(wait))
		retries--
	}

	return errors.New("TNC is not running")
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
	if path.Base(ex) != runtime.GOOS {
		return "", fmt.Errorf("%s executable in unknown location: %s", path.Base(ex), err)
	}
	ex = path.Dir(ex)
	if path.Base(ex) != binaryPrefix {
		return "", fmt.Errorf("%s executable in unknown location: %s", path.Base(ex), err)
	}
	return path.Dir(ex), nil
}
