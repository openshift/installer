package ignition

import (
	"github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Exists: resourceUserExists,
		Read:   resourceUserRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password_hash": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ssh_authorized_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"gecos": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"home_dir": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"no_create_home": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"primary_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"no_user_group": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"no_log_init": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"shell": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	id, err := buildUser(d, globalCache)
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildUser(d, globalCache)
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func buildUser(d *schema.ResourceData, c *cache) (string, error) {
	user := types.PasswdUser{
		Name:   d.Get("name").(string),
		UID:    getInt(d, "uid"),
		Groups: castSliceInterfaceToPasswdUserGroup(d.Get("groups").([]interface{})),
		SSHAuthorizedKeys: castSliceInterfaceToSSHAuthorizedKey(
			d.Get("ssh_authorized_keys").([]interface{}),
		),
	}

	pwd := d.Get("password_hash").(string)
	if pwd != "" {
		user.PasswordHash = &pwd
	}

	gecos := d.Get("gecos").(string)
	if gecos != "" {
		user.Gecos = &gecos
	}

	homedir := d.Get("home_dir").(string)
	if homedir != "" {
		user.HomeDir = &homedir
	}

	primaryGroup := d.Get("primary_group").(string)
	if primaryGroup != "" {
		user.PrimaryGroup = &primaryGroup
	}

	shell := d.Get("shell").(string)
	if shell != "" {
		user.Shell = &shell
	}

	nocreatehome, hasNocreatehome := d.GetOk("no_create_home")
	if hasNocreatehome {
		bnocreatehome := nocreatehome.(bool)
		user.NoCreateHome = &bnocreatehome
	}

	nousergroup, hasNousergroup := d.GetOk("no_user_group")
	if hasNousergroup {
		bnousergroup := nousergroup.(bool)
		user.NoUserGroup = &bnousergroup
	}

	nologinit, hasNologinit := d.GetOk("no_log_init")
	if hasNologinit {
		bnologinit := nologinit.(bool)
		user.NoLogInit = &bnologinit
	}

	system, hasSystem := d.GetOk("system")
	if hasSystem {
		bsystem := system.(bool)
		user.System = &bsystem
	}

	return c.addUser(&user), nil
}

func castSliceInterfaceToPasswdUserGroup(i []interface{}) []types.Group {
	var res []types.Group
	for _, g := range i {
		if g == nil {
			continue
		}

		res = append(res, types.Group(g.(string)))
	}
	return res
}

func castSliceInterfaceToSSHAuthorizedKey(i []interface{}) []types.SSHAuthorizedKey {
	var res []types.SSHAuthorizedKey
	for _, k := range i {
		if k == nil {
			continue
		}

		res = append(res, types.SSHAuthorizedKey(k.(string)))
	}
	return res
}
