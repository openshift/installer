package client

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

// GetService returns the Service object for the given namespace and name.
func (c *Client) GetService(namespace, name string) (*v1.Service, error) {
	glog.V(4).Infof("[GET Service]: %s:%s", namespace, name)
	return c.Services(namespace).Get(name)
}

// UpdateService updates the Service object.
func (c *Client) UpdateService(s *v1.Service) (*v1.Service, error) {
	glog.V(4).Infof("[UPDATE Service]: %s:%s", s.GetNamespace(), s.GetName())
	// If the cluster IP is not specified, then an update/replace operation will
	// fail. In that case, try to patch the service. Unfortunately, patches only
	// work with JSON, so first convert the service to JSON and then send it off.
	if s.Spec.ClusterIP == "" {
		data, err := json.Marshal(s)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal service: %v", err)
		}
		return c.Core().Services(s.Namespace).Patch(s.Name, api.StrategicMergePatchType, data)
	}
	return c.Services(s.GetNamespace()).Update(s)
}
