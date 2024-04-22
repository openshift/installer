package mirror

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/distribution/distribution/v3/registry/client/auth"
	digest "github.com/opencontainers/go-digest"

	"github.com/openshift/oc/pkg/cli/image/imagesource"
)

// ErrAlreadyExists may be returned by the blob Create function to indicate that the blob already exists.
var ErrAlreadyExists = fmt.Errorf("blob already exists in the target location")

type Mapping struct {
	Source      imagesource.TypedImageReference
	Destination imagesource.TypedImageReference
	// Name is an optional field for identifying uniqueness within the mappings
	Name string
}

func parseArgs(args []string, overlap map[string]string, expandFn func(s imagesource.TypedImageReference) ([]imagesource.TypedImageReference, error)) ([]Mapping, error) {
	var remainingArgs []string
	var mappingParts [][]string
	for _, s := range args {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			remainingArgs = append(remainingArgs, s)
			continue
		}
		if len(parts[0]) == 0 || len(parts[1]) == 0 {
			return nil, fmt.Errorf("all arguments must be valid SRC=DST mappings")
		}
		mappingParts = append(mappingParts, parts)
	}

	switch {
	case len(remainingArgs) > 1 && len(mappingParts) == 0:
		for i := 1; i < len(remainingArgs); i++ {
			if len(remainingArgs[i]) == 0 {
				continue
			}
			mappingParts = append(mappingParts, []string{remainingArgs[0], remainingArgs[i]})
		}
	case len(remainingArgs) == 1 && len(mappingParts) == 0:
		return nil, fmt.Errorf("all arguments must be valid SRC=DST mappings, or you must specify one SRC argument and one or more DST arguments")
	}

	var mappings []Mapping
	for _, parts := range mappingParts {
		sources, err := imagesource.ParseSourceReference(parts[0], expandFn)
		if err != nil {
			return nil, err
		}
		dst, err := imagesource.ParseDestinationReference(parts[1])
		if err != nil {
			return nil, err
		}
		if len(sources) > 1 && (len(dst.Ref.Tag) > 0 || len(dst.Ref.ID) > 0) {
			return nil, fmt.Errorf("when source contains wildcards, the destination must be a repository")
		}

		for _, src := range sources {
			if len(src.Ref.Tag) == 0 && len(src.Ref.ID) == 0 {
				return nil, fmt.Errorf("you must specify a tag or digest for SRC")
			}
			copied := dst
			if len(dst.Ref.Tag) == 0 && len(src.Ref.Tag) > 0 {
				copied.Ref.Tag = src.Ref.Tag
			}
			if _, ok := overlap[copied.String()]; ok {
				return nil, fmt.Errorf("each destination tag may only be specified once: %s", copied.String())
			}
			overlap[copied.String()] = src.String()

			mappings = append(mappings, Mapping{Source: src, Destination: copied})
		}
	}

	return mappings, nil
}

func parseFile(filename string, overlap map[string]string, in io.Reader, expandFn func(s imagesource.TypedImageReference) ([]imagesource.TypedImageReference, error)) ([]Mapping, error) {
	var fileMappings []Mapping
	if filename != "-" {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		in = f
	}
	s := bufio.NewScanner(in)
	lineNumber := 0
	for s.Scan() {
		line := s.Text()
		lineNumber++

		// remove comments and whitespace
		if i := strings.Index(line, "#"); i != -1 {
			line = line[0:i]
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		args := strings.Split(line, " ")
		mappings, err := parseArgs(args, overlap, expandFn)
		if err != nil {
			return nil, fmt.Errorf("file %s, line %d: %v", filename, lineNumber, err)
		}
		fileMappings = append(fileMappings, mappings...)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return fileMappings, nil
}

type key struct {
	t          imagesource.DestinationType
	registry   string
	repository string
}

type destination struct {
	ref  imagesource.TypedImageReference
	tags []string
}

func contextKeyForReference(t imagesource.TypedImageReference) contextKey {
	return contextKey{t: t.Type, registry: t.Ref.Registry}
}

type pushTargets map[key]destination

type destinations struct {
	ref imagesource.TypedImageReference

	lock    sync.Mutex
	tags    map[string]pushTargets
	digests map[string]pushTargets
}

func (d *destinations) mergeIntoDigests(srcDigest digest.Digest, target pushTargets) {
	d.lock.Lock()
	defer d.lock.Unlock()
	srcKey := srcDigest.String()
	current, ok := d.digests[srcKey]
	if !ok {
		d.digests[srcKey] = target
		return
	}
	for repo, dst := range target {
		existing, ok := current[repo]
		if !ok {
			current[repo] = dst
			continue
		}
		existing.tags = append(existing.tags, dst.tags...)
	}
}

type targetTree map[key]*destinations

func buildTargetTree(mappings []Mapping) targetTree {
	tree := make(targetTree)
	for _, m := range mappings {
		srcKey := key{t: m.Source.Type, registry: m.Source.Ref.Registry, repository: m.Source.Ref.RepositoryName()}
		dstKey := key{t: m.Destination.Type, registry: m.Destination.Ref.Registry, repository: m.Destination.Ref.RepositoryName()}

		src, ok := tree[srcKey]
		if !ok {
			src = &destinations{}
			src.ref = imagesource.TypedImageReference{Ref: m.Source.Ref.AsRepository(), Type: m.Source.Type}
			src.digests = make(map[string]pushTargets)
			src.tags = make(map[string]pushTargets)
			tree[srcKey] = src
		}

		var current pushTargets
		if id := m.Source.Ref.ID; len(id) > 0 {
			current = src.digests[m.Source.Ref.ID]
			if current == nil {
				current = make(pushTargets)
				src.digests[m.Source.Ref.ID] = current
			}
		} else {
			tag := m.Source.Ref.Tag
			current = src.tags[tag]
			if current == nil {
				current = make(pushTargets)
				src.tags[tag] = current
			}
		}

		dst, ok := current[dstKey]
		if !ok {
			dst.ref = imagesource.TypedImageReference{Ref: m.Destination.Ref.AsRepository(), Type: m.Destination.Type}
		}
		if len(m.Destination.Ref.Tag) > 0 {
			dst.tags = append(dst.tags, m.Destination.Ref.Tag)
		}
		current[dstKey] = dst
	}
	return tree
}

func addDockerRegistryScopes(scopes map[contextKey]map[string]bool, targets map[string]pushTargets, srcKey key) {
	for _, target := range targets {
		for dstKey, t := range target {
			dstContextKey := contextKey{t: dstKey.t, registry: dstKey.registry}
			m := scopes[dstContextKey]
			if m == nil {
				m = make(map[string]bool)
				scopes[dstContextKey] = m
			}
			m[dstKey.repository] = true
			if t.ref.Type != imagesource.DestinationRegistry || dstKey.registry != srcKey.registry || dstKey.repository == srcKey.repository {
				continue
			}
			srcContextKey := contextKey{t: srcKey.t, registry: srcKey.registry}
			m = scopes[srcContextKey]
			if m == nil {
				m = make(map[string]bool)
				scopes[srcContextKey] = m
			}
			if _, ok := m[srcKey.repository]; !ok {
				m[srcKey.repository] = false
			}
		}
	}
}

func calculateDockerRegistryScopes(tree targetTree) map[contextKey][]auth.Scope {
	scopes := make(map[contextKey]map[string]bool)
	for srcKey, dst := range tree {
		addDockerRegistryScopes(scopes, dst.tags, srcKey)
		addDockerRegistryScopes(scopes, dst.digests, srcKey)
	}
	uniqueScopes := make(map[contextKey][]auth.Scope)
	for registry, repos := range scopes {
		var repoScopes []auth.Scope
		for name, push := range repos {
			if push {
				repoScopes = append(repoScopes, auth.RepositoryScope{Repository: name, Actions: []string{"pull", "push"}})
			} else {
				repoScopes = append(repoScopes, auth.RepositoryScope{Repository: name, Actions: []string{"pull"}})
			}
		}
		uniqueScopes[registry] = repoScopes
	}
	return uniqueScopes
}
