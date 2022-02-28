/*
Package apiversions provides information and interaction with the different
API versions for the Shared File System service, code-named Manila.

Example to List API Versions

	allPages, err := apiversions.List(client).AllPages()
	if err != nil {
		panic(err)
	}

	allVersions, err := apiversions.ExtractAPIVersions(allPages)
	if err != nil {
		panic(err)
	}

	for _, version := range allVersions {
		fmt.Printf("%+v\n", version)
	}

Example to Get an API Version

	version, err := apiVersions.Get(client, "v2.1").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", version)
*/
package apiversions
