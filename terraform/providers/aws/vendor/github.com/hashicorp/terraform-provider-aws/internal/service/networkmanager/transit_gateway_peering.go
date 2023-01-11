package networkmanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/networkmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceTransitGatewayPeering() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTransitGatewayPeeringCreate,
		ReadWithoutTimeout:   resourceTransitGatewayPeeringRead,
		UpdateWithoutTimeout: resourceTransitGatewayPeeringUpdate,
		DeleteWithoutTimeout: resourceTransitGatewayPeeringDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_network_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"edge_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peering_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"transit_gateway_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"transit_gateway_peering_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTransitGatewayPeeringCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	coreNetworkID := d.Get("core_network_id").(string)
	transitGatewayARN := d.Get("transit_gateway_arn").(string)
	input := &networkmanager.CreateTransitGatewayPeeringInput{
		CoreNetworkId:     aws.String(coreNetworkID),
		TransitGatewayArn: aws.String(transitGatewayARN),
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	log.Printf("[DEBUG] Creating Network Manager Transit Gateway Peering: %s", input)
	output, err := conn.CreateTransitGatewayPeeringWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating Network Manager Transit Gateway (%s) Peering (%s): %s", transitGatewayARN, coreNetworkID, err)
	}

	d.SetId(aws.StringValue(output.TransitGatewayPeering.Peering.PeeringId))

	if _, err := waitTransitGatewayPeeringCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("waiting for Network Manager Transit Gateway Peering (%s) create: %s", d.Id(), err)
	}

	return resourceTransitGatewayPeeringRead(ctx, d, meta)
}

func resourceTransitGatewayPeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	transitGatewayPeering, err := FindTransitGatewayPeeringByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Network Manager Transit Gateway Peering %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Network Manager Transit Gateway Peering (%s): %s", d.Id(), err)
	}

	p := transitGatewayPeering.Peering
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "networkmanager",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("peering/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("core_network_arn", p.CoreNetworkArn)
	d.Set("core_network_id", p.CoreNetworkId)
	d.Set("edge_location", p.EdgeLocation)
	d.Set("owner_account_id", p.OwnerAccountId)
	d.Set("peering_type", p.PeeringType)
	d.Set("resource_arn", p.ResourceArn)
	d.Set("transit_gateway_arn", transitGatewayPeering.TransitGatewayArn)
	d.Set("transit_gateway_peering_attachment_id", transitGatewayPeering.TransitGatewayPeeringAttachmentId)

	tags := KeyValueTags(p.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return diag.Errorf("Setting tags: %s", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return diag.Errorf("setting tags_all: %s", err)
	}

	return nil
}

func resourceTransitGatewayPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTagsWithContext(ctx, conn, d.Get("arn").(string), o, n); err != nil {
			return diag.Errorf("updating Network Manager Transit Gateway Peering (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceTransitGatewayPeeringRead(ctx, d, meta)
}

func resourceTransitGatewayPeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).NetworkManagerConn

	log.Printf("[DEBUG] Deleting Network Manager Transit Gateway Peering: %s", d.Id())
	_, err := conn.DeletePeeringWithContext(ctx, &networkmanager.DeletePeeringInput{
		PeeringId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting Network Manager Transit Gateway Peering (%s): %s", d.Id(), err)
	}

	if _, err := waitTransitGatewayPeeringDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("waiting for Network Manager Transit Gateway Peering (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func FindTransitGatewayPeeringByID(ctx context.Context, conn *networkmanager.NetworkManager, id string) (*networkmanager.TransitGatewayPeering, error) {
	input := &networkmanager.GetTransitGatewayPeeringInput{
		PeeringId: aws.String(id),
	}

	output, err := conn.GetTransitGatewayPeeringWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, networkmanager.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.TransitGatewayPeering == nil || output.TransitGatewayPeering.Peering == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.TransitGatewayPeering, nil
}

func StatusTransitGatewayPeeringState(ctx context.Context, conn *networkmanager.NetworkManager, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindTransitGatewayPeeringByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Peering.State), nil
	}
}

func waitTransitGatewayPeeringCreated(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.TransitGatewayPeering, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{networkmanager.PeeringStateCreating},
		Target:  []string{networkmanager.PeeringStateAvailable},
		Timeout: timeout,
		Refresh: StatusTransitGatewayPeeringState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.TransitGatewayPeering); ok {
		return output, err
	}

	return nil, err
}

func waitTransitGatewayPeeringDeleted(ctx context.Context, conn *networkmanager.NetworkManager, id string, timeout time.Duration) (*networkmanager.TransitGatewayPeering, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{networkmanager.PeeringStateDeleting},
		Target:  []string{},
		Timeout: timeout,
		Refresh: StatusTransitGatewayPeeringState(ctx, conn, id),
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*networkmanager.TransitGatewayPeering); ok {
		return output, err
	}

	return nil, err
}
