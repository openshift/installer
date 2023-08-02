package sns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/attrmap"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

var (
	subscriptionSchema = map[string]*schema.Schema{
		"arn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"confirmation_timeout_in_minutes": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"confirmation_was_authenticated": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"delivery_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: SuppressEquivalentTopicSubscriptionDeliveryPolicy,
		},
		"endpoint": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"endpoint_auto_confirms": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"filter_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
			StateFunc: func(v interface{}) string {
				json, _ := structure.NormalizeJsonString(v)
				return json
			},
		},
		"filter_policy_scope": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true, // When filter_policy is set, this defaults to MessageAttributes.
			ValidateFunc: validation.StringInSlice(SubscriptionFilterPolicyScope_Values(), false),
		},
		"owner_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pending_confirmation": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(SubscriptionProtocol_Values(), false),
		},
		"raw_message_delivery": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"redrive_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: verify.SuppressEquivalentJSONDiffs,
		},
		"subscription_role_arn": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: verify.ValidARN,
		},
		"topic_arn": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: verify.ValidARN,
		},
	}

	subscriptionAttributeMap = attrmap.New(map[string]string{
		"arn":                            SubscriptionAttributeNameSubscriptionARN,
		"confirmation_was_authenticated": SubscriptionAttributeNameConfirmationWasAuthenticated,
		"delivery_policy":                SubscriptionAttributeNameDeliveryPolicy,
		"endpoint":                       SubscriptionAttributeNameEndpoint,
		"filter_policy":                  SubscriptionAttributeNameFilterPolicy,
		"filter_policy_scope":            SubscriptionAttributeNameFilterPolicyScope,
		"owner_id":                       SubscriptionAttributeNameOwner,
		"pending_confirmation":           SubscriptionAttributeNamePendingConfirmation,
		"protocol":                       SubscriptionAttributeNameProtocol,
		"raw_message_delivery":           SubscriptionAttributeNameRawMessageDelivery,
		"redrive_policy":                 SubscriptionAttributeNameRedrivePolicy,
		"subscription_role_arn":          SubscriptionAttributeNameSubscriptionRoleARN,
		"topic_arn":                      SubscriptionAttributeNameTopicARN,
	}, subscriptionSchema).WithMissingSetToNil("*")
)

// @SDKResource("aws_sns_topic_subscription")
func ResourceTopicSubscription() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTopicSubscriptionCreate,
		ReadWithoutTimeout:   resourceTopicSubscriptionRead,
		UpdateWithoutTimeout: resourceTopicSubscriptionUpdate,
		DeleteWithoutTimeout: resourceTopicSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: resourceTopicSubscriptionCustomizeDiff,

		Schema: subscriptionSchema,
	}
}

func resourceTopicSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SNSConn(ctx)

	attributes, err := subscriptionAttributeMap.ResourceDataToAPIAttributesCreate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Endpoint, Protocol and TopicArn are not passed in Attributes.
	delete(attributes, SubscriptionAttributeNameEndpoint)
	delete(attributes, SubscriptionAttributeNameProtocol)
	delete(attributes, SubscriptionAttributeNameTopicARN)

	protocol := d.Get("protocol").(string)
	input := &sns.SubscribeInput{
		Attributes:            aws.StringMap(attributes),
		Endpoint:              aws.String(d.Get("endpoint").(string)),
		Protocol:              aws.String(protocol),
		ReturnSubscriptionArn: aws.Bool(true), // even if not confirmed, will get ARN
		TopicArn:              aws.String(d.Get("topic_arn").(string)),
	}

	output, err := conn.SubscribeWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating SNS Topic Subscription: %s", err)
	}

	d.SetId(aws.StringValue(output.SubscriptionArn))

	waitForConfirmation := true

	if !d.Get("endpoint_auto_confirms").(bool) && strings.Contains(protocol, "http") {
		waitForConfirmation = false
	}

	if strings.Contains(protocol, "email") {
		waitForConfirmation = false
	}

	if waitForConfirmation {
		timeout := subscriptionPendingConfirmationTimeout
		if strings.Contains(protocol, "http") {
			timeout = time.Duration(int64(d.Get("confirmation_timeout_in_minutes").(int)) * int64(time.Minute))
		}

		if _, err := waitSubscriptionConfirmed(ctx, conn, d.Id(), timeout); err != nil {
			return diag.Errorf("waiting for SNS Topic Subscription (%s) confirmation: %s", d.Id(), err)
		}
	}

	return resourceTopicSubscriptionRead(ctx, d, meta)
}

func resourceTopicSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SNSConn(ctx)

	outputRaw, err := tfresource.RetryWhenNewResourceNotFound(ctx, subscriptionCreateTimeout, func() (interface{}, error) {
		return FindSubscriptionAttributesByARN(ctx, conn, d.Id())
	}, d.IsNewResource())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] SNS Topic Subscription %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading SNS Topic Subscription (%s): %s", d.Id(), err)
	}

	attributes := outputRaw.(map[string]string)

	return diag.FromErr(subscriptionAttributeMap.APIAttributesToResourceData(attributes, d))
}

func resourceTopicSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SNSConn(ctx)

	attributes, err := subscriptionAttributeMap.ResourceDataToAPIAttributesUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	err = putSubscriptionAttributes(ctx, conn, d.Id(), attributes)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTopicSubscriptionRead(ctx, d, meta)
}

func resourceTopicSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SNSConn(ctx)

	log.Printf("[DEBUG] Deleting SNS Topic Subscription: %s", d.Id())
	_, err := conn.UnsubscribeWithContext(ctx, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(d.Id()),
	})

	if tfawserr.ErrMessageContains(err, sns.ErrCodeInvalidParameterException, "Cannot unsubscribe a subscription that is pending confirmation") {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting SNS Topic Subscription (%s): %s", d.Id(), err)
	}

	if _, err := waitSubscriptionDeleted(ctx, conn, d.Id(), subscriptionDeleteTimeout); err != nil {
		return diag.Errorf("waiting for SNS Topic Subscription (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func putSubscriptionAttributes(ctx context.Context, conn *sns.SNS, arn string, attributes map[string]string) error {
	// Filter policy order matters
	filterPolicyScope, ok := attributes[SubscriptionAttributeNameFilterPolicyScope]

	if ok {
		delete(attributes, SubscriptionAttributeNameFilterPolicyScope)
	}

	// MessageBody is backwards-compatible so it should always be applied first
	if filterPolicyScope == SubscriptionFilterPolicyScopeMessageBody {
		err := putSubscriptionAttribute(ctx, conn, arn, SubscriptionAttributeNameFilterPolicyScope, filterPolicyScope)
		if err != nil {
			return err
		}
	}

	for name, value := range attributes {
		err := putSubscriptionAttribute(ctx, conn, arn, name, value)

		if err != nil {
			return err
		}
	}

	// MessageAttributes isn't compatible with nested policies, so it should always be last
	// in case the update also includes a change from a nested policy to a flat policy
	if filterPolicyScope == SubscriptionFilterPolicyScopeMessageAttributes {
		err := putSubscriptionAttribute(ctx, conn, arn, SubscriptionAttributeNameFilterPolicyScope, filterPolicyScope)

		if err != nil {
			return err
		}
	}

	return nil
}

func putSubscriptionAttribute(ctx context.Context, conn *sns.SNS, arn string, name, value string) error {
	// https://docs.aws.amazon.com/sns/latest/dg/message-filtering.html#message-filtering-policy-remove
	if name == SubscriptionAttributeNameFilterPolicy && value == "" {
		value = "{}"
	}

	input := &sns.SetSubscriptionAttributesInput{
		AttributeName:   aws.String(name),
		AttributeValue:  aws.String(value),
		SubscriptionArn: aws.String(arn),
	}

	// The AWS API requires a non-empty string value or nil for the RedrivePolicy attribute,
	// else throws an InvalidParameter error.
	if name == SubscriptionAttributeNameRedrivePolicy && value == "" {
		input.AttributeValue = nil
	}

	_, err := conn.SetSubscriptionAttributesWithContext(ctx, input)

	if err != nil {
		return fmt.Errorf("setting SNS Topic Subscription (%s) attribute (%s): %w", arn, name, err)
	}

	return nil
}

func FindSubscriptionAttributesByARN(ctx context.Context, conn *sns.SNS, arn string) (map[string]string, error) {
	input := &sns.GetSubscriptionAttributesInput{
		SubscriptionArn: aws.String(arn),
	}

	output, err := conn.GetSubscriptionAttributesWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, sns.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Attributes) == 0 {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return aws.StringValueMap(output.Attributes), nil
}

func statusSubscriptionPendingConfirmation(ctx context.Context, conn *sns.SNS, arn string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindSubscriptionAttributesByARN(ctx, conn, arn)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, output[SubscriptionAttributeNamePendingConfirmation], nil
	}
}

const (
	subscriptionCreateTimeout              = 2 * time.Minute
	subscriptionPendingConfirmationTimeout = 2 * time.Minute
	subscriptionDeleteTimeout              = 2 * time.Minute
)

func waitSubscriptionConfirmed(ctx context.Context, conn *sns.SNS, arn string, timeout time.Duration) (map[string]string, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"true"},
		Target:  []string{"false"},
		Refresh: statusSubscriptionPendingConfirmation(ctx, conn, arn),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(map[string]string); ok {
		return output, err
	}

	return nil, err
}

func waitSubscriptionDeleted(ctx context.Context, conn *sns.SNS, arn string, timeout time.Duration) (map[string]string, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"false", "true"},
		Target:  []string{},
		Refresh: statusSubscriptionPendingConfirmation(ctx, conn, arn),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(map[string]string); ok {
		return output, err
	}

	return nil, err
}

type TopicSubscriptionDeliveryPolicy struct {
	Guaranteed         bool                                                 `json:"guaranteed,omitempty"`
	HealthyRetryPolicy *TopicSubscriptionDeliveryPolicyHealthyRetryPolicy   `json:"healthyRetryPolicy,omitempty"`
	SicklyRetryPolicy  *snsTopicSubscriptionDeliveryPolicySicklyRetryPolicy `json:"sicklyRetryPolicy,omitempty"`
	ThrottlePolicy     *snsTopicSubscriptionDeliveryPolicyThrottlePolicy    `json:"throttlePolicy,omitempty"`
}

func (s TopicSubscriptionDeliveryPolicy) String() string {
	return awsutil.Prettify(s)
}

func (s TopicSubscriptionDeliveryPolicy) GoString() string {
	return s.String()
}

type TopicSubscriptionDeliveryPolicyHealthyRetryPolicy struct {
	BackoffFunction    string `json:"backoffFunction,omitempty"`
	MaxDelayTarget     int    `json:"maxDelayTarget,omitempty"`
	MinDelayTarget     int    `json:"minDelayTarget,omitempty"`
	NumMaxDelayRetries int    `json:"numMaxDelayRetries,omitempty"`
	NumMinDelayRetries int    `json:"numMinDelayRetries,omitempty"`
	NumNoDelayRetries  int    `json:"numNoDelayRetries,omitempty"`
	NumRetries         int    `json:"numRetries,omitempty"`
}

func (s TopicSubscriptionDeliveryPolicyHealthyRetryPolicy) String() string {
	return awsutil.Prettify(s)
}

func (s TopicSubscriptionDeliveryPolicyHealthyRetryPolicy) GoString() string {
	return s.String()
}

type snsTopicSubscriptionDeliveryPolicySicklyRetryPolicy struct {
	BackoffFunction    string `json:"backoffFunction,omitempty"`
	MaxDelayTarget     int    `json:"maxDelayTarget,omitempty"`
	MinDelayTarget     int    `json:"minDelayTarget,omitempty"`
	NumMaxDelayRetries int    `json:"numMaxDelayRetries,omitempty"`
	NumMinDelayRetries int    `json:"numMinDelayRetries,omitempty"`
	NumNoDelayRetries  int    `json:"numNoDelayRetries,omitempty"`
	NumRetries         int    `json:"numRetries,omitempty"`
}

func (s snsTopicSubscriptionDeliveryPolicySicklyRetryPolicy) String() string {
	return awsutil.Prettify(s)
}

func (s snsTopicSubscriptionDeliveryPolicySicklyRetryPolicy) GoString() string {
	return s.String()
}

type snsTopicSubscriptionDeliveryPolicyThrottlePolicy struct {
	MaxReceivesPerSecond int `json:"maxReceivesPerSecond,omitempty"`
}

func (s snsTopicSubscriptionDeliveryPolicyThrottlePolicy) String() string {
	return awsutil.Prettify(s)
}

func (s snsTopicSubscriptionDeliveryPolicyThrottlePolicy) GoString() string {
	return s.String()
}

type TopicSubscriptionRedrivePolicy struct {
	DeadLetterTargetArn string `json:"deadLetterTargetArn,omitempty"`
}

func SuppressEquivalentTopicSubscriptionDeliveryPolicy(k, old, new string, d *schema.ResourceData) bool {
	ob, err := normalizeTopicSubscriptionDeliveryPolicy(old)
	if err != nil {
		log.Print(err)
		return false
	}

	nb, err := normalizeTopicSubscriptionDeliveryPolicy(new)
	if err != nil {
		log.Print(err)
		return false
	}

	return verify.JSONBytesEqual(ob, nb)
}

func normalizeTopicSubscriptionDeliveryPolicy(policy string) ([]byte, error) {
	var deliveryPolicy TopicSubscriptionDeliveryPolicy

	if err := json.Unmarshal([]byte(policy), &deliveryPolicy); err != nil {
		return nil, fmt.Errorf("[WARN] Unable to unmarshal SNS Topic Subscription delivery policy JSON: %s", err)
	}

	normalizedDeliveryPolicy, err := json.Marshal(deliveryPolicy)

	if err != nil {
		return nil, fmt.Errorf("[WARN] Unable to marshal SNS Topic Subscription delivery policy back to JSON: %s", err)
	}

	b := bytes.NewBufferString("")
	if err := json.Compact(b, normalizedDeliveryPolicy); err != nil {
		return nil, fmt.Errorf("[WARN] Unable to marshal SNS Topic Subscription delivery policy back to JSON: %s", err)
	}

	return b.Bytes(), nil
}

func resourceTopicSubscriptionCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	hasPolicy := diff.Get("filter_policy").(string) != ""
	hasScope := !diff.GetRawConfig().GetAttr("filter_policy_scope").IsNull()
	hadScope := diff.Get("filter_policy_scope").(string) != ""

	if hasPolicy && !hasScope {
		if !hadScope {
			// When the filter_policy_scope hasn't been read back from the API,
			// don't attempt to set a value. Either the default will be computed
			// on the next read, or this is a partition that doesn't support it.
			return nil
		}

		// When the scope is removed from configuration, the API will
		// continue reading back the last value so long as the policy
		// itself still exists. The expected result would be to revert
		// to the default value of the attribute (MessageAttributes).
		return diff.SetNew("filter_policy_scope", SubscriptionFilterPolicyScopeMessageAttributes)
	}

	if !hasPolicy && !hasScope {
		// When the policy is not set, the API silently drops the scope.
		return diff.Clear("filter_policy_scope")
	}

	if !hasPolicy && hasScope {
		// Make it explicit that the scope doesn't exist without a policy.
		return errors.New("filter_policy is required when filter_policy_scope is set")
	}

	return nil
}
