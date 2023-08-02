package wafregional

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfwaf "github.com/hashicorp/terraform-provider-aws/internal/service/waf"
)

// @SDKResource("aws_wafregional_sql_injection_match_set")
func ResourceSQLInjectionMatchSet() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSQLInjectionMatchSetCreate,
		ReadWithoutTimeout:   resourceSQLInjectionMatchSetRead,
		UpdateWithoutTimeout: resourceSQLInjectionMatchSetUpdate,
		DeleteWithoutTimeout: resourceSQLInjectionMatchSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sql_injection_match_tuple": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      resourceSQLInjectionMatchSetTupleHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_to_match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data": {
										Type:     schema.TypeString,
										Optional: true,
										StateFunc: func(v interface{}) string {
											value := v.(string)
											return strings.ToLower(value)
										},
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"text_transformation": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSQLInjectionMatchSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	log.Printf("[INFO] Creating Regional WAF SQL Injection Match Set: %s", d.Get("name").(string))

	wr := NewRetryer(conn, region)
	out, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
		params := &waf.CreateSqlInjectionMatchSetInput{
			ChangeToken: token,
			Name:        aws.String(d.Get("name").(string)),
		}

		return conn.CreateSqlInjectionMatchSetWithContext(ctx, params)
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Regional WAF SQL Injection Match Set: %s", err)
	}
	resp := out.(*waf.CreateSqlInjectionMatchSetOutput)
	d.SetId(aws.StringValue(resp.SqlInjectionMatchSet.SqlInjectionMatchSetId))

	return append(diags, resourceSQLInjectionMatchSetUpdate(ctx, d, meta)...)
}

func resourceSQLInjectionMatchSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	log.Printf("[INFO] Reading Regional WAF SQL Injection Match Set: %s", d.Get("name").(string))
	params := &waf.GetSqlInjectionMatchSetInput{
		SqlInjectionMatchSetId: aws.String(d.Id()),
	}

	resp, err := conn.GetSqlInjectionMatchSetWithContext(ctx, params)
	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
		log.Printf("[WARN] Regional WAF SQL Injection Match Set (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting Regional WAF SQL Injection Match Set (%s): %s", d.Id(), err)
	}

	d.Set("name", resp.SqlInjectionMatchSet.Name)
	d.Set("sql_injection_match_tuple", flattenSQLInjectionMatchTuples(resp.SqlInjectionMatchSet.SqlInjectionMatchTuples))

	return diags
}

func resourceSQLInjectionMatchSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	if d.HasChange("sql_injection_match_tuple") {
		o, n := d.GetChange("sql_injection_match_tuple")
		oldT, newT := o.(*schema.Set).List(), n.(*schema.Set).List()

		err := updateSQLInjectionMatchSetResourceWR(ctx, d.Id(), oldT, newT, conn, region)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Regional WAF SQL Injection Match Set (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceSQLInjectionMatchSetRead(ctx, d, meta)...)
}

func resourceSQLInjectionMatchSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	oldTuples := d.Get("sql_injection_match_tuple").(*schema.Set).List()

	if len(oldTuples) > 0 {
		noTuples := []interface{}{}
		err := updateSQLInjectionMatchSetResourceWR(ctx, d.Id(), oldTuples, noTuples, conn, region)

		if tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
			return diags
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Regional WAF SQL Injection Match Set (%s): %s", d.Id(), err)
		}
	}

	wr := NewRetryer(conn, region)
	_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
		req := &waf.DeleteSqlInjectionMatchSetInput{
			ChangeToken:            token,
			SqlInjectionMatchSetId: aws.String(d.Id()),
		}

		return conn.DeleteSqlInjectionMatchSetWithContext(ctx, req)
	})
	if tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Regional WAF SQL Injection Match Set (%s): %s", d.Id(), err)
	}

	return diags
}

func updateSQLInjectionMatchSetResourceWR(ctx context.Context, id string, oldT, newT []interface{}, conn *wafregional.WAFRegional, region string) error {
	wr := NewRetryer(conn, region)
	_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
		req := &waf.UpdateSqlInjectionMatchSetInput{
			ChangeToken:            token,
			SqlInjectionMatchSetId: aws.String(id),
			Updates:                diffSQLInjectionMatchTuplesWR(oldT, newT),
		}

		log.Printf("[INFO] Updating Regional WAF SQL Injection Match Set: %s", req)
		return conn.UpdateSqlInjectionMatchSetWithContext(ctx, req)
	})

	return err
}

func diffSQLInjectionMatchTuplesWR(oldT, newT []interface{}) []*waf.SqlInjectionMatchSetUpdate {
	updates := make([]*waf.SqlInjectionMatchSetUpdate, 0)

	for _, od := range oldT {
		tuple := od.(map[string]interface{})

		if idx, contains := sliceContainsMap(newT, tuple); contains {
			newT = append(newT[:idx], newT[idx+1:]...)
			continue
		}

		ftm := tuple["field_to_match"].([]interface{})

		updates = append(updates, &waf.SqlInjectionMatchSetUpdate{
			Action: aws.String(waf.ChangeActionDelete),
			SqlInjectionMatchTuple: &waf.SqlInjectionMatchTuple{
				FieldToMatch:       tfwaf.ExpandFieldToMatch(ftm[0].(map[string]interface{})),
				TextTransformation: aws.String(tuple["text_transformation"].(string)),
			},
		})
	}

	for _, nd := range newT {
		tuple := nd.(map[string]interface{})
		ftm := tuple["field_to_match"].([]interface{})

		updates = append(updates, &waf.SqlInjectionMatchSetUpdate{
			Action: aws.String(waf.ChangeActionInsert),
			SqlInjectionMatchTuple: &waf.SqlInjectionMatchTuple{
				FieldToMatch:       tfwaf.ExpandFieldToMatch(ftm[0].(map[string]interface{})),
				TextTransformation: aws.String(tuple["text_transformation"].(string)),
			},
		})
	}
	return updates
}

func resourceSQLInjectionMatchSetTupleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["field_to_match"]; ok {
		ftms := v.([]interface{})
		ftm := ftms[0].(map[string]interface{})

		if v, ok := ftm["data"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
		}
		buf.WriteString(fmt.Sprintf("%s-", ftm["type"].(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["text_transformation"].(string)))

	return create.StringHashcode(buf.String())
}

func flattenSQLInjectionMatchTuples(ts []*waf.SqlInjectionMatchTuple) []interface{} {
	out := make([]interface{}, len(ts))
	for i, t := range ts {
		m := make(map[string]interface{})
		m["text_transformation"] = aws.StringValue(t.TextTransformation)
		m["field_to_match"] = tfwaf.FlattenFieldToMatch(t.FieldToMatch)
		out[i] = m
	}

	return out
}
