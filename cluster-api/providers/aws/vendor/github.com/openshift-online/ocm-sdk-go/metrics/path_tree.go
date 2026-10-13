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

// This file contains the type that describes trees of URL paths used to translate request paths
// into labes suitalbe for use as Prometheus labels.

package metrics

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// pathTree defines a tree of URL paths that will be used to transform request paths into labels
// suitable for use in Prometheus metrics. For example, a server that has these URL paths:
//
//	/api
//	/api/clusters_mgmt
//	/api/clusters_mgmt/v1
//	/api/clusters_mgmt/v1/clusters
//	/api/clusters_mgmt/v1/clusters/{cluster_id}
//	/api/clusters_mgmt/v1/clusters/{cluster_id}/groups
//	/api/clusters_mgmt/v1/clusters/{cluster_id}/groups/{group_id}
//
// Will be described with a tree like this:
//
//	var pathRoot = pathTree{
//		"api": {
//			"clusters_mgmt": {
//				"v1": {
//					"clusters": {
//						"-": {
//							"groups": {
//								"-": nil,
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//
// Path variables are represented with a dash.
type pathTree map[string]pathTree

// copy creates a deep copy of this tree.
func (t pathTree) copy() pathTree {
	if t == nil {
		return nil
	}
	tree := pathTree{}
	for label, child := range t {
		tree[label] = child.copy()
	}
	return tree
}

// add adds the given branch to this tree.
func (t pathTree) add(path string) {
	path = t.clean(path)
	if len(path) == 0 {
		return
	}
	segments := strings.Split(path, "/")
	t.addSegments(segments)
}

func (t pathTree) addSegments(segments []string) {
	if len(segments) == 0 {
		return
	}
	head := segments[0]
	tail := segments[1:]
	next := t[head]
	if next == nil {
		if len(tail) > 0 {
			next = pathTree{}
		}
		t[head] = next
	}
	next.addSegments(tail)
}

func (t pathTree) clean(path string) string {
	for len(path) > 0 && strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	for len(path) > 0 && strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	return path
}

// pathRoot is the root of the URL path tree.
var pathRoot pathTree

func init() {
	err := jsoniter.Unmarshal([]byte(pathTreeData), &pathRoot)
	if err != nil {
		panic(err)
	}
}
