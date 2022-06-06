// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIbmSatelliteClusterWorkerPoolZoneAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSatelliteClusterWorkerPoolZoneAttachmentCreate,
		ReadContext:   resourceIbmSatelliteClusterWorkerPoolZoneAttachmentRead,
		DeleteContext: resourceIbmSatelliteClusterWorkerPoolZoneAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"worker_pool": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the resource group that the Satellite location is in. To list the resource group ID of the location, use the `GET /v2/satellite/getController` API method.",
			},
			"autobalance_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"messages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Filter features by a list of comma separated collections.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"worker_count": {
				Description: "Number of workers",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceIbmSatelliteClusterWorkerPoolZoneAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := d.Get("cluster").(string)
	workerpool := d.Get("worker_pool").(string)
	zone := d.Get("zone").(string)
	createSatelliteWorkerPoolZoneOptions := &kubernetesserviceapiv1.CreateSatelliteWorkerPoolZoneOptions{}
	createSatelliteWorkerPoolZoneOptions.SetCluster(cluster)
	createSatelliteWorkerPoolZoneOptions.SetWorkerpool(workerpool)
	createSatelliteWorkerPoolZoneOptions.SetID(zone)

	if _, ok := d.GetOk("resource_group_id"); ok {
		createSatelliteWorkerPoolZoneOptions.SetXAuthResourceGroup(d.Get("resource_group_id").(string))
	}

	response, err := satClient.CreateSatelliteWorkerPoolZone(createSatelliteWorkerPoolZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSatelliteWorkerPoolZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateSatelliteWorkerPoolZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", cluster, workerpool, zone))

	return resourceIbmSatelliteClusterWorkerPoolZoneAttachmentRead(context, d, meta)
}

func resourceIbmSatelliteClusterWorkerPoolZoneAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	getWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{}
	getWorkerPoolOptions.SetCluster(parts[0])
	getWorkerPoolOptions.SetWorkerpool(parts[1])

	getWorkerPoolResponse, response, err := satClient.GetWorkerPoolWithContext(context, getWorkerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetWorkerPool1WithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetWorkerPool1WithContext failed %s\n%s", err, response))
	}

	if getWorkerPoolResponse != nil && getWorkerPoolResponse.Zones != nil {
		d.Set("cluster", parts[0])
		d.Set("worker_pool", parts[1])

		for _, z := range getWorkerPoolResponse.Zones {
			if parts[2] == *z.ID {
				d.Set("zone", *z.ID)
				d.Set("autobalance_enabled", *z.AutobalanceEnabled)
				d.Set("messages", *&z.Messages)
				d.Set("worker_count", *z.WorkerCount)
			}
		}
	}

	return nil
}

func resourceIbmSatelliteClusterWorkerPoolZoneAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	removeWorkerPoolZoneOptions := &kubernetesserviceapiv1.RemoveWorkerPoolZoneOptions{}
	removeWorkerPoolZoneOptions.SetIdOrName(parts[0])
	removeWorkerPoolZoneOptions.SetPoolidOrName(parts[1])
	removeWorkerPoolZoneOptions.SetZoneid(parts[2])

	response, err := satClient.RemoveWorkerPoolZoneWithContext(context, removeWorkerPoolZoneOptions)
	if err != nil {
		log.Printf("[DEBUG] RemoveWorkerPoolZoneWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("RemoveWorkerPoolZoneWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
