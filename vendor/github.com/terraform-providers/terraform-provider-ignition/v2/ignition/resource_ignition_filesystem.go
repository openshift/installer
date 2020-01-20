package ignition

import (
	"encoding/json"

	"github.com/coreos/ignition/v2/config/v3_0/types"
	vcontext_path "github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceFilesystem() *schema.Resource {
	return &schema.Resource{
		Exists: resourceFilesystemExists,
		Read:   resourceFilesystemRead,
		Schema: map[string]*schema.Schema{
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wipe_filesystem": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFilesystemRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildFilesystem(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceFilesystemExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildFilesystem(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildFilesystem(d *schema.ResourceData) (string, error) {
	fs := &types.Filesystem{
		Device: d.Get("device").(string),
	}
	path, hasPath := d.GetOk("path")
	if hasPath {
		str := path.(string)
		fs.Path = &str
	}

	format, hasFormat := d.GetOk("format")
	if hasFormat {
		str := format.(string)
		fs.Format = &str
	}

	wipe, hasWipeFilesystem := d.GetOk("wipe_filesystem")
	if hasWipeFilesystem {
		wp := wipe.(bool)
		fs.WipeFilesystem = &wp
	}

	label, hasLabel := d.GetOk("label")
	if hasLabel {
		str := label.(string)
		fs.Label = &str
	}

	uuid, hasUUID := d.GetOk("uuid")
	if hasUUID {
		str := uuid.(string)
		fs.UUID = &str
	}

	options, hasOptions := d.GetOk("options")
	if hasOptions {
		fs.Options = castSliceInterfaceToMountOption(options.([]interface{}))
	}

	b, err := json.Marshal(fs)
	if err != nil {
		return "", err
	}
	d.Set("rendered", string(b))

	return hash(string(b)), handleReport(fs.Validate(vcontext_path.ContextPath{}))
}

func castSliceInterfaceToMountOption(i []interface{}) []types.FilesystemOption {
	var o []types.FilesystemOption
	for _, value := range i {
		if value == nil {
			continue
		}

		o = append(o, types.FilesystemOption(value.(string)))
	}

	return o
}
