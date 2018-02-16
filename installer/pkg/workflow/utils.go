package workflow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/tectonic-installer/installer/pkg/config"
	"github.com/coreos/tectonic-installer/installer/pkg/config-generator"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	assetsStep     = "assets"
	bootstrapStep  = "bootstrap"
	configFileName = "config.yaml"
	joinStep       = "joining"
	kubeConfigPath = "generated/auth/kubeconfig"
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
	return terraformExec(clusterDir, "destroy", "-force", fmt.Sprintf("-state=%s.tfstate", bootstrapStep), "-target=aws_route53_record.tectonic_ncg", findTemplatesForStep(bootstrapStep))
}

// TODO: Handle errors for Getwd
func findTemplatesForStep(step string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "steps", step)
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
	if err := terraformExec(m.clusterDir, "import", fmt.Sprintf("-state=%s.tfstate", joinStep), fmt.Sprintf("-config=%s", findTemplatesForStep(joinStep)), "aws_autoscaling_group.masters", fmt.Sprintf("%s-masters", m.cluster.Name)); err != nil {
		return err
	}

	return terraformExec(m.clusterDir, "import", fmt.Sprintf("-state=%s.tfstate", joinStep), fmt.Sprintf("-config=%s", findTemplatesForStep(joinStep)), "aws_autoscaling_group.workers", fmt.Sprintf("%s-workers", m.cluster.Name))

}

func readClusterConfig(configFilePath string) (*config.Cluster, error) {
	config, err := config.ParseFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid config file: %s", configFilePath, err)
	}

	return &config.Clusters[0], nil
}

func readClusterConfigStep(m *metadata) error {
	var configFilePath string

	if m.configFilePath != "" {
		configFilePath = m.configFilePath
	} else {
		configFilePath = filepath.Join(m.clusterDir, configFileName)
	}

	cluster, err := readClusterConfig(configFilePath)
	if err != nil {
		return err
	}

	m.cluster = *cluster

	return nil
}

func waitForNCG(m *metadata) error {
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
		ds, _ := client.DaemonSets("kube-system").Get("ncg")
		log.Printf("Waiting for NCG to be running, this might take a while...")
		if ds.Status.NumberReady >= 1 {
			return nil
		}
		time.Sleep(time.Second * time.Duration(wait))
		retries--
	}

	return errors.New("NCG is not running")
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
