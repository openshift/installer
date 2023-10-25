/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains the the implementation of the configuration object.

package configuration

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Builder contains the data and logic needed to populate a configuration object. Don't create
// instances of this type directly, use the New function instead.
type Builder struct {
	// sources contains the list of sources where the configuration should be loaded from.
	sources []interface{}

	// tags contains for each tag information describing it, like its name and the kind of nodes
	// that it can process.
	tags []tagRegistryEntry

	// titles contains the title for each node. This will usually be the name of the file that
	// the node was loaded from. This is used to generate error messages that include the name
	// of the file.
	titles map[*yaml.Node]string
}

// Object contains configuration data.
type Object struct {
	tree *yaml.Node
}

// New creates a new builder that can be use to populate a configuration object.
func New() *Builder {
	return &Builder{}
}

// Load adds the given objects as sources where the configuration will be loaded from.
//
// If a source is a string ending in `.yaml` or `.yml` it will be interpreted as the name of a
// file containing the YAML text.
//
// If a source is a string ending in `.d` it will be interpreted as a directory containing YAML
// files. The `.yaml` or `.yml` files inside that directory will be loaded in alphabetical order.
//
// Any string not ending in `.yaml`, `.yml` or `.d` will be interprested as actual YAML text. In
// order to simplify embedding these strings in Go programs leading tabs will be removed from all
// the lines of that YAML text.
//
// If a source is an array of bytes it will be interpreted as actual YAML text.
//
// If a source implements the io.Reader interface, then it will be used to read in memory the YAML
// text.
//
// If the source can also be a yaml.Node or another configuration Object. In those cases the
// content of the source will be copied.
//
// If the source is any other kind of object then it will be serialized as YAML and then loaded.
func (b *Builder) Load(sources ...interface{}) *Builder {
	b.sources = append(b.sources, sources...)
	return b
}

// Build uses the information stored in the builder to create and populate a configuration
// object.
func (b *Builder) Build() (object *Object, err error) {
	// Add the builtin tags to the tag registry:
	b.registerTag("boolean", yaml.ScalarNode, b.processBooleanTag)
	b.registerTag("file", yaml.ScalarNode, b.processFileTag)
	b.registerTag("float", yaml.ScalarNode, b.processFloatTag)
	b.registerTag("integer", yaml.ScalarNode, b.processIntegerTag)
	b.registerTag("script", yaml.ScalarNode, b.processScriptTag)
	b.registerTag("string", yaml.ScalarNode, b.processStringTag)
	b.registerTag("trim", yaml.ScalarNode, b.processTrimTag)
	b.registerTag("variable", yaml.ScalarNode, b.processVariableTag)
	b.registerTag("yaml", yaml.ScalarNode, b.processYamlTag)

	// Initialize the titles index:
	b.titles = map[*yaml.Node]string{}

	// Merge the sources:
	tree := &yaml.Node{}
	for _, current := range b.sources {
		switch source := current.(type) {
		case string:
			err = b.mergeString(source, tree)
		case []byte:
			err = b.mergeBytes("", source, tree)
		case io.Reader:
			err = b.mergeReader(source, tree)
		case yaml.Node:
			err = b.mergeNode(&source, tree)
		case *yaml.Node:
			err = b.mergeNode(source, tree)
		case *Object:
			err = b.mergeNode(source.tree, tree)
		case Object:
			err = b.mergeNode(source.tree, tree)
		default:
			err = b.mergeAny(source, tree)
		}
		if err != nil {
			return
		}
	}

	// Create and populate the object:
	object = &Object{
		tree: tree,
	}

	return
}

func (b *Builder) mergeString(src string, dst *yaml.Node) error {
	ext := filepath.Ext(src)
	if ext == ".yaml" || ext == ".yml" || ext == ".d" {
		return b.mergeFile(src, dst)
	}
	src = b.removeLeadingTabs(src)
	return b.mergeBytes("", []byte(src), dst)
}

func (b *Builder) mergeReader(src io.Reader, dst *yaml.Node) error {
	buffer, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	return b.mergeBytes("", buffer, dst)
}

func (b *Builder) mergeFile(src string, dst *yaml.Node) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return b.mergeDir(src, dst)
	}
	buffer, err := os.ReadFile(src) // #nosec G304
	if err != nil {
		return err
	}
	return b.mergeBytes(src, buffer, dst)
}

func (b *Builder) mergeDir(src string, dst *yaml.Node) error {
	infos, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(infos))
	for _, info := range infos {
		name := info.Name()
		ext := filepath.Ext(name)
		if ext == ".yaml" || ext == ".yml" {
			files = append(files, filepath.Join(src, name))
		}
	}
	sort.Strings(files)
	for _, file := range files {
		err = b.mergeFile(file, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) mergeAny(src interface{}, dst *yaml.Node) error {
	buffer, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	return b.mergeBytes("", buffer, dst)
}

func (b *Builder) mergeBytes(title string, src []byte, dst *yaml.Node) error {
	var tree yaml.Node
	err := yaml.Unmarshal(src, &tree)
	if err != nil {
		return err
	}
	b.titleTree(title, &tree)
	err = b.processTreeTags(&tree)
	if err != nil {
		return err
	}
	return b.mergeNode(&tree, dst)
}

func (b *Builder) mergeNode(src, dst *yaml.Node) error {
	if src.Kind != dst.Kind {
		b.deepCopy(src, dst)
		return nil
	}
	switch src.Kind {
	case 0:
		return b.mergeEmpty(src, dst)
	case yaml.DocumentNode:
		return b.mergeDocument(src, dst)
	case yaml.SequenceNode:
		return b.mergeSequence(src, dst)
	case yaml.MappingNode:
		return b.mergeMapping(src, dst)
	case yaml.ScalarNode:
		return b.mergeScalar(src, dst)
	case yaml.AliasNode:
		return b.mergeAlias(src, dst)
	default:
		return fmt.Errorf("don't know how to handle YAML node of type %d", src.Kind)
	}
}

func (b *Builder) mergeDocument(src, dst *yaml.Node) error {
	return b.mergeNode(src.Content[0], dst.Content[0])
}

func (b *Builder) mergeEmpty(src, dst *yaml.Node) error {
	return nil
}

func (b *Builder) mergeSequence(src, dst *yaml.Node) error {
	size := len(src.Content)
	nodes := make([]*yaml.Node, size)
	for i := 0; i < size; i++ {
		nodes[i] = &yaml.Node{}
		b.deepCopy(src.Content[i], nodes[i])
	}
	dst.Content = append(dst.Content, nodes...)
	return nil
}

func (b *Builder) mergeMapping(src, dst *yaml.Node) error {
	srcSize := len(src.Content) / 2
	i := 0
	for i < srcSize {
		srcKey := src.Content[2*i]
		srcValue := src.Content[2*i+1]
		dstSize := len(dst.Content) / 2
		j := 0
		for j < dstSize {
			dstKey := dst.Content[2*j]
			dstValue := dst.Content[2*j+1]
			if srcKey.Value == dstKey.Value {
				err := b.mergeNode(srcValue, dstValue)
				if err != nil {
					return err
				}
				break
			}
			j++
		}
		if j == dstSize {
			dstKey := &yaml.Node{}
			b.deepCopy(srcKey, dstKey)
			dstValue := &yaml.Node{}
			b.deepCopy(srcValue, dstValue)
			dst.Content = append(dst.Content, dstKey, dstValue)
		}
		i++
	}
	return nil
}

func (b *Builder) mergeScalar(src, dst *yaml.Node) error {
	b.deepCopy(src, dst)
	return nil
}

func (b *Builder) mergeAlias(src, dst *yaml.Node) error {
	b.deepCopy(src, dst)
	return nil
}

func (b *Builder) deepCopy(src, dst *yaml.Node) {
	// Copy the title:
	b.titles[dst] = b.titles[src]

	// Copy the content:
	*dst = *src
	if src.Content != nil {
		size := len(src.Content)
		dst.Content = make([]*yaml.Node, size)
		for i := 0; i < size; i++ {
			dst.Content[i] = &yaml.Node{}
			b.deepCopy(src.Content[i], dst.Content[i])
		}
	}
}

func (b *Builder) nodeError(node *yaml.Node, format string, a ...interface{}) error {
	format = "%s:%d:%d: " + format
	title := b.titles[node]
	if title == "" {
		title = "unknown"
	}
	v := make([]interface{}, len(a)+3)
	v[0], v[1], v[2] = title, node.Line, node.Column
	copy(v[3:], a)
	return fmt.Errorf(format, v...)
}

func (b *Builder) titleTree(title string, tree *yaml.Node) {
	b.titles[tree] = title
	for _, node := range tree.Content {
		b.titleTree(title, node)
	}
}

func (b *Builder) quoteForError(s string) string {
	r := strconv.Quote(s)
	return r[1 : len(r)-1]
}

func (b *Builder) kindString(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "document"
	case yaml.SequenceNode:
		return "sequence"
	case yaml.MappingNode:
		return "mapping"
	case yaml.ScalarNode:
		return "scalar"
	case yaml.AliasNode:
		return "alias"
	}
	return ""
}

// removeLeadingTabs removes the leading tabs from the lines of the given string.
func (b *Builder) removeLeadingTabs(s string) string {
	return leadingTabsRE.ReplaceAllString(s, "")
}

// leadingTabsRE is the regular expression used to remove leading tabs from strings generated with
// the EvaluateTemplate function.
var leadingTabsRE = regexp.MustCompile(`(?m)^\t*`)

// Populate populates the given destination object with the information stored in this
// configuration object. The destination object should be a pointer to a variable containing
// the same tags used by the yaml.Unmarshal method of the YAML library.
func (o *Object) Populate(v interface{}) error {
	return o.tree.Decode(v)
}

// Effective returns an array of bytes containing the YAML representation of the configuration
// after processing all the tags.
func (o *Object) Effective() (out []byte, err error) {
	return yaml.Marshal(o.tree)
}

// MarshalYAML is the implementation of the yaml.Marshaller interface. This is intended to be able
// use the type for fields inside other structs. Refrain from calling this method for any other
// use.
func (o *Object) MarshalYAML() (result interface{}, err error) {
	result = o.tree
	return
}

// UnmarshalYAML is the implementation of the yaml.Unmarshaller interface. This is intended to be
// able to use the type for fields inside structs. Refraim from calling this method for any other
// use.
func (o *Object) UnmarshalYAML(value *yaml.Node) error {
	o.tree = value
	return nil
}
