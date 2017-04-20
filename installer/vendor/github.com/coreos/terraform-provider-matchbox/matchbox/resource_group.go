package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"selector": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				ForceNew: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				ForceNew: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)

	selectors := map[string]string{}
	for k, v := range d.Get("selector").(map[string]interface{}) {
		selectors[k] = v.(string)
	}

	richGroup := &storagepb.RichGroup{
		Id:       name,
		Profile:  d.Get("profile").(string),
		Selector: selectors,
		Metadata: d.Get("metadata").(map[string]interface{}),
	}
	group, err := richGroup.ToGroup()
	if err != nil {
		return err
	}

	_, err = client.Groups.GroupPut(ctx, &serverpb.GroupPutRequest{
		Group: group,
	})
	if err != nil {
		return err
	}

	d.SetId(group.GetId())
	return err
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)
	_, err := client.Groups.GroupGet(ctx, &serverpb.GroupGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist anymore
		d.SetId("")
		return nil
	}
	return err
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)
	_, err := client.Groups.GroupDelete(ctx, &serverpb.GroupDeleteRequest{
		Id: name,
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
