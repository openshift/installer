/*
Package snapshots provides information and interaction with snapshots in the
OpenStack Block Storage service. A snapshot is a point in time copy of the
data contained in an external storage volume, and can be controlled
programmatically.

Example to list Snapshots

	allPages, err := snapshots.List(client, snapshots.ListOpts{}).AllPages(context.TODO())
	if err != nil{
		panic(err)
	}
	snapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil{
		panic(err)
	}
	for _,s := range snapshots{
		fmt.Println(s)
	}

Example to get a Snapshot

	snapshotID := "4a584cae-e4ce-429b-9154-d4c9eb8fda4c"
	snapshot, err := snapshots.Get(context.TODO(), client, snapshotID).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(snapshot)

Example to create a Snapshot

	snapshot, err := snapshots.Create(context.TODO(), client, snapshots.CreateOpts{
		Name:"snapshot_001",
		VolumeID:"5aa119a8-d25b-45a7-8d1b-88e127885635",
	}).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(snapshot)

Example to delete a Snapshot

	snapshotID := "4a584cae-e4ce-429b-9154-d4c9eb8fda4c"
	err := snapshots.Delete(context.TODO(), client, snapshotID).ExtractErr()
	if err != nil{
		panic(err)
	}

Example to update a Snapshot

	snapshotID := "4a584cae-e4ce-429b-9154-d4c9eb8fda4c"
	snapshot, err = snapshots.Update(context.TODO(), client, snapshotID, snapshots.UpdateOpts{
		Name: "snapshot_002",
		Description:"description_002",
	}).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(snapshot)
*/
package snapshots
