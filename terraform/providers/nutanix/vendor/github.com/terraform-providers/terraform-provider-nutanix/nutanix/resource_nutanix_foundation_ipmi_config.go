package nutanix

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
)

func resourceNutanixFoundationIPMIConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoundationIPMIConfigCreate,
		ReadContext:   resourceFoundationIPMIConfigRead,
		DeleteContext: resourceFoundationIPMIConfigDelete,
		Schema: map[string]*schema.Schema{
			"ipmi_user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipmi_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipmi_netmask": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipmi_gateway": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
									"ipmi_mac": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"ipmi_configure_now": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
									"ipmi_ip": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"ipmi_configure_successful": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ipmi_message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"block_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceFoundationIPMIConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// get foundation client
	conn := meta.(*Client).FoundationClientAPI

	//api input spec
	inpSpec := &foundation.IPMIConfigAPIInput{}

	//setting required values which will be common for all nodes
	ipmiUser, ok := d.GetOk("ipmi_user")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("ipmi_user"))
	}
	ipmiPass, ok := d.GetOk("ipmi_password")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("ipmi_password"))
	}
	ipmiNetmask, ok := d.GetOk("ipmi_netmask")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("ipmi_netmask"))
	}
	ipmiGateway, ok := d.GetOk("ipmi_gateway")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("ipmi_gateway"))
	}
	inpSpec.IpmiUser = ipmiUser.(string)
	inpSpec.IpmiPassword = ipmiPass.(string)
	inpSpec.IpmiNetmask = ipmiNetmask.(string)
	inpSpec.IpmiGateway = ipmiGateway.(string)

	// get blocks details
	b, ok := d.GetOk("blocks")
	if !ok {
		return diag.Errorf(getRequiredErrorMessage("blocks"))
	}
	// create struct for blocks
	blocks := b.([]interface{})
	blocksInput := make([]foundation.IPMIConfigBlockInput, len(blocks))
	for k, v := range blocks {
		blockSpec := foundation.IPMIConfigBlockInput{}
		val := v.(map[string]interface{})
		if blockID, ok := val["block_id"]; ok {
			blockSpec.BlockID = blockID.(string)
		}

		n, ok := val["nodes"]
		if !ok {
			return diag.Errorf("Please provide node ipmi details for %dth block", k)
		}

		// create struct for node list for every block
		nodes := n.([]interface{})
		nodeList := make([]foundation.IPMIConfigNodeInput, len(nodes))
		for k1, v1 := range nodes {
			nodeSpec := foundation.IPMIConfigNodeInput{}
			val1 := v1.(map[string]interface{})

			ipmiMac, ok := val1["ipmi_mac"]
			if !ok {
				return diag.Errorf("Please provide ipmi_mac for %dth block and %dth node", k, k1)
			}
			ipmiIP, ok := val1["ipmi_ip"]
			if !ok {
				return diag.Errorf("Please provide ipmi_ip for %dth block and %dth node", k, k1)
			}
			ipmiConfigureNow, ok := val1["ipmi_configure_now"]
			if !ok {
				return diag.Errorf("Please provide ipmi_configure_now for %dth block and %dth node", k, k1)
			}
			nodeSpec.IpmiMac = ipmiMac.(string)
			nodeSpec.IpmiIP = ipmiIP.(string)
			nodeSpec.IpmiConfigureNow = ipmiConfigureNow.(bool)
			nodeList[k1] = nodeSpec
		}
		blockSpec.Nodes = nodeList
		blocksInput[k] = blockSpec
	}
	inpSpec.Blocks = blocksInput

	// call networking service method ConfigureIPMI to configure IPMI to given block of nodes
	resp, err := conn.Networking.ConfigureIPMI(ctx, inpSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	// incase of errored response
	if resp.Error != nil {
		// check if error details exists for every ipmi IP
		var diags diag.Diagnostics
		for k, v := range resp.Error.Details {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("IPMI config failed for IPMI IP: %s", k),
				Detail:   v,
			})
		}
		if len(diags) > 0 {
			return diags
		}
		// incase there is no error details
		return diag.Errorf(resp.Error.Message)
	}

	if setErr := d.Set("ipmi_user", resp.IpmiUser); setErr != nil {
		return diag.FromErr(err)
	}
	if setErr := d.Set("ipmi_password", resp.IpmiPassword); setErr != nil {
		return diag.FromErr(err)
	}
	if setErr := d.Set("ipmi_gateway", resp.IpmiGateway); setErr != nil {
		return diag.FromErr(err)
	}
	if setErr := d.Set("ipmi_netmask", resp.IpmiNetmask); setErr != nil {
		return diag.FromErr(err)
	}

	// flatten blocks and nodes list for every block and set appropriately
	blocksResponse := make([]map[string]interface{}, len(resp.Blocks))
	for k, v := range resp.Blocks {
		block := make(map[string]interface{})
		block["block_id"] = v.BlockID

		//nodes
		nodes := make([]map[string]interface{}, len(v.Nodes))
		for k1, v1 := range v.Nodes {
			node := make(map[string]interface{})
			node["ipmi_configure_successful"] = v1.IpmiConfigureSuccessful
			node["ipmi_configure_now"] = true
			node["ipmi_ip"] = v1.IpmiIP
			node["ipmi_mac"] = v1.IpmiMac
			node["ipmi_message"] = v1.IpmiMessage
			nodes[k1] = node
		}
		block["nodes"] = nodes
		blocksResponse[k] = block
	}

	if setErr := d.Set("blocks", blocksResponse); setErr != nil {
		return diag.FromErr(setErr)
	}
	d.SetId(resource.UniqueId())
	return nil
}

// getRequiredErrorMessage returns required attribute error as per attribute name given
func getRequiredErrorMessage(attr string) string {
	return fmt.Sprintf("please provide the value of required attribute: %s", attr)
}

// there is no read API for IPIMI config hence skipping read
func resourceFoundationIPMIConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// there is no delete API for IPIMI config hence skipping delete
func resourceFoundationIPMIConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
