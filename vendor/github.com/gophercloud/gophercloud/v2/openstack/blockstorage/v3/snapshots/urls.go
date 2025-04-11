package snapshots

import "github.com/gophercloud/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("snapshots")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return createURL(c)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func metadataURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id, "metadata")
}

func updateMetadataURL(c *gophercloud.ServiceClient, id string) string {
	return metadataURL(c, id)
}

func resetStatusURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}

func updateStatusURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}

func forceDeleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id, "action")
}
