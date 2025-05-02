// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	//	"github.com/IBM-Cloud/bluemix-go/api/globaltagging/globaltaggingv3"
	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
	"github.com/IBM/go-sdk-core/v5/core"
)

const (
	databaseInstanceSuccessStatus      = "active"
	databaseInstanceProvisioningStatus = "provisioning"
	databaseInstanceProgressStatus     = "in progress"
	databaseInstanceInactiveStatus     = "inactive"
	databaseInstanceFailStatus         = "failed"
	databaseInstanceRemovedStatus      = "removed"
	databaseInstanceReclamation        = "pending_reclamation"
)

const (
	databaseTaskSuccessStatus  = "completed"
	databaseTaskProgressStatus = "running"
	databaseTaskFailStatus     = "failed"
)

type userChange struct {
	Old, New map[string]interface{}
}

func retry(f func() error) (err error) {
	attempts := 3

	for i := 0; ; i++ {
		sleep := time.Duration(10*i*i) * time.Second
		time.Sleep(sleep)

		err = f()
		if err == nil {
			return nil
		}

		if i == attempts {
			return err
		}

		log.Println("retrying after error:", err)
	}
}

func ResourceIBMDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMDatabaseInstanceCreate,
		ReadContext:   resourceIBMDatabaseInstanceRead,
		UpdateContext: resourceIBMDatabaseInstanceUpdate,
		DeleteContext: resourceIBMDatabaseInstanceDelete,
		Exists:        resourceIBMDatabaseInstanceExists,

		CustomizeDiff: customdiff.All(
			resourceIBMDatabaseInstanceDiff,
			checkV5Groups),

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Resource instance name for example, my Database instance",
				Type:        schema.TypeString,
				Required:    true,
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The id of the resource group in which the Database instance is present",
				ValidateFunc: validate.InvokeValidator(
					"ibm_database",
					"resource_group_id"),
			},

			"location": {
				Description: "The location or the region in which Database instance exists",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeValidator(
					"ibm_database",
					"location"),
			},

			"service": {
				Description:  "The name of the Cloud Internet database service",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_database", "service"),
			},
			"plan": {
				Description:  "The plan type of the Database instance",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_database", "plan"),
				ForceNew:     true,
			},

			"status": {
				Description: "The resource instance status",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of resource instance",
			},

			"adminuser": {
				Description: "The admin user id for the instance",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"adminpassword": {
				Description:  "The admin user password for the instance",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(10, 32),
				Sensitive:    true,
				// DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				//  return true
				// },
			},
			"configuration": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					json, err := flex.NormalizeJSONString(v)
					if err != nil {
						return fmt.Sprintf("%q", err.Error())
					}
					return json
				},
				Description: "The configuration in JSON format",
			},
			"configuration_schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration schema in JSON format",
			},
			"version": {
				Description: "The database version to provision if specified",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"members_memory_allocation_mb": {
				Description:   "Memory allocation required for cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"members_disk_allocation_mb": {
				Description:   "Disk allocation required for cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"members_cpu_allocation_count": {
				Description:   "CPU allocation required for cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"node_count": {
				Description:   "Total number of nodes in the cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"node_memory_allocation_mb": {
				Description:   "Memory allocation per node",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"node_disk_allocation_mb": {
				Description:   "Disk allocation per node",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"node_cpu_allocation_count": {
				Description:   "CPU allocation per node",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count", "group"},
				Deprecated:    "use group instead",
			},
			"plan_validation": {
				Description: "For elasticsearch and postgres perform database parameter validation during the plan phase. Otherwise, database parameter validation happens in apply phase.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return true
					}
					return false
				},
			},
			"service_endpoints": {
				Description:  "Types of the service endpoints. Possible values are 'public', 'private', 'public-and-private'.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public",
				ValidateFunc: validate.InvokeValidator("ibm_database", "service_endpoints"),
			},
			"backup_id": {
				Description: "The CRN of backup source database",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"remote_leader_id": {
				Description:      "The CRN of leader database",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
			},
			"key_protect_instance": {
				Description: "The CRN of Key protect instance",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"key_protect_key": {
				Description: "The CRN of Key protect key",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"backup_encryption_key_crn": {
				Description: "The Backup Encryption Key CRN",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_database", "tags")},
				Set:      flex.ResourceIBMVPCHash,
			},
			"point_in_time_recovery_deployment_id": {
				Description:      "The CRN of source instance",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
			},
			"point_in_time_recovery_time": {
				Description:      "The point in time recovery time stamp of the deployed instance",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.ApplyOnce,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "User name",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(4, 32),
						},
						"password": {
							Description:  "User password",
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(10, 32),
						},
						"type": {
							Description:  "User type",
							Type:         schema.TypeString,
							Default:      "database",
							Optional:     true,
							Sensitive:    false,
							ValidateFunc: validation.StringInSlice([]string{"database", "ops_manager", "read_only_replica"}, false),
						},
						"role": {
							Description:  "User role. Only available for ops_manager user type.",
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    false,
							ValidateFunc: validation.StringInSlice([]string{"group_read_only", "group_data_access_admin"}, false),
						},
					},
				},
			},
			"connectionstrings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "User name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"composed": {
							Description: "Connection string",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"scheme": {
							Description: "DB scheme",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"certname": {
							Description: "Certificate Name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"certbase64": {
							Description: "Certificate in base64 encoding",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"bundlename": {
							Description: "Cassandra Bundle Name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"bundlebase64": {
							Description: "Cassandra base64 encoding",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"password": {
							Description: "Password",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"queryoptions": {
							Description: "DB query options",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"database": {
							Description: "DB name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path": {
							Description: "DB path",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hosts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": {
										Description: "DB host name",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"port": {
										Description: "DB port",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
				Deprecated: "This field is deprecated, please use ibm_database_connection instead",
			},
			"whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description:  "Whitelist IP address in CIDR notation",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.ValidateCIDR,
						},
						"description": {
							Description:  "Unique white list description",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
						},
					},
				},
				Deprecated:    "Whitelist is deprecated please use allowlist",
				ConflictsWith: []string{"allowlist"},
			},
			"allowlist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description:  "Allowlist IP address in CIDR notation",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.ValidateCIDR,
						},
						"description": {
							Description:  "Unique allow list description",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
						},
					},
				},
				ConflictsWith: []string{"whitelist"},
			},
			"logical_replication_slot": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Logical Replication Slot name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"database_name": {
							Description: "Database Name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"plugin_type": {
							Description: "Plugin Type",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"group": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count", "node_count"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validate.InvokeValidator("ibm_database", "group_id"),
						},
						"members": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_count": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"memory": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_mb": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"disk": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_mb": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"cpu": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_count": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
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
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk": {
							Type:        schema.TypeList,
							Description: "Disk Auto Scaling",
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capacity_enabled": {
										Description: "Auto Scaling Scalar: Capacity Enabled",
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
									},
									"free_space_less_than_percent": {
										Description: "Auto Scaling Scalar: Capacity Free Space Less Than Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"io_enabled": {
										Description: "Auto Scaling Scalar: IO Utilization Enabled",
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
									},

									"io_over_period": {
										Description: "Auto Scaling Scalar: IO Utilization Over Period",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"io_above_percent": {
										Description: "Auto Scaling Scalar: IO Utilization Above Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_limit_mb_per_member": {
										Description: "Auto Scaling Rate: Limit mb per member",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
						"memory": {
							Type:        schema.TypeList,
							Description: "Memory Auto Scaling",
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"io_enabled": {
										Description: "Auto Scaling Scalar: IO Utilization Enabled",
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
									},

									"io_over_period": {
										Description: "Auto Scaling Scalar: IO Utilization Over Period",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
									"io_above_percent": {
										Description: "Auto Scaling Scalar: IO Utilization Above Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_limit_mb_per_member": {
										Description: "Auto Scaling Rate: Limit mb per member",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
						"cpu": {
							Type:        schema.TypeList,
							Description: "CPU Auto Scaling",
							Deprecated:  "This field is deprecated, auto scaling cpu is unsupported by IBM Cloud Databases",
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rate_increase_percent": {
										Description: "Auto Scaling Rate: Increase Percent",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_period_seconds": {
										Description: "Auto Scaling Rate: Period Seconds",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_limit_count_per_member": {
										Description: "Auto Scaling Rate: Limit count per number",
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
									},
									"rate_units": {
										Description: "Auto Scaling Rate: Units ",
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
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
func ResourceIBMICDValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
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
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "databases-for-etcd, databases-for-postgresql, databases-for-redis, databases-for-elasticsearch, databases-for-mongodb, messages-for-rabbitmq, databases-for-mysql, databases-for-cassandra, databases-for-enterprisedb",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "plan",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "standard, enterprise",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_endpoints",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "public, private, public-and-private",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "group_id",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "member, analytics, bi_connector, search",
			Required:                   true})

	ibmICDResourceValidator := validate.ResourceValidator{ResourceName: "ibm_database", Schema: validateSchema}
	return &ibmICDResourceValidator
}

type Params struct {
	Version             string  `json:"version,omitempty"`
	KeyProtectKey       string  `json:"disk_encryption_key_crn,omitempty"`
	BackUpEncryptionCRN string  `json:"backup_encryption_key_crn,omitempty"`
	Memory              int     `json:"members_memory_allocation_mb,omitempty"`
	Disk                int     `json:"members_disk_allocation_mb,omitempty"`
	CPU                 int     `json:"members_cpu_allocation_count,omitempty"`
	KeyProtectInstance  string  `json:"disk_encryption_instance_crn,omitempty"`
	ServiceEndpoints    string  `json:"service-endpoints,omitempty"`
	BackupID            string  `json:"backup-id,omitempty"`
	RemoteLeaderID      string  `json:"remote_leader_id,omitempty"`
	PITRDeploymentID    string  `json:"point_in_time_recovery_deployment_id,omitempty"`
	PITRTimeStamp       *string `json:"point_in_time_recovery_time,omitempty"`
}

type Group struct {
	ID      string
	Members *GroupResource
	Memory  *GroupResource
	Disk    *GroupResource
	CPU     *GroupResource
}

type GroupResource struct {
	Units        string
	Allocation   int
	Minimum      int
	Maximum      int
	StepSize     int
	IsAdjustable bool
	IsOptional   bool
	CanScaleDown bool
}

func getDefaultScalingGroups(_service string, _plan string, meta interface{}) (groups []clouddatabasesv5.Group, err error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return groups, fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	re := regexp.MustCompile("(?:messages|databases)-for-([a-z]+)")
	match := re.FindStringSubmatch(_service)

	if match == nil {
		return groups, fmt.Errorf("[ERROR] Error invalid service name: %s", _service)
	}

	service := match[1]

	if service == "cassandra" {
		service = "datastax_enterprise_full"
	}

	if service == "mongodb" && _plan == "enterprise" {
		service = "mongodbee"
	}

	getDefaultScalingGroupsOptions := cloudDatabasesClient.NewGetDefaultScalingGroupsOptions(service)

	getDefaultScalingGroupsResponse, _, err := cloudDatabasesClient.GetDefaultScalingGroups(getDefaultScalingGroupsOptions)
	if err != nil {
		return groups, err
	}

	return getDefaultScalingGroupsResponse.Groups, nil
}

func getDatabaseServiceDefaults(service string, meta interface{}) (*icdv4.Group, error) {
	icdClient, err := meta.(conns.ClientSession).ICDAPI()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	var dbType string
	if service == "databases-for-cassandra" {
		dbType = "datastax_enterprise_full"
	} else if strings.HasPrefix(service, "messages-for-") {
		dbType = service[len("messages-for-"):]
	} else {
		dbType = service[len("databases-for-"):]
	}

	groupDefaults, err := icdClient.Groups().GetDefaultGroups(dbType)
	if err != nil {
		return nil, fmt.Errorf("ICD API is down for plan validation, set plan_validation=false %s", err)
	}
	return &groupDefaults.Groups[0], nil
}

func getInitialNodeCount(service string, plan string, meta interface{}) (int, error) {
	groups, err := getDefaultScalingGroups(service, plan, meta)

	if err != nil {
		return 0, err
	}

	for _, g := range groups {
		if *g.ID == "member" {
			return int(*g.Members.MinimumCount), nil
		}
	}

	return 0, fmt.Errorf("getInitialNodeCount failed for member group")
}

func getGroups(instanceID string, meta interface{}) (groups []clouddatabasesv5.Group, err error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return nil, err
	}

	listDeploymentScalingGroupsOptions := &clouddatabasesv5.ListDeploymentScalingGroupsOptions{
		ID: &instanceID,
	}

	groupsResponse, _, err := cloudDatabasesClient.ListDeploymentScalingGroups(listDeploymentScalingGroupsOptions)
	if err != nil {
		return groups, err
	}

	return groupsResponse.Groups, nil
}

// V5 Groups
func checkGroupScaling(groupId string, resourceName string, value int, resource *GroupResource, nodeCount int) error {
	if nodeCount == 0 {
		nodeCount = 1
	}
	if resource.StepSize == 0 {
		return fmt.Errorf("%s group must have members scaled > 0 before scaling %s", groupId, resourceName)
	}
	if value < resource.Minimum/nodeCount || value > resource.Maximum/nodeCount || value%(resource.StepSize/nodeCount) != 0 {
		if !(value == 0 && resource.IsOptional) {
			return fmt.Errorf("%s group %s must be >= %d and <= %d in increments of %d", groupId, resourceName, resource.Minimum/nodeCount, resource.Maximum/nodeCount, resource.StepSize/nodeCount)
		}
	}
	if value != resource.Allocation/nodeCount && !resource.IsAdjustable {
		return fmt.Errorf("%s can not change %s value after create", groupId, resourceName)
	}
	if value < resource.Allocation/nodeCount && !resource.CanScaleDown {
		return fmt.Errorf("can not scale %s group %s below %d to %d", groupId, resourceName, resource.Allocation/nodeCount, value)
	}
	return nil
}

func checkGroupValue(name string, limits GroupResource, divider int, diff *schema.ResourceDiff) error {
	if diff.HasChange(name) {
		oldSetting, newSetting := diff.GetChange(name)
		old := oldSetting.(int)
		new := newSetting.(int)

		if new < limits.Minimum/divider || new > limits.Maximum/divider || new%(limits.StepSize/divider) != 0 {
			if !(new == 0 && limits.IsOptional) {
				return fmt.Errorf("%s must be >= %d and <= %d in increments of %d", name, limits.Minimum/divider, limits.Maximum/divider/divider, limits.StepSize/divider)
			}
		}
		if old != new && !limits.IsAdjustable {
			return fmt.Errorf("%s can not change value after create", name)
		}
		if new < old && !limits.CanScaleDown {
			return fmt.Errorf("%s can not scale down from %d to %d", name, old, new)
		}
		return nil
	}
	return nil
}

type CountLimit struct {
	Units           string
	AllocationCount int
	MinimumCount    int
	MaximumCount    int
	StepSizeCount   int
	IsAdjustable    bool
	IsOptional      bool
	CanScaleDown    bool
}

func checkCountValue(name string, limits CountLimit, divider int, diff *schema.ResourceDiff) error {
	groupLimit := GroupResource{
		Units:        limits.Units,
		Allocation:   limits.AllocationCount,
		Minimum:      limits.MinimumCount,
		Maximum:      limits.MaximumCount,
		StepSize:     limits.StepSizeCount,
		IsAdjustable: limits.IsAdjustable,
		IsOptional:   limits.IsOptional,
		CanScaleDown: limits.CanScaleDown,
	}
	return checkGroupValue(name, groupLimit, divider, diff)
}

type MbLimit struct {
	Units        string
	AllocationMb int
	MinimumMb    int
	MaximumMb    int
	StepSizeMb   int
	IsAdjustable bool
	IsOptional   bool
	CanScaleDown bool
}

func checkMbValue(name string, limits MbLimit, divider int, diff *schema.ResourceDiff) error {
	groupLimit := GroupResource{
		Units:        limits.Units,
		Allocation:   limits.AllocationMb,
		Minimum:      limits.MinimumMb,
		Maximum:      limits.MaximumMb,
		StepSize:     limits.StepSizeMb,
		IsAdjustable: limits.IsAdjustable,
		IsOptional:   limits.IsOptional,
		CanScaleDown: limits.CanScaleDown,
	}
	return checkGroupValue(name, groupLimit, divider, diff)
}

func resourceIBMDatabaseInstanceDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	err = flex.ResourceTagsCustomizeDiff(diff)
	if err != nil {
		return err
	}

	service := diff.Get("service").(string)
	planPhase := diff.Get("plan_validation").(bool)

	if service == "databases-for-postgresql" ||
		service == "databases-for-elasticsearch" ||
		service == "databases-for-cassandra" ||
		service == "databases-for-enterprisedb" {
		if planPhase {
			groupDefaults, err := getDatabaseServiceDefaults(service, meta)
			if err != nil {
				return err
			}

			err = checkMbValue("members_memory_allocation_mb", MbLimit(groupDefaults.Memory), 1, diff)
			if err != nil {
				return err
			}

			err = checkMbValue("members_disk_allocation_mb", MbLimit(groupDefaults.Disk), 1, diff)
			if err != nil {
				return err
			}

			err = checkCountValue("members_cpu_allocation_count", CountLimit(groupDefaults.Cpu), 1, diff)
			if err != nil {
				return err
			}

			err = checkCountValue("node_count", CountLimit(groupDefaults.Members), 1, diff)
			if err != nil {
				return err
			}

			var divider = groupDefaults.Members.MinimumCount
			err = checkMbValue("node_memory_allocation_mb", MbLimit(groupDefaults.Memory), divider, diff)
			if err != nil {
				return err
			}

			err = checkMbValue("node_disk_allocation_mb", MbLimit(groupDefaults.Disk), divider, diff)
			if err != nil {
				return err
			}

			if diff.HasChange("node_cpu_allocation_count") {
				err = checkCountValue("node_cpu_allocation_count", CountLimit(groupDefaults.Cpu), divider, diff)
				if err != nil {
					return err
				}
			} else if diff.HasChange("node_count") {
				if _, ok := diff.GetOk("node_cpu_allocation_count"); !ok {
					_, newSetting := diff.GetChange("node_count")
					min := groupDefaults.Cpu.MinimumCount / divider
					if newSetting != min {
						return fmt.Errorf("node_cpu_allocation_count must be set when node_count is greater then the minimum %d", min)
					}
				}
			}
		}
	} else if diff.HasChange("node_count") || diff.HasChange("node_memory_allocation_mb") || diff.HasChange("node_disk_allocation_mb") || diff.HasChange("node_cpu_allocation_count") {
		return fmt.Errorf("[ERROR] node_count, node_memory_allocation_mb, node_disk_allocation_mb, node_cpu_allocation_count only supported for postgresql, elasticsearch and cassandra")
	}

	_, logicalReplicationSet := diff.GetOk("logical_replication_slot")

	if service != "databases-for-postgresql" && logicalReplicationSet {
		return fmt.Errorf("[ERROR] logical_replication_slot is only supported for databases-for-postgresql")
	}

	configJSON, configOk := diff.GetOk("configuration")

	if configOk {
		var rawConfig map[string]json.RawMessage
		err = json.Unmarshal([]byte(configJSON.(string)), &rawConfig)
		if err != nil {
			return fmt.Errorf("[ERROR] configuration JSON invalid\n%s", err)
		}

		var unmarshalFn func(m map[string]json.RawMessage, result interface{}) (err error)

		var configuration clouddatabasesv5.ConfigurationIntf = new(clouddatabasesv5.Configuration)

		switch service {
		case "databases-for-postgresql":
			unmarshalFn = clouddatabasesv5.UnmarshalConfigurationPgConfiguration
		case "databases-for-enterprisedb":
			unmarshalFn = clouddatabasesv5.UnmarshalConfigurationPgConfiguration
		case "databases-for-redis":
			unmarshalFn = clouddatabasesv5.UnmarshalConfigurationRedisConfiguration
		case "databases-for-mysql":
			unmarshalFn = clouddatabasesv5.UnmarshalConfigurationMySQLConfiguration
		case "messages-for-rabbitmq":
			unmarshalFn = clouddatabasesv5.UnmarshalConfigurationRabbitMqConfiguration
		default:
			return fmt.Errorf("[ERROR] configuration is not supported for %s", service)
		}

		err = core.UnmarshalModel(rawConfig, "", &configuration, unmarshalFn)
		if err != nil {
			return fmt.Errorf("[ERROR] configuration is invalid\n%s", err)
		}

		b, _ := json.Marshal(configuration)
		var result map[string]json.RawMessage
		json.Unmarshal(b, &result)

		invalidFields := []string{}
		for k, _ := range rawConfig {
			if _, ok := result[k]; !ok {
				invalidFields = append(invalidFields, k)
			}
		}

		if len(invalidFields) != 0 {
			return fmt.Errorf("[ERROR] configuration contained invalid field(s): %s", invalidFields)
		}
	}

	return nil
}

// Replace with func wrapper for resourceIBMResourceInstanceCreate specifying serviceName := "database......."
func resourceIBMDatabaseInstanceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	rsInst := rc.CreateResourceInstanceOptions{
		Name: &name,
	}

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving database service offering: %s", err))
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving plan: %s", err))
	}
	rsInst.ResourcePlanID = &servicePlan

	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err))
	}
	if len(deployments) == 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] No deployment found for service plan : %s", plan))
	}
	deployments, supportedLocations := filterDatabaseDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return diag.FromErr(fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\nValid location(s) are: %q", plan, location, locationList))
	}
	catalogCRN := deployments[0].CatalogCRN
	rsInst.Target = &catalogCRN

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rgID := rsGrpID.(string)
		rsInst.ResourceGroup = &rgID
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return diag.FromErr(err)
		}
		rsInst.ResourceGroup = &defaultRg
	}

	initialNodeCount, err := getInitialNodeCount(serviceName, plan, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	params := Params{}
	if group, ok := d.GetOk("group"); ok {
		groups := expandGroups(group.(*schema.Set).List())
		var memberGroup *Group
		for _, g := range groups {
			if g.ID == "member" {
				memberGroup = g
				break
			}
		}

		if memberGroup != nil {
			if memberGroup.Memory != nil {
				params.Memory = memberGroup.Memory.Allocation * initialNodeCount
			}

			if memberGroup.Disk != nil {
				params.Disk = memberGroup.Disk.Allocation * initialNodeCount
			}

			if memberGroup.CPU != nil {
				params.CPU = memberGroup.CPU.Allocation * initialNodeCount
			}
		}
	}
	if memory, ok := d.GetOk("members_memory_allocation_mb"); ok {
		params.Memory = memory.(int)
	}
	if memory, ok := d.GetOk("node_memory_allocation_mb"); ok {
		params.Memory = memory.(int) * initialNodeCount
	}
	if disk, ok := d.GetOk("members_disk_allocation_mb"); ok {
		params.Disk = disk.(int)
	}
	if disk, ok := d.GetOk("node_disk_allocation_mb"); ok {
		params.Disk = disk.(int) * initialNodeCount
	}
	if cpu, ok := d.GetOk("members_cpu_allocation_count"); ok {
		params.CPU = cpu.(int)
	}
	if cpu, ok := d.GetOk("node_cpu_allocation_count"); ok {
		params.CPU = cpu.(int) * initialNodeCount
	}
	if version, ok := d.GetOk("version"); ok {
		params.Version = version.(string)
	}
	if keyProtect, ok := d.GetOk("key_protect_key"); ok {
		params.KeyProtectKey = keyProtect.(string)
	}
	if keyProtectInstance, ok := d.GetOk("key_protect_instance"); ok {
		params.KeyProtectInstance = keyProtectInstance.(string)
	}
	if backupID, ok := d.GetOk("backup_id"); ok {
		params.BackupID = backupID.(string)
	}
	if backUpEncryptionKey, ok := d.GetOk("backup_encryption_key_crn"); ok {
		params.BackUpEncryptionCRN = backUpEncryptionKey.(string)
	}
	if remoteLeader, ok := d.GetOk("remote_leader_id"); ok {
		params.RemoteLeaderID = remoteLeader.(string)
	}

	if pitrID, ok := d.GetOk("point_in_time_recovery_deployment_id"); ok {
		params.PITRDeploymentID = pitrID.(string)
	}

	pitrOk := !d.GetRawConfig().AsValueMap()["point_in_time_recovery_time"].IsNull()
	if pitrTime, ok := d.GetOk("point_in_time_recovery_time"); pitrOk {
		if !ok {
			pitrTime = ""
		}

		pitrTimeTrimmed := strings.TrimSpace(pitrTime.(string))
		params.PITRTimeStamp = &pitrTimeTrimmed
	}

	serviceEndpoint := d.Get("service_endpoints").(string)
	params.ServiceEndpoints = serviceEndpoint
	parameters, _ := json.Marshal(params)
	var raw map[string]interface{}
	json.Unmarshal(parameters, &raw)
	//paramString := string(parameters[:])
	rsInst.Parameters = raw

	instance, response, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil {
		return diag.FromErr(
			fmt.Errorf("[ERROR] Error creating database instance: %s %s", err, response))
	}
	d.SetId(*instance.ID)

	_, err = waitForDatabaseInstanceCreate(d, meta, *instance.ID)
	if err != nil {
		return diag.FromErr(
			fmt.Errorf(
				"[ERROR] Error waiting for create database instance (%s) to complete: %s", *instance.ID, err))
	}

	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	if node_count, ok := d.GetOk("node_count"); ok {
		if initialNodeCount != node_count {
			icdClient, err := meta.(conns.ClientSession).ICDAPI()
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
			}

			err = horizontalScale(d, meta, icdClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if group, ok := d.GetOk("group"); ok {
		groups := expandGroups(group.(*schema.Set).List())
		groupsResponse, err := getGroups(*instance.ID, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		currentGroups := normalizeGroups(groupsResponse)

		for _, g := range groups {
			groupScaling := &clouddatabasesv5.GroupScaling{}
			var currentGroup *Group
			var nodeCount int

			for _, cg := range currentGroups {
				if cg.ID == g.ID {
					currentGroup = &cg
					nodeCount = currentGroup.Members.Allocation
				}
			}

			if g.ID == "member" && (g.Members == nil || g.Members.Allocation == nodeCount) {
				// No Horizontal Scaling needed
				continue
			}

			if g.Members != nil && g.Members.Allocation != currentGroup.Members.Allocation {
				groupScaling.Members = &clouddatabasesv5.GroupScalingMembers{AllocationCount: core.Int64Ptr(int64(g.Members.Allocation))}
				nodeCount = g.Members.Allocation
			}
			if g.Memory != nil && g.Memory.Allocation*nodeCount != currentGroup.Memory.Allocation {
				groupScaling.Memory = &clouddatabasesv5.GroupScalingMemory{AllocationMb: core.Int64Ptr(int64(g.Memory.Allocation * nodeCount))}
			}
			if g.Disk != nil && g.Disk.Allocation*nodeCount != currentGroup.Disk.Allocation {
				groupScaling.Disk = &clouddatabasesv5.GroupScalingDisk{AllocationMb: core.Int64Ptr(int64(g.Disk.Allocation * nodeCount))}
			}
			if g.CPU != nil && g.CPU.Allocation*nodeCount != currentGroup.CPU.Allocation {
				groupScaling.CPU = &clouddatabasesv5.GroupScalingCPU{AllocationCount: core.Int64Ptr(int64(g.CPU.Allocation * nodeCount))}
			}

			if groupScaling.Members != nil || groupScaling.Memory != nil || groupScaling.Disk != nil || groupScaling.CPU != nil {
				setDeploymentScalingGroupOptions := &clouddatabasesv5.SetDeploymentScalingGroupOptions{
					ID:      instance.ID,
					GroupID: &g.ID,
					Group:   groupScaling,
				}

				setDeploymentScalingGroupResponse, _, err := cloudDatabasesClient.SetDeploymentScalingGroup(setDeploymentScalingGroupOptions)

				taskIDLink := *setDeploymentScalingGroupResponse.Task.ID

				_, err = waitForDatabaseTaskComplete(taskIDLink, d, meta, d.Timeout(schema.TimeoutCreate))

				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of ibm database (%s) tags: %s", d.Id(), err)
		}
	}

	instanceID := *instance.ID
	icdId := flex.EscapeUrlParm(instanceID)

	if pw, ok := d.GetOk("adminpassword"); ok {
		adminPassword := pw.(string)

		getDeploymentInfoOptions := &clouddatabasesv5.GetDeploymentInfoOptions{
			ID: core.StringPtr(instanceID),
		}
		getDeploymentInfoResponse, response, err := cloudDatabasesClient.GetDeploymentInfo(getDeploymentInfoOptions)

		if err != nil {
			if response.StatusCode == 404 {
				return diag.FromErr(fmt.Errorf("[ERROR] The database instance was not found in the region set for the Provider, or the default of us-south. Specify the correct region in the provider definition, or create a provider alias for the correct region. %v", err))
			}
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting database config while updating adminpassword for: %s with error %s", instanceID, err))
		}
		deployment := getDeploymentInfoResponse.Deployment

		adminUser := deployment.AdminUsernames["database"]

		user := &clouddatabasesv5.APasswordSettingUser{
			Password: &adminPassword,
		}

		changeUserPasswordOptions := &clouddatabasesv5.ChangeUserPasswordOptions{
			ID:       core.StringPtr(instanceID),
			UserType: core.StringPtr("database"),
			Username: core.StringPtr(adminUser),
			User:     user,
		}

		changeUserPasswordResponse, response, err := cloudDatabasesClient.ChangeUserPassword(changeUserPasswordOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] ChangeUserPassword (%s) failed %s\n%s", *changeUserPasswordOptions.Username, err, response))
		}

		taskID := *changeUserPasswordResponse.Task.ID
		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutCreate))

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database admin password: %s", err))
		}
	}

	_, hasWhitelist := d.GetOk("whitelist")
	_, hasAllowlist := d.GetOk("allowlist")

	if hasWhitelist || hasAllowlist {
		var ipAddresses *schema.Set
		if hasWhitelist {
			ipAddresses = d.Get("whitelist").(*schema.Set)
		} else {
			ipAddresses = d.Get("allowlist").(*schema.Set)
		}

		entries := flex.ExpandAllowlist(ipAddresses)

		setAllowlistOptions := &clouddatabasesv5.SetAllowlistOptions{
			ID:          &instanceID,
			IPAddresses: entries,
		}

		setAllowlistResponse, _, err := cloudDatabasesClient.SetAllowlist(setAllowlistOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database allowlists: %s", err))
		}

		taskId := *setAllowlistResponse.Task.ID

		_, err = waitForDatabaseTaskComplete(taskId, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for update of database (%s) allowlist task to complete: %s", instanceID, err))
		}
	}

	if _, ok := d.GetOk("auto_scaling.0"); ok {
		autoscalingSetGroupAutoscaling := &clouddatabasesv5.AutoscalingSetGroupAutoscaling{}

		if diskRecord, ok := d.GetOk("auto_scaling.0.disk"); ok {
			diskGroup, err := expandAutoscalingDiskGroup(d, diskRecord)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error in getting diskGroup from expandAutoscalingDiskGroup %s", err))
			}
			autoscalingSetGroupAutoscaling.Disk = diskGroup
		}

		if memoryRecord, ok := d.GetOk("auto_scaling.0.memory"); ok {
			memoryGroup, err := expandAutoscalingMemoryGroup(d, memoryRecord)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error in getting memoryBody from expandAutoscalingMemoryGroup %s", err))
			}

			autoscalingSetGroupAutoscaling.Memory = memoryGroup
		}

		if autoscalingSetGroupAutoscaling.Disk != nil || autoscalingSetGroupAutoscaling.Memory != nil {
			setAutoscalingConditionsOptions := &clouddatabasesv5.SetAutoscalingConditionsOptions{
				ID:          &instanceID,
				GroupID:     core.StringPtr("member"),
				Autoscaling: autoscalingSetGroupAutoscaling,
			}

			setAutoscalingConditionsResponse, _, err := cloudDatabasesClient.SetAutoscalingConditions(setAutoscalingConditionsOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error updating database auto_scaling: %s", err))
			}

			taskId := *setAutoscalingConditionsResponse.Task.ID

			_, err = waitForDatabaseTaskComplete(taskId, d, meta, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for database (%s) memory auto_scaling group update task to complete: %s", instanceID, err))
			}
		}
	}

	if userList, ok := d.GetOk("users"); ok {
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
		}

		for _, user := range userList.(*schema.Set).List() {
			userEl := user.(map[string]interface{})
			err := userUpdateCreate(userEl, instanceID, meta, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if config, ok := d.GetOk("configuration"); ok {
		var rawConfig map[string]json.RawMessage
		err = json.Unmarshal([]byte(config.(string)), &rawConfig)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] configuration JSON invalid\n%s", err))
		}

		var configuration clouddatabasesv5.ConfigurationIntf = new(clouddatabasesv5.Configuration)
		err = core.UnmarshalModel(rawConfig, "", &configuration, clouddatabasesv5.UnmarshalConfiguration)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] database configuration is invalid"))
		}

		updateDatabaseConfigurationOptions := &clouddatabasesv5.UpdateDatabaseConfigurationOptions{
			ID:            &instanceID,
			Configuration: configuration,
		}

		updateDatabaseConfigurationResponse, response, err := cloudDatabasesClient.UpdateDatabaseConfiguration(updateDatabaseConfigurationOptions)

		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error updating database configuration failed %s\n%s", err, response))
		}

		taskID := *updateDatabaseConfigurationResponse.Task.ID

		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for database (%s) configuration update task to complete: %s", icdId, err))
		}
	}

	if _, ok := d.GetOk("logical_replication_slot"); ok {
		service := d.Get("service").(string)
		if service != "databases-for-postgresql" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Logical Replication can only be set for databases-for-postgresql instances"))
		}

		cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
		}

		_, logicalReplicationList := d.GetChange("logical_replication_slot")

		add := logicalReplicationList.(*schema.Set).List()

		for _, entry := range add {
			newEntry := entry.(map[string]interface{})
			logicalReplicationSlot := &clouddatabasesv5.LogicalReplicationSlot{
				Name:         core.StringPtr(newEntry["name"].(string)),
				DatabaseName: core.StringPtr(newEntry["database_name"].(string)),
				PluginType:   core.StringPtr(newEntry["plugin_type"].(string)),
			}

			createLogicalReplicationOptions := &clouddatabasesv5.CreateLogicalReplicationSlotOptions{
				ID:                     &instanceID,
				LogicalReplicationSlot: logicalReplicationSlot,
			}

			createLogicalRepSlotResponse, response, err := cloudDatabasesClient.CreateLogicalReplicationSlot(createLogicalReplicationOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] CreateLogicalReplicationSlot (%s) failed %s\n%s", *createLogicalReplicationOptions.LogicalReplicationSlot.Name, err, response))
			}

			taskID := *createLogicalRepSlotResponse.Task.ID
			_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(fmt.Errorf(
					"[ERROR] Error waiting for database (%s) logical replication slot (%s) create task to complete: %s", instanceID, *createLogicalReplicationOptions.LogicalReplicationSlot.Name, err))
			}
		}
	}

	return resourceIBMDatabaseInstanceRead(context, d, meta)
}

func resourceIBMDatabaseInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	connectionEndpoint := "public"
	rsInst := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := rsConClient.GetResourceInstance(&rsInst)
	if err != nil {
		if strings.Contains(err.Error(), "Object not found") ||
			strings.Contains(err.Error(), "status code: 404") {
			log.Printf("[WARN] Removing record from state because it's not found via the API")
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving resource instance: %s %s", err, response))
	}
	if strings.Contains(*instance.State, "removed") {
		log.Printf("[WARN] Removing instance from TF state because it's now in removed state")
		d.SetId("")
		return nil
	}

	tags, err := flex.GetTagsUsingCRN(meta, *instance.CRN)
	if err != nil {
		log.Printf(
			"Error on get of ibm Database tags (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)
	d.Set("name", *instance.Name)
	d.Set("status", *instance.State)
	d.Set("resource_group_id", *instance.ResourceGroupID)
	if instance.CRN != nil {
		location := strings.Split(*instance.CRN, ":")
		if len(location) > 5 {
			d.Set("location", location[5])
		}
	}
	d.Set("guid", *instance.GUID)

	if instance.Parameters != nil {
		if endpoint, ok := instance.Parameters["service-endpoints"]; ok {
			if endpoint == "private" {
				connectionEndpoint = "private"
			}
			d.Set("service_endpoints", endpoint)
		}

	}

	d.Set(flex.ResourceName, *instance.Name)
	d.Set(flex.ResourceCRN, *instance.CRN)
	d.Set(flex.ResourceStatus, *instance.State)
	d.Set(flex.ResourceGroupName, *instance.ResourceGroupCRN)

	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(flex.ResourceControllerURL, rcontroller+"/services/"+url.QueryEscape(*instance.CRN))

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.GetServiceName(*instance.ResourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving service offering: %s", err))
	}

	d.Set("service", serviceOff)

	servicePlan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving plan: %s", err))
	}
	d.Set("plan", servicePlan)

	icdClient, err := meta.(conns.ClientSession).ICDAPI()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
	}

	icdId := flex.EscapeUrlParm(instanceID)

	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
	}

	getDeploymentInfoOptions := &clouddatabasesv5.GetDeploymentInfoOptions{
		ID: core.StringPtr(instanceID),
	}
	getDeploymentInfoResponse, response, err := cloudDatabasesClient.GetDeploymentInfo(getDeploymentInfoOptions)

	if err != nil {
		if response.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("[ERROR] The database instance was not found in the region set for the Provider, or the default of us-south. Specify the correct region in the provider definition, or create a provider alias for the correct region. %v", err))
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database config while updating adminpassword for: %s with error %s", instanceID, err))
	}

	deployment := getDeploymentInfoResponse.Deployment

	d.Set("adminuser", deployment.AdminUsernames["database"])
	d.Set("version", deployment.Version)

	groupList, err := icdClient.Groups().GetGroups(icdId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database groups: %s", err))
	}
	if groupList.Groups[0].Members.AllocationCount == 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] This database appears to have have 0 members. Unable to proceed"))
	}

	d.Set("groups", flex.FlattenIcdGroups(groupList))
	d.Set("node_count", groupList.Groups[0].Members.AllocationCount)

	d.Set("members_memory_allocation_mb", groupList.Groups[0].Memory.AllocationMb)
	d.Set("node_memory_allocation_mb", groupList.Groups[0].Memory.AllocationMb/groupList.Groups[0].Members.AllocationCount)

	d.Set("members_disk_allocation_mb", groupList.Groups[0].Disk.AllocationMb)
	d.Set("node_disk_allocation_mb", groupList.Groups[0].Disk.AllocationMb/groupList.Groups[0].Members.AllocationCount)

	d.Set("members_cpu_allocation_count", groupList.Groups[0].Cpu.AllocationCount)
	d.Set("node_cpu_allocation_count", groupList.Groups[0].Cpu.AllocationCount/groupList.Groups[0].Members.AllocationCount)

	getAutoscalingConditionsOptions := &clouddatabasesv5.GetAutoscalingConditionsOptions{
		ID:      instance.ID,
		GroupID: core.StringPtr("member"),
	}

	autoscalingGroup, _, err := cloudDatabasesClient.GetAutoscalingConditions(getAutoscalingConditionsOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database autoscaling groups: %s", err))
	}
	d.Set("auto_scaling", flattenAutoScalingGroup(*autoscalingGroup))

	_, hasWhitelist := d.GetOk("whitelist")

	alEntry := &clouddatabasesv5.GetAllowlistOptions{
		ID: &instanceID,
	}

	allowlist, _, err := cloudDatabasesClient.GetAllowlist(alEntry)

	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database allowlist: %s", err))
	}

	if hasWhitelist {
		d.Set("whitelist", flex.FlattenAllowlist(allowlist.IPAddresses))
	} else {
		d.Set("allowlist", flex.FlattenAllowlist(allowlist.IPAddresses))
	}

	var connectionStrings []flex.CsEntry
	//ICD does not implement a GetUsers API. Users populated from tf configuration.
	tfusers := d.Get("users").(*schema.Set)
	users := flex.ExpandUsers(tfusers)
	user := icdv4.User{
		UserName: deployment.AdminUsernames["database"],
	}
	users = append(users, user)
	for _, user := range users {
		userName := user.UserName
		csEntry, err := getConnectionString(d, userName, connectionEndpoint, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting user connection string for user (%s): %s", userName, err))
		}
		connectionStrings = append(connectionStrings, csEntry)
	}
	d.Set("connectionstrings", flex.FlattenConnectionStrings(connectionStrings))

	if serviceOff == "databases-for-postgresql" || serviceOff == "databases-for-redis" || serviceOff == "databases-for-enterprisedb" {
		configSchema, err := icdClient.Configurations().GetConfiguration(icdId)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting database (%s) configuration schema : %s", icdId, err))
		}
		s, err := json.Marshal(configSchema)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error marshalling the database configuration schema: %s", err))
		}

		if err = d.Set("configuration_schema", string(s)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting the database configuration schema: %s", err))
		}
	}
	return nil
}

func resourceIBMDatabaseInstanceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	updateReq := rc.UpdateResourceInstanceOptions{
		ID: &instanceID,
	}
	update := false
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateReq.Name = &name
		update = true
	}
	if d.HasChange("service_endpoints") {
		params := Params{}
		params.ServiceEndpoints = d.Get("service_endpoints").(string)
		parameters, _ := json.Marshal(params)
		var raw map[string]interface{}
		json.Unmarshal(parameters, &raw)
		updateReq.Parameters = raw
		update = true
	}

	if update {
		_, response, err := rsConClient.UpdateResourceInstance(&updateReq)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating resource instance: %s %s", err, response))
		}

		_, err = waitForDatabaseInstanceUpdate(d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for update of resource instance (%s) to complete: %s", d.Id(), err))
		}
	}

	if d.HasChange("tags") {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, instanceID)
		if err != nil {
			log.Printf(
				"[ERROR] Error on update of Database (%s) tags: %s", d.Id(), err)
		}
	}

	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
	}

	icdClient, err := meta.(conns.ClientSession).ICDAPI()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
	}
	icdId := flex.EscapeUrlParm(instanceID)

	if d.HasChange("node_count") {
		err = horizontalScale(d, meta, icdClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("configuration") {
		if config, ok := d.GetOk("configuration"); ok {
			var rawConfig map[string]json.RawMessage
			err = json.Unmarshal([]byte(config.(string)), &rawConfig)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] configuration JSON invalid\n%s", err))
			}

			var configuration clouddatabasesv5.ConfigurationIntf = new(clouddatabasesv5.Configuration)
			err = core.UnmarshalModel(rawConfig, "", &configuration, clouddatabasesv5.UnmarshalConfiguration)
			if err != nil {
				return diag.FromErr(err)
			}

			updateDatabaseConfigurationOptions := &clouddatabasesv5.UpdateDatabaseConfigurationOptions{
				ID:            &instanceID,
				Configuration: configuration,
			}

			updateDatabaseConfigurationResponse, response, err := cloudDatabasesClient.UpdateDatabaseConfiguration(updateDatabaseConfigurationOptions)

			if err != nil {
				return diag.FromErr(fmt.Errorf(
					"[ERROR] Error updating database configuration failed %s\n%s", err, response))
			}

			taskID := *updateDatabaseConfigurationResponse.Task.ID

			_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(fmt.Errorf(
					"[ERROR] Error waiting for database (%s) configuration update task to complete: %s", icdId, err))
			}
		}
	}

	if d.HasChange("members_memory_allocation_mb") || d.HasChange("members_disk_allocation_mb") || d.HasChange("members_cpu_allocation_count") || d.HasChange("node_memory_allocation_mb") || d.HasChange("node_disk_allocation_mb") || d.HasChange("node_cpu_allocation_count") {
		params := icdv4.GroupReq{}
		if d.HasChange("members_memory_allocation_mb") {
			memory := d.Get("members_memory_allocation_mb").(int)
			memoryReq := icdv4.MemoryReq{AllocationMb: memory}
			params.GroupBdy.Memory = &memoryReq
		}
		if d.HasChange("node_memory_allocation_mb") || d.HasChange("node_count") {
			memory := d.Get("node_memory_allocation_mb").(int)
			count := d.Get("node_count").(int)
			memoryReq := icdv4.MemoryReq{AllocationMb: memory * count}
			params.GroupBdy.Memory = &memoryReq
		}
		if d.HasChange("members_disk_allocation_mb") {
			disk := d.Get("members_disk_allocation_mb").(int)
			diskReq := icdv4.DiskReq{AllocationMb: disk}
			params.GroupBdy.Disk = &diskReq
		}
		if d.HasChange("node_disk_allocation_mb") || d.HasChange("node_count") {
			disk := d.Get("node_disk_allocation_mb").(int)
			count := d.Get("node_count").(int)
			diskReq := icdv4.DiskReq{AllocationMb: disk * count}
			params.GroupBdy.Disk = &diskReq
		}
		if d.HasChange("members_cpu_allocation_count") {
			cpu := d.Get("members_cpu_allocation_count").(int)
			cpuReq := icdv4.CpuReq{AllocationCount: cpu}
			params.GroupBdy.Cpu = &cpuReq
		}
		if d.HasChange("node_cpu_allocation_mb") || d.HasChange("node_count") {
			cpu := d.Get("node_cpu_allocation_count").(int)
			count := d.Get("node_count").(int)
			CpuReq := icdv4.CpuReq{AllocationCount: cpu * count}
			params.GroupBdy.Cpu = &CpuReq
		}

		task, err := icdClient.Groups().UpdateGroup(icdId, "member", params)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database scaling group: %s", err))
		}

		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for database (%s) scaling group update task to complete: %s", icdId, err))
		}
	}

	if d.HasChange("group") {
		oldGroup, newGroup := d.GetChange("group")
		if oldGroup == nil {
			oldGroup = new(schema.Set)
		}
		if newGroup == nil {
			newGroup = new(schema.Set)
		}

		os := oldGroup.(*schema.Set)
		ns := newGroup.(*schema.Set)

		groupChanges := expandGroups(ns.Difference(os).List())

		groupsResponse, err := getGroups(instanceID, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error geting group (%s) scaling group update task to complete: %s", icdId, err))
		}

		currentGroups := normalizeGroups(groupsResponse)

		for _, group := range groupChanges {
			groupScaling := &clouddatabasesv5.GroupScaling{}
			var currentGroup *Group
			for _, g := range currentGroups {
				if g.ID == group.ID {
					currentGroup = &g
					break
				}
			}

			if currentGroup == nil {
				return diag.FromErr(fmt.Errorf(
					"[ERROR]  (%s) group does not exist: %s", icdId, err))
			}
			nodeCount := currentGroup.Members.Allocation

			if group.Members != nil && group.Members.Allocation != currentGroup.Members.Allocation {
				groupScaling.Members = &clouddatabasesv5.GroupScalingMembers{AllocationCount: core.Int64Ptr(int64(group.Members.Allocation))}
				nodeCount = group.Members.Allocation
			}
			if group.Memory != nil && group.Memory.Allocation*nodeCount != currentGroup.Memory.Allocation {
				groupScaling.Memory = &clouddatabasesv5.GroupScalingMemory{AllocationMb: core.Int64Ptr(int64(group.Memory.Allocation * nodeCount))}
			}
			if group.Disk != nil && group.Disk.Allocation*nodeCount != currentGroup.Disk.Allocation {
				groupScaling.Disk = &clouddatabasesv5.GroupScalingDisk{AllocationMb: core.Int64Ptr(int64(group.Disk.Allocation * nodeCount))}
			}
			if group.CPU != nil && group.CPU.Allocation*nodeCount != currentGroup.CPU.Allocation {
				groupScaling.CPU = &clouddatabasesv5.GroupScalingCPU{AllocationCount: core.Int64Ptr(int64(group.CPU.Allocation * nodeCount))}
			}

			if groupScaling.Members != nil || groupScaling.Memory != nil || groupScaling.Disk != nil || groupScaling.CPU != nil {
				setDeploymentScalingGroupOptions := &clouddatabasesv5.SetDeploymentScalingGroupOptions{
					ID:      &instanceID,
					GroupID: &group.ID,
					Group:   groupScaling,
				}

				setDeploymentScalingGroupResponse, response, err := cloudDatabasesClient.SetDeploymentScalingGroup(setDeploymentScalingGroupOptions)

				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] SetDeploymentScalingGroup (%s) failed %s\n%s", group.ID, err, response))
				}

				// API may return HTTP 204 No Content if no change made
				if response.StatusCode == 202 {
					taskIDLink := *setDeploymentScalingGroupResponse.Task.ID

					_, err = waitForDatabaseTaskComplete(taskIDLink, d, meta, d.Timeout(schema.TimeoutCreate))

					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	if d.HasChange("auto_scaling.0") {
		autoscalingSetGroupAutoscaling := &clouddatabasesv5.AutoscalingSetGroupAutoscaling{}

		if d.HasChange("auto_scaling.0.disk") {
			diskRecord := d.Get("auto_scaling.0.disk")

			diskBody, err := expandAutoscalingDiskGroup(d, diskRecord)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error in getting diskBody from expandAutoscalingDiskGroup %s", err))
			}

			autoscalingSetGroupAutoscaling.Disk = diskBody
		}

		if d.HasChange("auto_scaling.0.memory") {
			memoryRecord := d.Get("auto_scaling.0.memory")
			memoryBody, err := expandAutoscalingMemoryGroup(d, memoryRecord)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error in getting memoryBody from expandAutoscalingMemoryGroup %s", err))
			}

			autoscalingSetGroupAutoscaling.Memory = memoryBody
		}

		setAutoscalingConditionsOptions := &clouddatabasesv5.SetAutoscalingConditionsOptions{
			ID:          &instanceID,
			GroupID:     core.StringPtr("member"),
			Autoscaling: autoscalingSetGroupAutoscaling,
		}

		setAutoscalingConditionsResponse, _, err := cloudDatabasesClient.SetAutoscalingConditions(setAutoscalingConditionsOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database memory auto_scaling group: %s", err))
		}

		taskId := *setAutoscalingConditionsResponse.Task.ID

		_, err = waitForDatabaseTaskComplete(taskId, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for database (%s) auto scaling group update task to complete: %s", instanceID, err))
		}
	}

	if d.HasChange("adminpassword") {
		adminUser := d.Get("adminuser").(string)
		password := d.Get("adminpassword").(string)
		user := &clouddatabasesv5.APasswordSettingUser{
			Password: &password,
		}

		changeUserPasswordOptions := &clouddatabasesv5.ChangeUserPasswordOptions{
			ID:       core.StringPtr(instanceID),
			UserType: core.StringPtr("database"),
			Username: core.StringPtr(adminUser),
			User:     user,
		}

		changeUserPasswordResponse, response, err := cloudDatabasesClient.ChangeUserPassword(changeUserPasswordOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] ChangeUserPassword (%s) failed %s\n%s", *changeUserPasswordOptions.Username, err, response))
		}

		taskID := *changeUserPasswordResponse.Task.ID
		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database admin password: %s", err))
		}
	}

	if d.HasChange("whitelist") || d.HasChange("allowlist") {
		_, hasAllowlist := d.GetOk("allowlist")
		_, hasWhitelist := d.GetOk("whitelist")

		var entries interface{}

		if hasWhitelist {
			_, entries = d.GetChange("whitelist")
		} else if hasAllowlist {
			_, entries = d.GetChange("allowlist")
		}

		if entries == nil {
			entries = new(schema.Set)
		}

		allowlistEntries := flex.ExpandAllowlist(entries.(*schema.Set))

		setAllowlistOptions := &clouddatabasesv5.SetAllowlistOptions{
			ID:          &instanceID,
			IPAddresses: allowlistEntries,
		}

		setAllowlistResponse, _, err := cloudDatabasesClient.SetAllowlist(setAllowlistOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database allowlist entry: %s", err))
		}

		taskId := *setAllowlistResponse.Task.ID

		_, err = waitForDatabaseTaskComplete(taskId, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(fmt.Errorf(
				"[ERROR] Error waiting for update of database (%s) whitelist task to complete: %s", instanceID, err))
		}
	}

	if d.HasChange("users") {
		oldUsers, newUsers := d.GetChange("users")
		userChanges := make(map[string]*userChange)
		userKey := func(raw map[string]interface{}) string {
			if raw["role"].(string) != "" {
				return fmt.Sprintf("%s-%s-%s", raw["type"].(string), raw["role"].(string), raw["name"].(string))
			} else {
				return fmt.Sprintf("%s-%s", raw["type"].(string), raw["name"].(string))
			}
		}

		for _, raw := range oldUsers.(*schema.Set).List() {
			user := raw.(map[string]interface{})
			k := userKey(user)
			userChanges[k] = &userChange{Old: user}
		}

		for _, raw := range newUsers.(*schema.Set).List() {
			user := raw.(map[string]interface{})
			k := userKey(user)
			if _, ok := userChanges[k]; !ok {
				userChanges[k] = &userChange{}
			}
			userChanges[k].New = user
		}

		for _, change := range userChanges {
			// Delete Old User
			if change.Old != nil && change.New == nil {
				deleteDatabaseUserOptions := &clouddatabasesv5.DeleteDatabaseUserOptions{
					ID:       &instanceID,
					UserType: core.StringPtr(change.Old["type"].(string)),
					Username: core.StringPtr(change.Old["name"].(string)),
				}

				deleteDatabaseUserResponse, response, err := cloudDatabasesClient.DeleteDatabaseUser(deleteDatabaseUserOptions)

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] DeleteDatabaseUser (%s) failed %s\n%s", *deleteDatabaseUserOptions.Username, err, response))

				}

				taskID := *deleteDatabaseUserResponse.Task.ID
				_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] Error waiting for database (%s) user (%s) delete task to complete: %s", icdId, *deleteDatabaseUserOptions.Username, err))
				}

				continue
			}

			if change.New != nil {
				// No change
				if change.Old != nil && change.Old["password"].(string) == change.New["password"].(string) && change.Old["name"].(string) == change.New["name"].(string) {
					continue
				}

				err := userUpdateCreate(change.New, instanceID, meta, d)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("logical_replication_slot") {
		service := d.Get("service").(string)
		if service != "databases-for-postgresql" {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Logical Replication can only be set for databases-for-postgresql instances"))
		}

		cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
		}

		oldList, newList := d.GetChange("logical_replication_slot")

		if oldList == nil {
			oldList = new(schema.Set)
		}
		if newList == nil {
			newList = new(schema.Set)
		}
		os := oldList.(*schema.Set)
		ns := newList.(*schema.Set)

		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		// Create New Logical Rep Slot
		if len(add) > 0 {
			for _, entry := range add {
				newEntry := entry.(map[string]interface{})
				logicalReplicationSlot := &clouddatabasesv5.LogicalReplicationSlot{
					Name:         core.StringPtr(newEntry["name"].(string)),
					DatabaseName: core.StringPtr(newEntry["database_name"].(string)),
					PluginType:   core.StringPtr(newEntry["plugin_type"].(string)),
				}

				createLogicalReplicationOptions := &clouddatabasesv5.CreateLogicalReplicationSlotOptions{
					ID:                     &instanceID,
					LogicalReplicationSlot: logicalReplicationSlot,
				}

				createLogicalRepSlotResponse, response, err := cloudDatabasesClient.CreateLogicalReplicationSlot(createLogicalReplicationOptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] CreateLogicalReplicationSlot (%s) failed %s\n%s", *createLogicalReplicationOptions.LogicalReplicationSlot.Name, err, response))
				}

				taskID := *createLogicalRepSlotResponse.Task.ID
				_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] Error waiting for database (%s) logical replication slot (%s) create task to complete: %s", instanceID, *createLogicalReplicationOptions.LogicalReplicationSlot.Name, err))
				}
			}
		}

		// Delete Old Logical Rep Slot
		if len(remove) > 0 {
			for _, entry := range remove {
				newEntry := entry.(map[string]interface{})
				deleteDatabaseUserOptions := &clouddatabasesv5.DeleteLogicalReplicationSlotOptions{
					ID:   &instanceID,
					Name: core.StringPtr(newEntry["name"].(string)),
				}

				deleteDatabaseUserResponse, response, err := cloudDatabasesClient.DeleteLogicalReplicationSlot(deleteDatabaseUserOptions)

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] DeleteDatabaseUser (%s) failed %s\n%s", *deleteDatabaseUserOptions.Name, err, response))
				}

				taskID := *deleteDatabaseUserResponse.Task.ID
				_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] Error waiting for database (%s) logical replication slot (%s) delete task to complete: %s", icdId, *deleteDatabaseUserOptions.Name, err))
				}
			}
		}
	}

	return resourceIBMDatabaseInstanceRead(context, d, meta)
}

func horizontalScale(d *schema.ResourceData, meta interface{}, icdClient icdv4.ICDServiceAPI) error {
	params := icdv4.GroupReq{}

	icdId := flex.EscapeUrlParm(d.Id())

	members := d.Get("node_count").(int)
	membersReq := icdv4.MembersReq{AllocationCount: members}
	params.GroupBdy.Members = &membersReq

	_, err := icdClient.Groups().UpdateGroup(icdId, "member", params)

	if err != nil {
		return fmt.Errorf("[ERROR] Error updating database scaling group: %s", err)
	}

	// ScaleOut is handled with an ICD API call, however, the check is is on the instance status
	_, err = waitForDatabaseInstanceUpdate(d, meta)
	if err != nil {
		return fmt.Errorf(
			"[ERROR] Error waiting for database (%s) horizontal scale to complete: %s", d.Id(), err)
	}

	return nil
}

func getConnectionString(d *schema.ResourceData, userName, connectionEndpoint string, meta interface{}) (flex.CsEntry, error) {
	csEntry := flex.CsEntry{}
	icdClient, err := meta.(conns.ClientSession).ICDAPI()
	if err != nil {
		return csEntry, fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	icdId := d.Id()
	connection, err := icdClient.Connections().GetConnection(icdId, userName, connectionEndpoint)
	if err != nil {
		return csEntry, fmt.Errorf("[ERROR] Error getting database user connection string via ICD API: %s", err)
	}

	service := d.Get("service")
	dbConnection := icdv4.Uri{}
	var cassandraConnection icdv4.CassandraUri

	switch service {
	case "databases-for-postgresql":
		dbConnection = connection.Postgres
	case "databases-for-redis":
		dbConnection = connection.Rediss
	case "databases-for-mongodb":
		dbConnection = connection.Mongo
	case "databases-for-mysql":
		dbConnection = connection.Mysql
	case "databases-for-elasticsearch":
		dbConnection = connection.Https
	case "databases-for-cassandra":
		cassandraConnection = connection.Secure
	case "databases-for-etcd":
		dbConnection = connection.Grpc
	case "messages-for-rabbitmq":
		dbConnection = connection.Amqps
	case "databases-for-enterprisedb":
		dbConnection = connection.Postgres
	default:
		return csEntry, fmt.Errorf("[ERROR] Unrecognised database type during connection string lookup: %s", service)
	}

	if !reflect.DeepEqual(cassandraConnection, icdv4.CassandraUri{}) {
		csEntry = flex.CsEntry{
			Name:         userName,
			Hosts:        cassandraConnection.Hosts,
			BundleName:   cassandraConnection.Bundle.Name,
			BundleBase64: cassandraConnection.Bundle.BundleBase64,
		}
	} else {
		csEntry = flex.CsEntry{
			Name:     userName,
			Password: "",
			// Populate only first 'composed' connection string as an example
			Composed:     dbConnection.Composed[0],
			CertName:     dbConnection.Certificate.Name,
			CertBase64:   dbConnection.Certificate.CertificateBase64,
			Hosts:        dbConnection.Hosts,
			Scheme:       dbConnection.Scheme,
			Path:         dbConnection.Path,
			QueryOptions: dbConnection.QueryOptions.(map[string]interface{}),
		}

		// Postgres DB name is of type string, Redis is json.Number, others are nil
		if dbConnection.Database != nil {
			switch v := dbConnection.Database.(type) {
			default:
				return csEntry, fmt.Errorf("Unexpected data type: %T", v)
			case json.Number:
				csEntry.Database = dbConnection.Database.(json.Number).String()
			case string:
				csEntry.Database = dbConnection.Database.(string)
			}
		} else {
			csEntry.Database = ""
		}
	}

	return csEntry, nil
}

func resourceIBMDatabaseInstanceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	recursive := true
	deleteReq := rc.DeleteResourceInstanceOptions{
		Recursive: &recursive,
		ID:        &id,
	}
	response, err := rsConClient.DeleteResourceInstance(&deleteReq)
	if err != nil {
		// If prior delete occurs, instance is not immediately deleted, but remains in "removed" state"
		// RC 410 with "Gone" returned as error
		if strings.Contains(err.Error(), "Gone") ||
			strings.Contains(err.Error(), "status code: 410") {
			log.Printf("[WARN] Resource instance already deleted %s\n ", err)
			err = nil
		} else {
			return diag.FromErr(fmt.Errorf("[ERROR] Error deleting resource instance: %s %s ", err, response))
		}
	}

	_, err = waitForDatabaseInstanceDelete(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"[ERROR] Error waiting for resource instance (%s) to be deleted: %s", d.Id(), err))
	}

	d.SetId("")

	return nil
}
func resourceIBMDatabaseInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	rsInst := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := rsConClient.GetResourceInstance(&rsInst)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting database: %s %s", err, response)
	}
	if instance != nil && (strings.Contains(*instance.State, "removed") || strings.Contains(*instance.State, databaseInstanceReclamation)) {
		log.Printf("[WARN] Removing instance from state because it's in removed or pending_reclamation state")
		d.SetId("")
		return false, nil
	}

	return *instance.ID == instanceID, nil
}

func waitForICDReady(meta interface{}, instanceID string) error {
	icdId := flex.EscapeUrlParm(instanceID)
	icdClient, clientErr := meta.(conns.ClientSession).ICDAPI()
	if clientErr != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %s", clientErr)
	}

	// Wait for ICD Interface
	err := retry(func() (err error) {
		_, cdbErr := icdClient.Cdbs().GetCdb(icdId)
		if cdbErr != nil {
			if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
				return fmt.Errorf("[ERROR] The database instance was not found in the region set for the Provider, or the default of us-south. Specify the correct region in the provider definition, or create a provider alias for the correct region. %v", err)
			}
			return fmt.Errorf("[ERROR] Error getting database config for: %s with error %s\n", icdId, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func waitForDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}, instanceID string) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{databaseInstanceProgressStatus, databaseInstanceInactiveStatus, databaseInstanceProvisioningStatus},
		Target:  []string{databaseInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil || instance == nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %s %s", d.Id(), err, response)
				}
				return nil, "", fmt.Errorf("[ERROR] GetResourceInstance on %s failed with error %s %s", d.Id(), err, response)
			}
			if *instance.State == databaseInstanceFailStatus {
				return *instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed: %s %s", d.Id(), err, response)
			}
			return *instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	waitErr := waitForICDReady(meta, instanceID)
	if waitErr != nil {
		return false, fmt.Errorf("[ERROR] Error ICD interface not ready after create: %s with error %s\n", instanceID, waitErr)

	}

	return stateConf.WaitForState()
}

func waitForDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending: []string{databaseInstanceProgressStatus, databaseInstanceInactiveStatus},
		Target:  []string{databaseInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %s %s", d.Id(), err, response)
				}
				return nil, "", fmt.Errorf("[ERROR] GetResourceInstance on %s failed with error %s %s", d.Id(), err, response)
			}
			if *instance.State == databaseInstanceFailStatus {
				return *instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed: %s %s", d.Id(), err, response)
			}
			return *instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	waitErr := waitForICDReady(meta, instanceID)
	if waitErr != nil {
		return false, fmt.Errorf("[ERROR] Error ICD interface not ready after update: %s with error %s\n", instanceID, waitErr)

	}

	return stateConf.WaitForState()
}

func waitForDatabaseTaskComplete(taskId string, d *schema.ResourceData, meta interface{}, t time.Duration) (bool, error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	delayDuration := 5 * time.Second

	timeout := time.After(t)
	delay := time.Tick(delayDuration)
	getTaskOptions := &clouddatabasesv5.GetTaskOptions{
		ID: &taskId,
	}

	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("[Error] Time out waiting for database task to complete")
		case <-delay:
			getTaskResponse, _, err := cloudDatabasesClient.GetTask(getTaskOptions)

			if err != nil {
				return false, fmt.Errorf("[ERROR] Database Task errored: %v", err)
			}

			if getTaskResponse.Task == nil {
				return true, nil
			}

			switch *getTaskResponse.Task.Status {
			case "failed":
				return false, fmt.Errorf("[Error] Database Task failed")
			case "complete", "":
				return true, nil
			case "queued", "running":
				break
			}
		}
	}
}

func waitForDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending: []string{databaseInstanceProgressStatus, databaseInstanceInactiveStatus, databaseInstanceSuccessStatus},
		Target:  []string{databaseInstanceRemovedStatus, databaseInstanceReclamation},
		Refresh: func() (interface{}, string, error) {
			rsInst := rc.GetResourceInstanceOptions{
				ID: &instanceID,
			}
			instance, response, err := rsConClient.GetResourceInstance(&rsInst)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return instance, databaseInstanceSuccessStatus, nil
				}
				return nil, "", fmt.Errorf("[ERROR] GetResourceInstance on %s failed with error %s %s", d.Id(), err, response)
			}
			if *instance.State == databaseInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance %s failed to delete: %s %s", d.Id(), err, response)
			}
			return *instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func filterDatabaseDeployments(deployments []models.ServiceDeployment, location string) ([]models.ServiceDeployment, map[string]bool) {
	supportedDeployments := []models.ServiceDeployment{}
	supportedLocations := make(map[string]bool)
	for _, d := range deployments {
		if d.Metadata.RCCompatible {
			deploymentLocation := d.Metadata.Deployment.Location
			supportedLocations[deploymentLocation] = true
			if deploymentLocation == location {
				supportedDeployments = append(supportedDeployments, d)
			}
		}
	}
	return supportedDeployments, supportedLocations
}

func expandAutoscalingDiskGroup(d *schema.ResourceData, asRecord interface{}) (autoscalingDiskGroup *clouddatabasesv5.AutoscalingDiskGroupDisk, err error) {
	autoscalingRecord := asRecord.([]interface{})[0].(map[string]interface{})

	autoscalingDiskGroup = &clouddatabasesv5.AutoscalingDiskGroupDisk{
		Scalers: &clouddatabasesv5.AutoscalingDiskGroupDiskScalers{
			Capacity:      &clouddatabasesv5.AutoscalingDiskGroupDiskScalersCapacity{},
			IoUtilization: &clouddatabasesv5.AutoscalingDiskGroupDiskScalersIoUtilization{},
		},
		Rate: &clouddatabasesv5.AutoscalingDiskGroupDiskRate{},
	}

	if _, ok := autoscalingRecord["capacity_enabled"]; ok {
		autoscalingDiskGroup.Scalers.Capacity.Enabled = core.BoolPtr(autoscalingRecord["capacity_enabled"].(bool))
	}
	if _, ok := autoscalingRecord["free_space_less_than_percent"]; ok {
		autoscalingDiskGroup.Scalers.Capacity.FreeSpaceLessThanPercent = core.Int64Ptr(int64(autoscalingRecord["free_space_less_than_percent"].(int)))
	}

	// IO Payload
	if _, ok := autoscalingRecord["io_enabled"]; ok {
		autoscalingDiskGroup.Scalers.IoUtilization.Enabled = core.BoolPtr(autoscalingRecord["io_enabled"].(bool))
	}
	if _, ok := autoscalingRecord["io_over_period"]; ok {
		autoscalingDiskGroup.Scalers.IoUtilization.OverPeriod = core.StringPtr(autoscalingRecord["io_over_period"].(string))
	}
	if _, ok := autoscalingRecord["io_above_percent"]; ok {
		autoscalingDiskGroup.Scalers.IoUtilization.AbovePercent = core.Int64Ptr(int64(autoscalingRecord["io_above_percent"].(int)))
	}

	// Rate Payload
	if _, ok := autoscalingRecord["rate_increase_percent"]; ok {
		autoscalingDiskGroup.Rate.IncreasePercent = core.Float64Ptr(float64(autoscalingRecord["rate_increase_percent"].(int)))
	}
	if _, ok := autoscalingRecord["rate_period_seconds"]; ok {
		autoscalingDiskGroup.Rate.PeriodSeconds = core.Int64Ptr(int64(autoscalingRecord["rate_period_seconds"].(int)))
	}
	if _, ok := autoscalingRecord["rate_limit_mb_per_member"]; ok {
		autoscalingDiskGroup.Rate.LimitMbPerMember = core.Float64Ptr(float64(autoscalingRecord["rate_limit_mb_per_member"].(int)))
	}
	if _, ok := autoscalingRecord["rate_units"]; ok {
		autoscalingDiskGroup.Rate.Units = core.StringPtr(autoscalingRecord["rate_units"].(string))
	}

	return
}

func expandAutoscalingMemoryGroup(d *schema.ResourceData, asRecord interface{}) (autoscalingMemoryGroup *clouddatabasesv5.AutoscalingMemoryGroupMemory, err error) {
	autoscalingRecord := asRecord.([]interface{})[0].(map[string]interface{})
	autoscalingMemoryGroup = &clouddatabasesv5.AutoscalingMemoryGroupMemory{
		Scalers: &clouddatabasesv5.AutoscalingMemoryGroupMemoryScalers{
			IoUtilization: &clouddatabasesv5.AutoscalingMemoryGroupMemoryScalersIoUtilization{},
		},
		Rate: &clouddatabasesv5.AutoscalingMemoryGroupMemoryRate{},
	}

	// IO Payload
	if _, ok := autoscalingRecord["io_enabled"]; ok {
		autoscalingMemoryGroup.Scalers.IoUtilization.Enabled = core.BoolPtr(autoscalingRecord["io_enabled"].(bool))
	}
	if _, ok := autoscalingRecord["io_over_period"]; ok {
		autoscalingMemoryGroup.Scalers.IoUtilization.OverPeriod = core.StringPtr(autoscalingRecord["io_over_period"].(string))
	}
	if _, ok := autoscalingRecord["io_above_percent"]; ok {
		autoscalingMemoryGroup.Scalers.IoUtilization.AbovePercent = core.Int64Ptr(int64(autoscalingRecord["io_above_percent"].(int)))
	}

	// Rate Payload
	if _, ok := autoscalingRecord["rate_increase_percent"]; ok {
		autoscalingMemoryGroup.Rate.IncreasePercent = core.Float64Ptr(float64(autoscalingRecord["rate_increase_percent"].(int)))
	}
	if _, ok := autoscalingRecord["rate_period_seconds"]; ok {
		autoscalingMemoryGroup.Rate.PeriodSeconds = core.Int64Ptr(int64(autoscalingRecord["rate_period_seconds"].(int)))
	}
	if _, ok := autoscalingRecord["rate_limit_mb_per_member"]; ok {
		autoscalingMemoryGroup.Rate.LimitMbPerMember = core.Float64Ptr(float64(autoscalingRecord["rate_limit_mb_per_member"].(int)))
	}
	if _, ok := autoscalingRecord["rate_units"]; ok {
		autoscalingMemoryGroup.Rate.Units = core.StringPtr(autoscalingRecord["rate_units"].(string))
	}

	return
}

func flattenAutoScalingGroup(autoScalingGroup clouddatabasesv5.AutoscalingGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	memorys := make([]map[string]interface{}, 0)
	memory := make(map[string]interface{})

	if autoScalingGroup.Autoscaling.Memory.Scalers.IoUtilization != nil {
		memoryIO := *autoScalingGroup.Autoscaling.Memory.Scalers.IoUtilization
		memory["io_enabled"] = memoryIO.Enabled
		memory["io_over_period"] = memoryIO.OverPeriod
		memory["io_above_percent"] = memoryIO.AbovePercent
	}

	if &autoScalingGroup.Autoscaling.Memory.Rate != nil {
		ip := autoScalingGroup.Autoscaling.Memory.Rate.IncreasePercent
		memory["rate_increase_percent"] = *ip
		memory["rate_period_seconds"] = autoScalingGroup.Autoscaling.Memory.Rate.PeriodSeconds

		lmp := autoScalingGroup.Autoscaling.Memory.Rate.LimitMbPerMember
		memory["rate_limit_mb_per_member"] = *lmp
		memory["rate_units"] = autoScalingGroup.Autoscaling.Memory.Rate.Units
	}
	memorys = append(memorys, memory)

	cpus := make([]map[string]interface{}, 0)
	cpu := make(map[string]interface{})

	if &autoScalingGroup.Autoscaling.CPU.Rate != nil {
		ip := autoScalingGroup.Autoscaling.CPU.Rate.IncreasePercent
		cpu["rate_increase_percent"] = *ip
		cpu["rate_period_seconds"] = autoScalingGroup.Autoscaling.CPU.Rate.PeriodSeconds
		cpu["rate_limit_count_per_member"] = autoScalingGroup.Autoscaling.CPU.Rate.LimitCountPerMember
		cpu["rate_units"] = autoScalingGroup.Autoscaling.CPU.Rate.Units
	}
	cpus = append(cpus, cpu)

	disks := make([]map[string]interface{}, 0)
	disk := make(map[string]interface{})

	if autoScalingGroup.Autoscaling.Disk.Scalers.Capacity != nil {
		diskCapacity := *autoScalingGroup.Autoscaling.Disk.Scalers.Capacity
		disk["capacity_enabled"] = diskCapacity.Enabled
		disk["free_space_less_than_percent"] = diskCapacity.FreeSpaceLessThanPercent
	}

	if autoScalingGroup.Autoscaling.Disk.Scalers.IoUtilization != nil {
		diskIO := *autoScalingGroup.Autoscaling.Disk.Scalers.IoUtilization
		disk["io_enabled"] = diskIO.Enabled
		disk["io_over_period"] = diskIO.OverPeriod
		disk["io_above_percent"] = diskIO.AbovePercent
	}

	if &autoScalingGroup.Autoscaling.Disk.Rate != nil {
		ip := autoScalingGroup.Autoscaling.Disk.Rate.IncreasePercent
		disk["rate_increase_percent"] = ip
		disk["rate_period_seconds"] = autoScalingGroup.Autoscaling.Disk.Rate.PeriodSeconds

		lpm := autoScalingGroup.Autoscaling.Disk.Rate.LimitMbPerMember
		disk["rate_limit_mb_per_member"] = lpm
		disk["rate_units"] = autoScalingGroup.Autoscaling.Disk.Rate.Units
	}

	disks = append(disks, disk)
	as := map[string]interface{}{
		"memory": memorys,
		"cpu":    cpus,
		"disk":   disks,
	}

	result = append(result, as)
	return result
}

func normalizeGroups(_groups []clouddatabasesv5.Group) (groups []Group) {
	groups = make([]Group, len(_groups))
	for _, g := range _groups {
		group := Group{ID: *g.ID}

		group.Members = &GroupResource{
			Units:        *g.Members.Units,
			Allocation:   int(*g.Members.AllocationCount),
			Minimum:      int(*g.Members.MinimumCount),
			Maximum:      int(*g.Members.MaximumCount),
			StepSize:     int(*g.Members.StepSizeCount),
			IsAdjustable: *g.Members.IsAdjustable,
			IsOptional:   *g.Members.IsOptional,
			CanScaleDown: *g.Members.CanScaleDown,
		}

		group.Memory = &GroupResource{
			Units:        *g.Memory.Units,
			Allocation:   int(*g.Memory.AllocationMb),
			Minimum:      int(*g.Memory.MinimumMb),
			Maximum:      int(*g.Memory.MaximumMb),
			StepSize:     int(*g.Memory.StepSizeMb),
			IsAdjustable: *g.Memory.IsAdjustable,
			IsOptional:   *g.Memory.IsOptional,
			CanScaleDown: *g.Memory.CanScaleDown,
		}

		group.Disk = &GroupResource{
			Units:        *g.Disk.Units,
			Allocation:   int(*g.Disk.AllocationMb),
			Minimum:      int(*g.Disk.MinimumMb),
			Maximum:      int(*g.Disk.MaximumMb),
			StepSize:     int(*g.Disk.StepSizeMb),
			IsAdjustable: *g.Disk.IsAdjustable,
			IsOptional:   *g.Disk.IsOptional,
			CanScaleDown: *g.Disk.CanScaleDown,
		}

		group.CPU = &GroupResource{
			Units:        *g.CPU.Units,
			Allocation:   int(*g.CPU.AllocationCount),
			Minimum:      int(*g.CPU.MinimumCount),
			Maximum:      int(*g.CPU.MaximumCount),
			StepSize:     int(*g.CPU.StepSizeCount),
			IsAdjustable: *g.CPU.IsAdjustable,
			IsOptional:   *g.CPU.IsOptional,
			CanScaleDown: *g.CPU.CanScaleDown,
		}

		groups = append(groups, group)
	}

	return groups
}

func expandGroups(_groups []interface{}) []*Group {
	if len(_groups) == 0 {
		return nil
	}

	groups := make([]*Group, 0, len(_groups))

	for _, groupRaw := range _groups {
		if tfGroup, ok := groupRaw.(map[string]interface{}); ok {
			group := Group{ID: tfGroup["group_id"].(string)}

			if membersSet, ok := tfGroup["members"].(*schema.Set); ok {
				members := membersSet.List()
				if len(members) != 0 {
					group.Members = &GroupResource{Allocation: members[0].(map[string]interface{})["allocation_count"].(int)}
				}
			}

			if memorySet, ok := tfGroup["memory"].(*schema.Set); ok {
				memory := memorySet.List()
				if len(memory) != 0 {
					group.Memory = &GroupResource{Allocation: memory[0].(map[string]interface{})["allocation_mb"].(int)}
				}
			}

			if diskSet, ok := tfGroup["disk"].(*schema.Set); ok {
				disk := diskSet.List()
				if len(disk) != 0 {
					group.Disk = &GroupResource{Allocation: disk[0].(map[string]interface{})["allocation_mb"].(int)}
				}
			}

			if cpuSet, ok := tfGroup["cpu"].(*schema.Set); ok {
				cpu := cpuSet.List()
				if len(cpu) != 0 {
					group.CPU = &GroupResource{Allocation: cpu[0].(map[string]interface{})["allocation_count"].(int)}
				}
			}

			groups = append(groups, &group)
		}
	}

	// analytics must be created before bi_connector
	sortPriority := map[string]int{
		"members":      10,
		"analytics":    2,
		"bi_connector": 1,
	}

	sort.SliceStable(groups, func(i, j int) bool {
		return sortPriority[groups[i].ID] > sortPriority[groups[j].ID]
	})

	return groups
}

func checkV5Groups(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	instanceID := diff.Id()
	service := diff.Get("service").(string)
	plan := diff.Get("plan").(string)

	if group, ok := diff.GetOk("group"); ok {
		var currentGroups []Group
		var groupList []clouddatabasesv5.Group
		var groupIds []string

		if instanceID != "" {
			groupList, err = getGroups(instanceID, meta)
		} else {
			groupList, err = getDefaultScalingGroups(service, plan, meta)
		}

		if err != nil {
			return err
		}

		currentGroups = normalizeGroups(groupList)

		tfGroups := expandGroups(group.(*schema.Set).List())

		// Check group_ids are unique
		groupIds = make([]string, 0, len(tfGroups))
		for _, g := range tfGroups {
			groupIds = append(groupIds, g.ID)
		}

		for n1, i1 := range groupIds {
			for n2, i2 := range groupIds {
				if i1 == i2 && n1 != n2 {
					return fmt.Errorf("found 2 or more instances of group with group_id %v", i1)
				}
			}
		}

		// Get default or current group scaling values
		for _, group := range tfGroups {
			if group == nil {
				break
			}
			groupId := group.ID
			var groupDefaults *Group
			for _, g := range currentGroups {
				if g.ID == groupId {
					groupDefaults = &g
					break
				}
			}

			// set current nodeCount
			nodeCount := groupDefaults.Members.Allocation

			if group.Members != nil {
				err = checkGroupScaling(groupId, "members", group.Members.Allocation, groupDefaults.Members, 1)
				if err != nil {
					return err
				}
			}

			if group.Memory != nil {
				err = checkGroupScaling(groupId, "memory", group.Memory.Allocation, groupDefaults.Memory, nodeCount)
				if err != nil {
					return err
				}
			}

			if group.Disk != nil {
				err = checkGroupScaling(groupId, "disk", group.Disk.Allocation, groupDefaults.Disk, nodeCount)
				if err != nil {
					return err
				}
			}

			if group.CPU != nil {
				err = checkGroupScaling(groupId, "cpu", group.CPU.Allocation, groupDefaults.CPU, nodeCount)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Updates and creates users. Because we cannot get users, we first attempt to update the users, then create them
func userUpdateCreate(userData map[string]interface{}, instanceID string, meta interface{}, d *schema.ResourceData) (err error) {
	cloudDatabasesClient, _ := meta.(conns.ClientSession).CloudDatabasesV5()
	// Attempt to update user password
	passwordSettingUser := &clouddatabasesv5.APasswordSettingUser{
		Password: core.StringPtr(userData["password"].(string)),
	}

	changeUserPasswordOptions := &clouddatabasesv5.ChangeUserPasswordOptions{
		ID:       &instanceID,
		UserType: core.StringPtr(userData["type"].(string)),
		Username: core.StringPtr(userData["name"].(string)),
		User:     passwordSettingUser,
	}

	changeUserPasswordResponse, response, err := cloudDatabasesClient.ChangeUserPassword(changeUserPasswordOptions)

	// user was found but an error occurs while triggering task
	if response.StatusCode != 404 && err != nil {
		return fmt.Errorf("[ERROR] ChangeUserPassword (%s) failed %s\n%s", *changeUserPasswordOptions.Username, err, response)
	}

	taskID := *changeUserPasswordResponse.Task.ID
	updatePass, err := waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		log.Printf("[ERROR] Error waiting for database (%s) user (%s) password update task to complete: %s", instanceID, *changeUserPasswordOptions.Username, err)
	}

	// Updating the password has failed
	if !updatePass {
		//Attempt to create user
		userEntry := &clouddatabasesv5.User{
			Username: core.StringPtr(userData["name"].(string)),
			Password: core.StringPtr(userData["password"].(string)),
		}

		// User Role only for ops_manager user type
		if userData["type"].(string) == "ops_manager" && userData["role"].(string) != "" {
			userEntry.Role = core.StringPtr(userData["role"].(string))
		}

		createDatabaseUserOptions := &clouddatabasesv5.CreateDatabaseUserOptions{
			ID:       &instanceID,
			UserType: core.StringPtr(userData["type"].(string)),
			User:     userEntry,
		}

		createDatabaseUserResponse, response, err := cloudDatabasesClient.CreateDatabaseUser(createDatabaseUserOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] CreateDatabaseUser (%s) failed %s\n%s", *userEntry.Username, err, response)
		}

		taskID := *createDatabaseUserResponse.Task.ID
		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"[ERROR] Error waiting for database (%s) user (%s) create task to complete: %s", instanceID, *userEntry.Username, err)
		}
	}

	return nil
}
