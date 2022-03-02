/*
Package zones provides information and interaction with the zone API
resource for the OpenStack DNS service.

Example to List Zone Transfer Requests

        allPages, err := transferRequests.List(dnsClient, nil).AllPages()
        if err != nil {
        	panic(err)
        }

        allTransferRequests, err := transferRequests.ExtractTransferRequests(allPages)
        if err != nil {
        	panic(err)
        }

        for _, transferRequest := range allTransferRequests {
        	fmt.Printf("%+v\n", transferRequest)
        }

Example to Create a Zone Transfer Request

        zoneID := "99d10f68-5623-4491-91a0-6daafa32b60e"
        targetProjectID := "f977bd7c-6485-4385-b04f-b5af0d186fcc"
        createOpts := transferRequests.CreateOpts{
                TargetProjectID: targetProjectID,
        	Description: "This is a zone transfer request.",
        }
        transferRequest, err := transferRequests.Create(dnsClient, zoneID, createOpts).Extract()
        if err != nil {
        	panic(err)
        }

Example to Delete a Zone Transfer Request

        transferID := "99d10f68-5623-4491-91a0-6daafa32b60e"
        err := transferRequests.Delete(dnsClient, transferID).ExtractErr()
        if err != nil {
        	panic(err)
        }
*/
package request
