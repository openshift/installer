package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ignv2 "github.com/coreos/ignition/config/v2_4"
	"github.com/coreos/ign-converter/translate/v24tov31"
)

func newIgnitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ignition",
		Short: "Operations on Ignition configs",
		Long: `Operations on Ignition configs

See sub-operations for details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newIgnitionConvertCmd())
	return cmd
}

var (
	ignitionConvertOpts struct {
		source string
		dest   string
	}
)

func newIgnitionConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert a Spec2 config to Spec 3",
		Long: `Convert a Spec2 config to Spec 3

OpenShift versions 4.6 and above use Ignition Spec 3.  This
command can be used to translate the "pointer config" generated
by a prior version of openshift-install to Spec 3, so that
it can be used with newer bootimages.
`,
		Args: cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			dataIn, err := ioutil.ReadFile(args[0])
			if err != nil {
				return errors.Wrapf(err, "failed to read %s", args[0])
			}
			dest := args[1]
			cfg, rpt, err := ignv2.Parse(dataIn)
			fmt.Fprintf(os.Stderr, "%s", rpt.String())
			if err != nil || rpt.IsFatal() {
				return errors.Errorf("Error parsing spec v2 config: %v\n%v", err, rpt)
			}

			newCfg, err := v24tov31.Translate(cfg, nil)
			if err != nil {
				return errors.Wrapf(err, "translation failed", err)
			}
			dataOut, err := json.Marshal(newCfg)
			if err != nil {
				return errors.Wrapf(err, "failed to marshal JSON")
			}
			return ioutil.WriteFile(dest, dataOut, 0666)
		},
	}
	return cmd
}