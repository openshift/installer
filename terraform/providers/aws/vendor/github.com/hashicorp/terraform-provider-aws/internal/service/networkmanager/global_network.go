package networkmanager

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceGlobalNetwork() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGlobalNetworkCreate,
		ReadWithoutTimeout:   resourceGlobalNetworkRead,
		UpdateWithoutTimeout: resourceGlobalNetworkUpdate,
		DeleteWithoutTimeout: resourceGlobalNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
	}
}

func resourceGlobalNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &networkmanager.CreateGlobalNetworkInput{}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating Network Manager Global Network: %s", input)
	output, err := conn.CreateGlobalNetworkWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("error creating Network Manager Global Network: %s", err)
	}

	d.SetId(aws.StringValue(output.GlobalNetwork.GlobalNetworkId))

	if _, err := waitGlobalNetworkCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for Network Manager Global Network (%s) create: %s", d.Id(), err)
	}

	return resourceGlobalNetworkRead(ctx, d, meta)
}

func resourceGlobalNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	globalNetwork, err := FindGlobalNetworkByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Network Manager Global Network %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading Network Manager Global Network (%s): %s", d.Id(), err)
	}

	d.Set("arn", globalNetwork.GlobalNetworkArn)
	d.Set("description", globalNetwork.Description)

	tags := KeyValueTags(globalNetwork.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.Errorf("error setting tags: %s", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.Errorf("error setting tags_all: %s", err)
	}

	return nil
}

func resourceGlobalNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn

	if d.HasChangesExcept("tags", "tags_all") {
		input := &networkmanager.UpdateGlobalNetworkInput{
			Description:     aws.String(d.Get("description").(string)),
			GlobalNetworkId: aws.String(d.Id()),
		}

		log.Printf("[DEBUG] Updating Network Manager Global Network: %s", input)
		_, err := conn.UpdateGlobalNetworkWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("error updating Network Manager Global Network (%s): %s", d.Id(), err)
		}

		if _, err := waitGlobalNetworkUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for Network Manager Global Network (%s) update: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTagsWithContext(ctx, conn, d.Get("arn").(string), o, n); err != nil {
			return diag.Errorf("error updating Network Manager Global Network (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceGlobalNetworkRead(ctx, d, meta)
}

func resourceGlobalNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn

	if diags := disassociateCustomerGateways(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); diags.HasError() {
		return diags
	}

	if diags := disassociateTransitGatewayConnectPeers(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); diags.HasError() {
		return diags
	}

	if diags := deregisterTransitGateways(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); diags.HasError() {
		return diags
	}

	log.Printf("[DEBUG] Deleting Network Manager Global Network: %s", d.Id())
	_, err := tfresource.RetryWhenContext(ctx, globalNetworkValidationExceptionTimeout,
		func() (interface{}, error) {
			return conn.DeleteGlobalNetworkWithContext(ctx, &networkmanager.DeleteGlobalNetworkInput{
				GlobalNetworkId: aws.String(d.Id()),
			})
		},
		func(err error) (bool, error) {
			if tfawserr.ErrMessageContains(err, networkmanager.ErrCodeValidationException, "cannot be deleted due to existing devices, sites, or links") {
				return true, err
			}

			return false, err
		},
	)

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("error deleting Network Manager Global Network (%s): %s", d.Id(), err)
	}

	if _, err := waitGlobalNetworkDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for Network Manager Global Network (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func deregisterTransitGateways(ctx context.Context, conn *networkmanager.NetworkManager, globalNetworkID string, timeout time.Duration) diag.Diagnostics {
	output, err := FindTransitGatewayRegistrations(ctx, conn, &networkmanager.GetTransitGatewayRegistrationsInput{
		GlobalNetworkId: aws.String(globalNetworkID),
	})

	if tfresource.NotFound(err) {
		err = nil
	}

	if err != nil {
		return diag.Errorf("error listing Network Manager Transit Gateway Registrations (%s): %s", globalNetworkID, err)
	}

	var diags diag.Diagnostics

	for _, v := range output {
		err := deregisterTransitGateway(ctx, conn, globalNetworkID, aws.StringValue(v.TransitGatewayArn), timeout)

		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	return nil
}

func disassociateCustomerGateways(ctx context.Context, conn *networkmanager.NetworkManager, globalNetworkID string, timeout time.Duration) diag.Diagnostics {
	output, err := FindCustomerGatewayAssociations(ctx, conn, &networkmanager.GetCustomerGatewayAssociationsInput{
		GlobalNetworkId: aws.String(globalNetworkID),
	})

	if tfresource.NotFound(err) {
		err = nil
	}

	if err != nil {
		return diag.Errorf("error listing Network Manager Customer Gateway Associations (%s): %s", globalNetworkID, err)
	}

	var diags diag.Diagnostics

	for _, v := range output {
		err := disassociateCustomerGateway(ctx, conn, globalNetworkID, aws.StringValue(v.CustomerGatewayArn), timeout)

		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	return nil
}

func disassociateTransitGatewayConnectPeers(ctx context.Context, conn *networkmanager.NetworkManager, globalNetworkID string, timeout time.Duration) diag.Diagnostics {
	output, err := FindTransitGatewayConnectPeerAssociations(ctx, conn, &networkmanager.GetTransitGatewayConnectPeerAssociationsInput{
		GlobalNetworkId: aws.String(globalNetworkID),
	})

	if tfresource.NotFound(err) {
		err = nil
	}

	if err != nil {
		return diag.Errorf("error listing Network Manager Transit Gateway Connect Peer Associations (%s): %s", globalNetworkID, err)
	}

	var diags diag.Diagnostics

	for _, v := range output {
		err := disassociateTransitGatewayConnectPeer(ctx, conn, globalNetworkID, aws.StringValue(v.TransitGatewayConnectPeerArn), timeout)

		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	if diags.HasError() {
		return diags
	}

	return nil
}

func globalNetworkIDNotFoundError(err error) bool {
	return validationExceptionMessageContains(err, networkmanager.ValidationExceptionReasonFieldValidationFailed, "Global network not found")
}

func FindGlobalNetwork(ctx context.Context, conn *networkmanager.NetworkManager, input *networkmanager.DescribeGlobalNetworksInput) (*networkmanager.GlobalNetwork, error) {
	output, err := FindGlobalNetworks(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	if len(output) == 0 || output[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output[0], nil
}

func FindGlobalNetworks(ctx context.Context, conn *networkmanager.NetworkManager, input *networkmanager.DescribeGlobalNetworksInput) ([]*networkmanager.GlobalNetwork, error) {
	var output []*networkmanager.GlobalNetwork

	err := conn.DescribeGlobalNetworksPagesWithContext(ctx, input, func(page *networkmanager.DescribeGlobalNetworksOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.GlobalNetworks {
			if v == nil {
				continue
			}

			output = append(output, v)
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

func FindGlobalNetworkByID(ctx context.Context, conn *networkmanager.NetworkManager, id string) (*networkmanager.GlobalNetwork, error) {
	input := &networkmanager.DescribeGlobalNetworksInput{
		GlobalNetworkIds: aws.StringSlice([]string{id}),
	}

	output, err := FindGlobalNetwork(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	// Eventual consistency check.
	if aws.StringValue(output.GlobalNetworkId) != id {
		return nil, &resource.NotFoundError{
			LastRequest: input,
		}
	}

	return output, nil
}

func statusGlobalNetworkState(ctx context.Context, conn *networkmanager.NetworkManager, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindGlobalNetworkByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.State), nil
	}
}

func waitGlobalNetworkCreated(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.GlobalNetwork, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{networkmanager.GlobalNetworkStatePending},
		Target:  []string{networkmanager.GlobalNetworkStateAvailable},
		Timeout: timeout,
		Refresh: statusGlobalNetworkState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.GlobalNetwork); ok {
		return output, err
	}

	return nil, err
}

func waitGlobalNetworkDeleted(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.GlobalNetwork, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{networkmanager.GlobalNetworkStateDeleting},
		Target:         []string{},
		Timeout:        timeout,
		Refresh:        statusGlobalNetworkState(ctx, conn, id),
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.GlobalNetwork); ok {
		return output, err
	}

	return nil, err
}

func waitGlobalNetworkUpdated(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.GlobalNetwork, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{networkmanager.GlobalNetworkStateUpdating},
		Target:  []string{networkmanager.GlobalNetworkStateAvailable},
		Timeout: timeout,
		Refresh: statusGlobalNetworkState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.GlobalNetwork); ok {
		return output, err
	}

	return nil, err
}

const (
	globalNetworkValidationExceptionTimeout = 2 * time.Minute
)
