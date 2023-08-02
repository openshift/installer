package location

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/locationservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

// @SDKResource("aws_location_geofence_collection", name="Geofence Collection")
// @Tags(identifierAttribute="collection_arn")
func ResourceGeofenceCollection() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGeofenceCollectionCreate,
		ReadWithoutTimeout:   resourceGeofenceCollectionRead,
		UpdateWithoutTimeout: resourceGeofenceCollectionUpdate,
		DeleteWithoutTimeout: resourceGeofenceCollectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"collection_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"collection_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 2048),
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameGeofenceCollection = "Geofence Collection"
)

func resourceGeofenceCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LocationConn(ctx)

	in := &locationservice.CreateGeofenceCollectionInput{
		CollectionName: aws.String(d.Get("collection_name").(string)),
		Tags:           GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("description"); ok && v != "" {
		in.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok && v != "" {
		in.KmsKeyId = aws.String(v.(string))
	}

	out, err := conn.CreateGeofenceCollectionWithContext(ctx, in)
	if err != nil {
		return create.DiagError(names.Location, create.ErrActionCreating, ResNameGeofenceCollection, d.Get("collection_name").(string), err)
	}

	if out == nil {
		return create.DiagError(names.Location, create.ErrActionCreating, ResNameGeofenceCollection, d.Get("collection_name").(string), errors.New("empty output"))
	}

	d.SetId(aws.StringValue(out.CollectionName))

	return resourceGeofenceCollectionRead(ctx, d, meta)
}

func resourceGeofenceCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LocationConn(ctx)

	out, err := findGeofenceCollectionByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Location GeofenceCollection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.Location, create.ErrActionReading, ResNameGeofenceCollection, d.Id(), err)
	}

	d.Set("collection_arn", out.CollectionArn)
	d.Set("collection_name", out.CollectionName)
	d.Set("create_time", aws.TimeValue(out.CreateTime).Format(time.RFC3339))
	d.Set("description", out.Description)
	d.Set("kms_key_id", out.KmsKeyId)
	d.Set("update_time", aws.TimeValue(out.UpdateTime).Format(time.RFC3339))

	return nil
}

func resourceGeofenceCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LocationConn(ctx)

	update := false

	in := &locationservice.UpdateGeofenceCollectionInput{
		CollectionName: aws.String(d.Id()),
	}

	if d.HasChange("description") {
		in.Description = aws.String(d.Get("description").(string))
		update = true
	}

	if !update {
		return nil
	}

	log.Printf("[DEBUG] Updating Location GeofenceCollection (%s): %#v", d.Id(), in)
	_, err := conn.UpdateGeofenceCollectionWithContext(ctx, in)
	if err != nil {
		return create.DiagError(names.Location, create.ErrActionUpdating, ResNameGeofenceCollection, d.Id(), err)
	}

	return resourceGeofenceCollectionRead(ctx, d, meta)
}

func resourceGeofenceCollectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LocationConn(ctx)

	log.Printf("[INFO] Deleting Location GeofenceCollection %s", d.Id())

	_, err := conn.DeleteGeofenceCollectionWithContext(ctx, &locationservice.DeleteGeofenceCollectionInput{
		CollectionName: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, locationservice.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return create.DiagError(names.Location, create.ErrActionDeleting, ResNameGeofenceCollection, d.Id(), err)
	}

	return nil
}

func findGeofenceCollectionByName(ctx context.Context, conn *locationservice.LocationService, name string) (*locationservice.DescribeGeofenceCollectionOutput, error) {
	in := &locationservice.DescribeGeofenceCollectionInput{
		CollectionName: aws.String(name),
	}

	out, err := conn.DescribeGeofenceCollectionWithContext(ctx, in)
	if tfawserr.ErrCodeEquals(err, locationservice.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}

	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
