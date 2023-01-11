package cloudhsmv2

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudhsmv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"cluster_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_csr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aws_hardware_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hsm_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manufacturer_hardware_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).CloudHSMV2Conn

	clusterId := d.Get("cluster_id").(string)
	filters := []*string{&clusterId}
	log.Printf("[DEBUG] Reading CloudHSM v2 Cluster %s", clusterId)
	result := int64(1)
	input := &cloudhsmv2.DescribeClustersInput{
		Filters: map[string][]*string{
			"clusterIds": filters,
		},
		MaxResults: &result,
	}
	state := d.Get("cluster_state").(string)
	states := []*string{&state}
	if len(state) > 0 {
		input.Filters["states"] = states
	}
	out, err := conn.DescribeClusters(input)

	if err != nil {
		return fmt.Errorf("error describing CloudHSM v2 Cluster: %w", err)
	}

	var cluster *cloudhsmv2.Cluster
	for _, c := range out.Clusters {
		if aws.StringValue(c.ClusterId) == clusterId {
			cluster = c
			break
		}
	}

	if cluster == nil {
		return fmt.Errorf("cluster with id %s not found", clusterId)
	}

	d.SetId(clusterId)
	d.Set("vpc_id", cluster.VpcId)
	d.Set("security_group_id", cluster.SecurityGroup)
	d.Set("cluster_state", cluster.State)
	if err := d.Set("cluster_certificates", readClusterCertificates(cluster)); err != nil {
		return fmt.Errorf("error setting cluster_certificates: %w", err)
	}

	var subnets []string
	for _, sn := range cluster.SubnetMapping {
		subnets = append(subnets, *sn)
	}

	if err := d.Set("subnet_ids", subnets); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Subnet IDs to state for CloudHSM v2 Cluster (%s): %w", d.Id(), err)
	}

	return nil
}
