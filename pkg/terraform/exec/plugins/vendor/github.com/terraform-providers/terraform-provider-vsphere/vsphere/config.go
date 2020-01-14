package vsphere

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/debug"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/vic/pkg/vsphere/tags"
)

// VSphereClient is the client connection manager for the vSphere provider. It
// holds the connections to the various API endpoints we need to interface
// with, such as the VMODL API through govmomi, and the REST SDK through
// alternate libraries.
type VSphereClient struct {
	// The VIM/govmomi client.
	vimClient *govmomi.Client

	// The specialized tags client SDK imported from vmware/vic.
	tagsClient *tags.RestClient
}

// TagsClient returns the embedded REST client used for tags, after determining
// if the connection is eligible:
//
// * The connection information in vimClient is valid vCenter connection
// * The provider has a connection to the CIS REST client. This is true if
// tagsClient != nil.
//
// This function should be used whenever possible to return the client from the
// provider meta variable for use, to determine if it can be used at all.
//
// The nil value that is returned on an unsupported connection can be
// considered stable behavior for read purposes on resources that need to be
// able to read tags if they are present. You can use the snippet below in a
// Read call to determine if tags are supported on this connection, and if they
// are, read them from the object and save them in the resource:
//
//   if tagsClient, _ := meta.(*VSphereClient).TagsClient(); tagsClient != nil {
//     if err := readTagsForResource(tagsClient, obj, d); err != nil {
//       return err
//     }
//   }
func (c *VSphereClient) TagsClient() (*tags.RestClient, error) {
	if err := viapi.ValidateVirtualCenter(c.vimClient); err != nil {
		return nil, err
	}
	if c.tagsClient == nil {
		return nil, fmt.Errorf("tags require %s or higher", tagsMinVersion)
	}
	return c.tagsClient, nil
}

// Config holds the provider configuration, and delivers a populated
// VSphereClient based off the contained settings.
type Config struct {
	InsecureFlag    bool
	Debug           bool
	Persist         bool
	User            string
	Password        string
	VSphereServer   string
	DebugPath       string
	DebugPathRun    string
	VimSessionPath  string
	RestSessionPath string
	KeepAlive       int
}

// NewConfig returns a new Config from a supplied ResourceData.
func NewConfig(d *schema.ResourceData) (*Config, error) {
	// Handle backcompat support for vcenter_server; once that is removed,
	// vsphere_server can just become a Required field that is referenced inline
	// in Config below.
	server := d.Get("vsphere_server").(string)

	if server == "" {
		server = d.Get("vcenter_server").(string)
	}

	if server == "" {
		return nil, fmt.Errorf("one of vsphere_server or [deprecated] vcenter_server must be provided")
	}

	c := &Config{
		User:            d.Get("user").(string),
		Password:        d.Get("password").(string),
		InsecureFlag:    d.Get("allow_unverified_ssl").(bool),
		VSphereServer:   server,
		Debug:           d.Get("client_debug").(bool),
		DebugPathRun:    d.Get("client_debug_path_run").(string),
		DebugPath:       d.Get("client_debug_path").(string),
		Persist:         d.Get("persist_session").(bool),
		VimSessionPath:  d.Get("vim_session_path").(string),
		RestSessionPath: d.Get("rest_session_path").(string),
		KeepAlive:       d.Get("vim_keep_alive").(int),
	}

	return c, nil
}

// vimURL returns a URL to pass to the VIM SOAP client.
func (c *Config) vimURL() (*url.URL, error) {
	u, err := url.Parse("https://" + c.VSphereServer + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("Error parse url: %s", err)
	}

	u.User = url.UserPassword(c.User, c.Password)

	return u, nil
}

// Client returns a new client for accessing VMWare vSphere.
func (c *Config) Client() (*VSphereClient, error) {
	client := new(VSphereClient)

	u, err := c.vimURL()
	if err != nil {
		return nil, fmt.Errorf("Error generating SOAP endpoint url: %s", err)
	}

	err = c.EnableDebug()
	if err != nil {
		return nil, fmt.Errorf("Error setting up client debug: %s", err)
	}

	// Set up the VIM/govmomi client connection, or load a previous session
	client.vimClient, err = c.SavedVimSessionOrNew(u)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] VMWare vSphere Client configured for URL: %s", c.VSphereServer)

	if isEligibleTagEndpoint(client.vimClient) {
		// Connect to the CIS REST endpoint for tagging, or load a previous session
		client.tagsClient, err = c.SavedRestSessionOrNew(u)
		if err != nil {
			return nil, err
		}
		log.Println("[DEBUG] CIS REST client configuration successful")
	} else {
		// Just print a log message so that we know that tags are not available on
		// this connection.
		log.Printf("[DEBUG] Connected endpoint does not support tags (%s)", viapi.ParseVersionFromClient(client.vimClient))
	}

	// Done, save sessions if we need to and return
	if err := c.SaveVimClient(client.vimClient); err != nil {
		return nil, fmt.Errorf("error persisting SOAP session to disk: %s", err)
	}
	if err := c.SaveRestClient(client.tagsClient); err != nil {
		return nil, fmt.Errorf("error persisting REST session to disk: %s", err)
	}

	return client, nil
}

// EnableDebug turns on govmomi API operation logging, if appropriate settings
// are set on the provider.
func (c *Config) EnableDebug() error {
	if !c.Debug {
		return nil
	}

	// Base path for storing debug logs.
	r := c.DebugPath
	if r == "" {
		r = filepath.Join(os.Getenv("HOME"), ".govmomi")
	}
	r = filepath.Join(r, "debug")

	// Path for this particular run.
	run := c.DebugPathRun
	if run == "" {
		now := time.Now().Format("2006-01-02T15-04-05.999999999")
		r = filepath.Join(r, now)
	} else {
		// reuse the same path
		r = filepath.Join(r, run)
		_ = os.RemoveAll(r)
	}

	err := os.MkdirAll(r, 0700)
	if err != nil {
		log.Printf("[ERROR] Client debug setup failed: %v", err)
		return err
	}

	p := debug.FileProvider{
		Path: r,
	}

	debug.SetProvider(&p)
	return nil
}

func (c *Config) vimURLWithoutPassword() (*url.URL, error) {
	u, err := c.vimURL()
	if err != nil {
		return nil, err
	}
	withoutCredentials := u
	withoutCredentials.User = url.User(u.User.Username())
	return withoutCredentials, nil
}

// sessionFile is a helper that generates a unique hash of the client's URL
// to use as the session file name.
//
// This is the same logic used as part of govmomi and is designed to be
// consistent so that sessions can be shared if possible between both tools.
func (c *Config) sessionFile() (string, error) {
	u, err := c.vimURLWithoutPassword()
	if err != nil {
		return "", err
	}

	// Key session file off of full URI and insecure setting.
	// Hash key to get a predictable, canonical format.
	key := fmt.Sprintf("%s#insecure=%t", u.String(), c.InsecureFlag)
	name := fmt.Sprintf("%040x", sha1.Sum([]byte(key)))
	return name, nil
}

// vimSessionFile is takes the session file name generated by sessionFile and
// then prefixes the SOAP client session path to it.
func (c *Config) vimSessionFile() (string, error) {
	p, err := c.sessionFile()
	if err != nil {
		return "", err
	}
	return filepath.Join(c.VimSessionPath, p), nil
}

// restSessionFile is takes the session file name generated by sessionFile and
// then prefixes the REST client session path to it.
func (c *Config) restSessionFile() (string, error) {
	p, err := c.sessionFile()
	if err != nil {
		return "", err
	}
	return filepath.Join(c.RestSessionPath, p), nil
}

// SaveVimClient saves a client to the supplied path. This facilitates re-use of
// the session at a later date.
//
// Note the logic in this function has been largely adapted from govc and is
// designed to be compatible with it.
func (c *Config) SaveVimClient(client *govmomi.Client) error {
	if !c.Persist {
		return nil
	}

	p, err := c.vimSessionFile()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Will persist SOAP client session data to %q", p)
	err = os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("[DEBUG] Error closing SOAP client session file %q: %s", p, err)
		}
	}()

	err = json.NewEncoder(f).Encode(client.Client)
	if err != nil {
		return err
	}

	return nil
}

// SaveRestClient saves the REST client session ID to the supplied path. This
// facilitates re-use of the session at a later date.
func (c *Config) SaveRestClient(client *tags.RestClient) error {
	if !c.Persist {
		return nil
	}

	p, err := c.restSessionFile()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Will persist REST client session data to %q", p)
	err = os.MkdirAll(filepath.Dir(p), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, []byte(client.SessionID()), 0600)
	if err != nil {
		return err
	}

	return nil
}

// restoreVimClient loads the saved session from disk. Note that this is a helper
// function to LoadVimClient and should not be called directly.
func (c *Config) restoreVimClient(client *vim25.Client) (bool, error) {
	if !c.Persist {
		return false, nil
	}

	p, err := c.vimSessionFile()
	if err != nil {
		return false, err
	}
	log.Printf("[DEBUG] Attempting to locate SOAP client session data in %q", p)
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] SOAP client session data not found in %q", p)
			return false, nil
		}

		return false, err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Printf("[DEBUG] Error closing SOAP client session file %q: %s", p, err)
		}
	}()

	dec := json.NewDecoder(f)
	err = dec.Decode(client)
	if err != nil {
		return false, err
	}

	return true, nil
}

// readRestSessionID reads a saved REST session ID and returns it. An empty
// string is returned if session does not exist.
func (c *Config) readRestSessionID() (string, error) {
	if !c.Persist {
		return "", nil
	}

	p, err := c.restSessionFile()
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Attempting to locate REST client session data in %q", p)
	id, err := ioutil.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] REST client session data not found in %q", p)
			return "", nil
		}

		return "", err
	}

	return string(id), nil
}

// LoadVimClient loads a saved vSphere SOAP API session from disk, previously
// saved by SaveVimClient, checking it for validity before returning it. A nil
// client means that the session is no longer valid and should be created from
// scratch.
//
// Note the logic in this function has been largely adapted from govc and is
// designed to be compatible with it - if a session has already been saved with
// govc, Terraform will attempt to use that session first.
func (c *Config) LoadVimClient() (*govmomi.Client, error) {
	client := new(vim25.Client)
	ok, err := c.restoreVimClient(client)
	if err != nil {
		return nil, err
	}

	if !ok || !client.Valid() {
		log.Println("[DEBUG] Cached SOAP client session data not valid or persistence not enabled, new session necessary")
		return nil, nil
	}

	m := session.NewManager(client)
	u, err := m.UserSession(context.TODO())
	if err != nil {
		if soap.IsSoapFault(err) {
			fault := soap.ToSoapFault(err).VimFault()
			// If the PropertyCollector is not found, the saved session for this URL is not valid
			if _, ok := fault.(types.ManagedObjectNotFound); ok {
				log.Println("[DEBUG] Cached SOAP client session missing property collector, new session necessary")
				return nil, nil
			}
		}

		return nil, err
	}

	// If the session is nil, the client is not authenticated
	if u == nil {
		log.Println("[DEBUG] Unauthenticated session, new session necessary")
		return nil, nil
	}

	log.Println("[DEBUG] Cached SOAP client session loaded successfully")
	return &govmomi.Client{
		Client:         client,
		SessionManager: m,
	}, nil
}

// LoadRestClient loads a saved vSphere REST API session from disk, previously
// saved by SaveRestClient, checking it for validity before returning it. If
// it's not valid, false is returned as the third return value, but the client
// can still be technically used for logging in by calling Login on the client.
func (c *Config) LoadRestClient(ctx context.Context, u *url.URL) (*tags.RestClient, bool, error) {
	id, err := c.readRestSessionID()
	if err != nil {
		return nil, false, err
	}

	client := tags.NewClientWithSessionID(u, c.InsecureFlag, "", id)

	if id == "" {
		log.Println("[DEBUG] No cached REST session data found or persistence not enabled, new session necessary")
		return client, false, nil
	}

	if !client.Valid(ctx) {
		log.Println("[DEBUG] Cached REST client session data not valid, new session necessary")
		return client, false, nil
	}

	log.Println("[DEBUG] Cached REST client session loaded successfully")
	return client, true, nil
}

// SavedVimSessionOrNew either loads a saved SOAP session from disk, or creates
// a new one.
func (c *Config) SavedVimSessionOrNew(u *url.URL) (*govmomi.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()

	client, err := c.LoadVimClient()
	if err != nil {
		return nil, fmt.Errorf("error trying to load vSphere SOAP session from disk: %s", err)
	}
	if client == nil {
		log.Printf("[DEBUG] Creating new SOAP API session on endpoint %s", c.VSphereServer)
		client, err = newClientWithKeepAlive(ctx, u, c.InsecureFlag, c.KeepAlive)
		if err != nil {
			return nil, fmt.Errorf("error setting up new vSphere SOAP client: %s", err)
		}
		log.Println("[DEBUG] SOAP API session creation successful")
	}
	return client, nil
}

func newClientWithKeepAlive(ctx context.Context, u *url.URL, insecure bool, keepAlive int) (*govmomi.Client, error) {
	soapClient := soap.NewClient(u, insecure)
	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
	}

	c := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	k := session.KeepAlive(c.Client.RoundTripper, time.Duration(keepAlive)*time.Minute)
	c.Client.RoundTripper = k

	// Only login if the URL contains user information.
	if u.User != nil {
		err = c.Login(ctx, u.User)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SavedRestSessionOrNew either loads a saved REST session from disk, or creates
// a new one.
func (c *Config) SavedRestSessionOrNew(u *url.URL) (*tags.RestClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()

	client, valid, err := c.LoadRestClient(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error trying to load vSphere REST session from disk: %s", err)
	}
	if !valid {
		log.Printf("[DEBUG] Creating new CIS REST API session on endpoint %s", c.VSphereServer)
		if err := client.Login(ctx); err != nil {
			return nil, fmt.Errorf("Error connecting to CIS REST endpoint: %s", err)
		}
		log.Println("[DEBUG] CIS REST API session creation successful")
	}
	return client, nil
}
