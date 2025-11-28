/*
Package projects manages and retrieves Projects in the OpenStack Identity
Service.

Example to List Projects

	listOpts := projects.ListOpts{
		Enabled: gophercloud.Enabled,
	}

	allPages, err := projects.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		panic(err)
	}

	for _, project := range allProjects {
		fmt.Printf("%+v\n", project)
	}

Example to Create a Project

	createOpts := projects.CreateOpts{
		Name:        "project_name",
		Description: "Project Description",
		Tags:        []string{"FirstTag", "SecondTag"},
	}

	project, err := projects.Create(context.TODO(), identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Project

	projectID := "966b3c7d36a24facaf20b7e458bf2192"

	updateOpts := projects.UpdateOpts{
		Enabled: gophercloud.Disabled,
	}

	project, err := projects.Update(context.TODO(), identityClient, projectID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

	updateOpts = projects.UpdateOpts{
		Tags: &[]string{"FirstTag"},
	}

	project, err = projects.Update(context.TODO(), identityClient, projectID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Project

	projectID := "966b3c7d36a24facaf20b7e458bf2192"
	err := projects.Delete(context.TODO(), identityClient, projectID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List all tags of a Project

	projectID := "966b3c7d36a24facaf20b7e458bf2192"
	err := projects.ListTags(context.TODO(), identityClient, projectID).Extract()
	if err != nil {
		panic(err)
	}

Example to modify all tags of a Project

	projectID := "966b3c7d36a24facaf20b7e458bf2192"
	tags := ["foo", "bar"]
	projects, err := projects.ModifyTags(context.TODO(), identityClient, projectID, tags).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete all tags of a Project

	projectID := "966b3c7d36a24facaf20b7e458bf2192"
	err := projects.DeleteTags(context.TODO(), identityClient, projectID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package projects
