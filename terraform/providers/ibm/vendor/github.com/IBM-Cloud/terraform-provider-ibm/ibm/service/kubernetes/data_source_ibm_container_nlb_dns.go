// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMContainerNLBDNS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMContainerNLBDNSRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name of the cluster",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_container_nlb_dns",
					"cluster"),
			},
			"nlb_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of nlb config of cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the secret.",
						},
						"secret_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of Secret.",
						},
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Id.",
						},
						"dns_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of DNS.",
						},
						"lb_hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host Name of load Balancer.",
						},
						"nlb_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: " NLB IPs.",
						},
						"nlb_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NLB Sub-Domain.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " Nlb Type.",
						},
						"secret_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace of Secret.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMContainerNLBDNSValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	iBMContainerNLBDNSValidator := validate.ResourceValidator{ResourceName: "ibm_container_nlb_dns", Schema: validateSchema}
	return &iBMContainerNLBDNSValidator
}
func dataSourceIBMContainerNLBDNSRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	name := d.Get("cluster").(string)

	kubeClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	nlbData, err := kubeClient.NlbDns().GetNLBDNSList(name)
	if err != nil || nlbData == nil || len(nlbData) < 1 {
		return diag.FromErr(fmt.Errorf("[ERROR] Error Listing NLB DNS (%s): %s", name, err))
	}
	d.SetId(name)
	d.Set("cluster", name)
	d.Set("nlb_config", flex.FlattenNlbConfigs(nlbData))
	return nil
}
