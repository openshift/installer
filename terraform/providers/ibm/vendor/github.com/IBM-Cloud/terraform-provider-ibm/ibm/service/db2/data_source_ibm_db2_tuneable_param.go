// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package db2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/cloud-db2-go-sdk/db2saasv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func DataSourceIbmDb2TuneableParam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmDb2TuneableParamRead,

		Schema: map[string]*schema.Schema{
			"tuneable_param": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": &schema.Schema{
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Tunable parameters related to the Db2 database instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"act_sortmem_limit": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"alt_collate": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"appgroup_mem_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"applheapsz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"appl_memory": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"app_ctl_heap_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"archretrydelay": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"authn_cache_duration": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"autorestart": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_cg_stats": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_maint": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_reorg": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_reval": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_runstats": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_sampling": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_stats_views": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_stmt_stats": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"auto_tbl_maint": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"avg_appls": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"catalogcache_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"chngpgs_thresh": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"cur_commit": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"database_memory": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dbheap": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_collname": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_mem_thresh": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"ddl_compression_def": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"ddl_constraint_def": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"decflt_rounding": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dec_arithmetic": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dec_to_char_fmt": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_degree": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_extent_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_loadrec_ses": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mttb_types": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_prefetch_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_queryopt": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_refresh_age": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_schemas_dcc": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_sqlmathwarn": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_table_org": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dlchktime": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable_xmlchar": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"extended_row_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"groupheap_ratio": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"indexrec": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"large_aggregation": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"locklist": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"locktimeout": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"logindexbuild": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_appl_info": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_ddl_stmts": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"log_disk_cap": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"maxappls": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"maxfilop": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"maxlocks": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"min_dec_div_3": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_act_metrics": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_deadlock": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_lck_msg_lvl": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_locktimeout": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_lockwait": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_lw_thresh": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_obj_metrics": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_pkglist_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_req_metrics": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_rtn_data": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_rtn_execlist": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_uow_data": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_uow_execlist": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_uow_pkglist": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"nchar_mapping": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_freqvalues": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_iocleaners": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_ioservers": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_log_span": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_quantiles": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"opt_buffpage": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"opt_direct_wrkld": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"opt_locklist": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"opt_maxlocks": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"opt_sortheap": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"page_age_trgt_gcr": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"page_age_trgt_mcr": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"pckcachesz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"pl_stack_trace": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"self_tuning_mem": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"seqdetect": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"sheapthres_shr": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"softmax": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"sortheap": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"sql_ccflags": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"stat_heap_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"stmtheap": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"stmt_conc": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"string_units": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"systime_period_adj": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"trackmod": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"util_heap_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_admission_ctrl": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_agent_load_trgt": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_cpu_limit": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_cpu_shares": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_cpu_share_mode": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dbm": &schema.Schema{
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Tunable parameters related to the Db2 instance manager (dbm).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comm_bandwidth": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpuspeed": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_bufpool": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_lock": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_sort": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_stmt": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_table": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_timestamp": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"dft_mon_uow": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"diaglevel": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"federated_async": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"indexrec": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"intra_parallel": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"keepfenced": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_connretries": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_querydegree": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mon_heap_sz": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"multipartsizemb": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"notifylevel": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_initagents": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_initfenced": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"num_poolagents": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"resync_interval": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"rqrioblk": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_stop_time": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"util_impact_lim": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_dispatcher": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_disp_concur": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_disp_cpu_shares": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"wlm_disp_min_util": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"registry": &schema.Schema{
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Tunable parameters related to the Db2 registry.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db2_bidi": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_compopt": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_lock_to_rb": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_stmm": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_alternate_authz_behaviour": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_antijoin": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_ats_enable": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_deferred_prepare_semantics": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_evaluncommitted": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_extended_optimization": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_index_pctfree_default": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_inlist_to_nljn": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_minimize_listprefetch": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_object_table_entries": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_optprofile": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_optstats_log": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_opt_max_temp_size": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_parallel_io": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_reduced_optimization": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_selectivity": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_skipdeleted": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_skipinserted": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_sync_release_lock_attributes": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_truncate_reusestorage": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_use_alternate_page_cleaning": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_view_reopt_values": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_wlm_settings": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"db2_workload": &schema.Schema{
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

func dataSourceIbmDb2TuneableParamRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	db2saasClient, err := meta.(conns.ClientSession).Db2saasV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_db2_tuneable_param", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDb2SaasTuneableParamOptions := &db2saasv1.GetDb2SaasTuneableParamOptions{}

	successTuneableParams, _, err := db2saasClient.GetDb2SaasTuneableParamWithContext(context, getDb2SaasTuneableParamOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDb2SaasTuneableParamWithContext failed: %s", err.Error()), "(Data) ibm_db2_tuneable_param", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmDb2TuneableParamID(d))

	if !core.IsNil(successTuneableParams.TuneableParam) {
		tuneableParam := []map[string]interface{}{}
		tuneableParamMap, err := DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamToMap(successTuneableParams.TuneableParam)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ", "read", "tuneable_param-to-map").GetDiag()
		}
		tuneableParam = append(tuneableParam, tuneableParamMap)
		if err = d.Set("tuneable_param", tuneableParam); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tuneable_param: %s", err), "(Data) ibm_db2_tuneable_param", "read", "set-tuneable_param").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmDb2SaasTuneableParamID returns a reasonable ID for the list.
func dataSourceIbmDb2TuneableParamID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamToMap(model *db2saasv1.SuccessTuneableParamsTuneableParam) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Db != nil {
		dbMap, err := DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamDbToMap(model.Db)
		if err != nil {
			return modelMap, err
		}
		modelMap["db"] = []map[string]interface{}{dbMap}
	}
	if model.Dbm != nil {
		dbmMap, err := DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamDbmToMap(model.Dbm)
		if err != nil {
			return modelMap, err
		}
		modelMap["dbm"] = []map[string]interface{}{dbmMap}
	}
	if model.Registry != nil {
		registryMap, err := DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamRegistryToMap(model.Registry)
		if err != nil {
			return modelMap, err
		}
		modelMap["registry"] = []map[string]interface{}{registryMap}
	}
	return modelMap, nil
}

func DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamDbToMap(model *db2saasv1.SuccessTuneableParamsTuneableParamDb) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ACTSORTMEMLIMIT != nil {
		modelMap["act_sortmem_limit"] = *model.ACTSORTMEMLIMIT
	}
	if model.ALTCOLLATE != nil {
		modelMap["alt_collate"] = *model.ALTCOLLATE
	}
	if model.APPGROUPMEMSZ != nil {
		modelMap["appgroup_mem_sz"] = *model.APPGROUPMEMSZ
	}
	if model.APPLHEAPSZ != nil {
		modelMap["applheapsz"] = *model.APPLHEAPSZ
	}
	if model.APPLMEMORY != nil {
		modelMap["appl_memory"] = *model.APPLMEMORY
	}
	if model.APPCTLHEAPSZ != nil {
		modelMap["app_ctl_heap_sz"] = *model.APPCTLHEAPSZ
	}
	if model.ARCHRETRYDELAY != nil {
		modelMap["archretrydelay"] = *model.ARCHRETRYDELAY
	}
	if model.AUTHNCACHEDURATION != nil {
		modelMap["authn_cache_duration"] = *model.AUTHNCACHEDURATION
	}
	if model.AUTORESTART != nil {
		modelMap["autorestart"] = *model.AUTORESTART
	}
	if model.AUTOCGSTATS != nil {
		modelMap["auto_cg_stats"] = *model.AUTOCGSTATS
	}
	if model.AUTOMAINT != nil {
		modelMap["auto_maint"] = *model.AUTOMAINT
	}
	if model.AUTOREORG != nil {
		modelMap["auto_reorg"] = *model.AUTOREORG
	}
	if model.AUTOREVAL != nil {
		modelMap["auto_reval"] = *model.AUTOREVAL
	}
	if model.AUTORUNSTATS != nil {
		modelMap["auto_runstats"] = *model.AUTORUNSTATS
	}
	if model.AUTOSAMPLING != nil {
		modelMap["auto_sampling"] = *model.AUTOSAMPLING
	}
	if model.AUTOSTATSVIEWS != nil {
		modelMap["auto_stats_views"] = *model.AUTOSTATSVIEWS
	}
	if model.AUTOSTMTSTATS != nil {
		modelMap["auto_stmt_stats"] = *model.AUTOSTMTSTATS
	}
	if model.AUTOTBLMAINT != nil {
		modelMap["auto_tbl_maint"] = *model.AUTOTBLMAINT
	}
	if model.AVGAPPLS != nil {
		modelMap["avg_appls"] = *model.AVGAPPLS
	}
	if model.CATALOGCACHESZ != nil {
		modelMap["catalogcache_sz"] = *model.CATALOGCACHESZ
	}
	if model.CHNGPGSTHRESH != nil {
		modelMap["chngpgs_thresh"] = *model.CHNGPGSTHRESH
	}
	if model.CURCOMMIT != nil {
		modelMap["cur_commit"] = *model.CURCOMMIT
	}
	if model.DATABASEMEMORY != nil {
		modelMap["database_memory"] = *model.DATABASEMEMORY
	}
	if model.DBHEAP != nil {
		modelMap["dbheap"] = *model.DBHEAP
	}
	if model.DBCOLLNAME != nil {
		modelMap["db_collname"] = *model.DBCOLLNAME
	}
	if model.DBMEMTHRESH != nil {
		modelMap["db_mem_thresh"] = *model.DBMEMTHRESH
	}
	if model.DDLCOMPRESSIONDEF != nil {
		modelMap["ddl_compression_def"] = *model.DDLCOMPRESSIONDEF
	}
	if model.DDLCONSTRAINTDEF != nil {
		modelMap["ddl_constraint_def"] = *model.DDLCONSTRAINTDEF
	}
	if model.DECFLTROUNDING != nil {
		modelMap["decflt_rounding"] = *model.DECFLTROUNDING
	}
	if model.DECARITHMETIC != nil {
		modelMap["dec_arithmetic"] = *model.DECARITHMETIC
	}
	if model.DECTOCHARFMT != nil {
		modelMap["dec_to_char_fmt"] = *model.DECTOCHARFMT
	}
	if model.DFTDEGREE != nil {
		modelMap["dft_degree"] = *model.DFTDEGREE
	}
	if model.DFTEXTENTSZ != nil {
		modelMap["dft_extent_sz"] = *model.DFTEXTENTSZ
	}
	if model.DFTLOADRECSES != nil {
		modelMap["dft_loadrec_ses"] = *model.DFTLOADRECSES
	}
	if model.DFTMTTBTYPES != nil {
		modelMap["dft_mttb_types"] = *model.DFTMTTBTYPES
	}
	if model.DFTPREFETCHSZ != nil {
		modelMap["dft_prefetch_sz"] = *model.DFTPREFETCHSZ
	}
	if model.DFTQUERYOPT != nil {
		modelMap["dft_queryopt"] = *model.DFTQUERYOPT
	}
	if model.DFTREFRESHAGE != nil {
		modelMap["dft_refresh_age"] = *model.DFTREFRESHAGE
	}
	if model.DFTSCHEMASDCC != nil {
		modelMap["dft_schemas_dcc"] = *model.DFTSCHEMASDCC
	}
	if model.DFTSQLMATHWARN != nil {
		modelMap["dft_sqlmathwarn"] = *model.DFTSQLMATHWARN
	}
	if model.DFTTABLEORG != nil {
		modelMap["dft_table_org"] = *model.DFTTABLEORG
	}
	if model.DLCHKTIME != nil {
		modelMap["dlchktime"] = *model.DLCHKTIME
	}
	if model.ENABLEXMLCHAR != nil {
		modelMap["enable_xmlchar"] = *model.ENABLEXMLCHAR
	}
	if model.EXTENDEDROWSZ != nil {
		modelMap["extended_row_sz"] = *model.EXTENDEDROWSZ
	}
	if model.GROUPHEAPRATIO != nil {
		modelMap["groupheap_ratio"] = *model.GROUPHEAPRATIO
	}
	if model.INDEXREC != nil {
		modelMap["indexrec"] = *model.INDEXREC
	}
	if model.LARGEAGGREGATION != nil {
		modelMap["large_aggregation"] = *model.LARGEAGGREGATION
	}
	if model.LOCKLIST != nil {
		modelMap["locklist"] = *model.LOCKLIST
	}
	if model.LOCKTIMEOUT != nil {
		modelMap["locktimeout"] = *model.LOCKTIMEOUT
	}
	if model.LOGINDEXBUILD != nil {
		modelMap["logindexbuild"] = *model.LOGINDEXBUILD
	}
	if model.LOGAPPLINFO != nil {
		modelMap["log_appl_info"] = *model.LOGAPPLINFO
	}
	if model.LOGDDLSTMTS != nil {
		modelMap["log_ddl_stmts"] = *model.LOGDDLSTMTS
	}
	if model.LOGDISKCAP != nil {
		modelMap["log_disk_cap"] = *model.LOGDISKCAP
	}
	if model.MAXAPPLS != nil {
		modelMap["maxappls"] = *model.MAXAPPLS
	}
	if model.MAXFILOP != nil {
		modelMap["maxfilop"] = *model.MAXFILOP
	}
	if model.MAXLOCKS != nil {
		modelMap["maxlocks"] = *model.MAXLOCKS
	}
	if model.MINDECDIV3 != nil {
		modelMap["min_dec_div_3"] = *model.MINDECDIV3
	}
	if model.MONACTMETRICS != nil {
		modelMap["mon_act_metrics"] = *model.MONACTMETRICS
	}
	if model.MONDEADLOCK != nil {
		modelMap["mon_deadlock"] = *model.MONDEADLOCK
	}
	if model.MONLCKMSGLVL != nil {
		modelMap["mon_lck_msg_lvl"] = *model.MONLCKMSGLVL
	}
	if model.MONLOCKTIMEOUT != nil {
		modelMap["mon_locktimeout"] = *model.MONLOCKTIMEOUT
	}
	if model.MONLOCKWAIT != nil {
		modelMap["mon_lockwait"] = *model.MONLOCKWAIT
	}
	if model.MONLWTHRESH != nil {
		modelMap["mon_lw_thresh"] = *model.MONLWTHRESH
	}
	if model.MONOBJMETRICS != nil {
		modelMap["mon_obj_metrics"] = *model.MONOBJMETRICS
	}
	if model.MONPKGLISTSZ != nil {
		modelMap["mon_pkglist_sz"] = *model.MONPKGLISTSZ
	}
	if model.MONREQMETRICS != nil {
		modelMap["mon_req_metrics"] = *model.MONREQMETRICS
	}
	if model.MONRTNDATA != nil {
		modelMap["mon_rtn_data"] = *model.MONRTNDATA
	}
	if model.MONRTNEXECLIST != nil {
		modelMap["mon_rtn_execlist"] = *model.MONRTNEXECLIST
	}
	if model.MONUOWDATA != nil {
		modelMap["mon_uow_data"] = *model.MONUOWDATA
	}
	if model.MONUOWEXECLIST != nil {
		modelMap["mon_uow_execlist"] = *model.MONUOWEXECLIST
	}
	if model.MONUOWPKGLIST != nil {
		modelMap["mon_uow_pkglist"] = *model.MONUOWPKGLIST
	}
	if model.NCHARMAPPING != nil {
		modelMap["nchar_mapping"] = *model.NCHARMAPPING
	}
	if model.NUMFREQVALUES != nil {
		modelMap["num_freqvalues"] = *model.NUMFREQVALUES
	}
	if model.NUMIOCLEANERS != nil {
		modelMap["num_iocleaners"] = *model.NUMIOCLEANERS
	}
	if model.NUMIOSERVERS != nil {
		modelMap["num_ioservers"] = *model.NUMIOSERVERS
	}
	if model.NUMLOGSPAN != nil {
		modelMap["num_log_span"] = *model.NUMLOGSPAN
	}
	if model.NUMQUANTILES != nil {
		modelMap["num_quantiles"] = *model.NUMQUANTILES
	}
	if model.OPTBUFFPAGE != nil {
		modelMap["opt_buffpage"] = *model.OPTBUFFPAGE
	}
	if model.OPTDIRECTWRKLD != nil {
		modelMap["opt_direct_wrkld"] = *model.OPTDIRECTWRKLD
	}
	if model.OPTLOCKLIST != nil {
		modelMap["opt_locklist"] = *model.OPTLOCKLIST
	}
	if model.OPTMAXLOCKS != nil {
		modelMap["opt_maxlocks"] = *model.OPTMAXLOCKS
	}
	if model.OPTSORTHEAP != nil {
		modelMap["opt_sortheap"] = *model.OPTSORTHEAP
	}
	if model.PAGEAGETRGTGCR != nil {
		modelMap["page_age_trgt_gcr"] = *model.PAGEAGETRGTGCR
	}
	if model.PAGEAGETRGTMCR != nil {
		modelMap["page_age_trgt_mcr"] = *model.PAGEAGETRGTMCR
	}
	if model.PCKCACHESZ != nil {
		modelMap["pckcachesz"] = *model.PCKCACHESZ
	}
	if model.PLSTACKTRACE != nil {
		modelMap["pl_stack_trace"] = *model.PLSTACKTRACE
	}
	if model.SELFTUNINGMEM != nil {
		modelMap["self_tuning_mem"] = *model.SELFTUNINGMEM
	}
	if model.SEQDETECT != nil {
		modelMap["seqdetect"] = *model.SEQDETECT
	}
	if model.SHEAPTHRESSHR != nil {
		modelMap["sheapthres_shr"] = *model.SHEAPTHRESSHR
	}
	if model.SOFTMAX != nil {
		modelMap["softmax"] = *model.SOFTMAX
	}
	if model.SORTHEAP != nil {
		modelMap["sortheap"] = *model.SORTHEAP
	}
	if model.SQLCCFLAGS != nil {
		modelMap["sql_ccflags"] = *model.SQLCCFLAGS
	}
	if model.STATHEAPSZ != nil {
		modelMap["stat_heap_sz"] = *model.STATHEAPSZ
	}
	if model.STMTHEAP != nil {
		modelMap["stmtheap"] = *model.STMTHEAP
	}
	if model.STMTCONC != nil {
		modelMap["stmt_conc"] = *model.STMTCONC
	}
	if model.STRINGUNITS != nil {
		modelMap["string_units"] = *model.STRINGUNITS
	}
	if model.SYSTIMEPERIODADJ != nil {
		modelMap["systime_period_adj"] = *model.SYSTIMEPERIODADJ
	}
	if model.TRACKMOD != nil {
		modelMap["trackmod"] = *model.TRACKMOD
	}
	if model.UTILHEAPSZ != nil {
		modelMap["util_heap_sz"] = *model.UTILHEAPSZ
	}
	if model.WLMADMISSIONCTRL != nil {
		modelMap["wlm_admission_ctrl"] = *model.WLMADMISSIONCTRL
	}
	if model.WLMAGENTLOADTRGT != nil {
		modelMap["wlm_agent_load_trgt"] = *model.WLMAGENTLOADTRGT
	}
	if model.WLMCPULIMIT != nil {
		modelMap["wlm_cpu_limit"] = *model.WLMCPULIMIT
	}
	if model.WLMCPUSHARES != nil {
		modelMap["wlm_cpu_shares"] = *model.WLMCPUSHARES
	}
	if model.WLMCPUSHAREMODE != nil {
		modelMap["wlm_cpu_share_mode"] = *model.WLMCPUSHAREMODE
	}
	return modelMap, nil
}

func DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamDbmToMap(model *db2saasv1.SuccessTuneableParamsTuneableParamDbm) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.COMMBANDWIDTH != nil {
		modelMap["comm_bandwidth"] = *model.COMMBANDWIDTH
	}
	if model.CPUSPEED != nil {
		modelMap["cpuspeed"] = *model.CPUSPEED
	}
	if model.DFTMONBUFPOOL != nil {
		modelMap["dft_mon_bufpool"] = *model.DFTMONBUFPOOL
	}
	if model.DFTMONLOCK != nil {
		modelMap["dft_mon_lock"] = *model.DFTMONLOCK
	}
	if model.DFTMONSORT != nil {
		modelMap["dft_mon_sort"] = *model.DFTMONSORT
	}
	if model.DFTMONSTMT != nil {
		modelMap["dft_mon_stmt"] = *model.DFTMONSTMT
	}
	if model.DFTMONTABLE != nil {
		modelMap["dft_mon_table"] = *model.DFTMONTABLE
	}
	if model.DFTMONTIMESTAMP != nil {
		modelMap["dft_mon_timestamp"] = *model.DFTMONTIMESTAMP
	}
	if model.DFTMONUOW != nil {
		modelMap["dft_mon_uow"] = *model.DFTMONUOW
	}
	if model.DIAGLEVEL != nil {
		modelMap["diaglevel"] = *model.DIAGLEVEL
	}
	if model.FEDERATEDASYNC != nil {
		modelMap["federated_async"] = *model.FEDERATEDASYNC
	}
	if model.INDEXREC != nil {
		modelMap["indexrec"] = *model.INDEXREC
	}
	if model.INTRAPARALLEL != nil {
		modelMap["intra_parallel"] = *model.INTRAPARALLEL
	}
	if model.KEEPFENCED != nil {
		modelMap["keepfenced"] = *model.KEEPFENCED
	}
	if model.MAXCONNRETRIES != nil {
		modelMap["max_connretries"] = *model.MAXCONNRETRIES
	}
	if model.MAXQUERYDEGREE != nil {
		modelMap["max_querydegree"] = *model.MAXQUERYDEGREE
	}
	if model.MONHEAPSZ != nil {
		modelMap["mon_heap_sz"] = *model.MONHEAPSZ
	}
	if model.MULTIPARTSIZEMB != nil {
		modelMap["multipartsizemb"] = *model.MULTIPARTSIZEMB
	}
	if model.NOTIFYLEVEL != nil {
		modelMap["notifylevel"] = *model.NOTIFYLEVEL
	}
	if model.NUMINITAGENTS != nil {
		modelMap["num_initagents"] = *model.NUMINITAGENTS
	}
	if model.NUMINITFENCED != nil {
		modelMap["num_initfenced"] = *model.NUMINITFENCED
	}
	if model.NUMPOOLAGENTS != nil {
		modelMap["num_poolagents"] = *model.NUMPOOLAGENTS
	}
	if model.RESYNCINTERVAL != nil {
		modelMap["resync_interval"] = *model.RESYNCINTERVAL
	}
	if model.RQRIOBLK != nil {
		modelMap["rqrioblk"] = *model.RQRIOBLK
	}
	if model.STARTSTOPTIME != nil {
		modelMap["start_stop_time"] = *model.STARTSTOPTIME
	}
	if model.UTILIMPACTLIM != nil {
		modelMap["util_impact_lim"] = *model.UTILIMPACTLIM
	}
	if model.WLMDISPATCHER != nil {
		modelMap["wlm_dispatcher"] = *model.WLMDISPATCHER
	}
	if model.WLMDISPCONCUR != nil {
		modelMap["wlm_disp_concur"] = *model.WLMDISPCONCUR
	}
	if model.WLMDISPCPUSHARES != nil {
		modelMap["wlm_disp_cpu_shares"] = *model.WLMDISPCPUSHARES
	}
	if model.WLMDISPMINUTIL != nil {
		modelMap["wlm_disp_min_util"] = *model.WLMDISPMINUTIL
	}
	return modelMap, nil
}

func DataSourceIbmDb2TuneableParamSuccessTuneableParamsTuneableParamRegistryToMap(model *db2saasv1.SuccessTuneableParamsTuneableParamRegistry) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DB2BIDI != nil {
		modelMap["db2_bidi"] = *model.DB2BIDI
	}
	if model.DB2COMPOPT != nil {
		modelMap["db2_compopt"] = *model.DB2COMPOPT
	}
	if model.DB2LOCKTORB != nil {
		modelMap["db2_lock_to_rb"] = *model.DB2LOCKTORB
	}
	if model.DB2STMM != nil {
		modelMap["db2_stmm"] = *model.DB2STMM
	}
	if model.DB2ALTERNATEAUTHZBEHAVIOUR != nil {
		modelMap["db2_alternate_authz_behaviour"] = *model.DB2ALTERNATEAUTHZBEHAVIOUR
	}
	if model.DB2ANTIJOIN != nil {
		modelMap["db2_antijoin"] = *model.DB2ANTIJOIN
	}
	if model.DB2ATSENABLE != nil {
		modelMap["db2_ats_enable"] = *model.DB2ATSENABLE
	}
	if model.DB2DEFERREDPREPARESEMANTICS != nil {
		modelMap["db2_deferred_prepare_semantics"] = *model.DB2DEFERREDPREPARESEMANTICS
	}
	if model.DB2EVALUNCOMMITTED != nil {
		modelMap["db2_evaluncommitted"] = *model.DB2EVALUNCOMMITTED
	}
	if model.DB2EXTENDEDOPTIMIZATION != nil {
		modelMap["db2_extended_optimization"] = *model.DB2EXTENDEDOPTIMIZATION
	}
	if model.DB2INDEXPCTFREEDEFAULT != nil {
		modelMap["db2_index_pctfree_default"] = *model.DB2INDEXPCTFREEDEFAULT
	}
	if model.DB2INLISTTONLJN != nil {
		modelMap["db2_inlist_to_nljn"] = *model.DB2INLISTTONLJN
	}
	if model.DB2MINIMIZELISTPREFETCH != nil {
		modelMap["db2_minimize_listprefetch"] = *model.DB2MINIMIZELISTPREFETCH
	}
	if model.DB2OBJECTTABLEENTRIES != nil {
		modelMap["db2_object_table_entries"] = *model.DB2OBJECTTABLEENTRIES
	}
	if model.DB2OPTPROFILE != nil {
		modelMap["db2_optprofile"] = *model.DB2OPTPROFILE
	}
	if model.DB2OPTSTATSLOG != nil {
		modelMap["db2_optstats_log"] = *model.DB2OPTSTATSLOG
	}
	if model.DB2OPTMAXTEMPSIZE != nil {
		modelMap["db2_opt_max_temp_size"] = *model.DB2OPTMAXTEMPSIZE
	}
	if model.DB2PARALLELIO != nil {
		modelMap["db2_parallel_io"] = *model.DB2PARALLELIO
	}
	if model.DB2REDUCEDOPTIMIZATION != nil {
		modelMap["db2_reduced_optimization"] = *model.DB2REDUCEDOPTIMIZATION
	}
	if model.DB2SELECTIVITY != nil {
		modelMap["db2_selectivity"] = *model.DB2SELECTIVITY
	}
	if model.DB2SKIPDELETED != nil {
		modelMap["db2_skipdeleted"] = *model.DB2SKIPDELETED
	}
	if model.DB2SKIPINSERTED != nil {
		modelMap["db2_skipinserted"] = *model.DB2SKIPINSERTED
	}
	if model.DB2SYNCRELEASELOCKATTRIBUTES != nil {
		modelMap["db2_sync_release_lock_attributes"] = *model.DB2SYNCRELEASELOCKATTRIBUTES
	}
	if model.DB2TRUNCATEREUSESTORAGE != nil {
		modelMap["db2_truncate_reusestorage"] = *model.DB2TRUNCATEREUSESTORAGE
	}
	if model.DB2USEALTERNATEPAGECLEANING != nil {
		modelMap["db2_use_alternate_page_cleaning"] = *model.DB2USEALTERNATEPAGECLEANING
	}
	if model.DB2VIEWREOPTVALUES != nil {
		modelMap["db2_view_reopt_values"] = *model.DB2VIEWREOPTVALUES
	}
	if model.DB2WLMSETTINGS != nil {
		modelMap["db2_wlm_settings"] = *model.DB2WLMSETTINGS
	}
	if model.DB2WORKLOAD != nil {
		modelMap["db2_workload"] = *model.DB2WORKLOAD
	}
	return modelMap, nil
}
