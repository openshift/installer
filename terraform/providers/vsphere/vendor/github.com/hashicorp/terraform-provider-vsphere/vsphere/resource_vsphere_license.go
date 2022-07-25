package vsphere

import (
	"errors"
	"fmt"
	"log"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	// ErrNoSuchKeyFound is an error primarily thrown by the Read method of the resource.
	// The error doesn't display the key itself for security reasons.
	ErrNoSuchKeyFound = errors.New("The key was not found")
	// ErrKeyCannotBeDeleted is an error which occurs when a key that is used by VMs is
	// being removed
	ErrKeyCannotBeDeleted = errors.New("The key wasn't deleted")
)

func resourceVSphereLicense() *schema.Resource {
	return &schema.Resource{

		SchemaVersion: 1,

		Create: resourceVSphereLicenseCreate,
		Read:   resourceVSphereLicenseRead,
		Update: resourceVSphereLicenseUpdate,
		Delete: resourceVSphereLicenseDelete,

		Schema: map[string]*schema.Schema{
			"license_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// computed properties returned by the API
			"edition_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVSphereLicenseCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] Running the create method")

	client := meta.(*Client).vimClient
	manager := license.NewManager(client.Client)

	key := d.Get("license_key").(string)

	log.Println(" [INFO] Reading the key from the resource data")
	var labelMap map[string]interface{}
	if labels, ok := d.GetOk("labels"); ok {
		labelMap = labels.(map[string]interface{})
	}

	var info types.LicenseManagerLicenseInfo
	var err error
	switch t := client.ServiceContent.About.ApiType; t {
	case "HostAgent":
		// Labels are not allowed in ESXi
		if len(labelMap) != 0 {
			return errors.New("Labels are not allowed in ESXi")
		}
		info, err = manager.Update(context.TODO(), key, nil)

	case "VirtualCenter":
		info, err = manager.Add(context.TODO(), key, nil)
		if err != nil {
			return err
		}
		err = updateLabels(manager, key, labelMap)

	default:
		return fmt.Errorf("unsupported ApiType: %s", t)
	}

	if err != nil {
		return err
	}

	if err = DecodeError(info); err != nil {
		return err
	}

	// This can be used in the read method to set the computed parameters
	d.SetId(info.LicenseKey)

	return resourceVSphereLicenseRead(d, meta)
}

func resourceVSphereLicenseRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] Running the read method")

	client := meta.(*Client).vimClient
	manager := license.NewManager(client.Client)

	if info := getLicenseInfoFromKey(d.Get("license_key").(string), manager); info != nil {
		log.Println("[INFO] Setting the values")
		_ = d.Set("edition_key", info.EditionKey)
		_ = d.Set("total", info.Total)
		_ = d.Set("used", info.Used)
		_ = d.Set("name", info.Name)
		_ = d.Set("labels", keyValuesToMap(info.Labels))
	} else {
		return ErrNoSuchKeyFound
	}

	return nil
}

// resourceVSphereLicenseUpdate check for change in labels of the key and updates them.
func resourceVSphereLicenseUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] Running the update method")

	client := meta.(*Client).vimClient
	manager := license.NewManager(client.Client)

	if key, ok := d.GetOk("license_key"); ok {
		licenseKey := key.(string)
		if !isKeyPresent(licenseKey, manager) {
			return ErrNoSuchKeyFound
		}

		if d.HasChange("labels") {
			labelMap := d.Get("labels").(map[string]interface{})

			err := updateLabels(manager, licenseKey, labelMap)
			if err != nil {
				return err
			}
		}
	}

	return resourceVSphereLicenseRead(d, meta)
}

func updateLabels(manager *license.Manager, licenseKey string, labelMap map[string]interface{}) error {
	for key, value := range labelMap {
		err := UpdateLabel(context.TODO(), manager, licenseKey, key, value.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceVSphereLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] Running the delete method")

	client := meta.(*Client).vimClient
	manager := license.NewManager(client.Client)

	if key := d.Get("license_key").(string); isKeyPresent(key, manager) {
		err := manager.Remove(context.TODO(), key)

		if err != nil {
			return err
		}

		// if the key is still present
		if isKeyPresent(key, manager) {
			return ErrKeyCannotBeDeleted
		}
		d.SetId("")
		return nil
	}
	return ErrNoSuchKeyFound
}

func getLicenseInfoFromKey(key string, manager *license.Manager) *types.LicenseManagerLicenseInfo {
	// Use of decode is not returning labels so using list instead
	// Issue - https://github.com/vmware/govmomi/issues/797
	infoList, _ := manager.List(context.TODO())
	for _, info := range infoList {
		if info.LicenseKey == key {
			return &info
		}
	}
	return nil
}

// isKeyPresent iterates over the InfoList to check if the license is present or not.
func isKeyPresent(key string, manager *license.Manager) bool {
	infoList, _ := manager.List(context.TODO())

	for _, info := range infoList {
		if info.LicenseKey == key {
			return true
		}
	}

	return false
}

// UpdateLabel provides a wrapper around the UpdateLabel data objects
func UpdateLabel(ctx context.Context, m *license.Manager, licenseKey string, key string, val string) error {
	req := types.UpdateLicenseLabel{
		This:       m.Reference(),
		LicenseKey: licenseKey,
		LabelKey:   key,
		LabelValue: val,
	}

	_, err := methods.UpdateLicenseLabel(ctx, m.Client(), &req)
	return err
}

// DecodeError tries to find a specific error which occurs when an invalid key is passed
// to the server
func DecodeError(info types.LicenseManagerLicenseInfo) error {
	for _, property := range info.Properties {
		if property.Key == "diagnostic" {
			return errors.New(property.Value.(string))
		}
	}

	return nil
}

func keyValuesToMap(keyValues []types.KeyValue) map[string]interface{} {
	KVMap := make(map[string]interface{})
	for _, keyValue := range keyValues {
		KVMap[keyValue.Key] = keyValue.Value
	}
	return KVMap
}
