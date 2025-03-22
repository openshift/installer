// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package db2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"
	"github.com/IBM/go-sdk-core/v5/core"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RsInstanceSuccessStatus       = "active"
	RsInstanceProgressStatus      = "in progress"
	RsInstanceProvisioningStatus  = "provisioning"
	RsInstanceInactiveStatus      = "inactive"
	RsInstanceFailStatus          = "failed"
	RsInstanceRemovedStatus       = "removed"
	RsInstanceReclamation         = "pending_reclamation"
	RsInstanceUpdateSuccessStatus = "succeeded"
	PerformanceSubscription       = "PerformanceSubscription"
)

func ResourceIBMDb2Instance() *schema.Resource {
	riSchema := resourcecontroller.ResourceIBMResourceInstance().Schema

	riSchema["high_availability"] = &schema.Schema{
		Description: "If you require high availability, please choose this option",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["instance_type"] = &schema.Schema{
		Description: "Available machine type flavours (default selection will assume smallest configuration)",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["backup_location"] = &schema.Schema{
		Description: "Cross Regional backups can be stored across multiple regions in a zone. Regional backups are stored in only specific region.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["disk_encryption_instance_crn"] = &schema.Schema{
		Description: "Cross Regional disk encryption crn",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["disk_encryption_key_crn"] = &schema.Schema{
		Description: "Cross Regional disk encryption crn",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["oracle_compatibility"] = &schema.Schema{
		Description: "Indicates whether is has compatibility for oracle or not",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["subscription_id"] = &schema.Schema{
		Description: "For PerformanceSubscription plans a Subscription ID is required. It is not required for Performance plans.",
		Optional:    true,
		Type:        schema.TypeString,
	}

	riSchema["autoscale_config"] = &schema.Schema{
		Description: "The db2 auto scaling config",
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_scaling_threshold": {
					Description: "The auto_scaling_threshold of the instance",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"auto_scaling_over_time_period": {
					Description: "The auto_scaling_over_time_period of the instance",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"auto_scaling_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates if automatic scaling is enabled or not.",
				},
				"auto_scaling_allow_plan_limit": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates the maximum number of scaling actions that are allowed within a specified time period.",
				},
				"auto_scaling_pause_limit": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Specifies the duration to pause auto-scaling actions after a scaling event has occurred.",
				},
			},
		},
	}

	riSchema["custom_setting_config"] = &schema.Schema{
		Description: "Db and Dm configurations",
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Tunable parameters related to the Db2 database instance.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"act_sortmem_limit": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"alt_collate": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"appgroup_mem_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"applheapsz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"appl_memory": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"app_ctl_heap_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"archretrydelay": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"authn_cache_duration": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"autorestart": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_cg_stats": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_maint": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_reorg": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_reval": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_runstats": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_sampling": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_stats_views": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_stmt_stats": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"auto_tbl_maint": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"avg_appls": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"catalogcache_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"chngpgs_thresh": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"cur_commit": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"database_memory": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dbheap": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db_collname": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db_mem_thresh": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"ddl_compression_def": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"ddl_constraint_def": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"decflt_rounding": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dec_arithmetic": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dec_to_char_fmt": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_degree": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_extent_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_loadrec_ses": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mttb_types": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_prefetch_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_queryopt": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_refresh_age": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_schemas_dcc": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_sqlmathwarn": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_table_org": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dlchktime": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"enable_xmlchar": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"extended_row_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"groupheap_ratio": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"indexrec": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"large_aggregation": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"locklist": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"locktimeout": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"logindexbuild": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"log_appl_info": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"log_ddl_stmts": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"log_disk_cap": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"maxappls": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"maxfilop": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"maxlocks": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"min_dec_div_3": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_act_metrics": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_deadlock": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_lck_msg_lvl": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_locktimeout": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_lockwait": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_lw_thresh": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_obj_metrics": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_pkglist_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_req_metrics": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_rtn_data": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_rtn_execlist": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_uow_data": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_uow_execlist": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_uow_pkglist": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"nchar_mapping": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_freqvalues": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_iocleaners": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_ioservers": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_log_span": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_quantiles": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"opt_buffpage": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"opt_direct_wrkld": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"opt_locklist": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"opt_maxlocks": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"opt_sortheap": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"page_age_trgt_gcr": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"page_age_trgt_mcr": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"pckcachesz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"pl_stack_trace": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"self_tuning_mem": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"seqdetect": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"sheapthres_shr": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"softmax": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"sortheap": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"sql_ccflags": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"stat_heap_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"stmtheap": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"stmt_conc": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"string_units": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"systime_period_adj": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"trackmod": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"util_heap_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_admission_ctrl": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_agent_load_trgt": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_cpu_limit": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_cpu_shares": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_cpu_share_mode": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"dbm": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Tunable parameters related to the Db2 instance manager (dbm).",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"comm_bandwidth": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"cpuspeed": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_bufpool": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_lock": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_sort": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_stmt": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_table": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_timestamp": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"dft_mon_uow": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"diaglevel": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"federated_async": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"indexrec": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"intra_parallel": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"keepfenced": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"max_connretries": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"max_querydegree": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"mon_heap_sz": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"multipartsizemb": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"notifylevel": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_initagents": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_initfenced": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"num_poolagents": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"resync_interval": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"rqrioblk": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"start_stop_time": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"util_impact_lim": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_dispatcher": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_disp_concur": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_disp_cpu_shares": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"wlm_disp_min_util": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"registry": &schema.Schema{
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Tunable parameters related to the Db2 registry.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"db2_bidi": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_compopt": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_lock_to_rb": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_stmm": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_alternate_authz_behaviour": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_antijoin": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_ats_enable": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_deferred_prepare_semantics": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_evaluncommitted": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_extended_optimization": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_index_pctfree_default": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_inlist_to_nljn": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_minimize_listprefetch": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_object_table_entries": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_optprofile": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_optstats_log": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_opt_max_temp_size": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_parallel_io": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_reduced_optimization": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_selectivity": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_skipdeleted": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_skipinserted": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_sync_release_lock_attributes": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_truncate_reusestorage": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_use_alternate_page_cleaning": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_view_reopt_values": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_wlm_settings": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"db2_workload": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}

	return &schema.Resource{
		Create:   resourceIBMDb2InstanceCreate,
		Read:     resourcecontroller.ResourceIBMResourceInstanceRead,
		Update:   resourcecontroller.ResourceIBMResourceInstanceUpdate,
		Delete:   resourcecontroller.ResourceIBMResourceInstanceDelete,
		Exists:   resourcecontroller.ResourceIBMResourceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: riSchema,
	}
}

func resourceIBMDb2InstanceCreate(d *schema.ResourceData, meta interface{}) error {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
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

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	if metadata, ok := serviceOff[0].Metadata.(*models.ServiceResourceMetadata); ok {
		if !metadata.Service.RCProvisionable {
			return fmt.Errorf("%s cannot be provisioned by resource controller", serviceName)
		}
	} else {
		return fmt.Errorf("[ERROR] Cannot create instance of resource %s\nUse 'ibm_service_instance' if the resource is a Cloud Foundry service", serviceName)
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	rsInst.ResourcePlanID = &servicePlan

	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err)
	}
	if len(deployments) == 0 {
		return fmt.Errorf("[ERROR] No deployment found for service plan : %s", plan)
	}
	deployments, supportedLocations := resourcecontroller.FilterDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\nValid location(s) are: %q.\nUse 'ibm_service_instance' if the service is a Cloud Foundry service", plan, location, locationList)
	}

	rsInst.Target = &deployments[0].CatalogCRN

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rg := rsGrpID.(string)
		rsInst.ResourceGroup = &rg
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInst.ResourceGroup = &defaultRg
	}

	params := map[string]interface{}{}

	if serviceEndpoints, ok := d.GetOk("service_endpoints"); ok {
		params["service-endpoints"] = serviceEndpoints.(string)
	}
	if highAvailability, ok := d.GetOk("high_availability"); ok {
		params["high_availability"] = highAvailability.(string)
	}
	if instanceType, ok := d.GetOk("instance_type"); ok {
		params["instance_type"] = instanceType.(string)
	}
	if backupLocation, ok := d.GetOk("backup_location"); ok {
		params["backup-locations"] = backupLocation.(string)
	}

	if diskEncryptionInstanceCrn, ok := d.GetOk("disk_encryption_instance_crn"); ok {
		params["disk_encryption_instance_crn"] = diskEncryptionInstanceCrn.(string)
	}

	if diskEncryptionKeyCrn, ok := d.GetOk("disk_encryption_key_crn"); ok {
		params["disk_encryption_key_crn"] = diskEncryptionKeyCrn.(string)
	}

	if oracleCompatibility, ok := d.GetOk("oracle_compatibility"); ok {
		params["oracle_compatibility"] = oracleCompatibility.(string)
	}

	if plan == PerformanceSubscription {
		if subscriptionId, ok := d.GetOk("subscription_id"); ok {
			params["subscription_id"] = subscriptionId.(string)
		} else {
			return fmt.Errorf("[ERROR] Missing required field 'subscription_id' while creating an instance for plan: %s", plan)
		}
	}

	if parameters, ok := d.GetOk("parameters"); ok {
		temp := parameters.(map[string]interface{})
		for k, v := range temp {
			if v == "true" || v == "false" {
				b, _ := strconv.ParseBool(v.(string))
				params[k] = b
			} else if strings.HasPrefix(v.(string), "[") && strings.HasSuffix(v.(string), "]") {
				//transform v.(string) to be []string
				arrayString := v.(string)
				result := []string{}
				trimLeft := strings.TrimLeft(arrayString, "[")
				trimRight := strings.TrimRight(trimLeft, "]")
				if len(trimRight) == 0 {
					params[k] = result
				} else {
					array := strings.Split(trimRight, ",")
					for _, a := range array {
						result = append(result, strings.Trim(a, "\""))
					}
					params[k] = result
				}
			} else {
				params[k] = v
			}
		}

	}

	if s, ok := d.GetOk("parameters_json"); ok {
		json.Unmarshal([]byte(s.(string)), &params)
	}

	rsInst.Parameters = params

	//Start to create resource instance
	instance, resp, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil {
		log.Printf(
			"Error when creating resource instance: %s, Instance info  NAME->%s, LOCATION->%s, GROUP_ID->%s, PLAN_ID->%s",
			err, *rsInst.Name, *rsInst.Target, *rsInst.ResourceGroup, *rsInst.ResourcePlanID)
		return fmt.Errorf("[ERROR] Error when creating resource instance: %s with resp code: %s", err, resp)
	}

	d.SetId(*instance.ID)

	_, err = waitForResourceInstanceCreate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for create resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	db2SaasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		return err
	}

	encodedCRN := url.QueryEscape(*instance.CRN)

	if autoscaleConfigRaw, ok := d.GetOk("autoscale_config"); ok {
		if autoscaleConfigRaw == nil || reflect.ValueOf(autoscaleConfigRaw).IsNil() {
			fmt.Println("No autoscaling config is provided, Skipping.")
		} else {
			autoscalingConfig := autoscaleConfigRaw.([]interface{})[0].(map[string]interface{})

			var (
				autoScalingThreshold      int
				autoScalingOverTimePeriod int
				autoScalingPauseLimit     int
			)

			if autoscalingConfig["auto_scaling_threshold"] != nil {
				autoScalingThreshold, err = strconv.Atoi(autoscalingConfig["auto_scaling_threshold"].(string))
				if err != nil {
					fmt.Println("Error in converting auto_scaling_threshold to string", err)
					return err
				}
			}

			if autoscalingConfig["auto_scaling_over_time_period"] != nil {
				autoScalingOverTimePeriod, err = strconv.Atoi(autoscalingConfig["auto_scaling_over_time_period"].(string))
				if err != nil {
					fmt.Println("Error in converting auto_scaling_over_time_period to string", err)
					return err
				}
			}

			if autoscalingConfig["auto_scaling_pause_limit"] != nil {
				switch v := autoscalingConfig["auto_scaling_pause_limit"].(type) {
				case string:
					autoScalingPauseLimit, err = strconv.Atoi(v)
					if err != nil {
						fmt.Println("Error in converting auto_scaling_pause_limit to int:", err)
						return err
					}
				case int:
					autoScalingPauseLimit = v
				}
			}

			input := &db2saasv1.PutDb2SaasAutoscaleOptions{
				XDbProfile:                core.StringPtr(encodedCRN),
				AutoScalingEnabled:        core.StringPtr("YES"),
				AutoScalingAllowPlanLimit: core.StringPtr("YES"),
				AutoScalingThreshold:      core.Int64Ptr(int64(autoScalingThreshold)),
				AutoScalingOverTimePeriod: core.Float64Ptr(float64(autoScalingOverTimePeriod)),
				AutoScalingPauseLimit:     core.Int64Ptr(int64(autoScalingPauseLimit)),
			}

			result, response, err := db2SaasClient.PutDb2SaasAutoscale(input)
			if err != nil {
				log.Printf("Error while updating autoscaling config to DB2Saas: %s", err)
			} else {
				log.Printf("StatusCode of response %d", response.StatusCode)
				log.Printf("Success result %v", result)
			}
		}

	}

	if customSettingRaw, ok := d.GetOk("custom_setting_config"); ok {
		if customSettingRaw == nil || reflect.ValueOf(customSettingRaw).IsNil() {
			fmt.Println("No custom setting configs provided; skipping.")
		} else {
			customerSettingConfig := customSettingRaw.([]interface{})[0].(map[string]interface{})

			var (
				registryConfig, dbConfig, dbmConfig map[string]interface{}
				registry                            *db2saasv1.CreateCustomSettingsRegistry
				db                                  *db2saasv1.CreateCustomSettingsDb
				dbm                                 *db2saasv1.CreateCustomSettingsDbm
			)

			registryConfigRaw := customerSettingConfig["registry"]
			if registryConfigRaw == nil || reflect.ValueOf(registryConfigRaw).IsNil() {
				fmt.Println("No custom setting registry configs provided; skipping.")
			} else {
				if reflect.TypeOf(registryConfigRaw).Kind() == reflect.Slice {
					registryConfig = registryConfigRaw.([]interface{})[0].(map[string]interface{})

					registry = &db2saasv1.CreateCustomSettingsRegistry{
						DB2BIDI:                      checkStringNilValue(registryConfig, "DB2BIDI"),
						DB2COMPOPT:                   checkStringNilValue(registryConfig, "DB2COMPOPT"),
						DB2LOCKTORB:                  checkStringNilValue(registryConfig, "DB2LOCK_TO_RB"),
						DB2STMM:                      checkStringNilValue(registryConfig, "DB2STMM"),
						DB2ALTERNATEAUTHZBEHAVIOUR:   checkStringNilValue(registryConfig, "DB2_ALTERNATE_AUTHZ_BEHAVIOUR"),
						DB2ANTIJOIN:                  checkStringNilValue(registryConfig, "DB2_ANTIJOIN"),
						DB2ATSENABLE:                 checkStringNilValue(registryConfig, "DB2_ATS_ENABLE"),
						DB2DEFERREDPREPARESEMANTICS:  checkStringNilValue(registryConfig, "DB2_DEFERRED_PREPARE_SEMANTICS"),
						DB2EVALUNCOMMITTED:           checkStringNilValue(registryConfig, "DB2_EVALUNCOMMITTED"),
						DB2EXTENDEDOPTIMIZATION:      checkStringNilValue(registryConfig, "DB2_EXTENDED_OPTIMIZATION"),
						DB2INDEXPCTFREEDEFAULT:       checkStringNilValue(registryConfig, "DB2_INDEX_PCTFREE_DEFAULT"),
						DB2INLISTTONLJN:              checkStringNilValue(registryConfig, "DB2_INLIST_TO_NLJN"),
						DB2MINIMIZELISTPREFETCH:      checkStringNilValue(registryConfig, "DB2_MINIMIZE_LISTPREFETCH"),
						DB2OBJECTTABLEENTRIES:        checkStringNilValue(registryConfig, "DB2_OBJECT_TABLE_ENTRIES"),
						DB2OPTPROFILE:                checkStringNilValue(registryConfig, "DB2_OPTPROFILE"),
						DB2OPTSTATSLOG:               checkStringNilValue(registryConfig, "DB2_OPTSTATS_LOG"),
						DB2OPTMAXTEMPSIZE:            checkStringNilValue(registryConfig, "DB2_OPT_MAX_TEMP_SIZE"),
						DB2PARALLELIO:                checkStringNilValue(registryConfig, "DB2_PARALLEL_IO"),
						DB2REDUCEDOPTIMIZATION:       checkStringNilValue(registryConfig, "DB2_REDUCED_OPTIMIZATION"),
						DB2SELECTIVITY:               checkStringNilValue(registryConfig, "DB2_SELECTIVITY"),
						DB2SKIPDELETED:               checkStringNilValue(registryConfig, "DB2_SKIPDELETED"),
						DB2SKIPINSERTED:              checkStringNilValue(registryConfig, "DB2_SKIPINSERTED"),
						DB2SYNCRELEASELOCKATTRIBUTES: checkStringNilValue(registryConfig, "DB2_SYNC_RELEASE_LOCK_ATTRIBUTES"),
						DB2TRUNCATEREUSESTORAGE:      checkStringNilValue(registryConfig, "DB2_TRUNCATE_REUSESTORAGE"),
						DB2USEALTERNATEPAGECLEANING:  checkStringNilValue(registryConfig, "DB2_USE_ALTERNATE_PAGE_CLEANING"),
						DB2VIEWREOPTVALUES:           checkStringNilValue(registryConfig, "DB2_VIEW_REOPT_VALUES"),
						DB2WLMSETTINGS:               checkStringNilValue(registryConfig, "DB2_WLM_SETTINGS"),
						DB2WORKLOAD:                  checkStringNilValue(registryConfig, "DB2_WORKLOAD"),
					}
				} else {
					fmt.Println("Expected an array for registryConfig, but found another type")
				}

			}

			dbConfigRaw := customerSettingConfig["db"]
			if dbConfigRaw == nil || reflect.ValueOf(dbConfigRaw).IsNil() {
				fmt.Println("No custom setting db configs provided; skipping.")
			} else {
				if reflect.TypeOf(dbConfigRaw).Kind() == reflect.Slice {
					dbConfig = dbConfigRaw.([]interface{})[0].(map[string]interface{})

					db = &db2saasv1.CreateCustomSettingsDb{
						ACTSORTMEMLIMIT:    checkStringNilValue(dbConfig, "ACT_SORTMEM_LIMIT"),
						ALTCOLLATE:         checkStringNilValue(dbConfig, "ALT_COLLATE"),
						APPGROUPMEMSZ:      checkStringNilValue(dbConfig, "APPGROUP_MEM_SZ"),
						APPLHEAPSZ:         checkStringNilValue(dbConfig, "APPLHEAPSZ"),
						APPLMEMORY:         checkStringNilValue(dbConfig, "APPL_MEMORY"),
						APPCTLHEAPSZ:       checkStringNilValue(dbConfig, "APP_CTL_HEAP_SZ"),
						ARCHRETRYDELAY:     checkStringNilValue(dbConfig, "ARCHRETRYDELAY"),
						AUTHNCACHEDURATION: checkStringNilValue(dbConfig, "AUTHN_CACHE_DURATION"),
						AUTORESTART:        checkStringNilValue(dbConfig, "AUTORESTART"),
						AUTOCGSTATS:        checkStringNilValue(dbConfig, "AUTO_CG_STATS"),
						AUTOMAINT:          checkStringNilValue(dbConfig, "AUTO_MAINT"),
						AUTOREORG:          checkStringNilValue(dbConfig, "AUTO_REORG"),
						AUTOREVAL:          checkStringNilValue(dbConfig, "AUTO_REVAL"),
						AUTORUNSTATS:       checkStringNilValue(dbConfig, "AUTO_RUNSTATS"),
						AUTOSAMPLING:       checkStringNilValue(dbConfig, "AUTO_SAMPLING"),
						AUTOSTATSVIEWS:     checkStringNilValue(dbConfig, "AUTO_STATS_VIEWS"),
						AUTOSTMTSTATS:      checkStringNilValue(dbConfig, "AUTO_STMT_STATS"),
						AUTOTBLMAINT:       checkStringNilValue(dbConfig, "AUTO_TBL_MAINT"),
						AVGAPPLS:           checkStringNilValue(dbConfig, "AVG_APPLS"),
						CATALOGCACHESZ:     checkStringNilValue(dbConfig, "CATALOGCACHE_SZ"),
						CHNGPGSTHRESH:      checkStringNilValue(dbConfig, "CHNGPGS_THRESH"),
						CURCOMMIT:          checkStringNilValue(dbConfig, "CUR_COMMIT"),
						DATABASEMEMORY:     checkStringNilValue(dbConfig, "DATABASE_MEMORY"),
						DBHEAP:             checkStringNilValue(dbConfig, "DBHEAP"),
						DBCOLLNAME:         checkStringNilValue(dbConfig, "DB_COLLNAME"),
						DBMEMTHRESH:        checkStringNilValue(dbConfig, "DB_MEM_THRESH"),
						DDLCOMPRESSIONDEF:  checkStringNilValue(dbConfig, "DDL_COMPRESSION_DEF"),
						DDLCONSTRAINTDEF:   checkStringNilValue(dbConfig, "DDL_CONSTRAINT_DEF"),
						DECFLTROUNDING:     checkStringNilValue(dbConfig, "DECFLT_ROUNDING"),
						DECARITHMETIC:      checkStringNilValue(dbConfig, "DEC_ARITHMETIC"),
						DECTOCHARFMT:       checkStringNilValue(dbConfig, "DEC_TO_CHAR_FMT"),
						DFTDEGREE:          checkStringNilValue(dbConfig, "DFT_DEGREE"),
						DFTEXTENTSZ:        checkStringNilValue(dbConfig, "DFT_EXTENT_SZ"),
						DFTLOADRECSES:      checkStringNilValue(dbConfig, "DFT_LOADREC_SES"),
						DFTMTTBTYPES:       checkStringNilValue(dbConfig, "DFT_MTTB_TYPES"),
						DFTPREFETCHSZ:      checkStringNilValue(dbConfig, "DFT_PREFETCH_SZ"),
						DFTQUERYOPT:        checkStringNilValue(dbConfig, "DFT_QUERYOPT"),
						DFTREFRESHAGE:      checkStringNilValue(dbConfig, "DFT_REFRESH_AGE"),
						DFTSCHEMASDCC:      checkStringNilValue(dbConfig, "DFT_SCHEMAS_DCC"),
						DFTSQLMATHWARN:     checkStringNilValue(dbConfig, "DFT_SQLMATHWARN"),
						DFTTABLEORG:        checkStringNilValue(dbConfig, "DFT_TABLE_ORG"),
						DLCHKTIME:          checkStringNilValue(dbConfig, "DLCHKTIME"),
						ENABLEXMLCHAR:      checkStringNilValue(dbConfig, "ENABLE_XMLCHAR"),
						EXTENDEDROWSZ:      checkStringNilValue(dbConfig, "EXTENDED_ROW_SZ"),
						GROUPHEAPRATIO:     checkStringNilValue(dbConfig, "GROUPHEAP_RATIO"),
						INDEXREC:           checkStringNilValue(dbConfig, "INDEXREC"),
						LARGEAGGREGATION:   checkStringNilValue(dbConfig, "LARGE_AGGREGATION"),
						LOCKLIST:           checkStringNilValue(dbConfig, "LOCKLIST"),
						LOCKTIMEOUT:        checkStringNilValue(dbConfig, "LOCKTIMEOUT"),
						LOGINDEXBUILD:      checkStringNilValue(dbConfig, "LOGINDEXBUILD"),
						LOGAPPLINFO:        checkStringNilValue(dbConfig, "LOG_APPL_INFO"),
						LOGDDLSTMTS:        checkStringNilValue(dbConfig, "LOG_DDL_STMTS"),
						LOGDISKCAP:         checkStringNilValue(dbConfig, "LOG_DISK_CAP"),
						MAXAPPLS:           checkStringNilValue(dbConfig, "MAXAPPLS"),
						MAXFILOP:           checkStringNilValue(dbConfig, "MAXFILOP"),
						MAXLOCKS:           checkStringNilValue(dbConfig, "MAXLOCKS"),
						MINDECDIV3:         checkStringNilValue(dbConfig, "MIN_DEC_DIV_3"),
						MONACTMETRICS:      checkStringNilValue(dbConfig, "MON_ACT_METRICS"),
						MONDEADLOCK:        checkStringNilValue(dbConfig, "MON_DEADLOCK"),
						MONLCKMSGLVL:       checkStringNilValue(dbConfig, "MON_LCK_MSG_LVL"),
						MONLOCKTIMEOUT:     checkStringNilValue(dbConfig, "MON_LOCKTIMEOUT"),
						MONLOCKWAIT:        checkStringNilValue(dbConfig, "MON_LOCKWAIT"),
						MONLWTHRESH:        checkStringNilValue(dbConfig, "MON_LW_THRESH"),
						MONOBJMETRICS:      checkStringNilValue(dbConfig, "MON_OBJ_METRICS"),
						MONPKGLISTSZ:       checkStringNilValue(dbConfig, "MON_PKGLIST_SZ"),
						MONREQMETRICS:      checkStringNilValue(dbConfig, "MON_REQ_METRICS"),
						MONRTNDATA:         checkStringNilValue(dbConfig, "MON_RTN_DATA"),
						MONRTNEXECLIST:     checkStringNilValue(dbConfig, "MON_RTN_EXECLIST"),
						MONUOWDATA:         checkStringNilValue(dbConfig, "MON_UOW_DATA"),
						MONUOWEXECLIST:     checkStringNilValue(dbConfig, "MON_UOW_EXECLIST"),
						MONUOWPKGLIST:      checkStringNilValue(dbConfig, "MON_UOW_PKGLIST"),
						NCHARMAPPING:       checkStringNilValue(dbConfig, "NCHAR_MAPPING"),
						NUMFREQVALUES:      checkStringNilValue(dbConfig, "NUM_FREQVALUES"),
						NUMIOCLEANERS:      checkStringNilValue(dbConfig, "NUM_IOCLEANERS"),
						NUMIOSERVERS:       checkStringNilValue(dbConfig, "NUM_IOSERVERS"),
						NUMLOGSPAN:         checkStringNilValue(dbConfig, "NUM_LOG_SPAN"),
						NUMQUANTILES:       checkStringNilValue(dbConfig, "NUM_QUANTILES"),
						OPTBUFFPAGE:        checkStringNilValue(dbConfig, "OPT_BUFFPAGE"),
						OPTDIRECTWRKLD:     checkStringNilValue(dbConfig, "OPT_DIRECT_WRKLD"),
						OPTLOCKLIST:        checkStringNilValue(dbConfig, "OPT_LOCKLIST"),
						OPTMAXLOCKS:        checkStringNilValue(dbConfig, "OPT_MAXLOCKS"),
						OPTSORTHEAP:        checkStringNilValue(dbConfig, "OPT_SORTHEAP"),
						PAGEAGETRGTGCR:     checkStringNilValue(dbConfig, "PAGE_AGE_TRGT_GCR"),
						PAGEAGETRGTMCR:     checkStringNilValue(dbConfig, "PAGE_AGE_TRGT_MCR"),
						PCKCACHESZ:         checkStringNilValue(dbConfig, "PCKCACHESZ"),
						PLSTACKTRACE:       checkStringNilValue(dbConfig, "PL_STACK_TRACE"),
						SELFTUNINGMEM:      checkStringNilValue(dbConfig, "SELF_TUNING_MEM"),
						SEQDETECT:          checkStringNilValue(dbConfig, "SEQDETECT"),
						SHEAPTHRESSHR:      checkStringNilValue(dbConfig, "SHEAPTHRES_SHR"),
						SOFTMAX:            checkStringNilValue(dbConfig, "SOFTMAX"),
						SORTHEAP:           checkStringNilValue(dbConfig, "SORTHEAP"),
						SQLCCFLAGS:         checkStringNilValue(dbConfig, "SQL_CCFLAGS"),
						STATHEAPSZ:         checkStringNilValue(dbConfig, "STAT_HEAP_SZ"),
						STMTHEAP:           checkStringNilValue(dbConfig, "STMTHEAP"),
						STMTCONC:           checkStringNilValue(dbConfig, "STMT_CONC"),
						STRINGUNITS:        checkStringNilValue(dbConfig, "STRING_UNITS"),
						SYSTIMEPERIODADJ:   checkStringNilValue(dbConfig, "SYSTIME_PERIOD_ADJ"),
						TRACKMOD:           checkStringNilValue(dbConfig, "TRACKMOD"),
						UTILHEAPSZ:         checkStringNilValue(dbConfig, "UTIL_HEAP_SZ"),
						WLMADMISSIONCTRL:   checkStringNilValue(dbConfig, "WLM_ADMISSION_CTRL"),
						WLMAGENTLOADTRGT:   checkStringNilValue(dbConfig, "WLM_AGENT_LOAD_TRGT"),
						WLMCPULIMIT:        checkStringNilValue(dbConfig, "WLM_CPU_LIMIT"),
						WLMCPUSHARES:       checkStringNilValue(dbConfig, "WLM_CPU_SHARES"),
						WLMCPUSHAREMODE:    checkStringNilValue(dbConfig, "WLM_CPU_SHARE_MODE"),
					}
				} else {
					fmt.Println("Expected an array for dbConfig, but found another type")
				}
			}

			dbmConfigRaw := customerSettingConfig["dbm"]
			if dbmConfigRaw == nil || reflect.ValueOf(dbmConfigRaw).IsNil() {
				fmt.Println("No custom setting dbm configs provided; skipping.")
			} else {
				if reflect.TypeOf(dbmConfigRaw).Kind() == reflect.Slice {
					dbmConfig = dbmConfigRaw.([]interface{})[0].(map[string]interface{})

					dbm = &db2saasv1.CreateCustomSettingsDbm{
						COMMBANDWIDTH:    checkStringNilValue(dbmConfig, "COMM_BANDWIDTH"),
						CPUSPEED:         checkStringNilValue(dbmConfig, "CPUSPEED"),
						DFTMONBUFPOOL:    checkStringNilValue(dbmConfig, "DFT_MON_BUFPOOL"),
						DFTMONLOCK:       checkStringNilValue(dbmConfig, "DFT_MON_LOCK"),
						DFTMONSORT:       checkStringNilValue(dbmConfig, "DFT_MON_SORT"),
						DFTMONSTMT:       checkStringNilValue(dbmConfig, "DFT_MON_STMT"),
						DFTMONTABLE:      checkStringNilValue(dbmConfig, "DFT_MON_TABLE"),
						DFTMONTIMESTAMP:  checkStringNilValue(dbmConfig, "DFT_MON_TIMESTAMP"),
						DFTMONUOW:        checkStringNilValue(dbmConfig, "DFT_MON_UOW"),
						DIAGLEVEL:        checkStringNilValue(dbmConfig, "DIAGLEVEL"),
						FEDERATEDASYNC:   checkStringNilValue(dbmConfig, "FEDERATED_ASYNC"),
						INDEXREC:         checkStringNilValue(dbmConfig, "INDEXREC"),
						INTRAPARALLEL:    checkStringNilValue(dbmConfig, "INTRA_PARALLEL"),
						KEEPFENCED:       checkStringNilValue(dbmConfig, "KEEPFENCED"),
						MAXCONNRETRIES:   checkStringNilValue(dbmConfig, "MAX_CONNRETRIES"),
						MAXQUERYDEGREE:   checkStringNilValue(dbmConfig, "MAX_QUERYDEGREE"),
						MONHEAPSZ:        checkStringNilValue(dbmConfig, "MON_HEAP_SZ"),
						MULTIPARTSIZEMB:  checkStringNilValue(dbmConfig, "MULTIPARTSIZEMB"),
						NOTIFYLEVEL:      checkStringNilValue(dbmConfig, "NOTIFYLEVEL"),
						NUMINITAGENTS:    checkStringNilValue(dbmConfig, "NUM_INITAGENTS"),
						NUMINITFENCED:    checkStringNilValue(dbmConfig, "NUM_INITFENCED"),
						NUMPOOLAGENTS:    checkStringNilValue(dbmConfig, "NUM_POOLAGENTS"),
						RESYNCINTERVAL:   checkStringNilValue(dbmConfig, "RESYNC_INTERVAL"),
						RQRIOBLK:         checkStringNilValue(dbmConfig, "RQRIOBLK"),
						STARTSTOPTIME:    checkStringNilValue(dbmConfig, "START_STOP_TIME"),
						UTILIMPACTLIM:    checkStringNilValue(dbmConfig, "UTIL_IMPACT_LIM"),
						WLMDISPATCHER:    checkStringNilValue(dbmConfig, "WLM_DISPATCHER"),
						WLMDISPCONCUR:    checkStringNilValue(dbmConfig, "WLM_DISP_CONCUR"),
						WLMDISPCPUSHARES: checkStringNilValue(dbmConfig, "WLM_DISP_CPU_SHARES"),
						WLMDISPMINUTIL:   checkStringNilValue(dbmConfig, "WLM_DISP_MIN_UTIL"),
					}
				} else {
					fmt.Println("Expected an array for dbmConfig, but found another type")
				}
			}

			input := &db2saasv1.PostDb2SaasDbConfigurationOptions{
				XDbProfile: core.StringPtr(encodedCRN),
				Registry:   registry,
				Db:         db,
				Dbm:        dbm,
			}

			result, response, err := db2SaasClient.PostDb2SaasDbConfiguration(input)
			if err != nil {
				log.Printf("Error while posting DB configuration to DB2Saas: %s", err)
			} else {
				log.Printf("StatusCode of response %d", response.StatusCode)
				log.Printf("Success result %v", result)
			}
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}

	return resourcecontroller.ResourceIBMResourceInstanceRead(d, meta)
}

func waitForResourceInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{RsInstanceProgressStatus, RsInstanceInactiveStatus, RsInstanceProvisioningStatus},
		Target:  []string{RsInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get the resource instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance '%s' creation failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}

func checkStringNilValue(config map[string]interface{}, key string) *string {
	if value, ok := config[key]; ok && value != nil {
		strValue, isString := value.(string)
		if isString {
			return core.StringPtr(strValue)
		}
	}

	return nil
}
