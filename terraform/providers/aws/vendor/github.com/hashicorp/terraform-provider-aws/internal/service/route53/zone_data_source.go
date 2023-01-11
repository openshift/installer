package route53

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_zone": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"caller_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"resource_record_set_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name_servers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"linked_service_principal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"linked_service_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).Route53Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name, nameExists := d.GetOk("name")
	name = name.(string)
	id, idExists := d.GetOk("zone_id")
	vpcId, vpcIdExists := d.GetOk("vpc_id")
	tags := tftags.New(d.Get("tags").(map[string]interface{})).IgnoreAWS()

	if nameExists && idExists {
		return fmt.Errorf("zone_id and name arguments can't be used together")
	}

	if !nameExists && !idExists {
		return fmt.Errorf("Either name or zone_id must be set")
	}

	var nextMarker *string

	var hostedZoneFound *route53.HostedZone
	// We loop through all hostedzone
	for allHostedZoneListed := false; !allHostedZoneListed; {
		req := &route53.ListHostedZonesInput{}
		if nextMarker != nil {
			req.Marker = nextMarker
		}
		log.Printf("[DEBUG] Reading Route53 Zone: %s", req)
		resp, err := conn.ListHostedZones(req)

		if err != nil {
			return fmt.Errorf("Error finding Route 53 Hosted Zone: %w", err)
		}
		for _, hostedZone := range resp.HostedZones {
			hostedZoneId := CleanZoneID(aws.StringValue(hostedZone.Id))
			if idExists && hostedZoneId == id.(string) {
				hostedZoneFound = hostedZone
				break
				// we check if the name is the same as requested and if private zone field is the same as requested or if there is a vpc_id
			} else if (TrimTrailingPeriod(aws.StringValue(hostedZone.Name)) == TrimTrailingPeriod(name)) && (aws.BoolValue(hostedZone.Config.PrivateZone) == d.Get("private_zone").(bool) || (aws.BoolValue(hostedZone.Config.PrivateZone) && vpcIdExists)) {
				matchingVPC := false
				if vpcIdExists {
					reqHostedZone := &route53.GetHostedZoneInput{}
					reqHostedZone.Id = aws.String(hostedZoneId)

					respHostedZone, errHostedZone := conn.GetHostedZone(reqHostedZone)
					if errHostedZone != nil {
						return fmt.Errorf("Error finding Route 53 Hosted Zone: %w", errHostedZone)
					}
					// we go through all VPCs
					for _, vpc := range respHostedZone.VPCs {
						if aws.StringValue(vpc.VPCId) == vpcId.(string) {
							matchingVPC = true
							break
						}
					}
				} else {
					matchingVPC = true
				}
				// we check if tags match
				matchingTags := true
				if len(tags) > 0 {
					listTags, err := ListTags(conn, hostedZoneId, route53.TagResourceTypeHostedzone)

					if err != nil {
						return fmt.Errorf("Error finding Route 53 Hosted Zone: %w", err)
					}
					matchingTags = listTags.ContainsAll(tags)
				}

				if matchingTags && matchingVPC {
					if hostedZoneFound != nil {
						return fmt.Errorf("multiple Route53Zone found please use vpc_id option to filter")
					}

					hostedZoneFound = hostedZone
				}
			}
		}
		if *resp.IsTruncated {
			nextMarker = resp.NextMarker
		} else {
			allHostedZoneListed = true
		}
	}
	if hostedZoneFound == nil {
		return fmt.Errorf("no matching Route53Zone found")
	}

	idHostedZone := CleanZoneID(aws.StringValue(hostedZoneFound.Id))
	d.SetId(idHostedZone)
	d.Set("zone_id", idHostedZone)
	// To be consistent with other AWS services (e.g. ACM) that do not accept a trailing period,
	// we remove the suffix from the Hosted Zone Name returned from the API
	d.Set("name", TrimTrailingPeriod(aws.StringValue(hostedZoneFound.Name)))
	d.Set("comment", hostedZoneFound.Config.Comment)
	d.Set("private_zone", hostedZoneFound.Config.PrivateZone)
	d.Set("caller_reference", hostedZoneFound.CallerReference)
	d.Set("resource_record_set_count", hostedZoneFound.ResourceRecordSetCount)
	if hostedZoneFound.LinkedService != nil {
		d.Set("linked_service_principal", hostedZoneFound.LinkedService.ServicePrincipal)
		d.Set("linked_service_description", hostedZoneFound.LinkedService.Description)
	}

	nameServers, err := hostedZoneNameServers(conn, idHostedZone, aws.StringValue(hostedZoneFound.Name))

	if err != nil {
		return fmt.Errorf("getting Route 53 Hosted Zone (%s) name servers: %w", idHostedZone, err)
	}

	if err := d.Set("name_servers", nameServers); err != nil {
		return fmt.Errorf("setting name_servers: %w", err)
	}

	tags, err = ListTags(conn, idHostedZone, route53.TagResourceTypeHostedzone)

	if err != nil {
		return fmt.Errorf("listing Route 53 Hosted Zone (%s) tags: %w", idHostedZone, err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "route53",
		Resource:  fmt.Sprintf("hostedzone/%s", d.Id()),
	}.String()
	d.Set("arn", arn)

	return nil
}

// used to retrieve name servers
func hostedZoneNameServers(conn *route53.Route53, id string, name string) ([]string, error) {
	input := &route53.GetHostedZoneInput{
		Id: aws.String(id),
	}

	output, err := conn.GetHostedZone(input)

	if err != nil {
		return nil, fmt.Errorf("getting Route 53 Hosted Zone (%s): %w", id, err)
	}

	if output == nil {
		return nil, fmt.Errorf("getting Route 53 Hosted Zone (%s): empty response", id)
	}

	if output.DelegationSet != nil {
		return aws.StringValueSlice(output.DelegationSet.NameServers), nil
	}

	if output.HostedZone != nil && output.HostedZone.Config != nil && aws.BoolValue(output.HostedZone.Config.PrivateZone) {
		nameServers, err := getNameServers(id, name, conn)

		if err != nil {
			return nil, fmt.Errorf("listing Route 53 Hosted Zone (%s) NS records: %w", id, err)
		}

		return nameServers, nil
	}

	return nil, nil
}
