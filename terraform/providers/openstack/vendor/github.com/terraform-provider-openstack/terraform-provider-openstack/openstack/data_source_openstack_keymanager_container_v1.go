package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
)

func dataSourceKeyManagerContainerV1() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerContainerV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secret_refs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_ref": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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

			"acl": {
				Type:     schema.TypeList,
				Computed: true,
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

func dataSourceKeyManagerContainerV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	listOpts := containers.ListOpts{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Containers List Options: %#v", listOpts)

	allPages, err := containers.List(kmClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_keymanager_container_v1 containers: %s", err)
	}

	allContainers, err := containers.ExtractContainers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_keymanager_container_v1 containers: %s", err)
	}

	if len(allContainers) < 1 {
		return diag.Errorf("Your query returned no openstack_keymanager_container_v1 results. " +
			"Please change your search criteria and try again.")
	}

	if len(allContainers) > 1 {
		log.Printf("[DEBUG] Multiple openstack_keymanager_container_v1 results found: %#v", allContainers)
		return diag.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	container := allContainers[0]

	log.Printf("[DEBUG] Retrieved openstack_keymanager_container_v1 %s: %#v", d.Id(), container)

	uuid := keyManagerContainerV1GetUUIDfromContainerRef(container.ContainerRef)

	d.SetId(uuid)
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
		log.Printf("[DEBUG] Unable to get %s container acls: %s", uuid, err)
	}
	d.Set("acl", flattenKeyManagerV1ACLs(acl))

	// Set the region
	d.Set("region", GetRegion(d, config))

	return nil
}
