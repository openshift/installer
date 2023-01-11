package mq

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mq"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/mitchellh/copystructure"
)

func ResourceBroker() *schema.Resource {
	return &schema.Resource{
		Create: resourceBrokerCreate,
		Read:   resourceBrokerRead,
		Update: resourceBrokerUpdate,
		Delete: resourceBrokerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"apply_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(mq.AuthenticationStrategy_Values(), true),
			},
			"auto_minor_version_upgrade": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"broker_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateBrokerName,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"revision": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"deployment_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      mq.DeploymentModeSingleInstance,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(mq.DeploymentMode_Values(), true),
			},
			"encryption_options": {
				Type:             schema.TypeList,
				Optional:         true,
				ForceNew:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"use_aws_owned_key": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
					},
				},
			},
			"engine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(mq.EngineType_Values(), true),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"console_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ldap_server_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hosts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"role_base": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_search_matching": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role_search_subtree": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"service_account_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"service_account_username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_base": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_role_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_search_matching": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_search_subtree": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"logs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// Ignore missing configuration block
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "1" && new == "0" {
						return true
					}
					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"general": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"audit": {
							Type:             nullable.TypeNullableBool,
							Optional:         true,
							ValidateFunc:     nullable.ValidateTypeStringNullableBool,
							DiffSuppressFunc: nullable.DiffSuppressNullableBoolFalseAsNull,
						},
					},
				},
			},
			"maintenance_window_start_time": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(mq.DayOfWeek_Values(), true),
						},
						"time_of_day": {
							Type:     schema.TypeString,
							Required: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"publicly_accessible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 5,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(mq.BrokerStorageType_Values(), true),
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"user": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceUserHash,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// AWS currently does not support updating the RabbitMQ users beyond resource creation.
					// User list is not returned back after creation.
					// Updates to users can only be in the RabbitMQ UI.
					if v := d.Get("engine_type").(string); strings.EqualFold(v, mq.EngineTypeRabbitmq) && d.Get("arn").(string) != "" {
						return true
					}

					return false
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"console_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"groups": {
							Type:     schema.TypeSet,
							MaxItems: 20,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(2, 100),
							},
							Optional: true,
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: ValidBrokerPassword,
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(2, 100),
						},
					},
				},
			},
		},

		CustomizeDiff: customdiff.All(
			verify.SetTagsDiff,
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				if strings.EqualFold(diff.Get("engine_type").(string), mq.EngineTypeRabbitmq) {
					if v, ok := diff.GetOk("logs.0.audit"); ok {
						if v, _, _ := nullable.Bool(v.(string)).Value(); v {
							return errors.New("logs.audit: Can not be configured when engine is RabbitMQ")
						}
					}
				}

				return nil
			},
		),
	}
}

func resourceBrokerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).MQConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	name := d.Get("broker_name").(string)
	engineType := d.Get("engine_type").(string)
	requestId := resource.PrefixedUniqueId(fmt.Sprintf("tf-%s", name))
	input := mq.CreateBrokerRequest{
		AutoMinorVersionUpgrade: aws.Bool(d.Get("auto_minor_version_upgrade").(bool)),
		BrokerName:              aws.String(name),
		CreatorRequestId:        aws.String(requestId),
		EngineType:              aws.String(engineType),
		EngineVersion:           aws.String(d.Get("engine_version").(string)),
		HostInstanceType:        aws.String(d.Get("host_instance_type").(string)),
		PubliclyAccessible:      aws.Bool(d.Get("publicly_accessible").(bool)),
		Users:                   expandUsers(d.Get("user").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("authentication_strategy"); ok {
		input.AuthenticationStrategy = aws.String(v.(string))
	}
	if v, ok := d.GetOk("configuration"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Configuration = expandConfigurationId(v.([]interface{}))
	}
	if v, ok := d.GetOk("deployment_mode"); ok {
		input.DeploymentMode = aws.String(v.(string))
	}
	if v, ok := d.GetOk("encryption_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.EncryptionOptions = expandEncryptionOptions(d.Get("encryption_options").([]interface{}))
	}
	if v, ok := d.GetOk("logs"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Logs = expandLogs(engineType, v.([]interface{}))
	}
	if v, ok := d.GetOk("ldap_server_metadata"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.LdapServerMetadata = expandLDAPServerMetadata(v.([]interface{}))
	}
	if v, ok := d.GetOk("maintenance_window_start_time"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.MaintenanceWindowStartTime = expandWeeklyStartTime(v.([]interface{}))
	}
	if v, ok := d.GetOk("security_groups"); ok && v.(*schema.Set).Len() > 0 {
		input.SecurityGroups = flex.ExpandStringSet(v.(*schema.Set))
	}
	if v, ok := d.GetOk("storage_type"); ok {
		input.StorageType = aws.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_ids"); ok {
		input.SubnetIds = flex.ExpandStringSet(v.(*schema.Set))
	}
	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[INFO] Creating MQ Broker: %s", input)
	out, err := conn.CreateBroker(&input)
	if err != nil {
		return err
	}

	d.SetId(aws.StringValue(out.BrokerId))
	d.Set("arn", out.BrokerArn)

	if _, err := WaitBrokerCreated(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for MQ Broker (%s) creation: %w", d.Id(), err)
	}

	return resourceBrokerRead(d, meta)
}

func resourceBrokerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).MQConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	output, err := conn.DescribeBroker(&mq.DescribeBrokerInput{
		BrokerId: aws.String(d.Id()),
	})

	if !d.IsNewResource() && (tfawserr.ErrCodeEquals(err, mq.ErrCodeNotFoundException) || tfawserr.ErrCodeEquals(err, mq.ErrCodeForbiddenException)) {
		log.Printf("[WARN] MQ broker (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading MQ broker (%s): %w", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("empty response while reading MQ broker (%s)", d.Id())
	}

	d.Set("arn", output.BrokerArn)
	d.Set("authentication_strategy", output.AuthenticationStrategy)
	d.Set("auto_minor_version_upgrade", output.AutoMinorVersionUpgrade)
	d.Set("broker_name", output.BrokerName)
	d.Set("deployment_mode", output.DeploymentMode)
	d.Set("engine_type", output.EngineType)
	d.Set("engine_version", output.EngineVersion)
	d.Set("host_instance_type", output.HostInstanceType)
	d.Set("instances", flattenBrokerInstances(output.BrokerInstances))
	d.Set("publicly_accessible", output.PubliclyAccessible)
	d.Set("security_groups", aws.StringValueSlice(output.SecurityGroups))
	d.Set("storage_type", output.StorageType)
	d.Set("subnet_ids", aws.StringValueSlice(output.SubnetIds))

	if err := d.Set("configuration", flattenConfiguration(output.Configurations)); err != nil {
		return fmt.Errorf("error setting configuration: %w", err)
	}

	if err := d.Set("encryption_options", flattenEncryptionOptions(output.EncryptionOptions)); err != nil {
		return fmt.Errorf("error setting encryption_options: %w", err)
	}

	var password string
	if v, ok := d.GetOk("ldap_server_metadata.0.service_account_password"); ok {
		password = v.(string)
	}

	if err := d.Set("ldap_server_metadata", flattenLDAPServerMetadata(output.LdapServerMetadata, password)); err != nil {
		return fmt.Errorf("error setting ldap_server_metadata: %w", err)
	}

	if err := d.Set("logs", flattenLogs(output.Logs)); err != nil {
		return fmt.Errorf("error setting logs: %w", err)
	}

	if err := d.Set("maintenance_window_start_time", flattenWeeklyStartTime(output.MaintenanceWindowStartTime)); err != nil {
		return fmt.Errorf("error setting maintenance_window_start_time: %w", err)
	}

	rawUsers, err := expandUsersForBroker(conn, d.Id(), output.Users)

	if err != nil {
		return fmt.Errorf("error retrieving user info for MQ broker (%s): %w", d.Id(), err)
	}

	if err := d.Set("user", flattenUsers(rawUsers, d.Get("user").(*schema.Set).List())); err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	tags := KeyValueTags(output.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceBrokerUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).MQConn

	requiresReboot := false

	if d.HasChange("security_groups") {
		_, err := conn.UpdateBroker(&mq.UpdateBrokerRequest{
			BrokerId:       aws.String(d.Id()),
			SecurityGroups: flex.ExpandStringSet(d.Get("security_groups").(*schema.Set)),
		})
		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) security groups: %w", d.Id(), err)
		}
	}

	if d.HasChanges("configuration", "logs", "engine_version") {
		engineType := d.Get("engine_type").(string)
		_, err := conn.UpdateBroker(&mq.UpdateBrokerRequest{
			BrokerId:      aws.String(d.Id()),
			Configuration: expandConfigurationId(d.Get("configuration").([]interface{})),
			Logs:          expandLogs(engineType, d.Get("logs").([]interface{})),
			EngineVersion: aws.String(d.Get("engine_version").(string)),
		})
		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) configuration: %w", d.Id(), err)
		}
		requiresReboot = true
	}

	if d.HasChange("user") {
		o, n := d.GetChange("user")
		var err error
		// d.HasChange("user") always reports a change when running resourceBrokerUpdate
		// updateBrokerUsers needs to be called to know if changes to user are actually made
		var usersUpdated bool
		usersUpdated, err = updateBrokerUsers(conn, d.Id(),
			o.(*schema.Set).List(), n.(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) user: %w", d.Id(), err)
		}

		if usersUpdated {
			requiresReboot = true
		}
	}

	if d.HasChange("host_instance_type") {
		_, err := conn.UpdateBroker(&mq.UpdateBrokerRequest{
			BrokerId:         aws.String(d.Id()),
			HostInstanceType: aws.String(d.Get("host_instance_type").(string)),
		})

		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) host instance type: %w", d.Id(), err)
		}

		requiresReboot = true
	}

	if d.HasChange("auto_minor_version_upgrade") {
		_, err := conn.UpdateBroker(&mq.UpdateBrokerRequest{
			BrokerId:                aws.String(d.Id()),
			AutoMinorVersionUpgrade: aws.Bool(d.Get("auto_minor_version_upgrade").(bool)),
		})

		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) auto minor version upgrade: %w", d.Id(), err)
		}

		requiresReboot = true
	}

	if d.HasChange("maintenance_window_start_time") {
		_, err := conn.UpdateBroker(&mq.UpdateBrokerRequest{
			BrokerId:                   aws.String(d.Id()),
			MaintenanceWindowStartTime: expandWeeklyStartTime(d.Get("maintenance_window_start_time").([]interface{})),
		})

		if err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) maintenance window start time: %w", d.Id(), err)
		}

		requiresReboot = true
	}

	if d.Get("apply_immediately").(bool) && requiresReboot {
		_, err := conn.RebootBroker(&mq.RebootBrokerInput{
			BrokerId: aws.String(d.Id()),
		})
		if err != nil {
			return fmt.Errorf("error rebooting MQ Broker (%s): %w", d.Id(), err)
		}

		if _, err := WaitBrokerRebooted(conn, d.Id()); err != nil {
			return fmt.Errorf("error waiting for MQ Broker (%s) reboot: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating MQ Broker (%s) tags: %w", d.Get("arn").(string), err)
		}
	}

	return nil
}

func resourceBrokerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).MQConn

	log.Printf("[INFO] Deleting MQ Broker: %s", d.Id())
	_, err := conn.DeleteBroker(&mq.DeleteBrokerInput{
		BrokerId: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}

	if _, err := WaitBrokerDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for MQ Broker (%s) deletion: %w", d.Id(), err)
	}

	return nil
}

func resourceUserHash(v interface{}) int {
	var buf bytes.Buffer

	m := v.(map[string]interface{})
	if ca, ok := m["console_access"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", ca.(bool)))
	} else {
		buf.WriteString("false-")
	}
	if g, ok := m["groups"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", g.(*schema.Set).List()))
	}
	if p, ok := m["password"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", p.(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["username"].(string)))

	return create.StringHashcode(buf.String())
}

func updateBrokerUsers(conn *mq.MQ, bId string, oldUsers, newUsers []interface{}) (bool, error) {
	// If there are any user creates/deletes/updates, updatedUsers will be set to true
	updatedUsers := false

	createL, deleteL, updateL, err := DiffBrokerUsers(bId, oldUsers, newUsers)
	if err != nil {
		return updatedUsers, err
	}

	for _, c := range createL {
		_, err := conn.CreateUser(c)
		updatedUsers = true
		if err != nil {
			return updatedUsers, err
		}
	}
	for _, d := range deleteL {
		_, err := conn.DeleteUser(d)
		updatedUsers = true
		if err != nil {
			return updatedUsers, err
		}
	}
	for _, u := range updateL {
		_, err := conn.UpdateUser(u)
		updatedUsers = true
		if err != nil {
			return updatedUsers, err
		}
	}

	return updatedUsers, nil
}

func DiffBrokerUsers(bId string, oldUsers, newUsers []interface{}) (
	cr []*mq.CreateUserRequest, di []*mq.DeleteUserInput, ur []*mq.UpdateUserRequest, e error) {

	existingUsers := make(map[string]interface{})
	for _, ou := range oldUsers {
		u := ou.(map[string]interface{})
		username := u["username"].(string)
		// Convert Set to slice to allow easier comparison
		if g, ok := u["groups"]; ok {
			groups := g.(*schema.Set).List()
			u["groups"] = groups
		}

		existingUsers[username] = u
	}

	for _, nu := range newUsers {
		// Still need access to the original map
		// because Set contents doesn't get copied
		// Likely related to https://github.com/mitchellh/copystructure/issues/17
		nuOriginal := nu.(map[string]interface{})

		// Create a mutable copy
		newUser, err := copystructure.Copy(nu)
		if err != nil {
			return cr, di, ur, err
		}

		newUserMap := newUser.(map[string]interface{})
		username := newUserMap["username"].(string)

		// Convert Set to slice to allow easier comparison
		var ng []interface{}
		if g, ok := nuOriginal["groups"]; ok {
			ng = g.(*schema.Set).List()
			newUserMap["groups"] = ng
		}

		if eu, ok := existingUsers[username]; ok {

			existingUserMap := eu.(map[string]interface{})

			if !reflect.DeepEqual(existingUserMap, newUserMap) {
				ur = append(ur, &mq.UpdateUserRequest{
					BrokerId:      aws.String(bId),
					ConsoleAccess: aws.Bool(newUserMap["console_access"].(bool)),
					Groups:        flex.ExpandStringList(ng),
					Password:      aws.String(newUserMap["password"].(string)),
					Username:      aws.String(username),
				})
			}

			// Delete after processing, so we know what's left for deletion
			delete(existingUsers, username)
		} else {
			cur := &mq.CreateUserRequest{
				BrokerId:      aws.String(bId),
				ConsoleAccess: aws.Bool(newUserMap["console_access"].(bool)),
				Password:      aws.String(newUserMap["password"].(string)),
				Username:      aws.String(username),
			}
			if len(ng) > 0 {
				cur.Groups = flex.ExpandStringList(ng)
			}
			cr = append(cr, cur)
		}
	}

	for username := range existingUsers {
		di = append(di, &mq.DeleteUserInput{
			BrokerId: aws.String(bId),
			Username: aws.String(username),
		})
	}

	return cr, di, ur, nil
}

func expandEncryptionOptions(l []interface{}) *mq.EncryptionOptions {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	encryptionOptions := &mq.EncryptionOptions{
		UseAwsOwnedKey: aws.Bool(m["use_aws_owned_key"].(bool)),
	}

	if v, ok := m["kms_key_id"].(string); ok && v != "" {
		encryptionOptions.KmsKeyId = aws.String(v)
	}

	return encryptionOptions
}

func flattenEncryptionOptions(encryptionOptions *mq.EncryptionOptions) []interface{} {
	if encryptionOptions == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"kms_key_id":        aws.StringValue(encryptionOptions.KmsKeyId),
		"use_aws_owned_key": aws.BoolValue(encryptionOptions.UseAwsOwnedKey),
	}

	return []interface{}{m}
}

func ValidBrokerPassword(v interface{}, k string) (ws []string, errors []error) {
	min := 12
	max := 250
	value := v.(string)
	unique := make(map[string]bool)

	for _, v := range value {
		if _, ok := unique[string(v)]; ok {
			continue
		}
		if string(v) == "," {
			errors = append(errors, fmt.Errorf("%q must not contain commas", k))
		}
		unique[string(v)] = true
	}
	if len(unique) < 4 {
		errors = append(errors, fmt.Errorf("%q must contain at least 4 unique characters", k))
	}
	if len(value) < min || len(value) > max {
		errors = append(errors, fmt.Errorf(
			"%q must be %d to %d characters long. provided string length: %d", k, min, max, len(value)))
	}
	return
}

func expandUsers(cfg []interface{}) []*mq.User {
	users := make([]*mq.User, len(cfg))
	for i, m := range cfg {
		u := m.(map[string]interface{})
		user := mq.User{
			Username: aws.String(u["username"].(string)),
			Password: aws.String(u["password"].(string)),
		}
		if v, ok := u["console_access"]; ok {
			user.ConsoleAccess = aws.Bool(v.(bool))
		}
		if v, ok := u["groups"]; ok {
			user.Groups = flex.ExpandStringSet(v.(*schema.Set))
		}
		users[i] = &user
	}
	return users
}

func expandUsersForBroker(conn *mq.MQ, brokerId string, input []*mq.UserSummary) ([]*mq.User, error) {
	var rawUsers []*mq.User

	for _, u := range input {
		if u == nil {
			continue
		}

		uOut, err := conn.DescribeUser(&mq.DescribeUserInput{
			BrokerId: aws.String(brokerId),
			Username: u.Username,
		})

		if err != nil {
			return nil, err
		}

		user := &mq.User{
			ConsoleAccess: uOut.ConsoleAccess,
			Groups:        uOut.Groups,
			Username:      uOut.Username,
		}

		rawUsers = append(rawUsers, user)
	}

	return rawUsers, nil
}

// We use cfgdUsers to get & set the password
func flattenUsers(users []*mq.User, cfgUsers []interface{}) *schema.Set {
	existingPairs := make(map[string]string)
	for _, u := range cfgUsers {
		user := u.(map[string]interface{})
		username := user["username"].(string)
		existingPairs[username] = user["password"].(string)
	}

	out := make([]interface{}, 0)
	for _, u := range users {
		m := map[string]interface{}{
			"username": aws.StringValue(u.Username),
		}
		password := ""
		if p, ok := existingPairs[aws.StringValue(u.Username)]; ok {
			password = p
		}
		if password != "" {
			m["password"] = password
		}
		if u.ConsoleAccess != nil {
			m["console_access"] = aws.BoolValue(u.ConsoleAccess)
		}
		if len(u.Groups) > 0 {
			m["groups"] = flex.FlattenStringSet(u.Groups)
		}
		out = append(out, m)
	}
	return schema.NewSet(resourceUserHash, out)
}

func expandWeeklyStartTime(cfg []interface{}) *mq.WeeklyStartTime {
	if len(cfg) < 1 {
		return nil
	}

	m := cfg[0].(map[string]interface{})
	return &mq.WeeklyStartTime{
		DayOfWeek: aws.String(m["day_of_week"].(string)),
		TimeOfDay: aws.String(m["time_of_day"].(string)),
		TimeZone:  aws.String(m["time_zone"].(string)),
	}
}

func flattenWeeklyStartTime(wst *mq.WeeklyStartTime) []interface{} {
	if wst == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})
	if wst.DayOfWeek != nil {
		m["day_of_week"] = aws.StringValue(wst.DayOfWeek)
	}
	if wst.TimeOfDay != nil {
		m["time_of_day"] = aws.StringValue(wst.TimeOfDay)
	}
	if wst.TimeZone != nil {
		m["time_zone"] = aws.StringValue(wst.TimeZone)
	}
	return []interface{}{m}
}

func expandConfigurationId(cfg []interface{}) *mq.ConfigurationId {
	if len(cfg) < 1 {
		return nil
	}

	m := cfg[0].(map[string]interface{})
	out := mq.ConfigurationId{
		Id: aws.String(m["id"].(string)),
	}
	if v, ok := m["revision"].(int); ok && v > 0 {
		out.Revision = aws.Int64(int64(v))
	}

	return &out
}

func flattenConfiguration(config *mq.Configurations) []interface{} {
	if config == nil || config.Current == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"id":       aws.StringValue(config.Current.Id),
		"revision": aws.Int64Value(config.Current.Revision),
	}

	return []interface{}{m}
}

func flattenBrokerInstances(instances []*mq.BrokerInstance) []interface{} {
	if len(instances) == 0 {
		return []interface{}{}
	}
	l := make([]interface{}, len(instances))
	for i, instance := range instances {
		m := make(map[string]interface{})
		if instance.ConsoleURL != nil {
			m["console_url"] = aws.StringValue(instance.ConsoleURL)
		}
		if len(instance.Endpoints) > 0 {
			m["endpoints"] = aws.StringValueSlice(instance.Endpoints)
		}
		if instance.IpAddress != nil {
			m["ip_address"] = aws.StringValue(instance.IpAddress)
		}
		l[i] = m
	}

	return l
}

func flattenLogs(logs *mq.LogsSummary) []interface{} {
	if logs == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if logs.General != nil {
		m["general"] = aws.BoolValue(logs.General)
	}

	if logs.Audit != nil {
		m["audit"] = strconv.FormatBool(aws.BoolValue(logs.Audit))
	}

	return []interface{}{m}
}

func expandLogs(engineType string, l []interface{}) *mq.Logs {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	logs := &mq.Logs{}

	if v, ok := m["general"]; ok {
		logs.General = aws.Bool(v.(bool))
	}

	// When the engine type is "RabbitMQ", the parameter audit cannot be set at all.
	if v, ok := m["audit"]; ok {
		if v, null, _ := nullable.Bool(v.(string)).Value(); !null {
			if !strings.EqualFold(engineType, mq.EngineTypeRabbitmq) {
				logs.Audit = aws.Bool(v)
			}
		}
	}

	return logs
}

func flattenLDAPServerMetadata(apiObject *mq.LdapServerMetadataOutput, password string) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Hosts; v != nil {
		tfMap["hosts"] = aws.StringValueSlice(v)
	}
	if v := apiObject.RoleBase; v != nil {
		tfMap["role_base"] = aws.StringValue(v)
	}
	if v := apiObject.RoleName; v != nil {
		tfMap["role_name"] = aws.StringValue(v)
	}
	if v := apiObject.RoleSearchMatching; v != nil {
		tfMap["role_search_matching"] = aws.StringValue(v)
	}
	if v := apiObject.RoleSearchSubtree; v != nil {
		tfMap["role_search_subtree"] = aws.BoolValue(v)
	}
	if password != "" {
		tfMap["service_account_password"] = password
	}
	if v := apiObject.ServiceAccountUsername; v != nil {
		tfMap["service_account_username"] = aws.StringValue(v)
	}
	if v := apiObject.UserBase; v != nil {
		tfMap["user_base"] = aws.StringValue(v)
	}
	if v := apiObject.UserRoleName; v != nil {
		tfMap["user_role_name"] = aws.StringValue(v)
	}
	if v := apiObject.UserSearchMatching; v != nil {
		tfMap["user_search_matching"] = aws.StringValue(v)
	}
	if v := apiObject.UserSearchSubtree; v != nil {
		tfMap["user_search_subtree"] = aws.BoolValue(v)
	}

	return []interface{}{tfMap}
}

func expandLDAPServerMetadata(tfList []interface{}) *mq.LdapServerMetadataInput {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	apiObject := &mq.LdapServerMetadataInput{}

	tfMap := tfList[0].(map[string]interface{})

	if v, ok := tfMap["hosts"]; ok && len(v.([]interface{})) > 0 {
		apiObject.Hosts = flex.ExpandStringList(v.([]interface{}))
	}
	if v, ok := tfMap["role_base"].(string); ok && v != "" {
		apiObject.RoleBase = aws.String(v)
	}
	if v, ok := tfMap["role_name"].(string); ok && v != "" {
		apiObject.RoleName = aws.String(v)
	}
	if v, ok := tfMap["role_search_matching"].(string); ok && v != "" {
		apiObject.RoleSearchMatching = aws.String(v)
	}
	if v, ok := tfMap["role_search_subtree"].(bool); ok {
		apiObject.RoleSearchSubtree = aws.Bool(v)
	}
	if v, ok := tfMap["service_account_password"].(string); ok && v != "" {
		apiObject.ServiceAccountPassword = aws.String(v)
	}
	if v, ok := tfMap["service_account_username"].(string); ok && v != "" {
		apiObject.ServiceAccountUsername = aws.String(v)
	}
	if v, ok := tfMap["user_base"].(string); ok && v != "" {
		apiObject.UserBase = aws.String(v)
	}
	if v, ok := tfMap["user_role_name"].(string); ok && v != "" {
		apiObject.UserRoleName = aws.String(v)
	}
	if v, ok := tfMap["user_search_matching"].(string); ok && v != "" {
		apiObject.UserSearchMatching = aws.String(v)
	}
	if v, ok := tfMap["user_search_subtree"].(bool); ok {
		apiObject.UserSearchSubtree = aws.Bool(v)
	}

	return apiObject
}

var ValidateBrokerName = validation.All(
	validation.StringLenBetween(1, 50),
	validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z_-]+$`), ""),
)
