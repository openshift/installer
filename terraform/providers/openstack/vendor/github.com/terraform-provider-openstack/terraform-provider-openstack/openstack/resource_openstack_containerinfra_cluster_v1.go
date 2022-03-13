package openstack

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
)

func resourceContainerInfraClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerInfraClusterV1Create,
		ReadContext:   resourceContainerInfraClusterV1Read,
		UpdateContext: resourceContainerInfraClusterV1Update,
		DeleteContext: resourceContainerInfraClusterV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"api_address": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"coe_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"cluster_template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_MAGNUM_CLUSTER_TEMPLATE", nil),
			},

			"container_version": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"create_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"discovery_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"docker_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"master_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"keypair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"merge_labels": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"master_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"master_addresses": {
				Type:     schema.TypeList,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"node_addresses": {
				Type:     schema.TypeList,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"stack_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Computed: true,
			},

			"fixed_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"fixed_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"floating_ip_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"kubeconfig": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContainerInfraClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	// Get and check labels map.
	rawLabels := d.Get("labels").(map[string]interface{})
	labels, err := expandContainerInfraV1LabelsMap(rawLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the flavors to use.
	// First check if it was set in the config.
	// If not, try using the appropriate environment variable.
	flavor, err := containerInfraClusterV1Flavor(d)
	if err != nil {
		return diag.Errorf("Unable to determine openstack_containerinfra_cluster_v1 flavor")
	}

	masterFlavor, err := containerInfraClusterV1MasterFlavor(d)
	if err != nil {
		return diag.Errorf("Unable to determine openstack_containerinfra_cluster_v1 master_flavor")
	}

	// Get boolean parameters that will be passed by reference.
	floatingIPEnabled := d.Get("floating_ip_enabled").(bool)

	createOpts := clusters.CreateOpts{
		ClusterTemplateID: d.Get("cluster_template_id").(string),
		DiscoveryURL:      d.Get("discovery_url").(string),
		FlavorID:          flavor,
		Keypair:           d.Get("keypair").(string),
		Labels:            labels,
		MasterFlavorID:    masterFlavor,
		Name:              d.Get("name").(string),
		FixedNetwork:      d.Get("fixed_network").(string),
		FixedSubnet:       d.Get("fixed_subnet").(string),
		FloatingIPEnabled: &floatingIPEnabled,
	}

	// Set int parameters that will be passed by reference.
	createTimeout := d.Get("create_timeout").(int)
	if createTimeout > 0 {
		createOpts.CreateTimeout = &createTimeout
	}

	dockerVolumeSize := d.Get("docker_volume_size").(int)
	if dockerVolumeSize > 0 {
		createOpts.DockerVolumeSize = &dockerVolumeSize
	}

	masterCount := d.Get("master_count").(int)
	if masterCount > 0 {
		createOpts.MasterCount = &masterCount
	}

	nodeCount := d.Get("node_count").(int)
	if nodeCount > 0 {
		createOpts.NodeCount = &nodeCount
	}

	mergeLabels := d.Get("merge_labels").(bool)
	if mergeLabels {
		createOpts.MergeLabels = &mergeLabels
	}

	s, err := clusters.Create(containerInfraClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_containerinfra_cluster_v1: %s", err)
	}

	// Store the Cluster ID.
	d.SetId(s)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Target:       []string{"CREATE_COMPLETE"},
		Refresh:      containerInfraClusterV1StateRefreshFunc(containerInfraClient, s),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_cluster_v1 %s to become ready: %s", s, err)
	}

	log.Printf("[DEBUG] Created openstack_containerinfra_cluster_v1 %s", s)

	return resourceContainerInfraClusterV1Read(ctx, d, meta)
}

func resourceContainerInfraClusterV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	s, err := clusters.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_containerinfra_cluster_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_containerinfra_cluster_v1 %s: %#v", d.Id(), s)

	if err := d.Set("labels", s.Labels); err != nil {
		return diag.Errorf("Unable to set openstack_containerinfra_cluster_v1 labels: %s", err)
	}

	d.Set("name", s.Name)
	d.Set("api_address", s.APIAddress)
	d.Set("coe_version", s.COEVersion)
	d.Set("cluster_template_id", s.ClusterTemplateID)
	d.Set("container_version", s.ContainerVersion)
	d.Set("create_timeout", s.CreateTimeout)
	d.Set("discovery_url", s.DiscoveryURL)
	d.Set("docker_volume_size", s.DockerVolumeSize)
	d.Set("flavor", s.FlavorID)
	d.Set("master_flavor", s.MasterFlavorID)
	d.Set("keypair", s.KeyPair)
	d.Set("master_count", s.MasterCount)
	d.Set("node_count", s.NodeCount)
	d.Set("master_addresses", s.MasterAddresses)
	d.Set("node_addresses", s.NodeAddresses)
	d.Set("stack_id", s.StackID)
	d.Set("fixed_network", s.FixedNetwork)
	d.Set("fixed_subnet", s.FixedSubnet)
	d.Set("floating_ip_enabled", s.FloatingIPEnabled)

	kubeconfig, err := flattenContainerInfraV1Kubeconfig(d, containerInfraClient)
	if err != nil {
		return diag.Errorf("Error building kubeconfig for openstack_containerinfra_cluster_v1 %s: %s", d.Id(), err)
	}
	d.Set("kubeconfig", kubeconfig)

	if err := d.Set("created_at", s.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_cluster_v1 created_at: %s", err)
	}
	if err := d.Set("updated_at", s.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_containerinfra_cluster_v1 updated_at: %s", err)
	}

	return nil
}

func resourceContainerInfraClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	updateOpts := []clusters.UpdateOptsBuilder{}

	if d.HasChange("node_count") {
		v := d.Get("node_count").(int)
		nodeCount := strconv.Itoa(v)
		updateOpts = append(updateOpts, clusters.UpdateOpts{
			Op:    clusters.ReplaceOp,
			Path:  strings.Join([]string{"/", "node_count"}, ""),
			Value: nodeCount,
		})
	}

	if len(updateOpts) > 0 {
		log.Printf(
			"[DEBUG] Updating openstack_containerinfra_cluster_v1 %s with options: %#v", d.Id(), updateOpts)

		_, err = clusters.Update(containerInfraClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_containerinfra_cluster_v1 %s: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"UPDATE_IN_PROGRESS"},
			Target:       []string{"UPDATE_COMPLETE"},
			Refresh:      containerInfraClusterV1StateRefreshFunc(containerInfraClient, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Minute,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_containerinfra_cluster_v1 %s to become updated: %s", d.Id(), err)
		}
	}
	return resourceContainerInfraClusterV1Read(ctx, d, meta)
}

func resourceContainerInfraClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerInfraV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack container infra client: %s", err)
	}

	if err := clusters.Delete(containerInfraClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_containerinfra_cluster_v1"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Target:       []string{"DELETE_COMPLETE"},
		Refresh:      containerInfraClusterV1StateRefreshFunc(containerInfraClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_containerinfra_cluster_v1 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
