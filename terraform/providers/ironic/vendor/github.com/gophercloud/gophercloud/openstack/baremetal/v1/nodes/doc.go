/*
Package nodes provides information and interaction with the nodes API
resource in the OpenStack Bare Metal service.

Example to List Nodes with Detail

	nodes.ListDetail(client, nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		nodeList, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		for _, n := range nodeList {
			// Do something
		}

		return true, nil
	})

Example to List Nodes

	listOpts := nodes.ListOpts{
		ProvisionState: nodes.Deploying,
		Fields:         []string{"name"},
	}

	nodes.List(client, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		nodeList, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		for _, n := range nodeList {
			// Do something
		}

		return true, nil
	})

Example to Create Node

	createOpts := nodes.CreateOpts
		Driver:        "ipmi",
		BootInterface: "pxe",
		Name:          "coconuts",
		DriverInfo: map[string]interface{}{
			"ipmi_port":      "6230",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "admin",
		},
	}

	createNode, err := nodes.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get Node

	showNode, err := nodes.Get(client, "c9afd385-5d89-4ecb-9e1c-68194da6b474").Extract()
	if err != nil {
		panic(err)
	}

Example to Update Node

	updateOpts := nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    ReplaceOp,
			Path:  "/maintenance",
			Value: "true",
		},
	}

	updateNode, err := nodes.Update(client, "c9afd385-5d89-4ecb-9e1c-68194da6b474", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete Node

	err = nodes.Delete(client, "c9afd385-5d89-4ecb-9e1c-68194da6b474").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Validate Node

	validation, err := nodes.Validate(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8").Extract()
	if err != nil {
		panic(err)
	}

Example to inject non-masking interrupts

	err := nodes.InjectNMI(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to get array of supported boot devices for a node

	bootDevices, err := nodes.GetSupportedBootDevices(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8").Extract()
	if err != nil {
		panic(err)
	}

Example to set boot device for a node

	bootOpts := nodes.BootDeviceOpts{
		BootDevice: "pxe",
		Persistent: false,
	}

	err := nodes.SetBootDevice(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8", bootOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to get boot device for a node

	bootDevice, err := nodes.GetBootDevice(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8").Extract()
	if err != nil {
		panic(err)
	}

Example to list all vendor passthru methods

	methods, err := nodes.GetVendorPassthruMethods(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8").Extract()
	if err != nil {
		panic(err)
	}

Example to list all subscriptions

	method := nodes.CallVendorPassthruOpts{
		Method: "get_all_subscriptions",
	}
	allSubscriptions, err := nodes.GetAllSubscriptions(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8", method).Extract()
	if err != nil {
		panic(err)
	}

Example to get a subscription

	method := nodes.CallVendorPassthruOpts{
		Method: "get_subscription",
	}
	subscriptionOpt := nodes.GetSubscriptionOpts{
		Id:     "subscription id",
	}

	subscription, err := nodes.GetSubscription(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8", method, subscriptionOpt).Extract()
	if err != nil {
		panic(err)
	}

Example to delete a subscription

	method := nodes.CallVendorPassthruOpts{
		Method: "delete_subscription",
	}
	subscriptionDeleteOpt := nodes.DeleteSubscriptionOpts{
		Id: "subscription id",
	}

	err := nodes.DeleteSubscription(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8", method, subscriptionDeleteOpt).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to create a subscription

	method := nodes.CallVendorPassthruOpts{
		Method: "create_subscription",
	}
	subscriptionCreateOpt := nodes.CreateSubscriptionOpts{
		Destination: "https://subscription_destination_url"
		Context:     "MyContext",
		Protocol:    "Redfish",
		EventTypes:  ["Alert"],
		HttpHeaders: [{"Key1":"Value1"}, {"Key2":"Value2"}],
	}

	newSubscription, err := nodes.CreateSubscription(client, "a62b8495-52e2-407b-b3cb-62775d04c2b8", method, subscriptionCreateOpt).Extract()
	if err != nil {
		panic(err)
	}
*/
package nodes
