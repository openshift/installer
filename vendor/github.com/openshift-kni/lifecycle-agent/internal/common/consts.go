/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"math"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Common constants mainly used by packages in lca-cli
const (
	VarFolder       = "/var"
	BackupDir       = "/var/tmp/backup"
	BackupCertsDir  = "/var/tmp/backupCertsDir"
	BackupChecksDir = "/var/tmp/checks"

	// Workload partitioning annotation key and value
	WorkloadManagementAnnotationKey   = "target.workload.openshift.io/management"
	WorkloadManagementAnnotationValue = `{"effect": "PreferredDuringScheduling"}`

	// ImageRegistryAuthFile is the pull secret. Written by the machine-config-operator
	ImageRegistryAuthFile = "/var/lib/kubelet/config.json"
	KubeconfigFile        = "/etc/kubernetes/static-pod-resources/kube-apiserver-certs/secrets/node-kubeconfigs/lb-ext.kubeconfig"

	RecertImageEnvKey      = "RELATED_IMAGE_RECERT_IMAGE"
	DefaultRecertImage     = "quay.io/edge-infrastructure/recert:v0"
	EtcdStaticPodFile      = "/etc/kubernetes/manifests/etcd-pod.yaml"
	EtcdStaticPodContainer = "etcd"
	EtcdDefaultEndpoint    = "localhost:2379"

	OvnIcEtcFolder = "/var/lib/ovn-ic/etc"
	OvnNodeCerts   = OvnIcEtcFolder + "/ovnkube-node-certs"

	MultusCerts = "/etc/cni/multus/certs"

	MCDCurrentConfig = "/etc/machine-config-daemon/currentconfig"

	InstallationConfigurationFilesDir = "/usr/local/installation_configuration_files"
	OptOpenshift                      = "/opt/openshift"
	SeedDataDir                       = "/var/seed_data"
	KubeconfigCryptoDir               = "kubeconfig-crypto"
	ClusterConfigDir                  = "cluster-configuration"
	SeedClusterInfoFileName           = "manifest.json"
	SeedReconfigurationFileName       = "manifest.json"
	ManifestsDir                      = "manifests"
	ExtraManifestsDir                 = "extra-manifests"
	EtcdContainerName                 = "recert_etcd"
	LvmConfigDir                      = "lvm-configuration"
	LvmDevicesPath                    = "/etc/lvm/devices/system.devices"
	CABundleFilePath                  = "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"

	LCAConfigDir                                    = "/var/lib/lca"
	IBUAutoRollbackConfigFile                       = LCAConfigDir + "/autorollback_config.json"
	IBUAutoRollbackInitMonitorTimeoutDefaultSeconds = 1800
	IBUInitMonitorService                           = "lca-init-monitor.service"
	IBUInitMonitorServiceFile                       = "/etc/systemd/system/" + IBUInitMonitorService
	// AutoRollbackOnFailurePostRebootConfigAnnotation configure automatic rollback when the reconfiguration of the cluster fails upon the first reboot.
	// Only acceptable value is AutoRollbackDisableValue. Any other value is treated as "Enabled".
	AutoRollbackOnFailurePostRebootConfigAnnotation = "auto-rollback-on-failure.lca.openshift.io/post-reboot-config"
	// AutoRollbackOnFailureUpgradeCompletionAnnotation configure automatic rollback after the Lifecycle Agent reports a failed upgrade upon completion.
	// Only acceptable value is AutoRollbackOnFailureDisableValue. Any other value is treated as "Enabled".
	AutoRollbackOnFailureUpgradeCompletionAnnotation = "auto-rollback-on-failure.lca.openshift.io/upgrade-completion"
	// AutoRollbackOnFailureInitMonitorAnnotation configure automatic rollback LCA Init Monitor watchdog, which triggers auto-rollback if timeout occurs before upgrade completion
	// Only acceptable value is AutoRollbackDisableValue. Any other value is treated as "Enabled".
	AutoRollbackOnFailureInitMonitorAnnotation = "auto-rollback-on-failure.lca.openshift.io/init-monitor"
	// AutoRollbackDisableValue value that decides if rollback is disabled
	AutoRollbackDisableValue = "Disabled"

	LcaNamespace = "openshift-lifecycle-agent"
	Host         = "/host"

	CsvDeploymentName      = "cluster-version-operator"
	CsvDeploymentNamespace = "openshift-cluster-version"
	// InstallConfigCM cm name
	InstallConfigCM = "cluster-config-v1"
	// InstallConfigCMNamespace cm namespace
	InstallConfigCMNamespace = "kube-system"
	// InstallConfigCMNamespace data key
	InstallConfigCMInstallConfigDataKey = "install-config"
	OpenshiftInfraCRName                = "cluster"
	OpenshiftProxyCRName                = "cluster"

	// Env var to configure auto rollback for post-reboot config failure
	IBUPostRebootConfigAutoRollbackOnFailureEnv = "LCA_IBU_AUTO_ROLLBACK_ON_CONFIG_FAILURE"

	// Bump this every time the seed format changes in a backwards incompatible way
	SeedFormatVersion  = 3
	SeedFormatOCILabel = "com.openshift.lifecycle-agent.seed_format_version"

	SeedClusterInfoOCILabel = "com.openshift.lifecycle-agent.seed_cluster_info"

	PullSecretName           = "pull-secret"
	PullSecretEmptyData      = "{\"auths\":{\"registry.connect.redhat.com\":{\"username\":\"empty\",\"password\":\"empty\",\"auth\":\"ZW1wdHk6ZW1wdHk=\",\"email\":\"\"}}}" //nolint:gosec
	OpenshiftConfigNamespace = "openshift-config"

	NMConnectionFolder = "/etc/NetworkManager/system-connections"
	NetworkDir         = "network-configuration"
)

// Annotation names and values related to extra manifest
const (
	ApplyWaveAnn        = "lca.openshift.io/apply-wave"
	defaultApplyWave    = math.MaxInt32 // 2147483647, an enough large number
	ApplyTypeAnnotation = "lca.openshift.io/apply-type"
	ApplyTypeReplace    = "replace" // default if annotation doesn't exist
	ApplyTypeMerge      = "merge"
)

var (
	BackupGvk  = schema.GroupVersionKind{Group: "velero.io", Kind: "Backup", Version: "v1"}
	RestoreGvk = schema.GroupVersionKind{Group: "velero.io", Kind: "Restore", Version: "v1"}
)

// CertPrefixes is the list of certificate prefixes to be backed up
// before creating the seed image
var CertPrefixes = []string{
	"loadbalancer-serving-signer",
	"localhost-serving-signer",
	"service-network-serving-signer",
}

var TarOpts = []string{"--selinux", "--xattrs", "--xattrs-include=*", "--acls"}
