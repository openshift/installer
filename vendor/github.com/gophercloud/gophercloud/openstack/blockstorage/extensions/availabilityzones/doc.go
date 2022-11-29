/*
Package availabilityzones provides the ability to get lists of
available volume availability zones.

Example of Get Availability Zone Information

		allPages, err := availabilityzones.List(volumeClient).AllPages()
		if err != nil {
			panic(err)
		}

		availabilityZoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
		if err != nil {
			panic(err)
		}

		for _, zoneInfo := range availabilityZoneInfo {
	  		fmt.Printf("%+v\n", zoneInfo)
		}
*/
package availabilityzones
