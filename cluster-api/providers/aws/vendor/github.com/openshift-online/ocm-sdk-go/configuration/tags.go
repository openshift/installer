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

// This file contains the code that is used to process tags like `!variable` and `!file` inside
// configuration files.

package configuration

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

// tagRegistryEntry stores the description of one tag. When a tag is detected in a node the
// corresponding function will be called passing the node, so that the function can modify
// it as needed.
type tagRegistryEntry struct {
	name    string
	kind    yaml.Kind
	process func(*yaml.Node) error
}

// registerTag adds a tag to the registry.
func (b *Builder) registerTag(name string, kind yaml.Kind, process func(*yaml.Node) error) {
	b.tags = append(b.tags, tagRegistryEntry{
		name:    name,
		kind:    kind,
		process: process,
	})
}

// lookupTag tries to find an entry in the tag registry corresponding to the given name.
func (b *Builder) lookupTag(name string) (result tagRegistryEntry, ok bool) {
	for _, entry := range b.tags {
		if strings.HasPrefix(entry.name, name) {
			result = entry
			ok = true
			break
		}
	}
	return
}

// processTreeTags process recursively the tags present in the given YAML tree and returns the
// result.
func (b *Builder) processTreeTags(tree *yaml.Node) error {
	err := b.processNodeTags(tree)
	if err != nil {
		return err
	}
	for _, node := range tree.Content {
		err = b.processTreeTags(node)
		if err != nil {
			return err
		}
	}
	return nil
}

// processNodeTags process the tags present in the given YAML node.
func (b *Builder) processNodeTags(node *yaml.Node) error {
	tag := node.ShortTag()
	if strings.HasPrefix(tag, "!!") {
		return nil
	}
	names := b.parseTag(tag)
	for _, name := range names {
		entry, ok := b.lookupTag(name)
		if ok {
			if entry.kind != node.Kind {
				return b.nodeError(
					node,
					"tag '%s' expects %s node but found %s",
					name,
					b.kindString(entry.kind),
					b.kindString(node.Kind),
				)
			}
			err := entry.process(node)
			if err != nil {
				return err
			}
		} else {
			return b.nodeError(node, "unsupported tag '%s'", name)
		}
	}
	return nil
}

// processVariableTag is the implementation fo the `variable` tag: replaces an environment variable
// reference with its value.
func (b *Builder) processVariableTag(node *yaml.Node) error {
	variable := node.Value
	result, ok := os.LookupEnv(variable)
	if !ok {
		return b.nodeError(node, "can't find environment variable '%s'", variable)
	}
	node.SetString(result)
	return nil
}

// processFileTag is the implementation of the `file` tag: replaces a file reference with the
// content of the file.
func (b *Builder) processFileTag(node *yaml.Node) error {
	file := node.Value
	data, err := os.ReadFile(file) // #nosec G304
	if err != nil {
		return b.nodeError(node, "%w", err)
	}
	result := string(data)
	node.SetString(result)
	b.titleTree(file, node)
	return nil
}

// processScriptTag is the implementation of the `script` tag: replaces a script script with
// the result of executing it.
func (b *Builder) processScriptTag(node *yaml.Node) error {
	shell, ok := os.LookupEnv("SHELL")
	if !ok {
		shell = "/bin/sh"
	}
	script := node.Value
	stdin := strings.NewReader(script)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := &exec.Cmd{
		Path:   shell,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	err := cmd.Run()
	_, exit := err.(*exec.ExitError)
	if exit {
		err = nil
	}
	if err != nil {
		return b.nodeError(node, "%w", err)
	}
	code := cmd.ProcessState.ExitCode()
	if code != 0 || stderr.Len() > 0 {
		return b.nodeError(
			node,
			"script '%s' finished with exit code %d, wrote '%s' to stdout and '%s' "+
				"to stderr",
			b.quoteForError(script),
			code,
			b.quoteForError(stdout.String()),
			b.quoteForError(stderr.String()),
		)
	}
	result := stdout.String()
	node.SetString(result)
	return nil
}

// processTrimTag is the implementation of the `trim` tag: trims leading and trailing white
// space, including line breaks and tabs.
func (b *Builder) processTrimTag(node *yaml.Node) error {
	value := node.Value
	value = strings.TrimSpace(value)
	node.SetString(value)
	return nil
}

// processYamlTag is the implementation of the `yaml` tag: parses the value as a YAML
// document and replaces the current node with the result.
func (b *Builder) processYamlTag(node *yaml.Node) error {
	buffer := []byte(node.Value)
	var parsed yaml.Node
	err := yaml.Unmarshal(buffer, &parsed)
	if err != nil {
		return b.nodeError(node, "can't parse '%s': %w", buffer, err)
	}
	b.titleTree(b.titles[node], &parsed)
	err = b.processTreeTags(&parsed)
	if err != nil {
		return err
	}
	*node = *parsed.Content[0]
	return nil
}

// processStrTag is the implementation of the `string` tag: it explicitly sets the node kind
// to `!!str`.
func (b *Builder) processStringTag(node *yaml.Node) error {
	node.Tag = "!!str"
	return nil
}

// processIntegerTag is the implementation of the `integer` tag: it explicitly sets the node kind
// to `!!int`.
func (b *Builder) processIntegerTag(node *yaml.Node) error {
	node.Tag = "!!int"
	return nil
}

// processBooleanTag is the implementation of the `bool` tag: it explicitly sets the node kind
// to `!!bool`.
func (b *Builder) processBooleanTag(node *yaml.Node) error {
	node.Tag = "!!bool"
	return nil
}

// processFloatTag is the implementation of the `float` tag: it explicitly sets the node type
// to `!!float`.
func (b *Builder) processFloatTag(node *yaml.Node) error {
	node.Tag = "!!float"
	return nil
}

// parseTag extracts the coponents from the given tags. For example, the tag `!file/int` will be
// parsed as a slice containing `file` and `int`.
func (b *Builder) parseTag(tag string) []string {
	if !strings.HasPrefix(tag, "!") {
		return nil
	}
	return strings.Split(tag[1:], "/")
}
