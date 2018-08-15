package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/go-log/log/print"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

var (
	graphOpts struct {
		outputFile string
	}
)

func newGraphCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Outputs the internal dependency graph for installer",
		Long:  "",
		RunE:  runGraphCmd,
	}
	cmd.PersistentFlags().StringVar(&graphOpts.outputFile, "output-file", "", "file where the graph is written, if empty prints the graph to Stdout.")
	return cmd
}

func runGraphCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	g := gographviz.NewGraph()
	g.SetName("G")
	g.SetDir(true)
	g.SetStrict(true)

	installerAssets := installerassets.New()
	err := installerAssets.Read(ctx, rootOpts.dir, installerassets.GetDefault, print.New(logrus.StandardLogger()))
	if err != nil {
		logrus.Fatal(err)
	}

	root, err := installerAssets.GetByHash(ctx, installerAssets.Root.Hash)
	if err != nil {
		logrus.Fatal(err)
	}

	directories := make(map[string]float32)
	seen := make(map[string]bool)
	stack := []*assets.Asset{&root}
	for len(stack) > 0 {
		asset := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[asset.Name] {
			continue
		}

		directories[path.Dir(asset.Name)] = 1
		for _, reference := range asset.Parents {
			parent, err := installerAssets.GetByHash(ctx, reference.Hash)
			if err != nil {
				logrus.Fatal(err)
			}
			stack = append(stack, &parent)
		}
	}
	dirSlice := make([]string, 0, len(directories))
	for dir := range directories {
		dirSlice = append(dirSlice, dir)
	}
	sort.Strings(dirSlice)
	for i, dir := range dirSlice {
		directories[dir] = float32(i) / float32(len(dirSlice))
	}

	added := make(map[string]bool)
	err = addNodes(ctx, g, &root, installerAssets.GetByHash, added, directories)
	if err != nil {
		logrus.Fatal(err)
	}

	out := os.Stdout
	if graphOpts.outputFile != "" {
		f, err := os.Create(graphOpts.outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	if _, err := io.WriteString(out, g.String()); err != nil {
		return err
	}

	var unused []string
	for key := range installerassets.Defaults {
		if _, ok := added[key]; !ok {
			unused = append(unused, key)
		}
	}
	for key := range installerassets.Rebuilders {
		if _, ok := added[key]; !ok {
			unused = append(unused, key)
		}
	}
	sort.Strings(unused)
	if unused != nil {
		logrus.Warnf("potentially unused asset(s): %s", strings.Join(unused, ", "))
	}

	return nil
}

func addNodes(ctx context.Context, g *gographviz.Graph, asset *assets.Asset, getByHash assets.GetByBytes, added map[string]bool, directories map[string]float32) (err error) {
	_, ok := added[asset.Name]
	if ok {
		return nil
	}
	added[asset.Name] = true

	assetName := fmt.Sprintf("%q", asset.Name)
	attrs := make(map[string]string)
	hue, ok := directories[path.Dir(asset.Name)]
	if ok {
		saturation := 0.1
		if asset.Name == "tls/kubelet-client.crt" {
			saturation = 0.3
		}
		attrs[string(gographviz.FillColor)] = fmt.Sprintf("\"%.2f %.2f 1\"", hue, saturation)

		attrs[string(gographviz.Style)] = "filled"
	}
	g.AddNode("G", assetName, attrs)

	for _, parentReference := range asset.Parents {
		parent, err := getByHash(ctx, parentReference.Hash)
		if err != nil {
			return err
		}

		err = addNodes(ctx, g, &parent, getByHash, added, directories)
		if err != nil {
			return err
		}

		parentName := fmt.Sprintf("%q", parent.Name)
		g.AddEdge(parentName, assetName, true, nil)
	}

	return nil
}
