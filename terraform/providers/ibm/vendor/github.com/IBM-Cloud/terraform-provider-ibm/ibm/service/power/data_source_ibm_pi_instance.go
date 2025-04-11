// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DedicatedHostID: {
				Computed:    true,
				Description: "The dedicated host ID where the shared processor pool resides.",
				Type:        schema.TypeString,
			},
			Attr_DeploymentType: {
				Computed:    true,
				Description: "The custom deployment type.",
				Type:        schema.TypeString,
			},
			Attr_Fault: {
				Computed:    true,
				Description: "Fault information.",
				Type:        schema.TypeMap,
			},
			Attr_HealthStatus: {
				Computed:    true,
				Description: "The health of the instance.",
				Type:        schema.TypeString,
			},
			Attr_IBMiCSS: {
				Computed:    true,
				Description: "IBMi Cloud Storage Solution",
				Type:        schema.TypeBool,
			},
			Attr_IBMiPHA: {
				Computed:    true,
				Description: "IBMi Power High Availability",
				Type:        schema.TypeBool,
			},
			Attr_IBMiRDS: {
				Computed:    true,
				Description: "IBMi Rational Dev Studio",
				Type:        schema.TypeBool,
			},
			Attr_IBMiRDSUsers: {
				Computed:    true,
				Description: "IBMi Rational Dev Studio Number of User Licenses",
				Type:        schema.TypeInt,
			},
			Attr_LicenseRepositoryCapacity: {
				Computed:    true,
				Deprecated:  "This field is deprecated.",
				Description: "The VTL license repository capacity TB value.",
				Type:        schema.TypeInt,
			},
			Attr_MaxMem: {
				Computed:    true,
				Description: "The maximum amount of memory that can be allocated to the instance without shutting down or rebooting the LPAR.",
				Type:        schema.TypeFloat,
			},
			Attr_MaxProc: {
				Computed:    true,
				Description: "The maximum number of processors that can be allocated to the instance without shutting down or rebooting the LPAR.",
				Type:        schema.TypeFloat,
			},
			Attr_MaxVirtualCores: {
				Computed:    true,
				Description: "The maximum number of virtual cores that can be assigned without rebooting the instance.",
				Type:        schema.TypeInt,
			},
			Attr_Memory: {
				Computed:    true,
				Description: "The amount of memory that is allocated to the instance.",
				Type:        schema.TypeFloat,
			},
			Attr_MinMem: {
				Computed:    true,
				Description: "The minimum amount of memory that must be allocated to the instance.",
				Type:        schema.TypeFloat,
			},
			Attr_MinProc: {
				Computed:    true,
				Description: "The minimum number of processors that must be allocated to the instance.",
				Type:        schema.TypeFloat,
			},
			Attr_MinVirtualCores: {
				Computed:    true,
				Description: "The minimum number of virtual cores that can be assigned without rebooting the instance.",
				Type:        schema.TypeInt,
			},
			Attr_Networks: {
				Computed:    true,
				Description: "List of networks associated with this instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ExternalIP: {
							Computed:    true,
							Description: "The external IP address of the instance.",
							Type:        schema.TypeString,
						},
						Attr_IP: {
							Computed:    true,
							Description: "The IP address of the instance.",
							Type:        schema.TypeString,
						},
						Attr_MacAddress: {
							Computed:    true,
							Description: "The MAC address of the instance.",
							Type:        schema.TypeString,
						},
						Attr_NetworkID: {
							Computed:    true,
							Description: "The network ID of the instance.",
							Type:        schema.TypeString,
						},
						Attr_NetworkInterfaceID: {
							Computed:    true,
							Description: "ID of the network interface.",
							Type:        schema.TypeString,
						},
						Attr_NetworkName: {
							Computed:    true,
							Description: "The network name of the instance.",
							Type:        schema.TypeString,
						},
						Attr_NetworkSecurityGroupIDs: {
							Computed:    true,
							Description: "IDs of the network necurity groups that the network interface is a member of.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeSet,
						},
						Attr_NetworkSecurityGroupsHref: {
							Computed:    true,
							Description: "Links to the network security groups that the network interface is a member of.",
							Elem:        &schema.Schema{Type: schema.TypeString},
							Type:        schema.TypeList,
						},
						Attr_Type: {
							Computed:    true,
							Description: "The type of the network.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_PinPolicy: {
				Computed:    true,
				Description: "The pinning policy of the instance.",
				Type:        schema.TypeString,
			},
			Attr_PlacementGroupID: {
				Computed:    true,
				Description: "The ID of the placement group that the instance is a member.",
				Type:        schema.TypeString,
			},
			Attr_Processors: {
				Computed:    true,
				Description: "The number of processors that are allocated to the instance.",
				Type:        schema.TypeFloat,
			},
			Attr_ProcType: {
				Computed:    true,
				Description: "The procurement type of the instance. Supported values are shared and dedicated.",
				Type:        schema.TypeString,
			},
			Attr_ServerName: {
				Computed:    true,
				Description: "The name of the instance.",
				Type:        schema.TypeString,
			},
			Attr_SharedProcessorPool: {
				Computed:    true,
				Description: "The name of the shared processor pool for the instance.",
				Type:        schema.TypeString,
			},
			Attr_SharedProcessorPoolID: {
				Computed:    true,
				Description: "The ID of the shared processor pool for the instance.",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the instance.",
				Type:        schema.TypeString,
			},
			Attr_StorageConnection: {
				Computed:    true,
				Description: "The storage connection type.",
				Type:        schema.TypeString,
			},
			Attr_StoragePool: {
				Computed:    true,
				Description: "The storage Pool where server is deployed.",
				Type:        schema.TypeString,
			},
			Attr_StoragePoolAffinity: {
				Computed:    true,
				Description: "Indicates if all volumes attached to the server must reside in the same storage pool.",
				Type:        schema.TypeBool,
			},
			Attr_StorageType: {
				Computed:    true,
				Description: "The storage type where server is deployed.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
			Attr_VirtualCoresAssigned: {
				Computed:    true,
				Description: "The virtual cores that are assigned to the instance.",
				Type:        schema.TypeInt,
			},
			Attr_VirtualSerialNumber: {
				Computed:    true,
				Description: "Virtual Serial Number information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Computed:    true,
							Description: "Description of the Virtual Serial Number",
							Type:        schema.TypeString,
						},
						Attr_Serial: {
							Computed:    true,
							Description: "Virtual serial number.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Volumes: {
				Computed:    true,
				Description: "List of volume IDs that are attached to the instance.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := powerC.Get(d.Get(Arg_InstanceName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	pvminstanceid := *powervmdata.PvmInstanceID
	d.SetId(pvminstanceid)
	if powervmdata.Crn != "" {
		d.Set(Attr_CRN, powervmdata.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(powervmdata.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of pi instance (%s) user_tags: %s", *powervmdata.PvmInstanceID, err)
		}
		d.Set(Attr_UserTags, tags)
	}
	d.Set(Attr_DedicatedHostID, powervmdata.DedicatedHostID)
	d.Set(Attr_DeploymentType, powervmdata.DeploymentType)
	d.Set(Attr_LicenseRepositoryCapacity, powervmdata.LicenseRepositoryCapacity)
	d.Set(Attr_MaxMem, powervmdata.Maxmem)
	d.Set(Attr_MaxProc, powervmdata.Maxproc)
	d.Set(Attr_MaxVirtualCores, powervmdata.VirtualCores.Max)
	d.Set(Attr_Memory, powervmdata.Memory)
	d.Set(Attr_MinMem, powervmdata.Minmem)
	d.Set(Attr_MinProc, powervmdata.Minproc)
	d.Set(Attr_MinVirtualCores, powervmdata.VirtualCores.Min)
	d.Set(Attr_Networks, flattenPvmInstanceNetworks(powervmdata.Networks))
	d.Set(Attr_PinPolicy, powervmdata.PinPolicy)
	d.Set(Attr_Processors, powervmdata.Processors)
	d.Set(Attr_ProcType, powervmdata.ProcType)
	d.Set(Attr_ServerName, powervmdata.ServerName)
	d.Set(Attr_SharedProcessorPool, powervmdata.SharedProcessorPool)
	d.Set(Attr_SharedProcessorPoolID, powervmdata.SharedProcessorPoolID)
	d.Set(Attr_Status, powervmdata.Status)
	d.Set(Attr_StorageConnection, powervmdata.StorageConnection)
	d.Set(Attr_StorageType, powervmdata.StorageType)
	d.Set(Attr_StoragePool, powervmdata.StoragePool)
	d.Set(Attr_StoragePoolAffinity, powervmdata.StoragePoolAffinity)
	d.Set(Attr_Volumes, powervmdata.VolumeIDs)
	d.Set(Attr_VirtualCoresAssigned, powervmdata.VirtualCores.Assigned)

	if *powervmdata.PlacementGroup != "none" {
		d.Set(Attr_PlacementGroupID, powervmdata.PlacementGroup)
	}

	if powervmdata.Health != nil {
		d.Set(Attr_HealthStatus, powervmdata.Health.Status)
	}

	if powervmdata.SoftwareLicenses != nil {
		d.Set(Attr_IBMiCSS, powervmdata.SoftwareLicenses.IbmiCSS)
		d.Set(Attr_IBMiPHA, powervmdata.SoftwareLicenses.IbmiPHA)
		d.Set(Attr_IBMiRDS, powervmdata.SoftwareLicenses.IbmiRDS)
		if *powervmdata.SoftwareLicenses.IbmiRDS {
			d.Set(Attr_IBMiRDSUsers, powervmdata.SoftwareLicenses.IbmiRDSUsers)
		} else {
			d.Set(Attr_IBMiRDSUsers, 0)
		}
	}
	if powervmdata.Fault != nil {
		d.Set(Attr_Fault, flattenPvmInstanceFault(powervmdata.Fault))
	}
	if powervmdata.VirtualSerialNumber != nil {
		d.Set(Attr_VirtualSerialNumber, flattenVirtualSerialNumberToList(powervmdata.VirtualSerialNumber))
	}

	return nil
}
