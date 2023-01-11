package cloudtrail

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// See http://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-supported-regions.html
// See https://docs.aws.amazon.com/govcloud-us/latest/ug-east/verifying-cloudtrail.html
// See https://docs.aws.amazon.com/govcloud-us/latest/ug-west/verifying-cloudtrail.html

var ServiceAccountPerRegionMap = map[string]string{
	endpoints.AfSouth1RegionID:     "525921808201",
	endpoints.ApEast1RegionID:      "119688915426",
	endpoints.ApNortheast1RegionID: "216624486486",
	endpoints.ApNortheast2RegionID: "492519147666",
	endpoints.ApNortheast3RegionID: "765225791966",
	endpoints.ApSouth1RegionID:     "977081816279",
	endpoints.ApSoutheast1RegionID: "903692715234",
	endpoints.ApSoutheast2RegionID: "284668455005",
	endpoints.ApSoutheast3RegionID: "069019280451",
	endpoints.CaCentral1RegionID:   "819402241893",
	endpoints.CnNorth1RegionID:     "193415116832",
	endpoints.CnNorthwest1RegionID: "681348832753",
	endpoints.EuCentral1RegionID:   "035351147821",
	endpoints.EuNorth1RegionID:     "829690693026",
	endpoints.EuSouth1RegionID:     "669305197877",
	endpoints.EuWest1RegionID:      "859597730677",
	endpoints.EuWest2RegionID:      "282025262664",
	endpoints.EuWest3RegionID:      "262312530599",
	endpoints.MeCentral1RegionID:   "585772288577",
	endpoints.MeSouth1RegionID:     "034638983726",
	endpoints.SaEast1RegionID:      "814480443879",
	endpoints.UsEast1RegionID:      "086441151436",
	endpoints.UsEast2RegionID:      "475085895292",
	endpoints.UsGovEast1RegionID:   "608710470296",
	endpoints.UsGovWest1RegionID:   "608710470296",
	endpoints.UsWest1RegionID:      "388731089494",
	endpoints.UsWest2RegionID:      "113285607260",
}

func DataSourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	if accid, ok := ServiceAccountPerRegionMap[region]; ok {
		d.SetId(accid)
		arn := arn.ARN{
			Partition: meta.(*conns.AWSClient).Partition,
			Service:   "iam",
			AccountID: accid,
			Resource:  "root",
		}.String()
		d.Set("arn", arn)

		return nil
	}

	return fmt.Errorf("Unknown region (%q)", region)
}
