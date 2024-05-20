package openstack

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

func dataSourceComputeAggregateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeAggregateV2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},

			"hosts": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceComputeAggregateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	allPages, err := aggregates.List(computeClient).AllPages()
	if err != nil {
		return diag.Errorf("Error listing compute aggregates: %s", err)
	}

	allAggregates, err := aggregates.ExtractAggregates(allPages)
	if err != nil {
		return diag.Errorf("Error extracting compute aggregates: %s", err)
	}

	name := d.Get("name").(string)

	var refinedAggregates []aggregates.Aggregate
	for _, aggregate := range allAggregates {
		if aggregate.Name == name {
			refinedAggregates = append(refinedAggregates, aggregate)
		}
	}

	if len(refinedAggregates) < 1 {
		return diag.Errorf("Could not find any host aggregate with this name: %s", name)
	}
	if len(refinedAggregates) > 1 {
		return diag.Errorf("More than one object found with this name: %s", name)
	}

	aggr := refinedAggregates[0]

	// Metadata is redundant with Availability Zone
	metadata := aggr.Metadata
	_, ok := metadata["availability_zone"]
	if ok {
		delete(metadata, "availability_zone")
	}

	idStr := strconv.Itoa(aggr.ID)
	d.SetId(idStr)
	d.Set("name", aggr.Name)
	d.Set("zone", aggr.AvailabilityZone)
	d.Set("hosts", aggr.Hosts)
	d.Set("metadata", metadata)

	return nil
}
