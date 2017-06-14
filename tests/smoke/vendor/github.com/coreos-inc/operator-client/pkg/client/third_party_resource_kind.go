package client

import (
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api/errors"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/util/wait"
)

// GetThirdPartyResourceKind gets the thirdparty resource kind.
func (c *Client) GetThirdPartyResourceKind(name string) (*v1beta1extensions.ThirdPartyResource, error) {
	glog.V(4).Infof("[GET TPR Kind]: %s", name)
	return c.ExtensionsV1beta1().ThirdPartyResources().Get(name)
}

// CreateThirdPartyResourceKind creates the thirdparty resource kind.
func (c *Client) CreateThirdPartyResourceKind(tpr *v1beta1extensions.ThirdPartyResource) error {
	glog.V(4).Infof("[CREATE TPR Kind]: %s", tpr.Name)
	_, err := c.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
	return err
}

// DeleteThirdPartyResourceKind deletes the thirdparty resource kind.
func (c *Client) DeleteThirdPartyResourceKind(name string) error {
	glog.V(4).Infof("[DELETE TPR Kind]: %s", name)
	return c.ExtensionsV1beta1().ThirdPartyResources().Delete(name, nil)
}

// EnsureThirdPartyResourceKind will creates the TPR kind if it doesn't exist.
// On success, the TPR kind is guaranteed to be existed.
func (c *Client) EnsureThirdPartyResourceKind(tpr *v1beta1extensions.ThirdPartyResource) error {
	return wait.PollInfinite(time.Second, func() (bool, error) {
		_, err := c.GetThirdPartyResourceKind(tpr.Name)
		if err == nil {
			return true, nil
		}

		if !errors.IsNotFound(err) {
			return false, err
		}

		if err := c.CreateThirdPartyResourceKind(tpr); err != nil {
			glog.Errorf("Failed to create TPR Kind %q: %v", tpr.Name, err)
			return false, err
		}
		return false, nil
	})
}
