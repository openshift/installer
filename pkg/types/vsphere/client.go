package vsphere

import (
	"context"
	"net/url"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string) (*vim25.Client, *rest.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, err
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(ctx, u, false)

	if err != nil {
		return nil, nil, err
	}

	restClient := rest.NewClient(c.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		return nil, nil, err
	}

	return c.Client, restClient, nil
}
