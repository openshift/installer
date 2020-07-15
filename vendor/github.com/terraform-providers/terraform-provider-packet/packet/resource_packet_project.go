package packet

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

var uuidRE = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func resourcePacketProject() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketProjectCreate,
		Read:   resourcePacketProjectRead,
		Update: resourcePacketProjectUpdate,
		Delete: resourcePacketProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_transfer": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"payment_method_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(strings.Trim(old, `"`)) == strings.ToLower(strings.Trim(new, `"`))
				},
				ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
			},

			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(strings.Trim(old, `"`)) == strings.ToLower(strings.Trim(new, `"`))
				},
				ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
			},
			"bgp_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"local", "global"}, false),
						},
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"md5": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_prefix": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func expandBGPConfig(d *schema.ResourceData) packngo.CreateBGPConfigRequest {
	bgpCreateRequest := packngo.CreateBGPConfigRequest{
		DeploymentType: d.Get("bgp_config.0.deployment_type").(string),
		Asn:            d.Get("bgp_config.0.asn").(int),
	}
	md5, ok := d.GetOk("bgp_config.0.md5")
	if ok {
		bgpCreateRequest.Md5 = md5.(string)
	}

	return bgpCreateRequest

}

func resourcePacketProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.ProjectCreateRequest{
		Name:           d.Get("name").(string),
		OrganizationID: d.Get("organization_id").(string),
	}

	project, _, err := client.Projects.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(project.ID)

	_, hasBGPConfig := d.GetOk("bgp_config")
	if hasBGPConfig {
		bgpCR := expandBGPConfig(d)
		_, err := client.BGPConfig.Create(project.ID, bgpCR)
		if err != nil {
			return friendlyError(err)
		}
	}

	backendTransfer := d.Get("backend_transfer").(bool)
	if backendTransfer {
		pur := packngo.ProjectUpdateRequest{BackendTransfer: &backendTransfer}
		_, _, err := client.Projects.Update(project.ID, &pur)
		if err != nil {
			return friendlyError(err)
		}
	}
	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	proj, _, err := client.Projects.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(proj.ID)
	d.Set("payment_method_id", path.Base(proj.PaymentMethod.URL))
	d.Set("name", proj.Name)
	d.Set("organization_id", path.Base(proj.Organization.URL))
	d.Set("created", proj.Created)
	d.Set("updated", proj.Updated)
	d.Set("backend_transfer", proj.BackendTransfer)

	bgpConf, _, err := client.BGPConfig.Get(proj.ID, nil)

	if (err == nil) && (bgpConf != nil) {
		// guard against an empty struct
		if bgpConf.ID != "" {
			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
			if err != nil {
				err = friendlyError(err)
				return err
			}
		}
	}
	return nil
}

func flattenBGPConfig(l *packngo.BGPConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if l == nil {
		return nil
	}

	r := make(map[string]interface{})

	if l.Status != "" {
		r["status"] = l.Status
	}
	if l.DeploymentType != "" {
		r["deployment_type"] = l.DeploymentType
	}
	if l.Md5 != "" {
		r["md5"] = l.Md5
	}
	if l.Asn != 0 {
		r["asn"] = l.Asn
	}
	if l.MaxPrefix != 0 {
		r["max_prefix"] = l.MaxPrefix
	}

	result = append(result, r)

	return result
}

func resourcePacketProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	updateRequest := &packngo.ProjectUpdateRequest{}
	if d.HasChange("name") {
		pName := d.Get("name").(string)
		updateRequest.Name = &pName
	}
	if d.HasChange("payment_method_id") {
		pPayment := d.Get("payment_method_id").(string)
		updateRequest.PaymentMethodID = &pPayment
	}
	if d.HasChange("backend_transfer") {
		pBT := d.Get("backend_transfer").(bool)
		updateRequest.BackendTransfer = &pBT
	}
	if d.HasChange("bgp_config") {
		o, n := d.GetChange("bgp_config")
		oldarr := o.([]interface{})
		newarr := n.([]interface{})
		if len(newarr) == 1 {
			bgpCreateRequest := expandBGPConfig(d)
			_, err := client.BGPConfig.Create(d.Id(), bgpCreateRequest)
			if err != nil {
				return friendlyError(err)
			}
		} else {
			if len(oldarr) == 1 {
				m := oldarr[0].(map[string]interface{})

				bgpConfStr := fmt.Sprintf(
					"bgp_config {\n"+
						"  deployment_type = \"%s\"\n"+
						"  md5 = \"%s\"\n"+
						"  asn = %d\n"+
						"}", m["deployment_type"].(string), m["md5"].(string),
					m["asn"].(int))

				errStr := fmt.Errorf("BGP Config can not be removed from a project, please add back\n%s", bgpConfStr)
				return friendlyError(errStr)
			}
		}
	} else {
		_, _, err := client.Projects.Update(d.Id(), updateRequest)
		if err != nil {
			return friendlyError(err)
		}
	}

	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.Projects.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
