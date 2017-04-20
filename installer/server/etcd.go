package server

import (
	"fmt"
	"net/url"
	"strings"
)

// etcdNodeName returns the expected etcd name value.
func etcdNodeName(i int) string {
	return fmt.Sprintf("node%d", i)
}

// etcdInitialCluster returns a comma separated list of name=http://peer:2380
// values for an etcd cluster with the given nodes. This value is typically
// used for the etcd -initial-cluster flag or ETCD_INITIAL_CLUSTER variable.
func etcdInitialCluster(nodes []Node) string {
	peers := make([]string, len(nodes))
	for i, node := range nodes {
		peers[i] = fmt.Sprintf("%s=http://%s:2380", etcdNodeName(i), node.Name)
	}
	return strings.Join(peers, ",")
}

// etcdEndpoints returns the comma separated list of etcd client endpoints which
// can be polled to determine availability of the etcd cluster with the
// given nodes.
func etcdEndpoints(nodes []Node) string {
	peers := make([]string, len(nodes))
	for i, node := range nodes {
		peers[i] = fmt.Sprintf("%s:2379", node.Name)
	}
	return strings.Join(peers, ",")
}

// etcdURLs returns the list of etcd client URLs.
func etcdURLs(nodes []Node) ([]*url.URL, error) {
	urls := make([]*url.URL, len(nodes))
	for i, node := range nodes {
		raw := fmt.Sprintf("http://%s:2379", node.Name)
		u, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}
		urls[i] = u
	}
	return urls, nil
}
