// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
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

const (
	databaseUserSpecialChars   = "_-"
	opsManagerUserSpecialChars = "~!@#$%^&*()=+[]{}|;:,.<>/?_-"
)

const (
	redisRBACRoleRegexPattern = `[+-]@(?P<category>[a-z]+)`
)

type DatabaseUser struct {
	Username string
	Password string
	Role     *string
	Type     string
}

type databaseUserValidationError struct {
	user *DatabaseUser
	errs []error
}

func (e *databaseUserValidationError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}

	return fmt.Sprintf("database user (%s) validation error:\n%s", e.user.Username, string(b))
}

func (e *databaseUserValidationError) Unwrap() error {
	if e == nil || len(e.errs) == 0 {
		return nil
	}

	// only return the first
	return e.errs[0]
}

type userChange struct {
	Old, New *DatabaseUser
}

func redisRBACAllowedRoles() []string {
	return []string{"all", "admin", "read", "write"}
}

func opsManagerRoles() []string {
	return []string{"group_read_only", "group_data_access_admin"}
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
			validateGroupsDiff,
			validateUsersDiff,
			validateRemoteLeaderIDDiff),

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
				Description: "The admin user password for the instance",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(15, 72),
					DatabaseUserPasswordValidator("database"),
				),
				Sensitive: true,
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
			"service_endpoints": {
				Description:  "Types of the service endpoints. Possible values are 'public', 'private', 'public-and-private'.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_database", "service_endpoints"),
			},
			"backup_id": {
				Description: "The CRN of backup source database",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"remote_leader_id": {
				Description: "The CRN of leader database",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"skip_initial_backup": {
				Description: "Option to skip the initial backup when promoting a read-only replica. Skipping the initial backup means that your replica becomes available more quickly, but there is no immediate backup available.",
				Type:        schema.TypeBool,
				Optional:    true,
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
			"offline_restore": {
				Description:      "Set offline restore mode for MongoDB Enterprise Edition",
				Type:             schema.TypeBool,
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
							ValidateFunc: validation.StringLenBetween(15, 32),
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
							Description: "User role. Only available for ops_manager user type and Redis 6.0 and above.",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   false,
						},
					},
				},
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
				Type:     schema.TypeSet,
				Optional: true,
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
						"host_flavor": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"multitenant",
											"b3c.4x16.encrypted",
											"b3c.8x32.encrypted",
											"m3c.8x64.encrypted",
											"b3c.16x64.encrypted",
											"b3c.32x128.encrypted",
											"m3c.30x240.encrypted"}, false),
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
									"cpu_enforcement_ratio_ceiling_mb": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The amount of memory required before the cpu/memory ratio is no longer enforced. (multitenant only).",
									},
									"cpu_enforcement_ratio_mb": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum memory allowed per CPU until the ratio ceiling is reached. (multitenant only).",
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
						"host_flavor": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host flavor id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host flavor name",
									},
									"hosting_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The host flavor size",
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
			flex.DeletionProtection: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether Terraform will be prevented from destroying the instance",
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
			AllowedValues:              "databases-for-etcd, databases-for-postgresql, databases-for-redis, databases-for-elasticsearch, databases-for-mongodb, messages-for-rabbitmq, databases-for-mysql, databases-for-enterprisedb",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "plan",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "standard, enterprise, enterprise-sharding, platinum",
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_endpoints",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "public, private, public-and-private",
			Required:                   false})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "group_id",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			AllowedValues:              "member, analytics, bi_connector",
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
	HostFlavor          string  `json:"members_host_flavor,omitempty"`
	KeyProtectInstance  string  `json:"disk_encryption_instance_crn,omitempty"`
	ServiceEndpoints    string  `json:"service-endpoints,omitempty"`
	BackupID            string  `json:"backup-id,omitempty"`
	RemoteLeaderID      string  `json:"remote_leader_id,omitempty"`
	PITRDeploymentID    string  `json:"point_in_time_recovery_deployment_id,omitempty"`
	PITRTimeStamp       *string `json:"point_in_time_recovery_time,omitempty"`
	OfflineRestore      bool    `json:"offline_restore,omitempty"`
}

type Group struct {
	ID         string
	Members    *GroupResource
	Memory     *GroupResource
	Disk       *GroupResource
	CPU        *GroupResource
	HostFlavor *HostFlavorGroupResource
}

type GroupResource struct {
	Units                        string
	Allocation                   int
	Minimum                      int
	Maximum                      int
	StepSize                     int
	IsAdjustable                 bool
	IsOptional                   bool
	CanScaleDown                 bool
	CPUEnforcementRatioCeilingMb int
	CPUEnforcementRatioMb        int
}

type HostFlavorGroupResource struct {
	ID string
}

func getDefaultScalingGroups(_service string, _plan string, _hostFlavor string, meta interface{}) (groups []clouddatabasesv5.Group, err error) {
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

	if service == "mongodb" && _plan == "enterprise" {
		service = "mongodbee"
	}

	if service == "elasticsearch" && _plan == "platinum" {
		service = "elasticsearchpl"
	}

	if service == "mongodb" && _plan == "enterprise-sharding" {
		service = "mongodbees"
	}

	getDefaultScalingGroupsOptions := cloudDatabasesClient.NewGetDefaultScalingGroupsOptions(service)
	if _hostFlavor != "" {
		getDefaultScalingGroupsOptions.SetHostFlavor(_hostFlavor)
	}

	getDefaultScalingGroupsResponse, response, err := cloudDatabasesClient.GetDefaultScalingGroups(getDefaultScalingGroupsOptions)
	if err != nil {
		if response.StatusCode == 422 {
			return groups, fmt.Errorf("%s is not available on multitenant", service)
		}
		return groups, err
	}

	return getDefaultScalingGroupsResponse.Groups, nil
}

func getInitialNodeCount(service string, plan string, meta interface{}) (int, error) {
	groups, err := getDefaultScalingGroups(service, plan, "", meta)

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

func resourceIBMDatabaseInstanceDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	err = flex.ResourceTagsCustomizeDiff(diff)
	if err != nil {
		return err
	}

	service := diff.Get("service").(string)
	plan := diff.Get("plan").(string)

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

	_, offlineRestoreOk := diff.GetOk("offline_restore")
	if offlineRestoreOk && service != "databases-for-mongodb" && plan != "enterprise" {
		return fmt.Errorf("[ERROR] offline_restore is only supported for databases-for-mongodb enterprise")
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
		if serviceName == "databases-for-mongodb" && plan == "enterprise-sharding" {
			return diag.FromErr(fmt.Errorf("%s %s is not available yet in this region", serviceName, plan))
		} else {
			return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err))
		}
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

	params := Params{}

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

	if offlineRestore, ok := d.GetOk("offline_restore"); ok {
		params.OfflineRestore = offlineRestore.(bool)
	}

	var initialNodeCount int
	var sourceCRN string

	if params.PITRDeploymentID != "" {
		sourceCRN = params.PITRDeploymentID
	}

	if params.RemoteLeaderID != "" {
		sourceCRN = params.RemoteLeaderID
	}

	if sourceCRN != "" {
		group, err := getMemberGroup(sourceCRN, meta)

		if err != nil {
			return diag.FromErr(
				fmt.Errorf("[ERROR] Error fetching source formation group: %s", err)) // raise error
		}

		if group != nil {
			initialNodeCount = group.Members.Allocation
		}
	} else {
		initialNodeCount, err = getInitialNodeCount(serviceName, plan, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var memberGroup *Group

	if group, ok := d.GetOk("group"); ok {
		groups := expandGroups(group.(*schema.Set).List())

		for _, g := range groups {
			if g.ID == "member" {
				memberGroup = g
				break
			}
		}
	}

	if memberGroup != nil && memberGroup.Memory != nil {
		params.Memory = memberGroup.Memory.Allocation * initialNodeCount
	}

	if memberGroup != nil && memberGroup.Disk != nil {
		params.Disk = memberGroup.Disk.Allocation * initialNodeCount
	}

	if memberGroup != nil && memberGroup.CPU != nil {
		params.CPU = memberGroup.CPU.Allocation * initialNodeCount
	}

	if memberGroup != nil && memberGroup.HostFlavor != nil {
		params.HostFlavor = memberGroup.HostFlavor.ID
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
			if g.HostFlavor != nil {
				groupScaling.HostFlavor = &clouddatabasesv5.GroupScalingHostFlavor{ID: core.StringPtr(g.HostFlavor.ID)}
			}

			if groupScaling.Members != nil || groupScaling.Memory != nil || groupScaling.Disk != nil || groupScaling.CPU != nil || groupScaling.HostFlavor != nil {
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

		user := &clouddatabasesv5.UserUpdatePasswordSetting{
			Password: &adminPassword,
		}

		updateUserOptions := &clouddatabasesv5.UpdateUserOptions{
			ID:       core.StringPtr(instanceID),
			UserType: core.StringPtr("database"),
			Username: core.StringPtr(adminUser),
			User:     user,
		}

		updateUserResponse, response, err := cloudDatabasesClient.UpdateUser(updateUserOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] UpdateUser (%s) failed %s\n%s", *updateUserOptions.Username, err, response))
		}

		taskID := *updateUserResponse.Task.ID
		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutCreate))

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database admin password: %s", err))
		}
	}

	_, hasAllowlist := d.GetOk("allowlist")

	if hasAllowlist {
		var ipAddresses *schema.Set

		ipAddresses = d.Get("allowlist").(*schema.Set)

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

		users := expandUsers(userList.(*schema.Set).List())
		for _, user := range users {
			// Note: Some db users exist after provisioning (i.e. admin, repl)
			// so we must attempt both methods
			err := user.Update(instanceID, d, meta)

			if err != nil {
				err = user.Create(instanceID, d, meta)
			}

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

	listDeploymentScalingGroupsOptions := &clouddatabasesv5.ListDeploymentScalingGroupsOptions{
		ID: core.StringPtr(instanceID),
	}
	groupList, _, err := cloudDatabasesClient.ListDeploymentScalingGroups(listDeploymentScalingGroupsOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database groups: %s", err))
	}
	if len(groupList.Groups) == 0 || groupList.Groups[0].Members == nil || groupList.Groups[0].Members.AllocationCount == nil || *groupList.Groups[0].Members.AllocationCount == 0 {
		return diag.FromErr(fmt.Errorf("[ERROR] This database appears to have have 0 members. Unable to proceed"))
	}

	d.Set("groups", flex.FlattenIcdGroups(groupList))

	getAutoscalingConditionsOptions := &clouddatabasesv5.GetAutoscalingConditionsOptions{
		ID:      instance.ID,
		GroupID: core.StringPtr("member"),
	}

	autoscalingGroup, _, err := cloudDatabasesClient.GetAutoscalingConditions(getAutoscalingConditionsOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database autoscaling groups: %s\n Hint: Check if there is a mismatch between your database location and IBMCLOUD_REGION", err))
	}
	d.Set("auto_scaling", flattenAutoScalingGroup(*autoscalingGroup))

	alEntry := &clouddatabasesv5.GetAllowlistOptions{
		ID: &instanceID,
	}

	allowlist, _, err := cloudDatabasesClient.GetAllowlist(alEntry)

	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database allowlist: %s", err))
	}

	d.Set("allowlist", flex.FlattenAllowlist(allowlist.IPAddresses))

	//ICD does not implement a GetUsers API. Users populated from tf configuration.
	tfusers := d.Get("users").(*schema.Set)
	users := flex.ExpandUsers(tfusers)
	user := icdv4.User{
		UserName: deployment.AdminUsernames["database"],
	}
	users = append(users, user)

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

	// This can be removed any time after August once all old multitenant instances are switched over to the new multitenant
	if groupList.Groups[0].HostFlavor == nil && (groupList.Groups[0].CPU != nil && *groupList.Groups[0].CPU.AllocationCount == 0) {
		return appendSwitchoverWarning()
	}

	endpoint, _ := instance.Parameters["service-endpoints"]
	if endpoint == "public" || endpoint == "public-and-private" {
		return publicServiceEndpointsWarning()
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

	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting database client settings: %s", err))
	}
	icdId := flex.EscapeUrlParm(instanceID)

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
			if group.HostFlavor != nil {
				groupScaling.HostFlavor = &clouddatabasesv5.GroupScalingHostFlavor{ID: core.StringPtr(group.HostFlavor.ID)}
			}

			if groupScaling.Members != nil || groupScaling.Memory != nil || groupScaling.Disk != nil || groupScaling.CPU != nil || groupScaling.HostFlavor != nil {
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

		user := &clouddatabasesv5.UserUpdatePasswordSetting{
			Password: &password,
		}

		updateUserOptions := &clouddatabasesv5.UpdateUserOptions{
			ID:       core.StringPtr(instanceID),
			UserType: core.StringPtr("database"),
			Username: core.StringPtr(adminUser),
			User:     user,
		}

		updateUserResponse, response, err := cloudDatabasesClient.UpdateUser(updateUserOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] UpdateUser (%s) failed %s\n%s", *updateUserOptions.Username, err, response))
		}

		taskID := *updateUserResponse.Task.ID
		_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating database admin password: %s", err))
		}
	}

	if d.HasChange("allowlist") {
		_, hasAllowlist := d.GetOk("allowlist")

		var entries interface{}

		if hasAllowlist {
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
				"[ERROR] Error waiting for update of database (%s) allowlist task to complete: %s", instanceID, err))
		}
	}

	if d.HasChange("users") {
		oldUsers, newUsers := d.GetChange("users")
		userChanges := expandUserChanges(oldUsers.(*schema.Set).List(), newUsers.(*schema.Set).List())

		for _, change := range userChanges {
			// Delete User
			if change.isDelete() {
				// Delete Old User
				err = change.Old.Delete(instanceID, d, meta)

				if err != nil {
					return diag.FromErr(err)
				}
			}

			if change.isCreate() || change.isUpdate() {

				// Note: User Update is not supported for ops_manager user type
				// Delete (ignoring errors), then re-create
				if change.isUpdate() && !change.New.isUpdatable() {
					change.Old.Delete(instanceID, d, meta)

					err = change.New.Create(instanceID, d, meta)
				} else {
					// Note: Some db users exist after provisioning (i.e. admin, repl)
					// so we must attempt both methods
					err = change.New.Update(instanceID, d, meta)

					// Create User if Update failed
					if err != nil {
						err = change.New.Create(instanceID, d, meta)
					}
				}

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
				deleteLogicalReplicationSlotOptions := &clouddatabasesv5.DeleteLogicalReplicationSlotOptions{
					ID:   &instanceID,
					Name: core.StringPtr(newEntry["name"].(string)),
				}

				deleteLogicalReplicationSlotResponse, response, err := cloudDatabasesClient.DeleteLogicalReplicationSlot(deleteLogicalReplicationSlotOptions)

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] DeleteLogicalReplicationSlot (%s) failed %s\n%s", *deleteLogicalReplicationSlotOptions.Name, err, response))
				}

				taskID := *deleteLogicalReplicationSlotResponse.Task.ID
				_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

				if err != nil {
					return diag.FromErr(fmt.Errorf(
						"[ERROR] Error waiting for database (%s) logical replication slot (%s) delete task to complete: %s", icdId, *deleteLogicalReplicationSlotOptions.Name, err))
				}
			}
		}
	}

	if d.HasChange("remote_leader_id") {
		remoteLeaderId := d.Get("remote_leader_id").(string)

		if remoteLeaderId == "" {
			skipInitialBackup := false
			if skip, ok := d.GetOk("skip_initial_backup"); ok {
				skipInitialBackup = skip.(bool)
			}

			promoteReadOnlyReplicaOptions := &clouddatabasesv5.PromoteReadOnlyReplicaOptions{
				ID: &instanceID,
				Promotion: map[string]interface{}{
					"skip_initial_backup": skipInitialBackup,
				},
			}

			promoteReadReplicaResponse, response, err := cloudDatabasesClient.PromoteReadOnlyReplica(promoteReadOnlyReplicaOptions)

			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error promoting read replica: %s\n%s", err, response))
			}

			taskID := *promoteReadReplicaResponse.Task.ID
			_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error promoting read replica: %s", err))
			}
		}
	}

	return resourceIBMDatabaseInstanceRead(context, d, meta)
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
			return fmt.Errorf("[ERROR] Error getting database config for: %s with error %s\n", icdId, cdbErr)
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
	ticker := time.NewTicker(delayDuration)
	delay := ticker.C
	defer ticker.Stop()

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
			Units:                        *g.Memory.Units,
			Allocation:                   int(*g.Memory.AllocationMb),
			Minimum:                      int(*g.Memory.MinimumMb),
			Maximum:                      int(*g.Memory.MaximumMb),
			StepSize:                     int(*g.Memory.StepSizeMb),
			IsAdjustable:                 *g.Memory.IsAdjustable,
			IsOptional:                   *g.Memory.IsOptional,
			CanScaleDown:                 *g.Memory.CanScaleDown,
			CPUEnforcementRatioCeilingMb: getCPUEnforcementRatioCeilingMb(g.Memory),
			CPUEnforcementRatioMb:        getCPUEnforcementRatioMb(g.Memory),
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

		group.HostFlavor = &HostFlavorGroupResource{}
		if g.HostFlavor != nil {
			group.HostFlavor.ID = *g.HostFlavor.ID
		}

		groups = append(groups, group)
	}

	return groups
}

func getCPUEnforcementRatioCeilingMb(groupMemory *clouddatabasesv5.GroupMemory) int {
	if groupMemory.CPUEnforcementRatioCeilingMb != nil {
		return int(*groupMemory.CPUEnforcementRatioCeilingMb)
	}

	return 0
}

func getCPUEnforcementRatioMb(groupMemory *clouddatabasesv5.GroupMemory) int {
	if groupMemory.CPUEnforcementRatioMb != nil {
		return int(*groupMemory.CPUEnforcementRatioMb)
	}

	return 0
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

			if hostflavorSet, ok := tfGroup["host_flavor"].(*schema.Set); ok {
				hostflavor := hostflavorSet.List()
				if len(hostflavor) != 0 {
					group.HostFlavor = &HostFlavorGroupResource{ID: hostflavor[0].(map[string]interface{})["id"].(string)}
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

func validateGroupScaling(groupId string, resourceName string, value int, resource *GroupResource, nodeCount int) error {
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

func validateGroupHostFlavor(groupId string, resourceName string, group *Group) error {
	if group.CPU != nil || group.Memory != nil {
		return fmt.Errorf("%s must not be set with cpu and memory", resourceName)
	}
	return nil
}

func validateMultitenantMemoryCpu(resourceDefaults *Group, group *Group, cpuEnforcementRatioCeiling int, cpuEnforcementRatioMb int) error {
	// TODO: Replace this with  cpuEnforcementRatioCeiling when it is fixed
	cpuEnforcementRatioCeilingTemp := 16384

	if group.CPU == nil || group.CPU.Allocation == 0 {
		return nil
	}

	if group.CPU.Allocation >= cpuEnforcementRatioCeilingTemp/cpuEnforcementRatioMb {
		return nil
	} else {
		return fmt.Errorf("The current cpu allocation of %d is not valid for your current configuration.", group.CPU.Allocation)
	}
}

// This can be removed any time after August once all old multitenant instances are switched over to the new multitenant
func appendSwitchoverWarning() diag.Diagnostics {
	var diags diag.Diagnostics

	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Note: IBM Cloud Databases released new Hosting Models on May 1. All existing multi-tenant instances will have their resources adjusted to Shared Compute allocations during August 2024. To monitor your current resource needs, and learn about how the transition to Shared Compute will impact your instance, see our documentation https://cloud.ibm.com/docs/cloud-databases?topic=cloud-databases-hosting-models",
	}

	diags = append(diags, warning)

	return diags
}

func publicServiceEndpointsWarning() diag.Diagnostics {
	var diags diag.Diagnostics

	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "IBM recommends using private endpoints only to improve security by restricting access to your database to the IBM Cloud private network. For more information, please refer to our security best practices, https://cloud.ibm.com/docs/cloud-databases?topic=cloud-databases-manage-security-compliance.",
	}

	diags = append(diags, warning)

	return diags
}

func validateGroupsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	instanceID := diff.Id()
	service := diff.Get("service").(string)
	plan := diff.Get("plan").(string)

	if group, ok := diff.GetOk("group"); ok {
		var currentGroups []Group
		var groupList []clouddatabasesv5.Group
		var groupIds []string
		groups := expandGroups(group.(*schema.Set).List())
		var memberGroup *Group
		for _, g := range groups {
			if g.ID == "member" {
				memberGroup = g
				break
			}
		}

		if instanceID != "" {
			groupList, err = getGroups(instanceID, meta)
		} else {
			if memberGroup.HostFlavor != nil {
				groupList, err = getDefaultScalingGroups(service, plan, memberGroup.HostFlavor.ID, meta)
			} else {
				groupList, err = getDefaultScalingGroups(service, plan, "", meta)
			}
		}

		if err != nil {
			return err
		}

		currentGroups = normalizeGroups(groupList)

		tfGroups := expandGroups(group.(*schema.Set).List())

		cpuEnforcementRatioCeiling, cpuEnforcementRatioMb := 0, 0

		if memberGroup.HostFlavor != nil && memberGroup.HostFlavor.ID == "multitenant" {
			err, cpuEnforcementRatioCeiling, cpuEnforcementRatioMb = getCpuEnforcementRatios(service, plan, memberGroup.HostFlavor.ID, meta, group)

			if err != nil {
				return err
			}
		}

		// validate group_ids are unique
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

			if group.HostFlavor != nil && group.HostFlavor.ID != "" && group.HostFlavor.ID != "multitenant" {
				err = validateGroupHostFlavor(groupId, "host_flavor", group)
				if err != nil {
					return err
				}
			}

			if group.Members != nil {
				err = validateGroupScaling(groupId, "members", group.Members.Allocation, groupDefaults.Members, 1)
				if err != nil {
					return err
				}
			}

			if group.Memory != nil {
				err = validateGroupScaling(groupId, "memory", group.Memory.Allocation, groupDefaults.Memory, nodeCount)
				if err != nil {
					return err
				}

				if group.HostFlavor != nil && group.HostFlavor.ID != "" && group.HostFlavor.ID == "multitenant" && cpuEnforcementRatioCeiling != 0 && cpuEnforcementRatioMb != 0 {
					err = validateMultitenantMemoryCpu(groupDefaults, group, cpuEnforcementRatioCeiling, cpuEnforcementRatioMb)
					if err != nil {
						return err
					}
				}
			}

			if group.Disk != nil {
				err = validateGroupScaling(groupId, "disk", group.Disk.Allocation, groupDefaults.Disk, nodeCount)
				if err != nil {
					return err
				}
			}

			if group.CPU != nil {
				err = validateGroupScaling(groupId, "cpu", group.CPU.Allocation, groupDefaults.CPU, nodeCount)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getCpuEnforcementRatios(service string, plan string, hostFlavor string, meta interface{}, group interface{}) (err error, cpuEnforcementRatioCeiling int, cpuEnforcementRatioMb int) {
	var currentGroups []Group
	defaultList, err := getDefaultScalingGroups(service, plan, hostFlavor, meta)

	if err != nil {
		return err, 0, 0
	}

	currentGroups = normalizeGroups(defaultList)
	tfGroups := expandGroups(group.(*schema.Set).List())

	// Get default or current group scaling values
	for _, group := range tfGroups {
		if group == nil {
			continue
		}

		groupId := group.ID
		var groupDefaults *Group
		for _, g := range currentGroups {
			if g.ID == groupId {
				groupDefaults = &g
				break
			}
		}

		if groupDefaults.Memory != nil {
			return nil, groupDefaults.Memory.CPUEnforcementRatioCeilingMb, groupDefaults.Memory.CPUEnforcementRatioMb
		}

	}

	return nil, 0, 0
}

func getMemberGroup(instanceCRN string, meta interface{}) (*Group, error) {
	groupsResponse, err := getGroups(instanceCRN, meta)

	if err != nil {
		return nil, err
	}

	currentGroups := normalizeGroups(groupsResponse)

	for _, cg := range currentGroups {
		if cg.ID == "member" {
			return &cg, nil
		}
	}

	return nil, nil
}

func validateUsersDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	service := diff.Get("service").(string)

	var versionStr string
	var version int

	if _version, ok := diff.GetOk("version"); ok {
		versionStr = _version.(string)
	}

	if versionStr == "" {
		// Latest Version
		version = 0
	} else {
		_v, err := strconv.ParseFloat(versionStr, 64)

		if err != nil {
			return fmt.Errorf("invalid version: %s", versionStr)
		}

		version = int(_v)
	}

	oldUsers, newUsers := diff.GetChange("users")
	userChanges := expandUserChanges(oldUsers.(*schema.Set).List(), newUsers.(*schema.Set).List())

	for _, change := range userChanges {
		if change.isDelete() {
			continue
		}

		if change.isCreate() || change.isUpdate() {
			err = change.New.ValidatePassword()

			if err != nil {
				return err
			}

			// TODO: Use Capability API
			// RBAC roles supported for Redis 6.0 and above
			if (service == "databases-for-redis") && !(version > 0 && version < 6) {
				err = change.New.ValidateRBACRole()
			} else if service == "databases-for-mongodb" && change.New.Type == "ops_manager" {
				err = change.New.ValidateOpsManagerRole()
			} else {
				if change.New.Role != nil {
					if *change.New.Role != "" {
						err = errors.New("role is not supported for this deployment or user type")
						err = &databaseUserValidationError{user: change.New, errs: []error{err}}
					}
				}
			}

			if err != nil {
				return err
			}
		}
	}

	return
}

func expandUsers(_users []interface{}) []*DatabaseUser {
	if len(_users) == 0 {
		return nil
	}

	users := make([]*DatabaseUser, 0, len(_users))

	for _, userRaw := range _users {
		if tfUser, ok := userRaw.(map[string]interface{}); ok {

			user := DatabaseUser{
				Username: tfUser["name"].(string),
				Password: tfUser["password"].(string),
				Type:     tfUser["type"].(string),
			}

			// NOTE: cannot differentiate nil vs empty string
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/741
			if role, ok := tfUser["role"].(string); ok {
				if tfUser["role"] != "" {
					user.Role = &role
				}
			}

			users = append(users, &user)
		}
	}

	return users
}

func expandUserChanges(_oldUsers []interface{}, _newUsers []interface{}) (userChanges []*userChange) {
	oldUsers := expandUsers(_oldUsers)
	newUsers := expandUsers(_newUsers)

	userChangeMap := make(map[string]*userChange)

	for _, user := range oldUsers {
		userChangeMap[user.ID()] = &userChange{Old: user}
	}

	for _, user := range newUsers {
		if _, ok := userChangeMap[user.ID()]; !ok {
			userChangeMap[user.ID()] = &userChange{}
		}
		userChangeMap[user.ID()].New = user
	}

	userChanges = make([]*userChange, 0, len(userChangeMap))

	for _, change := range userChangeMap {
		userChanges = append(userChanges, change)
	}

	return userChanges
}

func validateRemoteLeaderIDDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) (err error) {
	_, remoteLeaderIdOk := diff.GetOk("remote_leader_id")
	service := diff.Get("service").(string)
	crn := diff.Get("resource_crn").(string)

	if remoteLeaderIdOk && (service != "databases-for-postgresql" && service != "databases-for-mysql" && service != "databases-for-enterprisedb") {
		return fmt.Errorf("[ERROR] remote_leader_id is only supported for databases-for-postgresql, databases-for-enterprisedb and databases-for-mysql")
	}

	oldValue, newValue := diff.GetChange("remote_leader_id")

	if crn != "" && oldValue == "" && newValue != "" {
		return fmt.Errorf("[ERROR] You cannot convert an existing instance to a read replica")
	}

	return nil
}

func (c *userChange) isDelete() bool {
	return c.Old != nil && c.New == nil
}

func (c *userChange) isCreate() bool {
	return c.Old == nil && c.New != nil
}

func (c *userChange) isUpdate() bool {
	return c.New != nil &&
		c.Old != nil &&
		((c.Old.Password != c.New.Password) ||
			(c.Old.Role != c.New.Role))
}

func (u *DatabaseUser) ID() (id string) {
	return fmt.Sprintf("%s-%s", u.Type, u.Username)
}

func (u *DatabaseUser) Create(instanceID string, d *schema.ResourceData, meta interface{}) (err error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %w", err)
	}

	//Attempt to create user
	userEntry := &clouddatabasesv5.User{
		Username: core.StringPtr(u.Username),
		Password: core.StringPtr(u.Password),
	}

	// User Role only for ops_manager user type and Redis 6.0 and above
	if u.Role != nil {
		userEntry.Role = u.Role
	}

	createDatabaseUserOptions := &clouddatabasesv5.CreateDatabaseUserOptions{
		ID:       &instanceID,
		UserType: core.StringPtr(u.Type),
		User:     userEntry,
	}

	createDatabaseUserResponse, response, err := cloudDatabasesClient.CreateDatabaseUser(createDatabaseUserOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] CreateDatabaseUser (%s) failed %w\n%s", *userEntry.Username, err, response)
	}

	taskID := *createDatabaseUserResponse.Task.ID
	_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf(
			"[ERROR] Error waiting for database (%s) user (%s) create task to complete: %w", instanceID, *userEntry.Username, err)
	}

	return nil
}

func (u *DatabaseUser) Update(instanceID string, d *schema.ResourceData, meta interface{}) (err error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	// Attempt to update user password
	user := &clouddatabasesv5.UserUpdate{
		Password: core.StringPtr(u.Password),
	}

	if u.Role != nil {
		user.Role = u.Role
	}

	updateUserOptions := &clouddatabasesv5.UpdateUserOptions{
		ID:       &instanceID,
		UserType: core.StringPtr(u.Type),
		Username: core.StringPtr(u.Username),
		User:     user,
	}

	updateUserResponse, response, err := cloudDatabasesClient.UpdateUser(updateUserOptions)

	// user was found but an error occurs while triggering task
	if err != nil || (response.StatusCode < 200 || response.StatusCode >= 300) {
		return fmt.Errorf("[ERROR] UpdateUser (%s) failed %w\n%s", *updateUserOptions.Username, err, response)
	}

	taskID := *updateUserResponse.Task.ID
	_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf(
			"[ERROR] Error waiting for database (%s) user (%s) create task to complete: %w", instanceID, *updateUserOptions.Username, err)
	}

	return nil
}

func (u *DatabaseUser) Delete(instanceID string, d *schema.ResourceData, meta interface{}) (err error) {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting database client settings: %s", err)
	}

	deleteDatabaseUserOptions := &clouddatabasesv5.DeleteDatabaseUserOptions{
		ID:       &instanceID,
		UserType: core.StringPtr(u.Type),
		Username: core.StringPtr(u.Username),
	}

	deleteDatabaseUserResponse, response, err := cloudDatabasesClient.DeleteDatabaseUser(deleteDatabaseUserOptions)

	if err != nil {
		return fmt.Errorf(
			"[ERROR] DeleteDatabaseUser (%s) failed %s\n%s", *deleteDatabaseUserOptions.Username, err, response)

	}

	taskID := *deleteDatabaseUserResponse.Task.ID
	_, err = waitForDatabaseTaskComplete(taskID, d, meta, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf(
			"[ERROR] Error waiting for database (%s) user (%s) delete task to complete: %s", instanceID, *deleteDatabaseUserOptions.Username, err)
	}

	return nil
}

func (u *DatabaseUser) isUpdatable() bool {
	return u.Type != "ops_manager"
}

func (u *DatabaseUser) ValidatePassword() (err error) {
	var errs []error

	var specialChars string
	switch u.Type {
	case "ops_manager":
		specialChars = opsManagerUserSpecialChars
	default:
		specialChars = databaseUserSpecialChars
	}

	// Format for regexp
	var specialCharPattern string
	var bs strings.Builder
	for i, c := range strings.Split(specialChars, "") {
		if i > 0 {
			bs.WriteByte('|')
		}
		bs.WriteString(regexp.QuoteMeta(c))
	}

	specialCharPattern = bs.String()

	var allowedCharacters = regexp.MustCompile(fmt.Sprintf("^(?:[a-zA-Z0-9]|%s)+$", specialCharPattern))
	var beginWithSpecialChar = regexp.MustCompile(fmt.Sprintf("^(?:%s)", specialCharPattern))
	var containsLower = regexp.MustCompile("[a-z]")
	var containsUpper = regexp.MustCompile("[A-Z]")
	var containsNumber = regexp.MustCompile("[0-9]")
	var containsSpecialChar = regexp.MustCompile(fmt.Sprintf("(?:%s)", specialCharPattern))

	if u.Type == "ops_manager" && !containsSpecialChar.MatchString(u.Password) {
		errs = append(errs, fmt.Errorf(
			"password must contain at least one special character (%s)", specialChars))
	}

	if u.Type == "database" && beginWithSpecialChar.MatchString(u.Password) {
		errs = append(errs, fmt.Errorf(
			"password must not begin with a special character (%s)", specialChars))
	}

	if !containsLower.MatchString(u.Password) {
		errs = append(errs, errors.New("password must contain at least one lower case letter"))
	}

	if !containsUpper.MatchString(u.Password) {
		errs = append(errs, errors.New("password must contain at least one upper case letter"))
	}

	if !containsNumber.MatchString(u.Password) {
		errs = append(errs, errors.New("password must contain at least one number"))
	}

	if !allowedCharacters.MatchString(u.Password) {
		errs = append(errs, errors.New("password must not contain invalid characters"))
	}

	if len(errs) == 0 {
		return
	}

	return &databaseUserValidationError{user: u, errs: errs}
}

func (u *DatabaseUser) ValidateRBACRole() (err error) {
	var errs []error

	if u.Role == nil || *u.Role == "" {
		return
	}

	if u.Type != "database" {
		errs = append(errs, errors.New("role is only allowed for the database user"))
		return &databaseUserValidationError{user: u, errs: errs}
	}

	redisRBACCategoryRegex := regexp.MustCompile(redisRBACRoleRegexPattern)
	redisRBACRoleRegex := regexp.MustCompile(fmt.Sprintf(`^(%s\s?)+$`, redisRBACRoleRegexPattern))

	if !redisRBACRoleRegex.MatchString(*u.Role) {
		errs = append(errs, errors.New("role must be in the format +@category or -@category"))
	}

	matches := redisRBACCategoryRegex.FindAllStringSubmatch(*u.Role, -1)

	for _, match := range matches {
		valid := false
		role := match[1]
		for _, allowed := range redisRBACAllowedRoles() {
			if role == allowed {
				valid = true
				break
			}
		}

		if !valid {
			errs = append(errs, fmt.Errorf("role must contain only allowed categories: %s", strings.Join(redisRBACAllowedRoles()[:], ",")))
			break
		}
	}

	if len(errs) == 0 {
		return
	}

	return &databaseUserValidationError{user: u, errs: errs}
}

func (u *DatabaseUser) ValidateOpsManagerRole() (err error) {
	if u.Role == nil {
		return
	}

	if u.Type != "ops_manager" {
		return
	}

	if *u.Role == "" {
		return
	}

	for _, str := range opsManagerRoles() {
		if *u.Role == str {
			return
		}
	}

	err = fmt.Errorf("role must be a valid ops_manager role: %s", strings.Join(opsManagerRoles()[:], ","))

	return &databaseUserValidationError{user: u, errs: []error{err}}
}

func DatabaseUserPasswordValidator(userType string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		user := &DatabaseUser{Username: "admin", Type: userType, Password: i.(string)}
		err := user.ValidatePassword()
		if err != nil {
			errors = append(errors, err)
		}
		return
	}
}
