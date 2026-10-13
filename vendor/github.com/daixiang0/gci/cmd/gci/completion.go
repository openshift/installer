package gci

import (
	"strings"

	"github.com/spf13/cobra"
)

func subCommandOrGoFileCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var commandAliases []string
	for _, subCmd := range cmd.Commands() {
		commandAliases = append(commandAliases, subCmd.Name())
		commandAliases = append(commandAliases, subCmd.Aliases...)
	}
	for _, subCmdStr := range commandAliases {
		if strings.HasPrefix(subCmdStr, toComplete) {
			// completion for commands is already provided by cobra
			// do not return file completion
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
	}
	return goFileCompletion(cmd, args, toComplete)
}

func goFileCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"go"}, cobra.ShellCompDirectiveFilterFileExt
}
