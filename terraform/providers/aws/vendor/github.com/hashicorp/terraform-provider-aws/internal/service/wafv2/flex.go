package wafv2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func expandRules(l []interface{}) []*wafv2.Rule {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	rules := make([]*wafv2.Rule, 0)

	for _, rule := range l {
		if rule == nil {
			continue
		}
		rules = append(rules, expandRule(rule.(map[string]interface{})))
	}

	return rules
}

func expandRule(m map[string]interface{}) *wafv2.Rule {
	if m == nil {
		return nil
	}

	rule := &wafv2.Rule{
		Name:             aws.String(m["name"].(string)),
		Priority:         aws.Int64(int64(m["priority"].(int))),
		Action:           expandRuleAction(m["action"].([]interface{})),
		Statement:        expandRuleGroupRootStatement(m["statement"].([]interface{})),
		VisibilityConfig: expandVisibilityConfig(m["visibility_config"].([]interface{})),
	}

	if v, ok := m["rule_label"].(*schema.Set); ok && v.Len() > 0 {
		rule.RuleLabels = expandRuleLabels(v.List())
	}

	return rule
}

func expandRuleLabels(l []interface{}) []*wafv2.Label {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	labels := make([]*wafv2.Label, 0)

	for _, label := range l {
		if label == nil {
			continue
		}
		m := label.(map[string]interface{})
		labels = append(labels, &wafv2.Label{
			Name: aws.String(m["name"].(string)),
		})
	}

	return labels
}

func expandRuleAction(l []interface{}) *wafv2.RuleAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	action := &wafv2.RuleAction{}

	if v, ok := m["allow"]; ok && len(v.([]interface{})) > 0 {
		action.Allow = expandAllowAction(v.([]interface{}))
	}

	if v, ok := m["block"]; ok && len(v.([]interface{})) > 0 {
		action.Block = expandBlockAction(v.([]interface{}))
	}

	if v, ok := m["captcha"]; ok && len(v.([]interface{})) > 0 {
		action.Captcha = expandCaptchaAction(v.([]interface{}))
	}

	if v, ok := m["count"]; ok && len(v.([]interface{})) > 0 {
		action.Count = expandCountAction(v.([]interface{}))
	}

	return action
}

func expandAllowAction(l []interface{}) *wafv2.AllowAction {
	action := &wafv2.AllowAction{}

	if len(l) == 0 || l[0] == nil {
		return action
	}

	m, ok := l[0].(map[string]interface{})
	if !ok {
		return action
	}

	if v, ok := m["custom_request_handling"].([]interface{}); ok && len(v) > 0 {
		action.CustomRequestHandling = expandCustomRequestHandling(v)
	}

	return action
}

func expandBlockAction(l []interface{}) *wafv2.BlockAction {
	action := &wafv2.BlockAction{}

	if len(l) == 0 || l[0] == nil {
		return action
	}

	m, ok := l[0].(map[string]interface{})
	if !ok {
		return action
	}

	if v, ok := m["custom_response"].([]interface{}); ok && len(v) > 0 {
		action.CustomResponse = expandCustomResponse(v)
	}

	return action
}

func expandCaptchaAction(l []interface{}) *wafv2.CaptchaAction {
	action := &wafv2.CaptchaAction{}

	if len(l) == 0 || l[0] == nil {
		return action
	}

	m, ok := l[0].(map[string]interface{})
	if !ok {
		return action
	}

	if v, ok := m["custom_request_handling"].([]interface{}); ok && len(v) > 0 {
		action.CustomRequestHandling = expandCustomRequestHandling(v)
	}

	return action
}

func expandCountAction(l []interface{}) *wafv2.CountAction {
	action := &wafv2.CountAction{}

	if len(l) == 0 || l[0] == nil {
		return action
	}

	m, ok := l[0].(map[string]interface{})
	if !ok {
		return action
	}

	if v, ok := m["custom_request_handling"].([]interface{}); ok && len(v) > 0 {
		action.CustomRequestHandling = expandCustomRequestHandling(v)
	}

	return action
}

func expandCustomResponseBodies(m []interface{}) map[string]*wafv2.CustomResponseBody {
	if len(m) == 0 {
		return nil
	}

	customResponseBodies := make(map[string]*wafv2.CustomResponseBody, len(m))

	for _, v := range m {
		vm := v.(map[string]interface{})
		key := vm["key"].(string)
		customResponseBodies[key] = &wafv2.CustomResponseBody{
			Content:     aws.String(vm["content"].(string)),
			ContentType: aws.String(vm["content_type"].(string)),
		}
	}

	return customResponseBodies
}

func expandCustomResponse(l []interface{}) *wafv2.CustomResponse {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}

	customResponse := &wafv2.CustomResponse{}

	if v, ok := m["custom_response_body_key"].(string); ok && v != "" {
		customResponse.CustomResponseBodyKey = aws.String(v)
	}
	if v, ok := m["response_code"].(int); ok && v > 0 {
		customResponse.ResponseCode = aws.Int64(int64(v))
	}
	if v, ok := m["response_header"].(*schema.Set); ok && len(v.List()) > 0 {
		customResponse.ResponseHeaders = expandCustomHeaders(v.List())
	}

	return customResponse
}

func expandCustomRequestHandling(l []interface{}) *wafv2.CustomRequestHandling {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	requestHandling := &wafv2.CustomRequestHandling{}

	if v, ok := m["insert_header"].(*schema.Set); ok && len(v.List()) > 0 {
		requestHandling.InsertHeaders = expandCustomHeaders(v.List())
	}

	return requestHandling
}

func expandCustomHeaders(l []interface{}) []*wafv2.CustomHTTPHeader {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	headers := make([]*wafv2.CustomHTTPHeader, 0)

	for _, header := range l {
		if header == nil {
			continue
		}
		m := header.(map[string]interface{})
		headers = append(headers, &wafv2.CustomHTTPHeader{
			Name:  aws.String(m["name"].(string)),
			Value: aws.String(m["value"].(string)),
		})
	}

	return headers
}

func expandVisibilityConfig(l []interface{}) *wafv2.VisibilityConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	configuration := &wafv2.VisibilityConfig{}

	if v, ok := m["cloudwatch_metrics_enabled"]; ok {
		configuration.CloudWatchMetricsEnabled = aws.Bool(v.(bool))
	}

	if v, ok := m["metric_name"]; ok && len(v.(string)) > 0 {
		configuration.MetricName = aws.String(v.(string))
	}

	if v, ok := m["sampled_requests_enabled"]; ok {
		configuration.SampledRequestsEnabled = aws.Bool(v.(bool))
	}

	return configuration
}

func expandRuleGroupRootStatement(l []interface{}) *wafv2.Statement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return expandStatement(m)
}

func expandStatements(l []interface{}) []*wafv2.Statement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	statements := make([]*wafv2.Statement, 0)

	for _, statement := range l {
		if statement == nil {
			continue
		}
		statements = append(statements, expandStatement(statement.(map[string]interface{})))
	}

	return statements
}

func expandStatement(m map[string]interface{}) *wafv2.Statement {
	if m == nil {
		return nil
	}

	statement := &wafv2.Statement{}

	if v, ok := m["and_statement"]; ok {
		statement.AndStatement = expandAndStatement(v.([]interface{}))
	}

	if v, ok := m["byte_match_statement"]; ok {
		statement.ByteMatchStatement = expandByteMatchStatement(v.([]interface{}))
	}

	if v, ok := m["ip_set_reference_statement"]; ok {
		statement.IPSetReferenceStatement = expandIPSetReferenceStatement(v.([]interface{}))
	}

	if v, ok := m["geo_match_statement"]; ok {
		statement.GeoMatchStatement = expandGeoMatchStatement(v.([]interface{}))
	}

	if v, ok := m["label_match_statement"]; ok {
		statement.LabelMatchStatement = expandLabelMatchStatement(v.([]interface{}))
	}

	if v, ok := m["not_statement"]; ok {
		statement.NotStatement = expandNotStatement(v.([]interface{}))
	}

	if v, ok := m["or_statement"]; ok {
		statement.OrStatement = expandOrStatement(v.([]interface{}))
	}

	if v, ok := m["rate_based_statement"]; ok {
		statement.RateBasedStatement = expandRateBasedStatement(v.([]interface{}))
	}

	if v, ok := m["regex_match_statement"]; ok {
		statement.RegexMatchStatement = expandRegexMatchStatement(v.([]interface{}))
	}

	if v, ok := m["regex_pattern_set_reference_statement"]; ok {
		statement.RegexPatternSetReferenceStatement = expandRegexPatternSetReferenceStatement(v.([]interface{}))
	}

	if v, ok := m["size_constraint_statement"]; ok {
		statement.SizeConstraintStatement = expandSizeConstraintStatement(v.([]interface{}))
	}

	if v, ok := m["sqli_match_statement"]; ok {
		statement.SqliMatchStatement = expandSQLiMatchStatement(v.([]interface{}))
	}

	if v, ok := m["xss_match_statement"]; ok {
		statement.XssMatchStatement = expandXSSMatchStatement(v.([]interface{}))
	}

	return statement
}

func expandAndStatement(l []interface{}) *wafv2.AndStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.AndStatement{
		Statements: expandStatements(m["statement"].([]interface{})),
	}
}

func expandByteMatchStatement(l []interface{}) *wafv2.ByteMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.ByteMatchStatement{
		FieldToMatch:         expandFieldToMatch(m["field_to_match"].([]interface{})),
		PositionalConstraint: aws.String(m["positional_constraint"].(string)),
		SearchString:         []byte(m["search_string"].(string)),
		TextTransformations:  expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandFieldToMatch(l []interface{}) *wafv2.FieldToMatch {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	f := &wafv2.FieldToMatch{}

	if v, ok := m["all_query_arguments"]; ok && len(v.([]interface{})) > 0 {
		f.AllQueryArguments = &wafv2.AllQueryArguments{}
	}

	if v, ok := m["body"]; ok && len(v.([]interface{})) > 0 {
		f.Body = &wafv2.Body{}
	}

	if v, ok := m["cookies"]; ok && len(v.([]interface{})) > 0 {
		f.Cookies = expandCookies(m["cookies"].([]interface{}))
	}

	if v, ok := m["headers"]; ok && len(v.([]interface{})) > 0 {
		f.Headers = expandHeaders(m["headers"].([]interface{}))
	}

	if v, ok := m["json_body"]; ok && len(v.([]interface{})) > 0 {
		f.JsonBody = expandJSONBody(v.([]interface{}))
	}

	if v, ok := m["method"]; ok && len(v.([]interface{})) > 0 {
		f.Method = &wafv2.Method{}
	}

	if v, ok := m["query_string"]; ok && len(v.([]interface{})) > 0 {
		f.QueryString = &wafv2.QueryString{}
	}

	if v, ok := m["single_header"]; ok && len(v.([]interface{})) > 0 {
		f.SingleHeader = expandSingleHeader(m["single_header"].([]interface{}))
	}

	if v, ok := m["single_query_argument"]; ok && len(v.([]interface{})) > 0 {
		f.SingleQueryArgument = expandSingleQueryArgument(m["single_query_argument"].([]interface{}))
	}

	if v, ok := m["uri_path"]; ok && len(v.([]interface{})) > 0 {
		f.UriPath = &wafv2.UriPath{}
	}

	return f
}

func expandForwardedIPConfig(l []interface{}) *wafv2.ForwardedIPConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.ForwardedIPConfig{
		FallbackBehavior: aws.String(m["fallback_behavior"].(string)),
		HeaderName:       aws.String(m["header_name"].(string)),
	}
}

func expandIPSetForwardedIPConfig(l []interface{}) *wafv2.IPSetForwardedIPConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.IPSetForwardedIPConfig{
		FallbackBehavior: aws.String(m["fallback_behavior"].(string)),
		HeaderName:       aws.String(m["header_name"].(string)),
		Position:         aws.String(m["position"].(string)),
	}
}

func expandCookies(l []interface{}) *wafv2.Cookies {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	cookies := &wafv2.Cookies{
		MatchScope:       aws.String(m["match_scope"].(string)),
		OversizeHandling: aws.String(m["oversize_handling"].(string)),
	}

	if v, ok := m["match_pattern"]; ok && len(v.([]interface{})) > 0 {
		cookies.MatchPattern = expandCookieMatchPattern(v.([]interface{}))
	}

	return cookies
}

func expandCookieMatchPattern(l []interface{}) *wafv2.CookieMatchPattern {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	CookieMatchPattern := &wafv2.CookieMatchPattern{}

	if v, ok := m["included_cookies"]; ok && len(v.([]interface{})) > 0 {
		CookieMatchPattern.IncludedCookies = flex.ExpandStringList(v.([]interface{}))
	}

	if v, ok := m["excluded_cookies"]; ok && len(v.([]interface{})) > 0 {
		CookieMatchPattern.ExcludedCookies = flex.ExpandStringList(v.([]interface{}))
	}

	if v, ok := m["all"].([]interface{}); ok && len(v) > 0 {
		CookieMatchPattern.All = &wafv2.All{}
	}

	return CookieMatchPattern
}

func expandJSONBody(l []interface{}) *wafv2.JsonBody {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	jsonBody := &wafv2.JsonBody{
		MatchScope:       aws.String(m["match_scope"].(string)),
		OversizeHandling: aws.String(m["oversize_handling"].(string)),
		MatchPattern:     expandJSONMatchPattern(m["match_pattern"].([]interface{})),
	}

	if v, ok := m["invalid_fallback_behavior"].(string); ok && v != "" {
		jsonBody.InvalidFallbackBehavior = aws.String(v)
	}

	return jsonBody
}

func expandJSONMatchPattern(l []interface{}) *wafv2.JsonMatchPattern {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	jsonMatchPattern := &wafv2.JsonMatchPattern{}

	if v, ok := m["all"].([]interface{}); ok && len(v) > 0 {
		jsonMatchPattern.All = &wafv2.All{}
	}

	if v, ok := m["included_paths"]; ok && len(v.([]interface{})) > 0 {
		jsonMatchPattern.IncludedPaths = flex.ExpandStringList(v.([]interface{}))
	}

	return jsonMatchPattern
}

func expandSingleHeader(l []interface{}) *wafv2.SingleHeader {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.SingleHeader{
		Name: aws.String(m["name"].(string)),
	}
}

func expandSingleQueryArgument(l []interface{}) *wafv2.SingleQueryArgument {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.SingleQueryArgument{
		Name: aws.String(m["name"].(string)),
	}
}

func expandTextTransformations(l []interface{}) []*wafv2.TextTransformation {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	rules := make([]*wafv2.TextTransformation, 0)

	for _, rule := range l {
		if rule == nil {
			continue
		}
		rules = append(rules, expandTextTransformation(rule.(map[string]interface{})))
	}

	return rules
}

func expandTextTransformation(m map[string]interface{}) *wafv2.TextTransformation {
	if m == nil {
		return nil
	}

	return &wafv2.TextTransformation{
		Priority: aws.Int64(int64(m["priority"].(int))),
		Type:     aws.String(m["type"].(string)),
	}
}

func expandIPSetReferenceStatement(l []interface{}) *wafv2.IPSetReferenceStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	statement := &wafv2.IPSetReferenceStatement{
		ARN: aws.String(m["arn"].(string)),
	}

	if v, ok := m["ip_set_forwarded_ip_config"]; ok {
		statement.IPSetForwardedIPConfig = expandIPSetForwardedIPConfig(v.([]interface{}))
	}

	return statement
}

func expandGeoMatchStatement(l []interface{}) *wafv2.GeoMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	statement := &wafv2.GeoMatchStatement{
		CountryCodes: flex.ExpandStringList(m["country_codes"].([]interface{})),
	}

	if v, ok := m["forwarded_ip_config"]; ok {
		statement.ForwardedIPConfig = expandForwardedIPConfig(v.([]interface{}))
	}

	return statement
}

func expandLabelMatchStatement(l []interface{}) *wafv2.LabelMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	statement := &wafv2.LabelMatchStatement{
		Key:   aws.String(m["key"].(string)),
		Scope: aws.String(m["scope"].(string)),
	}

	return statement
}

func expandNotStatement(l []interface{}) *wafv2.NotStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	s := m["statement"].([]interface{})

	if len(s) == 0 || s[0] == nil {
		return nil
	}

	m = s[0].(map[string]interface{})

	return &wafv2.NotStatement{
		Statement: expandStatement(m),
	}
}

func expandOrStatement(l []interface{}) *wafv2.OrStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.OrStatement{
		Statements: expandStatements(m["statement"].([]interface{})),
	}
}

func expandRegexMatchStatement(l []interface{}) *wafv2.RegexMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.RegexMatchStatement{
		RegexString:         aws.String(m["regex_string"].(string)),
		FieldToMatch:        expandFieldToMatch(m["field_to_match"].([]interface{})),
		TextTransformations: expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandRegexPatternSetReferenceStatement(l []interface{}) *wafv2.RegexPatternSetReferenceStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.RegexPatternSetReferenceStatement{
		ARN:                 aws.String(m["arn"].(string)),
		FieldToMatch:        expandFieldToMatch(m["field_to_match"].([]interface{})),
		TextTransformations: expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandSizeConstraintStatement(l []interface{}) *wafv2.SizeConstraintStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.SizeConstraintStatement{
		ComparisonOperator:  aws.String(m["comparison_operator"].(string)),
		FieldToMatch:        expandFieldToMatch(m["field_to_match"].([]interface{})),
		Size:                aws.Int64(int64(m["size"].(int))),
		TextTransformations: expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandSQLiMatchStatement(l []interface{}) *wafv2.SqliMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.SqliMatchStatement{
		FieldToMatch:        expandFieldToMatch(m["field_to_match"].([]interface{})),
		TextTransformations: expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandXSSMatchStatement(l []interface{}) *wafv2.XssMatchStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.XssMatchStatement{
		FieldToMatch:        expandFieldToMatch(m["field_to_match"].([]interface{})),
		TextTransformations: expandTextTransformations(m["text_transformation"].(*schema.Set).List()),
	}
}

func expandHeaders(l []interface{}) *wafv2.Headers {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.Headers{
		MatchPattern:     expandHeaderMatchPattern(m["match_pattern"].([]interface{})),
		MatchScope:       aws.String(m["match_scope"].(string)),
		OversizeHandling: aws.String(m["oversize_handling"].(string)),
	}
}

func expandHeaderMatchPattern(l []interface{}) *wafv2.HeaderMatchPattern {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	f := &wafv2.HeaderMatchPattern{}

	if v, ok := m["all"]; ok && len(v.([]interface{})) > 0 {
		f.All = &wafv2.All{}
	}

	if v, ok := m["included_headers"]; ok && len(v.([]interface{})) > 0 {
		f.IncludedHeaders = flex.ExpandStringList(m["included_headers"].([]interface{}))
	}

	if v, ok := m["excluded_headers"]; ok && len(v.([]interface{})) > 0 {
		f.ExcludedHeaders = flex.ExpandStringList(m["excluded_headers"].([]interface{}))
	}

	return f
}

func expandWebACLRules(l []interface{}) []*wafv2.Rule {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	rules := make([]*wafv2.Rule, 0)

	for _, rule := range l {
		if rule == nil {
			continue
		}
		rules = append(rules, expandWebACLRule(rule.(map[string]interface{})))
	}

	return rules
}

func expandWebACLRule(m map[string]interface{}) *wafv2.Rule {
	if m == nil {
		return nil
	}

	rule := &wafv2.Rule{
		Name:             aws.String(m["name"].(string)),
		Priority:         aws.Int64(int64(m["priority"].(int))),
		Action:           expandRuleAction(m["action"].([]interface{})),
		OverrideAction:   expandOverrideAction(m["override_action"].([]interface{})),
		Statement:        expandWebACLRootStatement(m["statement"].([]interface{})),
		VisibilityConfig: expandVisibilityConfig(m["visibility_config"].([]interface{})),
	}

	if v, ok := m["rule_label"].(*schema.Set); ok && v.Len() > 0 {
		rule.RuleLabels = expandRuleLabels(v.List())
	}

	return rule
}

func expandOverrideAction(l []interface{}) *wafv2.OverrideAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	action := &wafv2.OverrideAction{}

	if v, ok := m["count"]; ok && len(v.([]interface{})) > 0 {
		action.Count = &wafv2.CountAction{}
	}

	if v, ok := m["none"]; ok && len(v.([]interface{})) > 0 {
		action.None = &wafv2.NoneAction{}
	}

	return action
}

func expandDefaultAction(l []interface{}) *wafv2.DefaultAction {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	action := &wafv2.DefaultAction{}

	if v, ok := m["allow"]; ok && len(v.([]interface{})) > 0 {
		action.Allow = expandAllowAction(v.([]interface{}))
	}

	if v, ok := m["block"]; ok && len(v.([]interface{})) > 0 {
		action.Block = expandBlockAction(v.([]interface{}))
	}

	return action
}

func expandWebACLRootStatement(l []interface{}) *wafv2.Statement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return expandWebACLStatement(m)
}

func expandWebACLStatement(m map[string]interface{}) *wafv2.Statement {
	if m == nil {
		return nil
	}

	statement := &wafv2.Statement{}

	if v, ok := m["and_statement"]; ok {
		statement.AndStatement = expandAndStatement(v.([]interface{}))
	}

	if v, ok := m["byte_match_statement"]; ok {
		statement.ByteMatchStatement = expandByteMatchStatement(v.([]interface{}))
	}

	if v, ok := m["ip_set_reference_statement"]; ok {
		statement.IPSetReferenceStatement = expandIPSetReferenceStatement(v.([]interface{}))
	}

	if v, ok := m["geo_match_statement"]; ok {
		statement.GeoMatchStatement = expandGeoMatchStatement(v.([]interface{}))
	}

	if v, ok := m["label_match_statement"]; ok {
		statement.LabelMatchStatement = expandLabelMatchStatement(v.([]interface{}))
	}

	if v, ok := m["managed_rule_group_statement"]; ok {
		statement.ManagedRuleGroupStatement = expandManagedRuleGroupStatement(v.([]interface{}))
	}

	if v, ok := m["not_statement"]; ok {
		statement.NotStatement = expandNotStatement(v.([]interface{}))
	}

	if v, ok := m["or_statement"]; ok {
		statement.OrStatement = expandOrStatement(v.([]interface{}))
	}

	if v, ok := m["rate_based_statement"]; ok {
		statement.RateBasedStatement = expandRateBasedStatement(v.([]interface{}))
	}

	if v, ok := m["regex_match_statement"]; ok {
		statement.RegexMatchStatement = expandRegexMatchStatement(v.([]interface{}))
	}

	if v, ok := m["regex_pattern_set_reference_statement"]; ok {
		statement.RegexPatternSetReferenceStatement = expandRegexPatternSetReferenceStatement(v.([]interface{}))
	}

	if v, ok := m["rule_group_reference_statement"]; ok {
		statement.RuleGroupReferenceStatement = expandRuleGroupReferenceStatement(v.([]interface{}))
	}

	if v, ok := m["size_constraint_statement"]; ok {
		statement.SizeConstraintStatement = expandSizeConstraintStatement(v.([]interface{}))
	}

	if v, ok := m["sqli_match_statement"]; ok {
		statement.SqliMatchStatement = expandSQLiMatchStatement(v.([]interface{}))
	}

	if v, ok := m["xss_match_statement"]; ok {
		statement.XssMatchStatement = expandXSSMatchStatement(v.([]interface{}))
	}

	return statement
}

func expandManagedRuleGroupStatement(l []interface{}) *wafv2.ManagedRuleGroupStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	r := &wafv2.ManagedRuleGroupStatement{
		ExcludedRules: expandExcludedRules(m["excluded_rule"].([]interface{})),
		Name:          aws.String(m["name"].(string)),
		VendorName:    aws.String(m["vendor_name"].(string)),
	}

	if s, ok := m["scope_down_statement"].([]interface{}); ok && len(s) > 0 && s[0] != nil {
		r.ScopeDownStatement = expandStatement(s[0].(map[string]interface{}))
	}

	if v, ok := m["version"]; ok && v != "" {
		r.Version = aws.String(v.(string))
	}

	return r
}

func expandRateBasedStatement(l []interface{}) *wafv2.RateBasedStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	r := &wafv2.RateBasedStatement{
		AggregateKeyType: aws.String(m["aggregate_key_type"].(string)),
		Limit:            aws.Int64(int64(m["limit"].(int))),
	}

	if v, ok := m["forwarded_ip_config"]; ok {
		r.ForwardedIPConfig = expandForwardedIPConfig(v.([]interface{}))
	}

	s := m["scope_down_statement"].([]interface{})
	if len(s) > 0 && s[0] != nil {
		r.ScopeDownStatement = expandStatement(s[0].(map[string]interface{}))
	}

	return r
}

func expandRuleGroupReferenceStatement(l []interface{}) *wafv2.RuleGroupReferenceStatement {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	return &wafv2.RuleGroupReferenceStatement{
		ARN:           aws.String(m["arn"].(string)),
		ExcludedRules: expandExcludedRules(m["excluded_rule"].([]interface{})),
	}
}

func expandExcludedRules(l []interface{}) []*wafv2.ExcludedRule {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	rules := make([]*wafv2.ExcludedRule, 0)

	for _, rule := range l {
		if rule == nil {
			continue
		}
		rules = append(rules, expandExcludedRule(rule.(map[string]interface{})))
	}

	return rules
}

func expandExcludedRule(m map[string]interface{}) *wafv2.ExcludedRule {
	if m == nil {
		return nil
	}

	return &wafv2.ExcludedRule{
		Name: aws.String(m["name"].(string)),
	}
}

func expandRegexPatternSet(l []interface{}) []*wafv2.Regex {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	regexPatterns := make([]*wafv2.Regex, 0)
	for _, regexPattern := range l {
		if regexPattern == nil {
			continue
		}
		regexPatterns = append(regexPatterns, expandRegex(regexPattern.(map[string]interface{})))
	}

	return regexPatterns
}

func expandRegex(m map[string]interface{}) *wafv2.Regex {
	if m == nil {
		return nil
	}

	return &wafv2.Regex{
		RegexString: aws.String(m["regex_string"].(string)),
	}
}

func flattenRules(r []*wafv2.Rule) interface{} {
	out := make([]map[string]interface{}, len(r))
	for i, rule := range r {
		m := make(map[string]interface{})
		m["action"] = flattenRuleAction(rule.Action)
		m["name"] = aws.StringValue(rule.Name)
		m["priority"] = int(aws.Int64Value(rule.Priority))
		m["rule_label"] = flattenRuleLabels(rule.RuleLabels)
		m["statement"] = flattenRuleGroupRootStatement(rule.Statement)
		m["visibility_config"] = flattenVisibilityConfig(rule.VisibilityConfig)
		out[i] = m
	}

	return out
}

func flattenRuleAction(a *wafv2.RuleAction) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if a.Allow != nil {
		m["allow"] = flattenAllow(a.Allow)
	}

	if a.Block != nil {
		m["block"] = flattenBlock(a.Block)
	}

	if a.Captcha != nil {
		m["captcha"] = flattenCaptcha(a.Captcha)
	}

	if a.Count != nil {
		m["count"] = flattenCount(a.Count)
	}

	return []interface{}{m}
}

func flattenAllow(a *wafv2.AllowAction) []interface{} {
	if a == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{}

	if a.CustomRequestHandling != nil {
		m["custom_request_handling"] = flattenCustomRequestHandling(a.CustomRequestHandling)
	}

	return []interface{}{m}
}

func flattenBlock(a *wafv2.BlockAction) []interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if a.CustomResponse != nil {
		m["custom_response"] = flattenCustomResponse(a.CustomResponse)
	}

	return []interface{}{m}
}

func flattenCaptcha(a *wafv2.CaptchaAction) []interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if a.CustomRequestHandling != nil {
		m["custom_request_handling"] = flattenCustomRequestHandling(a.CustomRequestHandling)
	}

	return []interface{}{m}
}

func flattenCount(a *wafv2.CountAction) []interface{} {
	if a == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{}

	if a.CustomRequestHandling != nil {
		m["custom_request_handling"] = flattenCustomRequestHandling(a.CustomRequestHandling)
	}

	return []interface{}{m}
}

func flattenCustomResponseBodies(b map[string]*wafv2.CustomResponseBody) interface{} {
	if len(b) == 0 {
		return make([]map[string]interface{}, 0)
	}

	out := make([]map[string]interface{}, len(b))
	i := 0
	for key, body := range b {
		out[i] = map[string]interface{}{
			"key":          key,
			"content":      aws.StringValue(body.Content),
			"content_type": aws.StringValue(body.ContentType),
		}
		i += 1
	}

	return out
}

func flattenCustomRequestHandling(c *wafv2.CustomRequestHandling) []interface{} {
	if c == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"insert_header": flattenCustomHeaders(c.InsertHeaders),
	}

	return []interface{}{m}
}

func flattenCustomResponse(r *wafv2.CustomResponse) []interface{} {
	if r == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"response_code":   int(aws.Int64Value(r.ResponseCode)),
		"response_header": flattenCustomHeaders(r.ResponseHeaders),
	}

	if r.CustomResponseBodyKey != nil {
		m["custom_response_body_key"] = aws.StringValue(r.CustomResponseBodyKey)
	}

	return []interface{}{m}
}

func flattenCustomHeaders(h []*wafv2.CustomHTTPHeader) []interface{} {
	out := make([]interface{}, len(h))
	for i, header := range h {
		out[i] = flattenCustomHeader(header)
	}

	return out
}

func flattenCustomHeader(h *wafv2.CustomHTTPHeader) map[string]interface{} {
	if h == nil {
		return map[string]interface{}{}
	}

	m := map[string]interface{}{
		"name":  aws.StringValue(h.Name),
		"value": aws.StringValue(h.Value),
	}

	return m
}

func flattenRuleLabels(l []*wafv2.Label) []interface{} {
	if len(l) == 0 {
		return nil
	}

	out := make([]interface{}, len(l))
	for i, label := range l {
		out[i] = map[string]interface{}{
			"name": aws.StringValue(label.Name),
		}
	}

	return out
}

func flattenRuleGroupRootStatement(s *wafv2.Statement) interface{} {
	if s == nil {
		return []interface{}{}
	}

	return []interface{}{flattenStatement(s)}
}

func flattenStatements(s []*wafv2.Statement) interface{} {
	out := make([]interface{}, len(s))
	for i, statement := range s {
		out[i] = flattenStatement(statement)
	}

	return out
}

func flattenStatement(s *wafv2.Statement) map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}

	m := map[string]interface{}{}

	if s.AndStatement != nil {
		m["and_statement"] = flattenAndStatement(s.AndStatement)
	}

	if s.ByteMatchStatement != nil {
		m["byte_match_statement"] = flattenByteMatchStatement(s.ByteMatchStatement)
	}

	if s.IPSetReferenceStatement != nil {
		m["ip_set_reference_statement"] = flattenIPSetReferenceStatement(s.IPSetReferenceStatement)
	}

	if s.GeoMatchStatement != nil {
		m["geo_match_statement"] = flattenGeoMatchStatement(s.GeoMatchStatement)
	}

	if s.LabelMatchStatement != nil {
		m["label_match_statement"] = flattenLabelMatchStatement(s.LabelMatchStatement)
	}

	if s.NotStatement != nil {
		m["not_statement"] = flattenNotStatement(s.NotStatement)
	}

	if s.OrStatement != nil {
		m["or_statement"] = flattenOrStatement(s.OrStatement)
	}

	if s.RateBasedStatement != nil {
		m["rate_based_statement"] = flattenRateBasedStatement(s.RateBasedStatement)
	}

	if s.RegexMatchStatement != nil {
		m["regex_match_statement"] = flattenRegexMatchStatement(s.RegexMatchStatement)
	}

	if s.RegexPatternSetReferenceStatement != nil {
		m["regex_pattern_set_reference_statement"] = flattenRegexPatternSetReferenceStatement(s.RegexPatternSetReferenceStatement)
	}

	if s.SizeConstraintStatement != nil {
		m["size_constraint_statement"] = flattenSizeConstraintStatement(s.SizeConstraintStatement)
	}

	if s.SqliMatchStatement != nil {
		m["sqli_match_statement"] = flattenSQLiMatchStatement(s.SqliMatchStatement)
	}

	if s.XssMatchStatement != nil {
		m["xss_match_statement"] = flattenXSSMatchStatement(s.XssMatchStatement)
	}

	return m
}

func flattenAndStatement(a *wafv2.AndStatement) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"statement": flattenStatements(a.Statements),
	}

	return []interface{}{m}
}

func flattenByteMatchStatement(b *wafv2.ByteMatchStatement) interface{} {
	if b == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"field_to_match":        flattenFieldToMatch(b.FieldToMatch),
		"positional_constraint": aws.StringValue(b.PositionalConstraint),
		"search_string":         string(b.SearchString),
		"text_transformation":   flattenTextTransformations(b.TextTransformations),
	}

	return []interface{}{m}
}

func flattenFieldToMatch(f *wafv2.FieldToMatch) interface{} {
	if f == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if f.AllQueryArguments != nil {
		m["all_query_arguments"] = make([]map[string]interface{}, 1)
	}

	if f.Body != nil {
		m["body"] = make([]map[string]interface{}, 1)
	}

	if f.Cookies != nil {
		m["cookies"] = flattenCookies(f.Cookies)
	}

	if f.Headers != nil {
		m["headers"] = flattenHeaders(f.Headers)
	}

	if f.JsonBody != nil {
		m["json_body"] = flattenJSONBody(f.JsonBody)
	}

	if f.Method != nil {
		m["method"] = make([]map[string]interface{}, 1)
	}

	if f.QueryString != nil {
		m["query_string"] = make([]map[string]interface{}, 1)
	}

	if f.SingleHeader != nil {
		m["single_header"] = flattenSingleHeader(f.SingleHeader)
	}

	if f.SingleQueryArgument != nil {
		m["single_query_argument"] = flattenSingleQueryArgument(f.SingleQueryArgument)
	}

	if f.UriPath != nil {
		m["uri_path"] = make([]map[string]interface{}, 1)
	}

	return []interface{}{m}
}

func flattenForwardedIPConfig(f *wafv2.ForwardedIPConfig) interface{} {
	if f == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"fallback_behavior": aws.StringValue(f.FallbackBehavior),
		"header_name":       aws.StringValue(f.HeaderName),
	}

	return []interface{}{m}
}

func flattenIPSetForwardedIPConfig(i *wafv2.IPSetForwardedIPConfig) interface{} {
	if i == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"fallback_behavior": aws.StringValue(i.FallbackBehavior),
		"header_name":       aws.StringValue(i.HeaderName),
		"position":          aws.StringValue(i.Position),
	}

	return []interface{}{m}
}

func flattenCookies(c *wafv2.Cookies) interface{} {
	if c == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"match_scope":       aws.StringValue(c.MatchScope),
		"oversize_handling": aws.StringValue(c.OversizeHandling),
		"match_pattern":     flattenCookiesMatchPattern(c.MatchPattern),
	}

	return []interface{}{m}
}

func flattenCookiesMatchPattern(c *wafv2.CookieMatchPattern) interface{} {
	if c == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"included_cookies": aws.StringValueSlice(c.IncludedCookies),
		"excluded_cookies": aws.StringValueSlice(c.ExcludedCookies),
	}

	if c.All != nil {
		m["all"] = make([]map[string]interface{}, 1)
	}

	return []interface{}{m}
}

func flattenJSONBody(b *wafv2.JsonBody) interface{} {
	if b == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"invalid_fallback_behavior": aws.StringValue(b.InvalidFallbackBehavior),
		"match_pattern":             flattenJSONMatchPattern(b.MatchPattern),
		"match_scope":               aws.StringValue(b.MatchScope),
		"oversize_handling":         aws.StringValue(b.OversizeHandling),
	}

	return []interface{}{m}
}

func flattenJSONMatchPattern(p *wafv2.JsonMatchPattern) []interface{} {
	if p == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"included_paths": flex.FlattenStringList(p.IncludedPaths),
	}

	if p.All != nil {
		m["all"] = make([]map[string]interface{}, 1)
	}

	return []interface{}{m}
}

func flattenSingleHeader(s *wafv2.SingleHeader) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"name": aws.StringValue(s.Name),
	}

	return []interface{}{m}
}

func flattenSingleQueryArgument(s *wafv2.SingleQueryArgument) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"name": aws.StringValue(s.Name),
	}

	return []interface{}{m}
}

func flattenTextTransformations(l []*wafv2.TextTransformation) []interface{} {
	out := make([]interface{}, len(l))
	for i, t := range l {
		m := make(map[string]interface{})
		m["priority"] = int(aws.Int64Value(t.Priority))
		m["type"] = aws.StringValue(t.Type)
		out[i] = m
	}
	return out
}

func flattenIPSetReferenceStatement(i *wafv2.IPSetReferenceStatement) interface{} {
	if i == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"arn":                        aws.StringValue(i.ARN),
		"ip_set_forwarded_ip_config": flattenIPSetForwardedIPConfig(i.IPSetForwardedIPConfig),
	}

	return []interface{}{m}
}

func flattenGeoMatchStatement(g *wafv2.GeoMatchStatement) interface{} {
	if g == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"country_codes":       flex.FlattenStringList(g.CountryCodes),
		"forwarded_ip_config": flattenForwardedIPConfig(g.ForwardedIPConfig),
	}

	return []interface{}{m}
}

func flattenLabelMatchStatement(l *wafv2.LabelMatchStatement) interface{} {
	if l == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"key":   aws.StringValue(l.Key),
		"scope": aws.StringValue(l.Scope),
	}

	return []interface{}{m}
}

func flattenNotStatement(a *wafv2.NotStatement) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"statement": []interface{}{flattenStatement(a.Statement)},
	}

	return []interface{}{m}
}

func flattenOrStatement(a *wafv2.OrStatement) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"statement": flattenStatements(a.Statements),
	}

	return []interface{}{m}
}

func flattenRegexMatchStatement(r *wafv2.RegexMatchStatement) interface{} {
	if r == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"regex_string":        aws.StringValue(r.RegexString),
		"field_to_match":      flattenFieldToMatch(r.FieldToMatch),
		"text_transformation": flattenTextTransformations(r.TextTransformations),
	}

	return []interface{}{m}
}

func flattenRegexPatternSetReferenceStatement(r *wafv2.RegexPatternSetReferenceStatement) interface{} {
	if r == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"arn":                 aws.StringValue(r.ARN),
		"field_to_match":      flattenFieldToMatch(r.FieldToMatch),
		"text_transformation": flattenTextTransformations(r.TextTransformations),
	}

	return []interface{}{m}
}

func flattenSizeConstraintStatement(s *wafv2.SizeConstraintStatement) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"comparison_operator": aws.StringValue(s.ComparisonOperator),
		"field_to_match":      flattenFieldToMatch(s.FieldToMatch),
		"size":                int(aws.Int64Value(s.Size)),
		"text_transformation": flattenTextTransformations(s.TextTransformations),
	}

	return []interface{}{m}
}

func flattenSQLiMatchStatement(s *wafv2.SqliMatchStatement) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"field_to_match":      flattenFieldToMatch(s.FieldToMatch),
		"text_transformation": flattenTextTransformations(s.TextTransformations),
	}

	return []interface{}{m}
}

func flattenXSSMatchStatement(s *wafv2.XssMatchStatement) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"field_to_match":      flattenFieldToMatch(s.FieldToMatch),
		"text_transformation": flattenTextTransformations(s.TextTransformations),
	}

	return []interface{}{m}
}

func flattenVisibilityConfig(config *wafv2.VisibilityConfig) interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"cloudwatch_metrics_enabled": aws.BoolValue(config.CloudWatchMetricsEnabled),
		"metric_name":                aws.StringValue(config.MetricName),
		"sampled_requests_enabled":   aws.BoolValue(config.SampledRequestsEnabled),
	}

	return []interface{}{m}
}

func flattenHeaders(s *wafv2.Headers) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"match_scope":       aws.StringValue(s.MatchScope),
		"match_pattern":     flattenHeaderMatchPattern(s.MatchPattern),
		"oversize_handling": aws.StringValue(s.OversizeHandling),
	}

	return []interface{}{m}
}

func flattenHeaderMatchPattern(s *wafv2.HeaderMatchPattern) interface{} {
	if s == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if s.All != nil {
		m["all"] = make([]map[string]interface{}, 1)
	}

	if s.ExcludedHeaders != nil {
		m["excluded_headers"] = flex.FlattenStringList(s.ExcludedHeaders)
	}

	if s.IncludedHeaders != nil {
		m["included_headers"] = flex.FlattenStringList(s.IncludedHeaders)
	}

	return []interface{}{m}
}

func flattenWebACLRootStatement(s *wafv2.Statement) interface{} {
	if s == nil {
		return []interface{}{}
	}

	return []interface{}{flattenWebACLStatement(s)}
}

func flattenWebACLStatement(s *wafv2.Statement) map[string]interface{} {
	if s == nil {
		return map[string]interface{}{}
	}

	m := map[string]interface{}{}

	if s.AndStatement != nil {
		m["and_statement"] = flattenAndStatement(s.AndStatement)
	}

	if s.ByteMatchStatement != nil {
		m["byte_match_statement"] = flattenByteMatchStatement(s.ByteMatchStatement)
	}

	if s.IPSetReferenceStatement != nil {
		m["ip_set_reference_statement"] = flattenIPSetReferenceStatement(s.IPSetReferenceStatement)
	}

	if s.GeoMatchStatement != nil {
		m["geo_match_statement"] = flattenGeoMatchStatement(s.GeoMatchStatement)
	}

	if s.LabelMatchStatement != nil {
		m["label_match_statement"] = flattenLabelMatchStatement(s.LabelMatchStatement)
	}

	if s.ManagedRuleGroupStatement != nil {
		m["managed_rule_group_statement"] = flattenManagedRuleGroupStatement(s.ManagedRuleGroupStatement)
	}

	if s.NotStatement != nil {
		m["not_statement"] = flattenNotStatement(s.NotStatement)
	}

	if s.OrStatement != nil {
		m["or_statement"] = flattenOrStatement(s.OrStatement)
	}

	if s.RateBasedStatement != nil {
		m["rate_based_statement"] = flattenRateBasedStatement(s.RateBasedStatement)
	}

	if s.RegexMatchStatement != nil {
		m["regex_match_statement"] = flattenRegexMatchStatement(s.RegexMatchStatement)
	}

	if s.RegexPatternSetReferenceStatement != nil {
		m["regex_pattern_set_reference_statement"] = flattenRegexPatternSetReferenceStatement(s.RegexPatternSetReferenceStatement)
	}

	if s.RuleGroupReferenceStatement != nil {
		m["rule_group_reference_statement"] = flattenRuleGroupReferenceStatement(s.RuleGroupReferenceStatement)
	}

	if s.SizeConstraintStatement != nil {
		m["size_constraint_statement"] = flattenSizeConstraintStatement(s.SizeConstraintStatement)
	}

	if s.SqliMatchStatement != nil {
		m["sqli_match_statement"] = flattenSQLiMatchStatement(s.SqliMatchStatement)
	}

	if s.XssMatchStatement != nil {
		m["xss_match_statement"] = flattenXSSMatchStatement(s.XssMatchStatement)
	}

	return m
}

func flattenWebACLRules(r []*wafv2.Rule) interface{} {
	out := make([]map[string]interface{}, len(r))
	for i, rule := range r {
		m := make(map[string]interface{})
		m["action"] = flattenRuleAction(rule.Action)
		m["override_action"] = flattenOverrideAction(rule.OverrideAction)
		m["name"] = aws.StringValue(rule.Name)
		m["priority"] = int(aws.Int64Value(rule.Priority))
		m["rule_label"] = flattenRuleLabels(rule.RuleLabels)
		m["statement"] = flattenWebACLRootStatement(rule.Statement)
		m["visibility_config"] = flattenVisibilityConfig(rule.VisibilityConfig)
		out[i] = m
	}

	return out
}

func flattenOverrideAction(a *wafv2.OverrideAction) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if a.Count != nil {
		m["count"] = make([]map[string]interface{}, 1)
	}

	if a.None != nil {
		m["none"] = make([]map[string]interface{}, 1)
	}

	return []interface{}{m}
}

func flattenDefaultAction(a *wafv2.DefaultAction) interface{} {
	if a == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if a.Allow != nil {
		m["allow"] = flattenAllow(a.Allow)
	}

	if a.Block != nil {
		m["block"] = flattenBlock(a.Block)
	}

	return []interface{}{m}
}

func flattenManagedRuleGroupStatement(apiObject *wafv2.ManagedRuleGroupStatement) interface{} {
	if apiObject == nil {
		return []interface{}{}
	}

	tfMap := map[string]interface{}{}

	if apiObject.ExcludedRules != nil {
		tfMap["excluded_rule"] = flattenExcludedRules(apiObject.ExcludedRules)
	}

	if apiObject.Name != nil {
		tfMap["name"] = aws.StringValue(apiObject.Name)
	}

	if apiObject.ScopeDownStatement != nil {
		tfMap["scope_down_statement"] = []interface{}{flattenStatement(apiObject.ScopeDownStatement)}
	}

	if apiObject.VendorName != nil {
		tfMap["vendor_name"] = aws.StringValue(apiObject.VendorName)
	}

	if apiObject.Version != nil {
		tfMap["version"] = aws.StringValue(apiObject.Version)
	}

	return []interface{}{tfMap}
}

func flattenRateBasedStatement(apiObject *wafv2.RateBasedStatement) interface{} {
	if apiObject == nil {
		return []interface{}{}
	}

	tfMap := map[string]interface{}{}

	if apiObject.AggregateKeyType != nil {
		tfMap["aggregate_key_type"] = aws.StringValue(apiObject.AggregateKeyType)
	}

	if apiObject.ForwardedIPConfig != nil {
		tfMap["forwarded_ip_config"] = flattenForwardedIPConfig(apiObject.ForwardedIPConfig)
	}

	if apiObject.Limit != nil {
		tfMap["limit"] = int(aws.Int64Value(apiObject.Limit))
	}

	if apiObject.ScopeDownStatement != nil {
		tfMap["scope_down_statement"] = []interface{}{flattenStatement(apiObject.ScopeDownStatement)}
	}

	return []interface{}{tfMap}
}

func flattenRuleGroupReferenceStatement(r *wafv2.RuleGroupReferenceStatement) interface{} {
	if r == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"excluded_rule": flattenExcludedRules(r.ExcludedRules),
		"arn":           aws.StringValue(r.ARN),
	}

	return []interface{}{m}
}

func flattenExcludedRules(r []*wafv2.ExcludedRule) interface{} {
	out := make([]map[string]interface{}, len(r))
	for i, rule := range r {
		m := make(map[string]interface{})
		m["name"] = aws.StringValue(rule.Name)
		out[i] = m
	}

	return out
}

func flattenRegexPatternSet(r []*wafv2.Regex) interface{} {
	if r == nil {
		return []interface{}{}
	}

	regexPatterns := make([]interface{}, 0)

	for _, regexPattern := range r {
		if regexPattern == nil {
			continue
		}
		d := map[string]interface{}{
			"regex_string": aws.StringValue(regexPattern.RegexString),
		}
		regexPatterns = append(regexPatterns, d)
	}

	return regexPatterns
}
