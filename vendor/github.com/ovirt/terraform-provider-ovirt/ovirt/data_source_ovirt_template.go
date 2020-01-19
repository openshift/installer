package ovirt

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func dataSourceOvirtTemplates() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceOvirtTemplatesRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchemaWith(
				"max", "criteria", "case_sensitive"),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_shares": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The memory of VM which associated with oVirt Template, in Megabytes(MB)",
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceOvirtTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	templatesReq := conn.SystemService().TemplatesService().List()

	search, searchOK := d.GetOk("search")
	nameRegex, nameRegexOK := d.GetOk("name_regex")

	if searchOK {
		searchMap := search.(map[string]interface{})
		searchCriteria, searchCriteriaOK := searchMap["criteria"]
		searchMax, searchMaxOK := searchMap["max"]
		searchCaseSensitive, searchCaseSensitiveOK := searchMap["case_sensitive"]
		if !searchCriteriaOK && !searchMaxOK && !searchCaseSensitiveOK {
			return fmt.Errorf("One of criteria or max or case_sensitive in search must be assigned")
		}

		if searchCriteriaOK {
			templatesReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			templatesReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			templatesReq.CaseSensitive(csBool)
		}
	}
	templatesResp, err := templatesReq.Send()
	if err != nil {
		return err
	}
	templates, ok := templatesResp.Templates()
	if !ok || len(templates.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredTemplates []*ovirtsdk4.Template
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, t := range templates.Slice() {
			if r.MatchString(t.MustName()) {
				filteredTemplates = append(filteredTemplates, t)
			}
		}
	} else {
		filteredTemplates = templates.Slice()[:]
	}

	if len(filteredTemplates) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return templatesDescriptionAttributes(d, filteredTemplates, meta)
}

func templatesDescriptionAttributes(d *schema.ResourceData, templates []*ovirtsdk4.Template, meta interface{}) error {
	var s []map[string]interface{}

	for _, v := range templates {
		mapping := map[string]interface{}{
			"id":            v.MustId(),
			"name":          v.MustName(),
			"cpu_shares":    v.MustCpuShares(),
			"memory":        v.MustMemory() / int64(math.Pow(2, 20)),
			"creation_time": v.MustCreationTime().Format(time.RFC3339),
			"status":        v.MustStatus(),
		}
		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("templates", s)
}
