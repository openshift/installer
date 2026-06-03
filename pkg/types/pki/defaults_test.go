package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
)

func TestDefaultPKIProfile(t *testing.T) {
	profile := DefaultPKIProfile()

	assert.Equal(t, configv1alpha1.KeyAlgorithmRSA, profile.Defaults.Key.Algorithm)
	assert.Equal(t, int32(4096), profile.Defaults.Key.RSA.KeySize)
	assert.Equal(t, configv1alpha1.KeyAlgorithmRSA, profile.SignerCertificates.Key.Algorithm)
	assert.Equal(t, int32(4096), profile.SignerCertificates.Key.RSA.KeySize)
}
