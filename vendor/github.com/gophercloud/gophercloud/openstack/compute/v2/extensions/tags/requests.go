package tags

import "github.com/gophercloud/gophercloud"

// List all tags on a server.
func List(client *gophercloud.ServiceClient, serverID string) (r ListResult) {
	url := listURL(client, serverID)
	resp, err := client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Check if a tag exists on a server.
func Check(client *gophercloud.ServiceClient, serverID, tag string) (r CheckResult) {
	url := checkURL(client, serverID, tag)
	resp, err := client.Get(url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ReplaceAllOptsBuilder allows to add additional parameters to the ReplaceAll request.
type ReplaceAllOptsBuilder interface {
	ToTagsReplaceAllMap() (map[string]interface{}, error)
}

// ReplaceAllOpts provides options used to replace Tags on a server.
type ReplaceAllOpts struct {
	Tags []string `json:"tags" required:"true"`
}

// ToTagsReplaceAllMap formats a ReplaceALlOpts into the body of the ReplaceAll request.
func (opts ReplaceAllOpts) ToTagsReplaceAllMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ReplaceAll replaces all Tags on a server.
func ReplaceAll(client *gophercloud.ServiceClient, serverID string, opts ReplaceAllOptsBuilder) (r ReplaceAllResult) {
	b, err := opts.ToTagsReplaceAllMap()
	url := replaceAllURL(client, serverID)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(url, &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Add adds a new Tag on a server.
func Add(client *gophercloud.ServiceClient, serverID, tag string) (r AddResult) {
	url := addURL(client, serverID, tag)
	resp, err := client.Put(url, nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a tag from a server.
func Delete(client *gophercloud.ServiceClient, serverID, tag string) (r DeleteResult) {
	url := deleteURL(client, serverID, tag)
	resp, err := client.Delete(url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteAll removes all tag from a server.
func DeleteAll(client *gophercloud.ServiceClient, serverID string) (r DeleteResult) {
	url := deleteAllURL(client, serverID)
	resp, err := client.Delete(url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
