// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isVolumeProfiles = "profiles"
)

func dataSourceIBMISVolumeProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISVolumeProfilesRead,

		Schema: map[string]*schema.Schema{

			isVolumeProfiles: {
				Type:        schema.TypeList,
				Description: "List of Volume profile maps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISVolumeProfilesRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classicVolumeProfilesList(d, meta)
		if err != nil {
			return err
		}
	} else {
		err := volumeProfilesList(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicVolumeProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcclassicv1.VolumeProfile{}
	for {
		listVolumeProfilesOptions := &vpcclassicv1.ListVolumeProfilesOptions{}
		if start != "" {
			listVolumeProfilesOptions.Start = &start
		}
		availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
		if err != nil {
			return fmt.Errorf("Error Fetching Volume Profiles %s\n%s", err, response)
		}
		start = GetNext(availableProfiles.Next)
		allrecs = append(allrecs, availableProfiles.Profiles...)
		if start == "" {
			break
		}
	}
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range allrecs {

		l := map[string]interface{}{
			"name":   *profile.Name,
			"family": *profile.Family,
		}
		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMISVolumeProfilesID(d))
	d.Set(isVolumeProfiles, profilesInfo)
	return nil
}

func volumeProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{}
	availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Volume Profiles %s\n%s", err, response)
	}
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range availableProfiles.Profiles {

		l := map[string]interface{}{
			"name":   *profile.Name,
			"family": *profile.Family,
		}
		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMISVolumeProfilesID(d))
	d.Set(isVolumeProfiles, profilesInfo)
	return nil
}

// dataSourceIBMISVolumeProfilesID returns a reasonable ID for a Volume Profile list.
func dataSourceIBMISVolumeProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
