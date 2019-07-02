package aws

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/pricing"
	awsutil "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultInstanceClass(t *testing.T) {
	preferredInstanceClasses := []string{"m4", "m5"} // decreasing precedence

	ssn, err := awsutil.GetSession()
	if err != nil {
		t.Fatal(err)
	}

	exists := struct{}{}
	pricingInstanceClasses := map[string]map[string]struct{}{}

	pricingClient := pricing.New(ssn, aws.NewConfig().WithRegion("us-east-1"))
	err = pricingClient.GetProductsPages(
		&pricing.GetProductsInput{
			ServiceCode: aws.String("AmazonEC2"),
			Filters: []*pricing.Filter{
				{
					Field: aws.String("tenancy"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Shared"),
				},
				{
					Field: aws.String("productFamily"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Compute Instance"),
				},
				{
					Field: aws.String("operatingSystem"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("Linux"),
				},
				{
					Field: aws.String("instanceFamily"),
					Type:  aws.String("TERM_MATCH"),
					Value: aws.String("General purpose"),
				},
			},
		},
		func(result *pricing.GetProductsOutput, lastPage bool) bool {
			for _, priceList := range result.PriceList {
				product := priceList["product"].(map[string]interface{})
				attr := product["attributes"].(map[string]interface{})
				location := attr["location"].(string)
				instanceType := attr["instanceType"].(string)
				instanceClassSlice := strings.Split(instanceType, ".")
				instanceClass := instanceClassSlice[0]
				_, ok := pricingInstanceClasses[location]
				if ok {
					pricingInstanceClasses[location][instanceClass] = exists
				} else {
					pricingInstanceClasses[location] = map[string]struct{}{instanceClass: exists}
				}
			}
			return !lastPage
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	regions := map[string]string{ // seed with locations that don't match AWS's usual names
		"AWS GovCloud (US)":          "us-gov-west-1",
		"AWS GovCloud (US-East)":     "us-gov-east-1",
		"Asia Pacific (Hong Kong)":   "ap-east-1",
		"Asia Pacific (Osaka-Local)": "ap-northeast-3",
		"EU (Stockholm)":             "eu-north-1",
		"South America (Sao Paulo)":  "sa-east-1",
	}

	for location, classes := range pricingInstanceClasses {
		t.Run(location, func(t *testing.T) {
			region, ok := regions[location]
			if !ok {
				for slug, name := range validation.Regions {
					if strings.Contains(location, name) {
						regions[location] = slug
						region = slug
						break
					}
				}
				if region == "" {
					t.Fatal("not a recognized region")
				}
			}

			ec2Client := ec2.New(ssn, aws.NewConfig().WithRegion(region))
			zonesResponse, err := ec2Client.DescribeAvailabilityZones(nil)
			if err != nil {
				t.Logf("no direct access to region, assuming full support: %v", err)

				var match string
				for _, instanceClass := range preferredInstanceClasses {
					if _, ok := classes[instanceClass]; ok {
						match = instanceClass
						break
					}
				}

				if match == "" {
					t.Fatalf("none of the preferred instance classes are priced: %v", classes)
				}

				t.Log(classes)
				assert.Equal(t, defaults.InstanceClass(region), match)
				return
			}

			zones := make(map[string]struct{}, len(zonesResponse.AvailabilityZones))
			for _, zone := range zonesResponse.AvailabilityZones {
				zones[*zone.ZoneName] = exists
			}

			available := make(map[string]map[string]struct{}, len(preferredInstanceClasses))
			var allowed []string

			for _, instanceClass := range preferredInstanceClasses {
				if _, ok := classes[instanceClass]; !ok {
					t.Logf("skip the unpriced %s", instanceClass)
					continue
				}

				available[instanceClass] = make(map[string]struct{}, len(zones))
				exampleInstanceType := fmt.Sprintf("%s.large", instanceClass)
				err := ec2Client.DescribeReservedInstancesOfferingsPages(
					&ec2.DescribeReservedInstancesOfferingsInput{
						Filters: []*ec2.Filter{
							{Name: aws.String("scope"), Values: []*string{aws.String("Availability Zone")}},
						},
						InstanceTenancy:    aws.String("default"),
						InstanceType:       &exampleInstanceType,
						ProductDescription: aws.String("Linux/UNIX"),
					},
					func(results *ec2.DescribeReservedInstancesOfferingsOutput, lastPage bool) bool {
						for _, offering := range results.ReservedInstancesOfferings {
							if offering.AvailabilityZone == nil {
								continue
							}

							available[instanceClass][*offering.AvailabilityZone] = exists
						}

						return !lastPage
					},
				)
				if err != nil {
					t.Fatal(err)
				}

				if reflect.DeepEqual(available[instanceClass], zones) {
					allowed = append(allowed, instanceClass)
				}
			}

			if len(allowed) == 0 {
				t.Fatalf("none of the preferred instance classes are fully supported: %v", available)
			}

			t.Log(available)
			assert.Contains(t, allowed, defaults.InstanceClass(region))
		})
	}
}
