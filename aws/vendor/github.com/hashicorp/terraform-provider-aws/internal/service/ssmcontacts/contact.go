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
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ssmcontacts_contact", name="Context")
// @Tags(identifierAttribute="id")
func ResourceContact() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceContactCreate,
		ReadWithoutTimeout:   resourceContactRead,
		UpdateWithoutTimeout: resourceContactUpdate,
		DeleteWithoutTimeout: resourceContactDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameContact = "Contact"
)

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	input := &ssmcontacts.CreateContactInput{
		Alias:       aws.String(d.Get("alias").(string)),
		DisplayName: aws.String(d.Get("display_name").(string)),
		Plan:        &types.Plan{Stages: []types.Stage{}},
		Tags:        GetTagsIn(ctx),
		Type:        types.ContactType(d.Get("type").(string)),
	}

	output, err := client.CreateContact(ctx, input)
	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionCreating, ResNameContact, d.Get("alias").(string), err)
	}

	if output == nil {
		return create.DiagError(names.SSMContacts, create.ErrActionCreating, ResNameContact, d.Get("alias").(string), errors.New("empty output"))
	}

	d.SetId(aws.ToString(output.ContactArn))

	return resourceContactRead(ctx, d, meta)
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	out, err := findContactByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] SSMContacts Contact (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionReading, ResNameContact, d.Id(), err)
	}

	if err := setContactResourceData(d, out); err != nil {
		return create.DiagError(names.SSMContacts, create.ErrActionSetting, ResNameContact, d.Id(), err)
	}

	return nil
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	if d.HasChanges("display_name") {
		in := &ssmcontacts.UpdateContactInput{
			ContactId:   aws.String(d.Id()),
			DisplayName: aws.String(d.Get("display_name").(string)),
		}

		_, err := conn.UpdateContact(ctx, in)
		if err != nil {
			return create.DiagError(names.SSMContacts, create.ErrActionUpdating, ResNameContact, d.Id(), err)
		}
	}

	return resourceContactRead(ctx, d, meta)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SSMContactsClient(ctx)

	log.Printf("[INFO] Deleting SSMContacts Contact %s", d.Id())

	_, err := conn.DeleteContact(ctx, &ssmcontacts.DeleteContactInput{
		ContactId: aws.String(d.Id()),
	})

	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.SSMContacts, create.ErrActionDeleting, ResNameContact, d.Id(), err)
	}
	return nil
}
