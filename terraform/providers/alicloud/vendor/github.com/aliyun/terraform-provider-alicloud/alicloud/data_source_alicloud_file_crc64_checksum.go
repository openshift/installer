package alicloud

import (
	"fmt"
	"hash/crc64"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudFileCRC64Checksum() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFileCRC64ChecksumRead,

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlicloudFileCRC64ChecksumRead(d *schema.ResourceData, meta interface{}) error {
	filename := d.Get("filename")
	file, err := loadFileContent(filename.(string))
	if err != nil {
		return WrapError(err)
	}
	table := crc64.MakeTable(crc64.ECMA)
	checkSum := fmt.Sprintf("%d", crc64.Checksum(file, table))
	d.Set("checksum", checkSum)
	d.SetId(checkSum)
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), map[string]string{"checksum": checkSum})
	}
	return nil
}
