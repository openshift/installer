package aws

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	awsutil "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultInstanceClass(t *testing.T) {
	ssn, err := awsutil.GetSession()
	if err != nil {
		t.Fatal(err)
	}

	exists := struct{}{}
	instanceClasses := map[string]map[string]struct{}{}

	client := pricing.New(ssn, aws.NewConfig().WithRegion("us-east-1"))

	err = client.GetProductsPages(
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
				_, ok := instanceClasses[location]
				if ok {
					instanceClasses[location][instanceClass] = exists
				} else {
					instanceClasses[location] = map[string]struct{}{instanceClass: exists}
				}
			}
			return !lastPage
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	regions := map[string]string{ // seed with locations that don't match AWS's usual names
		"South America (Sao Paulo)": "sa-east-1",
		"AWS GovCloud (US)":         "us-gov-west-1",
	}

	for location, classes := range instanceClasses {
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

			class := ""
			// ordered list of prefered instance classes
			for _, preferredClass := range []string{"m4", "m5"} {
				if _, ok := classes[preferredClass]; ok {
					class = preferredClass
					break
				}
			}
			if class == "" {
				t.Fatalf("does not support any preferred classes: %v", classes)
			}
			defaultClass := defaults.InstanceClass(region)
			assert.Equal(t, defaultClass, class)
		})
	}
}
