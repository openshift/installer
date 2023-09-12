/*
Package volumetypes provides information and interaction with volume types in the
OpenStack Block Storage service. A volume type is a collection of specs used to
define the volume capabilities.

Example to list Volume Types

	allPages, err := volumetypes.List(client, volumetypes.ListOpts{}).AllPages()
	if err != nil{
		panic(err)
	}
	volumeTypes, err := volumetypes.ExtractVolumeTypes(allPages)
	if err != nil{
		panic(err)
	}
	for _,vt := range volumeTypes{
		fmt.Println(vt)
	}

Example to show a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	volumeType, err := volumetypes.Get(client, typeID).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(volumeType)

Example to create a Volume Type

	volumeType, err := volumetypes.Create(client, volumetypes.CreateOpts{
		Name:"volume_type_001",
		IsPublic:true,
		Description:"description_001",
	}).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(volumeType)

Example to delete a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	err := volumetypes.Delete(client, typeID).ExtractErr()
	if err != nil{
		panic(err)
	}

Example to update a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	volumetype, err = volumetypes.Update(client, typeID, volumetypes.UpdateOpts{
		Name: "volume_type_002",
		Description:"description_002",
		IsPublic:false,
	}).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(volumetype)

Example to Create Extra Specs for a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"

	createOpts := volumetypes.ExtraSpecsOpts{
		"capabilities": "gpu",
	}
	createdExtraSpecs, err := volumetypes.CreateExtraSpecs(client, typeID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", createdExtraSpecs)

Example to Get Extra Specs for a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"

	extraSpecs, err := volumetypes.ListExtraSpecs(client, typeID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", extraSpecs)

Example to Get specific Extra Spec for a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"

	extraSpec, err := volumetypes.GetExtraSpec(client, typeID, "capabilities").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", extraSpec)

Example to Update Extra Specs for a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"

	updateOpts := volumetypes.ExtraSpecsOpts{
		"capabilities": "capabilities-updated",
	}
	updatedExtraSpec, err := volumetypes.UpdateExtraSpec(client, typeID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", updatedExtraSpec)

Example to Delete an Extra Spec for a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	err := volumetypes.DeleteExtraSpec(client, typeID, "capabilities").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List Volume Type Access

	typeID := "e91758d6-a54a-4778-ad72-0c73a1cb695b"

	allPages, err := volumetypes.ListAccesses(client, typeID).AllPages()
	if err != nil {
		panic(err)
	}

	allAccesses, err := volumetypes.ExtractAccesses(allPages)
	if err != nil {
		panic(err)
	}

	for _, access := range allAccesses {
		fmt.Printf("%+v", access)
	}

Example to Grant Access to a Volume Type

	typeID := "e91758d6-a54a-4778-ad72-0c73a1cb695b"

	accessOpts := volumetypes.AddAccessOpts{
		Project: "15153a0979884b59b0592248ef947921",
	}

	err := volumetypes.AddAccess(client, typeID, accessOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Remove/Revoke Access to a Volume Type

	typeID := "e91758d6-a54a-4778-ad72-0c73a1cb695b"

	accessOpts := volumetypes.RemoveAccessOpts{
		Project: "15153a0979884b59b0592248ef947921",
	}

	err := volumetypes.RemoveAccess(client, typeID, accessOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Create the Encryption of a Volume Type

	  typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
		volumeType, err := volumetypes.CreateEncryption(client, typeID, .CreateEncryptionOpts{
			KeySize:      256,
			Provider:    "luks",
			ControlLocation: "front-end",
			Cipher:  "aes-xts-plain64",
		}).Extract()
		if err != nil{
			panic(err)
		}
		fmt.Println(volumeType)

Example to Delete the Encryption of a Volume Type

		typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	  encryptionID := ""81e069c6-7394-4856-8df7-3b237ca61f74
		err := volumetypes.DeleteEncryption(client, typeID, encryptionID).ExtractErr()
		if err != nil{
			panic(err)
		}

Example to Update the Encryption of a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	volumetype, err = volumetypes.UpdateEncryption(client, typeID, volumetypes.UpdateEncryptionOpts{
		KeySize:      256,
		Provider:    "luks",
		ControlLocation: "front-end",
		Cipher:  "aes-xts-plain64",
	}).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(volumetype)

Example to Show an Encryption of a Volume Type

	typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	volumeType, err := volumetypes.GetEncrytpion(client, typeID).Extract()
	if err != nil{
		panic(err)
	}
	fmt.Println(volumeType)

Example to Show an Encryption Spec of a Volume Type

		typeID := "7ffaca22-f646-41d4-b79d-d7e4452ef8cc"
	  key := "cipher"
		volumeType, err := volumetypes.GetEncrytpionSpec(client, typeID).Extract()
		if err != nil{
			panic(err)
		}
		fmt.Println(volumeType)
*/
package volumetypes
