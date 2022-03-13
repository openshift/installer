package qos

import "github.com/gophercloud/gophercloud"

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id)
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("qos-specs")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("qos-specs")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("qos-specs", id)
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id)
}

func deleteKeysURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id, "delete_keys")
}

func associateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id, "associate")
}

func disassociateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id, "disassociate")
}

func disassociateAllURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id, "disassociate_all")
}

func listAssociationsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("qos-specs", id, "associations")
}
