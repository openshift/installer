/*
Package portstrustedvif provides information and interaction with the port
trusted vif extension which adds trusted attributes to the port resource
for the OpenStack Networking service.

Example to Get a Port with a Port Trusted VIF

	var portWithPortTrustedVIFExtension struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	portID := "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2"

	err := ports.Get(context.TODO(), networkingClient, portID).ExtractInto(&portWithPortTrustedVIFExtension)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", portWithPortTrustedVIFExtension)

Example to Create a Port With Port Trusted VIF set as true

	var portWithPortTrustedVIFExtension struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	iTrue := true
	networkID := "4e8e5957-649f-477b-9e5b-f1f75b21c03c"
	subnetID := "a87cc70a-3e15-4acf-8205-9b711a3531b7"

	portCreateOpts := ports.CreateOpts{
		NetworkID: networkID,
		FixedIPs:  []ports.IP{ports.IP{SubnetID: subnetID}},
	}

	createOpts := portstrustedvif.PortCreateOptsExt{
		CreateOptsBuilder:   portCreateOpts,
		PortTrustedVIF: &iTrue,
	}

	err := ports.Create(context.TODO(), networkingClient, createOpts).ExtractInto(&portWithPortTrustedVIFExtension)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", portWithPortTrustedVIFExtension)

Example to Update the status of the Port Trusted VIF to false on an Existing Port

	var portWithPortTrustedVIDExtension struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	iFalse := false
	portID := "65c0ee9f-d634-4522-8954-51021b570b0d"

	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := portstrustedvif.PortUpdateOptsExt{
		UpdateOptsBuilder:   portUpdateOpts,
		PortTrustedVIF: &iFalse,
	}

	err := ports.Update(context.TODO(), networkingClient, portID, updateOpts).ExtractInto(&portWithPortTrustedVIFExtension)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", portWithPortTrustedVIFExtension)
*/

package portstrustedvif
