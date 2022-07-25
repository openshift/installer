package endpoints

import (
	"regexp"
)

func Partitions() []Partition {
	ps := make([]Partition, len(partitions))
	for i := 0; i < len(partitions); i++ {
		ps[i] = partitions[i].Partition()
	}
	return ps
}

type Partition struct {
	id string
	p  *partition
}

func (p Partition) Regions() []string {
	rs := make([]string, len(p.p.regions))
	for i, v := range p.p.regions {
		rs[i] = v
	}
	return rs
}

func PartitionForRegion(regionID string) string {
	for _, p := range partitions {
		if p.regionRegex.MatchString(regionID) {
			return p.id
		}
	}

	return ""
}

type partition struct {
	id          string
	regionRegex *regexp.Regexp
	regions     []string
}

func (p partition) Partition() Partition {
	return Partition{
		id: p.id,
		p:  &p,
	}
}

// TODO: this should be generated from the AWS SDK source data
// Data from https://github.com/aws/aws-sdk-go/blob/main/models/endpoints/endpoints.json.
var partitions = []partition{
	{
		id:          "aws",
		regionRegex: regexp.MustCompile(`^(us|eu|ap|sa|ca|me|af)\-\w+\-\d+$`),
		regions: []string{
			"af-south-1",     // Africa (Cape Town).
			"ap-east-1",      // Asia Pacific (Hong Kong).
			"ap-northeast-1", // Asia Pacific (Tokyo).
			"ap-northeast-2", // Asia Pacific (Seoul).
			"ap-northeast-3", // Asia Pacific (Osaka).
			"ap-south-1",     // Asia Pacific (Mumbai).
			"ap-southeast-1", // Asia Pacific (Singapore).
			"ap-southeast-2", // Asia Pacific (Sydney).
			"ap-southeast-3", // Asia Pacific (Jakarta).
			"ca-central-1",   // Canada (Central).
			"eu-central-1",   // Europe (Frankfurt).
			"eu-north-1",     // Europe (Stockholm).
			"eu-south-1",     // Europe (Milan).
			"eu-west-1",      // Europe (Ireland).
			"eu-west-2",      // Europe (London).
			"eu-west-3",      // Europe (Paris).
			"me-south-1",     // Middle East (Bahrain).
			"sa-east-1",      // South America (Sao Paulo).
			"us-east-1",      // US East (N. Virginia).
			"us-east-2",      // US East (Ohio).
			"us-west-1",      // US West (N. California).
			"us-west-2",      // US West (Oregon).
		},
	},
	{
		id:          "aws-cn",
		regionRegex: regexp.MustCompile(`^cn\-\w+\-\d+$`),
		regions: []string{
			"cn-north-1",     // China (Beijing).
			"cn-northwest-1", // China (Ningxia).
		},
	},
	{
		id:          "aws-us-gov",
		regionRegex: regexp.MustCompile(`^us\-gov\-\w+\-\d+$`),
		regions: []string{
			"us-gov-east-1", // AWS GovCloud (US-East).
			"us-gov-west-1", // AWS GovCloud (US-West).
		},
	},
	{
		id:          "aws-iso",
		regionRegex: regexp.MustCompile(`^us\-iso\-\w+\-\d+$`),
		regions: []string{
			"us-iso-east-1", // US ISO East.
			"us-iso-west-1", // US ISO WEST.
		},
	},
	{
		id:          "aws-iso-b",
		regionRegex: regexp.MustCompile(`^us\-isob\-\w+\-\d+$`),
		regions: []string{
			"us-isob-east-1", // US ISOB East (Ohio).
		},
	},
}
