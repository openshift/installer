package ce

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ce_cost_allocation_tag")
func ResourceCostAllocationTag() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCostAllocationTagUpdate,
		ReadWithoutTimeout:   resourceCostAllocationTagRead,
		UpdateWithoutTimeout: resourceCostAllocationTagUpdate,
		DeleteWithoutTimeout: resourceCostAllocationTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(costexplorer.CostAllocationTagStatus_Values(), false),
			},
			"tag_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCostAllocationTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)

	costAllocTag, err := FindCostAllocationTagByKey(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		create.LogNotFoundRemoveState(names.CE, create.ErrActionReading, ResNameCostAllocationTag, d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.CE, create.ErrActionReading, ResNameCostAllocationTag, d.Id(), err)
	}

	d.Set("tag_key", costAllocTag.TagKey)
	d.Set("status", costAllocTag.Status)
	d.Set("type", costAllocTag.Type)

	return nil
}

func resourceCostAllocationTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	key := d.Get("tag_key").(string)

	updateTagStatus(ctx, d, meta, false)

	d.SetId(key)

	return resourceCostAllocationTagRead(ctx, d, meta)
}

func resourceCostAllocationTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateTagStatus(ctx, d, meta, true)
}

func updateTagStatus(ctx context.Context, d *schema.ResourceData, meta interface{}, delete bool) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).CEConn(ctx)

	key := d.Get("tag_key").(string)
	tagStatus := &costexplorer.CostAllocationTagStatusEntry{
		TagKey: aws.String(key),
		Status: aws.String(d.Get("status").(string)),
	}

	if delete {
		tagStatus.Status = aws.String(costexplorer.CostAllocationTagStatusInactive)
	}

	input := &costexplorer.UpdateCostAllocationTagsStatusInput{
		CostAllocationTagsStatus: []*costexplorer.CostAllocationTagStatusEntry{tagStatus},
	}

	_, err := conn.UpdateCostAllocationTagsStatusWithContext(ctx, input)

	if err != nil {
		return create.DiagError(names.CE, create.ErrActionUpdating, ResNameCostAllocationTag, d.Id(), err)
	}

	return nil
}
