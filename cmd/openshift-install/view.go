package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
)

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "view [filename]",
		Short:   "View files created by the assets",
		Long:    "View files created by the assets.\nIf run with a filename, view the contents of that file.\nIf run without a filename, list all of the files for all of the generated assets.",
		Example: "  openshift-install view\n  openshift-install view install-config.yaml",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				listFiles(rootOpts.dir)
			} else {
				viewFile(rootOpts.dir, args[0])
			}
		},
	}
}

func listFiles(directory string) {
	assetStore, err := loadAssets(directory)
	if err != nil {
		logrus.Error(err)
		return
	}
	files := assetStore.GetAllFiles()
	sort.Strings(files)
	for _, f := range files {
		fmt.Println(f)
	}
}

func viewFile(directory string, filename string) {
	assetStore, err := loadAssets(directory)
	if err != nil {
		logrus.Error(err)
		return
	}
	file := assetStore.GetFile(filename)
	if file == nil {
		logrus.Warnf("None of the generated assets write out the %q file", filename)
		return
	}
	os.Stdout.Write(file.Data)
}

func loadAssets(directory string) (asset.Store, error) {
	assetStore, err := asset.NewStore(directory)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create asset store")
	}
	_, err = assetStore.Load(&cluster.Cluster{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load assets")
	}
	return assetStore, nil
}
