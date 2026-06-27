package pki

import (
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

// CertificateConfigToUpstream converts the installer-local CertificateConfig
// to the upstream configv1alpha1.CertificateConfig for use in the PKI CR manifest.
func CertificateConfigToUpstream(local types.CertificateConfig) configv1alpha1.CertificateConfig {
	return configv1alpha1.CertificateConfig{
		Key: configv1alpha1.KeyConfig{
			Algorithm: configv1alpha1.KeyAlgorithm(local.Key.Algorithm),
			RSA:       configv1alpha1.RSAKeyConfig{KeySize: local.Key.RSA.KeySize},
			ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurve(local.Key.ECDSA.Curve)},
		},
	}
}
