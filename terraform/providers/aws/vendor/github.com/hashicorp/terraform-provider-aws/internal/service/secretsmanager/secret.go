package secretsmanager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfiam "github.com/hashicorp/terraform-provider-aws/internal/service/iam"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_secretsmanager_secret", name="Secret")
// @Tags(identifierAttribute="id")
func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecretCreate,
		ReadWithoutTimeout:   resourceSecretRead,
		UpdateWithoutTimeout: resourceSecretUpdate,
		DeleteWithoutTimeout: resourceSecretDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force_overwrite_replica_secret": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validSecretName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validSecretNamePrefix,
			},
			"policy": {
				Type:                  schema.TypeString,
				Optional:              true,
				Computed:              true,
				ValidateFunc:          validation.StringIsJSON,
				DiffSuppressFunc:      verify.SuppressEquivalentPolicyDiffs,
				DiffSuppressOnRefresh: true,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			"recovery_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				ValidateFunc: validation.Any(
					validation.IntBetween(7, 30),
					validation.IntInSlice([]int{0}),
				),
			},
			"replica": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      secretReplicaHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"last_accessed_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn(ctx)

	secretName := create.Name(d.Get("name").(string), d.Get("name_prefix").(string))
	input := &secretsmanager.CreateSecretInput{
		Description:                 aws.String(d.Get("description").(string)),
		ForceOverwriteReplicaSecret: aws.Bool(d.Get("force_overwrite_replica_secret").(bool)),
		Name:                        aws.String(secretName),
		Tags:                        GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("replica"); ok && v.(*schema.Set).Len() > 0 {
		input.AddReplicaRegions = expandSecretReplicas(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] Creating Secrets Manager Secret: %s", input)

	// Retry for secret recreation after deletion
	var output *secretsmanager.CreateSecretOutput
	err := retry.RetryContext(ctx, PropagationTimeout, func() *retry.RetryError {
		var err error
		output, err = conn.CreateSecretWithContext(ctx, input)
		// Temporarily retry on these errors to support immediate secret recreation:
		// InvalidRequestException: You can’t perform this operation on the secret because it was deleted.
		// InvalidRequestException: You can't create this secret because a secret with this name is already scheduled for deletion.
		if tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "scheduled for deletion") || tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "was deleted") {
			return retry.RetryableError(err)
		}
		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		output, err = conn.CreateSecretWithContext(ctx, input)
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Secrets Manager Secret: %s", err)
	}

	d.SetId(aws.StringValue(output.ARN))

	if v, ok := d.GetOk("policy"); ok && v.(string) != "" && v.(string) != "{}" {
		policy, err := structure.NormalizeJsonString(v.(string))
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "policy (%s) is invalid JSON: %s", v.(string), err)
		}

		input := &secretsmanager.PutResourcePolicyInput{
			ResourcePolicy: aws.String(policy),
			SecretId:       aws.String(d.Id()),
		}

		err = retry.RetryContext(ctx, PropagationTimeout, func() *retry.RetryError {
			_, err := conn.PutResourcePolicyWithContext(ctx, input)
			if tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeMalformedPolicyDocumentException,
				"This resource policy contains an unsupported principal") {
				return retry.RetryableError(err)
			}
			if err != nil {
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			_, err = conn.PutResourcePolicyWithContext(ctx, input)
		}
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting Secrets Manager Secret %q policy: %s", d.Id(), err)
		}
	}

	return append(diags, resourceSecretRead(ctx, d, meta)...)
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn(ctx)

	outputRaw, err := tfresource.RetryWhenNewResourceNotFound(ctx, PropagationTimeout, func() (interface{}, error) {
		return FindSecretByID(ctx, conn, d.Id())
	}, d.IsNewResource())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Secrets Manager Secret (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret (%s): %s", d.Id(), err)
	}

	output := outputRaw.(*secretsmanager.DescribeSecretOutput)

	d.Set("arn", output.ARN)
	d.Set("description", output.Description)
	d.Set("kms_key_id", output.KmsKeyId)
	d.Set("name", output.Name)
	d.Set("name_prefix", create.NamePrefixFromName(aws.StringValue(output.Name)))

	if err := d.Set("replica", flattenSecretReplicas(output.ReplicationStatus)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting replica: %s", err)
	}

	var policy *secretsmanager.GetResourcePolicyOutput
	err = tfresource.Retry(ctx, PropagationTimeout, func() *retry.RetryError {
		var err error
		policy, err = conn.GetResourcePolicyWithContext(ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String(d.Id()),
		})
		if err != nil {
			return retry.NonRetryableError(err)
		}

		if policy.ResourcePolicy != nil {
			valid, err := tfiam.PolicyHasValidAWSPrincipals(aws.StringValue(policy.ResourcePolicy))
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if !valid {
				log.Printf("[DEBUG] Retrying because of invalid principals")
				return retry.RetryableError(errors.New("contains invalid principals"))
			}
		}

		return nil
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret (%s) policy: %s", d.Id(), err)
	} else if v := policy.ResourcePolicy; v != nil {
		policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), aws.StringValue(v))
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading Secrets Manager Secret (%s): %s", d.Id(), err)
		}

		d.Set("policy", policyToSet)
	} else {
		d.Set("policy", "")
	}

	SetTagsOut(ctx, output.Tags)

	return diags
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn(ctx)

	if d.HasChange("replica") {
		o, n := d.GetChange("replica")

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		err := removeSecretReplicas(ctx, conn, d.Id(), os.Difference(ns).List())

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deleting Secrets Manager Secret (%s) replica: %s", d.Id(), err)
		}

		err = addSecretReplicas(ctx, conn, d.Id(), d.Get("force_overwrite_replica_secret").(bool), ns.Difference(os).List())

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "adding Secrets Manager Secret (%s) replica: %s", d.Id(), err)
		}
	}

	if d.HasChanges("description", "kms_key_id") {
		input := &secretsmanager.UpdateSecretInput{
			Description: aws.String(d.Get("description").(string)),
			SecretId:    aws.String(d.Id()),
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		log.Printf("[DEBUG] Updating Secrets Manager Secret: %s", input)
		_, err := conn.UpdateSecretWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Secrets Manager Secret (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("policy") {
		if v, ok := d.GetOk("policy"); ok && v.(string) != "" && v.(string) != "{}" {
			policy, err := structure.NormalizeJsonString(v.(string))
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "policy contains an invalid JSON: %s", err)
			}

			input := &secretsmanager.PutResourcePolicyInput{
				ResourcePolicy: aws.String(policy),
				SecretId:       aws.String(d.Id()),
			}

			log.Printf("[DEBUG] Setting Secrets Manager Secret resource policy: %s", input)
			_, err = tfresource.RetryWhenAWSErrMessageContains(ctx, PropagationTimeout,
				func() (interface{}, error) {
					return conn.PutResourcePolicyWithContext(ctx, input)
				},
				secretsmanager.ErrCodeMalformedPolicyDocumentException, "This resource policy contains an unsupported principal")

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "setting Secrets Manager Secret (%s) policy: %s", d.Id(), err)
			}
		} else {
			log.Printf("[DEBUG] Removing Secrets Manager Secret policy: %s", d.Id())
			_, err := conn.DeleteResourcePolicyWithContext(ctx, &secretsmanager.DeleteResourcePolicyInput{
				SecretId: aws.String(d.Id()),
			})

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "removing Secrets Manager Secret (%s) policy: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceSecretRead(ctx, d, meta)...)
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SecretsManagerConn(ctx)

	if v, ok := d.GetOk("replica"); ok && v.(*schema.Set).Len() > 0 {
		err := removeSecretReplicas(ctx, conn, d.Id(), v.(*schema.Set).List())

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "deleting Secrets Manager Secret (%s) replica: %s", d.Id(), err)
		}
	}

	input := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(d.Id()),
	}

	recoveryWindowInDays := d.Get("recovery_window_in_days").(int)
	if recoveryWindowInDays == 0 {
		input.ForceDeleteWithoutRecovery = aws.Bool(true)
	} else {
		input.RecoveryWindowInDays = aws.Int64(int64(recoveryWindowInDays))
	}

	log.Printf("[DEBUG] Deleting Secrets Manager Secret: %s", d.Id())
	_, err := conn.DeleteSecretWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Secrets Manager Secret (%s): %s", d.Id(), err)
	}

	_, err = tfresource.RetryUntilNotFound(ctx, PropagationTimeout, func() (interface{}, error) {
		return FindSecretByID(ctx, conn, d.Id())
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Secrets Manager Secret (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func removeSecretReplicas(ctx context.Context, conn *secretsmanager.SecretsManager, id string, tfList []interface{}) error {
	if len(tfList) == 0 {
		return nil
	}

	input := &secretsmanager.RemoveRegionsFromReplicationInput{
		SecretId: aws.String(id),
	}

	var regions []string

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		regions = append(regions, tfMap["region"].(string))
	}

	input.RemoveReplicaRegions = aws.StringSlice(regions)

	log.Printf("[DEBUG] Removing Secrets Manager Secret Replicas: %s", input)

	_, err := conn.RemoveRegionsFromReplicationWithContext(ctx, input)

	if err != nil {
		if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
			return nil
		}

		return err
	}

	return nil
}

func addSecretReplicas(ctx context.Context, conn *secretsmanager.SecretsManager, id string, forceOverwrite bool, tfList []interface{}) error {
	if len(tfList) == 0 {
		return nil
	}

	input := &secretsmanager.ReplicateSecretToRegionsInput{
		SecretId:                    aws.String(id),
		ForceOverwriteReplicaSecret: aws.Bool(forceOverwrite),
		AddReplicaRegions:           expandSecretReplicas(tfList),
	}

	log.Printf("[DEBUG] Removing Secrets Manager Secret Replica: %s", input)

	_, err := conn.ReplicateSecretToRegionsWithContext(ctx, input)

	return err
}

func expandSecretReplica(tfMap map[string]interface{}) *secretsmanager.ReplicaRegionType {
	if tfMap == nil {
		return nil
	}

	apiObject := &secretsmanager.ReplicaRegionType{}

	if v, ok := tfMap["kms_key_id"].(string); ok && v != "" {
		apiObject.KmsKeyId = aws.String(v)
	}

	if v, ok := tfMap["region"].(string); ok && v != "" {
		apiObject.Region = aws.String(v)
	}

	return apiObject
}

func expandSecretReplicas(tfList []interface{}) []*secretsmanager.ReplicaRegionType {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*secretsmanager.ReplicaRegionType

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandSecretReplica(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenSecretReplica(apiObject *secretsmanager.ReplicationStatusType) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.KmsKeyId; v != nil {
		tfMap["kms_key_id"] = aws.StringValue(v)
	}

	if v := apiObject.LastAccessedDate; v != nil {
		tfMap["last_accessed_date"] = aws.TimeValue(v).Format(time.RFC3339)
	}

	if v := apiObject.Region; v != nil {
		tfMap["region"] = aws.StringValue(v)
	}

	if v := apiObject.Status; v != nil {
		tfMap["status"] = aws.StringValue(v)
	}

	if v := apiObject.StatusMessage; v != nil {
		tfMap["status_message"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenSecretReplicas(apiObjects []*secretsmanager.ReplicationStatusType) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenSecretReplica(apiObject))
	}

	return tfList
}

func secretReplicaHash(v interface{}) int {
	var buf bytes.Buffer

	m := v.(map[string]interface{})

	if v, ok := m["kms_key_id"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	if v, ok := m["region"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", v))
	}

	return create.StringHashcode(buf.String())
}

func findSecret(ctx context.Context, conn *secretsmanager.SecretsManager, input *secretsmanager.DescribeSecretInput) (*secretsmanager.DescribeSecretOutput, error) {
	output, err := conn.DescribeSecretWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindSecretByID(ctx context.Context, conn *secretsmanager.SecretsManager, id string) (*secretsmanager.DescribeSecretOutput, error) {
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	}

	output, err := findSecret(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	if output.DeletedDate != nil {
		return nil, &retry.NotFoundError{LastRequest: input}
	}

	return output, nil
}
