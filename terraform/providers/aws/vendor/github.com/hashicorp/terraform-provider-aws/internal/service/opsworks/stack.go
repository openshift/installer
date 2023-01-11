package opsworks

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	securityGroupsCreatedSleepTime = 30 * time.Second
	securityGroupsDeletedSleepTime = 30 * time.Second
)

func ResourceStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceStackCreate,
		Read:   resourceStackRead,
		Update: resourceStackUpdate,
		Delete: resourceStackDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"agent_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"berkshelf_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultBerkshelfVersion,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_manager_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Chef",
			},
			"configuration_manager_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "11.10",
			},
			"custom_cookbooks_source": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"revision": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssh_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(opsworks.SourceType_Values(), false),
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"custom_json": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
			},
			"default_availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vpc_id"},
			},
			"default_instance_profile_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_os": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Ubuntu 12.04 LTS",
			},
			"default_root_device_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "instance-store",
			},
			"default_ssh_key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"vpc_id"},
			},
			"hostname_theme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Layer_Dependent",
			},
			"manage_berkshelf": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"service_role_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stack_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"use_custom_cookbooks": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"use_opsworks_security_groups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"vpc_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"default_availability_zone"},
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceStackCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).OpsWorksConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("name").(string)
	region := d.Get("region").(string)
	input := &opsworks.CreateStackInput{
		ChefConfiguration: &opsworks.ChefConfiguration{
			ManageBerkshelf: aws.Bool(d.Get("manage_berkshelf").(bool)),
		},
		ConfigurationManager: &opsworks.StackConfigurationManager{
			Name:    aws.String(d.Get("configuration_manager_name").(string)),
			Version: aws.String(d.Get("configuration_manager_version").(string)),
		},
		DefaultInstanceProfileArn: aws.String(d.Get("default_instance_profile_arn").(string)),
		DefaultOs:                 aws.String(d.Get("default_os").(string)),
		HostnameTheme:             aws.String(d.Get("hostname_theme").(string)),
		Name:                      aws.String(name),
		Region:                    aws.String(region),
		ServiceRoleArn:            aws.String(d.Get("service_role_arn").(string)),
		UseCustomCookbooks:        aws.Bool(d.Get("use_custom_cookbooks").(bool)),
		UseOpsworksSecurityGroups: aws.Bool(d.Get("use_opsworks_security_groups").(bool)),
	}

	if v, ok := d.GetOk("agent_version"); ok {
		input.AgentVersion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("color"); ok {
		input.Attributes = aws.StringMap(map[string]string{
			opsworks.StackAttributesKeysColor: v.(string),
		})
	}

	if v, ok := d.GetOk("custom_cookbooks_source"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.CustomCookbooksSource = expandSource(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("custom_json"); ok {
		input.CustomJson = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_availability_zone"); ok {
		input.DefaultAvailabilityZone = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_root_device_type"); ok {
		input.DefaultRootDeviceType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_ssh_key_name"); ok {
		input.DefaultSshKeyName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("default_subnet_id"); ok {
		input.DefaultSubnetId = aws.String(v.(string))
	}

	if d.Get("manage_berkshelf").(bool) {
		input.ChefConfiguration.BerkshelfVersion = aws.String(d.Get("berkshelf_version").(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		input.VpcId = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating OpsWorks Stack: %s", input)
	outputRaw, err := tfresource.RetryWhen(d.Timeout(schema.TimeoutCreate),
		func() (interface{}, error) {
			return conn.CreateStack(input)
		},
		func(err error) (bool, error) {
			// If Terraform is also managing the service IAM role, it may have just been created and not yet be
			// propagated. AWS doesn't provide a machine-readable code for this specific error, so we're forced
			// to do fragile message matching.
			// The full error we're looking for looks something like the following:
			// Service Role Arn: [...] is not yet propagated, please try again in a couple of minutes
			if tfawserr.ErrMessageContains(err, opsworks.ErrCodeValidationException, "not yet propagated") ||
				tfawserr.ErrMessageContains(err, opsworks.ErrCodeValidationException, "not the necessary trust relationship") ||
				tfawserr.ErrMessageContains(err, opsworks.ErrCodeValidationException, "validate IAM role permission") {
				return true, err
			}

			return false, err
		},
	)

	if err != nil {
		return fmt.Errorf("creating OpsWorks Stack (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(outputRaw.(*opsworks.CreateStackOutput).StackId))

	if len(tags) > 0 {
		arn := arn.ARN{
			Partition: meta.(*conns.AWSClient).Partition,
			Service:   opsworks.ServiceName,
			Region:    region,
			AccountID: meta.(*conns.AWSClient).AccountID,
			Resource:  fmt.Sprintf("stack/%s/", d.Id()),
		}.String()

		if err := UpdateTags(conn, arn, nil, tags); err != nil {
			return fmt.Errorf("adding OpsWorks Stack (%s) tags: %w", arn, err)
		}
	}

	if aws.StringValue(input.VpcId) != "" && aws.BoolValue(input.UseOpsworksSecurityGroups) {
		// For VPC-based stacks, OpsWorks asynchronously creates some default
		// security groups which must exist before layers can be created.
		// Unfortunately it doesn't tell us what the ids of these are, so
		// we can't actually check for them. Instead, we just wait a nominal
		// amount of time for their creation to complete.
		log.Print("[INFO] Waiting for OpsWorks built-in security groups to be created")
		time.Sleep(securityGroupsCreatedSleepTime)
	}

	return resourceStackRead(d, meta)
}

func resourceStackRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	conn := meta.(*conns.AWSClient).OpsWorksConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	if v, ok := d.GetOk("stack_endpoint"); ok {
		conn, err = regionalConn(meta.(*conns.AWSClient), v.(string))

		if err != nil {
			return err
		}
	}

	stack, err := FindStackByID(conn, d.Id())

	if tfresource.NotFound(err) {
		// If it's not found in the the default region we're in, we check us-east-1
		// in the event this stack was created with Terraform before version 0.9.
		// See https://github.com/hashicorp/terraform/issues/12842.
		conn, err = regionalConn(meta.(*conns.AWSClient), endpoints.UsEast1RegionID)

		if err != nil {
			return err
		}

		stack, err = FindStackByID(conn, d.Id())
	}

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] OpsWorks Stack %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading OpsWorks Stack (%s): %w", d.Id(), err)
	}

	// If the stack was found, set the stack_endpoint.
	if v := aws.StringValue(conn.Config.Region); v != "" {
		d.Set("stack_endpoint", v)
	}

	d.Set("agent_version", stack.AgentVersion)
	arn := aws.StringValue(stack.Arn)
	d.Set("arn", arn)
	if stack.ChefConfiguration != nil {
		if v := aws.StringValue(stack.ChefConfiguration.BerkshelfVersion); v != "" {
			d.Set("berkshelf_version", v)
		} else {
			d.Set("berkshelf_version", defaultBerkshelfVersion)
		}
		d.Set("manage_berkshelf", stack.ChefConfiguration.ManageBerkshelf)
	}
	if color, ok := stack.Attributes[opsworks.StackAttributesKeysColor]; ok {
		d.Set("color", color)
	}
	if stack.ConfigurationManager != nil {
		d.Set("configuration_manager_name", stack.ConfigurationManager.Name)
		d.Set("configuration_manager_version", stack.ConfigurationManager.Version)
	}
	if stack.CustomCookbooksSource != nil {
		tfMap := flattenSource(stack.CustomCookbooksSource)

		// CustomCookbooksSource.Password and CustomCookbooksSource.SshKey will, on read, contain the placeholder string "*****FILTERED*****",
		// so we ignore it on read and let persist the value already in the state.
		if v, ok := d.GetOk("custom_cookbooks_source"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			v := v.([]interface{})[0].(map[string]interface{})

			tfMap["password"] = v["password"]
			tfMap["ssh_key"] = v["ssh_key"]
		}

		if err := d.Set("custom_cookbooks_source", []interface{}{tfMap}); err != nil {
			return fmt.Errorf("setting custom_cookbooks_source: %w", err)
		}
	} else {
		d.Set("custom_cookbooks_source", nil)
	}
	if stack.CustomJson != nil {
		d.Set("custom_json", stack.CustomJson)
	}
	d.Set("default_availability_zone", stack.DefaultAvailabilityZone)
	d.Set("default_instance_profile_arn", stack.DefaultInstanceProfileArn)
	d.Set("default_os", stack.DefaultOs)
	d.Set("default_root_device_type", stack.DefaultRootDeviceType)
	d.Set("default_ssh_key_name", stack.DefaultSshKeyName)
	d.Set("default_subnet_id", stack.DefaultSubnetId)
	d.Set("hostname_theme", stack.HostnameTheme)
	d.Set("name", stack.Name)
	d.Set("region", stack.Region)
	d.Set("service_role_arn", stack.ServiceRoleArn)
	d.Set("use_custom_cookbooks", stack.UseCustomCookbooks)
	d.Set("use_opsworks_security_groups", stack.UseOpsworksSecurityGroups)
	d.Set("vpc_id", stack.VpcId)

	tags, err := ListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("listing tags for OpsWorks Stack (%s): %w", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

func resourceStackUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	conn := meta.(*conns.AWSClient).OpsWorksConn

	if v, ok := d.GetOk("stack_endpoint"); ok {
		conn, err = regionalConn(meta.(*conns.AWSClient), v.(string))

		if err != nil {
			return err
		}
	}

	if d.HasChangesExcept("tags", "tags_all") {
		input := &opsworks.UpdateStackInput{
			StackId: aws.String(d.Id()),
		}

		if d.HasChange("agent_version") {
			input.AgentVersion = aws.String(d.Get("agent_version").(string))
		}

		if d.HasChanges("berkshelf_version", "manage_berkshelf") {
			input.ChefConfiguration = &opsworks.ChefConfiguration{
				ManageBerkshelf: aws.Bool(d.Get("manage_berkshelf").(bool)),
			}

			if d.Get("manage_berkshelf").(bool) {
				input.ChefConfiguration.BerkshelfVersion = aws.String(d.Get("berkshelf_version").(string))
			}
		}

		if d.HasChange("color") {
			input.Attributes = aws.StringMap(map[string]string{
				opsworks.StackAttributesKeysColor: d.Get("color").(string),
			})
		}

		if d.HasChanges("configuration_manager_name", "configuration_manager_version") {
			input.ConfigurationManager = &opsworks.StackConfigurationManager{
				Name:    aws.String(d.Get("configuration_manager_name").(string)),
				Version: aws.String(d.Get("configuration_manager_version").(string)),
			}
		}

		if d.HasChange("custom_json") {
			input.CustomJson = aws.String(d.Get("custom_json").(string))
		}

		if d.HasChange("custom_cookbooks_source") {
			if v, ok := d.GetOk("custom_cookbooks_source"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.CustomCookbooksSource = expandSource(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("default_availability_zone") {
			input.DefaultAvailabilityZone = aws.String(d.Get("default_availability_zone").(string))
		}

		if d.HasChange("default_instance_profile_arn") {
			input.DefaultInstanceProfileArn = aws.String(d.Get("default_instance_profile_arn").(string))
		}

		if d.HasChange("default_os") {
			input.DefaultOs = aws.String(d.Get("default_os").(string))
		}

		if d.HasChange("default_root_device_type") {
			input.DefaultRootDeviceType = aws.String(d.Get("default_root_device_type").(string))
		}

		if d.HasChange("default_ssh_key_name") {
			input.DefaultSshKeyName = aws.String(d.Get("default_ssh_key_name").(string))
		}

		if d.HasChange("default_subnet_id") {
			input.DefaultSubnetId = aws.String(d.Get("default_subnet_id").(string))
		}

		if d.HasChange("hostname_theme") {
			input.HostnameTheme = aws.String(d.Get("hostname_theme").(string))
		}

		if d.HasChange("name") {
			input.Name = aws.String(d.Get("name").(string))
		}

		if d.HasChange("service_role_arn") {
			input.ServiceRoleArn = aws.String(d.Get("service_role_arn").(string))
		}

		if d.HasChange("use_custom_cookbooks") {
			input.UseCustomCookbooks = aws.Bool(d.Get("use_custom_cookbooks").(bool))
		}

		if d.HasChange("use_opsworks_security_groups") {
			input.UseOpsworksSecurityGroups = aws.Bool(d.Get("use_opsworks_security_groups").(bool))
		}

		log.Printf("[DEBUG] Updating OpsWorks Stack: %s", input)
		_, err = conn.UpdateStack(input)

		if err != nil {
			return fmt.Errorf("updating OpsWorks Stack (%s): %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("updating OpsWorks Stack (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceStackRead(d, meta)
}

func resourceStackDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	conn := meta.(*conns.AWSClient).OpsWorksConn

	if v, ok := d.GetOk("stack_endpoint"); ok {
		conn, err = regionalConn(meta.(*conns.AWSClient), v.(string))

		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Deleting OpsWorks Stack: %s", d.Id())
	_, err = conn.DeleteStack(&opsworks.DeleteStackInput{
		StackId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, opsworks.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting OpsWork Stack (%s): %w", d.Id(), err)
	}

	// For a stack in a VPC, OpsWorks has created some default security groups
	// in the VPC, which it will now delete.
	// Unfortunately, the security groups are deleted asynchronously and there
	// is no robust way for us to determine when it is done. The VPC itself
	// isn't deletable until the security groups are cleaned up, so this could
	// make 'terraform destroy' fail if the VPC is also managed and we don't
	// wait for the security groups to be deleted.
	// There is no robust way to check for this, so we'll just wait a
	// nominal amount of time.
	if _, ok := d.GetOk("vpc_id"); ok {
		if _, ok := d.GetOk("use_opsworks_security_groups"); ok {
			log.Print("[INFO] Waiting for Opsworks built-in security groups to be deleted")
			time.Sleep(securityGroupsDeletedSleepTime)
		}
	}

	return nil
}

func FindStackByID(conn *opsworks.OpsWorks, id string) (*opsworks.Stack, error) {
	input := &opsworks.DescribeStacksInput{
		StackIds: aws.StringSlice([]string{id}),
	}

	output, err := conn.DescribeStacks(input)

	if tfawserr.ErrCodeEquals(err, opsworks.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Stacks) == 0 || output.Stacks[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.Stacks); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output.Stacks[0], nil
}

func expandSource(tfMap map[string]interface{}) *opsworks.Source {
	if tfMap == nil {
		return nil
	}

	apiObject := &opsworks.Source{}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		apiObject.Password = aws.String(v)
	}

	if v, ok := tfMap["revision"].(string); ok && v != "" {
		apiObject.Revision = aws.String(v)
	}

	if v, ok := tfMap["ssh_key"].(string); ok && v != "" {
		apiObject.SshKey = aws.String(v)
	}

	if v, ok := tfMap["type"].(string); ok && v != "" {
		apiObject.Type = aws.String(v)
	}

	if v, ok := tfMap["url"].(string); ok && v != "" {
		apiObject.Url = aws.String(v)
	}

	if v, ok := tfMap["username"].(string); ok && v != "" {
		apiObject.Username = aws.String(v)
	}

	return apiObject
}

func flattenSource(apiObject *opsworks.Source) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Password; v != nil {
		tfMap["password"] = aws.StringValue(v)
	}

	if v := apiObject.Revision; v != nil {
		tfMap["revision"] = aws.StringValue(v)
	}

	if v := apiObject.SshKey; v != nil {
		tfMap["ssh_key"] = aws.StringValue(v)
	}

	if v := apiObject.Type; v != nil {
		tfMap["type"] = aws.StringValue(v)
	}

	if v := apiObject.Url; v != nil {
		tfMap["url"] = aws.StringValue(v)
	}

	if v := apiObject.Username; v != nil {
		tfMap["username"] = aws.StringValue(v)
	}

	return tfMap
}

// opsworksConn will return a connection for the stack_endpoint in the
// configuration. Stacks can only be accessed or managed within the endpoint
// in which they are created, so we allow users to specify an original endpoint
// for Stacks created before multiple endpoints were offered (Terraform v0.9.0).
// See:
//   - https://github.com/hashicorp/terraform/pull/12688
//   - https://github.com/hashicorp/terraform/issues/12842
func regionalConn(client *conns.AWSClient, regionName string) (*opsworks.OpsWorks, error) {
	conn := client.OpsWorksConn

	// Regions are the same, no need to reconfigure.
	if aws.StringValue(conn.Config.Region) == regionName {
		return conn, nil
	}

	sess, err := conns.NewSessionForRegion(&conn.Config, regionName, client.TerraformVersion)

	if err != nil {
		return nil, fmt.Errorf("creating AWS session (%s): %w", regionName, err)
	}

	return opsworks.New(sess), nil
}
