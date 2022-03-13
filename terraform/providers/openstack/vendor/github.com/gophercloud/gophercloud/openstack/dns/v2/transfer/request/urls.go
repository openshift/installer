package request

import "github.com/gophercloud/gophercloud"

const (
	rootPath     = "zones"
	tasksPath    = "tasks"
	resourcePath = "transfer_requests"
)

func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath)
}

func createURL(c *gophercloud.ServiceClient, zoneID string) string {
	return c.ServiceURL(rootPath, zoneID, tasksPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, transferID string) string {
	return c.ServiceURL(rootPath, tasksPath, resourcePath, transferID)
}
