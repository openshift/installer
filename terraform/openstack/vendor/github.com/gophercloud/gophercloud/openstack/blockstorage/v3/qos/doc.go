/*
Package qos provides information and interaction with the QoS specifications
for the Openstack Blockstorage service.

Example to create a QoS specification

	createOpts := qos.CreateOpts{
		Name:     "test",
		Consumer: qos.ConsumerFront,
		Specs: map[string]string{
			"read_iops_sec": "20000",
		},
	}

	test, err := qos.Create(client, createOpts).Extract()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("QoS: %+v\n", test)

Example to delete a QoS specification

	qosID := "d6ae28ce-fcb5-4180-aa62-d260a27e09ae"

	deleteOpts := qos.DeleteOpts{
		Force: false,
	}

	err = qos.Delete(client, qosID, deleteOpts).ExtractErr()
	if err != nil {
		log.Fatal(err)
	}

Example to list QoS specifications

	listOpts := qos.ListOpts{}

	allPages, err := qos.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allQoS, err := qos.ExtractQoS(allPages)
	if err != nil {
		panic(err)
	}

	for _, qos := range allQoS {
		fmt.Printf("List: %+v\n", qos)
	}

Example to get a single QoS specification

	qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"

	singleQos, err := qos.Get(client, test.ID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Get: %+v\n", singleQos)

Example of updating QoSSpec

	qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"

	updateOpts := qos.UpdateOpts{
		Consumer: qos.ConsumerBack,
		Specs: map[string]string{
			"read_iops_sec": "40000",
		},
	}

	specs, err := qos.Update(client, qosID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", specs)


Example of deleting specific keys/specs from a QoS

	qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"

	keysToDelete := qos.DeleteKeysOpts{"read_iops_sec"}
	err = qos.DeleteKeys(client, qosID, keysToDelete).ExtractErr()
	if err != nil {
		panic(err)
	}

	Example of associating a QoS with a volume type

qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"
volID := "b596be6a-0ce9-43fa-804a-5c5e181ede76"

associateOpts := qos.AssociateOpts{
	VolumeTypeID: volID,
}

err = qos.Associate(client, qosID, associateOpts).ExtractErr()
if err != nil {
	panic(err)
}

Example of disassociating a QoS from a volume type

qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"
volID := "b596be6a-0ce9-43fa-804a-5c5e181ede76"

disassociateOpts := qos.DisassociateOpts{
	VolumeTypeID: volID,
}

err = qos.Disassociate(client, qosID, disassociateOpts).ExtractErr()
if err != nil {
	panic(err)
}

Example of disaassociating a Qos from all volume types

qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"

err = qos.DisassociateAll(client, qosID).ExtractErr()
if err != nil {
	panic(err)
}

Example of listing all associations of a QoS

qosID := "de075d5e-8afc-4e23-9388-b84a5183d1c0"

allQosAssociations, err := qos.ListAssociations(client, qosID).AllPages()
if err != nil {
	panic(err)
}

allAssociations, err := qos.ExtractAssociations(allQosAssociations)
if err != nil {
	panic(err)
}

for _, association := range allAssociations {
	fmt.Printf("Association: %+v\n", association)
}

*/
package qos
