package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVSpherePrivateTagAttach() *schema.Resource {
	return &schema.Resource{
		Create:        resourceVSpherePrivateTagAttachCreate,
		Read:          resourceVSpherePrivateTagAttachRead,
		Delete:        resourceVSpherePrivateTagAttachDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"objectid": {
				Type:         schema.TypeString,
				Description:  "The managed object id of the object to attach the tag to.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"objecttype": {
				Type:         schema.TypeString,
				Description:  "The type of the managed object.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"tagid": {
				Type:         schema.TypeString,
				Description:  "The id of tag.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceVSpherePrivateTagAttachCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[DEBUG] beginning of tag attachment")

	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	tagManager := tags.NewManager(meta.(*VSphereClient).restClient)
	moRef := types.ManagedObjectReference{
		Value: d.Get("objectid").(string),
		Type:  d.Get("objecttype").(string),
	}

	err := tagManager.AttachTag(ctx, d.Get("tagid").(string), moRef)

	if err != nil {
		return fmt.Errorf("unable to attach tag %s to object %s caused by %s", d.Get("tagid").(string), d.Get("objectid"), err)
	}

	// Creating a tag attachment does not result in a new id. Generating an id
	// from the existing tag id, existing object type and object id.
	id := fmt.Sprintf("%s_%s_%s", d.Get("tagid"), d.Get("objecttype"), d.Get("objectid"))

	d.SetId(id)

	log.Print("[DEBUG] Tag attachment complete.")
	return resourceVSpherePrivateTagAttachRead(d, meta)
}

func resourceVSpherePrivateTagAttachRead(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	client := meta.(*VSphereClient).restClient
	tagManager := tags.NewManager(client)
	id := strings.Split(d.Id(), "_")

	if len(id) != 3 {
		return fmt.Errorf("error tag attachment id incorrect length, id: %s", d.Id())
	}

	attachedObjects, err := tagManager.GetAttachedObjectsOnTags(ctx, []string{id[0]})
	if err != nil {
		return fmt.Errorf("error attached objects not found: %s", err)
	}

	found := false
	for _, attachedObjs := range attachedObjects {
		for _, objid := range attachedObjs.ObjectIDs {
			if objid.Reference().Value == id[2] && objid.Reference().Type == id[1] {
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("error attached objects not found: %s", d.Id())
	}

	return nil
}

func resourceVSpherePrivateTagAttachDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning tag detachment")
	id := strings.Split(d.Id(), "_")

	if len(id) != 3 {
		return fmt.Errorf("error tag attachment id incorrect length, id: %s", d.Id())
	}

	ctx, cancel := context.WithTimeout(context.TODO(), defaultAPITimeout)
	defer cancel()
	tagManager := tags.NewManager(meta.(*VSphereClient).restClient)
	moRef := types.ManagedObjectReference{
		Value: id[2],
		Type:  id[1],
	}

	err := tagManager.DetachTag(ctx, id[0], moRef)

	if err != nil {
		return fmt.Errorf("error detaching tag: %s", err)
	}

	d.SetId("")

	log.Print("[DEBUG] Tag Detachment complete.")

	return nil
}
