package nutanix

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixProtectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixProtectionRuleCreate,
		Read:   resourceNutanixProtectionRuleRead,
		Update: resourceNutanixProtectionRuleUpdate,
		Delete: resourceNutanixProtectionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_connectivity_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_availability_zone_index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"source_availability_zone_index": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"snapshot_schedule_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recovery_point_objective_secs": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"local_snapshot_retention_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"rollup_retention_policy_multiple": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"rollup_retention_policy_snapshot_interval_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"HOURLY", "DAILY", "WEEKLY", "MONTHLY", "YEARLY"}, false),
												},
											},
										},
									},
									"auto_suspend_timeout_secs": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"snapshot_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"remote_snapshot_retention_policy": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"rollup_retention_policy_multiple": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"rollup_retention_policy_snapshot_interval_type": {
													Type:         schema.TypeString,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.StringInSlice([]string{"HOURLY", "DAILY", "WEEKLY", "MONTHLY", "YEARLY"}, false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"ordered_availability_zone_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"availability_zone_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"category_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"CATEGORIES_MATCH_ALL", "CATEGORIES_MATCH_ANY"}, false),
						},
						"kind_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"params": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixProtectionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.ProtectionRuleInput{}
	spec := &v3.ProtectionRuleSpec{}
	metadata := &v3.Metadata{}
	protectionRule := &v3.ProtectionRuleResources{}

	n, nok := d.GetOk("name")
	_, azclok := d.GetOk("availability_zone_connectivity_list")
	_, oazlok := d.GetOk("ordered_availability_zone_list")
	desc, descok := d.GetOk("description")

	if !nok && !azclok && !oazlok {
		return fmt.Errorf("please provide the required attributes `name`, `availability_zone`, `ordered_availability_zone`")
	}

	if err := getMetadataAttributes(d, metadata, "protection_rule"); err != nil {
		return err
	}

	getProtectionRulesResources(d, protectionRule)

	if descok {
		spec.Description = desc.(string)
	}

	protectionUUID, err := resourceNutanixProtectionRulesExists(conn, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error checking if protection_rule already exists %+v", err)
	}

	if protectionUUID != nil {
		return fmt.Errorf("protection_rule already with name %s exists , UUID %s", d.Get("name").(string), *protectionUUID)
	}

	spec.Name = n.(string)
	spec.Resources = protectionRule
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateProtectionRule(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix ProtectionRules %s: %+v", spec.Name, err)
	}

	d.SetId(*resp.Metadata.UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the ProtectionRules to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for protection_rule id (%s) to create: %+v", id, err)
	}

	// Setting Description because in Get request is not present.
	d.Set("description", resp.Spec.Description)

	return resourceNutanixProtectionRuleRead(d, meta)
}

func resourceNutanixProtectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	resp, err := conn.V3.GetProtectionRule(id)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValuesList(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValuesList(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return err
	}
	if err := d.Set("description", resp.Spec.Description); err != nil {
		return err
	}
	if err := d.Set("start_time", resp.Spec.Resources.StartTime); err != nil {
		return err
	}
	if err := d.Set("category_filter", flattenCategoriesFilter(resp.Spec.Resources.CategoryFilter)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_connectivity_list",
		flattenAvailabilityZoneConnectivityList(resp.Spec.Resources.AvailabilityZoneConnectivityList)); err != nil {
		return err
	}
	if err := d.Set("ordered_availability_zone_list",
		flattenOrderAvailibilityList(resp.Spec.Resources.OrderedAvailabilityZoneList)); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixProtectionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.ProtectionRuleInput{}
	metadata := &v3.Metadata{}
	spec := &v3.ProtectionRuleSpec{}

	id := d.Id()
	response, err := conn.V3.GetProtectionRule(id)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for protection rule id (%s) :%+v", id, err)
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec
		if response.Spec.Resources != nil {
			spec.Resources = response.Spec.Resources
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").([]interface{})
		metadata.OwnerReference = validateRefList(or, utils.StringPtr("protection_rule"))
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").([]interface{})
		metadata.ProjectReference = validateRefList(pr, utils.StringPtr("protection_rule"))
	}
	if d.HasChange("name") {
		spec.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		spec.Description = d.Get("description").(string)
	}
	if d.HasChange("start_time") {
		spec.Resources.StartTime = d.Get("start_time").(string)
	}
	if d.HasChange("availability_zone_connectivity_list") {
		spec.Resources.AvailabilityZoneConnectivityList = expandAvailabilityZoneConnectivityList(d)
	}
	if d.HasChange("ordered_availability_zone_list") {
		spec.Resources.OrderedAvailabilityZoneList = expandOrderAvailibilityList(d)
	}
	if d.HasChange("category_filter") {
		spec.Resources.CategoryFilter = expandCategoryFilter(d)
	}

	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateProtectionRule(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error updating protection_rule id %s): %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for protection rule (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixProtectionRuleRead(d, meta)
}

func resourceNutanixProtectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.DeleteProtectionRule(d.Id())

	if err != nil {
		return fmt.Errorf("error deleting protection_rule id %s): %s", d.Id(), err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for protection_rule (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixProtectionRulesExists(conn *v3.Client, name string) (*string, error) {
	var uuid *string

	filter := fmt.Sprintf("name==%s", name)
	protectionList, err := conn.V3.ListAllProtectionRules(filter)

	if err != nil {
		return nil, err
	}

	for _, protection := range protectionList.Entities {
		if protection.Status.Name == name {
			uuid = protection.Metadata.UUID
		}
	}
	return uuid, nil
}

func getProtectionRulesResources(d *schema.ResourceData, pr *v3.ProtectionRuleResources) {
	if v, ok := d.GetOk("start_time"); ok {
		pr.StartTime = v.(string)
	}
	if _, ok := d.GetOk("category_filter"); ok {
		pr.CategoryFilter = expandCategoryFilter(d)
	}
	pr.AvailabilityZoneConnectivityList = expandAvailabilityZoneConnectivityList(d)
	pr.OrderedAvailabilityZoneList = expandOrderAvailibilityList(d)
}

func expandCategoryFilter(d *schema.ResourceData) *v3.CategoryFilter {
	cf := &v3.CategoryFilter{}
	if v, ok := d.GetOk("category_filter.0.type"); ok {
		cf.Type = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("category_filter.0.kind_list"); ok {
		cf.KindList = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("category_filter.0.params"); ok {
		fpl := v.(*schema.Set).List()
		fl := make(map[string][]string)
		for _, v := range fpl {
			item := v.(map[string]interface{})
			if i, ok := item["name"]; ok && i.(string) != "" {
				if k2, kok := item["values"]; kok && len(k2.([]interface{})) > 0 {
					var values []string
					for _, item := range k2.([]interface{}) {
						values = append(values, item.(string))
					}
					fl[i.(string)] = values
				}
			}
		}
		cf.Params = fl
	}

	return cf
}

func expandAvailabilityZoneConnectivityList(d *schema.ResourceData) []*v3.AvailabilityZoneConnectivityList {
	zoneConnectivityLists := make([]*v3.AvailabilityZoneConnectivityList, 0)
	if v, ok := d.GetOk("availability_zone_connectivity_list"); ok {
		zoneList := v.([]interface{})
		for _, zone := range zoneList {
			azcl := &v3.AvailabilityZoneConnectivityList{}
			v1 := zone.(map[string]interface{})
			if v2, ok := v1["destination_availability_zone_index"]; ok {
				azcl.DestinationAvailabilityZoneIndex = utils.Int64Ptr(cast.ToInt64(v2))
			}
			if v2, ok := v1["source_availability_zone_index"]; ok {
				azcl.SourceAvailabilityZoneIndex = utils.Int64Ptr(cast.ToInt64(v2))
			}
			if v2, ok := v1["snapshot_schedule_list"]; ok {
				azcl.SnapshotScheduleList = expandSnapshotScheduleList(v2)
			}
			zoneConnectivityLists = append(zoneConnectivityLists, azcl)
		}
	}

	return zoneConnectivityLists
}

func expandSnapshotScheduleList(d interface{}) []*v3.SnapshotScheduleList {
	scl := d.([]interface{})
	snapshots := make([]*v3.SnapshotScheduleList, 0)
	if len(scl) > 0 {
		for _, val := range scl {
			v := val.(map[string]interface{})
			az := &v3.SnapshotScheduleList{}
			if v1, ok1 := v["recovery_point_objective_secs"]; ok1 {
				az.RecoveryPointObjectiveSecs = utils.Int64Ptr(cast.ToInt64(v1))
			}
			if v1, ok1 := v["local_snapshot_retention_policy"].([]interface{}); ok1 && len(v1) > 0 {
				az.LocalSnapshotRetentionPolicy = expandRetentionPolicy(v1[0])
			}
			if v1, ok1 := v["auto_suspend_timeout_secs"]; ok1 {
				az.AutoSuspendTimeoutSecs = utils.Int64Ptr(cast.ToInt64(v1))
			}
			if v1, ok1 := v["snapshot_type"]; ok1 && v1.(string) != "" {
				az.SnapshotType = v1.(string)
			}
			if v1, ok1 := v["remote_snapshot_retention_policy"].([]interface{}); ok1 && len(v1) > 0 {
				az.RemoteSnapshotRetentionPolicy = expandRetentionPolicy(v1[0])
			}
			snapshots = append(snapshots, az)
		}
	}

	return snapshots
}

func expandRetentionPolicy(d interface{}) *v3.SnapshotRetentionPolicy {
	srp := &v3.SnapshotRetentionPolicy{}
	rrp := &v3.RollupRetentionPolicy{}
	if d != nil {
		v := d.(map[string]interface{})
		log.Printf("[DEGUG] expandRetentionPolicy: %+v", v)
		flagRollup := false
		if v1, ok := v["num_snapshots"]; ok && v1.(int) != 0 {
			srp.NumSnapshots = utils.Int64Ptr(cast.ToInt64(v1))
			log.Printf("[DEGUG] srp.NumSnapshots: %+v", srp.NumSnapshots)
		}
		if v1, ok := v["rollup_retention_policy_multiple"]; ok && v1.(int) != 0 {
			rrp.Multiple = utils.Int64Ptr(cast.ToInt64(v1))
			log.Printf("[DEGUG] rrp.Multiple: %+v", rrp.Multiple)
			flagRollup = true
		}
		if v1, ok := v["rollup_retention_policy_snapshot_interval_type"]; ok && v1.(string) != "" {
			rrp.SnapshotIntervalType = v1.(string)
			flagRollup = true
		}
		if flagRollup {
			srp.RollupRetentionPolicy = rrp
		}
	}

	return srp
}

func expandOrderAvailibilityList(d *schema.ResourceData) []*v3.OrderedAvailabilityZoneList {
	zoneList := make([]*v3.OrderedAvailabilityZoneList, 0)
	if v, ok := d.GetOk("ordered_availability_zone_list"); ok {
		orders := v.([]interface{})
		for _, order := range orders {
			azcl := &v3.OrderedAvailabilityZoneList{}
			v1 := order.(map[string]interface{})
			if v1, ok1 := v1["cluster_uuid"]; ok1 && v1.(string) != "" {
				azcl.ClusterUUID = v1.(string)
			}
			if v1, ok1 := v1["availability_zone_url"]; ok1 && v1.(string) != "" {
				azcl.AvailabilityZoneURL = v1.(string)
			}
			zoneList = append(zoneList, azcl)
		}
	}
	return zoneList
}

func flattenCategoriesFilter(categoryFilter *v3.CategoryFilter) []interface{} {
	categories := make(map[string]interface{})
	categories2 := make([]interface{}, 0)
	if categoryFilter != nil {
		if categoryFilter.KindList != nil {
			fkl := categoryFilter.KindList
			fkList := make([]string, len(fkl))
			for i, f := range fkl {
				fkList[i] = utils.StringValue(f)
			}
			categories["kind_list"] = fkList
		}
	}
	categories["type"] = utils.StringValue(categoryFilter.Type)
	categories["params"] = expandFilterParams(categoryFilter.Params)
	categories2 = append(categories2, categories)
	return categories2
}

func flattenAvailabilityZoneConnectivityList(azcl []*v3.AvailabilityZoneConnectivityList) []map[string]interface{} {
	availibilityList := make([]map[string]interface{}, 0)
	for _, v := range azcl {
		availability := make(map[string]interface{})

		availability["destination_availability_zone_index"] = utils.Int64Value(v.DestinationAvailabilityZoneIndex)
		availability["source_availability_zone_index"] = utils.Int64Value(v.SourceAvailabilityZoneIndex)
		availability["snapshot_schedule_list"] = flattenSnapshotScheduleList(v.SnapshotScheduleList)

		availibilityList = append(availibilityList, availability)
	}
	return availibilityList
}

func flattenSnapshotScheduleList(snaps []*v3.SnapshotScheduleList) []map[string]interface{} {
	snapshots := make([]map[string]interface{}, 0)
	if snaps != nil {
		snap := make(map[string]interface{})
		for _, v2 := range snaps {
			snap["recovery_point_objective_secs"] = utils.Int64Value(v2.RecoveryPointObjectiveSecs)
			snap["auto_suspend_timeout_secs"] = utils.Int64Value(v2.AutoSuspendTimeoutSecs)
			snap["snapshot_type"] = v2.SnapshotType
			snap["local_snapshot_retention_policy"] = flattenRetentionPolicy(v2.LocalSnapshotRetentionPolicy)
			snap["remote_snapshot_retention_policy"] = flattenRetentionPolicy(v2.RemoteSnapshotRetentionPolicy)
			snapshots = append(snapshots, snap)
		}
	}
	return snapshots
}

func flattenRetentionPolicy(policy *v3.SnapshotRetentionPolicy) []interface{} {
	policies := make([]interface{}, 0)
	rollup := make(map[string]interface{})
	if policy != nil {
		rollup["num_snapshots"] = utils.Int64Value(policy.NumSnapshots)
		if policy.RollupRetentionPolicy != nil {
			rollup["rollup_retention_policy_multiple"] = utils.Int64Value(policy.RollupRetentionPolicy.Multiple)
			rollup["rollup_retention_policy_snapshot_interval_type"] = policy.RollupRetentionPolicy.SnapshotIntervalType
		}
		policies = append(policies, rollup)
	}
	return policies
}

func flattenOrderAvailibilityList(zoneList []*v3.OrderedAvailabilityZoneList) []map[string]interface{} {
	orders := make([]map[string]interface{}, 0)
	for _, v := range zoneList {
		order := make(map[string]interface{})
		order["cluster_uuid"] = v.ClusterUUID
		order["availability_zone_url"] = v.AvailabilityZoneURL
		orders = append(orders, order)
	}

	return orders
}
