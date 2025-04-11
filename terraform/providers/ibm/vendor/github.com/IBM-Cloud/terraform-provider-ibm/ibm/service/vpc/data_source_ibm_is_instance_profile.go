// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceProfileName         = "name"
	isInstanceProfileFamily       = "family"
	isInstanceProfileArchitecture = "architecture"
	isInstanceVCPUArchitecture    = "vcpu_architecture"
	isInstanceVCPUManufacturer    = "vcpu_manufacturer"
)

func DataSourceIBMISInstanceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceProfileRead,

		Schema: map[string]*schema.Schema{

			isInstanceProfileName: {
				Type:     schema.TypeString,
				Required: true,
			},

			// cluster changes
			"cluster_network_attachment_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"default": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"supported_cluster_network_profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cluster network profiles that support this instance profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this cluster network profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this cluster network profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},

			"confidential_compute_modes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default confidential compute mode for this profile.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The supported confidential compute modes.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"secure_boot_modes": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The default secure boot mode for this profile.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The supported `enable_secure_boot` values for an instance using this profile.",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
					},
				},
			},

			isInstanceProfileFamily: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this virtual server instance profile belongs to.",
			},

			isInstanceProfileArchitecture: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default OS architecture for an instance with this profile.",
			},

			"architecture_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type for the OS architecture.",
			},

			"architecture_values": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The supported OS architecture(s) for an instance with this profile.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"bandwidth": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"gpu_count": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU count of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"gpu_manufacturer": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU manufacturer of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The possible GPU manufacturer(s) for an instance with this profile",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"gpu_memory": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU memory of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"gpu_model": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GPU model of this profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The possible GPU model(s) for an instance with this profile",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"reservation_terms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The type for this profile field",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The supported committed use terms for a reservation using this profile",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total_volume_bandwidth": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The amount of bandwidth (in megabits per second) allocated exclusively to instance storage volumes. An increase in this value will result in a corresponding decrease to total_network_bandwidth.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of the instance profile's disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"quantity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"size": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field.",
									},
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default value for this profile field.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum value for this profile field.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum value for this profile field.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The permitted values for this profile field.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"supported_interface_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk interface used for attaching the disk.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"values": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The supported disk interfaces used for attaching the disk.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this virtual server instance profile.",
			},
			"memory": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"network_interface_count": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
					},
				},
			},
			"network_attachment_count": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
					},
				},
			},
			"numa_count": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
					},
				},
			},
			"port_speed": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the instance profile.",
			},
			isInstanceVCPUArchitecture: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture for an instance with this profile.",
						},
					},
				},
			},
			isInstanceVCPUManufacturer: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU manufacturer for an instance with this profile.",
						},
					},
				},
			},
			"vcpu_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceProfileRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isInstanceProfileName).(string)
	err := instanceProfileGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func instanceProfileGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getInstanceProfileOptions := &vpcv1.GetInstanceProfileOptions{
		Name: &name,
	}
	profile, _, err := sess.GetInstanceProfile(getInstanceProfileOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from profile name.
	d.SetId(*profile.Name)
	d.Set(isInstanceProfileName, *profile.Name)
	d.Set(isInstanceProfileFamily, *profile.Family)
	if profile.OsArchitecture != nil {
		if profile.OsArchitecture.Default != nil {
			d.Set(isInstanceProfileArchitecture, *profile.OsArchitecture.Default)
		}
		if profile.OsArchitecture.Type != nil {
			d.Set("architecture_type", *profile.OsArchitecture.Type)
		}
		if profile.OsArchitecture.Values != nil {
			d.Set("architecture_values", *&profile.OsArchitecture.Values)
		}

	}
	if profile.Status != nil {
		d.Set("status", profile.Status)
	}

	// cluster changes
	clusterNetworkAttachmentCount := []map[string]interface{}{}
	clusterNetworkAttachmentCountMap, err := DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountToMap(profile.ClusterNetworkAttachmentCount)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_profile", "read", "cluster_network_attachment_count-to-map")
	}
	clusterNetworkAttachmentCount = append(clusterNetworkAttachmentCount, clusterNetworkAttachmentCountMap)
	if err = d.Set("cluster_network_attachment_count", clusterNetworkAttachmentCount); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cluster_network_attachment_count: %s", err), "(Data) ibm_is_instance_profile", "read", "set-cluster_network_attachment_count")
	}

	supportedClusterNetworkProfiles := []map[string]interface{}{}
	for _, supportedClusterNetworkProfilesItem := range profile.SupportedClusterNetworkProfiles {
		supportedClusterNetworkProfilesItemMap, err := DataSourceIBMIsInstanceProfileClusterNetworkProfileReferenceToMap(&supportedClusterNetworkProfilesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_profile", "read", "supported_cluster_network_profiles-to-map")
		}
		supportedClusterNetworkProfiles = append(supportedClusterNetworkProfiles, supportedClusterNetworkProfilesItemMap)
	}
	if err = d.Set("supported_cluster_network_profiles", supportedClusterNetworkProfiles); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting supported_cluster_network_profiles: %s", err), "(Data) ibm_is_instance_profile", "read", "set-supported_cluster_network_profiles")
	}

	confidentialComputeModes := []map[string]interface{}{}
	if profile.ConfidentialComputeModes != nil {
		modelMap, err := dataSourceIBMIsInstanceProfileInstanceProfileSupportedConfidentialComputeModesToMap(profile.ConfidentialComputeModes)
		if err != nil {
			return (err)
		}
		confidentialComputeModes = append(confidentialComputeModes, modelMap)
	}
	if err = d.Set("confidential_compute_modes", confidentialComputeModes); err != nil {
		return fmt.Errorf("Error setting confidential_compute_modes %s", err)
	}

	secureBootModes := []map[string]interface{}{}
	if profile.SecureBootModes != nil {
		modelMap, err := dataSourceIBMIsInstanceProfileInstanceProfileSupportedSecureBootModesToMap(profile.SecureBootModes)
		if err != nil {
			return err
		}
		secureBootModes = append(secureBootModes, modelMap)
	}
	if err = d.Set("secure_boot_modes", secureBootModes); err != nil {
		return fmt.Errorf("Error setting secure_boot_modes %s", err)
	}

	if profile.Bandwidth != nil {
		err = d.Set("bandwidth", dataSourceInstanceProfileFlattenBandwidth(*profile.Bandwidth.(*vpcv1.InstanceProfileBandwidth)))
		if err != nil {
			return err
		}
	}
	if profile.GpuCount != nil {
		err = d.Set("gpu_count", dataSourceInstanceProfileFlattenGPUCount(*profile.GpuCount.(*vpcv1.InstanceProfileGpu)))
		if err != nil {
			return err
		}
	}
	if profile.GpuMemory != nil {
		err = d.Set("gpu_memory", dataSourceInstanceProfileFlattenGPUMemory(*profile.GpuMemory.(*vpcv1.InstanceProfileGpuMemory)))
		if err != nil {
			return err
		}
	}
	if profile.GpuManufacturer != nil {
		err = d.Set("gpu_manufacturer", dataSourceInstanceProfileFlattenGPUManufacturer(*profile.GpuManufacturer))
		if err != nil {
			return err
		}
	}
	if profile.GpuModel != nil {
		err = d.Set("gpu_model", dataSourceInstanceProfileFlattenGPUModel(*profile.GpuModel))
		if err != nil {
			return err
		}
	}
	if profile.ReservationTerms != nil {
		err = d.Set("reservation_terms", dataSourceInstanceProfileFlattenReservationTerms(*profile.ReservationTerms))
		if err != nil {
			return err
		}
	}
	if profile.TotalVolumeBandwidth != nil {
		err = d.Set("total_volume_bandwidth", dataSourceInstanceProfileFlattenTotalVolumeBandwidth(*profile.TotalVolumeBandwidth.(*vpcv1.InstanceProfileVolumeBandwidth)))
		if err != nil {
			return err
		}
	}
	if profile.Disks != nil {
		err = d.Set("disks", dataSourceInstanceProfileFlattenDisks(profile.Disks))
		if err != nil {
			return err
		}
	}
	if err = d.Set("href", profile.Href); err != nil {
		return err
	}

	if profile.Memory != nil {
		err = d.Set("memory", dataSourceInstanceProfileFlattenMemory(*profile.Memory.(*vpcv1.InstanceProfileMemory)))
		if err != nil {
			return err
		}
	}
	if profile.NetworkInterfaceCount != nil {
		err = d.Set("network_interface_count", dataSourceInstanceProfileFlattenNetworkInterfaceCount(*profile.NetworkInterfaceCount.(*vpcv1.InstanceProfileNetworkInterfaceCount)))
		if err != nil {
			return err
		}
	}
	if profile.NetworkAttachmentCount != nil {
		err = d.Set("network_attachment_count", dataSourceInstanceProfileFlattenNetworkAttachmentCount(*profile.NetworkAttachmentCount.(*vpcv1.InstanceProfileNetworkAttachmentCount)))
		if err != nil {
			return err
		}
	}
	if profile.NumaCount != nil {
		err = d.Set("numa_count", dataSourceInstanceProfileFlattenNumaCount(*profile.NumaCount.(*vpcv1.InstanceProfileNumaCount)))
		if err != nil {
			return err
		}
	}
	if profile.PortSpeed != nil {
		err = d.Set("port_speed", dataSourceInstanceProfileFlattenPortSpeed(*profile.PortSpeed.(*vpcv1.InstanceProfilePortSpeed)))
		if err != nil {
			return err
		}
	}

	if profile.VcpuArchitecture != nil {
		err = d.Set(isInstanceVCPUArchitecture, dataSourceInstanceProfileFlattenVcpuArchitecture(*profile.VcpuArchitecture))
		if err != nil {
			return err
		}
	}

	// Manufacturer details added.
	if profile.VcpuManufacturer != nil {
		err = d.Set(isInstanceVCPUManufacturer, dataSourceInstanceProfileFlattenVcpuManufacture(*profile.VcpuManufacturer))
		if err != nil {
			return err
		}
	}

	if profile.VcpuCount != nil {
		err = d.Set("vcpu_count", dataSourceInstanceProfileFlattenVcpuCount(*profile.VcpuCount.(*vpcv1.InstanceProfileVcpu)))
		if err != nil {
			return err
		}
	}
	return nil
}

func dataSourceInstanceProfileFlattenBandwidth(result vpcv1.InstanceProfileBandwidth) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileBandwidthToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileBandwidthToMap(bandwidthItem vpcv1.InstanceProfileBandwidth) (bandwidthMap map[string]interface{}) {
	bandwidthMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		bandwidthMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		bandwidthMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		bandwidthMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		bandwidthMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		bandwidthMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		bandwidthMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		bandwidthMap["values"] = bandwidthItem.Values
	}

	return bandwidthMap
}

func dataSourceInstanceProfileFlattenGPUCount(result vpcv1.InstanceProfileGpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUCountToMap(bandwidthItem vpcv1.InstanceProfileGpu) (gpuCountMap map[string]interface{}) {
	gpuCountMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuCountMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		gpuCountMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		gpuCountMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		gpuCountMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		gpuCountMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		gpuCountMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		gpuCountMap["values"] = bandwidthItem.Values
	}

	return gpuCountMap
}

func dataSourceInstanceProfileFlattenGPUManufacturer(result vpcv1.InstanceProfileGpuManufacturer) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUManufacturerToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUManufacturerToMap(bandwidthItem vpcv1.InstanceProfileGpuManufacturer) (gpuManufactrerMap map[string]interface{}) {
	gpuManufactrerMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuManufactrerMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Values != nil {
		gpuManufactrerMap["values"] = bandwidthItem.Values
	}

	return gpuManufactrerMap
}

func dataSourceInstanceProfileFlattenGPUModel(result vpcv1.InstanceProfileGpuModel) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUModelToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUModelToMap(bandwidthItem vpcv1.InstanceProfileGpuModel) (gpuModel map[string]interface{}) {
	gpuModelMap := map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuModelMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Values != nil {
		gpuModelMap["values"] = bandwidthItem.Values
	}

	return gpuModelMap
}

func dataSourceInstanceProfileFlattenReservationTerms(result vpcv1.InstanceProfileReservationTerms) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileReservationTermsToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileReservationTermsToMap(resTermItem vpcv1.InstanceProfileReservationTerms) map[string]interface{} {
	resTermMap := map[string]interface{}{}

	if resTermItem.Type != nil {
		resTermMap["type"] = resTermItem.Type
	}
	if resTermItem.Values != nil {
		resTermMap["values"] = resTermItem.Values
	}

	return resTermMap
}

func dataSourceInstanceProfileFlattenGPUMemory(result vpcv1.InstanceProfileGpuMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileGPUMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileGPUMemoryToMap(bandwidthItem vpcv1.InstanceProfileGpuMemory) (gpuMemoryMap map[string]interface{}) {
	gpuMemoryMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		gpuMemoryMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		gpuMemoryMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		gpuMemoryMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		gpuMemoryMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		gpuMemoryMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		gpuMemoryMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		gpuMemoryMap["values"] = bandwidthItem.Values
	}

	return gpuMemoryMap
}

func dataSourceInstanceProfileFlattenMemory(result vpcv1.InstanceProfileMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileMemoryToMap(memoryItem vpcv1.InstanceProfileMemory) (memoryMap map[string]interface{}) {
	memoryMap = map[string]interface{}{}

	if memoryItem.Type != nil {
		memoryMap["type"] = memoryItem.Type
	}
	if memoryItem.Value != nil {
		memoryMap["value"] = memoryItem.Value
	}
	if memoryItem.Default != nil {
		memoryMap["default"] = memoryItem.Default
	}
	if memoryItem.Max != nil {
		memoryMap["max"] = memoryItem.Max
	}
	if memoryItem.Min != nil {
		memoryMap["min"] = memoryItem.Min
	}
	if memoryItem.Step != nil {
		memoryMap["step"] = memoryItem.Step
	}
	if memoryItem.Values != nil {
		memoryMap["values"] = memoryItem.Values
	}

	return memoryMap
}

func dataSourceInstanceProfileFlattenNetworkInterfaceCount(result vpcv1.InstanceProfileNetworkInterfaceCount) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileNetworkInterfaceCount(result)
	finalList = append(finalList, finalMap)

	return finalList
}
func dataSourceInstanceProfileFlattenNetworkAttachmentCount(result vpcv1.InstanceProfileNetworkAttachmentCount) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileNetworkAttachmentCount(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileNetworkAttachmentCount(networkAttachmentCountItem vpcv1.InstanceProfileNetworkAttachmentCount) (networkAttachmentCountMap map[string]interface{}) {
	networkAttachmentCountMap = map[string]interface{}{}

	if networkAttachmentCountItem.Max != nil {
		networkAttachmentCountMap["max"] = networkAttachmentCountItem.Max
	}
	if networkAttachmentCountItem.Min != nil {
		networkAttachmentCountMap["min"] = networkAttachmentCountItem.Min
	}
	if networkAttachmentCountItem.Type != nil {
		networkAttachmentCountMap["type"] = networkAttachmentCountItem.Type
	}
	return networkAttachmentCountMap
}

func dataSourceInstanceProfileNetworkInterfaceCount(networkInterfaceCountItem vpcv1.InstanceProfileNetworkInterfaceCount) (networkInterfaceCountMap map[string]interface{}) {
	networkInterfaceCountMap = map[string]interface{}{}

	if networkInterfaceCountItem.Max != nil {
		networkInterfaceCountMap["max"] = networkInterfaceCountItem.Max
	}
	if networkInterfaceCountItem.Min != nil {
		networkInterfaceCountMap["min"] = networkInterfaceCountItem.Min
	}
	if networkInterfaceCountItem.Type != nil {
		networkInterfaceCountMap["type"] = networkInterfaceCountItem.Type
	}
	return networkInterfaceCountMap
}

func dataSourceInstanceProfileFlattenPortSpeed(result vpcv1.InstanceProfilePortSpeed) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfilePortSpeedToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfilePortSpeedToMap(portSpeedItem vpcv1.InstanceProfilePortSpeed) (portSpeedMap map[string]interface{}) {
	portSpeedMap = map[string]interface{}{}

	if portSpeedItem.Type != nil {
		portSpeedMap["type"] = portSpeedItem.Type
	}
	if portSpeedItem.Value != nil {
		portSpeedMap["value"] = portSpeedItem.Value
	}

	return portSpeedMap
}

func dataSourceInstanceProfileFlattenVcpuArchitecture(result vpcv1.InstanceProfileVcpuArchitecture) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileVcpuArchitectureToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileVcpuArchitectureToMap(vcpuArchitectureItem vpcv1.InstanceProfileVcpuArchitecture) (vcpuArchitectureMap map[string]interface{}) {
	vcpuArchitectureMap = map[string]interface{}{}

	if vcpuArchitectureItem.Type != nil {
		vcpuArchitectureMap["type"] = vcpuArchitectureItem.Type
	}
	if vcpuArchitectureItem.Value != nil {
		vcpuArchitectureMap["value"] = vcpuArchitectureItem.Value
	}

	return vcpuArchitectureMap
}

/* Changes for the AMD Support VCPU Manufacturer */
func dataSourceInstanceProfileFlattenVcpuManufacture(result vpcv1.InstanceProfileVcpuManufacturer) (fl []map[string]interface{}) {
	fl = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileVcpuManufacturerToMap(result)
	fl = append(fl, finalMap)

	return fl
}

func dataSourceInstanceProfileVcpuManufacturerToMap(vcpuManufacutererItem vpcv1.InstanceProfileVcpuManufacturer) (vcpuManufacturerMap map[string]interface{}) {
	vcpuManufacturerMap = map[string]interface{}{}

	if vcpuManufacutererItem.Type != nil {
		vcpuManufacturerMap["type"] = vcpuManufacutererItem.Type
	}
	if vcpuManufacutererItem.Value != nil {
		vcpuManufacturerMap["value"] = vcpuManufacutererItem.Value
	}

	return vcpuManufacturerMap
}

func dataSourceInstanceProfileFlattenVcpuCount(result vpcv1.InstanceProfileVcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileVcpuCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileVcpuCountToMap(vcpuCountItem vpcv1.InstanceProfileVcpu) (vcpuCountMap map[string]interface{}) {
	vcpuCountMap = map[string]interface{}{}

	if vcpuCountItem.Type != nil {
		vcpuCountMap["type"] = vcpuCountItem.Type
	}
	if vcpuCountItem.Value != nil {
		vcpuCountMap["value"] = vcpuCountItem.Value
	}
	if vcpuCountItem.Default != nil {
		vcpuCountMap["default"] = vcpuCountItem.Default
	}
	if vcpuCountItem.Max != nil {
		vcpuCountMap["max"] = vcpuCountItem.Max
	}
	if vcpuCountItem.Min != nil {
		vcpuCountMap["min"] = vcpuCountItem.Min
	}
	if vcpuCountItem.Step != nil {
		vcpuCountMap["step"] = vcpuCountItem.Step
	}
	if vcpuCountItem.Values != nil {
		vcpuCountMap["values"] = vcpuCountItem.Values
	}

	return vcpuCountMap
}

func dataSourceInstanceProfileFlattenDisks(result []vpcv1.InstanceProfileDisk) (disks []map[string]interface{}) {
	for _, disksItem := range result {
		disks = append(disks, dataSourceInstanceProfileDisksToMap(disksItem))
	}

	return disks
}

func dataSourceInstanceProfileDisksToMap(disksItem vpcv1.InstanceProfileDisk) (disksMap map[string]interface{}) {
	disksMap = map[string]interface{}{}

	if disksItem.Quantity != nil {
		quantityList := []map[string]interface{}{}
		quantityMap := dataSourceInstanceProfileDisksQuantityToMap(*disksItem.Quantity.(*vpcv1.InstanceProfileDiskQuantity))
		quantityList = append(quantityList, quantityMap)
		disksMap["quantity"] = quantityList
	}
	if disksItem.Size != nil {
		sizeList := []map[string]interface{}{}
		sizeMap := dataSourceInstanceProfileDisksSizeToMap(*disksItem.Size.(*vpcv1.InstanceProfileDiskSize))
		sizeList = append(sizeList, sizeMap)
		disksMap["size"] = sizeList
	}
	if disksItem.SupportedInterfaceTypes != nil {
		supportedInterfaceTypesList := []map[string]interface{}{}
		supportedInterfaceTypesMap := dataSourceInstanceProfileDisksSupportedInterfaceTypesToMap(*disksItem.SupportedInterfaceTypes)
		supportedInterfaceTypesList = append(supportedInterfaceTypesList, supportedInterfaceTypesMap)
		disksMap["supported_interface_types"] = supportedInterfaceTypesList
	}

	return disksMap
}

func dataSourceInstanceProfileDisksQuantityToMap(quantityItem vpcv1.InstanceProfileDiskQuantity) (quantityMap map[string]interface{}) {
	quantityMap = map[string]interface{}{}

	if quantityItem.Type != nil {
		quantityMap["type"] = quantityItem.Type
	}
	if quantityItem.Value != nil {
		quantityMap["value"] = quantityItem.Value
	}
	if quantityItem.Default != nil {
		quantityMap["default"] = quantityItem.Default
	}
	if quantityItem.Max != nil {
		quantityMap["max"] = quantityItem.Max
	}
	if quantityItem.Min != nil {
		quantityMap["min"] = quantityItem.Min
	}
	if quantityItem.Step != nil {
		quantityMap["step"] = quantityItem.Step
	}
	if quantityItem.Values != nil {
		quantityMap["values"] = quantityItem.Values
	}

	return quantityMap
}

func dataSourceInstanceProfileDisksSizeToMap(sizeItem vpcv1.InstanceProfileDiskSize) (sizeMap map[string]interface{}) {
	sizeMap = map[string]interface{}{}

	if sizeItem.Type != nil {
		sizeMap["type"] = sizeItem.Type
	}
	if sizeItem.Value != nil {
		sizeMap["value"] = sizeItem.Value
	}
	if sizeItem.Default != nil {
		sizeMap["default"] = sizeItem.Default
	}
	if sizeItem.Max != nil {
		sizeMap["max"] = sizeItem.Max
	}
	if sizeItem.Min != nil {
		sizeMap["min"] = sizeItem.Min
	}
	if sizeItem.Step != nil {
		sizeMap["step"] = sizeItem.Step
	}
	if sizeItem.Values != nil {
		sizeMap["values"] = sizeItem.Values
	}

	return sizeMap
}

func dataSourceInstanceProfileDisksSupportedInterfaceTypesToMap(supportedInterfaceTypesItem vpcv1.InstanceProfileDiskSupportedInterfaces) (supportedInterfaceTypesMap map[string]interface{}) {
	supportedInterfaceTypesMap = map[string]interface{}{}

	if supportedInterfaceTypesItem.Default != nil {
		supportedInterfaceTypesMap["default"] = supportedInterfaceTypesItem.Default
	}
	if supportedInterfaceTypesItem.Type != nil {
		supportedInterfaceTypesMap["type"] = supportedInterfaceTypesItem.Type
	}
	if supportedInterfaceTypesItem.Values != nil {
		supportedInterfaceTypesMap["values"] = supportedInterfaceTypesItem.Values
	}

	return supportedInterfaceTypesMap
}

func dataSourceInstanceProfileFlattenTotalVolumeBandwidth(result vpcv1.InstanceProfileVolumeBandwidth) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileTotalVolumeBandwidthToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileTotalVolumeBandwidthToMap(bandwidthItem vpcv1.InstanceProfileVolumeBandwidth) (bandwidthMap map[string]interface{}) {
	bandwidthMap = map[string]interface{}{}

	if bandwidthItem.Type != nil {
		bandwidthMap["type"] = bandwidthItem.Type
	}
	if bandwidthItem.Value != nil {
		bandwidthMap["value"] = bandwidthItem.Value
	}
	if bandwidthItem.Default != nil {
		bandwidthMap["default"] = bandwidthItem.Default
	}
	if bandwidthItem.Max != nil {
		bandwidthMap["max"] = bandwidthItem.Max
	}
	if bandwidthItem.Min != nil {
		bandwidthMap["min"] = bandwidthItem.Min
	}
	if bandwidthItem.Step != nil {
		bandwidthMap["step"] = bandwidthItem.Step
	}
	if bandwidthItem.Values != nil {
		bandwidthMap["values"] = bandwidthItem.Values
	}

	return bandwidthMap
}

func dataSourceInstanceProfileFlattenNumaCount(result vpcv1.InstanceProfileNumaCount) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceInstanceProfileNumaCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceInstanceProfileNumaCountToMap(numaItem vpcv1.InstanceProfileNumaCount) (numaMap map[string]interface{}) {
	numaMap = map[string]interface{}{}

	if numaItem.Type != nil {
		numaMap["type"] = numaItem.Type
	}
	if numaItem.Value != nil {
		numaMap["value"] = numaItem.Value
	}

	return numaMap
}

func dataSourceIBMIsInstanceProfileInstanceProfileSupportedSecureBootModesToMap(model *vpcv1.InstanceProfileSupportedSecureBootModes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = model.Default
	modelMap["type"] = model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func dataSourceIBMIsInstanceProfileInstanceProfileSupportedConfidentialComputeModesToMap(model *vpcv1.InstanceProfileSupportedConfidentialComputeModes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["default"] = model.Default
	modelMap["type"] = model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountToMap(model vpcv1.InstanceProfileClusterNetworkAttachmentCountIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountDependent); ok {
		return DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountDependentToMap(model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountDependent))
	} else if _, ok := model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountEnum); ok {
		return DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountEnumToMap(model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountEnum))
	} else if _, ok := model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountRange); ok {
		return DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountRangeToMap(model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCountRange))
	} else if _, ok := model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCount); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.InstanceProfileClusterNetworkAttachmentCount)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Default != nil {
			modelMap["default"] = flex.IntValue(model.Default)
		}
		if model.Values != nil {
			modelMap["values"] = model.Values
		}
		if model.Max != nil {
			modelMap["max"] = flex.IntValue(model.Max)
		}
		if model.Min != nil {
			modelMap["min"] = flex.IntValue(model.Min)
		}
		if model.Step != nil {
			modelMap["step"] = flex.IntValue(model.Step)
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.InstanceProfileClusterNetworkAttachmentCountIntf subtype encountered")
	}
}

func DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountDependentToMap(model *vpcv1.InstanceProfileClusterNetworkAttachmentCountDependent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountEnumToMap(model *vpcv1.InstanceProfileClusterNetworkAttachmentCountEnum) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Default != nil {
		modelMap["default"] = flex.IntValue(model.Default)
	}
	modelMap["type"] = *model.Type
	modelMap["values"] = model.Values
	return modelMap, nil
}

func DataSourceIBMIsInstanceProfileInstanceProfileClusterNetworkAttachmentCountRangeToMap(model *vpcv1.InstanceProfileClusterNetworkAttachmentCountRange) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Max != nil {
		modelMap["max"] = flex.IntValue(model.Max)
	}
	if model.Min != nil {
		modelMap["min"] = flex.IntValue(model.Min)
	}
	if model.Step != nil {
		modelMap["step"] = flex.IntValue(model.Step)
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIBMIsInstanceProfileClusterNetworkProfileReferenceToMap(model *vpcv1.ClusterNetworkProfileReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
