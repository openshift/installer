package utils

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/openshift-kni/lifecycle-agent/api/seedreconfig"
	"github.com/openshift-kni/lifecycle-agent/internal/common"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	cryptoDirMode = 0o600
	pwHashFile    = "kubeadmin-password-hash.txt"
)

func SeedReconfigurationKubeconfigRetentionToCryptoDir(cryptoDir string, kubeconfigCryptoRetention *seedreconfig.KubeConfigCryptoRetention) error {
	if err := os.MkdirAll(cryptoDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating %s: %w", cryptoDir, err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, "admin-kubeconfig-client-ca.crt"),
		[]byte(kubeconfigCryptoRetention.KubeAPICrypto.ClientAuthCrypto.AdminCACertificate), cryptoDirMode); err != nil {
		return fmt.Errorf("error writing admin-kubeconfig-client-ca.crt: %w", err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, "loadbalancer-serving-signer.key"),
		[]byte(kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.LoadbalancerSignerPrivateKey), cryptoDirMode); err != nil {
		return fmt.Errorf("error writing loadbalancer-serving-signer.key: %w", err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, "localhost-serving-signer.key"),
		[]byte(kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.LocalhostSignerPrivateKey), cryptoDirMode); err != nil {
		return fmt.Errorf("error writing localhost-serving-signer.key: %w", err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, "service-network-serving-signer.key"),
		[]byte(kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.ServiceNetworkSignerPrivateKey), cryptoDirMode); err != nil {
		return fmt.Errorf("error writing service-network-serving-signer.key: %w", err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, "ingresskey-ingress-operator.key"),
		[]byte(kubeconfigCryptoRetention.IngresssCrypto.IngressCA), cryptoDirMode); err != nil {
		return fmt.Errorf("error writing ingresskey-ingress-operator.key: %w", err)
	}

	return nil
}

func SeedReconfigurationKubeconfigRetentionFromCluster(ctx context.Context, client runtimeclient.Client) (*seedreconfig.KubeConfigCryptoRetention, error) {
	var kubeconfigCryptoRetention seedreconfig.KubeConfigCryptoRetention

	adminKubeConfigClientCA, err := GetConfigMapData(ctx, "admin-kubeconfig-client-ca", "openshift-config", "ca-bundle.crt", client)
	if err != nil {
		return nil, err
	}

	loadbalancerServingSignerPrivateKey, err := GetSecretData(ctx, "loadbalancer-serving-signer", "openshift-kube-apiserver-operator", "tls.key", client)
	if err != nil {
		return nil, err
	}

	localhostServingSignerKey, err := GetSecretData(ctx, "localhost-serving-signer", "openshift-kube-apiserver-operator", "tls.key", client)
	if err != nil {
		return nil, err
	}

	serviceNetworkServingSignerKey, err := GetSecretData(ctx, "service-network-serving-signer", "openshift-kube-apiserver-operator", "tls.key", client)
	if err != nil {
		return nil, err
	}

	ingressKey, err := GetSecretData(ctx, "router-ca", "openshift-ingress-operator", "tls.key", client)
	if err != nil {
		return nil, err
	}

	kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.LoadbalancerSignerPrivateKey = seedreconfig.PEM(loadbalancerServingSignerPrivateKey)
	kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.LocalhostSignerPrivateKey = seedreconfig.PEM(localhostServingSignerKey)
	kubeconfigCryptoRetention.KubeAPICrypto.ServingCrypto.ServiceNetworkSignerPrivateKey = seedreconfig.PEM(serviceNetworkServingSignerKey)
	kubeconfigCryptoRetention.KubeAPICrypto.ClientAuthCrypto.AdminCACertificate = seedreconfig.PEM(adminKubeConfigClientCA)
	kubeconfigCryptoRetention.IngresssCrypto.IngressCA = seedreconfig.PEM(ingressKey)

	return &kubeconfigCryptoRetention, nil
}

func BackupKubeconfigCrypto(ctx context.Context, client runtimeclient.Client, cryptoDir string) error {
	if err := os.MkdirAll(cryptoDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating %s: %w", cryptoDir, err)
	}

	adminKubeConfigClientCA, err := GetConfigMapData(ctx, "admin-kubeconfig-client-ca", "openshift-config", "ca-bundle.crt", client)
	if err != nil {
		return fmt.Errorf("failed to get configMap data with adminKubeConfigClientCA: %w", err)
	}
	p := path.Join(cryptoDir, "admin-kubeconfig-client-ca.crt")
	if err := os.WriteFile(p, []byte(adminKubeConfigClientCA), cryptoDirMode); err != nil {
		return fmt.Errorf("failed to admin-kubeconfig-client-ca.crt to path %s: %w", p, err)
	}

	for _, cert := range common.CertPrefixes {
		servingSignerKey, err := GetSecretData(ctx, cert, "openshift-kube-apiserver-operator", "tls.key", client)
		if err != nil {
			return fmt.Errorf("failed to get secret data with servingSignerKey: %w", err)
		}
		curP := path.Join(cryptoDir, cert+".key")
		if err := os.WriteFile(curP, []byte(servingSignerKey), cryptoDirMode); err != nil {
			return fmt.Errorf("failed write to .key file to path %s: %w", curP, err)
		}
	}

	ingressOperatorKey, err := GetSecretData(ctx, "router-ca", "openshift-ingress-operator", "tls.key", client)
	if err != nil {
		return fmt.Errorf("failed to get secret data with ingressOperatorKey: %w", err)
	}
	p = path.Join(cryptoDir, "ingresskey-ingress-operator.key")
	if err := os.WriteFile(p, []byte(ingressOperatorKey), cryptoDirMode); err != nil {
		return fmt.Errorf("failed to ingresskey-ingress-operator.key to path %s: %w", p, err)
	}

	return nil
}

func BackupKubeadminPasswordHash(ctx context.Context, client runtimeclient.Client, cryptoDir string) (bool, error) {
	password, err := GetSecretData(ctx, "kubeadmin", "kube-system", "kubeadmin", client)
	if err != nil {
		if runtimeclient.IgnoreNotFound(err) == nil {
			return false, nil
		}

		return false, fmt.Errorf("failed to get kubeadmin password: %w", err)
	}

	if err := os.WriteFile(path.Join(cryptoDir, pwHashFile), []byte(password), cryptoDirMode); err != nil {
		return false, fmt.Errorf("failed to write kubeadmin password hash: %w", err)
	}

	return true, nil
}

func LoadKubeadminPasswordHash(cryptoDir string) (string, error) {
	if _, err := os.Stat(path.Join(cryptoDir, pwHashFile)); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", fmt.Errorf("failed to check if kubeadmin password hash exists: %w", err)
	}

	password, err := os.ReadFile(path.Join(cryptoDir, pwHashFile))
	if err != nil {
		return "", fmt.Errorf("failed to read kubeadmin password hash: %w", err)
	}

	return string(password), err
}
