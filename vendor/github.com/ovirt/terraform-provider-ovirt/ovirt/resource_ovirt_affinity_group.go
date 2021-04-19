package ovirt

import (
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtAffinityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtAffinityGroupCreate,
		Read:   resourceOvirtAffinityGroupRead,
		Update: resourceOvirtAffinityGroupUpdate,
		Delete: resourceOvirtAffinityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Name of the affinity group",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Description of the affinity group",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID where the affinity group is",
			},
			"priority": {
				Type:        schema.TypeFloat,
				Optional:    true,
				ForceNew:    false,
				Description: "Priority of the affinity group",
			},
			"vm_positive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Positive or negative affinity",
			},
			"vm_enforcing": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     false,
				Description: "Is the policy being enforced",
			},
			"vm_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List of VMs in the affinity group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"host_positive": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Positive or negative affinity",
			},
			"host_enforcing": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     false,
				Description: "Is the policy being enforced",
			},
			"host_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List of Hosts in the affinity group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceOvirtAffinityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	agBuilder := ovirtsdk4.NewAffinityGroupBuilder()

	agBuilder.Name(d.Get("name").(string))

	if desc, ok := d.GetOk("description"); ok {
		agBuilder.Description(desc.(string))
	}

	if priority, ok := d.GetOk("priority"); ok {
		agBuilder.Priority(priority.(float64))
	}

	vmRuleBuilder := ovirtsdk4.NewAffinityRuleBuilder()
	vmRuleBuilder.Enabled(false)
	if vmPositive, ok := d.GetOkExists("vm_positive"); ok {
		vmRuleBuilder.Enabled(true)
		vmRuleBuilder.Positive(vmPositive.(bool))
	}
	if vmEnforcing, ok := d.GetOk("vm_enforcing"); ok {
		vmRuleBuilder.Enabled(true)
		vmRuleBuilder.Enforcing(vmEnforcing.(bool))
	}

	agBuilder.VmsRule(vmRuleBuilder.MustBuild())

	hostRuleBuilder := ovirtsdk4.NewAffinityRuleBuilder()
	hostRuleBuilder.Enabled(false)
	if hostPositive, ok := d.GetOkExists("host_positive"); ok {
		hostRuleBuilder.Enabled(true)
		hostRuleBuilder.Positive(hostPositive.(bool))
	}
	if hostEnforcing, ok := d.GetOk("host_enforcing"); ok {
		hostRuleBuilder.Enabled(true)
		hostRuleBuilder.Enforcing(hostEnforcing.(bool))
	}

	agBuilder.HostsRule(hostRuleBuilder.MustBuild())

	log.Printf("Creating %#v", agBuilder.MustBuild())
	addResp, err := conn.SystemService().
		ClustersService().
		ClusterService(d.Get("cluster_id").(string)).
		AffinityGroupsService().
		Add().
		Group(agBuilder.MustBuild()).
		Send()

	if err != nil {
		log.Printf("Failed to create Affinity Group")
		return err
	}

	log.Printf("Successfully created %#v", agBuilder.MustBuild().MustName())
	d.SetId(addResp.MustGroup().MustId())

	// Add VMs to affinity group
	if vmList, ok := d.GetOk("vm_list"); ok {
		vmsService := conn.SystemService().
			ClustersService().
			ClusterService(d.Get("cluster_id").(string)).
			AffinityGroupsService().
			GroupService(addResp.MustGroup().MustId())

		if err := updateVmList(vmsService, vmList.([]interface{})); err != nil {
			return err
		}
	}

	// Add hosts to affinity group
	if hostList, ok := d.GetOk("host_list"); ok {
		groupService := conn.SystemService().
			ClustersService().
			ClusterService(d.Get("cluster_id").(string)).
			AffinityGroupsService().
			GroupService(addResp.MustGroup().MustId())

		if err := updateHostList(groupService, hostList.([]interface{})); err != nil {
			return err
		}
	}

	return resourceOvirtAffinityGroupRead(d, meta)
}

func resourceOvirtAffinityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	agService := conn.SystemService().
		ClustersService().
		ClusterService(d.Get("cluster_id").(string)).
		AffinityGroupsService().
		GroupService(d.Id())

	affinityGroupResp, err := agService.Get().Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	affinityGroup := affinityGroupResp.MustGroup()

	d.Set("name", affinityGroup.MustName())
	if desc, ok := affinityGroup.Description(); ok {
		d.Set("description", desc)
	}
	if priority, ok := affinityGroup.Priority(); ok {
		d.Set("priority", priority)
	}
	d.Set("host_enabled", affinityGroup.MustHostsRule().MustEnabled())
	d.Set("host_enforcing", affinityGroup.MustHostsRule().MustEnforcing())
	d.Set("host_positive", affinityGroup.MustHostsRule().MustPositive())
	d.Set("vm_enabled", affinityGroup.MustVmsRule().MustEnabled())
	d.Set("vm_enforcing", affinityGroup.MustVmsRule().MustEnforcing())
	d.Set("vm_positive", affinityGroup.MustVmsRule().MustPositive())

	hosts := affinityGroup.MustHosts().Slice()
	hostNames := make([]string, len(hosts))
	for i, h := range hosts {
		hostNames[i] = h.MustId()
	}
	sort.Strings(hostNames)
	d.Set("host_list", hostNames)

	vms := affinityGroup.MustVms().Slice()
	vmNames := make([]string, len(vms))
	for i, v := range vms {
		vmNames[i] = v.MustId()
	}
	sort.Strings(vmNames)
	d.Set("vm_list", vmNames)

	return nil
}

func resourceOvirtAffinityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	group := ovirtsdk4.NewAffinityGroupBuilder()
	attributeUpdate := false

	if d.HasChange("name") {
		group.Name(d.Get("name").(string))
		attributeUpdate = true
	}
	if d.HasChange("description") {
		group.Description(d.Get("description").(string))
		attributeUpdate = true
	}
	if d.HasChange("priority") {
		group.Priority(d.Get("priority").(float64))
		attributeUpdate = true
	}
	if d.HasChange("cluster_id") {
		group.Cluster(
			ovirtsdk4.NewClusterBuilder().
				Id(d.Get("cluster_id").(string)).
				MustBuild())
		attributeUpdate = true
	}

	vmRuleBuilder := ovirtsdk4.NewAffinityRuleBuilder()
	vmRuleUpdate := false
	if d.HasChange("vm_positive") {
		vmRuleBuilder.Positive(d.Get("vm_positive").(bool))
		vmRuleUpdate = true
	}
	if d.HasChange("vm_enforcing") {
		vmRuleBuilder.Enforcing(d.Get("vm_enforcing").(bool))
		vmRuleUpdate = true
	}
	if d.HasChange("vm_list") {
		vmRuleBuilder.Enabled(len(d.Get("vm_list").([]interface{})) > 0)
		vmRuleUpdate = true
	}
	if vmRuleUpdate {
		group.VmsRule(vmRuleBuilder.MustBuild())
		attributeUpdate = true
	}

	hostRuleBuilder := ovirtsdk4.NewAffinityRuleBuilder()
	hostRuleUpdate := false
	if d.HasChange("host_positive") {
		hostRuleBuilder.Positive(d.Get("host_positive").(bool))
		hostRuleUpdate = true
	}
	if d.HasChange("host_enforcing") {
		hostRuleBuilder.Enforcing(d.Get("host_enforcing").(bool))
		hostRuleUpdate = true
	}
	if d.HasChange("host_list") {
		hostRuleBuilder.Enabled(len(d.Get("host_list").([]interface{})) > 0)
		hostRuleUpdate = true
	}
	if hostRuleUpdate {
		group.HostsRule(hostRuleBuilder.MustBuild())
		attributeUpdate = true
	}

	if attributeUpdate {
		log.Printf("[DEBUG] Updating %#v", group.MustBuild())
		_, err := conn.SystemService().
			ClustersService().
			ClusterService(d.Get("cluster_id").(string)).
			AffinityGroupsService().
			GroupService(d.Id()).
			Update().
			Group(group.MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error updating affinity group (%s): %s", d.Id(), err)
			return err
		}
	}

	if d.HasChange("vm_list") {
		groupService := conn.SystemService().
			ClustersService().
			ClusterService(d.Get("cluster_id").(string)).
			AffinityGroupsService().
			GroupService(d.Id())

		if err := updateVmList(groupService, d.Get("vm_list").([]interface{})); err != nil {
			return err
		}
	}

	if d.HasChange("host_list") {
		groupService := conn.SystemService().
			ClustersService().
			ClusterService(d.Get("cluster_id").(string)).
			AffinityGroupsService().
			GroupService(d.Id())

		if err := updateHostList(groupService, d.Get("host_list").([]interface{})); err != nil {
			return err
		}
	}

	return resourceOvirtAffinityGroupRead(d, meta)
}

func resourceOvirtAffinityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	_, err := conn.SystemService().
		ClustersService().
		ClusterService(d.Get("cluster_id").(string)).
		AffinityGroupsService().
		GroupService(d.Id()).
		Remove().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			return nil
		}
		return err
	}
	return nil
}

func updateVmList(affinityGroupService *ovirtsdk4.AffinityGroupService, vmList []interface{}) error {
	vms := make([]*ovirtsdk4.Vm, len(vmList))
	for i, h := range vmList {
		vms[i] = ovirtsdk4.NewVmBuilder().Id(h.(string)).MustBuild()
	}

	agBuilder := ovirtsdk4.NewAffinityGroupBuilder()
	var vmSlice = new(ovirtsdk4.VmSlice)
	vmSlice.SetSlice(vms)
	agBuilder.Vms(vmSlice)

	_, err := affinityGroupService.Update().Group(agBuilder.MustBuild()).Send()
	if err != nil {
		return err
	}

	return nil
}

func updateHostList(affinityGroupService *ovirtsdk4.AffinityGroupService, hostList []interface{}) error {
	hosts := make([]*ovirtsdk4.Host, len(hostList))
	for i, h := range hostList {
		hosts[i] = ovirtsdk4.NewHostBuilder().Id(h.(string)).MustBuild()
	}

	agBuilder := ovirtsdk4.NewAffinityGroupBuilder()
	var hostSlice = new(ovirtsdk4.HostSlice)
	hostSlice.SetSlice(hosts)
	agBuilder.Hosts(hostSlice)

	_, err := affinityGroupService.Update().Group(agBuilder.MustBuild()).Send()
	if err != nil {
		return err
	}

	return nil
}
