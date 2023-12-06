package kubeconfig

import (
	"context"
	"fmt"
	"net"
)

type dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// CreateDialContext overrides the kubeconfig api server address with the provided ip address on the same port.
func CreateDialContext(d dialer, apiServerIPOverride string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		if network != "tcp" {
			return nil, fmt.Errorf("unimplemented network %q", network)
		}

		_, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		return d.DialContext(ctx, network, net.JoinHostPort(apiServerIPOverride, port))
	}
}
