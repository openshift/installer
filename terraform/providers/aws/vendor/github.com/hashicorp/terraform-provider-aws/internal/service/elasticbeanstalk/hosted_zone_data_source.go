package elasticbeanstalk

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// See http://docs.aws.amazon.com/general/latest/gr/rande.html#elasticbeanstalk_region

var HostedZoneIDs = map[string]string{
	endpoints.AfSouth1RegionID:     "Z1EI3BVKMKK4AM",
	endpoints.ApSoutheast1RegionID: "Z16FZ9L249IFLT",
	endpoints.ApSoutheast2RegionID: "Z2PCDNR3VC2G1N",
	endpoints.ApEast1RegionID:      "ZPWYUBWRU171A",
	endpoints.ApNortheast1RegionID: "Z1R25G3KIG2GBW",
	endpoints.ApNortheast2RegionID: "Z3JE5OI70TWKCP",
	endpoints.ApNortheast3RegionID: "ZNE5GEY1TIAGY",
	endpoints.ApSouth1RegionID:     "Z18NTBI3Y7N9TZ",
	endpoints.CaCentral1RegionID:   "ZJFCZL7SSZB5I",
	endpoints.EuCentral1RegionID:   "Z1FRNW7UH4DEZJ",
	endpoints.EuNorth1RegionID:     "Z23GO28BZ5AETM",
	endpoints.EuSouth1RegionID:     "Z10VDYYOA2JFKM",
	endpoints.EuWest1RegionID:      "Z2NYPWQ7DFZAZH",
	endpoints.EuWest2RegionID:      "Z1GKAAAUGATPF1",
	endpoints.EuWest3RegionID:      "Z5WN6GAYWG5OB",
	// endpoints.MeCentral1RegionID:   "",
	endpoints.MeSouth1RegionID:   "Z2BBTEKR2I36N2",
	endpoints.SaEast1RegionID:    "Z10X7K2B4QSOFV",
	endpoints.UsEast1RegionID:    "Z117KPS5GTRQ2G",
	endpoints.UsEast2RegionID:    "Z14LCN19Q5QHIC",
	endpoints.UsWest1RegionID:    "Z1LQECGX5PH1X",
	endpoints.UsWest2RegionID:    "Z38NKT9BP95V3O",
	endpoints.UsGovEast1RegionID: "Z35TSARG0EJ4VU",
	endpoints.UsGovWest1RegionID: "Z4KAURWC4UUUG",
}

func DataSourceHostedZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostedZoneRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceHostedZoneRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	zoneID, ok := HostedZoneIDs[region]

	if !ok {
		return fmt.Errorf("Unsupported region: %s", region)
	}

	d.SetId(zoneID)
	d.Set("region", region)
	return nil
}
