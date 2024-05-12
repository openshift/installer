package recert

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift-kni/lifecycle-agent/api/seedreconfig"
	"github.com/openshift-kni/lifecycle-agent/internal/common"
	"github.com/openshift-kni/lifecycle-agent/lca-cli/seedclusterinfo"
	"github.com/openshift-kni/lifecycle-agent/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	RecertConfigFile = "recert_config.json"
	SummaryFile      = "/var/tmp/recert-summary.yaml"
)

var staticDirs = []string{"/kubelet", "/kubernetes", "/machine-config-daemon"}

type RecertConfig struct {
	DryRun           bool   `json:"dry_run,omitempty"`
	ExtendExpiration bool   `json:"extend_expiration,omitempty"`
	ForceExpire      bool   `json:"force_expire,omitempty"`
	EtcdEndpoint     string `json:"etcd_endpoint,omitempty"`
	ClusterRename    string `json:"cluster_rename,omitempty"`
	Hostname         string `json:"hostname,omitempty"`
	IP               string `json:"ip,omitempty"`
	Proxy            string `json:"proxy,omitempty"`
	InstallConfig    string `json:"install_config,omitempty"`
	// We intentionally don't omitEmpty this field because an empty string here
	// means "delete the kubeadmin password secret" while a complete omission
	// of the field means "don't touch the secret". We never want the latter,
	// we either want to delete the secret or update it, never leave it as is.
	KubeadminPasswordHash string `json:"kubeadmin_password_hash"`
	// WARNING: You probably don't want use `SummaryFile`! This will leak
	// private keys and tokens!
	SummaryFile       string   `json:"summary_file,omitempty"`
	SummaryFileClean  string   `json:"summary_file_clean,omitempty"`
	StaticDirs        []string `json:"static_dirs,omitempty"`
	StaticFiles       []string `json:"static_files,omitempty"`
	CNSanReplaceRules []string `json:"cn_san_replace_rules,omitempty"`
	UseKeyRules       []string `json:"use_key_rules,omitempty"`
	UseCertRules      []string `json:"use_cert_rules,omitempty"`
	PullSecret        string   `json:"pull_secret,omitempty"`
}

func FormatRecertProxyFromSeedReconfigProxy(proxy, statusProxy *seedreconfig.Proxy) string {
	if proxy == nil || statusProxy == nil {
		// Both must be set, anything else is invalid
		return ""
	}
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		proxy.HTTPProxy, proxy.HTTPSProxy, proxy.NoProxy,
		statusProxy.HTTPProxy, statusProxy.HTTPSProxy, statusProxy.NoProxy,
	)
}

// CreateRecertConfigFile function to create recert config file
// those params will be provided to an installation script after reboot
// that will run recert command with them
func CreateRecertConfigFile(seedReconfig *seedreconfig.SeedReconfiguration, seedClusterInfo *seedclusterinfo.SeedClusterInfo, cryptoDir, recertConfigFolder string) error {
	config := createBasicEmptyRecertConfig()

	config.ClusterRename = fmt.Sprintf("%s:%s", seedReconfig.ClusterName, seedReconfig.BaseDomain)
	if seedReconfig.InfraID != "" {
		config.ClusterRename = fmt.Sprintf("%s:%s", config.ClusterRename, seedReconfig.InfraID)
	}

	if seedReconfig.Hostname != seedClusterInfo.SNOHostname {
		config.Hostname = seedReconfig.Hostname
	}

	if seedReconfig.NodeIP != seedClusterInfo.NodeIP {
		config.IP = seedReconfig.NodeIP
	}

	config.Proxy = FormatRecertProxyFromSeedReconfigProxy(seedReconfig.Proxy, seedReconfig.StatusProxy)

	config.InstallConfig = seedReconfig.InstallConfig

	config.SummaryFileClean = SummaryFile
	seedFullDomain := fmt.Sprintf("%s.%s", seedClusterInfo.ClusterName, seedClusterInfo.BaseDomain)
	clusterFullDomain := fmt.Sprintf("%s.%s", seedReconfig.ClusterName, seedReconfig.BaseDomain)
	config.ExtendExpiration = true
	config.CNSanReplaceRules = []string{
		fmt.Sprintf("system:node:%s,system:node:%s", seedClusterInfo.SNOHostname, seedReconfig.Hostname),
		fmt.Sprintf("%s,%s", seedClusterInfo.SNOHostname, seedReconfig.Hostname),
		fmt.Sprintf("%s,%s", seedClusterInfo.NodeIP, seedReconfig.NodeIP),
		fmt.Sprintf("api.%s,api.%s", seedFullDomain, clusterFullDomain),
		fmt.Sprintf("api-int.%s,api-int.%s", seedFullDomain, clusterFullDomain),
		fmt.Sprintf("*.apps.%s,*.apps.%s", seedFullDomain, clusterFullDomain),
	}
	config.KubeadminPasswordHash = seedReconfig.KubeadminPasswordHash
	config.PullSecret = seedReconfig.PullSecret

	if _, err := os.Stat(cryptoDir); err == nil {
		ingressFile, ingressCN, err := getIngressCNAndFile(cryptoDir)
		if err != nil {
			return err
		}
		config.UseKeyRules = []string{
			fmt.Sprintf("kube-apiserver-lb-signer %s/loadbalancer-serving-signer.key", cryptoDir),
			fmt.Sprintf("kube-apiserver-localhost-signer %s/localhost-serving-signer.key", cryptoDir),
			fmt.Sprintf("kube-apiserver-service-network-signer %s/service-network-serving-signer.key", cryptoDir),
			fmt.Sprintf("%s %s/%s", ingressCN, cryptoDir, ingressFile),
		}
		config.UseCertRules = []string{filepath.Join(cryptoDir, "admin-kubeconfig-client-ca.crt")}
	}

	p := filepath.Join(recertConfigFolder, RecertConfigFile)
	if err := utils.MarshalToFile(config, p); err != nil {
		return fmt.Errorf("failed to marshal recert config file to %s: %w", p, err)
	}

	return nil
}

func CreateRecertConfigFileForSeedCreation(path string, withPassword bool) error {
	config := createBasicEmptyRecertConfig()
	config.SummaryFileClean = "/kubernetes/recert-seed-creation-summary.yaml"
	config.ForceExpire = true

	config.KubeadminPasswordHash = ""
	if withPassword {
		bytes, err := generateDisposablePasswordHash()
		if err != nil {
			return fmt.Errorf("failed to generate password hash: %w", err)
		}
		config.KubeadminPasswordHash = string(bytes)
	}

	if err := utils.MarshalToFile(config, path); err != nil {
		return fmt.Errorf("failed create recert config file for sed creatation in %s: %w", path, err)
	}

	return nil
}

// generateDisposablePasswordHash generates a random password hash from a ridiculously
// long length password that is never meant to be known or used by anyone, but
// only to be used as a placeholder in the seed. It will be replaced or deleted
// during seed reconfiguration.
func generateDisposablePasswordHash() ([]byte, error) {
	bcryptLargestSupportedLength := 72
	token := make([]byte, bcryptLargestSupportedLength)
	_, err := rand.Read(token)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes for password: %w", err)
	}
	bytes, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}
	return bytes, nil
}

func CreateRecertConfigFileForSeedRestoration(path, originalPasswordHash string) error {
	config := createBasicEmptyRecertConfig()
	config.SummaryFileClean = "/kubernetes/recert-seed-restoration-summary.yaml"
	config.ExtendExpiration = true
	config.UseKeyRules = []string{
		fmt.Sprintf("kube-apiserver-lb-signer %s/loadbalancer-serving-signer.key", common.BackupCertsDir),
		fmt.Sprintf("kube-apiserver-localhost-signer %s/localhost-serving-signer.key", common.BackupCertsDir),
		fmt.Sprintf("kube-apiserver-service-network-signer %s/service-network-serving-signer.key", common.BackupCertsDir),
		fmt.Sprintf("ingresskey-ingress-operator %s/ingresskey-ingress-operator.key", common.BackupCertsDir),
	}
	config.UseCertRules = []string{filepath.Join(common.BackupCertsDir, "admin-kubeconfig-client-ca.crt")}
	config.KubeadminPasswordHash = originalPasswordHash

	if err := utils.MarshalToFile(config, path); err != nil {
		return fmt.Errorf("failed to marshall recert config file for seed restoration: %w", err)
	}
	return nil
}

func createBasicEmptyRecertConfig() RecertConfig {
	return RecertConfig{
		DryRun:       false,
		EtcdEndpoint: common.EtcdDefaultEndpoint,
		StaticDirs:   staticDirs,
		StaticFiles: []string{
			"/host-etc/mcs-machine-config-content.json",
			"/host-etc/mco/proxy.env",
		},
	}
}

func getIngressCNAndFile(certsFolder string) (string, string, error) {
	certsFiles, err := os.ReadDir(certsFolder)
	if err != nil {
		return "", "", fmt.Errorf("failed to list files in %s while searching for ingress cn, "+
			"err: %w", certsFolder, err)
	}

	for _, path := range certsFiles {
		if strings.HasPrefix(path.Name(), "ingresskey-") {
			return path.Name(), strings.Replace(path.Name(), "ingresskey-", "", 1), nil
		}
	}
	return "", "", fmt.Errorf("failed to find ingress key file")
}
