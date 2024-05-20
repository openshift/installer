/*
Package ec2credentials provides information and interaction with the EC2
credentials API resource for the OpenStack Identity service.

For more information, see:
https://docs.openstack.org/api-ref/identity/v2-ext/

Example to Create an EC2 credential

	createOpts := ec2credentials.CreateOpts{
		// project ID of the EC2 credential scope
		TenantID: projectID,
	}

	credential, err := ec2credentials.Create(identityClient, userID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package ec2credentials
