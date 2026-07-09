/*
Package domains manages and retrieves Domains in the OpenStack Identity Service.

Example to List Domains

	var iTrue = true
	listOpts := domains.ListOpts{
		Enabled: &iTrue,
	}

	allPages, err := domains.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allDomains, err := domains.ExtractDomains(allPages)
	if err != nil {
		panic(err)
	}

	for _, domain := range allDomains {
		fmt.Printf("%+v\n", domain)
	}

Example to Create a Domain

	createOpts := domains.CreateOpts{
		Name:             "domain name",
		Description:      "Test domain",
	}

	domain, err := domains.Create(context.TODO(), identityClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Domain

	domainID := "0fe36e73809d46aeae6705c39077b1b3"

	var iFalse = false
	updateOpts := domains.UpdateOpts{
		Enabled: &iFalse,
	}

	domain, err := domains.Update(context.TODO(), identityClient, domainID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Domain

	domainID := "0fe36e73809d46aeae6705c39077b1b3"
	err := domains.Delete(context.TODO(), identityClient, domainID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package domains
