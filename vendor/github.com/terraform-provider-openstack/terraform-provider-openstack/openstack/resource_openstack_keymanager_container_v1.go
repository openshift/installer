package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKeyManagerContainerV1() *schema.Resource {
	ret := &schema.Resource{
		Create: resourceKeyManagerContainerV1Create,
		Read:   resourceKeyManagerContainerV1Read,
		Update: resourceKeyManagerContainerV1Update,
		Delete: resourceKeyManagerContainerV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"generic", "rsa", "certificate",
				}, false),
			},

			"secret_refs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_ref": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"acl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
			},

			"container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	elem := &schema.Resource{
		Schema: make(map[string]*schema.Schema),
	}
	for _, aclOp := range getSupportedACLOperations() {
		elem.Schema[aclOp] = getACLSchema()
	}
	ret.Schema["acl"].Elem = elem

	return ret
}

func resourceKeyManagerContainerV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	containerType := keyManagerContainerV1Type(d.Get("type").(string))

	createOpts := containers.CreateOpts{
		Name:       d.Get("name").(string),
		Type:       containerType,
		SecretRefs: expandKeyManagerContainerV1SecretRefs(d.Get("secret_refs").(*schema.Set)),
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_container_v1: %#v", createOpts)

	container, err := containers.Create(kmClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_keymanager_container_v1: %s", err)
	}

	uuid := keyManagerContainerV1GetUUIDfromContainerRef(container.ContainerRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    keyManagerContainerV1WaitForContainerCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for openstack_keymanager_container_v1: %s", err)
	}

	d.SetId(uuid)

	d.Partial(true)

	// set the acl first before setting the secret refs
	if _, ok := d.GetOk("acl"); ok {
		setOpts := expandKeyManagerV1ACLs(d.Get("acl"))
		_, err = acls.SetContainerACL(kmClient, uuid, setOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error settings ACLs for the openstack_keymanager_container_v1: %s", err)
		}
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for openstack_keymanager_container_v1: %s", err)
	}

	d.Partial(false)

	return resourceKeyManagerContainerV1Read(d, meta)
}

func resourceKeyManagerContainerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	container, err := containers.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_keymanager_container_v1")
	}

	log.Printf("[DEBUG] Retrieved openstack_keymanager_container_v1 %s: %#v", d.Id(), container)

	d.Set("name", container.Name)

	d.Set("creator_id", container.CreatorID)
	d.Set("container_ref", container.ContainerRef)
	d.Set("type", container.Type)
	d.Set("status", container.Status)
	d.Set("created_at", container.Created.Format(time.RFC3339))
	d.Set("updated_at", container.Updated.Format(time.RFC3339))
	d.Set("consumers", flattenKeyManagerContainerV1Consumers(container.Consumers))

	d.Set("secret_refs", flattenKeyManagerContainerV1SecretRefs(container.SecretRefs))

	acl, err := acls.GetContainerACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s container acls: %s", d.Id(), err)
	}
	d.Set("acl", flattenKeyManagerV1ACLs(acl))

	// Set the region
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKeyManagerContainerV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	if d.HasChange("acl") {
		updateOpts := expandKeyManagerV1ACLs(d.Get("acl"))
		_, err := acls.UpdateContainerACL(kmClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_keymanager_container_v1 %s acl: %s", d.Id(), err)
		}
	}

	if d.HasChange("secret_refs") {
		o, n := d.GetChange("secret_refs")

		oldRefs, newRefs := o.(*schema.Set), n.(*schema.Set)
		addRefs := newRefs.Difference(oldRefs)
		delRefs := oldRefs.Difference(newRefs)

		// delete old references first
		for _, delRef := range expandKeyManagerContainerV1SecretRefs(delRefs) {
			res := containers.DeleteSecretRef(kmClient, d.Id(), delRef)
			if res.Err != nil {
				if _, ok := res.Err.(gophercloud.ErrDefault404); !ok {
					return fmt.Errorf("Error removing old %s secret reference from the %s container: %s", delRef.Name, d.Id(), res.Err)
				}
			}
		}

		// then add new ones
		for _, addRef := range expandKeyManagerContainerV1SecretRefs(addRefs) {
			res := containers.CreateSecretRef(kmClient, d.Id(), addRef)
			if res.Err != nil {
				return fmt.Errorf("Error adding new %s secret reference to the %s container: %s", addRef.Name, d.Id(), res.Err)
			}
		}
	}

	return resourceKeyManagerContainerV1Read(d, meta)
}

func resourceKeyManagerContainerV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DELETED"},
		Refresh:    keyManagerContainerV1WaitForContainerDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}
