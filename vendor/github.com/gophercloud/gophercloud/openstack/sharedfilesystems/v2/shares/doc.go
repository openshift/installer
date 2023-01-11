/*
Package shares provides information and interaction with the different
API versions for the Shared File System service, code-named Manila.

For more information, see:
https://docs.openstack.org/api-ref/shared-file-system/

Example to Revert a Share to a Snapshot ID

	opts := &shares.RevertOpts{
		// snapshot ID to revert to
		SnapshotID: "ddeac769-9742-497f-b985-5bcfa94a3fd6",
	}
	manilaClient.Microversion = "2.27"
	err := shares.Revert(manilaClient, shareID, opts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Reset a Share Status

	opts := &shares.ResetStatusOpts{
		// a new Share Status
		Status: "available",
	}
	manilaClient.Microversion = "2.7"
	err := shares.ResetStatus(manilaClient, shareID, opts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Force Delete a Share

	manilaClient.Microversion = "2.7"
	err := shares.ForceDelete(manilaClient, shareID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Unmanage a Share

	manilaClient.Microversion = "2.7"
	err := shares.Unmanage(manilaClient, shareID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package shares
