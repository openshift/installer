package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	outputJSON     = false
	hiddenFeatures = []string{
		"terraform-spot-masters",
	}
)

func newListFeaturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list-hidden-features",
		Short:  "List supported hidden features",
		Long:   "",
		Hidden: true,
		Args:   cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			var out string
			if outputJSON {
				outb, err := json.Marshal(hiddenFeatures)
				if err != nil {
					logrus.Fatalf("failed to marshal output: %s", err.Error())
				}
				out = string(outb)
			} else {
				out = strings.Join(hiddenFeatures, "\n")
			}
			fmt.Println(out)
		},
	}
	cmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "print output in json format")
	return cmd
}
