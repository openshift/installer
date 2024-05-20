/*
Package zones provides information and interaction with the zone API
resource for the OpenStack DNS service.

Example to List Zone Transfer Accepts

        // Optionaly you can provide Status as query parameter for filtering the result.
        allPages, err := transferAccepts.List(dnsClient, nil).AllPages()
        if err != nil {
        	panic(err)
        }

        allTransferAccepts, err := transferAccepts.ExtractTransferAccepts(allPages)
        if err != nil {
        	panic(err)
        }

        for _, transferAccept := range allTransferAccepts {
        	fmt.Printf("%+v\n", transferAccept)
        }

Example to Create a Zone Transfer Accept

        zoneTransferRequestID := "99d10f68-5623-4491-91a0-6daafa32b60e"
        key := "JKHGD2F7"
        createOpts := transferAccepts.CreateOpts{
        	ZoneTransferRequestID: zoneTransferRequestID,
        	Key: key,
        }
        transferAccept, err := transferAccepts.Create(dnsClient, createOpts).Extract()
        if err != nil {
        	panic(err)
        }

Example to Get a Zone Transfer Accept

        transferAcceptID := "99d10f68-5623-4491-91a0-6daafa32b60e"
        transferAccept, err := transferAccepts.Get(dnsClient, transferAcceptID).Extract()
        if err != nil {
        	panic(err)
        }
*/
package accept
