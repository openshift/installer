package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func dataSourceAwsRoute53Zone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsRoute53ZoneRead,

		Schema: map[string]*schema.Schema{
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
			"tags": tagsSchemaComputed(),
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

func dataSourceAwsRoute53ZoneRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).r53conn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	name, nameExists := d.GetOk("name")
	name = name.(string)
	id, idExists := d.GetOk("zone_id")
	vpcId, vpcIdExists := d.GetOk("vpc_id")
	tags := keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws()

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
			return fmt.Errorf("Error finding Route 53 Hosted Zone: %v", err)
		}
		for _, hostedZone := range resp.HostedZones {
			hostedZoneId := cleanZoneID(aws.StringValue(hostedZone.Id))
			if idExists && hostedZoneId == id.(string) {
				hostedZoneFound = hostedZone
				break
				// we check if the name is the same as requested and if private zone field is the same as requested or if there is a vpc_id
			} else if (trimTrailingPeriod(aws.StringValue(hostedZone.Name)) == trimTrailingPeriod(name)) && (aws.BoolValue(hostedZone.Config.PrivateZone) == d.Get("private_zone").(bool) || (aws.BoolValue(hostedZone.Config.PrivateZone) && vpcIdExists)) {
				matchingVPC := false
				if vpcIdExists {
					reqHostedZone := &route53.GetHostedZoneInput{}
					reqHostedZone.Id = aws.String(hostedZoneId)

					respHostedZone, errHostedZone := conn.GetHostedZone(reqHostedZone)
					if errHostedZone != nil {
						return fmt.Errorf("Error finding Route 53 Hosted Zone: %v", errHostedZone)
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
					listTags, err := keyvaluetags.Route53ListTags(conn, hostedZoneId, route53.TagResourceTypeHostedzone)

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

	idHostedZone := cleanZoneID(aws.StringValue(hostedZoneFound.Id))
	d.SetId(idHostedZone)
	d.Set("zone_id", idHostedZone)
	// To be consistent with other AWS services (e.g. ACM) that do not accept a trailing period,
	// we remove the suffix from the Hosted Zone Name returned from the API
	d.Set("name", trimTrailingPeriod(aws.StringValue(hostedZoneFound.Name)))
	d.Set("comment", hostedZoneFound.Config.Comment)
	d.Set("private_zone", hostedZoneFound.Config.PrivateZone)
	d.Set("caller_reference", hostedZoneFound.CallerReference)
	d.Set("resource_record_set_count", hostedZoneFound.ResourceRecordSetCount)
	if hostedZoneFound.LinkedService != nil {
		d.Set("linked_service_principal", hostedZoneFound.LinkedService.ServicePrincipal)
		d.Set("linked_service_description", hostedZoneFound.LinkedService.Description)
	}

	nameServers, err := hostedZoneNameServers(idHostedZone, conn)
	if err != nil {
		return fmt.Errorf("Error finding Route 53 Hosted Zone: %v", err)
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return fmt.Errorf("error setting name_servers: %w", err)
	}

	tags, err = keyvaluetags.Route53ListTags(conn, idHostedZone, route53.TagResourceTypeHostedzone)

	if err != nil {
		return fmt.Errorf("Error finding Route 53 Hosted Zone: %v", err)
	}

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

// used to retrieve name servers
func hostedZoneNameServers(id string, conn *route53.Route53) ([]string, error) {
	req := &route53.GetHostedZoneInput{}
	req.Id = aws.String(id)

	resp, err := conn.GetHostedZone(req)
	if err != nil {
		return []string{}, err
	}

	if resp.DelegationSet == nil {
		return []string{}, nil
	}

	servers := []string{}
	for _, server := range resp.DelegationSet.NameServers {
		if server != nil {
			servers = append(servers, aws.StringValue(server))
		}
	}
	return servers, nil
}
