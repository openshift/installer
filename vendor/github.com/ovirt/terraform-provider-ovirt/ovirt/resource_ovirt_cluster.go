// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtClusterCreate,
		Read:   resourceOvirtClusterRead,
		Update: resourceOvirtClusterUpdate,
		Delete: resourceOvirtClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Name of the cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Default:     "",
				Description: "Description of the cluster",
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter ID where the cluster resides",
			},
			"management_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Management network ID of the cluster to access cluster hosts",
			},
			"memory_policy_over_commit_percent": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Description: "This indicates the maximum value of memory over committing",
			},
			"ballooning": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Memory balloon is used to re-distribute / reclaim the host memory based on VM needs in a dynamic way",
			},
			"gluster": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "hosts in this cluster will be used as Gluster Storage server nodes",
			},
			"threads_as_cores": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "the exposed host threads would be treated as cores which can be utilized by virtual machines",
			},
			// CPU attributes
			"cpu_arch": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					string(ovirtsdk4.ARCHITECTURE_PPC64),
					string(ovirtsdk4.ARCHITECTURE_S390X),
					string(ovirtsdk4.ARCHITECTURE_UNDEFINED),
					string(ovirtsdk4.ARCHITECTURE_X86_64),
				}, false),
				Description: "CPU architecture of the cluster",
			},
			"cpu_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "CPU codename",
			},
			"compatibility_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The compatibility version of the cluster",
			},
		},
	}
}

func resourceOvirtClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	clusterBuilder := ovirtsdk4.NewClusterBuilder()

	clusterBuilder.Name(d.Get("name").(string)).
		DataCenter(
			ovirtsdk4.NewDataCenterBuilder().
				Id(d.Get("datacenter_id").(string)).
				MustBuild())

	if network, ok := d.GetOk("management_network_id"); ok {
		clusterBuilder.ManagementNetwork(
			ovirtsdk4.NewNetworkBuilder().
				Id(network.(string)).
				MustBuild())
	}

	// Extract memory attributes
	if percent, ok := d.GetOk("memory_policy_over_commit_percent"); ok {
		clusterBuilder.MemoryPolicy(
			ovirtsdk4.NewMemoryPolicyBuilder().
				OverCommit(
					ovirtsdk4.NewMemoryOverCommitBuilder().
						Percent(int64(percent.(int))).
						MustBuild()).
				MustBuild())
	}
	if be, ok := d.GetOkExists("ballooning"); ok {
		clusterBuilder.BallooningEnabled(be.(bool))
	}

	// Extract CPU attributes
	clusterBuilder.Cpu(expandClusterCPU(d))
	// Extract description
	clusterBuilder.Description(d.Get("description").(string))
	// Extract gluster setting
	if gs, ok := d.GetOkExists("gluster"); ok {
		clusterBuilder.GlusterService(gs.(bool))
	}
	// Extract threads_as_cores setting
	if tac, ok := d.GetOkExists("threads_as_cores"); ok {
		clusterBuilder.ThreadsAsCores(tac.(bool))
	}
	// Extract compatibility version
	version, err := expandClusterCompatibilityVersion(d)
	if err != nil {
		return err
	}
	clusterBuilder.Version(version)

	addResp, err := conn.SystemService().ClustersService().
		Add().
		Cluster(clusterBuilder.MustBuild()).
		Send()
	if err != nil {
		log.Printf("[DEBUG] Error adding new cluster: %s", err)
		return err
	}
	d.SetId(addResp.MustCluster().MustId())
	return resourceOvirtClusterRead(d, meta)
}

func resourceOvirtClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	getResp, err := conn.SystemService().
		ClustersService().
		ClusterService(d.Id()).
		Get().
		Send()
	if err != nil {
		if _, ok := err.(*ovirtsdk4.NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	cluster := getResp.MustCluster()
	d.Set("name", cluster.MustName())
	d.Set("datacenter_id", cluster.MustDataCenter().MustId())

	if nL, nLOK := cluster.Networks(); nLOK {
		for _, n := range nL.Slice() {
			for _, u := range n.MustUsages() {
				if u == ovirtsdk4.NETWORKUSAGE_MANAGEMENT {
					d.Set("management_network_id", n.MustId())
					goto AfterNetwork
				}
			}
		}
	}
AfterNetwork:

	if desc, ok := cluster.Description(); ok {
		d.Set("description", desc)
	}
	if mp, mpOK := cluster.MemoryPolicy(); mpOK {
		if oc, ocOK := mp.OverCommit(); ocOK {
			if p, pOK := oc.Percent(); pOK {
				d.Set("memory_policy_over_commit_percent", p)
			}
		}
	}
	if be, ok := cluster.BallooningEnabled(); ok {
		d.Set("ballooning", be)
	}
	if gs, ok := cluster.GlusterService(); ok {
		d.Set("gluster", gs)
	}
	if tac, ok := cluster.ThreadsAsCores(); ok {
		d.Set("threads_as_cores", tac)
	}
	if cpu, cpuOK := cluster.Cpu(); cpuOK {
		if arch, archOK := cpu.Architecture(); archOK {
			d.Set("cpu_arch", arch)
		}
		if tp, tpOK := cpu.Type(); tpOK {
			d.Set("cpu_type", tp)
		}
	}
	if ver, ok := cluster.Version(); ok {
		d.Set("compatibility_version",
			fmt.Sprintf("%d.%d", ver.MustMajor(), ver.MustMinor()))
	}

	return nil
}

func resourceOvirtClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	conn.SystemService()
	paramCluster := ovirtsdk4.NewClusterBuilder()
	attributeUpdate := false

	if d.HasChange("name") {
		paramCluster.Name(d.Get("name").(string))
		attributeUpdate = true
	}

	if d.HasChange("datacenter_id") {
		paramCluster.DataCenter(
			ovirtsdk4.NewDataCenterBuilder().
				Id(d.Get("datacenter_id").(string)).
				MustBuild())
		attributeUpdate = true
	}

	if d.HasChange("management_network_id") {
		if network, ok := d.GetOk("management_network_id"); ok {
			paramCluster.ManagementNetwork(
				ovirtsdk4.NewNetworkBuilder().
					Id(network.(string)).
					MustBuild())
		}
		attributeUpdate = true
	}

	if d.HasChange("description") {
		// `description` has default value ""
		paramCluster.Description(d.Get("description").(string))
		attributeUpdate = true
	}

	if d.HasChange("memory_policy_over_commit_percent") {
		if percent, ok := d.GetOk("memory_policy_over_commit_percent"); ok {
			paramCluster.MemoryPolicy(
				ovirtsdk4.NewMemoryPolicyBuilder().
					OverCommit(
						ovirtsdk4.NewMemoryOverCommitBuilder().
							Percent(int64(percent.(int))).
							MustBuild()).
					MustBuild())
		}
		attributeUpdate = true
	}

	if d.HasChange("ballooning") {
		if be, ok := d.GetOkExists("ballooning"); ok {
			paramCluster.BallooningEnabled(be.(bool))
		}
		attributeUpdate = true
	}

	if d.HasChange("gluster") {
		if gs, ok := d.GetOkExists("gluster"); ok {
			paramCluster.GlusterService(gs.(bool))
		}
		attributeUpdate = true
	}

	if d.HasChange("threads_as_cores") {
		if tac, ok := d.GetOkExists("threads_as_cores"); ok {
			paramCluster.ThreadsAsCores(tac.(bool))
		}
		attributeUpdate = true
	}

	if d.HasChange("compatibility_version") {
		version, err := expandClusterCompatibilityVersion(d)
		if err != nil {
			return err
		}
		paramCluster.Version(version)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := conn.SystemService().
			ClustersService().
			ClusterService(d.Id()).
			Update().
			Cluster(paramCluster.MustBuild()).
			Send()
		if err != nil {
			log.Printf("[DEBUG] Error updating Cluster (%s): %s", d.Id(), err)
			return err
		}
	}

	return resourceOvirtClusterRead(d, meta)
}

func resourceOvirtClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	_, err := conn.SystemService().ClustersService().
		ClusterService(d.Id()).
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

func expandClusterCPU(d *schema.ResourceData) *ovirtsdk4.Cpu {
	cpuBuilder := ovirtsdk4.NewCpuBuilder().
		Architecture(ovirtsdk4.Architecture(d.Get("cpu_arch").(string))).
		Type(d.Get("cpu_type").(string))
	return cpuBuilder.MustBuild()
}

func expandClusterCompatibilityVersion(d *schema.ResourceData) (*ovirtsdk4.Version, error) {
	if version, ok := d.GetOk("compatibility_version"); ok {
		major, minor, _ := extractSemanticVerion(version.(string))
		majorInt, err1 := strconv.ParseInt(major, 10, 64)
		minorInt, err2 := strconv.ParseInt(minor, 10, 64)
		if err1 != nil || err2 != nil || majorInt < 0 || minorInt < 0 {
			return nil, fmt.Errorf("Invalid version format")
		}
		return ovirtsdk4.NewVersionBuilder().
			Major(majorInt).
			Minor(minorInt).
			MustBuild(), nil
	}
	return nil, nil
}
