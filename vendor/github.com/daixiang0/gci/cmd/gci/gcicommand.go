package gci

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/daixiang0/gci/pkg/section"
)

type processingFunc = func(args []string, gciCfg config.Config) error

func (e *Executor) newGciCommand(use, short, long string, aliases []string, stdInSupport bool, processingFunc processingFunc) *cobra.Command {
	var noInlineComments, noPrefixComments, skipGenerated, skipVendor, customOrder, noLexOrder, debug *bool
	var sectionStrings, sectionSeparatorStrings *[]string
	cmd := cobra.Command{
		Use:               use,
		Aliases:           aliases,
		Short:             short,
		Long:              long,
		ValidArgsFunction: goFileCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmtCfg := config.BoolConfig{
				NoInlineComments: *noInlineComments,
				NoPrefixComments: *noPrefixComments,
				Debug:            *debug,
				SkipGenerated:    *skipGenerated,
				SkipVendor:       *skipVendor,
				CustomOrder:      *customOrder,
				NoLexOrder:       *noLexOrder,
			}
			gciCfg, err := config.YamlConfig{Cfg: fmtCfg, SectionStrings: *sectionStrings, SectionSeparatorStrings: *sectionSeparatorStrings}.Parse()
			if err != nil {
				return err
			}
			if *debug {
				log.SetLevel(zapcore.DebugLevel)
			}
			return processingFunc(args, *gciCfg)
		},
	}
	if !stdInSupport {
		cmd.Args = cobra.MinimumNArgs(1)
	}

	// register command as subcommand
	e.rootCmd.AddCommand(&cmd)

	debug = cmd.Flags().BoolP("debug", "d", false, "Enables debug output from the formatter")

	sectionHelp := `Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). The section order is standard > default > custom > blank > dot > alias > localmodule. The default value is [standard,default].
standard - standard section that Go provides officially, like "fmt"
Prefix(github.com/daixiang0) - custom section, groups all imports with the specified Prefix. Imports will be matched to the longest Prefix. Multiple custom prefixes may be provided, they will be rendered as distinct sections separated by newline. You can regroup multiple prefixes by separating them with comma: Prefix(github.com/daixiang0,gitlab.com/daixiang0,daixiang0)
default - default section, contains all rest imports
blank - blank section, contains all blank imports.
dot - dot section, contains all dot imports.
alias - alias section, contains all alias imports.
localmodule: localmodule section, contains all imports from local packages`

	skipGenerated = cmd.Flags().Bool("skip-generated", false, "Skip generated files")
	skipVendor = cmd.Flags().Bool("skip-vendor", false, "Skip files inside vendor directory")

	customOrder = cmd.Flags().Bool("custom-order", false, "Enable custom order of sections")
	noLexOrder = cmd.Flags().Bool("no-lex-order", false, "Drops lexical ordering for custom sections")
	sectionStrings = cmd.Flags().StringArrayP("section", "s", section.DefaultSections().String(), sectionHelp)

	// deprecated
	noInlineComments = cmd.Flags().Bool("NoInlineComments", false, "Drops inline comments while formatting")
	cmd.Flags().MarkDeprecated("NoInlineComments", "Drops inline comments while formatting")
	noPrefixComments = cmd.Flags().Bool("NoPrefixComments", false, "Drops comment lines above an import statement while formatting")
	cmd.Flags().MarkDeprecated("NoPrefixComments", "Drops inline comments while formatting")
	sectionSeparatorStrings = cmd.Flags().StringSliceP("SectionSeparator", "x", section.DefaultSectionSeparators().String(), "SectionSeparators are inserted between Sections")
	cmd.Flags().MarkDeprecated("SectionSeparator", "Drops inline comments while formatting")
	cmd.Flags().MarkDeprecated("x", "Drops inline comments while formatting")

	return &cmd
}
