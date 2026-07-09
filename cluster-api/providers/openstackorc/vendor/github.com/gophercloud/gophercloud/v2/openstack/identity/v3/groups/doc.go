/*
Package groups manages and retrieves Groups in the OpenStack Identity Service.

Example to List Groups

	listOpts := groups.ListOpts{
		DomainID: "default",
	}

	allPages, err := groups.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		panic(err)
	}

	for _, group := range allGroups {
		fmt.Printf("%+v\n", group)
	}

Example to Create a Group

	createOpts := groups.CreateOpts{
		Name:             "groupname",
		DomainID:         "default",
		Extra: map[string]any{
			"email": "groupname@example.com",
		}
	}

	group, err := groups.Create(context.TODO(), identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Group

	groupID := "0fe36e73809d46aeae6705c39077b1b3"

	updateOpts := groups.UpdateOpts{
		Description: "Updated Description for group",
	}

	group, err := groups.Update(context.TODO(), identityClient, groupID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Group

	groupID := "0fe36e73809d46aeae6705c39077b1b3"
	err := groups.Delete(context.TODO(), identityClient, groupID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package groups
