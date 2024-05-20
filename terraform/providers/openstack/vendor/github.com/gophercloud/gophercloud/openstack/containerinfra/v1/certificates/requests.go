package certificates

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder allows extensions to add additional parameters
// to the Create request.
type CreateOptsBuilder interface {
	ToCertificateCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a certificate.
type CreateOpts struct {
	ClusterUUID string `json:"cluster_uuid,omitempty" xor:"BayUUID"`
	BayUUID     string `json:"bay_uuid,omitempty" xor:"ClusterUUID"`
	CSR         string `json:"csr" required:"true"`
}

// ToCertificateCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToCertificateCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Get makes a request against the API to get details for a certificate.
func Get(client *gophercloud.ServiceClient, clusterID string) (r GetResult) {
	url := getURL(client, clusterID)
	resp, err := client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Create requests the creation of a new certificate.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCertificateCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Update will rotate the CA certificate for a cluster
func Update(client *gophercloud.ServiceClient, clusterID string) (r UpdateResult) {
	resp, err := client.Patch(updateURL(client, clusterID), nil, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
