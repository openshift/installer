// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	//	"github.com/IBM-Cloud/bluemix-go/api/globaltagging/globaltaggingv3"
	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"
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

type CsEntry struct {
	Name       string
	Password   string
	String     string
	Composed   string
	CertName   string
	CertBase64 string
	Hosts      []struct {
		HostName string `json:"hostname"`
		Port     int    `json:"port"`
	}
	Scheme       string
	QueryOptions map[string]interface{}
	Path         string
	Database     string
}

func resourceIBMDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create:        resourceIBMDatabaseInstanceCreate,
		Read:          resourceIBMDatabaseInstanceRead,
		Update:        resourceIBMDatabaseInstanceUpdate,
		Delete:        resourceIBMDatabaseInstanceDelete,
		Exists:        resourceIBMDatabaseInstanceExists,
		CustomizeDiff: resourceIBMDatabaseInstanceDiff,
		Importer:      &schema.ResourceImporter{},

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
			},

			"location": {
				Description: "The location or the region in which Database instance exists",
				Type:        schema.TypeString,
				Required:    true,
			},

			"service": {
				Description:  "The name of the Cloud Internet database service",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"databases-for-etcd", "databases-for-postgresql", "databases-for-redis", "databases-for-elasticsearch", "databases-for-mongodb", "messages-for-rabbitmq", "databases-for-mysql"}),
			},
			"plan": {
				Description:  "The plan type of the Database instance",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"standard"}),
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
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count"},
			},
			"members_disk_allocation_mb": {
				Description:   "Disk allocation required for cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count"},
			},
			"members_cpu_allocation_count": {
				Description:   "CPU allocation required for cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"node_count", "node_memory_allocation_mb", "node_disk_allocation_mb", "node_cpu_allocation_count"},
			},
			"node_count": {
				Description:   "Total number of nodes in the cluster",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count"},
			},
			"node_memory_allocation_mb": {
				Description: "Memory allocation per node",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,

				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count"},
			},
			"node_disk_allocation_mb": {
				Description:   "Disk allocation per node",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count"},
			},
			"node_cpu_allocation_count": {
				Description:   "CPU allocation per node",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"members_memory_allocation_mb", "members_disk_allocation_mb", "members_cpu_allocation_count"},
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
				ValidateFunc: validateAllowedStringValue([]string{"public", "private", "public-and-private"}),
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
				DiffSuppressFunc: applyOnce,
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
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_database", "tag")},
				Set:      resourceIBMVPCHash,
			},
			"point_in_time_recovery_deployment_id": {
				Description:      "The CRN of source instance",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
			},
			"point_in_time_recovery_time": {
				Description:      "The point in time recovery time stamp of the deployed instance",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description:  "User name",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(5, 32),
						},
						"password": {
							Description:  "User password",
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(10, 32),
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
							ValidateFunc: validateCIDR,
						},
						"description": {
							Description:  "Unique white list description",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
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
			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}
func resourceIBMICDValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmICDResourceValidator := ResourceValidator{ResourceName: "ibm_database", Schema: validateSchema}
	return &ibmICDResourceValidator
}

type Params struct {
	Version             string `json:"version,omitempty"`
	KeyProtectKey       string `json:"key_protect_key,omitempty"`
	BackUpEncryptionCRN string `json:"backup_encryption_key_crn,omitempty"`
	Memory              int    `json:"members_memory_allocation_mb,omitempty"`
	Disk                int    `json:"members_disk_allocation_mb,omitempty"`
	CPU                 int    `json:"members_cpu_allocation_count,omitempty"`
	KeyProtectInstance  string `json:"key_protect_instance,omitempty"`
	ServiceEndpoints    string `json:"service-endpoints,omitempty"`
	BackupID            string `json:"backup-id,omitempty"`
	RemoteLeaderID      string `json:"remote_leader_id,omitempty"`
	PITRDeploymentID    string `json:"point_in_time_recovery_deployment_id,omitempty"`
	PITRTimeStamp       string `json:"point_in_time_recovery_time,omitempty"`
}

func getDatabaseServiceDefaults(service string, meta interface{}) (*icdv4.Group, error) {
	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return nil, fmt.Errorf("Error getting database client settings: %s", err)
	}

	var dbType string
	if strings.HasPrefix(service, "messages-for-") {
		dbType = service[len("messages-for-"):]
	} else {
		dbType = service[len("databases-for-"):]
	}

	groupDefaults, err := icdClient.Groups().GetDefaultGroups(dbType)
	if err != nil {
		return nil, fmt.Errorf("ICD API is down for plan validation, set plan_validation=false")
	}
	return &groupDefaults.Groups[0], nil
}

func getInitialNodeCount(d *schema.ResourceData, meta interface{}) (int, error) {
	service := d.Get("service").(string)
	planPhase := d.Get("plan_validation").(bool)
	if planPhase {
		groupDefaults, err := getDatabaseServiceDefaults(service, meta)
		if err != nil {
			return 0, err
		}
		return groupDefaults.Members.MinimumCount, nil
	} else {
		if service == "databases-for-elasticsearch" {
			return 3, nil
		}
		return 2, nil
	}
}

type GroupLimit struct {
	Units        string
	Allocation   int
	Minimum      int
	Maximum      int
	StepSize     int
	IsAdjustable bool
	CanScaleDown bool
}

func checkGroupValue(name string, limits GroupLimit, divider int, diff *schema.ResourceDiff) error {
	if diff.HasChange(name) {
		oldSetting, newSetting := diff.GetChange(name)
		old := oldSetting.(int)
		new := newSetting.(int)

		if new < limits.Minimum/divider || new > limits.Maximum/divider || new%(limits.StepSize/divider) != 0 {
			return fmt.Errorf("%s must be >= %d and <= %d in increments of %d", name, limits.Minimum/divider, limits.Maximum/divider/divider, limits.StepSize/divider)
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
	CanScaleDown    bool
}

func checkCountValue(name string, limits CountLimit, divider int, diff *schema.ResourceDiff) error {
	groupLimit := GroupLimit{
		Units:        limits.Units,
		Allocation:   limits.AllocationCount,
		Minimum:      limits.MinimumCount,
		Maximum:      limits.MaximumCount,
		StepSize:     limits.StepSizeCount,
		IsAdjustable: limits.IsAdjustable,
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
	CanScaleDown bool
}

func checkMbValue(name string, limits MbLimit, divider int, diff *schema.ResourceDiff) error {
	groupLimit := GroupLimit{
		Units:        limits.Units,
		Allocation:   limits.AllocationMb,
		Minimum:      limits.MinimumMb,
		Maximum:      limits.MaximumMb,
		StepSize:     limits.StepSizeMb,
		IsAdjustable: limits.IsAdjustable,
		CanScaleDown: limits.CanScaleDown,
	}
	return checkGroupValue(name, groupLimit, divider, diff)
}

func resourceIBMDatabaseInstanceDiff(diff *schema.ResourceDiff, meta interface{}) error {

	err := resourceTagsCustomizeDiff(diff)
	if err != nil {
		return err
	}

	service := diff.Get("service").(string)
	if service == "databases-for-postgresql" || service == "databases-for-elasticsearch" {
		planPhase := diff.Get("plan_validation").(bool)

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
				_, newSetting := diff.GetChange("node_count")
				min := groupDefaults.Cpu.MinimumCount / divider
				if newSetting != min {
					return fmt.Errorf("node_cpu_allocation_count must be set when node_count is greater then the minimum %d", min)
				}
			}
		}
	} else if diff.HasChange("node_count") || diff.HasChange("node_memory_allocation_mb") || diff.HasChange("node_disk_allocation_mb") || diff.HasChange("node_cpu_allocation_count") {
		return fmt.Errorf("node_count, node_memory_allocation_mb, node_disk_allocation_mb, node_cpu_allocation_count only supported for postgresql and elasticsearch")
	}

	return nil
}

// Replace with func wrapper for resourceIBMResourceInstanceCreate specifying serviceName := "database......."
func resourceIBMDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	rsInst := rc.CreateResourceInstanceOptions{
		Name: &name,
	}

	rsCatClient, err := meta.(ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return fmt.Errorf("Error retrieving database service offering: %s", err)
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return fmt.Errorf("Error retrieving plan: %s", err)
	}
	rsInst.ResourcePlanID = &servicePlan

	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return fmt.Errorf("Error retrieving deployment for plan %s : %s", plan, err)
	}
	if len(deployments) == 0 {
		return fmt.Errorf("No deployment found for service plan : %s", plan)
	}
	deployments, supportedLocations := filterDatabaseDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return fmt.Errorf("No deployment found for service plan %s at location %s.\nValid location(s) are: %q.", plan, location, locationList)
	}
	catalogCRN := deployments[0].CatalogCRN
	rsInst.Target = &catalogCRN

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rgID := rsGrpID.(string)
		rsInst.ResourceGroup = &rgID
	} else {
		defaultRg, err := defaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInst.ResourceGroup = &defaultRg
	}

	initialNodeCount, err := getInitialNodeCount(d, meta)
	if err != nil {
		return err
	}

	params := Params{}
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
	if pitrTime, ok := d.GetOk("point_in_time_recovery_time"); ok {
		params.PITRTimeStamp = pitrTime.(string)
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
		return fmt.Errorf("Error creating database instance: %s %s", err, response)
	}

	// Moved d.SetId(instance.ID) to after waiting for resource to finish creation. Otherwise Terraform initates depedent tasks too early.
	// Original flow had SetId here as its required as input to waitForDatabaseInstanceCreate

	_, err = waitForDatabaseInstanceCreate(d, meta, *instance.ID)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for create database instance (%s) to complete: %s", d.Id(), err)
	}

	d.SetId(*instance.ID)

	if node_count, ok := d.GetOk("node_count"); ok {
		if initialNodeCount != node_count {
			icdClient, err := meta.(ClientSession).ICDAPI()
			if err != nil {
				return fmt.Errorf("Error getting database client settings: %s", err)
			}

			err = horizontalScale(d, meta, icdClient)
			if err != nil {
				return err
			}
		}
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of ibm database (%s) tags: %s", d.Id(), err)
		}
	}

	icdId := EscapeUrlParm(*instance.ID)
	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return fmt.Errorf("Error getting database client settings: %s", err)
	}

	if pw, ok := d.GetOk("adminpassword"); ok {
		adminPassword := pw.(string)
		cdb, err := icdClient.Cdbs().GetCdb(icdId)
		if err != nil {
			if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
				return fmt.Errorf("The database instance was not found in the region set for the Provider, or the default of us-south. Specify the correct region in the provider definition, or create a provider alias for the correct region. %v", err)
			}
			return fmt.Errorf("Error getting database config while updating adminpassword for: %s with error %s\n", icdId, err)
		}

		userParams := icdv4.UserReq{
			User: icdv4.User{
				Password: adminPassword,
			},
		}
		task, err := icdClient.Users().UpdateUser(icdId, cdb.AdminUser, userParams)
		if err != nil {
			return fmt.Errorf("Error updating database admin password: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for update of database (%s) admin password task to complete: %s", icdId, err)
		}
	}

	if wl, ok := d.GetOk("whitelist"); ok {
		whitelist := expandWhitelist(wl.(*schema.Set))
		for _, wlEntry := range whitelist {
			whitelistReq := icdv4.WhitelistReq{
				WhitelistEntry: icdv4.WhitelistEntry{
					Address:     wlEntry.Address,
					Description: wlEntry.Description,
				},
			}
			task, err := icdClient.Whitelists().CreateWhitelist(icdId, whitelistReq)
			if err != nil {
				return fmt.Errorf("Error updating database whitelist entry: %s", err)
			}
			_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return fmt.Errorf(
					"Error waiting for update of database (%s) whitelist task to complete: %s", icdId, err)
			}
		}
	}
	if cpuRecord, ok := d.GetOk("auto_scaling.0.cpu"); ok {
		params := icdv4.AutoscalingSetGroup{}
		cpuBody, err := expandICDAutoScalingGroup(d, cpuRecord, "cpu")
		if err != nil {
			return fmt.Errorf("Error in getting cpuBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.CPU = &cpuBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database cpu auto_scaling group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) cpu auto_scaling group update task to complete: %s", icdId, err)
		}

	}
	if diskRecord, ok := d.GetOk("auto_scaling.0.disk"); ok {
		params := icdv4.AutoscalingSetGroup{}
		diskBody, err := expandICDAutoScalingGroup(d, diskRecord, "disk")
		if err != nil {
			return fmt.Errorf("Error in getting diskBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.Disk = &diskBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database disk auto_scaling group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) disk auto_scaling group update task to complete: %s", icdId, err)
		}

	}
	if memoryRecord, ok := d.GetOk("auto_scaling.0.memory"); ok {
		params := icdv4.AutoscalingSetGroup{}
		memoryBody, err := expandICDAutoScalingGroup(d, memoryRecord, "memory")
		if err != nil {
			return fmt.Errorf("Error in getting memoryBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.Memory = &memoryBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database memory auto_scaling group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) memory auto_scaling group update task to complete: %s", icdId, err)
		}

	}

	if userlist, ok := d.GetOk("users"); ok {
		users := expandUsers(userlist.(*schema.Set))
		for _, user := range users {
			userReq := icdv4.UserReq{
				User: icdv4.User{
					UserName: user.UserName,
					Password: user.Password,
				},
			}
			task, err := icdClient.Users().CreateUser(icdId, userReq)
			if err != nil {
				return fmt.Errorf("Error updating database user (%s) entry: %s", user.UserName, err)
			}
			_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return fmt.Errorf(
					"Error waiting for update of database (%s) user (%s) create task to complete: %s", icdId, user.UserName, err)
			}
		}
	}

	return resourceIBMDatabaseInstanceRead(d, meta)
}

func resourceIBMDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
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
		return fmt.Errorf("Error retrieving resource instance: %s %s", err, response)
	}
	if strings.Contains(*instance.State, "removed") {
		log.Printf("[WARN] Removing instance from TF state because it's now in removed state")
		d.SetId("")
		return nil
	}

	tags, err := GetTagsUsingCRN(meta, *instance.CRN)
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

	d.Set(ResourceName, *instance.Name)
	d.Set(ResourceCRN, *instance.CRN)
	d.Set(ResourceStatus, *instance.State)
	d.Set(ResourceGroupName, *instance.ResourceGroupCRN)

	rcontroller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, rcontroller+"/services/"+url.QueryEscape(*instance.CRN))

	rsCatClient, err := meta.(ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.GetServiceName(*instance.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving service offering: %s", err)
	}

	d.Set("service", serviceOff)

	servicePlan, err := rsCatRepo.GetServicePlanName(*instance.ResourcePlanID)
	if err != nil {
		return fmt.Errorf("Error retrieving plan: %s", err)
	}
	d.Set("plan", servicePlan)

	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return fmt.Errorf("Error getting database client settings: %s", err)
	}

	icdId := EscapeUrlParm(instanceID)
	cdb, err := icdClient.Cdbs().GetCdb(icdId)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
			return fmt.Errorf("The database instance was not found in the region set for the Provider. Specify the correct region in the provider definition. %v", err)
		}
		return fmt.Errorf("Error getting database config for: %s with error %s\n", icdId, err)
	}
	d.Set("adminuser", cdb.AdminUser)
	d.Set("version", cdb.Version)

	groupList, err := icdClient.Groups().GetGroups(icdId)
	if err != nil {
		return fmt.Errorf("Error getting database groups: %s", err)
	}
	d.Set("groups", flattenIcdGroups(groupList))
	d.Set("node_count", groupList.Groups[0].Members.AllocationCount)

	d.Set("members_memory_allocation_mb", groupList.Groups[0].Memory.AllocationMb)
	d.Set("node_memory_allocation_mb", groupList.Groups[0].Memory.AllocationMb/groupList.Groups[0].Members.AllocationCount)

	d.Set("members_disk_allocation_mb", groupList.Groups[0].Disk.AllocationMb)
	d.Set("node_disk_allocation_mb", groupList.Groups[0].Disk.AllocationMb/groupList.Groups[0].Members.AllocationCount)

	d.Set("members_cpu_allocation_count", groupList.Groups[0].Cpu.AllocationCount)
	d.Set("node_cpu_allocation_count", groupList.Groups[0].Cpu.AllocationCount/groupList.Groups[0].Members.AllocationCount)

	autoSclaingGroup, err := icdClient.AutoScaling().GetAutoScaling(icdId, "member")
	if err != nil {
		return fmt.Errorf("Error getting database groups: %s", err)
	}
	d.Set("auto_scaling", flattenICDAutoScalingGroup(autoSclaingGroup))

	whitelist, err := icdClient.Whitelists().GetWhitelist(icdId)
	if err != nil {
		return fmt.Errorf("Error getting database whitelist: %s", err)
	}
	d.Set("whitelist", flattenWhitelist(whitelist))

	var connectionStrings []CsEntry
	//ICD does not implement a GetUsers API. Users populated from tf configuration.
	tfusers := d.Get("users").(*schema.Set)
	users := expandUsers(tfusers)
	user := icdv4.User{
		UserName: cdb.AdminUser,
	}
	users = append(users, user)
	for _, user := range users {
		userName := user.UserName
		csEntry, err := getConnectionString(d, userName, connectionEndpoint, meta)
		if err != nil {
			return fmt.Errorf("Error getting user connection string for user (%s): %s", userName, err)
		}
		connectionStrings = append(connectionStrings, csEntry)
	}
	d.Set("connectionstrings", flattenConnectionStrings(connectionStrings))

	return nil
}

func resourceIBMDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
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
			return fmt.Errorf("Error updating resource instance: %s %s", err, response)
		}

		_, err = waitForDatabaseInstanceUpdate(d, meta)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for update of resource instance (%s) to complete: %s", d.Id(), err)
		}

	}

	if d.HasChange("tags") {

		oldList, newList := d.GetChange("tags")
		err = UpdateTagsUsingCRN(oldList, newList, meta, instanceID)
		if err != nil {
			log.Printf(
				"Error on update of Database (%s) tags: %s", d.Id(), err)
		}
	}

	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return fmt.Errorf("Error getting database client settings: %s", err)
	}
	icdId := EscapeUrlParm(instanceID)

	if d.HasChange("node_count") {
		err = horizontalScale(d, meta, icdClient)
		if err != nil {
			return err
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
			return fmt.Errorf("Error updating database scaling group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) scaling group update task to complete: %s", icdId, err)
		}
	}

	if d.HasChange("auto_scaling.0.cpu") {
		cpuRecord := d.Get("auto_scaling.0.cpu")
		params := icdv4.AutoscalingSetGroup{}
		cpuBody, err := expandICDAutoScalingGroup(d, cpuRecord, "cpu")
		if err != nil {
			return fmt.Errorf("Error in getting cpuBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.CPU = &cpuBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database cpu auto_scaling group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) cpu auto_scaling group update task to complete: %s", icdId, err)
		}

	}
	if d.HasChange("auto_scaling.0.disk") {
		diskRecord := d.Get("auto_scaling.0.disk")
		params := icdv4.AutoscalingSetGroup{}
		diskBody, err := expandICDAutoScalingGroup(d, diskRecord, "disk")
		if err != nil {
			return fmt.Errorf("Error in getting diskBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.Disk = &diskBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database disk auto_scaling  group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) disk auto_scaling group update task to complete: %s", icdId, err)
		}

	}
	if d.HasChange("auto_scaling.0.memory") {
		memoryRecord := d.Get("auto_scaling.0.memory")
		params := icdv4.AutoscalingSetGroup{}
		memoryBody, err := expandICDAutoScalingGroup(d, memoryRecord, "memory")
		if err != nil {
			return fmt.Errorf("Error in getting memoryBody from expandICDAutoScalingGroup %s", err)
		}
		params.Autoscaling.Memory = &memoryBody
		task, err := icdClient.AutoScaling().SetAutoScaling(icdId, "member", params)
		if err != nil {
			return fmt.Errorf("Error updating database memory auto_scaling  group: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) memory auto_scaling group update task to complete: %s", icdId, err)
		}

	}

	if d.HasChange("adminpassword") {
		adminUser := d.Get("adminuser").(string)
		password := d.Get("adminpassword").(string)
		userParams := icdv4.UserReq{
			User: icdv4.User{
				Password: password,
			},
		}
		task, err := icdClient.Users().UpdateUser(icdId, adminUser, userParams)
		if err != nil {
			return fmt.Errorf("Error updating database admin password: %s", err)
		}
		_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for database (%s) admin password update task to complete: %s", icdId, err)
		}
	}

	if d.HasChange("whitelist") {
		oldList, newList := d.GetChange("whitelist")
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

		if len(add) > 0 {
			for _, entry := range add {
				newEntry := entry.(map[string]interface{})
				wlEntry := icdv4.WhitelistEntry{
					Address:     newEntry["address"].(string),
					Description: newEntry["description"].(string),
				}
				whitelistReq := icdv4.WhitelistReq{
					WhitelistEntry: wlEntry,
				}
				task, err := icdClient.Whitelists().CreateWhitelist(icdId, whitelistReq)
				if err != nil {
					return fmt.Errorf("Error updating database whitelist entry %v : %s", wlEntry.Address, err)
				}
				_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return fmt.Errorf(
						"Error waiting for database (%s) whitelist create task to complete for entry %s : %s", icdId, wlEntry.Address, err)
				}

			}

		}

		if len(remove) > 0 {
			for _, entry := range remove {
				newEntry := entry.(map[string]interface{})
				wlEntry := icdv4.WhitelistEntry{
					Address:     newEntry["address"].(string),
					Description: newEntry["description"].(string),
				}
				ipAddress := wlEntry.Address
				task, err := icdClient.Whitelists().DeleteWhitelist(icdId, ipAddress)
				if err != nil {
					return fmt.Errorf("Error deleting database whitelist entry: %s", err)
				}
				_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return fmt.Errorf(
						"Error waiting for database (%s) whitelist delete task to complete for ipAddress %s : %s", icdId, ipAddress, err)
				}

			}
		}
	}

	if d.HasChange("users") {
		oldList, newList := d.GetChange("users")
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

		if len(add) > 0 {
			for _, entry := range add {
				newEntry := entry.(map[string]interface{})
				userEntry := icdv4.User{
					UserName: newEntry["name"].(string),
					Password: newEntry["password"].(string),
				}
				userReq := icdv4.UserReq{
					User: userEntry,
				}
				task, err := icdClient.Users().CreateUser(icdId, userReq)
				if err != nil {
					// ICD does not report if error was due to user already being defined. Check if can
					// successfully update password by itself.
					userParams := icdv4.UserReq{
						User: icdv4.User{
							Password: newEntry["password"].(string),
						},
					}
					task, err := icdClient.Users().UpdateUser(icdId, newEntry["name"].(string), userParams)
					if err != nil {
						return fmt.Errorf("Error updating database user (%s) password: %s", newEntry["name"].(string), err)
					}
					_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return fmt.Errorf(
							"Error waiting for database (%s) user (%s) password update task to complete: %s", icdId, newEntry["name"].(string), err)
					}
				} else {
					_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
					if err != nil {
						return fmt.Errorf(
							"Error waiting for database (%s) user (%s) create task to complete: %s", icdId, newEntry["name"].(string), err)
					}
				}
			}

		}

		if len(remove) > 0 {
			for _, entry := range remove {
				newEntry := entry.(map[string]interface{})
				userEntry := icdv4.User{
					UserName: newEntry["name"].(string),
					Password: newEntry["password"].(string),
				}
				user := userEntry.UserName
				task, err := icdClient.Users().DeleteUser(icdId, user)
				if err != nil {
					return fmt.Errorf("Error deleting database user (%s) entry: %s", user, err)
				}
				_, err = waitForDatabaseTaskComplete(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return fmt.Errorf(
						"Error waiting for database (%s) user (%s) delete task to complete: %s", icdId, user, err)
				}
			}
		}
	}

	return resourceIBMDatabaseInstanceRead(d, meta)
}

func horizontalScale(d *schema.ResourceData, meta interface{}, icdClient icdv4.ICDServiceAPI) error {
	params := icdv4.GroupReq{}

	icdId := EscapeUrlParm(d.Id())

	members := d.Get("node_count").(int)
	membersReq := icdv4.MembersReq{AllocationCount: members}
	params.GroupBdy.Members = &membersReq

	//task, err := icdClient.Groups().UpdateGroup(icdId, "member", params)
	_, err := icdClient.Groups().UpdateGroup(icdId, "member", params)

	if err != nil {
		return fmt.Errorf("Error updating database scaling group: %s", err)
	}

	//_, err = waitForDatabaseTaskCompleteDuration(task.Id, d, meta, d.Timeout(schema.TimeoutUpdate))
	//if err != nil {
	//	return fmt.Errorf(
	//		"Error waiting for database (%s) scaling group update task to complete: %s", icdId, err)
	//}

	_, err = waitForDatabaseInstanceUpdate(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for database (%s) horizontal scale to complete: %s", d.Id(), err)
	}

	return nil
}

func getConnectionString(d *schema.ResourceData, userName, connectionEndpoint string, meta interface{}) (CsEntry, error) {
	csEntry := CsEntry{}
	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return csEntry, fmt.Errorf("Error getting database client settings: %s", err)
	}

	icdId := d.Id()
	connection, err := icdClient.Connections().GetConnection(icdId, userName, connectionEndpoint)
	if err != nil {
		return csEntry, fmt.Errorf("Error getting database user connection string via ICD API: %s", err)
	}

	service := d.Get("service")
	dbConnection := icdv4.Uri{}
	switch service {
	case "databases-for-postgresql":
		dbConnection = connection.Postgres
	case "databases-for-redis":
		dbConnection = connection.Rediss
	case "databases-for-mongodb":
		dbConnection = connection.Mongo
	// case "databases-for-mysql":
	// 	dbConnection = connection.Mysql
	case "databases-for-elasticsearch":
		dbConnection = connection.Https
	case "databases-for-etcd":
		dbConnection = connection.Grpc
	case "messages-for-rabbitmq":
		dbConnection = connection.Amqps
	default:
		return csEntry, fmt.Errorf("Unrecognised database type during connection string lookup: %s", service)
	}

	csEntry = CsEntry{
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
	return csEntry, nil
}

func resourceIBMDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
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
			return fmt.Errorf("Error deleting resource instance: %s %s ", err, response)
		}
	}

	_, err = waitForDatabaseInstanceDelete(d, meta)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for resource instance (%s) to be deleted: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
func resourceIBMDatabaseInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
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
		return false, fmt.Errorf("Error communicating with the API: %s %s", err, response)
	}
	if strings.Contains(*instance.State, "removed") {
		log.Printf("[WARN] Removing instance from state because it's in removed state")
		d.SetId("")
		return false, nil
	}

	return *instance.ID == instanceID, nil
}

func waitForDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}, instanceID string) (interface{}, error) {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
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
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The resource instance %s does not exist anymore: %v %s", d.Id(), err, response)
				}
				return nil, "", err
			}
			if *instance.State == databaseInstanceFailStatus {
				return *instance, *instance.State, fmt.Errorf("The resource instance %s failed: %v %s", d.Id(), err, response)
			}
			return *instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
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
					return nil, "", fmt.Errorf("The resource instance %s does not exist anymore: %v %s", d.Id(), err, response)
				}
				return nil, "", err
			}
			if *instance.State == databaseInstanceFailStatus {
				return *instance, *instance.State, fmt.Errorf("The resource instance %s failed: %v %s", d.Id(), err, response)
			}
			return *instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func waitForDatabaseTaskComplete(taskId string, d *schema.ResourceData, meta interface{}, t time.Duration) (bool, error) {
	icdClient, err := meta.(ClientSession).ICDAPI()
	if err != nil {
		return false, fmt.Errorf("Error getting database client settings: %s", err)
	}
	delayDuration := 5 * time.Second

	timeout := time.After(t)
	delay := time.Tick(delayDuration)
	innerTask := icdv4.Task{}

	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("[Error] Time out waiting for database task to complete")
		case <-delay:
			innerTask, err = icdClient.Tasks().GetTask(EscapeUrlParm(taskId))
			if err != nil {
				return false, fmt.Errorf("The ICD Get task on database update errored: %v", err)
			}
			if innerTask.Status == "failed" {
				return false, fmt.Errorf("[Error] Database task failed")
			}
			// Completed status could be returned as "" due to interaction between bluemix-go and icd task response
			// Otherwise Running an queued
			if innerTask.Status == "completed" || innerTask.Status == "" {
				return true, nil
			}
		}
	}
}

func waitForDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(ClientSession).ResourceControllerV2API()
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
				return nil, "", err
			}
			if *instance.State == databaseInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("The resource instance %s failed to delete: %v %s", d.Id(), err, response)
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

func expandICDAutoScalingGroup(d *schema.ResourceData, asRecord interface{}, asType string) (asgBody icdv4.ASGBody, err error) {

	asgRecord := asRecord.([]interface{})[0].(map[string]interface{})
	asgCapacity := icdv4.CapacityBody{}
	if _, ok := asgRecord["capacity_enabled"]; ok {
		asgCapacity.Enabled = asgRecord["capacity_enabled"].(bool)
		asgBody.Scalers.Capacity = &asgCapacity
	}
	if _, ok := asgRecord["free_space_less_than_percent"]; ok {
		asgCapacity.FreeSpaceLessThanPercent = asgRecord["free_space_less_than_percent"].(int)
		asgBody.Scalers.Capacity = &asgCapacity
	}

	// IO Payload
	asgIO := icdv4.IOBody{}
	if _, ok := asgRecord["io_enabled"]; ok {
		asgIO.Enabled = asgRecord["io_enabled"].(bool)
		asgBody.Scalers.IO = &asgIO
	}
	if _, ok := asgRecord["io_over_period"]; ok {
		asgIO.OverPeriod = asgRecord["io_over_period"].(string)
		asgBody.Scalers.IO = &asgIO
	}
	if _, ok := asgRecord["io_above_percent"]; ok {
		asgIO.AbovePercent = asgRecord["io_above_percent"].(int)
		asgBody.Scalers.IO = &asgIO
	}

	// Rate Payload
	asgRate := icdv4.RateBody{}
	if _, ok := asgRecord["rate_increase_percent"]; ok {
		asgRate.IncreasePercent = asgRecord["rate_increase_percent"].(int)
		asgBody.Rate = asgRate
	}
	if _, ok := asgRecord["rate_period_seconds"]; ok {
		asgRate.PeriodSeconds = asgRecord["rate_period_seconds"].(int)
		asgBody.Rate = asgRate
	}
	if _, ok := asgRecord["rate_limit_mb_per_member"]; ok {
		asgRate.LimitMBPerMember = asgRecord["rate_limit_mb_per_member"].(int)
		asgBody.Rate = asgRate
	}
	if _, ok := asgRecord["rate_limit_count_per_member"]; ok {
		asgRate.LimitCountPerMember = asgRecord["rate_limit_count_per_member"].(int)
		asgBody.Rate = asgRate
	}
	if _, ok := asgRecord["rate_units"]; ok {
		asgRate.Units = asgRecord["rate_units"].(string)
		asgBody.Rate = asgRate
	}

	return asgBody, nil
}

func flattenICDAutoScalingGroup(autoScalingGroup icdv4.AutoscalingGetGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	memorys := make([]map[string]interface{}, 0)
	memory := make(map[string]interface{})

	if autoScalingGroup.Autoscaling.Memory.Scalers.IO != nil {
		memoryIO := *autoScalingGroup.Autoscaling.Memory.Scalers.IO
		memory["io_enabled"] = memoryIO.Enabled
		memory["io_over_period"] = memoryIO.OverPeriod
		memory["io_above_percent"] = memoryIO.AbovePercent
	}
	if &autoScalingGroup.Autoscaling.Memory.Rate != nil {
		ip, _ := autoScalingGroup.Autoscaling.Memory.Rate.IncreasePercent.Float64()
		memory["rate_increase_percent"] = int(ip)
		memory["rate_period_seconds"] = autoScalingGroup.Autoscaling.Memory.Rate.PeriodSeconds
		lmp, _ := autoScalingGroup.Autoscaling.Memory.Rate.LimitMBPerMember.Float64()
		memory["rate_limit_mb_per_member"] = int(lmp)
		memory["rate_units"] = autoScalingGroup.Autoscaling.Memory.Rate.Units
	}
	memorys = append(memorys, memory)

	cpus := make([]map[string]interface{}, 0)
	cpu := make(map[string]interface{})

	if &autoScalingGroup.Autoscaling.CPU.Rate != nil {

		ip, _ := autoScalingGroup.Autoscaling.CPU.Rate.IncreasePercent.Float64()
		cpu["rate_increase_percent"] = int(ip)
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
	if autoScalingGroup.Autoscaling.Disk.Scalers.IO != nil {
		diskIO := *autoScalingGroup.Autoscaling.Disk.Scalers.IO
		disk["io_enabled"] = diskIO.Enabled
		disk["io_over_period"] = diskIO.OverPeriod
		disk["io_above_percent"] = diskIO.AbovePercent
	}
	if &autoScalingGroup.Autoscaling.Disk.Rate != nil {

		ip, _ := autoScalingGroup.Autoscaling.Disk.Rate.IncreasePercent.Float64()
		disk["rate_increase_percent"] = int(ip)
		disk["rate_period_seconds"] = autoScalingGroup.Autoscaling.Disk.Rate.PeriodSeconds
		lpm, _ := autoScalingGroup.Autoscaling.Disk.Rate.LimitMBPerMember.Float64()
		disk["rate_limit_mb_per_member"] = int(lpm)
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
