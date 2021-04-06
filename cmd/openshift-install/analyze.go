package main

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/gather/service"
)

var (
	analyzeOpts struct {
		gatherBundle string
	}
)

func newAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze debugging data for a given installation failure",
		Long: `Analyze debugging data for a given installation failure.

This command helps users to analyze the reasons for an installation that failed while bootstrapping.`,
		Args: cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			gatherBundle := analyzeOpts.gatherBundle
			if gatherBundle == "" {
				var err error
				gatherBundle, err = getGatherBundleFromAssetsDirectory()
				if err != nil {
					logrus.Fatal(err)
				}
			}
			if !filepath.IsAbs(gatherBundle) {
				gatherBundle = filepath.Join(rootOpts.dir, gatherBundle)
			}
			if err := service.AnalyzeGatherBundle(gatherBundle); err != nil {
				logrus.Fatal(err)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&analyzeOpts.gatherBundle, "file", "", "Filename of the bootstrap gather bundle; either absolute or relative to the assets directory")
	return cmd
}

func getGatherBundleFromAssetsDirectory() (string, error) {
	matches, err := filepath.Glob(filepath.Join(rootOpts.dir, "log-bundle-*.tar.gz"))
	if err != nil {
		return "", errors.Wrap(err, "could not find gather bundles in assets directory")
	}
	switch len(matches) {
	case 0:
		return "", errors.New("no bootstrap gather bundles found in assets directory")
	case 1:
		return matches[0], nil
	default:
		return "", errors.New("multiple bootstrap gather bundles found in assets directory; select specific gather bundle by using the --file flag")
	}
}
