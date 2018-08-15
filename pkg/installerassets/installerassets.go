// Package installerassets contains installer-specific helpers for the
// asset Merkle DAG.
package installerassets

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/openshift/installer/pkg/assets"
	"github.com/pkg/errors"
)

// Rebuilders registers installer asset rebuilders by name.  Use this
// to set up assets that have parents.  For example:
//
//   func yourRebuilder(getByName assets.GetByString) (asset *assets.Asset, err error) {
//     asset = &assets.Asset{
//       Name:          "tls/root-ca.crt",
//       RebuildHelper: rootCARebuilder,
//     }
//
//     parents, err := asset.GetParents(getByName, "tls/root-ca.key")
//     if err != nil {
//       return nil, err
//     }
//
//     // Assemble your data based on the parent content using your custom logic.
//     for name, parent := range parents {
//       asset.Data = append(asset.Data, parent.Data)
//     }
//
//     return asset, nil
//   }
//
// and then somewhere (e.g. an init() function), add it to the registry:
//
//   Rebuilders["your/asset"] = yourRebuilder
var Rebuilders = make(map[string]assets.Rebuild)

// Defaulter returns a default value.  This type is consumed by the
// Defaults registry and helpers interacting with that registry.
type Defaulter func(ctx context.Context) (data []byte, err error)

// Defaults registers installer asset default functions by name.  Use
// this to set up assets that do not have parents.  For example,
// constants:
//
//   Defaults["your/asset"] = ConstantDefault([]byte("your value"))
//
// or values populated from outside the asset graph:
//
//   Defaults["your/asset"] = func() ([]byte, error) {
//     value = os.Getenv("YOUR_ENVIRONMENT_VARIABLE")
//     return []byte(value), nil
//   }
var Defaults = make(map[string]Defaulter)

// New returns a new installer asset store.
func New() *assets.Assets {
	return &assets.Assets{
		Root: assets.Reference{
			Name: "cluster",
		},
		Rebuilders: Rebuilders,
	}
}

// GetDefault calculates defaults for missing assets.
func GetDefault(ctx context.Context, name string) ([]byte, error) {
	defaulter, ok := Defaults[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	return defaulter(ctx)
}

// ConstantDefault returns a Defaulter which returns a constant value.
func ConstantDefault(value []byte) Defaulter {
	return func(ctx context.Context) ([]byte, error) {
		return value, nil
	}
}

// PlatformOverrideRebuilder generates rebuilders for values that
// allow per-platform overrides.  It pulls the 'platform' parent and
// looks up the {platform}/{name} asset.  If that asset exists, it uses
// its value.  If that asset does not exist, it uses 'defaulter' to
// calculate a generic default.
func PlatformOverrideRebuilder(name string, defaulter Defaulter) assets.Rebuild {
	rebuild := func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          name,
			RebuildHelper: PlatformOverrideRebuilder(name, defaulter),
		}

		parents, err := asset.GetParents(ctx, getByName, "platform")
		if err != nil {
			return nil, err
		}

		platform := string(parents["platform"].Data)
		perPlatformName := path.Join(platform, name)
		parents, err = asset.GetParents(ctx, getByName, perPlatformName)
		if err == nil {
			asset.Data = parents[perPlatformName].Data
			return asset, nil
		} else if os.IsNotExist(errors.Cause(err)) {
			asset.Data, err = defaulter(ctx)
			if err != nil {
				return nil, err
			}
			return asset, nil
		}

		return nil, err
	}

	return rebuild
}

// TemplateRebuilder returns a rebuilder that pulls a template from
// {name}.template and renders it with the parameters map for context.
// Keys for 'parameters' are template properties and values are parent
// names.  You can also inject additional parameters directly without
// involving parent assets.  For example:
//
//   TemplateRebuilder(
//     "my/asset",
//     map[string]string{"my/parent": "MyParent"},
//     map[string]interface{}{"MyExtra": "my extra"},
//   )
//
// will pull the template from my/asset.template and the template can
// use {{.MyParent}} to access the data from my/parent and
// {{.MyExtra}} to access "my extra".
//
// The following functions are also available in the template:
//
// * add, which returns the sum of its arguments.  For example:
//
//     {{add $index 1}}
//
// * base64, which encodes its argument in base64.  For example:
//
//     {{.Key | base64}}
//
// * etcdURIs, which, when called with clusterName, baseDomain, and
//   count, returns a []string of
//   https://{clusterName}-etcd-{count}.{baseDomain}:2379 URIs.
//
// * indent, which takes a count and a string and appends count spaces
//   to any newlines in the string.
//
// * int, which converts a base-10 string into an int.
//
// * join, which, when called with a separator and a slice of strings,
//   returns the slice elements separated by the separator.
func TemplateRebuilder(name string, parameters map[string]string, extra map[string]interface{}) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          name,
			RebuildHelper: TemplateRebuilder(name, parameters, extra),
		}

		templateName := fmt.Sprintf("%s.template", name)
		parents, err := asset.GetParents(ctx, getByName, templateName)
		if err != nil {
			return nil, err
		}
		templateData := string(parents[templateName].Data)

		funcMap := template.FuncMap{
			"add": func(values ...int) int {
				sum := 0
				for _, value := range values {
					sum += value
				}
				return sum
			},
			"base64": func(data string) string {
				return base64.StdEncoding.EncodeToString([]byte(data))
			},
			"etcdURIs": func(clusterName string, baseDomain string, count int) []string {
				uris := make([]string, 0, count)
				for i := 0; i < count; i++ {
					uris = append(uris, fmt.Sprintf("https://%s-etcd-%d.%s:2379", clusterName, i, baseDomain))
				}
				return uris
			},
			"indent": func(indentation int, value string) string {
				newline := "\n" + strings.Repeat(" ", indentation)
				return strings.Replace(value, "\n", newline, -1)
			},
			"int": func(value string) (int, error) {
				integer, err := strconv.ParseInt(value, 10, 0)
				return int(integer), err
			},
			"join": func(separator string, slice []string) string {
				return strings.Join(slice, separator)
			},
		}

		tmpl, err := template.New(name).Funcs(funcMap).Parse(templateData)
		if err != nil {
			return nil, err
		}

		params := make(map[string]interface{}, len(parameters)+len(extra))
		for key, parentName := range parameters {
			parents, err = asset.GetParents(ctx, getByName, parentName)
			if err != nil {
				return nil, err
			}

			params[key] = string(parents[parentName].Data)
		}
		for key, value := range extra {
			params[key] = value
		}

		buf := &bytes.Buffer{}
		err = tmpl.Option("missingkey=error").Execute(buf, params)
		if err != nil {
			return nil, err
		}

		asset.Data = buf.Bytes()
		return asset, err
	}
}
