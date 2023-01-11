package apigatewayv2

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

func ResourceRouteResponse() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouteResponseCreate,
		Read:   resourceRouteResponseRead,
		Update: resourceRouteResponseUpdate,
		Delete: resourceRouteResponseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRouteResponseImport,
		},

		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model_selection_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"response_models": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_response_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRouteResponseCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn

	req := &apigatewayv2.CreateRouteResponseInput{
		ApiId:            aws.String(d.Get("api_id").(string)),
		RouteId:          aws.String(d.Get("route_id").(string)),
		RouteResponseKey: aws.String(d.Get("route_response_key").(string)),
	}
	if v, ok := d.GetOk("model_selection_expression"); ok {
		req.ModelSelectionExpression = aws.String(v.(string))
	}
	if v, ok := d.GetOk("response_models"); ok {
		req.ResponseModels = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	log.Printf("[DEBUG] Creating API Gateway v2 route response: %s", req)
	resp, err := conn.CreateRouteResponse(req)
	if err != nil {
		return fmt.Errorf("creating API Gateway v2 route response: %s", err)
	}

	d.SetId(aws.StringValue(resp.RouteResponseId))

	return resourceRouteResponseRead(d, meta)
}

func resourceRouteResponseRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn

	resp, err := conn.GetRouteResponse(&apigatewayv2.GetRouteResponseInput{
		ApiId:           aws.String(d.Get("api_id").(string)),
		RouteId:         aws.String(d.Get("route_id").(string)),
		RouteResponseId: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) && !d.IsNewResource() {
		log.Printf("[WARN] API Gateway v2 route response (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading API Gateway v2 route response: %s", err)
	}

	d.Set("model_selection_expression", resp.ModelSelectionExpression)
	if err := d.Set("response_models", flex.PointersMapToStringList(resp.ResponseModels)); err != nil {
		return fmt.Errorf("setting response_models: %s", err)
	}
	d.Set("route_response_key", resp.RouteResponseKey)

	return nil
}

func resourceRouteResponseUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn

	req := &apigatewayv2.UpdateRouteResponseInput{
		ApiId:           aws.String(d.Get("api_id").(string)),
		RouteId:         aws.String(d.Get("route_id").(string)),
		RouteResponseId: aws.String(d.Id()),
	}
	if d.HasChange("model_selection_expression") {
		req.ModelSelectionExpression = aws.String(d.Get("model_selection_expression").(string))
	}
	if d.HasChange("response_models") {
		req.ResponseModels = flex.ExpandStringMap(d.Get("response_models").(map[string]interface{}))
	}
	if d.HasChange("route_response_key") {
		req.RouteResponseKey = aws.String(d.Get("route_response_key").(string))
	}

	log.Printf("[DEBUG] Updating API Gateway v2 route response: %s", req)
	_, err := conn.UpdateRouteResponse(req)
	if err != nil {
		return fmt.Errorf("updating API Gateway v2 route response: %s", err)
	}

	return resourceRouteResponseRead(d, meta)
}

func resourceRouteResponseDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn

	log.Printf("[DEBUG] Deleting API Gateway v2 route response (%s)", d.Id())
	_, err := conn.DeleteRouteResponse(&apigatewayv2.DeleteRouteResponseInput{
		ApiId:           aws.String(d.Get("api_id").(string)),
		RouteId:         aws.String(d.Get("route_id").(string)),
		RouteResponseId: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("deleting API Gateway v2 route response: %s", err)
	}

	return nil
}

func resourceRouteResponseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of import ID (%s), use: 'api-id/route-id/route-response-id'", d.Id())
	}

	d.SetId(parts[2])
	d.Set("api_id", parts[0])
	d.Set("route_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
