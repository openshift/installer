package domains

import "github.com/gophercloud/gophercloud/v2"

func listAvailableURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("auth", "domains")
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("domains")
}

func getURL(client *gophercloud.ServiceClient, domainID string) string {
	return client.ServiceURL("domains", domainID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("domains")
}

func deleteURL(client *gophercloud.ServiceClient, domainID string) string {
	return client.ServiceURL("domains", domainID)
}

func updateURL(client *gophercloud.ServiceClient, domainID string) string {
	return client.ServiceURL("domains", domainID)
}
