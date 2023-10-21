package ssmcontacts

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssmcontacts"
	"github.com/aws/aws-sdk-go-v2/service/ssmcontacts/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ssmcontacts_contact_channel", name="Contact Channel")
func ResourceContactChannel() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceContactChannelCreate,
		ReadWithoutTimeout:   resourceContactChannelRead,
		UpdateWithoutTimeout: resourceContactChannelUpdate,
		DeleteWithoutTimeout: resourceContactChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"activation_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact_id": {
				ForceNew: true,
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_address": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"simple_address": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				ForceNew: true,
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

const (
	ResNameContactChannel = "Contact Channel"
)

func resourceContactChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	delivery_address := expandContactChannelAddress(d.Get("delivery_address").([]interface{}))
	in := &ssmcontacts.CreateContactChannelInput{
		ContactId:       aws.String(d.Get("contact_id").(string)),
		DeferActivation: aws.Bool(true),
		DeliveryAddress: delivery_address,
		Name:            aws.String(d.Get("name").(string)),
		Type:            types.ChannelType(d.Get("type").(string)),
	}

	out, err := conn.CreateContactChannel(ctx, in)
	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionCreating, ResNameContactChannel, d.Get("name").(string), err)
	}

	if out == nil {
		return create.DiagError(names.SSMContacts, create.ErrActionCreating, ResNameContactChannel, d.Get("name").(string), errors.New("empty output"))
	}

	d.SetId(aws.ToString(out.ContactChannelArn))

	return resourceContactChannelRead(ctx, d, meta)
}

func resourceContactChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	out, err := findContactChannelByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] SSMContacts ContactChannel (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionReading, ResNameContactChannel, d.Id(), err)
	}

	if err := setContactChannelResourceData(d, out); err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionSetting, ResNameContactChannel, d.Id(), err)
	}

	return nil
}

func resourceContactChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	update := false

	in := &ssmcontacts.UpdateContactChannelInput{
		ContactChannelId: aws.String(d.Id()),
	}

	if d.HasChanges("delivery_address") {
		in.DeliveryAddress = expandContactChannelAddress(d.Get("delivery_address").([]interface{}))
		update = true
	}

	if d.HasChanges("name") {
		in.Name = aws.String(d.Get("name").(string))
		update = true
	}

	if !update {
		return nil
	}

	log.Printf("[DEBUG] Updating SSMContacts ContactChannel (%s): %#v", d.Id(), in)
	_, err := conn.UpdateContactChannel(ctx, in)
	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionUpdating, ResNameContactChannel, d.Id(), err)
	}

	return resourceContactChannelRead(ctx, d, meta)
}

func resourceContactChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	log.Printf("[INFO] Deleting SSMContacts ContactChannel %s", d.Id())

	_, err := conn.DeleteContactChannel(ctx, &ssmcontacts.DeleteContactChannelInput{
		ContactChannelId: aws.String(d.Id()),
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.SSMContacts, create.ErrActionDeleting, ResNameContactChannel, d.Id(), err)
	}

	return nil
}
