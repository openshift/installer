package ec2credentials

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List enumerates the Credentials to which the current token has access.
func List(client *gophercloud.ServiceClient, userID string) pagination.Pager {
	url := listURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return CredentialPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single EC2 credential by ID.
func Get(client *gophercloud.ServiceClient, userID string, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, userID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts provides options used to create an EC2 credential.
type CreateOpts struct {
	// TenantID is the project ID scope of the EC2 credential.
	TenantID string `json:"tenant_id" required:"true"`
}

// ToCredentialCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToCredentialCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new EC2 Credential.
func Create(client *gophercloud.ServiceClient, userID string, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client, userID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes an EC2 credential.
func Delete(client *gophercloud.ServiceClient, userID string, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, userID, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
