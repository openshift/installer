package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ec2_managed_prefix_list", name="Managed Prefix List")
// @Tags(identifierAttribute="id")
func ResourceManagedPrefixList() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceManagedPrefixListCreate,
		ReadWithoutTimeout:   resourceManagedPrefixListRead,
		UpdateWithoutTimeout: resourceManagedPrefixListUpdate,
		DeleteWithoutTimeout: resourceManagedPrefixListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.Sequence(
			customdiff.ComputedIf("version", func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("entry")
			}),
			verify.SetTagsDiff,
		),

		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(managedPrefixListAddressFamily_Values(), false),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entry": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
					},
				},
			},
			"max_entries": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceManagedPrefixListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.CreateManagedPrefixListInput{
		TagSpecifications: getTagSpecificationsIn(ctx, ec2.ResourceTypePrefixList),
	}

	if v, ok := d.GetOk("address_family"); ok {
		input.AddressFamily = aws.String(v.(string))
	}

	if v, ok := d.GetOk("entry"); ok && v.(*schema.Set).Len() > 0 {
		input.Entries = expandAddPrefixListEntries(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("max_entries"); ok {
		input.MaxEntries = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("name"); ok {
		input.PrefixListName = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating EC2 Managed Prefix List: %s", input)
	output, err := conn.CreateManagedPrefixListWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("creating EC2 Managed Prefix List: %s", err)
	}

	d.SetId(aws.StringValue(output.PrefixList.PrefixListId))

	if _, err := WaitManagedPrefixListCreated(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for EC2 Managed Prefix List (%s) create: %s", d.Id(), err)
	}

	return resourceManagedPrefixListRead(ctx, d, meta)
}

func resourceManagedPrefixListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	pl, err := FindManagedPrefixListByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Managed Prefix List %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading EC2 Managed Prefix List (%s): %s", d.Id(), err)
	}

	prefixListEntries, err := FindManagedPrefixListEntriesByID(ctx, conn, d.Id())

	if err != nil {
		return diag.Errorf("reading EC2 Managed Prefix List (%s) Entries: %s", d.Id(), err)
	}

	d.Set("address_family", pl.AddressFamily)
	d.Set("arn", pl.PrefixListArn)
	if err := d.Set("entry", flattenPrefixListEntries(prefixListEntries)); err != nil {
		return diag.Errorf("setting entry: %s", err)
	}
	d.Set("max_entries", pl.MaxEntries)
	d.Set("name", pl.PrefixListName)
	d.Set("owner_id", pl.OwnerId)
	d.Set("version", pl.Version)

	SetTagsOut(ctx, pl.Tags)

	return nil
}

func resourceManagedPrefixListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	// MaxEntries & Entry cannot change in the same API call.
	//   If MaxEntry is increasing, complete before updating entry(s)
	//   If MaxEntry is decreasing, complete after updating entry(s)
	maxEntryChangedDecrease := false
	var newMaxEntryInt int64

	if d.HasChange("max_entries") {
		oldMaxEntry, newMaxEntry := d.GetChange("max_entries")
		newMaxEntryInt = int64(d.Get("max_entries").(int))

		if newMaxEntry.(int) < oldMaxEntry.(int) {
			maxEntryChangedDecrease = true
		} else {
			err := updateMaxEntry(ctx, conn, d.Id(), newMaxEntryInt)
			if err != nil {
				return diag.Errorf("updating EC2 Managed Prefix List (%s) increased MaxEntries : %s", d.Id(), err)
			}
		}
	}

	if d.HasChangesExcept("tags", "tags_all", "max_entries") {
		input := &ec2.ModifyManagedPrefixListInput{
			PrefixListId: aws.String(d.Id()),
		}

		input.PrefixListName = aws.String(d.Get("name").(string))
		currentVersion := int64(d.Get("version").(int))
		wait := false

		oldAttr, newAttr := d.GetChange("entry")
		os := oldAttr.(*schema.Set)
		ns := newAttr.(*schema.Set)

		if addEntries := ns.Difference(os); addEntries.Len() > 0 {
			input.AddEntries = expandAddPrefixListEntries(addEntries.List())
			input.CurrentVersion = aws.Int64(currentVersion)
			wait = true
		}

		if removeEntries := os.Difference(ns); removeEntries.Len() > 0 {
			input.RemoveEntries = expandRemovePrefixListEntries(removeEntries.List())
			input.CurrentVersion = aws.Int64(currentVersion)
			wait = true
		}

		// Prevent the following error on description-only updates:
		//   InvalidParameterValue: Request cannot contain Cidr #.#.#.#/# in both AddPrefixListEntries and RemovePrefixListEntries
		// Attempting to just delete the RemoveEntries item causes:
		//   InvalidRequest: The request received was invalid.
		// Therefore it seems we must issue two ModifyManagedPrefixList calls,
		// one with a collection of all description-only removals and the
		// second one will add them all back.
		if len(input.AddEntries) > 0 && len(input.RemoveEntries) > 0 {
			descriptionOnlyRemovals := []*ec2.RemovePrefixListEntry{}
			removals := []*ec2.RemovePrefixListEntry{}

			for _, removeEntry := range input.RemoveEntries {
				inAddAndRemove := false

				for _, addEntry := range input.AddEntries {
					if aws.StringValue(addEntry.Cidr) == aws.StringValue(removeEntry.Cidr) {
						inAddAndRemove = true
						break
					}
				}

				if inAddAndRemove {
					descriptionOnlyRemovals = append(descriptionOnlyRemovals, removeEntry)
				} else {
					removals = append(removals, removeEntry)
				}
			}

			if len(descriptionOnlyRemovals) > 0 {
				_, err := conn.ModifyManagedPrefixListWithContext(ctx, &ec2.ModifyManagedPrefixListInput{
					CurrentVersion: input.CurrentVersion,
					PrefixListId:   aws.String(d.Id()),
					RemoveEntries:  descriptionOnlyRemovals,
				})

				if err != nil {
					return diag.Errorf("updating EC2 Managed Prefix List (%s): %s", d.Id(), err)
				}

				managedPrefixList, err := WaitManagedPrefixListModified(ctx, conn, d.Id())

				if err != nil {
					return diag.Errorf("waiting for EC2 Managed Prefix List (%s) update: %s", d.Id(), err)
				}

				input.CurrentVersion = managedPrefixList.Version
			}

			if len(removals) > 0 {
				input.RemoveEntries = removals
			} else {
				// Prevent this error if RemoveEntries is list with no elements after removals:
				//   InvalidRequest: The request received was invalid.
				input.RemoveEntries = nil
			}
		}

		_, err := conn.ModifyManagedPrefixListWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating EC2 Managed Prefix List (%s): %s", d.Id(), err)
		}

		if wait {
			if _, err := WaitManagedPrefixListModified(ctx, conn, d.Id()); err != nil {
				return diag.Errorf("waiting for EC2 Managed Prefix List (%s) update: %s", d.Id(), err)
			}
		}
	}

	// Only decrease MaxEntries after entry(s) have had opportunity to be removed
	if maxEntryChangedDecrease {
		err := updateMaxEntry(ctx, conn, d.Id(), newMaxEntryInt)
		if err != nil {
			return diag.Errorf("updating EC2 Managed Prefix List (%s) decreased MaxEntries : %s", d.Id(), err)
		}
	}

	return resourceManagedPrefixListRead(ctx, d, meta)
}

func resourceManagedPrefixListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[INFO] Deleting EC2 Managed Prefix List: %s", d.Id())
	_, err := conn.DeleteManagedPrefixListWithContext(ctx, &ec2.DeleteManagedPrefixListInput{
		PrefixListId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidPrefixListIDNotFound) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting EC2 Managed Prefix List (%s): %s", d.Id(), err)
	}

	if _, err := WaitManagedPrefixListDeleted(ctx, conn, d.Id()); err != nil {
		return diag.Errorf("waiting for EC2 Managed Prefix List (%s) delete: %s", d.Id(), err)
	}

	return nil
}

func updateMaxEntry(ctx context.Context, conn *ec2.EC2, id string, maxEntries int64) error {
	_, err := conn.ModifyManagedPrefixListWithContext(ctx, &ec2.ModifyManagedPrefixListInput{
		PrefixListId: aws.String(id),
		MaxEntries:   aws.Int64(maxEntries),
	})

	if err != nil {
		return fmt.Errorf("updating MaxEntries for EC2 Managed Prefix List (%s): %s", id, err)
	}

	_, err = WaitManagedPrefixListModified(ctx, conn, id)

	if err != nil {
		return fmt.Errorf("waiting for EC2 Managed Prefix List (%s) MaxEntries update: %s", id, err)
	}

	return nil
}

func expandAddPrefixListEntry(tfMap map[string]interface{}) *ec2.AddPrefixListEntry {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.AddPrefixListEntry{}

	if v, ok := tfMap["cidr"].(string); ok && v != "" {
		apiObject.Cidr = aws.String(v)
	}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		apiObject.Description = aws.String(v)
	}

	return apiObject
}

func expandAddPrefixListEntries(tfList []interface{}) []*ec2.AddPrefixListEntry {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*ec2.AddPrefixListEntry

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandAddPrefixListEntry(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func expandRemovePrefixListEntry(tfMap map[string]interface{}) *ec2.RemovePrefixListEntry {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.RemovePrefixListEntry{}

	if v, ok := tfMap["cidr"].(string); ok && v != "" {
		apiObject.Cidr = aws.String(v)
	}

	return apiObject
}

func expandRemovePrefixListEntries(tfList []interface{}) []*ec2.RemovePrefixListEntry {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*ec2.RemovePrefixListEntry

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandRemovePrefixListEntry(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenPrefixListEntry(apiObject *ec2.PrefixListEntry) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Cidr; v != nil {
		tfMap["cidr"] = aws.StringValue(v)
	}

	if v := apiObject.Description; v != nil {
		tfMap["description"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenPrefixListEntries(apiObjects []*ec2.PrefixListEntry) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenPrefixListEntry(apiObject))
	}

	return tfList
}
