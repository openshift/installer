package installconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

func TestSkipPermsCheckForCredentialsMode(t *testing.T) {
	assert.False(t, skipPermsCheckForCredentialsMode(ibmcloud.Name, string(types.ManualCredentialsMode)),
		"IBM Cloud with Manual must still run PlatformPermsCheck")
	assert.True(t, skipPermsCheckForCredentialsMode(aws.Name, string(types.ManualCredentialsMode)),
		"AWS with credentialsMode set must skip PlatformPermsCheck")
	assert.True(t, skipPermsCheckForCredentialsMode(aws.Name, string(types.MintCredentialsMode)))
	assert.False(t, skipPermsCheckForCredentialsMode(aws.Name, ""),
		"AWS with unset credentialsMode must not skip")
	assert.False(t, skipPermsCheckForCredentialsMode(ibmcloud.Name, ""),
		"IBM Cloud with unset credentialsMode must not skip")
}
