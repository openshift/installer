package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/posener/complete"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	backendInit "github.com/hashicorp/terraform/backend/init"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/internal/getproviders"
	"github.com/hashicorp/terraform/internal/providercache"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
	tfversion "github.com/hashicorp/terraform/version"
)

// InitCommand is a Command implementation that takes a Terraform
// module and clones it to the working directory.
type InitCommand struct {
	Meta
}

func (c *InitCommand) Run(args []string) int {
	var flagFromModule string
	var flagBackend, flagGet, flagUpgrade, getPlugins bool
	var flagPluginPath FlagStringSlice
	var flagVerifyPlugins bool
	flagConfigExtra := newRawFlags("-backend-config")

	args = c.Meta.process(args)
	cmdFlags := c.Meta.extendedFlagSet("init")
	cmdFlags.BoolVar(&flagBackend, "backend", true, "")
	cmdFlags.Var(flagConfigExtra, "backend-config", "")
	cmdFlags.StringVar(&flagFromModule, "from-module", "", "copy the source of the given module into the directory before init")
	cmdFlags.BoolVar(&flagGet, "get", true, "")
	cmdFlags.BoolVar(&getPlugins, "get-plugins", true, "no-op flag, use provider_installation blocks to customize provider installation")
	cmdFlags.BoolVar(&c.forceInitCopy, "force-copy", false, "suppress prompts about copying state data")
	cmdFlags.BoolVar(&c.Meta.stateLock, "lock", true, "lock state")
	cmdFlags.DurationVar(&c.Meta.stateLockTimeout, "lock-timeout", 0, "lock timeout")
	cmdFlags.BoolVar(&c.reconfigure, "reconfigure", false, "reconfigure")
	cmdFlags.BoolVar(&flagUpgrade, "upgrade", false, "")
	cmdFlags.Var(&flagPluginPath, "plugin-dir", "plugin directory")
	cmdFlags.BoolVar(&flagVerifyPlugins, "verify-plugins", true, "verify plugins")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	var diags tfdiags.Diagnostics

	if len(flagPluginPath) > 0 {
		c.pluginPath = flagPluginPath
	}

	// If users are setting the no-op get-plugins command, give them a warning,
	// this should allow us to remove the flag in the future.
	if !getPlugins {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Warning,
			"No-op -get-plugins flag used",
			`As of Terraform 0.13+, the -get-plugins=false command is a no-op flag. If you would like to customize provider installation, use a provider_installation block or other available Terraform settings.`,
		))
	}

	// Validate the arg count
	args = cmdFlags.Args()
	if len(args) > 1 {
		c.Ui.Error("The init command expects at most one argument.\n")
		cmdFlags.Usage()
		return 1
	}

	if err := c.storePluginPath(c.pluginPath); err != nil {
		c.Ui.Error(fmt.Sprintf("Error saving -plugin-path values: %s", err))
		return 1
	}

	// Get our pwd. We don't always need it but always getting it is easier
	// than the logic to determine if it is or isn't needed.
	pwd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting pwd: %s", err))
		return 1
	}

	// If an argument is provided then it overrides our working directory.
	path := pwd
	if len(args) == 1 {
		path = args[0]
	}

	// This will track whether we outputted anything so that we know whether
	// to output a newline before the success message
	var header bool

	if flagFromModule != "" {
		src := flagFromModule

		empty, err := configs.IsEmptyDir(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error validating destination directory: %s", err))
			return 1
		}
		if !empty {
			c.Ui.Error(strings.TrimSpace(errInitCopyNotEmpty))
			return 1
		}

		c.Ui.Output(c.Colorize().Color(fmt.Sprintf(
			"[reset][bold]Copying configuration[reset] from %q...", src,
		)))
		header = true

		hooks := uiModuleInstallHooks{
			Ui:             c.Ui,
			ShowLocalPaths: false, // since they are in a weird location for init
		}

		initDiags := c.initDirFromModule(path, src, hooks)
		diags = diags.Append(initDiags)
		if initDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}

		c.Ui.Output("")
	}

	// If our directory is empty, then we're done. We can't get or setup
	// the backend with an empty directory.
	empty, err := configs.IsEmptyDir(path)
	if err != nil {
		diags = diags.Append(fmt.Errorf("Error checking configuration: %s", err))
		c.showDiagnostics(diags)
		return 1
	}
	if empty {
		c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitEmpty)))
		return 0
	}

	// For Terraform v0.12 we introduced a special loading mode where we would
	// use the 0.11-syntax-compatible "earlyconfig" package as a heuristic to
	// identify situations where it was likely that the user was trying to use
	// 0.11-only syntax that the upgrade tool might help with.
	//
	// However, as the language has moved on that is no longer a suitable
	// heuristic in Terraform 0.13 and later: other new additions to the
	// language can cause the main loader to disagree with earlyconfig, which
	// would lead us to give poor advice about how to respond.
	//
	// For that reason, we no longer use a different error message in that
	// situation, but for now we still use both codepaths because some of our
	// initialization functionality remains built around "earlyconfig" and
	// so we need to still load the module via that mechanism anyway until we
	// can do some more invasive refactoring here.
	rootMod, confDiags := c.loadSingleModule(path)
	rootModEarly, earlyConfDiags := c.loadSingleModuleEarly(path)
	if confDiags.HasErrors() {
		c.Ui.Error(c.Colorize().Color(strings.TrimSpace(errInitConfigError)))
		// TODO: It would be nice to check the version constraints in
		// rootModEarly.RequiredCore and print out a hint if the module is
		// declaring that it's not compatible with this version of Terraform,
		// though we're deferring that for now because we're intending to
		// refactor our use of "earlyconfig" here anyway and so whatever we
		// might do here right now would likely be invalidated by that.
		c.showDiagnostics(confDiags)
		return 1
	}
	// If _only_ the early loader encountered errors then that's unusual
	// (it should generally be a superset of the normal loader) but we'll
	// return those errors anyway since otherwise we'll probably get
	// some weird behavior downstream. Errors from the early loader are
	// generally not as high-quality since it has less context to work with.
	if earlyConfDiags.HasErrors() {
		c.Ui.Error(c.Colorize().Color(strings.TrimSpace(errInitConfigError)))
		// Errors from the early loader are generally not as high-quality since
		// it has less context to work with.
		diags = diags.Append(confDiags)
		c.showDiagnostics(diags)
		return 1
	}

	if flagGet {
		modsOutput, modsDiags := c.getModules(path, rootModEarly, flagUpgrade)
		diags = diags.Append(modsDiags)
		if modsDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		if modsOutput {
			header = true
		}
	}

	// With all of the modules (hopefully) installed, we can now try to load the
	// whole configuration tree.
	config, confDiags := c.loadConfig(path)
	diags = diags.Append(confDiags)
	if confDiags.HasErrors() {
		c.Ui.Error(strings.TrimSpace(errInitConfigError))
		c.showDiagnostics(diags)
		return 1
	}

	// Before we go further, we'll check to make sure none of the modules in the
	// configuration declare that they don't support this Terraform version, so
	// we can produce a version-related error message rather than
	// potentially-confusing downstream errors.
	versionDiags := terraform.CheckCoreVersionRequirements(config)
	diags = diags.Append(versionDiags)
	if versionDiags.HasErrors() {
		c.showDiagnostics(diags)
		return 1
	}

	var back backend.Backend
	if flagBackend {

		be, backendOutput, backendDiags := c.initBackend(rootMod, flagConfigExtra)
		diags = diags.Append(backendDiags)
		if backendDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		if backendOutput {
			header = true
		}
		back = be
	} else {
		// load the previously-stored backend config
		be, backendDiags := c.Meta.backendFromState()
		diags = diags.Append(backendDiags)
		if backendDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		back = be
	}

	if back == nil {
		// If we didn't initialize a backend then we'll try to at least
		// instantiate one. This might fail if it wasn't already initialized
		// by a previous run, so we must still expect that "back" may be nil
		// in code that follows.
		var backDiags tfdiags.Diagnostics
		back, backDiags = c.Backend(nil)
		if backDiags.HasErrors() {
			// This is fine. We'll proceed with no backend, then.
			back = nil
		}
	}

	var state *states.State

	// If we have a functional backend (either just initialized or initialized
	// on a previous run) we'll use the current state as a potential source
	// of provider dependencies.
	if back != nil {
		c.ignoreRemoteBackendVersionConflict(back)
		workspace, err := c.Workspace()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error selecting workspace: %s", err))
			return 1
		}
		sMgr, err := back.StateMgr(workspace)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error loading state: %s", err))
			return 1
		}

		if err := sMgr.RefreshState(); err != nil {
			c.Ui.Error(fmt.Sprintf("Error refreshing state: %s", err))
			return 1
		}

		state = sMgr.State()
	}

	// Now that we have loaded all modules, check the module tree for missing providers.
	providersOutput, providersAbort, providerDiags := c.getProviders(config, state, flagUpgrade, flagPluginPath)
	diags = diags.Append(providerDiags)
	if providersAbort || providerDiags.HasErrors() {
		c.showDiagnostics(diags)
		return 1
	}
	if providersOutput {
		header = true
	}

	// If we outputted information, then we need to output a newline
	// so that our success message is nicely spaced out from prior text.
	if header {
		c.Ui.Output("")
	}

	// If we accumulated any warnings along the way that weren't accompanied
	// by errors then we'll output them here so that the success message is
	// still the final thing shown.
	c.showDiagnostics(diags)
	c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitSuccess)))
	if !c.RunningInAutomation {
		// If we're not running in an automation wrapper, give the user
		// some more detailed next steps that are appropriate for interactive
		// shell usage.
		c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitSuccessCLI)))
	}
	return 0
}

func (c *InitCommand) getModules(path string, earlyRoot *tfconfig.Module, upgrade bool) (output bool, diags tfdiags.Diagnostics) {
	if len(earlyRoot.ModuleCalls) == 0 {
		// Nothing to do
		return false, nil
	}

	if upgrade {
		c.Ui.Output(c.Colorize().Color(fmt.Sprintf("[reset][bold]Upgrading modules...")))
	} else {
		c.Ui.Output(c.Colorize().Color(fmt.Sprintf("[reset][bold]Initializing modules...")))
	}

	hooks := uiModuleInstallHooks{
		Ui:             c.Ui,
		ShowLocalPaths: true,
	}
	instDiags := c.installModules(path, upgrade, hooks)
	diags = diags.Append(instDiags)

	// Since module installer has modified the module manifest on disk, we need
	// to refresh the cache of it in the loader.
	if c.configLoader != nil {
		if err := c.configLoader.RefreshModules(); err != nil {
			// Should never happen
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to read module manifest",
				fmt.Sprintf("After installing modules, Terraform could not re-read the manifest of installed modules. This is a bug in Terraform. %s.", err),
			))
		}
	}

	return true, diags
}

func (c *InitCommand) initBackend(root *configs.Module, extraConfig rawFlags) (be backend.Backend, output bool, diags tfdiags.Diagnostics) {
	c.Ui.Output(c.Colorize().Color(fmt.Sprintf("\n[reset][bold]Initializing the backend...")))

	var backendConfig *configs.Backend
	var backendConfigOverride hcl.Body
	if root.Backend != nil {
		backendType := root.Backend.Type
		bf := backendInit.Backend(backendType)
		if bf == nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported backend type",
				Detail:   fmt.Sprintf("There is no backend type named %q.", backendType),
				Subject:  &root.Backend.TypeRange,
			})
			return nil, true, diags
		}

		b := bf()
		backendSchema := b.ConfigSchema()
		backendConfig = root.Backend

		var overrideDiags tfdiags.Diagnostics
		backendConfigOverride, overrideDiags = c.backendConfigOverrideBody(extraConfig, backendSchema)
		diags = diags.Append(overrideDiags)
		if overrideDiags.HasErrors() {
			return nil, true, diags
		}
	} else {
		// If the user supplied a -backend-config on the CLI but no backend
		// block was found in the configuration, it's likely - but not
		// necessarily - a mistake. Return a warning.
		if !extraConfig.Empty() {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Warning,
				"Missing backend configuration",
				`-backend-config was used without a "backend" block in the configuration.

If you intended to override the default local backend configuration,
no action is required, but you may add an explicit backend block to your
configuration to clear this warning:

terraform {
  backend "local" {}
}

However, if you intended to override a defined backend, please verify that
the backend configuration is present and valid.
`,
			))
		}
	}

	opts := &BackendOpts{
		Config:         backendConfig,
		ConfigOverride: backendConfigOverride,
		Init:           true,
	}

	back, backDiags := c.Backend(opts)
	diags = diags.Append(backDiags)
	return back, true, diags
}

// Load the complete module tree, and fetch any missing providers.
// This method outputs its own Ui.
func (c *InitCommand) getProviders(config *configs.Config, state *states.State, upgrade bool, pluginDirs []string) (output, abort bool, diags tfdiags.Diagnostics) {
	// Dev overrides cause the result of "terraform init" to be irrelevant for
	// any overridden providers, so we'll warn about it to avoid later
	// confusion when Terraform ends up using a different provider than the
	// lock file called for.
	diags = diags.Append(c.providerDevOverrideInitWarnings())

	// First we'll collect all the provider dependencies we can see in the
	// configuration and the state.
	reqs, hclDiags := config.ProviderRequirements()
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return false, true, diags
	}
	if state != nil {
		stateReqs := state.ProviderRequirements()
		reqs = reqs.Merge(stateReqs)
	}

	for providerAddr := range reqs {
		if providerAddr.IsLegacy() {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid legacy provider address",
				fmt.Sprintf(
					"This configuration or its associated state refers to the unqualified provider %q.\n\nYou must complete the Terraform 0.13 upgrade process before upgrading to later versions.",
					providerAddr.Type,
				),
			))
		}
	}

	previousLocks, moreDiags := c.lockedDependencies()
	diags = diags.Append(moreDiags)

	if diags.HasErrors() {
		return false, true, diags
	}

	var inst *providercache.Installer
	if len(pluginDirs) == 0 {
		// By default we use a source that looks for providers in all of the
		// standard locations, possibly customized by the user in CLI config.
		inst = c.providerInstaller()
	} else {
		// If the user passes at least one -plugin-dir then that circumvents
		// the usual sources and forces Terraform to consult only the given
		// directories. Anything not available in one of those directories
		// is not available for installation.
		source := c.providerCustomLocalDirectorySource(pluginDirs)
		inst = c.providerInstallerCustomSource(source)

		// The default (or configured) search paths are logged earlier, in provider_source.go
		// Log that those are being overridden by the `-plugin-dir` command line options
		log.Println("[DEBUG] init: overriding provider plugin search paths")
		log.Printf("[DEBUG] will search for provider plugins in %s", pluginDirs)
	}

	// Installation can be aborted by interruption signals
	ctx, done := c.InterruptibleContext()
	defer done()

	// Because we're currently just streaming a series of events sequentially
	// into the terminal, we're showing only a subset of the events to keep
	// things relatively concise. Later it'd be nice to have a progress UI
	// where statuses update in-place, but we can't do that as long as we
	// are shimming our vt100 output to the legacy console API on Windows.
	missingProviders := make(map[addrs.Provider]struct{})
	evts := &providercache.InstallerEvents{
		PendingProviders: func(reqs map[addrs.Provider]getproviders.VersionConstraints) {
			c.Ui.Output(c.Colorize().Color(
				"\n[reset][bold]Initializing provider plugins...",
			))
		},
		ProviderAlreadyInstalled: func(provider addrs.Provider, selectedVersion getproviders.Version) {
			c.Ui.Info(fmt.Sprintf("- Using previously-installed %s v%s", provider.ForDisplay(), selectedVersion))
		},
		BuiltInProviderAvailable: func(provider addrs.Provider) {
			c.Ui.Info(fmt.Sprintf("- %s is built in to Terraform", provider.ForDisplay()))
		},
		BuiltInProviderFailure: func(provider addrs.Provider, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid dependency on built-in provider",
				fmt.Sprintf("Cannot use %s: %s.", provider.ForDisplay(), err),
			))
		},
		QueryPackagesBegin: func(provider addrs.Provider, versionConstraints getproviders.VersionConstraints, locked bool) {
			if locked {
				c.Ui.Info(fmt.Sprintf("- Reusing previous version of %s from the dependency lock file", provider.ForDisplay()))
			} else {
				if len(versionConstraints) > 0 {
					c.Ui.Info(fmt.Sprintf("- Finding %s versions matching %q...", provider.ForDisplay(), getproviders.VersionConstraintsString(versionConstraints)))
				} else {
					c.Ui.Info(fmt.Sprintf("- Finding latest version of %s...", provider.ForDisplay()))
				}
			}
		},
		LinkFromCacheBegin: func(provider addrs.Provider, version getproviders.Version, cacheRoot string) {
			c.Ui.Info(fmt.Sprintf("- Using %s v%s from the shared cache directory", provider.ForDisplay(), version))
		},
		FetchPackageBegin: func(provider addrs.Provider, version getproviders.Version, location getproviders.PackageLocation) {
			c.Ui.Info(fmt.Sprintf("- Installing %s v%s...", provider.ForDisplay(), version))
		},
		QueryPackagesFailure: func(provider addrs.Provider, err error) {
			// We track providers that had missing metadata because we might
			// generate additional hints for some of them at the end.
			missingProviders[provider] = struct{}{}

			switch errorTy := err.(type) {
			case getproviders.ErrProviderNotFound:
				sources := errorTy.Sources
				displaySources := make([]string, len(sources))
				for i, source := range sources {
					displaySources[i] = fmt.Sprintf("  - %s", source)
				}
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Failed to query available provider packages",
					fmt.Sprintf("Could not retrieve the list of available versions for provider %s: %s\n\n%s",
						provider.ForDisplay(), err, strings.Join(displaySources, "\n"),
					),
				))
			case getproviders.ErrRegistryProviderNotKnown:
				// We might be able to suggest an alternative provider to use
				// instead of this one.
				var suggestion string
				alternative := getproviders.MissingProviderSuggestion(ctx, provider, inst.ProviderSource())
				if alternative != provider {
					suggestion = fmt.Sprintf(
						"\n\nDid you intend to use %s? If so, you must specify that source address in each module which requires that provider. To see which modules are currently depending on %s, run the following command:\n    terraform providers",
						alternative.ForDisplay(), provider.ForDisplay(),
					)
				}

				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Failed to query available provider packages",
					fmt.Sprintf(
						"Could not retrieve the list of available versions for provider %s: %s\n\nIf you have just upgraded directly from Terraform v0.12 to Terraform v0.14 then please upgrade to Terraform v0.13 first and follow the upgrade guide for that release, which might help you address this problem.%s",
						provider.ForDisplay(), err, suggestion,
					),
				))
			case getproviders.ErrHostNoProviders:
				switch {
				case errorTy.Hostname == svchost.Hostname("github.com") && !errorTy.HasOtherVersion:
					// If a user copies the URL of a GitHub repository into
					// the source argument and removes the schema to make it
					// provider-address-shaped then that's one way we can end up
					// here. We'll use a specialized error message in anticipation
					// of that mistake. We only do this if github.com isn't a
					// provider registry, to allow for the (admittedly currently
					// rather unlikely) possibility that github.com starts being
					// a real Terraform provider registry in the future.
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider registry host",
						fmt.Sprintf("The given source address %q specifies a GitHub repository rather than a Terraform provider. Refer to the documentation of the provider to find the correct source address to use.",
							provider.String(),
						),
					))

				case errorTy.HasOtherVersion:
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider registry host",
						fmt.Sprintf("The host %q given in in provider source address %q does not offer a Terraform provider registry that is compatible with this Terraform version, but it may be compatible with a different Terraform version.",
							errorTy.Hostname, provider.String(),
						),
					))

				default:
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider registry host",
						fmt.Sprintf("The host %q given in in provider source address %q does not offer a Terraform provider registry.",
							errorTy.Hostname, provider.String(),
						),
					))
				}

			case getproviders.ErrRequestCanceled:
				// We don't attribute cancellation to any particular operation,
				// but rather just emit a single general message about it at
				// the end, by checking ctx.Err().

			default:
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Failed to query available provider packages",
					fmt.Sprintf("Could not retrieve the list of available versions for provider %s: %s",
						provider.ForDisplay(), err,
					),
				))
			}

		},
		QueryPackagesWarning: func(provider addrs.Provider, warnings []string) {
			displayWarnings := make([]string, len(warnings))
			for i, warning := range warnings {
				displayWarnings[i] = fmt.Sprintf("- %s", warning)
			}

			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Warning,
				"Additional provider information from registry",
				fmt.Sprintf("The remote registry returned warnings for %s:\n%s",
					provider.String(),
					strings.Join(displayWarnings, "\n"),
				),
			))
		},
		LinkFromCacheFailure: func(provider addrs.Provider, version getproviders.Version, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to install provider from shared cache",
				fmt.Sprintf("Error while importing %s v%s from the shared cache directory: %s.", provider.ForDisplay(), version, err),
			))
		},
		FetchPackageFailure: func(provider addrs.Provider, version getproviders.Version, err error) {
			const summaryIncompatible = "Incompatible provider version"
			switch err := err.(type) {
			case getproviders.ErrProtocolNotSupported:
				closestAvailable := err.Suggestion
				switch {
				case closestAvailable == getproviders.UnspecifiedVersion:
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						summaryIncompatible,
						fmt.Sprintf(errProviderVersionIncompatible, provider.String()),
					))
				case version.GreaterThan(closestAvailable):
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						summaryIncompatible,
						fmt.Sprintf(providerProtocolTooNew, provider.ForDisplay(),
							version, tfversion.String(), closestAvailable, closestAvailable,
							getproviders.VersionConstraintsString(reqs[provider]),
						),
					))
				default: // version is less than closestAvailable
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						summaryIncompatible,
						fmt.Sprintf(providerProtocolTooOld, provider.ForDisplay(),
							version, tfversion.String(), closestAvailable, closestAvailable,
							getproviders.VersionConstraintsString(reqs[provider]),
						),
					))
				}
			case getproviders.ErrPlatformNotSupported:
				switch {
				case err.MirrorURL != nil:
					// If we're installing from a mirror then it may just be
					// the mirror lacking the package, rather than it being
					// unavailable from upstream.
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						summaryIncompatible,
						fmt.Sprintf(
							"Your chosen provider mirror at %s does not have a %s v%s package available for your current platform, %s.\n\nProvider releases are separate from Terraform CLI releases, so this provider might not support your current platform. Alternatively, the mirror itself might have only a subset of the plugin packages available in the origin registry, at %s.",
							err.MirrorURL, err.Provider, err.Version, err.Platform,
							err.Provider.Hostname,
						),
					))
				default:
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						summaryIncompatible,
						fmt.Sprintf(
							"Provider %s v%s does not have a package available for your current platform, %s.\n\nProvider releases are separate from Terraform CLI releases, so not all providers are available for all platforms. Other versions of this provider may have different platforms supported.",
							err.Provider, err.Version, err.Platform,
						),
					))
				}

			case getproviders.ErrRequestCanceled:
				// We don't attribute cancellation to any particular operation,
				// but rather just emit a single general message about it at
				// the end, by checking ctx.Err().

			default:
				// We can potentially end up in here under cancellation too,
				// in spite of our getproviders.ErrRequestCanceled case above,
				// because not all of the outgoing requests we do under the
				// "fetch package" banner are source metadata requests.
				// In that case we will emit a redundant error here about
				// the request being cancelled, but we'll still detect it
				// as a cancellation after the installer returns and do the
				// normal cancellation handling.

				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Failed to install provider",
					fmt.Sprintf("Error while installing %s v%s: %s", provider.ForDisplay(), version, err),
				))
			}
		},
		FetchPackageSuccess: func(provider addrs.Provider, version getproviders.Version, localDir string, authResult *getproviders.PackageAuthenticationResult) {
			var keyID string
			if authResult != nil && authResult.ThirdPartySigned() {
				keyID = authResult.KeyID
			}
			if keyID != "" {
				keyID = c.Colorize().Color(fmt.Sprintf(", key ID [reset][bold]%s[reset]", keyID))
			}

			c.Ui.Info(fmt.Sprintf("- Installed %s v%s (%s%s)", provider.ForDisplay(), version, authResult, keyID))
		},
		ProvidersFetched: func(authResults map[addrs.Provider]*getproviders.PackageAuthenticationResult) {
			thirdPartySigned := false
			for _, authResult := range authResults {
				if authResult.ThirdPartySigned() {
					thirdPartySigned = true
					break
				}
			}
			if thirdPartySigned {
				c.Ui.Info(fmt.Sprintf("\nPartner and community providers are signed by their developers.\n" +
					"If you'd like to know more about provider signing, you can read about it here:\n" +
					"https://www.terraform.io/docs/cli/plugins/signing.html"))
			}
		},
		HashPackageFailure: func(provider addrs.Provider, version getproviders.Version, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to validate installed provider",
				fmt.Sprintf(
					"Validating provider %s v%s failed: %s",
					provider.ForDisplay(),
					version,
					err,
				),
			))
		},
	}
	ctx = evts.OnContext(ctx)

	mode := providercache.InstallNewProvidersOnly
	if upgrade {
		mode = providercache.InstallUpgrades
	}
	newLocks, err := inst.EnsureProviderVersions(ctx, previousLocks, reqs, mode)
	if ctx.Err() == context.Canceled {
		c.showDiagnostics(diags)
		c.Ui.Error("Provider installation was canceled by an interrupt signal.")
		return true, true, diags
	}
	if len(missingProviders) > 0 {
		// If we encountered requirements for one or more providers where we
		// weren't able to find any metadata, that _might_ be because a
		// user had previously (before 0.14) been incorrectly using the
		// .terraform/plugins directory as if it were a local filesystem
		// mirror, rather than as the main cache directory.
		//
		// We no longer allow that because it'd be ambiguous whether plugins in
		// there are explictly intended to be a local mirror or if they are
		// just leftover cache entries from provider installation in
		// Terraform 0.13.
		//
		// To help those users migrate we have a specialized warning message
		// for it, which we'll produce only if one of the missing providers can
		// be seen in the "legacy" cache directory, which is what we're now
		// considering .terraform/plugins to be. (The _current_ cache directory
		// is .terraform/providers.)
		//
		// This is only a heuristic, so it might potentially produce false
		// positives if a user happens to encounter another sort of error
		// while they are upgrading from Terraform 0.13 to 0.14. Aside from
		// upgrading users should not end up in here because they won't
		// have a legacy cache directory at all.
		legacyDir := c.providerLegacyCacheDir()
		if legacyDir != nil { // if the legacy directory is present at all
			for missingProvider := range missingProviders {
				if missingProvider.IsDefault() {
					// If we get here for a default provider then it's more
					// likely that something _else_ went wrong, like a network
					// problem, so we'll skip the warning in this case to
					// avoid potentially misleading the user into creating an
					// unnecessary local mirror for an official provider.
					continue
				}
				entry := legacyDir.ProviderLatestVersion(missingProvider)
				if entry == nil {
					continue
				}
				// If we get here then the missing provider was cached, which
				// implies that it might be an in-house provider the user
				// placed manually to try to make Terraform use it as if it
				// were a local mirror directory.
				wantDir := filepath.FromSlash(fmt.Sprintf("terraform.d/plugins/%s/%s/%s", missingProvider, entry.Version, getproviders.CurrentPlatform))
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Warning,
					"Missing provider is in legacy cache directory",
					fmt.Sprintf(
						"Terraform supports a number of local directories that can serve as automatic local filesystem mirrors, but .terraform/plugins is not one of them because Terraform v0.13 and earlier used this directory to cache copies of provider plugins retrieved from elsewhere.\n\nIf you intended to use this directory as a filesystem mirror for %s, place it instead in the following directory:\n  %s",
						missingProvider, wantDir,
					),
				))
			}
		}
	}
	if err != nil {
		// The errors captured in "err" should be redundant with what we
		// received via the InstallerEvents callbacks above, so we'll
		// just return those as long as we have some.
		if !diags.HasErrors() {
			diags = diags.Append(err)
		}

		return true, true, diags
	}

	// If the provider dependencies have changed since the last run then we'll
	// say a little about that in case the reader wasn't expecting a change.
	// (When we later integrate module dependencies into the lock file we'll
	// probably want to refactor this so that we produce one lock-file related
	// message for all changes together, but this is here for now just because
	// it's the smallest change relative to what came before it, which was
	// a hidden JSON file specifically for tracking providers.)
	if !newLocks.Equal(previousLocks) {
		if previousLocks.Empty() {
			// A change from empty to non-empty is special because it suggests
			// we're running "terraform init" for the first time against a
			// new configuration. In that case we'll take the opportunity to
			// say a little about what the dependency lock file is, for new
			// users or those who are upgrading from a previous Terraform
			// version that didn't have dependency lock files.
			c.Ui.Output(c.Colorize().Color(`
Terraform has created a lock file [bold].terraform.lock.hcl[reset] to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.`))
		} else {
			c.Ui.Output(c.Colorize().Color(`
Terraform has made some changes to the provider dependency selections recorded
in the .terraform.lock.hcl file. Review those changes and commit them to your
version control system if they represent changes you intended to make.`))
		}
	}

	// TODO: Check whether newLocks is different from previousLocks and mention
	// in the UI if so. We should emit a different message if previousLocks was
	// empty, because that indicates we were creating a lock file for the first
	// time and so we need to introduce the user to the idea of it.

	moreDiags = c.replaceLockedDependencies(newLocks)
	diags = diags.Append(moreDiags)

	return true, false, diags
}

// backendConfigOverrideBody interprets the raw values of -backend-config
// arguments into a hcl Body that should override the backend settings given
// in the configuration.
//
// If the result is nil then no override needs to be provided.
//
// If the returned diagnostics contains errors then the returned body may be
// incomplete or invalid.
func (c *InitCommand) backendConfigOverrideBody(flags rawFlags, schema *configschema.Block) (hcl.Body, tfdiags.Diagnostics) {
	items := flags.AllItems()
	if len(items) == 0 {
		return nil, nil
	}

	var ret hcl.Body
	var diags tfdiags.Diagnostics
	synthVals := make(map[string]cty.Value)

	mergeBody := func(newBody hcl.Body) {
		if ret == nil {
			ret = newBody
		} else {
			ret = configs.MergeBodies(ret, newBody)
		}
	}
	flushVals := func() {
		if len(synthVals) == 0 {
			return
		}
		newBody := configs.SynthBody("-backend-config=...", synthVals)
		mergeBody(newBody)
		synthVals = make(map[string]cty.Value)
	}

	if len(items) == 1 && items[0].Value == "" {
		// Explicitly remove all -backend-config options.
		// We do this by setting an empty but non-nil ConfigOverrides.
		return configs.SynthBody("-backend-config=''", synthVals), diags
	}

	for _, item := range items {
		eq := strings.Index(item.Value, "=")

		if eq == -1 {
			// The value is interpreted as a filename.
			newBody, fileDiags := c.loadHCLFile(item.Value)
			diags = diags.Append(fileDiags)
			if fileDiags.HasErrors() {
				continue
			}
			// Generate an HCL body schema for the backend block.
			var bodySchema hcl.BodySchema
			for name := range schema.Attributes {
				// We intentionally ignore the `Required` attribute here
				// because backend config override files can be partial. The
				// goal is to make sure we're not loading a file with
				// extraneous attributes or blocks.
				bodySchema.Attributes = append(bodySchema.Attributes, hcl.AttributeSchema{
					Name: name,
				})
			}
			for name, block := range schema.BlockTypes {
				var labelNames []string
				if block.Nesting == configschema.NestingMap {
					labelNames = append(labelNames, "key")
				}
				bodySchema.Blocks = append(bodySchema.Blocks, hcl.BlockHeaderSchema{
					Type:       name,
					LabelNames: labelNames,
				})
			}
			// Verify that the file body matches the expected backend schema.
			_, schemaDiags := newBody.Content(&bodySchema)
			diags = diags.Append(schemaDiags)
			if schemaDiags.HasErrors() {
				continue
			}
			flushVals() // deal with any accumulated individual values first
			mergeBody(newBody)
		} else {
			name := item.Value[:eq]
			rawValue := item.Value[eq+1:]
			attrS := schema.Attributes[name]
			if attrS == nil {
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid backend configuration argument",
					fmt.Sprintf("The backend configuration argument %q given on the command line is not expected for the selected backend type.", name),
				))
				continue
			}
			value, valueDiags := configValueFromCLI(item.String(), rawValue, attrS.Type)
			diags = diags.Append(valueDiags)
			if valueDiags.HasErrors() {
				continue
			}
			synthVals[name] = value
		}
	}

	flushVals()

	return ret, diags
}

func (c *InitCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictDirs("")
}

func (c *InitCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-backend":        completePredictBoolean,
		"-backend-config": complete.PredictFiles("*.tfvars"), // can also be key=value, but we can't "predict" that
		"-force-copy":     complete.PredictNothing,
		"-from-module":    completePredictModuleSource,
		"-get":            completePredictBoolean,
		"-input":          completePredictBoolean,
		"-lock":           completePredictBoolean,
		"-lock-timeout":   complete.PredictAnything,
		"-no-color":       complete.PredictNothing,
		"-plugin-dir":     complete.PredictDirs(""),
		"-reconfigure":    complete.PredictNothing,
		"-upgrade":        completePredictBoolean,
		"-verify-plugins": completePredictBoolean,
	}
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: terraform init [options] [DIR]

  Initialize a new or existing Terraform working directory by creating
  initial files, loading any remote state, downloading modules, etc.

  This is the first command that should be run for any new or existing
  Terraform configuration per machine. This sets up all the local data
  necessary to run Terraform that is typically not committed to version
  control.

  This command is always safe to run multiple times. Though subsequent runs
  may give errors, this command will never delete your configuration or
  state. Even so, if you have important information, please back it up prior
  to running this command, just in case.

  If no arguments are given, the configuration in this working directory
  is initialized.

Options:

  -backend=true        Configure the backend for this configuration.

  -backend-config=path This can be either a path to an HCL file with key/value
                       assignments (same format as terraform.tfvars) or a
                       'key=value' format. This is merged with what is in the
                       configuration file. This can be specified multiple
                       times. The backend type must be in the configuration
                       itself.

  -force-copy          Suppress prompts about copying state data. This is
                       equivalent to providing a "yes" to all confirmation
                       prompts.

  -from-module=SOURCE  Copy the contents of the given module into the target
                       directory before initialization.

  -get=true            Download any modules for this configuration.

  -get-plugins=true    Download any missing plugins for this configuration.
                       This command is a no-op in Terraform 0.13+: use
                       -plugin-dir settings or provider_installation blocks
                       instead.

  -input=true          Ask for input if necessary. If false, will error if
                       input was required.

  -lock=true           Lock the state file when locking is supported.

  -lock-timeout=0s     Duration to retry a state lock.

  -no-color            If specified, output won't contain any color.

  -plugin-dir          Directory containing plugin binaries. This overrides all
                       default search paths for plugins, and prevents the
                       automatic installation of plugins. This flag can be used
                       multiple times.

  -reconfigure         Reconfigure the backend, ignoring any saved
                       configuration.

  -upgrade=false       If installing modules (-get) or plugins, ignore
                       previously-downloaded objects and install the
                       latest version allowed within configured constraints.

  -verify-plugins=true Verify the authenticity and integrity of automatically
                       downloaded plugins.
`
	return strings.TrimSpace(helpText)
}

func (c *InitCommand) Synopsis() string {
	return "Prepare your working directory for other commands"
}

const errInitConfigError = `
[reset]There are some problems with the configuration, described below.

The Terraform configuration must be valid before initialization so that
Terraform can determine which modules and providers need to be installed.
`

const errInitCopyNotEmpty = `
The working directory already contains files. The -from-module option requires
an empty directory into which a copy of the referenced module will be placed.

To initialize the configuration already in this working directory, omit the
-from-module option.
`

const outputInitEmpty = `
[reset][bold]Terraform initialized in an empty directory![reset]

The directory has no Terraform configuration files. You may begin working
with Terraform immediately by creating Terraform configuration files.
`

const outputInitSuccess = `
[reset][bold][green]Terraform has been successfully initialized![reset][green]
`

const outputInitSuccessCLI = `[reset][green]
You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
`

const outputInitProvidersUnconstrained = `
The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, we recommend adding version constraints in a required_providers block
in your configuration, with the constraint strings suggested below.
`

// providerProtocolTooOld is a message sent to the CLI UI if the provider's
// supported protocol versions are too old for the user's version of terraform,
// but a newer version of the provider is compatible.
const providerProtocolTooOld = `Provider %q v%s is not compatible with Terraform %s.
Provider version %s is the latest compatible version. Select it with the following version constraint:
	version = %q

Terraform checked all of the plugin versions matching the given constraint:
	%s

Consult the documentation for this provider for more information on compatibility between provider and Terraform versions.
`

// providerProtocolTooNew is a message sent to the CLI UI if the provider's
// supported protocol versions are too new for the user's version of terraform,
// and the user could either upgrade terraform or choose an older version of the
// provider.
const providerProtocolTooNew = `Provider %q v%s is not compatible with Terraform %s.
You need to downgrade to v%s or earlier. Select it with the following constraint:
	version = %q

Terraform checked all of the plugin versions matching the given constraint:
	%s

Consult the documentation for this provider for more information on compatibility between provider and Terraform versions.
Alternatively, upgrade to the latest version of Terraform for compatibility with newer provider releases.
`

// No version of the provider is compatible.
const errProviderVersionIncompatible = `No compatible versions of provider %s were found.`
