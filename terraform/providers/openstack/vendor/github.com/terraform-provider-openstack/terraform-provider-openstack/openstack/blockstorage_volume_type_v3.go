package openstack

import (
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func expandBlockStorageVolumeTypeV3ExtraSpecs(raw map[string]interface{}) volumetypes.ExtraSpecsOpts {
	extraSpecs := make(volumetypes.ExtraSpecsOpts, len(raw))
	for k, v := range raw {
		extraSpecs[k] = v.(string)
	}

	return extraSpecs
}
