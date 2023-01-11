package gci

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/daixiang0/gci/pkg/section"
)

type Executor struct {
	rootCmd    *cobra.Command
	diffMode   *bool
	writeMode  *bool
	localFlags *[]string
}

func NewExecutor(version string) *Executor {
	log.InitLogger()
	defer log.L().Sync()

	e := Executor{}
	rootCmd := cobra.Command{
		Use:   "gci [-diff | -write] [--local localPackageURLs] path...",
		Short: "Gci controls golang package import order and makes it always deterministic",
		Long: "Gci enables automatic formatting of imports in a deterministic manner" +
			"\n" +
			"If you want to integrate this as part of your CI take a look at golangci-lint.",
		ValidArgsFunction: subCommandOrGoFileCompletion,
		Args:              cobra.MinimumNArgs(1),
		Version:           version,
		RunE:              e.runInCompatibilityMode,
	}
	e.rootCmd = &rootCmd
	e.diffMode = rootCmd.Flags().BoolP("diff", "d", false, "display diffs instead of rewriting files")
	e.writeMode = rootCmd.Flags().BoolP("write", "w", false, "write result to (source) file instead of stdout")
	e.localFlags = rootCmd.Flags().StringSliceP("local", "l", []string{}, "put imports beginning with this string after 3rd-party packages, separate imports by comma")
	e.initDiff()
	e.initPrint()
	e.initWrite()
	return &e
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

func (e *Executor) runInCompatibilityMode(cmd *cobra.Command, args []string) error {
	// Workaround since the Deprecation message in Cobra can not be printed to STDERR
	_, _ = fmt.Fprintln(os.Stderr, "Using the old parameters is deprecated, please use the named subcommands!")

	if *e.writeMode && *e.diffMode {
		return fmt.Errorf("diff and write must not be specified at the same time")
	}
	// generate section specification from old localFlags format
	sections := gci.LocalFlagsToSections(*e.localFlags)
	sectionSeparators := section.DefaultSectionSeparators()
	cfg := config.Config{
		BoolConfig: config.BoolConfig{
			NoInlineComments: false,
			NoPrefixComments: false,
			Debug:            false,
			SkipGenerated:    false,
		},
		Sections:          sections,
		SectionSeparators: sectionSeparators,
	}
	if *e.writeMode {
		return gci.WriteFormattedFiles(args, cfg)
	}
	if *e.diffMode {
		return gci.DiffFormattedFiles(args, cfg)
	}
	return gci.PrintFormattedFiles(args, cfg)
}
