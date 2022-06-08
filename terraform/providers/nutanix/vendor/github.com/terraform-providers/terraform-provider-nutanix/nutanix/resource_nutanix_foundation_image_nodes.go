package nutanix

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var (
	ImageMinTimeout          = 60 * time.Minute
	AggregatePercentComplete = 100
	ClusterURL               = "https://%s:9440/"
)

func resourceFoundationImageNodes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoundationImageNodesCreate,
		ReadContext:   resourceFoundationImageNodesRead,
		DeleteContext: resourceFoundationImageNodesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ImageMinTimeout),
		},

		Schema: map[string]*schema.Schema{
			"xs_master_label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipmi_netmask": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipmi_gateway": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipmi_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cvm_gateway": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hyperv_external_vnic": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"xen_config_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ucsm_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ucsm_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hypervisor_iso": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hyperv": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"kvm": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"xen": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"esx": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"checksum": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"unc_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hypervisor_netmask": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fc_settings": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fc_metadata": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fc_ip": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"api_key": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"foundation_central": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"xs_master_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"svm_rescue_args": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cvm_netmask": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"xs_master_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_ns": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"backplane_subnet": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cluster_init_successful": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"backplane_netmask": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"backplane_vlan": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cluster_external_ip": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cvm_ntp_servers": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"single_node_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"cluster_members": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cvm_dns_servers": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cluster_init_now": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},
						"hypervisor_ntp_servers": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"hyperv_external_vswitch": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hypervisor_nameserver": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hyperv_sku": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"eos_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"account_name": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"tests": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_syscheck": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"run_ncc": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"blocks": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipmi_netmask": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_gateway": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ipv6_address": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"node_position": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"image_delay": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"ucsm_params": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"native_vlan": {
													Type:     schema.TypeBool,
													Optional: true,
													ForceNew: true,
												},
												"keep_ucsm_settings": {
													Type:     schema.TypeBool,
													Optional: true,
													ForceNew: true,
												},
												"mac_pool": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"vlan_name": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
									"hypervisor_hostname": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"cvm_gb_ram": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"device_hint": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"bond_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"rdma_passthrough": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ucsm_node_serial": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"hypervisor_ip": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"node_serial": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_configure_now": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"image_successful": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"ipv6_interface": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"cvm_num_vcpus": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_mac": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"rdma_mac_addr": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"bond_uplinks": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"current_network_interface": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"hypervisor": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"vswitches": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lacp": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"bond_mode": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"uplinks": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"other_config": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"mtu": {
													Type:     schema.TypeInt,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
									"bond_lacp_rate": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"image_now": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},
									"ucsm_managed_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"current_cvm_vlan_tag": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"cvm_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"exlude_boot_serial": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"mitigate_low_boot_space": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_password": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"ipmi_user": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"block_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"hyperv_product_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"unc_username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"install_script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipmi_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hypervisor_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"unc_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"xs_master_username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"skip_hypervisor": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"hypervisor_gateway": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nos_package": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ucsm_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"layout_egg_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceFoundationImageNodesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// create connection
	conn := meta.(*Client).FoundationClientAPI
	// Prepare request
	request := &foundation.ImageNodesInput{}
	xsmasterlabel, ok := d.GetOk("xs_master_label")
	if ok {
		request.XsMasterLabel = (xsmasterlabel.(string))
	}
	layoutEggUUID, ok := d.GetOk("layout_egg_uuid")
	if ok {
		request.LayoutEggUUID = (layoutEggUUID.(string))
	}
	ipmiPass, ok := d.GetOk("ipmi_password")
	if ok {
		request.IpmiPassword = ipmiPass.(string)
	}
	if ipmiNetmask, ipmiNetOk := d.GetOk("ipmi_netmask"); ipmiNetOk {
		request.IPMINetmask = ipmiNetmask.(string)
	}
	if ipmiGateway, ipmiGateOk := d.GetOk("ipmi_gateway"); ipmiGateOk {
		request.IPMIGateway = ipmiGateway.(string)
	}
	cvmGateway, cvmgok := d.GetOk("cvm_gateway")
	if cvmgok {
		request.CvmGateway = (cvmGateway.(string))
	}

	hypervExternalVnic, hyExNicok := d.GetOk("hyperv_external_vnic")
	if hyExNicok {
		request.HypervExternalVnic = hypervExternalVnic.(string)
	}

	xenConfigType, ok := d.GetOk("xen_config_type")
	if ok {
		request.XenConfigType = (xenConfigType.(string))
	}

	ucsmIP, ok := d.GetOk("ucsm_ip")
	if ok {
		request.UcsmIP = (ucsmIP.(string))
	}

	ucsmPassword, ok := d.GetOk("ucsm_password")
	if ok {
		request.UcsmPassword = (ucsmPassword.(string))
	}

	uncPath, ok := d.GetOk("unc_path")
	if ok {
		request.UncPath = (uncPath.(string))
	}

	hypervisorNetmask, ok := d.GetOk("hypervisor_netmask")
	if ok {
		request.HypervisorNetmask = (hypervisorNetmask.(string))
	}

	xsMasterPassword, ok := d.GetOk("xs_master_password")
	if ok {
		request.XsMasterPassword = (xsMasterPassword.(string))
	}

	cvmNetmask, ok := d.GetOk("cvm_netmask")
	if ok {
		request.CvmNetmask = (cvmNetmask.(string))
	}

	xsMasterIP, ok := d.GetOk("xs_master_ip")
	if ok {
		request.XsMasterIP = (xsMasterIP.(string))
	}

	hypervExternalVswitch, ok := d.GetOk("hyperv_external_vswitch")
	if ok {
		request.HypervExternalVswitch = hypervExternalVswitch.(string)
	}

	hypervisorNameserver, ok := d.GetOk("hypervisor_nameserver")
	if ok {
		request.HypervisorNameserver = (hypervisorNameserver.(string))
	}

	hypervSku, ok := d.GetOk("hyperv_sku")
	if ok {
		request.HypervSku = (hypervSku.(string))
	}

	hypervProductKey, ok := d.GetOk("hyperv_product_key")
	if ok {
		request.HypervProductKey = (hypervProductKey.(string))
	}

	uncUsername, ok := d.GetOk("unc_username")
	if ok {
		request.UncUsername = (uncUsername.(string))
	}

	installScript, ok := d.GetOk("install_script")
	if ok {
		request.InstallScript = (installScript.(string))
	}

	ipmiUser, ok := d.GetOk("ipmi_user")
	if ok {
		request.IpmiUser = (ipmiUser.(string))
	}

	hypervisorPassword, ok := d.GetOk("hypervisor_password")
	if ok {
		request.HypervisorPassword = (hypervisorPassword.(string))
	}

	uncPassword, ok := d.GetOk("unc_password")
	if ok {
		request.UncPassword = (uncPassword.(string))
	}

	xsMasterUsername, ok := d.GetOk("xs_master_username")
	if ok {
		request.XsMasterUsername = (xsMasterUsername.(string))
	}

	skipHypervisor, ok := d.GetOk("skip_hypervisor")
	if ok {
		request.SkipHypervisor = utils.BoolPtr(skipHypervisor.(bool))
	}

	hypervisorGateway, ok := d.GetOk("hypervisor_gateway")
	if ok {
		request.HypervisorGateway = (hypervisorGateway.(string))
	}

	nosPackage, ok := d.GetOk("nos_package")
	if ok {
		request.NosPackage = (nosPackage.(string))
	}

	ucsmUser, ok := d.GetOk("ucsm_user")
	if ok {
		request.UcsmUser = (ucsmUser.(string))
	}

	fcSettings, err := expandFcSetting(d)
	if err == nil {
		request.FcSettings = fcSettings
	}

	eosMeta, err := expandEosMetadata(d)
	if err == nil {
		request.EosMetadata = eosMeta
	}

	tests, err := expandTests(d)
	if err == nil {
		request.Tests = tests
	}

	hypervisorIso, ok := d.GetOk("hypervisor_iso")
	if ok && len(hypervisorIso.([]interface{})) > 0 {
		request.HypervisorIso = expandHypervisorIso(hypervisorIso.([]interface{}))
	}

	cluster, err := expandCluster(d)
	if err == nil {
		request.Clusters = cluster
	}

	blocks, err := expandBlocks(d)
	if err == nil {
		request.Blocks = blocks
	}

	// call the client here
	resp, err := conn.NodeImaging.ImageNodes(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	// if node images gets errors out initially itself
	if resp.Error != nil {
		return diag.Errorf("Node imaging process failed due to error: %s", resp.Error.Message)
	}

	//poll for progress
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: foundationImageRefresh(ctx, conn, resp.SessionID),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   1 * time.Minute,
	}
	info, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for image (%s) to be ready: %v", resp.SessionID, err)
	}
	if progress, ok := info.(*foundation.ImageNodesProgressResponse); ok {
		if utils.Float64Value(progress.AggregatePercentComplete) < float64(AggregatePercentComplete) {
			return collectIndividualErrorDiagnostics(progress)
		}
	}

	d.SetId(resp.SessionID)
	d.Set("session_id", resp.SessionID)

	// set cluster urls in state file
	clusterURLs := make([]map[string]interface{}, len(request.Clusters))
	for k, v := range request.Clusters {
		c := map[string]interface{}{
			"cluster_url":  fmt.Sprintf(ClusterURL, v.ClusterMembers[0]),
			"cluster_name": v.ClusterName,
		}
		clusterURLs[k] = c
	}
	d.Set("cluster_urls", clusterURLs)
	return nil
}

func resourceFoundationImageNodesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceFoundationImageNodesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandTests(d *schema.ResourceData) (*foundation.Tests, error) {
	if test, ok := d.GetOk("tests"); ok && len(test.([]interface{})) > 0 {
		set := (test.([]interface{})[0]).(map[string]interface{})
		tests := &foundation.Tests{}
		if runsync, ok := set["run_ncc"]; ok {
			tests.RunSyscheck = utils.BoolPtr(runsync.(bool))
		}
		if runncc, ok := set["run_syscheck"]; ok {
			tests.RunNcc = utils.BoolPtr(runncc.(bool))
		}
		return tests, nil
	}
	return nil, nil
}

func expandEosMetadata(d *schema.ResourceData) (*foundation.EosMetadata, error) {
	if eos, ok := d.GetOk("eos_metadata"); ok && len(eos.([]interface{})) > 0 {
		eosmeta := (eos.([]interface{}))[0].(map[string]interface{})
		eosMeta := &foundation.EosMetadata{}
		if config, ok := eosmeta["config_id"]; ok {
			eosMeta.ConfigID = (config.(string))
		}

		if acname, ok := eosmeta["account_name"]; ok {
			ac := acname.([]interface{})
			accNames := make([]string, len(ac))
			for a := range ac {
				accNames[a] = ac[a].(string)
			}
			eosMeta.AccountName = accNames
		}

		if email, ok := eosmeta["email"]; ok {
			eosMeta.Email = (email.(string))
		}
		return eosMeta, nil
	}
	return nil, nil
}
func expandFcSetting(d *schema.ResourceData) (*foundation.FcSettings, error) {
	if fcset, ok := d.GetOk("fc_settings"); ok && len(fcset.([]interface{})) > 0 {
		fc := &foundation.FcSettings{}
		set := (fcset.([]interface{}))[0].(map[string]interface{})

		if val, ok2 := set["foundation_central"]; ok2 {
			fc.FoundationCentral = utils.BoolPtr(val.(bool))
		}
		// expand fc metadata
		if val, ok2 := set["fc_metadata"]; ok2 && len(val.([]interface{})) > 0 {
			fcmeta := (val.([]interface{}))[0].(map[string]interface{})
			fcMetadata := &foundation.FcMetadata{}
			if val, ok := fcmeta["fc_ip"]; ok {
				fcMetadata.FcIP = (val.(string))
			}
			if val, ok := fcmeta["api_key"]; ok {
				fcMetadata.APIKey = val.(string)
			}
			fc.FcMetadata = fcMetadata
		} else {
			return fc, fmt.Errorf("fc_metadata is required field for fc_settings block")
		}
		return fc, nil
	}
	return nil, nil
}

func expandHypervisorIso(pr []interface{}) foundation.HypervisorIso {
	iso := foundation.HypervisorIso{}

	hypervisors := pr[0].(map[string]interface{})
	if hyperv, ok := hypervisors["hyperv"]; ok && len(hyperv.([]interface{})) > 0 {
		iso.Hyperv = expandHypervisor(hyperv.([]interface{}))
	}
	if kvm, ok := hypervisors["kvm"]; ok && len(kvm.([]interface{})) > 0 {
		iso.Kvm = expandHypervisor(kvm.([]interface{}))
	}
	if xen, ok := hypervisors["xen"]; ok && len(xen.([]interface{})) > 0 {
		iso.Xen = expandHypervisor(xen.([]interface{}))
	}
	if esx, ok := hypervisors["esx"]; ok && len(esx.([]interface{})) > 0 {
		iso.Esx = expandHypervisor(esx.([]interface{}))
	}
	return iso
}

func expandHypervisor(pr []interface{}) *foundation.Hypervisor {
	hyp := &foundation.Hypervisor{}

	hypervisors := pr[0].(map[string]interface{})
	if checksum, ok := hypervisors["checksum"]; ok {
		hyp.Checksum = checksum.(string)
	}
	if filename, ok := hypervisors["filename"]; ok {
		hyp.Filename = filename.(string)
	}
	return hyp
}

func expandVswitches(pr interface{}) []*foundation.Vswitches {
	vswit := pr.([]interface{})
	outbound := make([]*foundation.Vswitches, len(vswit))
	for k, vs := range vswit {
		vs := vs.(map[string]interface{})
		vst := &foundation.Vswitches{}
		if lacp, ok := vs["lacp"]; ok {
			vst.Lacp = lacp.(string)
		}
		if bondmode, ok := vs["bond_mode"]; ok {
			vst.BondMode = bondmode.(string)
		}
		if mtu, ok := vs["mtu"]; ok {
			vst.Mtu = utils.Int64Ptr(int64(mtu.(int)))
		}
		if name, ok := vs["name"]; ok {
			vst.Name = name.(string)
		}

		if otherconf, ok := vs["other_config"]; ok {
			other := otherconf.([]interface{})

			otherConfigArray := make([]string, len(other))
			for o := range other {
				otherConfigArray[o] = other[o].(string)
			}
			vst.OtherConfig = otherConfigArray
		}
		if uplinks, ok := vs["uplinks"]; ok {
			ups := uplinks.([]interface{})
			uplinksArray := make([]string, len(ups))
			for o := range ups {
				uplinksArray[o] = ups[o].(string)
			}
			vst.Uplinks = uplinksArray
		}
		outbound[k] = vst
	}
	return outbound
}

func expandUcsmParams(pr interface{}) *foundation.UcsmParams {
	ucsm := pr.([]interface{})
	if len(ucsm) == 0 {
		return nil
	}
	UcsmParam := &foundation.UcsmParams{}

	for _, k := range ucsm {
		set := k.(map[string]interface{})

		if nativevlan, ok := set["native_vlan"]; ok {
			UcsmParam.NativeVlan = utils.BoolPtr(nativevlan.(bool))
		}
		if KeepUcsmSettings, ok := set["keep_ucsm_settings"]; ok {
			UcsmParam.KeepUcsmSettings = utils.BoolPtr(KeepUcsmSettings.(bool))
		}
		if macPool, ok := set["mac_pool"]; ok {
			UcsmParam.MacPool = macPool.(string)
		}
		if VlanName, ok := set["vlan_name"]; ok {
			UcsmParam.VlanName = VlanName.(string)
		}
	}
	return UcsmParam
}

func expandCluster(d *schema.ResourceData) ([]*foundation.Clusters, error) {
	clstr := make([]*foundation.Clusters, 0)
	if v, ok := d.GetOk("clusters"); ok {
		n := v.([]interface{})
		if len(n) > 0 {
			cls := make([]*foundation.Clusters, 0)

			for _, nc := range n {
				clst := nc.(map[string]interface{})

				clusterList := &foundation.Clusters{}
				if enablens, ok := clst["enable_ns"]; ok {
					clusterList.EnableNs = utils.BoolPtr(enablens.(bool))
				}
				if backplanesn, ok := clst["backplane_subnet"]; ok {
					clusterList.BackplaneSubnet = backplanesn.(string)
				}
				if clstinit, ok := clst["cluster_init_successful"]; ok {
					clusterList.ClusterInitSuccessful = utils.BoolPtr(clstinit.(bool))
				}
				if backplanenm, ok := clst["backplane_netmask"]; ok {
					clusterList.BackplaneNetmask = (backplanenm.(string))
				}
				if rf, ok := clst["redundancy_factor"]; ok {
					clusterList.RedundancyFactor = utils.Int64Ptr(int64(rf.(int)))
				}
				if backplanevlan, ok := clst["backplane_vlan"]; ok {
					clusterList.BackplaneVlan = (backplanevlan.(string))
				}
				if clustername, ok := clst["cluster_name"]; ok {
					clusterList.ClusterName = clustername.(string)
				}
				if clusterext, clusterextok := clst["cluster_external_ip"]; clusterextok && clusterext != "" {
					clusterList.ClusterExternalIP = utils.StringPtr(clusterext.(string))
				}
				if cvmntps, ok := clst["cvm_ntp_servers"]; ok {
					clusterList.CvmNtpServers = (cvmntps.(string))
				}
				if sncluster, ok := clst["single_node_cluster"]; ok {
					clusterList.SingleNodeCluster = utils.BoolPtr(sncluster.(bool))
				}
				if cvmdns, ok := clst["cvm_dns_servers"]; ok {
					clusterList.CvmDNSServers = (cvmdns.(string))
				}
				if clusterinitnow, ok := clst["cluster_init_now"]; ok {
					clusterList.ClusterInitNow = utils.BoolPtr(clusterinitnow.(bool))
				}
				if hypervntps, ok := clst["hypervisor_ntp_servers"]; ok {
					clusterList.HypervisorNtpServers = (hypervntps.(string))
				}
				if clsmembers, ok := clst["cluster_members"]; ok {
					clsm := clsmembers.([]interface{})
					res := []string{}
					for _, v := range clsm {
						res = append(res, v.(string))
					}
					clusterList.ClusterMembers = res
				}
				cls = append(cls, clusterList)
			}
			return cls, nil
		}
	}
	return clstr, nil
}

func expandNodes(pr interface{}) []*foundation.Node {
	nodesList := pr.([]interface{})
	nodes := make([]*foundation.Node, len(nodesList))

	for i, p := range nodesList {
		node := p.(map[string]interface{})
		nodeList := &foundation.Node{}
		if ipmiNetmask, ipmiNetOk := node["ipmi_netmask"]; ipmiNetOk {
			nodeList.IPMINetmask = ipmiNetmask.(string)
		}
		if ipmiGateway, ipmiGateOk := node["ipmi_gateway"]; ipmiGateOk {
			nodeList.IPMIGateway = ipmiGateway.(string)
		}
		if ipv6, ipv6ok := node["ipv6_address"]; ipv6ok {
			nodeList.Ipv6Address = (ipv6.(string))
		}
		if np, npok := node["node_position"]; npok {
			nodeList.NodePosition = (np.(string))
		}
		if imgd, imgdok := node["image_delay"]; imgdok && imgd.(int) != 0 {
			nodeList.ImageDelay = utils.Int64Ptr(int64(imgd.(int)))
		}
		if hypervhostname, hpyervhostnok := node["hypervisor_hostname"]; hpyervhostnok {
			nodeList.HypervisorHostname = (hypervhostname.(string))
		}
		if cvmram, cvmramok := node["cvm_gb_ram"]; cvmramok && cvmram.(int) != 0 {
			nodeList.CvmGbRAM = utils.Int64Ptr(int64(cvmram.(int)))
		}
		if devicehint, devicehintok := node["device_hint"]; devicehintok {
			nodeList.DeviceHint = (devicehint.(string))
		}
		if bondmode, bondmodeok := node["bond_mode"]; bondmodeok {
			nodeList.BondMode = (bondmode.(string))
		}
		if bondUl, bondUlOk := node["bond_uplinks"]; bondUlOk && len(bondUl.([]interface{})) > 0 {
			//convert to []string from []inerfaces
			var b []string
			for _, v := range bondUl.([]interface{}) {
				b = append(b, v.(string))
			}
			nodeList.BondUplinks = b
		}
		if rdmapass, rdmapassok := node["rdma_passthrough"]; rdmapassok && rdmapass.(bool) {
			nodeList.RdmaPassthrough = utils.BoolPtr(rdmapass.(bool))
		}
		if clsid, clsidok := node["cluster_id"]; clsidok {
			nodeList.ClusterID = (clsid.(string))
		}
		if ucsmns, ucsmnsok := node["ucsm_node_serial"]; ucsmnsok {
			nodeList.UcsmNodeSerial = (ucsmns.(string))
		}
		if hypervip, hypervipok := node["hypervisor_ip"]; hypervipok {
			nodeList.HypervisorIP = (hypervip.(string))
		}
		if ns, nsok := node["node_serial"]; nsok {
			nodeList.NodeSerial = (ns.(string))
		}
		if ipmicn, ipmicnok := node["ipmi_configure_now"]; ipmicnok && ipmicn.(bool) {
			nodeList.IpmiConfigureNow = utils.BoolPtr(ipmicn.(bool))
		}
		if imgsuc, imgsucok := node["image_successful"]; imgsucok && imgsuc.(bool) {
			nodeList.ImageSuccessful = utils.BoolPtr(imgsuc.(bool))
		}
		if ipv6i, ipv6iok := node["ipv6_interface"]; ipv6iok {
			nodeList.Ipv6Interface = (ipv6i.(string))
		}
		if cvmvcpu, cvmvcpuok := node["cvm_num_vcpus"]; cvmvcpuok && cvmvcpu.(int) != 0 {
			nodeList.CvmNumVcpus = utils.Int64Ptr(int64(cvmvcpu.(int)))
		}
		if ipmimac, ipmimacok := node["ipmi_mac"]; ipmimacok {
			nodeList.IpmiMac = (ipmimac.(string))
		}
		if rma, rmaok := node["rdma_mac_addr"]; rmaok {
			nodeList.RdmaMacAddr = (rma.(string))
		}
		if cni, cniok := node["current_network_interface"]; cniok {
			nodeList.CurrentNetworkInterface = (cni.(string))
		}
		if hypervip, hypervipok := node["hypervisor_ip"]; hypervipok {
			nodeList.HypervisorIP = (hypervip.(string))
		}
		if hyperv, hypervok := node["hypervisor"]; hypervok {
			nodeList.Hypervisor = (hyperv.(string))
		}
		if bondlacprate, bondlacprateok := node["bond_lacp_rate"]; bondlacprateok {
			nodeList.BondLacpRate = (bondlacprate.(string))
		}
		if imgnow, imgnowok := node["image_now"]; imgnowok {
			nodeList.ImageNow = utils.BoolPtr(imgnow.(bool))
		}
		if ucsmmode, ucsmmodeok := node["ucsm_managed_mode"]; ucsmmodeok {
			nodeList.UcsmManagedMode = (ucsmmode.(string))
		}
		if ipmi, ipmiok := node["ipmi_ip"]; ipmiok {
			nodeList.IpmiIP = (ipmi.(string))
		}
		if cvmvlantag, cvmvlantagok := node["current_cvm_vlan_tag"]; cvmvlantagok && cvmvlantag.(int) != 0 {
			nodeList.CurrentCvmVlanTag = utils.Int64Ptr(int64(cvmvlantag.(int)))
		}
		if cvmip, cvmipok := node["cvm_ip"]; cvmipok {
			nodeList.CvmIP = (cvmip.(string))
		}
		if exboots, exbootsok := node["exlude_boot_serial"]; exbootsok {
			nodeList.ExludeBootSerial = (exboots.(string))
		}
		if lbootspace, lbootspaceok := node["mitigate_low_boot_space"]; lbootspaceok && lbootspace.(bool) {
			nodeList.MitigateLowBootSpace = utils.BoolPtr(lbootspace.(bool))
		}
		if ucsmParams, ucsmParamsok := node["ucsm_params"]; ucsmParamsok {
			nodeList.UcsmParams = expandUcsmParams(ucsmParams)
		}
		if vswitch, vswitchesok := node["vswitches"]; vswitchesok {
			nodeList.Vswitches = expandVswitches(vswitch)
		}
		if ipmiUser, ok := node["ipmi_user"]; ok {
			nodeList.IpmiUser = (ipmiUser.(string))
		}
		if ipmiPassword, ok := node["ipmi_password"]; ok {
			nodeList.IpmiPassword = (ipmiPassword.(string))
		}
		nodes[i] = nodeList
	}
	return nodes
}

func expandBlocks(d *schema.ResourceData) ([]*foundation.Block, error) {
	if blocks, ok := d.GetOk("blocks"); ok {
		set := blocks.([]interface{})
		outbound := make([]*foundation.Block, len(set))

		for k, v := range set {
			block := &foundation.Block{}

			entry := v.(map[string]interface{})

			if nodes, nodesok := entry["nodes"]; nodesok {
				block.Nodes = expandNodes(nodes)
			}

			if blockid, blockidok := entry["block_id"]; blockidok {
				block.BlockID = (blockid.(string))
			}
			outbound[k] = block
		}
		return outbound, nil
	}
	return nil, nil
}

// This method will look into individual node and cluster creation status and create a collection of errors for errored out processes
func collectIndividualErrorDiagnostics(progress *foundation.ImageNodesProgressResponse) diag.Diagnostics {
	// create empty diagnostics
	var diags diag.Diagnostics

	// append errors for failed node imaging
	for _, v := range progress.Nodes {
		if utils.Float64Value(v.PercentComplete) < float64(AggregatePercentComplete) {
			message := ""
			for _, v1 := range v.Messages {
				message += v1 + ". "
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Node imaging for CVM IP: %s failed with error:  %s.", v.CvmIP, v.Status),
				Detail:   message,
			})
		}
	}

	// append errors for failed cluster creation
	for _, v := range progress.Clusters {
		if utils.Float64Value(v.PercentComplete) < float64(AggregatePercentComplete) {
			message := ""
			for _, v1 := range v.Messages {
				message += v1 + ". "
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cluster creation for Cluster : %s failed with error:  %s.", v.ClusterName, v.Status),
				Detail:   message,
			})
		}
	}

	return diags
}
