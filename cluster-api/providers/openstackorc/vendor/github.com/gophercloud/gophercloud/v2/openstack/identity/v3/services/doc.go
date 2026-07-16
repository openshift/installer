/*
Package services provides information and interaction with the services API
resource for the OpenStack Identity service.

Example to List Services

	listOpts := services.ListOpts{
		ServiceType: "compute",
	}

	allPages, err := services.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allServices, err := services.ExtractServices(allPages)
	if err != nil {
		panic(err)
	}

	for _, service := range allServices {
		fmt.Printf("%+v\n", service)
	}

Example to Create a Service

	createOpts := services.CreateOpts{
		Type: "compute",
		Extra: map[string]any{
			"name": "compute-service",
			"description": "Compute Service",
		},
	}

	service, err := services.Create(context.TODO(), identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Service

	serviceID :=  "3c7bbe9a6ecb453ca1789586291380ed"

	var iFalse = false
	updateOpts := services.UpdateOpts{
		Enabled: &iFalse,
		Extra: map[string]any{
			"description": "Disabled Compute Service"
		},
	}

	service, err := services.Update(context.TODO(), identityClient, serviceID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Service

	serviceID := "3c7bbe9a6ecb453ca1789586291380ed"
	err := services.Delete(context.TODO(), identityClient, serviceID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package services
