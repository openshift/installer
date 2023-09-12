package volumetypes

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("types")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("types", id)
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("types")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("types", id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("types", id)
}

func extraSpecsListURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "extra_specs")
}

func extraSpecsGetURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func extraSpecsCreateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "extra_specs")
}

func extraSpecUpdateURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func extraSpecDeleteURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("types", id, "extra_specs", key)
}

func accessURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "os-volume-type-access")
}

func accessActionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "action")
}

func createEncryptionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "encryption")
}

func deleteEncryptionURL(client *gophercloud.ServiceClient, id, encryptionID string) string {
	return client.ServiceURL("types", id, "encryption", encryptionID)
}

func getEncryptionURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("types", id, "encryption")
}

func getEncryptionSpecURL(client *gophercloud.ServiceClient, id, key string) string {
	return client.ServiceURL("types", id, "encryption", key)
}

func updateEncryptionURL(client *gophercloud.ServiceClient, id, encryptionID string) string {
	return client.ServiceURL("types", id, "encryption", encryptionID)
}
