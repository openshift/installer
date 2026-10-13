/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains the implementation of the server address parser.

package internal

import (
	"context"
	"fmt"
	neturl "net/url"
	"strings"
)

// ServerAddress contains a parsed URL and additional information extracted from int, like the
// network (tcp or unix) and the socket name (for Unix sockets).
type ServerAddress struct {
	// Text is the original text that was passed to the ParseServerAddress function to create
	// this server address.
	Text string

	// Network is the network that should be used to connect to the server. Possible values are
	// `tcp` and `unix`.
	Network string

	// Protocol is the application protocol used to connect to the server. Possible values are
	// `http`, `https` and `h2c`.
	Protocol string

	// Host is the name of the host used to connect to the server. This will be populated only
	// even when using Unix sockets, because clients will need it in order to populate the
	// `Host` header.
	Host string

	// Port is the port number used to connect to the server. This will only be populated when
	// using TCP. When using Unix sockets it will be zero.
	Port string

	// Socket is tha nem of the path of the Unix socket used to connect to the server.
	Socket string

	// URL is the regular URL calculated from this server address. The scheme will be `http` if
	// the protocol is `http` or `h2c` and will be `https` if the protocol is https.
	URL *neturl.URL
}

// ParseServerAddress parses the given text as a server address. Server addresses should be URLs
// with this format:
//
//	network+protocol://host:port/path?network=...&protocol=...&socket=...
//
// The `network` and `protocol` parts of the scheme are optional.
//
// Valid values for the `network` part of the scheme are `unix` and `tcp`. If not specified the
// default value is `tcp`.
//
// Valid values for the `protocol` part of the scheme are `http`, `https` and `h2c`. If not
// specified the default value is `http`.
//
// The `host` is mandatory even when using Unix sockets, because it is necessary to populate the
// `Host` header.
//
// The `port` part is optional. If not specified it will be 80 for HTTP and H2C and 443 for HTTPS.
//
// When using Unix sockets the `path` part will be used as the name of the Unix socket.
//
// The network protocol and Unix socket can alternatively be specified using the `network`,
// `protocol` and `socket` query parameters. This is useful specially for specifying the Unix
// sockets when the path of the URL has some other meaning. For example, in order to specify
// the OpenID token URL it is usually necessary to include a path, so to use a Unix socket it
// is necessary to put it in the `socket` parameter instead:
//
//	unix://my.sso.com/my/token/path?socket=/sockets/my.socket
//
// When the Unix socket is specified in the `socket` query parameter as in the above example
// the URL path will be ignored.
//
// Some examples of valid server addresses:
//
//   - http://my.server.com - HTTP on top of TCP.
//   - https://my.server.com - HTTPS on top of TCP.
//   - unix://my.server.com/sockets/my.socket - HTTP on top Unix socket.
//   - unix+https://my.server.com/sockets/my.socket - HTTPS on top of Unix socket.
//   - h2c+unix://my.server.com?socket=/sockets/my.socket - H2C on top of Unix.
func ParseServerAddress(ctx context.Context, text string) (result *ServerAddress, err error) {
	// Parse the URL:
	parsed, err := neturl.Parse(text)
	if err != nil {
		return
	}
	query := parsed.Query()

	// Extract the network and protocol from the scheme:
	networkFromScheme, protocolFromScheme, err := parseScheme(ctx, parsed.Scheme)
	if err != nil {
		return
	}

	// Check if the network is also specified with a query parameter. If it is it should not be
	// conflicting with the value specified in the scheme.
	var network string
	networkValues, ok := query["network"]
	if ok {
		if len(networkValues) != 1 {
			err = fmt.Errorf(
				"expected exactly one value for the 'network' query parameter "+
					"but found %d",
				len(networkValues),
			)
			return
		}
		networkFromQuery := strings.TrimSpace(strings.ToLower(networkValues[0]))
		err = checkNetwork(networkFromQuery)
		if err != nil {
			return
		}
		if networkFromScheme != "" && networkFromScheme != networkFromQuery {
			err = fmt.Errorf(
				"network '%s' from query parameter isn't compatible with "+
					"network '%s' from scheme",
				networkFromQuery, networkFromScheme,
			)
			return
		}
		network = networkFromQuery
	} else {
		network = networkFromScheme
	}

	// Check if the protocol is also specified with a query parameter. If it is it should not be
	// conflicting with the value specified in the scheme.
	var protocol string
	protocolValues, ok := query["protocol"]
	if ok {
		if len(protocolValues) != 1 {
			err = fmt.Errorf(
				"expected exactly one value for the 'protocol' query parameter "+
					"but found %d",
				len(protocolValues),
			)
			return
		}
		protocolFromQuery := strings.TrimSpace(strings.ToLower(protocolValues[0]))
		err = checkProtocol(protocolFromQuery)
		if err != nil {
			return
		}
		if protocolFromScheme != "" && protocolFromScheme != protocolFromQuery {
			err = fmt.Errorf(
				"protocol '%s' from query parameter isn't compatible with "+
					"protocol '%s' from scheme",
				protocolFromQuery, protocolFromScheme,
			)
			return
		}
		protocol = protocolFromQuery
	} else {
		protocol = protocolFromScheme
	}

	// Set default values for the network and protocol if needed:
	if network == "" {
		network = TCPNetwork
	}
	if protocol == "" {
		protocol = HTTPProtocol
	}

	// Get the host name. Note that the host name is mandatory even when using Unix sockets,
	// because it is used to populate the `Host` header.
	host := parsed.Hostname()
	if host == "" {
		err = fmt.Errorf("host name is mandatory, but it is empty")
		return
	}

	// Get the port number:
	port := parsed.Port()
	if port == "" {
		switch protocol {
		case HTTPProtocol, H2CProtocol:
			port = "80"
		case HTTPSProtocol:
			port = "443"
		}
	}

	// Get the socket from the `socket` query parameter or from the path:
	var socket string
	if network == UnixNetwork {
		socketValues, ok := query["socket"]
		if ok {
			if len(socketValues) != 1 {
				err = fmt.Errorf(
					"expected exactly one value for the 'socket' query "+
						"parameter but found %d",
					len(socketValues),
				)
				return
			}
			socket = socketValues[0]
		} else {
			socket = parsed.Path
		}
		if socket == "" {
			err = fmt.Errorf(
				"expected socket name in the 'socket' query parameter or in " +
					"the path but both are empty",
			)
			return
		}
	}

	// Calculate the URL:
	url := &neturl.URL{
		Host: host,
	}
	switch protocol {
	case HTTPProtocol, H2CProtocol:
		url.Scheme = "http"
		if port != "80" {
			url.Host = fmt.Sprintf("%s:%s", url.Host, port)
		}
	case HTTPSProtocol:
		url.Scheme = "https"
		if port != "443" {
			url.Host = fmt.Sprintf("%s:%s", url.Host, port)
		}
	}

	// Create and populate the result:
	result = &ServerAddress{
		Text:     text,
		Network:  network,
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Socket:   socket,
		URL:      url,
	}

	return
}

func parseScheme(ctx context.Context, scheme string) (network, protocol string,
	err error) {
	components := strings.Split(strings.ToLower(scheme), "+")
	if len(components) > 2 {
		err = fmt.Errorf(
			"scheme '%s' should have at most two components separated by '+', "+
				"but it has %d",
			scheme, len(components),
		)
		return
	}
	for _, component := range components {
		switch strings.TrimSpace(component) {
		case TCPNetwork, UnixNetwork:
			network = component
		case HTTPProtocol, HTTPSProtocol, H2CProtocol:
			protocol = component
		default:
			err = fmt.Errorf(
				"component '%s' of scheme '%s' doesn't correspond to any "+
					"supported network or protocol, supported networks "+
					"are 'tcp' and 'unix', supported protocols are 'http', "+
					"'https' and 'h2c'",
				component, scheme,
			)
			return
		}
	}
	return
}

func checkNetwork(value string) error {
	switch value {
	case UnixNetwork, TCPNetwork:
		return nil
	default:
		return fmt.Errorf(
			"network '%s' isn't valid, valid values are 'unix' and 'tcp'",
			value,
		)
	}
}

func checkProtocol(value string) error {
	switch value {
	case HTTPProtocol, HTTPSProtocol, H2CProtocol:
		return nil
	default:
		return fmt.Errorf(
			"protocol '%s' isn't valid, valid values are 'http', 'https' "+
				"and 'h2c'",
			value,
		)
	}
}

// Network names:
const (
	UnixNetwork = "unix"
	TCPNetwork  = "tcp"
)

// Protocol names:
const (
	HTTPProtocol  = "http"
	HTTPSProtocol = "https"
	H2CProtocol   = "h2c"
)
