package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceEBSEncryptionByDefault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEBSEncryptionByDefaultRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
func dataSourceEBSEncryptionByDefaultRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	res, err := conn.GetEbsEncryptionByDefault(&ec2.GetEbsEncryptionByDefaultInput{})
	if err != nil {
		return fmt.Errorf("Error reading default EBS encryption toggle: %w", err)
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("enabled", res.EbsEncryptionByDefault)

	return nil
}
