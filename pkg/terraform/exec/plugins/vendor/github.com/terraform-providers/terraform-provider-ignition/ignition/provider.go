package ignition

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/coreos/ignition/config/v2_1/types"
	"github.com/coreos/ignition/config/validate/report"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/peterbourgon/diskv"
)

// globalCache keeps the instances of the internal types of ignition generated
// by the different data resources with the goal to be reused by the
// ignition_config data resource. The key of the maps are a hash of the types
// calculated on the type serialized to JSON.
var globalCache = &cache{
	disks:         make(map[string]*types.Disk, 0),
	arrays:        make(map[string]*types.Raid, 0),
	filesystems:   make(map[string]*types.Filesystem, 0),
	files:         make(map[string]*types.File, 0),
	directories:   make(map[string]*types.Directory, 0),
	links:         make(map[string]*types.Link, 0),
	systemdUnits:  make(map[string]*types.Unit, 0),
	networkdUnits: make(map[string]*types.Networkdunit, 0),
	users:         make(map[string]*types.PasswdUser, 0),
	groups:        make(map[string]*types.PasswdGroup, 0),
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"ignition_config":        dataSourceConfig(),
			"ignition_disk":          dataSourceDisk(),
			"ignition_raid":          dataSourceRaid(),
			"ignition_filesystem":    dataSourceFilesystem(),
			"ignition_file":          dataSourceFile(),
			"ignition_directory":     dataSourceDirectory(),
			"ignition_link":          dataSourceLink(),
			"ignition_systemd_unit":  dataSourceSystemdUnit(),
			"ignition_networkd_unit": dataSourceNetworkdUnit(),
			"ignition_user":          dataSourceUser(),
			"ignition_group":         dataSourceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	basePath := filepath.Join(os.TempDir(), "terraform-provider-ignition")

	globalCache.d = diskv.New(diskv.Options{
		BasePath:     basePath,
		CacheSizeMax: 100 * 1024 * 1024, // 100MB
	})
	return nil, nil
}

type cache struct {
	// quick retrieval without marshal/unmarshal overhead
	disks         map[string]*types.Disk
	arrays        map[string]*types.Raid
	filesystems   map[string]*types.Filesystem
	files         map[string]*types.File
	directories   map[string]*types.Directory
	links         map[string]*types.Link
	systemdUnits  map[string]*types.Unit
	networkdUnits map[string]*types.Networkdunit
	users         map[string]*types.PasswdUser
	groups        map[string]*types.PasswdGroup

	// cached retrieval from disk and marshal/unmarshal overhead
	d *diskv.Diskv

	sync.Mutex
}

func (c *cache) addDisk(g *types.Disk) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(g)
	c.disks[id] = g
	c.d.Write(keyFromType("disks", id), raw)
	return id
}

func (c *cache) getDisk(id string) (*types.Disk, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.disks[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("disks", id))
	if err != nil {
		return nil, err
	}
	out := &types.Disk{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.disks[id] = out
	return out, nil
}

func (c *cache) addRaid(r *types.Raid) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(r)
	c.arrays[id] = r
	c.d.Write(keyFromType("arrays", id), raw)
	return id
}

func (c *cache) getRaid(id string) (*types.Raid, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.arrays[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("arrays", id))
	if err != nil {
		return nil, err
	}
	out := &types.Raid{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.arrays[id] = out
	return out, nil
}

func (c *cache) addFilesystem(f *types.Filesystem) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(f)
	c.filesystems[id] = f
	c.d.Write(keyFromType("fs", id), raw)
	return id
}

func (c *cache) getFilesystem(id string) (*types.Filesystem, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.filesystems[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("fs", id))
	if err != nil {
		return nil, err
	}
	out := &types.Filesystem{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.filesystems[id] = out
	return out, nil
}

func (c *cache) addFile(f *types.File) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(f)
	c.files[id] = f
	c.d.Write(keyFromType("files", id), raw)
	return id
}

func (c *cache) getFile(id string) (*types.File, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.files[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("files", id))
	if err != nil {
		return nil, err
	}
	out := &types.File{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.files[id] = out
	return out, nil
}

func (c *cache) addDirectory(d *types.Directory) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(d)
	c.directories[id] = d
	c.d.Write(keyFromType("dirs", id), raw)
	return id
}

func (c *cache) getDirectory(id string) (*types.Directory, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.directories[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("dirs", id))
	if err != nil {
		return nil, err
	}
	out := &types.Directory{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.directories[id] = out
	return out, nil
}

func (c *cache) addLink(l *types.Link) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(l)
	c.links[id] = l
	c.d.Write(keyFromType("links", id), raw)
	return id
}

func (c *cache) getLink(id string) (*types.Link, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.links[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("links", id))
	if err != nil {
		return nil, err
	}
	out := &types.Link{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.links[id] = out
	return out, nil
}

func (c *cache) addSystemdUnit(u *types.Unit) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(u)
	c.systemdUnits[id] = u
	c.d.Write(keyFromType("sunits", id), raw)
	return id
}

func (c *cache) getSystemdUnit(id string) (*types.Unit, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.systemdUnits[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("sunits", id))
	if err != nil {
		return nil, err
	}
	out := &types.Unit{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.systemdUnits[id] = out
	return out, nil
}

func (c *cache) addNetworkdUnit(u *types.Networkdunit) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(u)
	c.networkdUnits[id] = u
	c.d.Write(keyFromType("nunits", id), raw)
	return id
}

func (c *cache) getNetworkdunit(id string) (*types.Networkdunit, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.networkdUnits[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("nunits", id))
	if err != nil {
		return nil, err
	}
	out := &types.Networkdunit{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.networkdUnits[id] = out
	return out, nil
}

func (c *cache) addUser(u *types.PasswdUser) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(u)
	c.users[id] = u
	c.d.Write(keyFromType("users", id), raw)
	return id
}

func (c *cache) getUser(id string) (*types.PasswdUser, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.users[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("users", id))
	if err != nil {
		return nil, err
	}
	out := &types.PasswdUser{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.users[id] = out
	return out, nil
}

func (c *cache) addGroup(g *types.PasswdGroup) string {
	c.Lock()
	defer c.Unlock()

	id, raw := id(g)
	c.groups[id] = g
	c.d.Write(keyFromType("groups", id), raw)
	return id
}

func (c *cache) getGroup(id string) (*types.PasswdGroup, error) {
	c.Lock()
	defer c.Unlock()

	if v, ok := c.groups[id]; ok {
		return v, nil
	}
	raw, err := c.d.Read(keyFromType("groups", id))
	if err != nil {
		return nil, err
	}
	out := &types.PasswdGroup{}
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	c.groups[id] = out
	return out, nil
}

// id returns the hashed and raw identifier for the input
func id(input interface{}) (string, []byte) {
	b, _ := json.Marshal(input)
	return hash(string(b)), b
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

func keyFromType(t string, id string) string { return fmt.Sprintf("%s-%s", t, id) }

func castSliceInterface(i []interface{}) []string {
	var o []string
	for _, value := range i {
		if value == nil {
			continue
		}

		o = append(o, value.(string))
	}

	return o
}

func getInt(d *schema.ResourceData, key string) *int {
	var i *int
	if value, ok := d.GetOk(key); ok {
		n := value.(int)
		i = &n
	}

	return i
}

func handleReport(r report.Report) error {
	for _, e := range r.Entries {
		debug(e.String())
	}

	if r.IsFatal() {
		return fmt.Errorf("invalid configuration:\n%s", r.String())
	}

	return nil
}

func debug(format string, a ...interface{}) {
	log.Printf("[DEBUG] %s", fmt.Sprintf(format, a...))
}
