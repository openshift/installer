package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsTablesRead,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MaxItems: 4,
						},
						"time_to_live": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type OtsTableInfo struct {
	instanceName string
	tableName    string
	primaryKey   []*tablestore.PrimaryKeySchema
	timeToLive   int
	maxVersion   int
}

func dataSourceAlicloudOtsTablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	instanceName := d.Get("instance_name").(string)

	object, err := otsService.ListOtsTable(instanceName)
	if err != nil {
		return WrapError(err)
	}

	idsMap := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, x := range v.([]interface{}) {
			if x == nil {
				continue
			}
			idsMap[x.(string)] = true
		}
	}

	var nameReg *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		nameReg = regexp.MustCompile(v.(string))
	}

	var filteredTableNames []string
	for _, tableName := range object.TableNames {
		//name_regex mismatch
		if nameReg != nil && !nameReg.MatchString(tableName) {
			continue
		}
		// ids mismatch
		if len(idsMap) != 0 {
			id := fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName)
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}
		filteredTableNames = append(filteredTableNames, tableName)
	}

	// get full table info via DescribeTable
	var allTableInfos []OtsTableInfo
	for _, tableName := range filteredTableNames {
		object, err := otsService.DescribeOtsTable(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
		if err != nil {
			return WrapError(err)
		}
		allTableInfos = append(allTableInfos, OtsTableInfo{
			instanceName: instanceName,
			tableName:    object.TableMeta.TableName,
			primaryKey:   object.TableMeta.SchemaEntry,
			timeToLive:   object.TableOption.TimeToAlive,
			maxVersion:   object.TableOption.MaxVersion,
		})
	}

	return otsTablesDescriptionAttributes(d, allTableInfos, meta)
}

func otsTablesDescriptionAttributes(d *schema.ResourceData, tableInfos []OtsTableInfo, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, table := range tableInfos {
		id := fmt.Sprintf("%s:%s", table.instanceName, table.tableName)
		mapping := map[string]interface{}{
			"id":            id,
			"instance_name": table.instanceName,
			"table_name":    table.tableName,
			"time_to_live":  table.timeToLive,
			"max_version":   table.maxVersion,
		}
		var primaryKey []map[string]interface{}
		for _, pk := range table.primaryKey {
			pkColumn := make(map[string]interface{})
			pkColumn["name"] = *pk.Name
			pkColumn["type"] = otsService.convertPrimaryKeyType(*pk.Type)
			primaryKey = append(primaryKey, pkColumn)
		}
		mapping["primary_key"] = primaryKey

		names = append(names, table.tableName)
		ids = append(ids, id)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("tables", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
