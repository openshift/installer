package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchPoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateAzureRMBatchPoolName,
			},
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_tasks_per_node": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fixed_scale": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_dedicated_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_low_priority_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resize_timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"auto_scale": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"evaluation_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"formula": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"storage_image_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"offer": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"sku": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"node_agent_sku_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"store_location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CurrentUser",
								"LocalMachine",
							}, false),
						},
						"store_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"visibility": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"start_task": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_line": {
							Type:     schema.TypeString,
							Required: true,
						},

						"max_task_retry_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},

						"wait_for_success": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"environment": {
							Type:     schema.TypeMap,
							Optional: true,
						},

						"user_identity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_user": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"elevation_level": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"scope": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},

						"resource_file": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_storage_container_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"blob_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"http_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"storage_container_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceArmBatchPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchPoolClient

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch pool %q in account %q (Resource Group %q) was not found", name, accountName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch pool %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.PoolProperties; props != nil {
		d.Set("vm_size", props.VMSize)
		d.Set("max_tasks_per_node", props.MaxTasksPerNode)

		if scaleSettings := props.ScaleSettings; scaleSettings != nil {
			if err := d.Set("auto_scale", azure.FlattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
				return fmt.Errorf("Error flattening `auto_scale`: %+v", err)
			}
			if err := d.Set("fixed_scale", azure.FlattenBatchPoolFixedScaleSettings(scaleSettings.FixedScale)); err != nil {
				return fmt.Errorf("Error flattening `fixed_scale `: %+v", err)
			}
		}

		if props.DeploymentConfiguration != nil &&
			props.DeploymentConfiguration.VirtualMachineConfiguration != nil &&
			props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference != nil {

			imageReference := props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference

			d.Set("storage_image_reference", azure.FlattenBatchPoolImageReference(imageReference))
			d.Set("node_agent_sku_id", props.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
		}

		if err := d.Set("certificate", azure.FlattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
			return fmt.Errorf("error setting `certificate`: %v", err)
		}

		d.Set("start_task", azure.FlattenBatchPoolStartTask(props.StartTask))
	}

	return nil
}
