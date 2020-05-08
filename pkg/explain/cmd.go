package explain

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCmd returns a subcommand for explain
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "List the fields for supported InstallConfig versions",
		Long: `This command describes the fields associated with each supported InstallConfig API. Fields are identified via a simple
JSONPath identifier:
		
installconfig.<fieldName>[.<fieldName>]
`,
		Example: `
# Get the documentation of the resource and its fields
openshift-install explain installconfig

# Get the documentation of a AWS platform
openshift-install explain installconfig.platform.aws`,
		RunE: runCmd,
	}

	return cmd
}

func runCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.Errorf("You must specify the type of resource to explain\n")
	}
	if len(args) > 1 {
		return errors.Errorf("We accept only this format: explain RESOURCE\n")
	}

	file, err := data.Assets.Open(installConfigCRDFileName)
	if err != nil {
		return errors.Wrap(err, "failed to load InstallConfig CRD")
	}
	defer file.Close()

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read InstallConfig CRD")
	}

	resource, path := splitDotNotation(args[0])
	if resource != "installconfig" {
		return errors.Errorf("only installconfig resource is supported")
	}

	schema, err := loadSchema(raw)
	if err != nil {
		return errors.Wrap(err, "failed to load schema")
	}

	fschema, err := lookup(schema, path)
	if err != nil {
		return errors.Wrapf(err, "failed to load schema for the field %s", strings.Join(path, "."))
	}

	p := printer{Writer: os.Stdout}
	p.PrintKindAndVersion()
	p.PrintResource(fschema)
	p.PrintFields(fschema)
	return nil
}

func splitDotNotation(model string) (string, []string) {
	var fieldsPath []string

	// ignore trailing period
	model = strings.TrimSuffix(model, ".")

	dotModel := strings.Split(model, ".")
	if len(dotModel) >= 1 {
		fieldsPath = dotModel[1:]
	}
	return dotModel[0], fieldsPath
}
