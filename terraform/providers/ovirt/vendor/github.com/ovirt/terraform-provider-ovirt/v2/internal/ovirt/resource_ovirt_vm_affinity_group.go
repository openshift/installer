package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

var vmAffinityGroupSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"cluster_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "ID of the oVirt cluster the VM is in.",
		ValidateDiagFunc: validateUUID,
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "ID of the oVirt VM to assign to the affinity group.",
		ValidateDiagFunc: validateUUID,
	},
	"affinity_group_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "ID of the oVirt affinity group to assign the VM to.",
		ValidateDiagFunc: validateNonEmpty,
	},
}

func (p *provider) vmAffinityGroupResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmAffinityGroupCreate,
		ReadContext:   p.vmAffinityGroupRead,
		DeleteContext: p.vmAffinityGroupDelete,
		Schema:        vmAffinityGroupSchema,
		Description:   "The ovirt_vm_affinity_group resource assigns VMs to affinity groups in oVirt.",
	}
}

func (p *provider) vmAffinityGroupCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	clusterID := data.Get("cluster_id").(string)
	VMID := data.Get("vm_id").(string)
	affinityGroupID := data.Get("affinity_group_id").(string)

	err := client.AddVMToAffinityGroup(ovirtclient.ClusterID(clusterID), ovirtclient.VMID(VMID), ovirtclient.AffinityGroupID(affinityGroupID))
	if err != nil {
		return diag.Diagnostics{
			errorToDiag(fmt.Sprintf("add VM %s in cluster %s to affinity group %s", VMID, clusterID, affinityGroupID), err),
		}
	}

	data.SetId(VMID)

	return nil
}

func (p *provider) vmAffinityGroupRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	clusterID := data.Get("cluster_id").(string)
	VMID := data.Get("vm_id").(string)
	affinityGroupID := data.Get("affinity_group_id").(string)

	affinityGroup, err := client.GetAffinityGroup(ovirtclient.ClusterID(clusterID), ovirtclient.AffinityGroupID(affinityGroupID))
	if err != nil {
		return errorToDiags(fmt.Sprintf("get affinity group '%s' in cluster '%s'", affinityGroupID, clusterID), err)
	}

	for _, affinityGroupVMID := range affinityGroup.VMIDs() {
		if string(affinityGroupVMID) == VMID {
			data.SetId(VMID)
			return nil
		}
	}

	data.SetId("")

	return nil
}

func (p *provider) vmAffinityGroupDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	clusterID := data.Get("cluster_id").(string)
	VMID := data.Get("vm_id").(string)
	affinityGroupID := data.Get("affinity_group_id").(string)

	err := client.RemoveVMFromAffinityGroup(ovirtclient.ClusterID(clusterID), ovirtclient.VMID(VMID), ovirtclient.AffinityGroupID(affinityGroupID))
	if err != nil {
		if !isNotFound(err) {
			return errorToDiags(fmt.Sprintf("remove VM '%s' in cluster '%s' from affinity group '%s'", VMID, clusterID, affinityGroupID), err)
		}
	}

	data.SetId("")

	return nil
}
