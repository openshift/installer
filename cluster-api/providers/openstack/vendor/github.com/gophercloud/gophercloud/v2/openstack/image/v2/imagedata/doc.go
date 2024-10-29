/*
Package imagedata enables management of image data.

Example to Upload Image Data

	imageID := "da3b75d9-3f4a-40e7-8a2c-bfab23927dea"

	imageData, err := os.Open("/path/to/image/file")
	if err != nil {
		panic(err)
	}
	defer imageData.Close()

	err = imagedata.Upload(context.TODO(), imageClient, imageID, imageData).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Stage Image Data

	imageID := "da3b75d9-3f4a-40e7-8a2c-bfab23927dea"

	imageData, err := os.Open("/path/to/image/file")
	if err != nil {
	  panic(err)
	}
	defer imageData.Close()

	err = imagedata.Stage(context.TODO(), imageClient, imageID, imageData).ExtractErr()
	if err != nil {
	  panic(err)
	}

Example to Download Image Data

	imageID := "da3b75d9-3f4a-40e7-8a2c-bfab23927dea"

	image, err := imagedata.Download(context.TODO(), imageClient, imageID).Extract()
	if err != nil {
		panic(err)
	}

	// close the reader, when reading has finished
	defer image.Close()

	imageData, err := io.ReadAll(image)
	if err != nil {
		panic(err)
	}
*/
package imagedata
