/*
Package imageimport enables management of images import and retrieval of the
Image service Import API information.

Example to Get an information about the Import API

	importInfo, err := imageimport.Get(context.TODO(), imagesClient).Extract()
	if err != nil {
	  panic(err)
	}

	fmt.Printf("%+v\n", importInfo)

Example to Create a new image import

	createOpts := imageimport.CreateOpts{
	  Name: imageimport.WebDownloadMethod,
	  URI:  "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img",
	}
	imageID := "da3b75d9-3f4a-40e7-8a2c-bfab23927dea"

	err := imageimport.Create(context.TODO(), imagesClient, imageID, createOpts).ExtractErr()
	if err != nil {
	  panic(err)
	}
*/
package imageimport
