package vsphereprivate

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vapi/rest"
)

// VSphereClient - The VIM/govmomi client.
type VSphereClient struct {
	// vim client
	vimClient *govmomi.Client

	// rest client for tags
	restClient *rest.Client
}

// ConfigWrapper - wrapping the terraform-provider-vsphere Config struct
type ConfigWrapper struct {
	config *vsphere.Config
}

// NewConfig function
func NewConfig(d *schema.ResourceData) (*ConfigWrapper, error) {
	config, err := vsphere.NewConfig(d)
	if err != nil {
		return nil, err
	}
	return &ConfigWrapper{config}, nil
}

// vimURL returns a URL to pass to the VIM SOAP client.
func (cw *ConfigWrapper) vimURL() (*url.URL, error) {
	u, err := url.Parse("https://" + cw.config.VSphereServer + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("error parse url: %s", err)
	}

	u.User = url.UserPassword(cw.config.User, cw.config.Password)

	return u, nil
}

// restURL returns a URL to pass to the REST client.
func (cw *ConfigWrapper) restURL() (*cache.Session, error) {
	u, err := url.Parse("https://" + cw.config.VSphereServer)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	u.User = url.UserPassword(cw.config.User, cw.config.Password)
	s := &cache.Session{
		URL:      u,
		Insecure: cw.config.InsecureFlag,
	}
	return s, err
}

// Client returns a new client for accessing VMWare vSphere.
func (cw *ConfigWrapper) Client() (*VSphereClient, error) {
	client := new(VSphereClient)

	u, err := cw.vimURL()
	if err != nil {
		return nil, fmt.Errorf("error generating SOAP endpoint url: %s", err)
	}

	err = cw.config.EnableDebug()
	if err != nil {
		return nil, fmt.Errorf("error setting up client debug: %s", err)
	}

	// Set up the VIM/govmomi client connection, or load a previous session
	client.vimClient, err = cw.config.SavedVimSessionOrNew(u)
	if err != nil {
		return nil, err
	}

	s := new(cache.Session)
	s, err = cw.restURL()
	if err != nil {
		return nil, err
	}
	client.restClient, err = cw.config.SavedRestSessionOrNew(s)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] VMWare vSphere Client configured for URL: %s", cw.config.VSphereServer)

	return client, nil
}
