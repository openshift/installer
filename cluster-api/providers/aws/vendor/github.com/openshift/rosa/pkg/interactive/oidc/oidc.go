package oidc

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/rosa"
)

func GetOidcConfigID(r *rosa.Runtime, cmd *cobra.Command) string {
	oidcConfigs, err := r.OCMClient.ListOidcConfigs(r.Creator.AccountID)
	if err != nil {
		r.Reporter.Warnf("There was a problem retrieving OIDC Configurations "+
			"for your organization: %v", err)
		return ""
	}
	if len(oidcConfigs) == 0 {
		return ""
	}
	oidcConfigsIds := []string{}
	for _, oidcConfig := range oidcConfigs {
		oidcConfigsIds = append(oidcConfigsIds, fmt.Sprintf("%s | %s", oidcConfig.ID(), oidcConfig.IssuerUrl()))
	}
	oidcConfigId, err := interactive.GetOption(interactive.Input{
		Question: "OIDC Configuration ID",
		Help:     cmd.Flags().Lookup("oidc-config-id").Usage,
		Options:  oidcConfigsIds,
		Default:  oidcConfigsIds[0],
		Required: true,
	})
	if err != nil {
		r.Reporter.Errorf("Expected a valid OIDC Config ID: %s", err)
		os.Exit(1)
	}
	return strings.TrimSpace(strings.Split(oidcConfigId, "|")[0])
}
