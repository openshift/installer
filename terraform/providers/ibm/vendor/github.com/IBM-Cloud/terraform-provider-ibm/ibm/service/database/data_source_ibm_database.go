// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"fmt"
	"log"
	"net/url"

	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func DataSourceIBMDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDatabaseInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Resource instance name for example, my Database instance",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database",
					"name"),
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the resource group in which the Database instance is present",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database",
					"resource_group_id"),
			},

			"location": {
				Description: "The location or the region in which the Database instance exists",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database",
					"location"),
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of resource instance",
			},

			"service": {
				Description: "The name of the Cloud Database service",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"plan": {
				Description: "The plan type of the Database instance",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"status": {
				Description: "The resource instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"adminuser": {
				Description: "The admin user id for the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"adminpassword": {
				Description: "The admin user id for the instance",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"version": {
				Description: "The database version to provision if specified",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"platform_options": {
				Description: "Platform-specific options for this deployment.r",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_encryption_key_crn": {
							Description: "Disk encryption key crn",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"backup_encryption_key_crn": {
							Description: "Backup encryption key crn",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "User name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"password": {
							Description: "User password",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"allowlist": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description: "Allowlist IP address in CIDR notation",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Unique allowlist description",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Description: "Scaling group name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"count": {
							Description: "Count of scaling groups for the instance",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"memory": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The units memory is allocated in.",
									},
									"allocation_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current memory allocation for a group instance",
									},
									"minimum_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum memory size for a group instance",
									},
									"step_size_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The step size memory increases or decreases in.",
									},
									"is_adjustable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the memory size adjustable.",
									},
									"can_scale_down": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can memory scale down as well as up.",
									},
								},
							},
						},
						"cpu": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The .",
									},
									"allocation_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current cpu allocation count",
									},
									"minimum_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum number of cpus allowed",
									},
									"step_size_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of CPUs allowed to step up or down by",
									},
									"is_adjustable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Are the number of CPUs adjustable",
									},
									"can_scale_down": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can the number of CPUs be scaled down as well as up",
									},
								},
							},
						},
						"disk": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"units": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The units disk is allocated in",
									},
									"allocation_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current disk allocation",
									},
									"minimum_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The minimum disk size allowed",
									},
									"step_size_mb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The step size disk increases or decreases in",
									},
									"is_adjustable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the disk size adjustable",
									},
									"can_scale_down": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can the disk size be scaled down as well as up",
									},
								},
							},
						},
					},
				},
			},
			"auto_scaling": {
				Type:        schema.TypeList,
				Description: "ICD Auto Scaling",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk": {
							Type:        schema.TypeList,
							Description: "Disk Auto Scaling",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity_enabled": {
										Description: "Auto Scaling Scalar: Capacity Enabled",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"free_space_remaining_percent": {
										Description: "Auto Scaling Scalar: Capacity Free Space Remaining Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"free_space_less_than_percent": {
										Description: "Auto Scaling Scalar: Capacity Free Space Less Than Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"io_enabled": {
										Description: "Auto Scaling Scalar: IO Utilization Enabled",
										Type:        schema.TypeBool,
										Computed:    true,
									},

									"io_over_period": {
										Description: "Auto Scaling Scalar: IO Utilization Over Period",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"io_above_percent": {
										Description: "Auto Scaling Scalar: IO Utilization Above Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_mb_per_member": {
										Description: "Auto Scaling Rate: Limit mb per member",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_count_per_member": {
										Description: "Auto Scaling Rate: Limit count per number",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"memory": {
							Type:        schema.TypeList,
							Description: "Memory Auto Scaling",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"io_enabled": {
										Description: "Auto Scaling Scalar: IO Utilization Enabled",
										Type:        schema.TypeBool,
										Computed:    true,
									},

									"io_over_period": {
										Description: "Auto Scaling Scalar: IO Utilization Over Period",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"io_above_percent": {
										Description: "Auto Scaling Scalar: IO Utilization Above Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_mb_per_member": {
										Description: "Auto Scaling Rate: Limit mb per member",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_count_per_member": {
										Description: "Auto Scaling Rate: Limit count per number",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"cpu": {
							Type:        schema.TypeList,
							Description: "CPU Auto Scaling",
							Computed:    true,
							Deprecated:  "This field is deprecated, auto scaling cpu is unsupported by IBM Cloud Databases",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_mb_per_member": {
										Description: "Auto Scaling Rate: Limit mb per member",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_limit_count_per_member": {
										Description: "Auto Scaling Rate: Limit count per number",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"configuration_schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration schema in JSON format",
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}

func DataSourceIBMDatabaseInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cloud-database",
			CloudDataRange:             []string{"resolved_to:name"}})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_group",
			CloudDataRange:             []string{"resolved_to:id"},
			Optional:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "region",
			Optional:                   true})

	iBMDatabaseInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_database", Schema: validateSchema}
	return &iBMDatabaseInstanceValidator
}

func dataSourceIBMDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerAPIV2()
	if err != nil {
		return err
	}
	rsAPI := rsConClient.ResourceServiceInstanceV2()
	name := d.Get("name").(string)

	rsInstQuery := controllerv2.ServiceInstanceQuery{
		Name: name,
	}

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rsInstQuery.ResourceGroupID = rsGrpID.(string)
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInstQuery.ResourceGroupID = defaultRg
	}

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	if service, ok := d.GetOk("service"); ok {

		serviceOff, err := rsCatRepo.FindByName(service.(string), true)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving database offering: %s", err)
		}

		rsInstQuery.ServiceID = serviceOff[0].ID
	}

	var instances []models.ServiceInstanceV2

	instances, err = rsAPI.ListInstances(rsInstQuery)
	if err != nil {
		return err
	}
	var filteredInstances []models.ServiceInstanceV2
	var location string

	if loc, ok := d.GetOk("location"); ok {
		location = loc.(string)
		for _, instance := range instances {
			if flex.GetLocation(instance) == location {
				filteredInstances = append(filteredInstances, instance)
			}
		}
	} else {
		filteredInstances = instances
	}

	if len(filteredInstances) == 0 {
		return fmt.Errorf("[ERROR] No resource instance found with name [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or database", name)
	}

	var instance models.ServiceInstanceV2

	if len(filteredInstances) > 1 {
		return fmt.Errorf(
			"More than one resource instance found with name matching [%s]\nIf not specified please specify more filters like resource_group_id if instance doesn't exists in default group, location or database", name)
	}
	instance = filteredInstances[0]

	d.SetId(instance.ID)

	tags, err := flex.GetTagsUsingCRN(meta, d.Id())
	if err != nil {
		log.Printf(
			"Error on get of ibm Database tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	d.Set("name", instance.Name)
	d.Set("status", instance.State)
	d.Set("resource_group_id", instance.ResourceGroupID)
	d.Set("location", instance.RegionID)
	d.Set("guid", instance.Guid)

	serviceOff, err := rsCatRepo.GetServiceName(instance.ServiceID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	d.Set("service", serviceOff)

	servicePlan, err := rsCatRepo.GetServicePlanName(instance.ResourcePlanID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	d.Set("plan", servicePlan)

	d.Set(flex.ResourceName, instance.Name)
	d.Set(flex.ResourceCRN, instance.Crn.String())
	d.Set(flex.ResourceStatus, instance.State)
	d.Set(flex.ResourceGroupName, instance.ResourceGroupName)

	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/services/"+url.QueryEscape(instance.Crn.String()))

	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	icdClient, err := meta.(conns.ClientSession).ICDAPI()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	icdId := flex.EscapeUrlParm(instance.ID)
	getDeploymentInfoOptions := &clouddatabasesv5.GetDeploymentInfoOptions{
		ID: core.StringPtr(instance.ID),
	}
	getDeploymentInfoResponse, response, err := cloudDatabasesClient.GetDeploymentInfo(getDeploymentInfoOptions)
	if err != nil {
		if response.StatusCode == 404 {
			return fmt.Errorf("[ERROR] The database instance was not found in the region set for the Provider, or the default of us-south. Specify the correct region in the provider definition, or create a provider alias for the correct region. %v", err)
		}
		return fmt.Errorf("[ERROR] Error getting database config while updating adminpassword for: %s with error %s", instance.ID, err)
	}

	deployment := getDeploymentInfoResponse.Deployment
	adminUser := deployment.AdminUsernames["database"]

	d.Set("adminuser", adminUser)
	d.Set("version", deployment.Version)

	if deployment.PlatformOptions != nil {
		d.Set("platform_options", flex.ExpandPlatformOptions(*deployment))
	}

	groupList, err := icdClient.Groups().GetGroups(icdId)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database groups: %s", err)
	}
	d.Set("groups", flex.FlattenIcdGroups(groupList))

	getAutoscalingConditionsOptions := &clouddatabasesv5.GetAutoscalingConditionsOptions{
		ID:      &instance.ID,
		GroupID: core.StringPtr("member"),
	}

	autoscalingGroup, _, err := cloudDatabasesClient.GetAutoscalingConditions(getAutoscalingConditionsOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database autoscaling groups: %s", err)
	}
	d.Set("auto_scaling", flattenAutoScalingGroup(*autoscalingGroup))

	alEntry := &clouddatabasesv5.GetAllowlistOptions{
		ID: &instance.ID,
	}

	allowlist, _, err := cloudDatabasesClient.GetAllowlist(alEntry)

	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database allowlist: %s", err)
	}

	d.Set("allowlist", flex.FlattenAllowlist(allowlist.IPAddresses))

	return nil
}
