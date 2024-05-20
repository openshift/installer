package imagedata

import (
	"io"

	"github.com/gophercloud/gophercloud"
)

// Upload uploads an image file.
func Upload(client *gophercloud.ServiceClient, id string, data io.Reader) (r UploadResult) {
	resp, err := client.Put(uploadURL(client, id), data, nil, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/octet-stream"},
		OkCodes:     []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Stage performs PUT call on the existing image object in the Imageservice with
// the provided file.
// Existing image object must be in the "queued" status.
func Stage(client *gophercloud.ServiceClient, id string, data io.Reader) (r StageResult) {
	resp, err := client.Put(stageURL(client, id), data, nil, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/octet-stream"},
		OkCodes:     []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Download retrieves an image.
func Download(client *gophercloud.ServiceClient, id string) (r DownloadResult) {
	resp, err := client.Get(downloadURL(client, id), nil, &gophercloud.RequestOpts{
		KeepResponseBody: true,
	})
	r.Body, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
