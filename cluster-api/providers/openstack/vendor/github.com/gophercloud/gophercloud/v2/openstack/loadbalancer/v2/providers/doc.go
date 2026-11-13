/*
Package providers provides information about the supported providers
at OpenStack Octavia Load Balancing service.

Example to List Providers

	allPages, err := providers.List(lbClient).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allProviders, err := providers.ExtractProviders(allPages)
	if err != nil {
		panic(err)
	}

	for _, p := range allProviders {
		fmt.Printf("%+v\n", p)
	}

Example to List Provider Flavor Capabilities

	capabilities, err := providers.ListFlavorCapabilities(context.TODO(), lbClient, "amphora", nil).Extract()
	if err != nil {
		panic(err)
	}

	for _, capability := range capabilities {
		fmt.Printf("%+v\n", capability)
	}

Example to List Provider Availability Zone Capabilities

	capabilities, err := providers.ListAvailabilityZoneCapabilities(context.TODO(), lbClient, "amphora", nil).Extract()
	if err != nil {
		panic(err)
	}

	for _, capability := range capabilities {
		fmt.Printf("%+v\n", capability)
	}
*/
package providers
