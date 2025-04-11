package shares

import "github.com/gophercloud/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("shares")
}

func listDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("shares", "detail")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id)
}

func listExportLocationsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "export_locations")
}

func getExportLocationURL(c *gophercloud.ServiceClient, shareID, id string) string {
	return c.ServiceURL("shares", shareID, "export_locations", id)
}

func grantAccessURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func revokeAccessURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func listAccessRightsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func extendURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func shrinkURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func revertURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func resetStatusURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func forceDeleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func unmanageURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "action")
}

func getMetadataURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func getMetadatumURL(c *gophercloud.ServiceClient, id, key string) string {
	return c.ServiceURL("shares", id, "metadata", key)
}

func setMetadataURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func updateMetadataURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("shares", id, "metadata")
}

func deleteMetadatumURL(c *gophercloud.ServiceClient, id, key string) string {
	return c.ServiceURL("shares", id, "metadata", key)
}
