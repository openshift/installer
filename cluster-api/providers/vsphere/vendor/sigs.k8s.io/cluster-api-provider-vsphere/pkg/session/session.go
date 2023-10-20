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
	"github.com/go-logr/logr"
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

type Feature struct {
	EnableKeepAlive   bool
	KeepAliveDuration time.Duration
}

func DefaultFeature() Feature {
	return Feature{
		EnableKeepAlive: constants.DefaultEnableKeepAlive,
	}
}

type Params struct {
	server     string
	datacenter string
	userinfo   *url.Userinfo
	thumbprint string
	feature    Feature
}

func NewParams() *Params {
	return &Params{
		feature: DefaultFeature(),
	}
}

func (p *Params) WithServer(server string) *Params {
	p.server = server
	return p
}

func (p *Params) WithDatacenter(datacenter string) *Params {
	p.datacenter = datacenter
	return p
}

func (p *Params) WithUserInfo(username, password string) *Params {
	p.userinfo = url.UserPassword(username, password)
	return p
}

func (p *Params) WithThumbprint(thumbprint string) *Params {
	p.thumbprint = thumbprint
	return p
}

func (p *Params) WithFeatures(feature Feature) *Params {
	p.feature = feature
	return p
}

// GetOrCreate gets a cached session or creates a new one if one does not
// already exist.
func GetOrCreate(ctx context.Context, params *Params) (*Session, error) {
	logger := ctrl.LoggerFrom(ctx).WithName("session").WithValues(
		"server", params.server,
		"datacenter", params.datacenter,
		"username", params.userinfo.Username())

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

		vimSessionActive, err := s.SessionManager.SessionIsActive(ctx)
		if err != nil {
			logger.Error(err, "unable to check if vim session is active")
		}

		tagManagerSession, err := s.TagManager.Session(ctx)
		if err != nil {
			logger.Error(err, "unable to check if rest session is active")
		}

		if vimSessionActive && tagManagerSession != nil {
			logger.V(2).Info("found active cached vSphere client session")
			return s, nil
		}

		logger.V(2).Info("logout the session because it is inactive")
		if err := s.Client.Logout(ctx); err != nil {
			logger.Error(err, "unable to logout session")
		} else {
			logger.Info("logout session succeed")
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
		return nil, errors.Wrapf(err, "error parsing vSphere URL %q", params.server)
	}
	if soapURL == nil {
		return nil, errors.Errorf("error parsing vSphere URL %q", params.server)
	}

	soapURL.User = params.userinfo
	client, err := newClient(ctx, logger, sessionKey, soapURL, params.thumbprint, params.feature)
	if err != nil {
		return nil, err
	}

	session := Session{Client: client}
	session.UserAgent = infrav1.GroupVersion.String()

	// Assign the finder to the session.
	session.Finder = find.NewFinder(session.Client.Client, false)
	// Assign tag manager to the session.
	manager, err := newManager(ctx, logger, sessionKey, client.Client, soapURL.User, params.feature)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create tags manager")
	}
	session.TagManager = manager

	// Assign the datacenter if one was specified.
	if params.datacenter != "" {
		dc, err := session.Finder.Datacenter(ctx, params.datacenter)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find datacenter %q", params.datacenter)
		}
		session.datacenter = dc
		session.Finder.SetDatacenter(dc)
	}
	// Cache the session.
	sessionCache.Store(sessionKey, &session)

	logger.V(2).Info("cached vSphere client session", "server", params.server, "datacenter", params.datacenter)

	return &session, nil
}

func newClient(ctx context.Context, logger logr.Logger, sessionKey string, url *url.URL, thumbprint string, feature Feature) (*govmomi.Client, error) {
	insecure := thumbprint == ""
	soapClient := soap.NewClient(url, insecure)
	if !insecure {
		soapClient.SetThumbprint(url.Host, thumbprint)
	}

	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
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
				logger.Error(err, "failed to keep alive govmomi client")
				logger.Info("clearing the session")
				sessionCache.Delete(sessionKey)
			}
			return err
		})
	}

	if err := c.Login(ctx, url.User); err != nil {
		return nil, err
	}

	return c, nil
}

// newManager creates a Manager that encompasses the REST Client for the VSphere tagging API.
func newManager(ctx context.Context, logger logr.Logger, sessionKey string, client *vim25.Client, user *url.Userinfo, feature Feature) (*tags.Manager, error) {
	rc := rest.NewClient(client)
	if feature.EnableKeepAlive {
		rc.Transport = keepalive.NewHandlerREST(rc, feature.KeepAliveDuration, func() error {
			s, err := rc.Session(ctx)
			if err != nil {
				return err
			}
			if s != nil {
				return nil
			}

			logger.Info("rest client session expired, clearing session")
			sessionCache.Delete(sessionKey)
			return errors.New("rest client session expired")
		})
	}
	if err := rc.Login(ctx, user); err != nil {
		return nil, err
	}
	return tags.NewManager(rc), nil
}

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
