// Package routeapihelpers contains utilities for handling OpenShift route objects.
package routeapihelpers

import (
	"fmt"
	"net/url"

	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
)

// IngressURI calculates an admitted ingress URI.
// If 'host' is nonempty, only the ingress for that host is considered.
// If 'host' is empty, the first admitted ingress is used.
func IngressURI(route *routev1.Route, host string) (*url.URL, *routev1.RouteIngress, error) {
	scheme := "http"
	if route.Spec.TLS != nil {
		scheme = "https"
	}

	for _, ingress := range route.Status.Ingress {
		if host == "" || host == ingress.Host {
			uri := &url.URL{
				Scheme: scheme,
				Host:   ingress.Host,
			}

			for _, condition := range ingress.Conditions {
				if condition.Type == routev1.RouteAdmitted && condition.Status == corev1.ConditionTrue {
					return uri, &ingress, nil
				}
			}

			if host == ingress.Host {
				return uri, &ingress, fmt.Errorf("ingress for host %s in route %s in namespace %s is not admitted", host, route.ObjectMeta.Name, route.ObjectMeta.Namespace)
			}
		}
	}

	if host == "" {
		return nil, nil, fmt.Errorf("no admitted ingress for route %s in namespace %s", route.ObjectMeta.Name, route.ObjectMeta.Namespace)
	}
	return nil, nil, fmt.Errorf("no ingress for host %s in route %s in namespace %s", host, route.ObjectMeta.Name, route.ObjectMeta.Namespace)
}
