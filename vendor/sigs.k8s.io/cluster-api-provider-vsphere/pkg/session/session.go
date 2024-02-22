/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package session contains tools to create and retrieve a VCenter session.
package session

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/netip"
	"net/url"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/session/keepalive"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/constants"
)

var (
	// global Session map against sessionKeys in map[sessionKey]Session.
	sessionCache sync.Map

	// mutex to control access to the GetOrCreate function to avoid duplicate
	// session creations on startup.
	sessionMU sync.Mutex
)

// Session is a vSphere session with a configured Finder.
type Session struct {
	*govmomi.Client
	Finder     *find.Finder
	datacenter *object.Datacenter
	TagManager *tags.Manager
}

// Feature is a set of Features of the session.
type Feature struct {
	EnableKeepAlive   bool
	KeepAliveDuration time.Duration
}

// DefaultFeature sets the default values for features.
func DefaultFeature() Feature {
	return Feature{
		EnableKeepAlive: constants.DefaultEnableKeepAlive,
	}
}

// Params are the parameters of a VCenter session.
type Params struct {
	server     string
	datacenter string
	userinfo   *url.Userinfo
	thumbprint string
	feature    Feature
}

// NewParams returns an empty set of parameters with default features.
func NewParams() *Params {
	return &Params{
		feature: DefaultFeature(),
	}
}

// WithServer adds a server to parameters.
func (p *Params) WithServer(server string) *Params {
	p.server = server
	return p
}

// WithDatacenter adds a datacenter to parameters.
func (p *Params) WithDatacenter(datacenter string) *Params {
	p.datacenter = datacenter
	return p
}

// WithUserInfo adds userinfo to parameters.
func (p *Params) WithUserInfo(username, password string) *Params {
	p.userinfo = url.UserPassword(username, password)
	return p
}

// WithThumbprint adds a thumbprint to parameters.
func (p *Params) WithThumbprint(thumbprint string) *Params {
	p.thumbprint = thumbprint
	return p
}

// WithFeatures adds features to parameters.
func (p *Params) WithFeatures(feature Feature) *Params {
	p.feature = feature
	return p
}

// GetOrCreate gets a cached session or creates a new one if one does not
// already exist.
func GetOrCreate(ctx context.Context, params *Params) (*Session, error) {
	log := ctrl.LoggerFrom(ctx).WithValues(
		"server", params.server,
		"datacenter", params.datacenter,
		"username", params.userinfo.Username())
	ctx = ctrl.LoggerInto(ctx, log)

	sessionMU.Lock()
	defer sessionMU.Unlock()

	userPassword, _ := params.userinfo.Password()
	h := sha256.New()
	h.Write([]byte(userPassword))
	hashedUserPassword := h.Sum(nil)
	sessionKey := fmt.Sprintf("%s#%s#%s#%x", params.server, params.datacenter, params.userinfo.Username(),
		hashedUserPassword)
	if cachedSession, ok := sessionCache.Load(sessionKey); ok {
		s := cachedSession.(*Session)

		// Retrieve the current session from Managed Object.
		// The userSession is active when the value is not nil.
		userSession, err := s.SessionManager.UserSession(ctx)
		if err != nil {
			log.Error(err, "Failed to check if vim session is active")
		}

		tagManagerSession, err := s.TagManager.Session(ctx)
		if err != nil {
			log.Error(err, "Failed to check if REST session is active")
		}

		if userSession != nil && tagManagerSession != nil {
			log.Info("Found active cached vSphere client session")
			return s, nil
		}

		log.Info("Logout the REST session because it is inactive")
		if err := s.TagManager.Logout(ctx); err != nil {
			log.Error(err, "Failed to logout REST session")
		} else {
			log.Info("Logout REST session succeed")
		}

		log.Info("Logout the session because it is inactive")
		if err := s.Client.Logout(ctx); err != nil {
			log.Error(err, "Failed to logout session")
		} else {
			log.Info("Logout session succeed")
		}
	}

	// soap.ParseURL expects a valid URL. In the case of a bare, unbracketed
	// IPv6 address (e.g fd00::1) ParseURL will fail. Surround unbracketed IPv6
	// addresses with brackets.
	urlSafeServer := params.server
	ip, err := netip.ParseAddr(urlSafeServer)
	if err == nil && ip.Is6() {
		urlSafeServer = fmt.Sprintf("[%s]", urlSafeServer)
	}

	soapURL, err := soap.ParseURL(urlSafeServer)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create vCenter session: error parsing vSphere URL %q", params.server)
	}
	if soapURL == nil {
		return nil, errors.Errorf("failed to create vCenter session: error parsing vSphere URL %q: URL is nil", params.server)
	}

	soapURL.User = params.userinfo
	client, err := newClient(ctx, sessionKey, soapURL, params.thumbprint, params.feature)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create vCenter session")
	}

	session := Session{Client: client}
	session.UserAgent = infrav1.GroupVersion.String()

	// Assign the finder to the session.
	session.Finder = find.NewFinder(session.Client.Client, false)
	// Assign tag manager to the session.
	manager, err := newManager(ctx, sessionKey, client.Client, soapURL.User, params.feature)
	if err != nil {
		log.Error(err, "Failed to create tags manager, will logout")
		// Logout of previously logged session to not leak
		if errLogout := client.Logout(ctx); errLogout != nil {
			log.Error(errLogout, "Failed to logout of leading client session")
		}
		return nil, errors.Wrap(err, "failed to create vCenter session: failed to create tags manager")
	}
	session.TagManager = manager

	// Assign the datacenter if one was specified.
	if params.datacenter != "" {
		dc, err := session.Finder.Datacenter(ctx, params.datacenter)
		if err != nil {
			log.Error(err, "Failed to get datacenter, will logout")
			// Logout of previously logged session to not leak
			if errLogout := manager.Logout(ctx); errLogout != nil {
				log.Error(errLogout, "Failed to logout of leading REST session")
			}
			if errLogout := client.Logout(ctx); errLogout != nil {
				log.Error(errLogout, "Failed to logout of leading client session")
			}
			return nil, errors.Wrapf(err, "failed to create vCenter session: failed to find datacenter %q", params.datacenter)
		}
		session.datacenter = dc
		session.Finder.SetDatacenter(dc)
	}
	// Cache the session.
	sessionCache.Store(sessionKey, &session)

	log.Info("Created and cached vSphere client session")

	return &session, nil
}

func newClient(ctx context.Context, sessionKey string, url *url.URL, thumbprint string, feature Feature) (*govmomi.Client, error) {
	log := ctrl.LoggerFrom(ctx)

	insecure := thumbprint == ""
	soapClient := soap.NewClient(url, insecure)
	if !insecure {
		soapClient.SetThumbprint(url.Host, thumbprint)
	}

	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create client")
	}
	vimClient.UserAgent = "k8s-capv-useragent"

	c := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	if feature.EnableKeepAlive {
		vimClient.RoundTripper = session.KeepAliveHandler(vimClient.RoundTripper, feature.KeepAliveDuration, func(tripper soap.RoundTripper) error {
			_, err := methods.GetCurrentTime(ctx, tripper)
			if err != nil {
				log.Error(err, "Failed to keep alive govmomi client, Clearing the session now")
				if errLogout := c.Logout(ctx); errLogout != nil {
					log.Error(err, "Failed to logout keepalive failed session")
				}
				sessionCache.Delete(sessionKey)
			}
			return err
		})
	}

	if err := c.Login(ctx, url.User); err != nil {
		return nil, errors.Wrapf(err, "failed to create client: failed to login")
	}

	return c, nil
}

// newManager creates a Manager that encompasses the REST Client for the VSphere tagging API.
func newManager(ctx context.Context, sessionKey string, client *vim25.Client, user *url.Userinfo, feature Feature) (*tags.Manager, error) {
	log := ctrl.LoggerFrom(ctx)

	rc := rest.NewClient(client)
	if feature.EnableKeepAlive {
		rc.Transport = keepalive.NewHandlerREST(rc, feature.KeepAliveDuration, func() error {
			s, err := rc.Session(ctx)
			if s != nil && err == nil {
				return nil
			}

			if err != nil {
				log.Error(err, "Failed to keep alive REST client")
			}

			log.Info("REST client session expired, clearing session")
			if errLogout := rc.Logout(ctx); errLogout != nil {
				log.Error(err, "Failed to logout keepalive failed REST session")
			}
			sessionCache.Delete(sessionKey)
			return errors.New("REST client session expired")
		})
	}
	if err := rc.Login(ctx, user); err != nil {
		return nil, errors.Wrapf(err, "failed to create tags manager: failed to login REST client")
	}
	return tags.NewManager(rc), nil
}

// GetVersion returns the VCenterVersion.
func (s *Session) GetVersion() (infrav1.VCenterVersion, error) {
	svcVersion := s.ServiceContent.About.Version
	version, err := semver.New(svcVersion)
	if err != nil {
		return "", err
	}

	switch version.Major {
	case 6, 7, 8:
		return infrav1.NewVCenterVersion(svcVersion), nil
	default:
		return "", unidentifiedVCenterVersion{version: svcVersion}
	}
}

// Clear is meant to destroy all the cached sessions.
func Clear() {
	sessionCache.Range(func(key, s any) bool {
		cachedSession := s.(*Session)
		_ = cachedSession.Logout(context.Background())
		return true
	})
}

// FindByBIOSUUID finds an object by its BIOS UUID.
//
// To avoid comments about this function's name, please see the Golang
// WIKI https://github.com/golang/go/wiki/CodeReviewComments#initialisms.
// This function is named in accordance with the example "XMLHTTP".
func (s *Session) FindByBIOSUUID(ctx context.Context, uuid string) (object.Reference, error) {
	return s.findByUUID(ctx, uuid, false)
}

// FindByInstanceUUID finds an object by its instance UUID.
func (s *Session) FindByInstanceUUID(ctx context.Context, uuid string) (object.Reference, error) {
	return s.findByUUID(ctx, uuid, true)
}

func (s *Session) findByUUID(ctx context.Context, uuid string, findByInstanceUUID bool) (object.Reference, error) {
	if s.Client == nil {
		return nil, errors.New("vSphere client is not initialized")
	}
	si := object.NewSearchIndex(s.Client.Client)
	ref, err := si.FindByUuid(ctx, s.datacenter, uuid, true, &findByInstanceUUID)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding object by uuid %q", uuid)
	}
	return ref, nil
}
