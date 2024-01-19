package ignition

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/coreos/ignition/v2/config/v3_3/types"
	"github.com/coreos/vcontext/path"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Exists: resourceFileExists,
		Read:   resourceFileRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"content": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mime": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "text/plain",
						},

						"content": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"compression": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"verification": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
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

func resourceFileRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildFile(d)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceFileExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildFile(d)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildFile(d *schema.ResourceData) (string, error) {
	_, hasContent := d.GetOk("content")
	_, hasSource := d.GetOk("source")
	if hasContent && hasSource {
		return "", fmt.Errorf("content and source options are incompatible")
	}

	if !hasContent && !hasSource {
		return "", fmt.Errorf("content or source options must be present")
	}

	var contents types.Resource
	if hasContent {
		s := encodeDataURL(
			d.Get("content.0.mime").(string),
			d.Get("content.0.content").(string),
		)
		contents.Source = &s
	}

	if hasSource {
		src := d.Get("source.0.source").(string)
		if src != "" {
			contents.Source = &src
		}
		compression := d.Get("source.0.compression").(string)
		if compression != "" {
			contents.Compression = &compression
		}
		h := d.Get("source.0.verification").(string)
		if h != "" {
			contents.Verification.Hash = &h
		}
	}

	file := &types.File{}

	file.Path = d.Get("path").(string)

	overwrite := d.Get("overwrite").(bool)
	file.Overwrite = &overwrite

	file.Contents = contents

	mode, hasMode := d.GetOk("mode")
	if hasMode {
		imode := mode.(int)
		file.Mode = &imode
	}

	uid := d.Get("uid").(int)
	if uid != 0 {
		file.User = types.NodeUser{ID: &uid}
	}

	gid := d.Get("gid").(int)
	if gid != 0 {
		file.Group = types.NodeGroup{ID: &gid}
	}

	if err := handleReport(file.Validate(path.ContextPath{})); err != nil {
		return "", err
	}

	b, err := json.Marshal(file)
	if err != nil {
		return "", err
	}
	err = d.Set("rendered", string(b))
	if err != nil {
		return "", err
	}

	return hash(string(b)), nil
}

func encodeDataURL(mime, content string) string {
	base64 := base64.StdEncoding.EncodeToString([]byte(content))
	return fmt.Sprintf("data:%s;charset=utf-8;base64,%s", mime, base64)
}
