package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/awalterschulze/gographviz"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
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
	g := gographviz.NewGraph()
	g.SetName("G")
	g.SetDir(true)
	g.SetStrict(true)

	tNodeAttr := map[string]string{
		string(gographviz.Shape): "box",
		string(gographviz.Style): "filled",
	}
	for _, t := range targets {
		name := fmt.Sprintf(`"Target %s"`, t.name)
		g.AddNode("G", name, tNodeAttr)
		for _, dep := range t.assets {
			addEdge(g, name, dep)
		}
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
	return nil
}

func addEdge(g *gographviz.Graph, parent string, asset asset.Asset) {
	elem := reflect.TypeOf(asset).Elem()
	name := fmt.Sprintf(`"%s"`, elem.Name())

	if !g.IsNode(name) {
		logrus.Debugf("adding node %s", name)
		g.AddNode("G", name, nil)
	}
	if !isEdge(g, name, parent) {
		logrus.Debugf("adding edge %s -> %s", name, parent)
		g.AddEdge(name, parent, true, nil)
	}

	deps := asset.Dependencies()
	for _, dep := range deps {
		addEdge(g, name, dep)
	}
}

func isEdge(g *gographviz.Graph, src, dst string) bool {
	for _, edge := range g.Edges.Edges {
		if edge.Src == src && edge.Dst == dst {
			return true
		}
	}
	return false
}
