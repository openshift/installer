package vsphere

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVSphereFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereFolderCreate,
		Read:   resourceVSphereFolderRead,
		Update: resourceVSphereFolderUpdate,
		Delete: resourceVSphereFolderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereFolderImport,
		},

		SchemaVersion: 1,
		MigrateState:  resourceVSphereFolderMigrateState,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Description:  "The path of the folder and any parents, relative to the datacenter and folder type being defined.",
				Required:     true,
				StateFunc:    folder.NormalizePath,
				ValidateFunc: validation.NoZeroValues,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of the folder.",
				ForceNew:    true,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(folder.VSphereFolderTypeVM),
						string(folder.VSphereFolderTypeNetwork),
						string(folder.VSphereFolderTypeHost),
						string(folder.VSphereFolderTypeDatastore),
						string(folder.VSphereFolderTypeDatacenter),
					},
					false,
				),
			},
			"datacenter_id": {
				Type:        schema.TypeString,
				Description: "The ID of the datacenter. Can be ignored if creating a datacenter folder, otherwise required.",
				ForceNew:    true,
				Optional:    true,
			},
			// Tagging
			vSphereTagAttributeKey: tagsSchema(),
			// Custom Attributes
			customattribute.ConfigKey: customattribute.ConfigSchema(),
		},
	}
}

func resourceVSphereFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	ft := folder.VSphereFolderType(d.Get("type").(string))
	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		var err error
		dc, err = datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
	} else if ft != folder.VSphereFolderTypeDatacenter {
		return fmt.Errorf("datacenter_id cannot be empty when creating a targetFolder of type %s", ft)
	}

	p := d.Get("path").(string)

	// Determine the parent targetFolder
	parent, err := folder.ParentFromPath(client, p, ft, dc)
	if err != nil {
		return fmt.Errorf("error trying to determine parent targetFolder: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()

	targetFolder, err := parent.CreateFolder(ctx, path.Base(p))
	if err != nil {
		return fmt.Errorf("error creating targetFolder: %s", err)
	}

	d.SetId(targetFolder.Reference().Value)

	// Apply any pending tags now
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, targetFolder); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	// Set custom attributes
	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(targetFolder); err != nil {
			return fmt.Errorf("error setting custom attributes: %s", err)
		}
	}

	return resourceVSphereFolderRead(d, meta)
}

func resourceVSphereFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	fo, err := folder.FromID(client, d.Id())
	if err != nil {
		return fmt.Errorf("cannot locate folder: %s", err)
	}

	// Determine the folder type first. We use the folder as the source of truth
	// here versus the state so that we can support import.
	ft, err := folder.FindType(fo)
	if err != nil {
		return fmt.Errorf("cannot determine folder type: %s", err)
	}

	// Again, to support a clean import (which is done off of absolute path to
	// the folder), we discover the datacenter from the path (if it's a thing).
	var dc *object.Datacenter
	p := fo.InventoryPath
	if ft != folder.VSphereFolderTypeDatacenter {
		particle := folder.RootPathParticle(ft)
		dcp, err := particle.SplitDatacenter(p)
		if err != nil {
			return fmt.Errorf("cannot determine datacenter path: %s", err)
		}
		dc, err = getDatacenter(client, dcp)
		if err != nil {
			return fmt.Errorf("cannot find datacenter from path %q: %s", dcp, err)
		}
		relative, err := particle.SplitRelative(p)
		if err != nil {
			return fmt.Errorf("cannot determine relative folder path: %s", err)
		}
		p = relative
	}

	_ = d.Set("path", folder.NormalizePath(p))
	_ = d.Set("type", ft)
	if dc != nil {
		_ = d.Set("datacenter_id", dc.Reference().Value)
	}

	// Read tags if we have the ability to do so
	if tagsClient, _ := meta.(*Client).TagsManager(); tagsClient != nil {
		if err := readTagsForResource(tagsClient, fo, d); err != nil {
			return fmt.Errorf("error reading tags: %s", err)
		}
	}

	// Read custom attributes
	if customattribute.IsSupported(client) {
		moFolder, err := folder.Properties(fo)
		if err != nil {
			return err
		}
		customattribute.ReadFromResource(moFolder.Entity(), d)
	}

	return nil
}

func resourceVSphereFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	tagsClient, err := tagsManagerIfDefined(d, meta)
	if err != nil {
		return err
	}
	// Verify a proper vCenter before proceeding if custom attributes are defined
	attrsProcessor, err := customattribute.GetDiffProcessorIfAttributesDefined(client, d)
	if err != nil {
		return err
	}

	fo, err := folder.FromID(client, d.Id())
	if err != nil {
		return fmt.Errorf("cannot locate folder: %s", err)
	}

	// Apply any pending tags first as it's the lesser expensive of the two
	// operations
	if tagsClient != nil {
		if err := processTagDiff(tagsClient, d, fo); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	if attrsProcessor != nil {
		if err := attrsProcessor.ProcessDiff(fo); err != nil {
			return fmt.Errorf("error setting custom attributes: %s", err)
		}
	}

	var dc *object.Datacenter
	if dcID, ok := d.GetOk("datacenter_id"); ok {
		var err error
		dc, err = datacenterFromID(client, dcID.(string))
		if err != nil {
			return fmt.Errorf("cannot locate datacenter: %s", err)
		}
	}

	if d.HasChange("path") {
		// The path has changed, which could mean either a change in parent, a
		// change in name, or both.
		ft := folder.VSphereFolderType(d.Get("type").(string))
		oldp, newp := d.GetChange("path")
		oldpa, err := folder.ParentFromPath(client, oldp.(string), ft, dc)
		if err != nil {
			return fmt.Errorf("error parsing parent folder from path %q: %s", oldp.(string), err)
		}
		newpa, err := folder.ParentFromPath(client, newp.(string), ft, dc)
		if err != nil {
			return fmt.Errorf("error parsing parent folder from path %q: %s", newp.(string), err)
		}
		oldn := path.Base(oldp.(string))
		newn := path.Base(newp.(string))

		if oldn != newn {
			// Folder base name has changed and needs a rename
			if err := viapi.RenameObject(client, fo.Reference(), newn); err != nil {
				return fmt.Errorf("could not rename folder: %s", err)
			}
		}
		if oldpa.Reference().Value != newpa.Reference().Value {
			// The parent folder has changed - we need to move the folder into the
			// new path
			ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
			defer cancel()
			task, err := newpa.MoveInto(ctx, []types.ManagedObjectReference{fo.Reference()})
			if err != nil {
				return fmt.Errorf("could not move folder: %s", err)
			}
			tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
			defer tcancel()
			if err := task.Wait(tctx); err != nil {
				return fmt.Errorf("error on waiting for move task completion: %s", err)
			}
		}
	}

	return resourceVSphereFolderRead(d, meta)
}

func resourceVSphereFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	fo, err := folder.FromID(client, d.Id())
	if err != nil {
		return fmt.Errorf("cannot locate folder: %s", err)
	}

	// We don't destroy if the folder has children. This might be flaggable in
	// the future, but I don't think it's necessary at this point in time -
	// better to have hardcoded safe behavior than hardcoded unsafe behavior.
	ne, err := folder.HasChildren(fo)
	if err != nil {
		return fmt.Errorf("error checking for folder contents: %s", err)
	}
	if ne {
		return errors.New("folder is not empty, please remove all items before deleting")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	task, err := fo.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete folder: %s", err)
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return fmt.Errorf("error on waiting for deletion task completion: %s", err)
	}

	return nil
}

func resourceVSphereFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Our subject is the full path to a specific targetFolder, for which we just get
	// the MOID for and then pass off to Read. Easy peasy.
	p := d.Id()
	if !strings.HasPrefix(p, "/") {
		return nil, errors.New("path must start with a trailing slash")
	}
	client := meta.(*Client).vimClient
	p = folder.NormalizePath(p)
	targetFolder, err := folder.FromAbsolutePath(client, p)
	if err != nil {
		return nil, err
	}
	d.SetId(targetFolder.Reference().Value)
	return []*schema.ResourceData{d}, nil
}
