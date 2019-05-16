package ignition

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/coreos/ignition/config/v2_1/types"
)

var configReferenceResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"source": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"verification": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
		},
	},
}

func dataSourceConfig() *schema.Resource {
	return &schema.Resource{
		Exists: resourceIgnitionFileExists,
		Read:   resourceIgnitionFileRead,
		Schema: map[string]*schema.Schema{
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arrays": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filesystems": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"files": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"directories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"links": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"systemd": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"networkd": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"replace": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem:     configReferenceResource,
			},
			"append": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem:     configReferenceResource,
			},
			"rendered": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIgnitionFileRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderConfig(d, globalCache)
	if err != nil {
		return err
	}

	if err := d.Set("rendered", rendered); err != nil {
		return err
	}

	d.SetId(hash(rendered))
	return nil
}

func resourceIgnitionFileExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rendered, err := renderConfig(d, globalCache)
	if err != nil {
		return false, err
	}

	return hash(rendered) == d.Id(), nil
}

func renderConfig(d *schema.ResourceData, c *cache) (string, error) {
	i, err := buildConfig(d, c)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(i)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func buildConfig(d *schema.ResourceData, c *cache) (*types.Config, error) {
	var err error
	config := &types.Config{}
	config.Ignition, err = buildIgnition(d)
	if err != nil {
		return nil, err
	}

	config.Storage, err = buildStorage(d, c)
	if err != nil {
		return nil, err
	}

	config.Systemd, err = buildSystemd(d, c)
	if err != nil {
		return nil, err
	}

	config.Networkd, err = buildNetworkd(d, c)
	if err != nil {
		return nil, err
	}

	config.Passwd, err = buildPasswd(d, c)
	if err != nil {
		return nil, err
	}

	return config, handleReport(config.Validate())
}

func buildIgnition(d *schema.ResourceData) (types.Ignition, error) {
	var err error

	i := types.Ignition{}
	i.Version = types.MaxVersion.String()

	rr := d.Get("replace.0").(map[string]interface{})
	if len(rr) != 0 {
		i.Config.Replace, err = buildConfigReference(rr)
		if err != nil {
			return i, err
		}
	}

	ar := d.Get("append").([]interface{})
	if len(ar) != 0 {
		for _, rr := range ar {
			r, err := buildConfigReference(rr.(map[string]interface{}))
			if err != nil {
				return i, err
			}

			i.Config.Append = append(i.Config.Append, *r)
		}
	}

	return i, nil
}

func buildConfigReference(raw map[string]interface{}) (*types.ConfigReference, error) {
	r := &types.ConfigReference{}
	r.Source = raw["source"].(string)

	hash := raw["verification"].(string)
	if hash != "" {
		r.Verification.Hash = &hash
	}

	return r, nil
}

func buildStorage(d *schema.ResourceData, c *cache) (types.Storage, error) {
	storage := types.Storage{}

	for _, id := range d.Get("disks").([]interface{}) {
		if id == nil {
			continue
		}
		d, err := c.getDisk(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid disk %q, failed to get disk id: %v", id, err)
		}

		storage.Disks = append(storage.Disks, *d)
	}

	for _, id := range d.Get("arrays").([]interface{}) {
		if id == nil {
			continue
		}
		a, err := c.getRaid(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid raid %q, failed to get disk id: %v", id, err)
		}

		storage.Raid = append(storage.Raid, *a)
	}

	for _, id := range d.Get("filesystems").([]interface{}) {
		if id == nil {
			continue
		}
		f, err := c.getFilesystem(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid filesystem %q, failed to get filesystem id: %v", id, err)
		}

		storage.Filesystems = append(storage.Filesystems, *f)
	}

	for _, id := range d.Get("files").([]interface{}) {
		if id == nil {
			continue
		}
		f, err := c.getFile(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid file %q, failed to get file id: %v", id, err)
		}

		storage.Files = append(storage.Files, *f)
	}

	for _, id := range d.Get("directories").([]interface{}) {
		if id == nil {
			continue
		}
		f, err := c.getDirectory(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid file %q, failed to get directory id: %v", id, err)
		}

		storage.Directories = append(storage.Directories, *f)
	}

	for _, id := range d.Get("links").([]interface{}) {
		if id == nil {
			continue
		}
		f, err := c.getLink(id.(string))
		if err != nil {
			return storage, fmt.Errorf("invalid file %q, failed to get link id: %v", id, err)
		}

		storage.Links = append(storage.Links, *f)
	}

	return storage, nil

}

func buildSystemd(d *schema.ResourceData, c *cache) (types.Systemd, error) {
	systemd := types.Systemd{}

	for _, id := range d.Get("systemd").([]interface{}) {
		if id == nil {
			continue
		}

		u, err := c.getSystemdUnit(id.(string))
		if err != nil {
			return systemd, fmt.Errorf("invalid systemd unit %q, failed to get systemd unit id: %v", id, err)
		}

		systemd.Units = append(systemd.Units, *u)
	}

	return systemd, nil

}

func buildNetworkd(d *schema.ResourceData, c *cache) (types.Networkd, error) {
	networkd := types.Networkd{}

	for _, id := range d.Get("networkd").([]interface{}) {
		if id == nil {
			continue
		}

		u, err := c.getNetworkdunit(id.(string))
		if err != nil {
			return networkd, fmt.Errorf("invalid networkd unit %q, failed to get networkd unit id: %v", id, err)
		}

		networkd.Units = append(networkd.Units, *u)
	}

	return networkd, nil
}

func buildPasswd(d *schema.ResourceData, c *cache) (types.Passwd, error) {
	passwd := types.Passwd{}

	for _, id := range d.Get("users").([]interface{}) {
		if id == nil {
			continue
		}

		u, err := c.getUser(id.(string))
		if err != nil {
			return passwd, fmt.Errorf("invalid user %q, failed to get user id: %v", id, err)
		}

		passwd.Users = append(passwd.Users, *u)
	}

	for _, id := range d.Get("groups").([]interface{}) {
		if id == nil {
			continue
		}

		g, err := c.getGroup(id.(string))
		if err != nil {
			return passwd, fmt.Errorf("invalid group %q, failed to get group id: %v", id, err)
		}

		passwd.Groups = append(passwd.Groups, *g)
	}

	return passwd, nil

}
