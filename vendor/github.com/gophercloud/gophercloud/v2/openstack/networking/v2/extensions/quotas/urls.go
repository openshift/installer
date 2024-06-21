package quotas

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "quotas"
const resourcePathDetail = "details.json"

func resourceURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID)
}

func resourceDetailURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePath, projectID, resourcePathDetail)
}

func getURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceURL(c, projectID)
}

func getDetailURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceDetailURL(c, projectID)
}

func updateURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceURL(c, projectID)
}
