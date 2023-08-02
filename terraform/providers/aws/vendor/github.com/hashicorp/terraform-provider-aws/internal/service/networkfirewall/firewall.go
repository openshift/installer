package networkfirewall

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_networkfirewall_firewall", name="Firewall")
// @Tags(identifierAttribute="id")
func ResourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceFirewallCreate,
		ReadWithoutTimeout:   resourceFirewallRead,
		UpdateWithoutTimeout: resourceFirewallUpdate,
		DeleteWithoutTimeout: resourceFirewallDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.Sequence(
			customdiff.ComputedIf("firewall_status", func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("subnet_mapping")
			}),
			verify.SetTagsDiff,
		),

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_configuration": encryptionConfigurationSchema(),
			"firewall_policy_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"firewall_policy_change_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"firewall_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync_states": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attachment": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint_id": {
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
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_change_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"subnet_mapping": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(networkfirewall.IPAddressType_Values(), false),
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"update_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)

	name := d.Get("name").(string)
	input := &networkfirewall.CreateFirewallInput{
		FirewallName:      aws.String(name),
		FirewallPolicyArn: aws.String(d.Get("firewall_policy_arn").(string)),
		SubnetMappings:    expandSubnetMappings(d.Get("subnet_mapping").(*schema.Set).List()),
		Tags:              GetTagsIn(ctx),
		VpcId:             aws.String(d.Get("vpc_id").(string)),
	}

	if v, ok := d.GetOk("delete_protection"); ok {
		input.DeleteProtection = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("encryption_configuration"); ok {
		input.EncryptionConfiguration = expandEncryptionConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("firewall_policy_change_protection"); ok {
		input.FirewallPolicyChangeProtection = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("subnet_change_protection"); ok {
		input.SubnetChangeProtection = aws.Bool(v.(bool))
	}

	output, err := conn.CreateFirewallWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating NetworkFirewall Firewall (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.Firewall.FirewallArn))

	if _, err := waitFirewallCreated(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for NetworkFirewall Firewall (%s) create: %s", d.Id(), err)
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)

	output, err := FindFirewallByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] NetworkFirewall Firewall (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading NetworkFirewall Firewall (%s): %s", d.Id(), err)
	}

	firewall := output.Firewall
	d.Set("arn", firewall.FirewallArn)
	d.Set("delete_protection", firewall.DeleteProtection)
	d.Set("description", firewall.Description)
	if err := d.Set("encryption_configuration", flattenEncryptionConfiguration(firewall.EncryptionConfiguration)); err != nil {
		return diag.Errorf("setting encryption_configuration: %s", err)
	}
	d.Set("firewall_policy_arn", firewall.FirewallPolicyArn)
	d.Set("firewall_policy_change_protection", firewall.FirewallPolicyChangeProtection)
	if err := d.Set("firewall_status", flattenFirewallStatus(output.FirewallStatus)); err != nil {
		return diag.Errorf("setting firewall_status: %s", err)
	}
	d.Set("name", firewall.FirewallName)
	d.Set("subnet_change_protection", firewall.SubnetChangeProtection)
	if err := d.Set("subnet_mapping", flattenSubnetMappings(firewall.SubnetMappings)); err != nil {
		return diag.Errorf("setting subnet_mapping: %s", err)
	}
	d.Set("update_token", output.UpdateToken)
	d.Set("vpc_id", firewall.VpcId)

	SetTagsOut(ctx, firewall.Tags)

	return nil
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)
	updateToken := d.Get("update_token").(string)

	if d.HasChange("delete_protection") {
		input := &networkfirewall.UpdateFirewallDeleteProtectionInput{
			DeleteProtection: aws.Bool(d.Get("delete_protection").(bool)),
			FirewallArn:      aws.String(d.Id()),
			UpdateToken:      aws.String(updateToken),
		}

		output, err := conn.UpdateFirewallDeleteProtectionWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) delete protection: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	if d.HasChange("description") {
		input := &networkfirewall.UpdateFirewallDescriptionInput{
			Description: aws.String(d.Get("description").(string)),
			FirewallArn: aws.String(d.Id()),
			UpdateToken: aws.String(updateToken),
		}

		output, err := conn.UpdateFirewallDescriptionWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) description: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	if d.HasChange("encryption_configuration") {
		input := &networkfirewall.UpdateFirewallEncryptionConfigurationInput{
			EncryptionConfiguration: expandEncryptionConfiguration(d.Get("encryption_configuration").([]interface{})),
			FirewallArn:             aws.String(d.Id()),
			UpdateToken:             aws.String(updateToken),
		}

		output, err := conn.UpdateFirewallEncryptionConfigurationWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) encryption configuration: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	// Note: The *_change_protection fields below are handled before their respective fields
	// to account for disabling and subsequent changes

	if d.HasChange("firewall_policy_change_protection") {
		input := &networkfirewall.UpdateFirewallPolicyChangeProtectionInput{
			FirewallArn:                    aws.String(d.Id()),
			FirewallPolicyChangeProtection: aws.Bool(d.Get("firewall_policy_change_protection").(bool)),
			UpdateToken:                    aws.String(updateToken),
		}

		output, err := conn.UpdateFirewallPolicyChangeProtectionWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) firewall policy change protection: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	if d.HasChange("firewall_policy_arn") {
		input := &networkfirewall.AssociateFirewallPolicyInput{
			FirewallArn:       aws.String(d.Id()),
			FirewallPolicyArn: aws.String(d.Get("firewall_policy_arn").(string)),
			UpdateToken:       aws.String(updateToken),
		}

		output, err := conn.AssociateFirewallPolicyWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) firewall policy ARN: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	if d.HasChange("subnet_change_protection") {
		input := &networkfirewall.UpdateSubnetChangeProtectionInput{
			FirewallArn:            aws.String(d.Id()),
			SubnetChangeProtection: aws.Bool(d.Get("subnet_change_protection").(bool)),
			UpdateToken:            aws.String(updateToken),
		}

		output, err := conn.UpdateSubnetChangeProtectionWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating NetworkFirewall Firewall (%s) subnet change protection: %s", d.Id(), err)
		}

		updateToken = aws.StringValue(output.UpdateToken)
	}

	if d.HasChange("subnet_mapping") {
		o, n := d.GetChange("subnet_mapping")
		subnetsToRemove, subnetsToAdd := subnetMappingsDiff(o.(*schema.Set), n.(*schema.Set))

		if len(subnetsToAdd) > 0 {
			input := &networkfirewall.AssociateSubnetsInput{
				FirewallArn:    aws.String(d.Id()),
				SubnetMappings: subnetsToAdd,
				UpdateToken:    aws.String(updateToken),
			}

			_, err := conn.AssociateSubnetsWithContext(ctx, input)

			if err != nil {
				return diag.Errorf("associating NetworkFirewall Firewall (%s) subnets: %s", d.Id(), err)
			}

			updateToken, err = waitFirewallUpdated(ctx, conn, d.Id())

			if err != nil {
				return diag.Errorf("waiting for NetworkFirewall Firewall (%s) update: %s", d.Id(), err)
			}
		}

		if len(subnetsToRemove) > 0 {
			input := &networkfirewall.DisassociateSubnetsInput{
				FirewallArn: aws.String(d.Id()),
				SubnetIds:   aws.StringSlice(subnetsToRemove),
				UpdateToken: aws.String(updateToken),
			}

			_, err := conn.DisassociateSubnetsWithContext(ctx, input)

			if err == nil {
				/*updateToken*/ _, err = waitFirewallUpdated(ctx, conn, d.Id())

				if err != nil {
					return diag.Errorf("waiting for NetworkFirewall Firewall (%s) update: %s", d.Id(), err)
				}
			} else if !tfawserr.ErrMessageContains(err, networkfirewall.ErrCodeInvalidRequestException, "inaccessible") {
				return diag.Errorf("disassociating NetworkFirewall Firewall (%s) subnets: %s", d.Id(), err)
			}
		}
	}

	return resourceFirewallRead(ctx, d, meta)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)

	log.Printf("[DEBUG] Deleting NetworkFirewall Firewall: %s", d.Id())
	_, err := conn.DeleteFirewallWithContext(ctx, &networkfirewall.DeleteFirewallInput{
		FirewallArn: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, networkfirewall.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting NetworkFirewall Firewall (%s): %s", d.Id(), err)
	}

	if _, err := waitFirewallDeleted(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for NetworkFirewall Firewall (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindFirewallByARN(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.DescribeFirewallOutput, error) {
	input := &networkfirewall.DescribeFirewallInput{
		FirewallArn: aws.String(arn),
	}

	output, err := conn.DescribeFirewallWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, networkfirewall.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Firewall == nil || output.FirewallStatus == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func statusFirewall(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindFirewallByARN(ctx, conn, arn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.FirewallStatus.Status), nil
	}
}

const (
	firewallTimeout = 20 * time.Minute
)

func waitFirewallCreated(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.Firewall, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{networkfirewall.FirewallStatusValueProvisioning},
		Target:  []string{networkfirewall.FirewallStatusValueReady},
		Refresh: statusFirewall(ctx, conn, arn),
		Timeout: firewallTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkfirewall.DescribeFirewallOutput); ok {
		return output.Firewall, err
	}

	return nil, err
}

func waitFirewallUpdated(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (string, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{networkfirewall.FirewallStatusValueProvisioning},
		Target:  []string{networkfirewall.FirewallStatusValueReady},
		Refresh: statusFirewall(ctx, conn, arn),
		Timeout: firewallTimeout,
		// Delay added to account for Associate/DisassociateSubnet calls that return
		// a READY status immediately after the method is called instead of immediately
		// returning PROVISIONING
		Delay: 30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkfirewall.DescribeFirewallOutput); ok {
		return aws.StringValue(output.UpdateToken), err
	}

	return "", err
}

func waitFirewallDeleted(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.Firewall, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{networkfirewall.FirewallStatusValueDeleting},
		Target:  []string{},
		Refresh: statusFirewall(ctx, conn, arn),
		Timeout: firewallTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkfirewall.DescribeFirewallOutput); ok {
		return output.Firewall, err
	}

	return nil, err
}

func expandSubnetMappings(l []interface{}) []*networkfirewall.SubnetMapping {
	mappings := make([]*networkfirewall.SubnetMapping, 0, len(l))
	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}
		mapping := &networkfirewall.SubnetMapping{
			SubnetId: aws.String(tfMap["subnet_id"].(string)),
		}
		if v, ok := tfMap["ip_address_type"].(string); ok && v != "" {
			mapping.IPAddressType = aws.String(v)
		}
		mappings = append(mappings, mapping)
	}

	return mappings
}

func expandSubnetMappingIDs(l []interface{}) []string {
	var ids []string
	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := tfMap["subnet_id"].(string); ok && id != "" {
			ids = append(ids, id)
		}
	}

	return ids
}

func flattenFirewallStatus(status *networkfirewall.FirewallStatus) []interface{} {
	if status == nil {
		return nil
	}

	m := map[string]interface{}{
		"sync_states": flattenSyncStates(status.SyncStates),
	}

	return []interface{}{m}
}

func flattenSyncStates(s map[string]*networkfirewall.SyncState) []interface{} {
	if s == nil {
		return nil
	}

	syncStates := make([]interface{}, 0, len(s))
	for k, v := range s {
		m := map[string]interface{}{
			"availability_zone": k,
			"attachment":        flattenSyncStateAttachment(v.Attachment),
		}
		syncStates = append(syncStates, m)
	}

	return syncStates
}

func flattenSyncStateAttachment(a *networkfirewall.Attachment) []interface{} {
	if a == nil {
		return nil
	}

	m := map[string]interface{}{
		"endpoint_id": aws.StringValue(a.EndpointId),
		"subnet_id":   aws.StringValue(a.SubnetId),
	}

	return []interface{}{m}
}

func flattenSubnetMappings(sm []*networkfirewall.SubnetMapping) []interface{} {
	mappings := make([]interface{}, 0, len(sm))
	for _, s := range sm {
		m := map[string]interface{}{
			"subnet_id":       aws.StringValue(s.SubnetId),
			"ip_address_type": aws.StringValue(s.IPAddressType),
		}
		mappings = append(mappings, m)
	}

	return mappings
}

func subnetMappingsHash(v interface{}) int {
	var buf bytes.Buffer

	tfMap, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	if id, ok := tfMap["subnet_id"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}
	if id, ok := tfMap["ip_address_type"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return create.StringHashcode(buf.String())
}

func subnetMappingsDiff(old, new *schema.Set) ([]string, []*networkfirewall.SubnetMapping) {
	if old.Len() == 0 {
		return nil, expandSubnetMappings(new.List())
	}
	if new.Len() == 0 {
		return expandSubnetMappingIDs(old.List()), nil
	}

	oldHashedSet := schema.NewSet(subnetMappingsHash, old.List())
	newHashedSet := schema.NewSet(subnetMappingsHash, new.List())

	toRemove := oldHashedSet.Difference(newHashedSet)
	toAdd := new.Difference(old)

	subnetsToRemove := expandSubnetMappingIDs(toRemove.List())
	subnetsToAdd := expandSubnetMappings(toAdd.List())

	return subnetsToRemove, subnetsToAdd
}
