package ignition

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/coreos/ignition/v2/config/v3_3/types"
	"github.com/coreos/ignition/v2/config/validate"
)

var httpHeaderReferenceResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
	},
}

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
		"http_headers": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem:     httpHeaderReferenceResource,
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
			"kernel_arguments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replace": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem:     configReferenceResource,
			},
			"merge": {
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
	rendered, err := renderConfig(d)
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
	rendered, err := renderConfig(d)
	if err != nil {
		return false, err
	}

	return hash(rendered) == d.Id(), nil
}

func renderConfig(d *schema.ResourceData) (string, error) {
	i, err := buildConfig(d)
	if err != nil {
		return "", err
	}

	bytes, err := json.Marshal(i)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func buildConfig(d *schema.ResourceData) (*types.Config, error) {
	var err error
	config := &types.Config{}
	config.Ignition, err = buildIgnition(d)
	if err != nil {
		return nil, err
	}

	config.Storage, err = buildStorage(d)
	if err != nil {
		return nil, err
	}

	config.Systemd, err = buildSystemd(d)
	if err != nil {
		return nil, err
	}

	config.Passwd, err = buildPasswd(d)
	if err != nil {
		return nil, err
	}

	config.KernelArguments, err = buildKernelArguments(d)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return config, handleReport(validate.ValidateWithContext(new(*types.Config), b))
}

func buildIgnition(d *schema.ResourceData) (types.Ignition, error) {
	i := types.Ignition{}
	i.Version = types.MaxVersion.String()

	rr := d.Get("replace.0").(map[string]interface{})
	if len(rr) != 0 {
		r, err := buildConfigReference(rr)
		if err != nil {
			return i, err
		}

		i.Config.Replace = *r
	}

	ar := d.Get("merge").([]interface{})
	if len(ar) != 0 {
		for _, rr := range ar {
			r, err := buildConfigReference(rr.(map[string]interface{}))
			if err != nil {
				return i, err
			}

			i.Config.Merge = append(i.Config.Merge, *r)
		}
	}

	return i, nil
}

func buildConfigReference(raw map[string]interface{}) (*types.Resource, error) {
	r := &types.Resource{}
	source := raw["source"].(string)
	if source != "" {
		r.Source = &source
	}
	hash := raw["verification"].(string)
	if hash != "" {
		r.Verification.Hash = &hash
	}

	for _, hh := range raw["http_headers"].([]interface{}) {
		h, err := buildConfigHTTPHeaderReference(hh.(map[string]interface{}))
		if err != nil {
			return r, err
		}
		r.HTTPHeaders = append(r.HTTPHeaders, h)
	}

	return r, nil
}

func buildConfigHTTPHeaderReference(raw map[string]interface{}) (types.HTTPHeader, error) {
	h := types.HTTPHeader{}
	name := raw["name"].(string)
	if name != "" {
		h.Name = name
	}
	value := raw["value"].(string)
	if value != "" {
		h.Value = &value
	}

	return h, nil
}

func buildStorage(d *schema.ResourceData) (types.Storage, error) {
	storage := types.Storage{}

	for _, disk := range d.Get("disks").([]interface{}) {
		if disk == nil {
			continue
		}

		d := types.Disk{}
		err := json.Unmarshal([]byte(disk.(string)), &d)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Disks = append(storage.Disks, d)
	}

	for _, array := range d.Get("arrays").([]interface{}) {
		if array == nil {
			continue
		}

		a := types.Raid{}
		err := json.Unmarshal([]byte(array.(string)), &a)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Raid = append(storage.Raid, a)
	}

	for _, fs := range d.Get("filesystems").([]interface{}) {
		if fs == nil {
			continue
		}

		f := types.Filesystem{}
		err := json.Unmarshal([]byte(fs.(string)), &f)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Filesystems = append(storage.Filesystems, f)
	}

	for _, file := range d.Get("files").([]interface{}) {
		if file == nil {
			continue
		}

		f := types.File{}
		err := json.Unmarshal([]byte(file.(string)), &f)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Files = append(storage.Files, f)
	}

	for _, dir := range d.Get("directories").([]interface{}) {
		if dir == nil {
			continue
		}

		f := types.Directory{}
		err := json.Unmarshal([]byte(dir.(string)), &f)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Directories = append(storage.Directories, f)
	}

	for _, link := range d.Get("links").([]interface{}) {
		if link == nil {
			continue
		}

		f := types.Link{}
		err := json.Unmarshal([]byte(link.(string)), &f)
		if err != nil {
			return storage, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		storage.Links = append(storage.Links, f)
	}

	return storage, nil

}
func buildKernelArguments(d *schema.ResourceData) (types.KernelArguments, error) {
	kargs := types.KernelArguments{}

	k := d.Get("kernel_arguments").(string)
	if k == "" {
		return kargs, nil
	}

	err := json.Unmarshal([]byte(k), &kargs)
	if err != nil {
		return kargs, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
	}

	return kargs, nil
}

func buildSystemd(d *schema.ResourceData) (types.Systemd, error) {
	systemd := types.Systemd{}

	for _, unit := range d.Get("systemd").([]interface{}) {
		if unit == nil {
			continue
		}

		u := types.Unit{}
		err := json.Unmarshal([]byte(unit.(string)), &u)
		if err != nil {
			return systemd, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		systemd.Units = append(systemd.Units, u)
	}

	return systemd, nil

}

func buildPasswd(d *schema.ResourceData) (types.Passwd, error) {
	passwd := types.Passwd{}

	for _, user := range d.Get("users").([]interface{}) {
		if user == nil {
			continue
		}

		u := types.PasswdUser{}
		err := json.Unmarshal([]byte(user.(string)), &u)
		if err != nil {
			return passwd, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		passwd.Users = append(passwd.Users, u)
	}

	for _, group := range d.Get("groups").([]interface{}) {
		if group == nil {
			continue
		}

		g := types.PasswdGroup{}
		err := json.Unmarshal([]byte(group.(string)), &g)
		if err != nil {
			return passwd, errors.Wrap(err, "No valid JSON found, make sure you're using .rendered and not .id")
		}

		passwd.Groups = append(passwd.Groups, g)
	}

	return passwd, nil

}
