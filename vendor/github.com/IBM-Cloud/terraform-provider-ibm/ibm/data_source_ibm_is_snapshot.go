// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSnapshotRead,

		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSnapshotName, "identifier"},
				Description:  "Snapshot identifier",
				ValidateFunc: InvokeDataSourceValidator("ibm_is_snapshot", "identifier"),
			},

			isSnapshotName: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isSnapshotName, "identifier"},
				ValidateFunc: InvokeDataSourceValidator("ibm_is_snapshot", isSnapshotName),
				Description:  "Snapshot name",
			},

			isSnapshotResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group info",
			},

			isSnapshotSourceVolume: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot source volume id",
			},
			isSnapshotSourceImage: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If present, the image id from which the data on this volume was most directly provisioned.",
			},

			isSnapshotOperatingSystem: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name for the operating system included in this image",
			},

			isSnapshotBootable: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if a boot volume attachment can be created with a volume created from this snapshot",
			},

			isSnapshotLCState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot lifecycle state",
			},
			isSnapshotCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			isSnapshotEncryption: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encryption type of the snapshot",
			},
			isSnapshotHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for the snapshot",
			},

			isSnapshotMinCapacity: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum capacity of the snapshot",
			},
			isSnapshotResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type of the snapshot",
			},

			isSnapshotSize: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the snapshot",
			},
		},
	}
}

func dataSourceIBMISSnapshotValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isSnapshotName,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString})

	ibmISSnapshotDataSourceValidator := ResourceValidator{ResourceName: "ibm_is_snapshot", Schema: validateSchema}
	return &ibmISSnapshotDataSourceValidator
}

func dataSourceIBMISSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get(isSnapshotName).(string)
	id := d.Get("identifier").(string)
	err := snapshotGetByNameOrID(d, meta, name, id)
	if err != nil {
		return err
	}
	return nil
}

func snapshotGetByNameOrID(d *schema.ResourceData, meta interface{}, name, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if name != "" {
		start := ""
		allrecs := []vpcv1.Snapshot{}
		for {
			listSnapshotOptions := &vpcv1.ListSnapshotsOptions{}
			if start != "" {
				listSnapshotOptions.Start = &start
			}
			snapshots, response, err := sess.ListSnapshots(listSnapshotOptions)
			if err != nil {
				return fmt.Errorf("Error Fetching snapshots %s\n%s", err, response)
			}
			start = GetNext(snapshots.Next)
			allrecs = append(allrecs, snapshots.Snapshots...)
			if start == "" {
				break
			}
		}
		for _, snapshot := range allrecs {
			if *snapshot.Name == name || *snapshot.ID == id {
				d.SetId(*snapshot.ID)
				d.Set(isSnapshotName, *snapshot.Name)
				d.Set(isSnapshotHref, *snapshot.Href)
				d.Set(isSnapshotCRN, *snapshot.CRN)
				d.Set(isSnapshotMinCapacity, *snapshot.MinimumCapacity)
				d.Set(isSnapshotSize, *snapshot.Size)
				d.Set(isSnapshotEncryption, *snapshot.Encryption)
				d.Set(isSnapshotLCState, *snapshot.LifecycleState)
				d.Set(isSnapshotResourceType, *snapshot.ResourceType)
				d.Set(isSnapshotBootable, *snapshot.Bootable)
				if snapshot.ResourceGroup != nil && snapshot.ResourceGroup.ID != nil {
					d.Set(isSnapshotResourceGroup, *snapshot.ResourceGroup.ID)
				}
				if snapshot.SourceVolume != nil && snapshot.SourceVolume.ID != nil {
					d.Set(isSnapshotSourceVolume, *snapshot.SourceVolume.ID)
				}
				if snapshot.SourceImage != nil && snapshot.SourceImage.ID != nil {
					d.Set(isSnapshotSourceImage, *snapshot.SourceImage.ID)
				}
				if snapshot.OperatingSystem != nil && snapshot.OperatingSystem.Name != nil {
					d.Set(isSnapshotOperatingSystem, *snapshot.OperatingSystem.Name)
				}
				return nil
			}
		}
		return fmt.Errorf("No snapshot found with name %s", name)
	} else {
		getSnapshotOptions := &vpcv1.GetSnapshotOptions{
			ID: &id,
		}
		snapshot, response, err := sess.GetSnapshot(getSnapshotOptions)
		if err != nil {
			return fmt.Errorf("Error fetching snapshot %s\n%s", err, response)
		}
		if (response != nil && response.StatusCode == 404) || snapshot == nil {
			return fmt.Errorf("No snapshot found with id %s", id)
		}
		d.SetId(*snapshot.ID)
		d.Set(isSnapshotName, *snapshot.Name)
		d.Set(isSnapshotHref, *snapshot.Href)
		d.Set(isSnapshotCRN, *snapshot.CRN)
		d.Set(isSnapshotMinCapacity, *snapshot.MinimumCapacity)
		d.Set(isSnapshotSize, *snapshot.Size)
		d.Set(isSnapshotEncryption, *snapshot.Encryption)
		d.Set(isSnapshotLCState, *snapshot.LifecycleState)
		d.Set(isSnapshotResourceType, *snapshot.ResourceType)
		d.Set(isSnapshotBootable, *snapshot.Bootable)
		if snapshot.ResourceGroup != nil && snapshot.ResourceGroup.ID != nil {
			d.Set(isSnapshotResourceGroup, *snapshot.ResourceGroup.ID)
		}
		if snapshot.SourceVolume != nil && snapshot.SourceVolume.ID != nil {
			d.Set(isSnapshotSourceVolume, *snapshot.SourceVolume.ID)
		}
		if snapshot.SourceImage != nil && snapshot.SourceImage.ID != nil {
			d.Set(isSnapshotSourceImage, *snapshot.SourceImage.ID)
		}
		if snapshot.OperatingSystem != nil && snapshot.OperatingSystem.Name != nil {
			d.Set(isSnapshotOperatingSystem, *snapshot.OperatingSystem.Name)
		}
		return nil
	}
}
