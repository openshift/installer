package imageimport

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// ImportMethod represents valid Import API method.
type ImportMethod string

const (
	// GlanceDirectMethod represents glance-direct Import API method.
	GlanceDirectMethod ImportMethod = "glance-direct"

	// WebDownloadMethod represents web-download Import API method.
	WebDownloadMethod ImportMethod = "web-download"
)

// Get retrieves Import API information data.
func Get(ctx context.Context, c *gophercloud.ServiceClient) (r GetResult) {
	resp, err := c.Get(ctx, infoURL(c), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToImportCreateMap() (map[string]any, error)
}

// CreateOpts specifies parameters of a new image import.
type CreateOpts struct {
	Name ImportMethod `json:"name"`
	URI  string       `json:"uri"`
}

// ToImportCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToImportCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return map[string]any{"method": b}, nil
}

// Create requests the creation of a new image import on the server.
func Create(ctx context.Context, client *gophercloud.ServiceClient, imageID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToImportCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, importURL(client, imageID), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
