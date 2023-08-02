package ec2

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKDataSource("aws_instance")
func DataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceInstanceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ami": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credit_specification": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_credits": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"disable_api_stop": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_api_termination": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ebs_block_device": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tftags.TagsSchemaComputed(),
						"throughput": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				// This should not be necessary, but currently is (see #7198)
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["device_name"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["snapshot_id"].(string)))
					return create.StringHashcode(buf.String())
				},
			},
			"ebs_optimized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enclave_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"ephemeral_block_device": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"no_device": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"virtual_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"filter": CustomFiltersSchema(),
			"get_password_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"get_user_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_resource_group_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_instance_profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_tags": tftags.TagsSchemaComputed(),
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_addresses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_recovery": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_put_response_hop_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"http_tokens": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_metadata_tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"monitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outpost_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placement_partition_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns_name_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_resource_name_dns_aaaa_record": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_resource_name_dns_a_record": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hostname_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"public_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_block_device": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tftags.TagsSchemaComputed(),
						"throughput": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secondary_private_ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_dest_check": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// dataSourceInstanceRead performs the instanceID lookup
func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	// Build up search parameters
	input := &ec2.DescribeInstancesInput{}

	if tags, tagsOk := d.GetOk("instance_tags"); tagsOk {
		input.Filters = append(input.Filters, BuildTagFilterList(
			Tags(tftags.New(ctx, tags.(map[string]interface{}))),
		)...)
	}

	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	if v, ok := d.GetOk("instance_id"); ok {
		input.InstanceIds = aws.StringSlice([]string{v.(string)})
	}

	instance, err := FindInstance(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("EC2 Instance", err))
	}

	log.Printf("[DEBUG] aws_instance - Single Instance ID found: %s", aws.StringValue(instance.InstanceId))
	if err := instanceDescriptionAttributes(ctx, d, instance, conn, ignoreTagsConfig); err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Instance (%s): %s", aws.StringValue(instance.InstanceId), err)
	}

	if d.Get("get_password_data").(bool) {
		passwordData, err := getInstancePasswordData(ctx, aws.StringValue(instance.InstanceId), conn)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading EC2 Instance (%s): %s", aws.StringValue(instance.InstanceId), err)
		}
		d.Set("password_data", passwordData)
	}

	// ARN
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   ec2.ServiceName,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("instance/%s", d.Id()),
	}
	d.Set("arn", arn.String())

	return diags
}

// Populate instance attribute fields with the returned instance
func instanceDescriptionAttributes(ctx context.Context, d *schema.ResourceData, instance *ec2.Instance, conn *ec2.EC2, ignoreTagsConfig *tftags.IgnoreConfig) error {
	d.SetId(aws.StringValue(instance.InstanceId))

	instanceType := aws.StringValue(instance.InstanceType)
	instanceTypeInfo, err := FindInstanceTypeByName(ctx, conn, instanceType)

	if err != nil {
		return fmt.Errorf("reading EC2 Instance Type (%s): %w", instanceType, err)
	}

	// Set the easy attributes
	d.Set("instance_state", instance.State.Name)
	d.Set("availability_zone", instance.Placement.AvailabilityZone)
	d.Set("placement_group", instance.Placement.GroupName)
	d.Set("placement_partition_number", instance.Placement.PartitionNumber)
	d.Set("tenancy", instance.Placement.Tenancy)
	d.Set("host_id", instance.Placement.HostId)
	d.Set("host_resource_group_arn", instance.Placement.HostResourceGroupArn)

	d.Set("ami", instance.ImageId)
	d.Set("instance_type", instanceType)
	d.Set("key_name", instance.KeyName)
	d.Set("outpost_arn", instance.OutpostArn)
	d.Set("private_dns", instance.PrivateDnsName)
	d.Set("private_ip", instance.PrivateIpAddress)
	d.Set("public_dns", instance.PublicDnsName)
	d.Set("public_ip", instance.PublicIpAddress)

	if instance.IamInstanceProfile != nil && instance.IamInstanceProfile.Arn != nil {
		name, err := InstanceProfileARNToName(aws.StringValue(instance.IamInstanceProfile.Arn))

		if err != nil {
			return fmt.Errorf("setting iam_instance_profile: %w", err)
		}

		d.Set("iam_instance_profile", name)
	} else {
		d.Set("iam_instance_profile", nil)
	}

	// iterate through network interfaces, and set subnet, network_interface, public_addr
	if len(instance.NetworkInterfaces) > 0 {
		for _, ni := range instance.NetworkInterfaces {
			if aws.Int64Value(ni.Attachment.DeviceIndex) == 0 {
				d.Set("subnet_id", ni.SubnetId)
				d.Set("network_interface_id", ni.NetworkInterfaceId)
				d.Set("associate_public_ip_address", ni.Association != nil)

				secondaryIPs := make([]string, 0, len(ni.PrivateIpAddresses))
				for _, ip := range ni.PrivateIpAddresses {
					if !aws.BoolValue(ip.Primary) {
						secondaryIPs = append(secondaryIPs, aws.StringValue(ip.PrivateIpAddress))
					}
				}
				if err := d.Set("secondary_private_ips", secondaryIPs); err != nil {
					return fmt.Errorf("setting secondary_private_ips: %w", err)
				}

				ipV6Addresses := make([]string, 0, len(ni.Ipv6Addresses))
				for _, ip := range ni.Ipv6Addresses {
					ipV6Addresses = append(ipV6Addresses, aws.StringValue(ip.Ipv6Address))
				}
				if err := d.Set("ipv6_addresses", ipV6Addresses); err != nil {
					return fmt.Errorf("setting ipv6_addresses: %w", err)
				}
			}
		}
	} else {
		d.Set("subnet_id", instance.SubnetId)
		d.Set("network_interface_id", "")
	}

	d.Set("ebs_optimized", instance.EbsOptimized)
	if aws.StringValue(instance.SubnetId) != "" {
		d.Set("source_dest_check", instance.SourceDestCheck)
	}

	if instance.Monitoring != nil {
		monitoringState := aws.StringValue(instance.Monitoring.State)
		d.Set("monitoring", monitoringState == "enabled" || monitoringState == "pending")
	}

	if err := d.Set("tags", KeyValueTags(ctx, instance.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	// Security Groups
	if err := readSecurityGroups(ctx, d, instance, conn); err != nil {
		return fmt.Errorf("reading EC2 Instance (%s): %w", aws.StringValue(instance.InstanceId), err)
	}

	// Block devices
	if err := readBlockDevices(ctx, d, instance, conn); err != nil {
		return fmt.Errorf("reading EC2 Instance (%s): %w", aws.StringValue(instance.InstanceId), err)
	}
	if _, ok := d.GetOk("ephemeral_block_device"); !ok {
		d.Set("ephemeral_block_device", []interface{}{})
	}

	// Lookup and Set Instance Attributes
	{
		attr, err := conn.DescribeInstanceAttributeWithContext(ctx, &ec2.DescribeInstanceAttributeInput{
			Attribute:  aws.String(ec2.InstanceAttributeNameDisableApiStop),
			InstanceId: aws.String(d.Id()),
		})
		if err != nil {
			return fmt.Errorf("getting attribute (%s): %w", ec2.InstanceAttributeNameDisableApiStop, err)
		}
		d.Set("disable_api_stop", attr.DisableApiStop.Value)
	}
	{
		attr, err := conn.DescribeInstanceAttributeWithContext(ctx, &ec2.DescribeInstanceAttributeInput{
			Attribute:  aws.String(ec2.InstanceAttributeNameDisableApiTermination),
			InstanceId: aws.String(d.Id()),
		})
		if err != nil {
			return fmt.Errorf("getting attribute (%s): %w", ec2.InstanceAttributeNameDisableApiTermination, err)
		}
		d.Set("disable_api_termination", attr.DisableApiTermination.Value)
	}
	{
		attr, err := conn.DescribeInstanceAttributeWithContext(ctx, &ec2.DescribeInstanceAttributeInput{
			Attribute:  aws.String(ec2.InstanceAttributeNameUserData),
			InstanceId: aws.String(d.Id()),
		})
		if err != nil {
			return fmt.Errorf("getting attribute (%s): %w", ec2.InstanceAttributeNameUserData, err)
		}
		if attr != nil && attr.UserData != nil && attr.UserData.Value != nil {
			d.Set("user_data", userDataHashSum(aws.StringValue(attr.UserData.Value)))
			if d.Get("get_user_data").(bool) {
				d.Set("user_data_base64", attr.UserData.Value)
			}
		}
	}

	// AWS Standard will return InstanceCreditSpecification.NotSupported errors for EC2 Instance IDs outside T2 and T3 instance types
	// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/8055
	if aws.BoolValue(instanceTypeInfo.BurstablePerformanceSupported) {
		instanceCreditSpecification, err := FindInstanceCreditSpecificationByID(ctx, conn, d.Id())

		// Ignore UnsupportedOperation errors for AWS China and GovCloud (US).
		// Reference: https://github.com/hashicorp/terraform-provider-aws/pull/4362.
		if tfawserr.ErrCodeEquals(err, errCodeUnsupportedOperation) {
			err = nil
		}

		if err != nil {
			return fmt.Errorf("reading EC2 Instance (%s) credit specification: %w", d.Id(), err)
		}

		if instanceCreditSpecification != nil {
			if err := d.Set("credit_specification", []interface{}{flattenInstanceCreditSpecification(instanceCreditSpecification)}); err != nil {
				return fmt.Errorf("setting credit_specification: %w", err)
			}
		} else {
			d.Set("credit_specification", nil)
		}
	} else {
		d.Set("credit_specification", nil)
	}

	if err := d.Set("enclave_options", flattenEnclaveOptions(instance.EnclaveOptions)); err != nil {
		return fmt.Errorf("setting enclave_options: %w", err)
	}

	if instance.MaintenanceOptions != nil {
		if err := d.Set("maintenance_options", []interface{}{flattenInstanceMaintenanceOptions(instance.MaintenanceOptions)}); err != nil {
			return fmt.Errorf("setting maintenance_options: %w", err)
		}
	} else {
		d.Set("maintenance_options", nil)
	}

	if err := d.Set("metadata_options", flattenInstanceMetadataOptions(instance.MetadataOptions)); err != nil {
		return fmt.Errorf("setting metadata_options: %w", err)
	}

	if instance.PrivateDnsNameOptions != nil {
		if err := d.Set("private_dns_name_options", []interface{}{flattenPrivateDNSNameOptionsResponse(instance.PrivateDnsNameOptions)}); err != nil {
			return fmt.Errorf("setting private_dns_name_options: %w", err)
		}
	} else {
		d.Set("private_dns_name_options", nil)
	}

	return nil
}
