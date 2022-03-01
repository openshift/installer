/*
Package quotas provides the ability to retrieve and manage Load Balancer quotas

Example to Get project quotas

    projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"
    quotasInfo, err := quotas.Get(networkClient, projectID).Extract()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("quotas: %#v\n", quotasInfo)

Example to Update project quotas

    projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"

    updateOpts := quotas.UpdateOpts{
		Loadbalancer:  gophercloud.IntToPointer(20),
		Listener:      gophercloud.IntToPointer(40),
		Member:        gophercloud.IntToPointer(200),
		Pool:          gophercloud.IntToPointer(20),
		Healthmonitor: gophercloud.IntToPointer(1),
		L7Policy:      gophercloud.IntToPointer(50),
		L7Rule:        gophercloud.IntToPointer(100),
    }
    quotasInfo, err := quotas.Update(networkClient, projectID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("quotas: %#v\n", quotasInfo)
*/
package quotas
