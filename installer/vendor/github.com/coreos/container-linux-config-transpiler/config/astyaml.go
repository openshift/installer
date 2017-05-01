// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"io"

	yaml "github.com/ajeddeloh/yaml"
	"github.com/coreos/ignition/config/validate"
)

var (
	ErrNotDocumentNode = errors.New("Can only convert from document node")
)

type yamlNode struct {
	key yaml.Node
	yaml.Node
}

func fromYamlDocumentNode(n yaml.Node) (yamlNode, error) {
	if n.Kind != yaml.DocumentNode {
		return yamlNode{}, ErrNotDocumentNode
	}

	return yamlNode{
		key:  n,
		Node: *n.Children[0],
	}, nil
}

func (n yamlNode) ValueLineCol(source io.ReadSeeker) (int, int, string) {
	return n.Line, n.Column, ""
}

func (n yamlNode) KeyLineCol(source io.ReadSeeker) (int, int, string) {
	return n.key.Line, n.key.Column, ""
}

func (n yamlNode) LiteralValue() interface{} {
	return n.Value
}

func (n yamlNode) SliceChild(index int) (validate.AstNode, bool) {
	if n.Kind != yaml.SequenceNode {
		return nil, false
	}
	if index >= len(n.Children) {
		return nil, false
	}

	return yamlNode{
		key:  yaml.Node{},
		Node: *n.Children[index],
	}, true
}

func (n yamlNode) KeyValueMap() (map[string]validate.AstNode, bool) {
	if n.Kind != yaml.MappingNode {
		return nil, false
	}

	kvmap := map[string]validate.AstNode{}
	for i := 0; i < len(n.Children); i += 2 {
		key := *n.Children[i]
		value := *n.Children[i+1]
		kvmap[key.Value] = yamlNode{
			key:  key,
			Node: value,
		}
	}
	return kvmap, true
}

func (n yamlNode) Tag() string {
	return "yaml"
}
