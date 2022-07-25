package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteClusterWorkerPoolAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSatelliteClusterWorkerPoolAttachmentRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name or id of the cluster",
			},
			"worker_pool": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "worker pool name",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "worker pool zone name",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the resource group that the Satellite location is in. To list the resource group ID of the location, use the `GET /v2/satellite/getController` API method.",
			},
			"worker_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of workers",
			},
			"autobalance_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Auto enabled status",
			},
			"messages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Filter features by a list of comma separated collections.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMSatelliteClusterWorkerPoolAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := d.Get("cluster").(string)
	workerpool := d.Get("worker_pool").(string)
	zone := d.Get("zone").(string)

	getWorkerPoolOptions := &kubernetesserviceapiv1.GetWorkerPoolOptions{}
	getWorkerPoolOptions.SetCluster(cluster)
	getWorkerPoolOptions.SetWorkerpool(workerpool)

	if v, ok := d.GetOk("resource_group_id"); ok {
		getWorkerPoolOptions.SetXAuthResourceGroup(v.(string))
	}

	getWorkerPoolResponse, response, err := satClient.GetWorkerPoolWithContext(context, getWorkerPoolOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("[DEBUG] Satellite cluster workerpool zone attachment record is not found %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("Satellite cluster workerpool zone attachment record is not found %s\n%s", err, response))
		}
		log.Printf("[DEBUG] GetWorkerPoolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetWorkerPoolWithContext failed %s\n%s", err, response))
	}

	if getWorkerPoolResponse != nil && getWorkerPoolResponse.Zones != nil {
		for _, z := range getWorkerPoolResponse.Zones {
			if zone == *z.ID {
				d.Set("worker_count", *z.WorkerCount)
				d.Set("autobalance_enabled", *z.AutobalanceEnabled)
				d.Set("messages", *&z.Messages)
				d.SetId(fmt.Sprintf("%s/%s/%s", cluster, workerpool, zone))
			}
		}
	}
	return nil
}
