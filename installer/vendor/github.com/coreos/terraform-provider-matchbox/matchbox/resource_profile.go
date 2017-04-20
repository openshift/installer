package matchbox

import (
	"context"

	matchbox "github.com/coreos/matchbox/matchbox/client"
	"github.com/coreos/matchbox/matchbox/server/serverpb"
	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileCreate,
		Read:   resourceProfileRead,
		Delete: resourceProfileDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_linux_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kernel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"initrd": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"args": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// resourceProfileCreate creates a Profile and its associated configs. Partial
// creates do not modify state and can be retried safely.
func resourceProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	// Profile
	name := d.Get("name").(string)
	clcName := containerLinuxConfigName(name)

	var initrds []string
	for _, initrd := range d.Get("initrd").([]interface{}) {
		initrds = append(initrds, initrd.(string))
	}
	var args []string
	for _, arg := range d.Get("args").([]interface{}) {
		args = append(args, arg.(string))
	}
	profile := &storagepb.Profile{
		Id:         name,
		IgnitionId: clcName,
		Boot: &storagepb.NetBoot{
			Kernel: d.Get("kernel").(string),
			Initrd: initrds,
			Args:   args,
		},
	}
	_, err := client.Profiles.ProfilePut(ctx, &serverpb.ProfilePutRequest{
		Profile: profile,
	})
	if err != nil {
		return err
	}

	// Container Linux Config
	_, err = client.Ignition.IgnitionPut(ctx, &serverpb.IgnitionPutRequest{
		Name:   clcName,
		Config: []byte(d.Get("container_linux_config").(string)),
	})
	if err != nil {
		return err
	}

	d.SetId(profile.GetId())
	return err
}

func resourceProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	name := d.Get("name").(string)
	clcName := containerLinuxConfigName(name)

	_, err := client.Profiles.ProfileGet(ctx, &serverpb.ProfileGetRequest{
		Id: name,
	})
	if err != nil {
		// resource doesn't exist or is corrupted
		d.SetId("")
		return nil
	}

	_, err = client.Ignition.IgnitionGet(ctx, &serverpb.IgnitionGetRequest{
		Name: clcName,
	})
	if err != nil {
		// resource doesn't exist or is corrupted
		d.SetId("")
		return nil
	}

	return nil
}

// resourceProfileDelete deletes a Profile and its associated configs. Partial
// deletes leave state unchanged and can be retried (deleting resources which
// no longer exist is a no-op).
func resourceProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*matchbox.Client)
	ctx := context.TODO()

	// Profile
	name := d.Get("name").(string)
	_, err := client.Profiles.ProfileDelete(ctx, &serverpb.ProfileDeleteRequest{
		Id: name,
	})
	if err != nil {
		return err
	}

	// Container Linux Config
	clcName := containerLinuxConfigName(name)
	_, err = client.Ignition.IgnitionDelete(ctx, &serverpb.IgnitionDeleteRequest{
		Name: clcName,
	})
	if err != nil {
		return err
	}

	// resource can be destroyed in state
	d.SetId("")
	return nil
}

func containerLinuxConfigName(name string) string {
	return name + ".yaml.tmpl"
}
