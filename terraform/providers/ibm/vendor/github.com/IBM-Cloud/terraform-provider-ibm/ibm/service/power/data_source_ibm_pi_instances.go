// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesAllRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_PVMInstances: {
				Computed:    true,
				Description: "List of power virtual server instances for the respective cloud instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The minimum number of processors that must be allocated to the instance. ",
							Type:        schema.TypeFloat,
						},
						Attr_MinVirtualCores: {
							Computed:    true,
							Description: "The minimum number of virtual cores that can be assigned without rebooting the instance.",
							Type:        schema.TypeInt,
						},
						Attr_Networks: {
							Computed: true,
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
						Attr_PVMInstanceID: {
							Computed:    true,
							Description: "The unique identifier of the instance.",
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
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIInstancesAllRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()

	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)
	powervmdata, err := powerC.GetAll()

	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_PVMInstances, flattenPvmInstances(powervmdata.PvmInstances, meta))

	return nil
}

func flattenPvmInstances(list []*models.PVMInstanceReference, meta interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			Attr_DedicatedHostID:           i.DedicatedHostID,
			Attr_LicenseRepositoryCapacity: i.LicenseRepositoryCapacity,
			Attr_MaxMem:                    i.Maxmem,
			Attr_MaxProc:                   i.Maxproc,
			Attr_MaxVirtualCores:           i.VirtualCores.Max,
			Attr_Memory:                    *i.Memory,
			Attr_MinMem:                    i.Minmem,
			Attr_MinProc:                   i.Minproc,
			Attr_MinVirtualCores:           i.VirtualCores.Min,
			Attr_Networks:                  flattenPvmInstanceNetworks(i.Networks),
			Attr_PinPolicy:                 i.PinPolicy,
			Attr_PlacementGroupID:          i.PlacementGroup,
			Attr_Processors:                *i.Processors,
			Attr_ProcType:                  *i.ProcType,
			Attr_PVMInstanceID:             *i.PvmInstanceID,
			Attr_ServerName:                i.ServerName,
			Attr_SharedProcessorPool:       i.SharedProcessorPool,
			Attr_SharedProcessorPoolID:     i.SharedProcessorPoolID,
			Attr_Status:                    *i.Status,
			Attr_StorageConnection:         i.StorageConnection,
			Attr_StoragePool:               i.StoragePool,
			Attr_StoragePoolAffinity:       i.StoragePoolAffinity,
			Attr_StorageType:               i.StorageType,
			Attr_VirtualCoresAssigned:      i.VirtualCores.Assigned,
		}

		if i.Crn != "" {
			l[Attr_CRN] = i.Crn
			tags, err := flex.GetGlobalTagsUsingCRN(meta, string(i.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on get of pi instance (%s) user_tags: %s", *i.PvmInstanceID, err)
			}
			l[Attr_UserTags] = tags
		}

		if i.Health != nil {
			l[Attr_HealthStatus] = i.Health.Status
		}

		if i.Fault != nil {
			l[Attr_Fault] = flattenPvmInstanceFault(i.Fault)
		}

		if i.VirtualSerialNumber != nil {
			l[Attr_VirtualSerialNumber] = flattenVirtualSerialNumberToList(i.VirtualSerialNumber)
		}

		result = append(result, l)
	}
	return result
}

func flattenPvmInstanceNetworks(list []*models.PVMInstanceNetwork) (networks []map[string]interface{}) {
	if list != nil {
		networks = make([]map[string]interface{}, len(list))
		for i, pvmip := range list {
			p := make(map[string]interface{})
			p[Attr_ExternalIP] = pvmip.ExternalIP
			p[Attr_IP] = pvmip.IPAddress
			p[Attr_MacAddress] = pvmip.MacAddress
			p[Attr_NetworkID] = pvmip.NetworkID
			p[Attr_NetworkInterfaceID] = pvmip.NetworkInterfaceID
			p[Attr_NetworkName] = pvmip.NetworkName
			p[Attr_Type] = pvmip.Type
			if len(pvmip.NetworkSecurityGroupIDs) > 0 {
				p[Attr_NetworkSecurityGroupIDs] = pvmip.NetworkSecurityGroupIDs
			}
			if len(pvmip.NetworkSecurityGroupsHref) > 0 {
				p[Attr_NetworkSecurityGroupsHref] = pvmip.NetworkSecurityGroupsHref
			}

			networks[i] = p
		}
		return networks
	}
	return
}

func flattenPvmInstanceFault(fault *models.PVMInstanceFault) map[string]interface{} {
	faultMap := make(map[string]interface{})
	faultMap[Attr_Code] = strconv.FormatFloat(fault.Code, 'f', -1, 64)
	if !fault.Created.IsZero() {
		faultMap[Attr_Created] = fault.Created.String()
	}
	if fault.Details != "" {
		faultMap[Attr_Details] = fault.Details
	}
	if fault.Message != "" {
		faultMap[Attr_Message] = fault.Message
	}
	return faultMap
}
