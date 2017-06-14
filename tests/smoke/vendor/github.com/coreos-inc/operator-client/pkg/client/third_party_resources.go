package client

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/meta"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/wait"
)

// GetThirdPartyResource returns the third party resource as *runtime.Unstructured by the given third party resource name.
func (c *Client) GetThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string) (*runtime.Unstructured, error) {
	glog.V(4).Infof("[GET TPR]: %s:%s", namespace, resourceName)
	var object runtime.Unstructured

	b, err := c.GetThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &object); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TPR: %v", err)
	}
	return &object, nil
}

// GetThirdPartyResourceRaw returns the third party resource's raw body data by the given third party resource name.
func (c *Client) GetThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string) ([]byte, error) {
	glog.V(4).Infof("[GET TPR RAW]: %s:%s", namespace, resourceName)
	httpRestClient := c.ExtensionsV1beta1().RESTClient()
	uri := constructThirdPartyResourceURI(apiGroup, version, namespace, resourceKind, resourceName)
	glog.V(4).Infof("[GET]: %s", uri)

	return httpRestClient.Get().RequestURI(uri).DoRaw()
}

// CreateThirdPartyResource creates the third party resource.
func (c *Client) CreateThirdPartyResource(item *runtime.Unstructured) error {
	glog.V(4).Infof("[CREATE TPR]: %s:%s", item.GetNamespace(), item.GetName())
	kind := item.GetKind()
	namespace := item.GetNamespace()
	apiVersion := item.GetAPIVersion()
	apiGroup, version, err := parseAPIVersion(apiVersion)
	if err != nil {
		return err
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.CreateThirdPartyResourceRaw(apiGroup, version, namespace, kind, data)
}

// CreateThirdPartyResourceRaw creates the raw bytes of the third party resource.
func (c *Client) CreateThirdPartyResourceRaw(apiGroup, version, namespace, kind string, data []byte) error {
	glog.V(4).Infof("[CREATE TPR RAW]: %s:%s", namespace, kind)
	var statusCode int

	httpRestClient := c.ExtensionsV1beta1().RESTClient()
	uri := constructThirdPartyResourceKindURI(apiGroup, version, namespace, kind)
	glog.V(4).Infof("[POST]: %s", uri)
	result := httpRestClient.Post().RequestURI(uri).Body(data).Do()

	if result.Error() != nil {
		return result.Error()
	}

	result.StatusCode(&statusCode)
	glog.V(4).Infof("Written %s, status: %d", uri, statusCode)

	if statusCode != 201 {
		return fmt.Errorf("unexpected status code %d, expecting 201", statusCode)
	}
	return nil
}

// CreateThirdPartyResourceRawIfNotFound creates the raw bytes of the third party resource if it doesn't exist.
// It also returns a boolean to indicate whether a new TPR is created.
func (c *Client) CreateThirdPartyResourceRawIfNotFound(apiGroup, version, namespace, kind, name string, data []byte) (bool, error) {
	glog.V(4).Infof("[CREATE TPR RAW if not found]: %s:%s", namespace, name)
	_, err := c.GetThirdPartyResource(apiGroup, version, namespace, kind, name)
	if err == nil {
		return false, nil
	}
	if !errors.IsNotFound(err) {
		return false, err
	}
	err = c.CreateThirdPartyResourceRaw(apiGroup, version, namespace, kind, data)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateThirdPartyResource updates the third party resource.
// To do an atomic update, use AtomicModifyThirdPartyResource().
func (c *Client) UpdateThirdPartyResource(item *runtime.Unstructured) error {
	glog.V(4).Infof("[UPDATE TPR]: %s:%s", item.GetNamespace(), item.GetName())
	kind := item.GetKind()
	name := item.GetName()
	namespace := item.GetNamespace()
	apiVersion := item.GetAPIVersion()
	apiGroup, version, err := parseAPIVersion(apiVersion)
	if err != nil {
		return err
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.UpdateThirdPartyResourceRaw(apiGroup, version, namespace, kind, name, data)
}

// UpdateThirdPartyResourceRaw updates the thirdparty resource with the raw data.
func (c *Client) UpdateThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string, data []byte) error {
	glog.V(4).Infof("[UPDATE TPR RAW]: %s:%s", namespace, resourceName)
	var statusCode int

	httpRestClient := c.ExtensionsV1beta1().RESTClient()
	uri := constructThirdPartyResourceURI(apiGroup, version, namespace, resourceKind, resourceName)
	glog.V(4).Infof("[PUT]: %s", uri)
	result := httpRestClient.Put().RequestURI(uri).Body(data).Do()

	if result.Error() != nil {
		return result.Error()
	}

	result.StatusCode(&statusCode)
	glog.V(4).Infof("Updated %s, status: %d", uri, statusCode)

	if statusCode != 200 {
		return fmt.Errorf("unexpected status code %d, expecting 200", statusCode)
	}
	return nil
}

// CreateOrUpdateThirdpartyResourceRaw creates the TPR if it doesn't exist.
// If the TPR exists, it updates the existing one.
func (c *Client) CreateOrUpdateThirdpartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName string, data []byte) error {
	glog.V(4).Infof("[CREATE OR UPDATE UPDATE TPR RAW]: %s:%s", namespace, resourceName)
	_, err := c.GetThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return c.CreateThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, data)
	}
	return c.UpdateThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName, data)
}

// DeleteThirdPartyResource deletes the third party resource with the given name.
func (c *Client) DeleteThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string) error {
	glog.V(4).Infof("[DELETE TPR]: %s:%s", namespace, resourceName)
	httpRestClient := c.ExtensionsV1beta1().RESTClient()
	uri := constructThirdPartyResourceURI(apiGroup, version, namespace, resourceKind, resourceName)

	glog.V(4).Infof("[DELETE]: %s", uri)
	_, err := httpRestClient.Delete().RequestURI(uri).DoRaw()
	return err
}

// TPRModifier takes the TPR object, modifies it and returns the
// expecting result.
type TPRModifier func(*runtime.Unstructured, interface{}) error

// AtomicModifyThirdPartyResource gets the TPR, modifies it and writes it back.
// If it's modified by other writers, we will retry until it succeeds.
func (c *Client) AtomicModifyThirdPartyResource(apiGroup, version, namespace, resourceKind, resourceName string, f TPRModifier, data interface{}) error {
	glog.V(4).Infof("[ATOMIC MODIFY TPR]: %s:%s", namespace, resourceName)
	return wait.PollInfinite(time.Second, func() (bool, error) {
		var tpr runtime.Unstructured
		b, err := c.GetThirdPartyResourceRaw(apiGroup, version, namespace, resourceKind, resourceName)
		if err != nil {
			glog.Errorf("Failed to get TPR %q, kind:%q: %v", resourceName, resourceKind, err)
			return false, err
		}

		if err := json.Unmarshal(b, &tpr); err != nil {
			glog.Errorf("Failed to unmarshal TPR %q, kind:%q: %v", resourceName, resourceKind, err)
			return false, err
		}

		if err := f(&tpr, data); err != nil {
			glog.Errorf("Failed to modify the TPR %q, kind:%q: %v", resourceName, resourceKind, err)
			return false, err
		}

		if err := c.UpdateThirdPartyResource(&tpr); err != nil {
			if errors.IsConflict(err) {
				glog.Errorf("Failed to update TPR %q, kind:%q: %v, will retry", resourceName, resourceKind, err)
				return false, nil
			}
			glog.Errorf("Failed to update TPR %q, kind:%q: %v", resourceName, resourceKind, err)
			return false, err
		}

		return true, nil
	})
}

// constructThirdPartyResourceURI returns the URI for the thirdparty resource.
//
// Example of apiGroup: "coreos.com"
// Example of version: "v1"
// Example of namespace: "default"
// Example of resourceKind: "ChannelOperatorConfig"
// Example of resourceName: "test-config"
func constructThirdPartyResourceURI(apiGroup, version, namespace, resourceKind, resourceName string) string {
	if namespace == "" {
		namespace = v1.NamespaceDefault
	}
	plural, _ := meta.KindToResource(unversioned.GroupVersionKind{
		Group:   apiGroup,
		Version: version,
		Kind:    resourceKind,
	})
	return fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s/%s",
		strings.ToLower(apiGroup),
		strings.ToLower(version),
		strings.ToLower(namespace),
		strings.ToLower(plural.Resource),
		strings.ToLower(resourceName))
}

// constructThirdPartyResourceKindURI returns the URI for the thirdparty resource kind.
//
// Example of apiGroup: "coreos.com"
// Example of version: "v1"
// Example of namespace: "default"
// Example of resourceKind: "ChannelOperatorConfig"
func constructThirdPartyResourceKindURI(apiGroup, version, namespace, resourceKind string) string {
	if namespace == "" {
		namespace = v1.NamespaceDefault
	}
	plural, _ := meta.KindToResource(unversioned.GroupVersionKind{
		Group:   apiGroup,
		Version: version,
		Kind:    resourceKind,
	})
	return fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s",
		strings.ToLower(apiGroup),
		strings.ToLower(version),
		strings.ToLower(namespace),
		strings.ToLower(plural.Resource))
}

// parseAPIVersion splits "coreos.com/v1" into
// "coreos.com" and "v1".
func parseAPIVersion(apiVersion string) (apiGroup, version string, err error) {
	parts := strings.Split(apiVersion, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid format of api version %q, expecting APIGroup/Version", apiVersion)
	}
	return path.Join(parts[:len(parts)-1]...), parts[len(parts)-1], nil
}
