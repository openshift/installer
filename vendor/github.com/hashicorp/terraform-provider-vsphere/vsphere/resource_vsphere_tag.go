package vsphere

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/govmomi/vapi/tags"
)

func resourceVSphereTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereTagCreate,
		Read:   resourceVSphereTagRead,
		Update: resourceVSphereTagUpdate,
		Delete: resourceVSphereTagDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereTagImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The display name of the tag. The name must be unique within its category.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tag.",
				Optional:    true,
			},
			"category_id": {
				Type:        schema.TypeString,
				Description: "The unique identifier of the parent category in which this tag will be created.",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceVSphereTagCreate(d *schema.ResourceData, meta interface{}) error {
	tm, err := meta.(*VSphereClient).TagsManager()
	if err != nil {
		return err
	}
	spec := &tags.Tag{
		CategoryID:  d.Get("category_id").(string),
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	id, err := tm.CreateTag(ctx, spec)
	if err != nil {
		return fmt.Errorf("could not create tag: %s", err)
	}
	if id == "" {
		return errors.New("no ID was returned")
	}
	d.SetId(id)
	return resourceVSphereTagRead(d, meta)
}

func resourceVSphereTagRead(d *schema.ResourceData, meta interface{}) error {
	tm, err := meta.(*VSphereClient).TagsManager()
	if err != nil {
		return err
	}

	id := d.Id()

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	tag, err := tm.GetTag(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "com.vmware.vapi.std.errors.not_found") {
			log.Printf("[DEBUG] Tag %s: Resource has been deleted", id)
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("name", tag.Name)
	d.Set("description", tag.Description)
	d.Set("category_id", tag.CategoryID)

	return nil
}

func resourceVSphereTagUpdate(d *schema.ResourceData, meta interface{}) error {
	tm, err := meta.(*VSphereClient).TagsManager()
	if err != nil {
		return err
	}

	id := d.Id()
	spec := &tags.Tag{
		ID:          id,
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	err = tm.UpdateTag(ctx, spec)
	if err != nil {
		return fmt.Errorf("could not update tag with id %q: %s", id, err)
	}
	return resourceVSphereTagRead(d, meta)
}

func resourceVSphereTagDelete(d *schema.ResourceData, meta interface{}) error {
	tm, err := meta.(*VSphereClient).TagsManager()
	if err != nil {
		return err
	}

	id := d.Id()

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	tag, err := tm.GetTag(ctx, id)
	if err != nil {
		return err
	}
	err = tm.DeleteTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("could not delete tag with id %q: %s", id, err)
	}
	return nil
}

func resourceVSphereTagImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import takes the tag and category names through JSON to make sure we can
	// search on special characters, since there does not seem to be any sort of
	// prohibited kind of character when dealing with either tags or categories.
	//
	// We just decode to a map[string]string and handle the rest from there. We
	// don't care about any other kind of value, so we lean on JSON errors in
	// those cases.
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	categoryName, ok := data["category_name"]
	if !ok {
		return nil, errors.New("missing category_name in input data")
	}
	tagName, ok := data["tag_name"]
	if !ok {
		return nil, errors.New("missing tag_name in input data")
	}

	tm, err := meta.(*VSphereClient).TagsManager()
	if err != nil {
		return nil, err
	}

	categoryID, err := tagCategoryByName(tm, categoryName)
	if err != nil {
		return nil, err
	}
	tagID, err := tagByName(tm, tagName, categoryID)
	if err != nil {
		return nil, err
	}

	d.SetId(tagID)
	return []*schema.ResourceData{d}, nil
}
