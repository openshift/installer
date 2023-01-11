package redshiftserverless

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceEndpointAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointAccessCreate,
		Read:   resourceEndpointAccessRead,
		Update: resourceEndpointAccessUpdate,
		Delete: resourceEndpointAccessDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_endpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"workgroup_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"endpoint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 30),
			},
		},
	}
}

func resourceEndpointAccessCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn

	input := redshiftserverless.CreateEndpointAccessInput{
		WorkgroupName: aws.String(d.Get("workgroup_name").(string)),
		EndpointName:  aws.String(d.Get("endpoint_name").(string)),
	}

	if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
		input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("subnet_ids"); ok && v.(*schema.Set).Len() > 0 {
		input.SubnetIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	out, err := conn.CreateEndpointAccess(&input)

	if err != nil {
		return fmt.Errorf("error creating Redshift Serverless Endpoint Access: %w", err)
	}

	d.SetId(aws.StringValue(out.Endpoint.EndpointName))

	if _, err := waitEndpointAccessActive(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Redshift Serverless Endpoint Access (%s) to be created: %w", d.Id(), err)
	}

	return resourceEndpointAccessRead(d, meta)
}

func resourceEndpointAccessRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn

	out, err := FindEndpointAccessByName(conn, d.Id())
	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Redshift Serverless EndpointAccess (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Redshift Serverless Endpoint Access (%s): %w", d.Id(), err)
	}

	d.Set("address", out.Address)
	d.Set("port", out.Port)
	d.Set("arn", out.EndpointArn)
	d.Set("endpoint_name", out.EndpointName)
	d.Set("workgroup_name", out.WorkgroupName)
	d.Set("subnet_ids", flex.FlattenStringSet(out.SubnetIds))

	result := make([]*string, 0, len(out.VpcSecurityGroups))

	for _, v := range out.VpcSecurityGroups {
		result = append(result, v.VpcSecurityGroupId)
	}
	d.Set("vpc_security_group_ids", flex.FlattenStringSet(result))

	if err := d.Set("vpc_endpoint", []interface{}{flattenVPCEndpoint(out.VpcEndpoint)}); err != nil {
		return fmt.Errorf("setting vpc_endpoint: %w", err)
	}

	return nil
}

func resourceEndpointAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn

	input := &redshiftserverless.UpdateEndpointAccessInput{
		EndpointName: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("vpc_security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
		input.VpcSecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	_, err := conn.UpdateEndpointAccess(input)
	if err != nil {
		return fmt.Errorf("error updating Redshift Serverless Endpoint Access (%s): %w", d.Id(), err)
	}

	if _, err := waitEndpointAccessActive(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Redshift Serverless Endpoint Access (%s) to be updated: %w", d.Id(), err)
	}

	return resourceEndpointAccessRead(d, meta)
}

func resourceEndpointAccessDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).RedshiftServerlessConn

	deleteInput := redshiftserverless.DeleteEndpointAccessInput{
		EndpointName: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting Redshift Serverless EndpointAccess: %s", d.Id())
	_, err := conn.DeleteEndpointAccess(&deleteInput)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, redshiftserverless.ErrCodeResourceNotFoundException) {
			return nil
		}
		return err
	}

	if _, err := waitEndpointAccessDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Redshift Serverless Endpoint Access (%s) delete: %w", d.Id(), err)
	}

	return nil
}
