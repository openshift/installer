package ignition

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/openshift/installer/pkg/asset"
)

// tokenLen is how many bytes of random input go into generating the token.
// It should be a multiple of 3 to render nicely in base64
// encoding.  The current default is long; it's a lot of entropy.  The MCS
// will have mitigations against brute forcing.
const tokenLen = 30

// ProvisioningToken implements https://github.com/openshift/enhancements/pull/443/
type ProvisioningToken struct {
	Token string
}

var _ asset.Asset = (*ProvisioningToken)(nil)

// Dependencies returns no dependencies.
func (a *ProvisioningToken) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate the token
func (a *ProvisioningToken) Generate(asset.Parents) error {
	b := make([]byte, tokenLen)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	a.Token = base64.StdEncoding.EncodeToString(b)
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *ProvisioningToken) Name() string {
	return "Ignition Provisioning Password"
}
