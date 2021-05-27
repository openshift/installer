// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMComputeSSHKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMComputeSSHKeyCreate,
		Read:     resourceIBMComputeSSHKeyRead,
		Update:   resourceIBMComputeSSHKeyUpdate,
		Delete:   resourceIBMComputeSSHKeyDelete,
		Exists:   resourceIBMComputeSSHKeyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH Key label",
			},

			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "Plublic Key info",
			},

			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH key fingerprint",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Additional notes",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of tags for the resource",
			},
		},
	}
}

func resourceIBMComputeSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetSecuritySshKeyService(sess)

	// First check if the key exists by fingerprint
	// If so, set the Id (and fingerprint), but update notes and label (if any)
	key := d.Get("public_key").(string)
	label := d.Get("label").(string)

	fingerprint, err := computeSSHKeyFingerprint(key)
	if err != nil {
		return err
	}

	keys, err := services.GetAccountService(sess).
		Filter(filter.Path("sshKeys.fingerprint").Eq(fingerprint).Build()).
		GetSshKeys()
	if err == nil && len(keys) > 0 {
		slKey := keys[0]
		id := *slKey.Id
		slKey.Id = nil
		d.SetId(fmt.Sprintf("%d", id))
		d.Set("fingerprint", fingerprint)
		editKey := false

		notes := d.Get("notes").(string)
		if notes != "" && (slKey.Notes == nil || notes != *slKey.Notes) {
			slKey.Notes = sl.String(notes)
			editKey = true
		} else if slKey.Notes != nil {
			d.Set("notes", *slKey.Notes)
		}

		if label != *slKey.Label {
			slKey.Label = sl.String(label)
			editKey = true
		}

		if editKey {
			_, err = service.Id(id).EditObject(&slKey)
			return err
		}

		return nil
	} // End of "Import"

	// Build up our creation options
	opts := datatypes.Security_Ssh_Key{
		Label: sl.String(label),
		Key:   sl.String(key),
	}

	if notes, ok := d.GetOk("notes"); ok {
		opts.Notes = sl.String(notes.(string))
	}

	res, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating SSH Key: %s", err)
	}

	d.SetId(strconv.Itoa(*res.Id))
	log.Printf("[INFO] SSH Key: %d", *res.Id)

	return resourceIBMComputeSSHKeyRead(d, meta)
}

func resourceIBMComputeSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetSecuritySshKeyService(sess)

	keyID, _ := strconv.Atoi(d.Id())
	key, err := service.Id(keyID).GetObject()
	if err != nil {
		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if err, ok := err.(sl.Error); ok && err.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving SSH key: %s", err)
	}

	d.Set("label", key.Label)
	d.Set("public_key", key.Key)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("notes", key.Notes)
	return nil
}

func resourceIBMComputeSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetSecuritySshKeyService(sess)

	keyID, _ := strconv.Atoi(d.Id())

	key, err := service.Id(keyID).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving SSH key: %s", err)
	}

	if d.HasChange("label") {
		key.Label = sl.String(d.Get("label").(string))
	}

	if d.HasChange("notes") {
		key.Notes = sl.String(d.Get("notes").(string))
	}

	_, err = service.Id(keyID).EditObject(&key)
	if err != nil {
		return fmt.Errorf("Error editing SSH key: %s", err)
	}
	return resourceIBMComputeSSHKeyRead(d, meta)
}

func resourceIBMComputeSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetSecuritySshKeyService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SSH Key: %s", err)
	}

	log.Printf("[INFO] Deleting SSH key: %d", id)
	_, err = service.Id(id).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting SSH key: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMComputeSSHKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetSecuritySshKeyService(sess)

	keyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(keyID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}
	return result.Id != nil && *result.Id == keyID, nil
}

func computeSSHKeyFingerprint(key string) (fingerprint string, err error) {
	parts := strings.Fields(key)
	if len(parts) < 2 {
		return "", fmt.Errorf("Invalid public key specified :%s\nPlease check the value of public_key", key)
	}
	k, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("Error decoding the public key: %s\nPlease check the value of public_key", err)
	}
	fp := sha256.Sum256([]byte(k))
	prints := make([]string, len(fp))
	for i, b := range fp {
		prints[i] = fmt.Sprintf("%02x", b)
	}
	fingerprint = strings.Join(prints, ":")
	return
}
