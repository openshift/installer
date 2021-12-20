package nodes

/*
Package nodes provides utilities for working with Ironic's baremetal API.

* Building a config drive

As part of provisioning a node, you may need a config drive that contains user data, metadata, and network data
stored inside a base64-encoded gzipped ISO9660 file.  These utilities can create that for you.

For example:

	configDrive = nodes.ConfigDrive{
		UserData: nodes.UserDataMap{
		"ignition": map[string]string{
			"version": "2.2.0",
		},
		"systemd": map[string]interface{}{
			"units": []map[string]interface{}{{
				"name":    "example.service",
				"enabled": true,
			},
			},
		},
	}

Then to upload this to Ironic as a using gophercloud:

	err = nodes.ChangeProvisionState(client, uuid, nodes.ProvisionStateOpts{
		Target:      "active",
		ConfigDrive: configDrive.ToConfigDrive(),
	}).ExtractErr()

*/
