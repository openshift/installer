# v1.9.0

## Enhancements
* `RunListOptions` is generally available, and rename field (Name -> User) by @mjyocca [#472](https://github.com/hashicorp/go-tfe/pull/472)
* [Beta] Adds optional `JsonState` field to `StateVersionCreateOptions` by @megan07 [#514](https://github.com/hashicorp/go-tfe/pull/514)

## Bug Fixes
* Fixed invalid memory address error when using `TaskResults` field by @glennsarti [#517](https://github.com/hashicorp/go-tfe/pull/517)

# v1.8.0

## Enhancements

* Adds support for reading and listing Agents by @laurenolivia [#456](https://github.com/hashicorp/go-tfe/pull/456)
* It was previously logged that we added an `Include` param field to `PolicySetListOptions` to allow policy list to include related resource data such as workspaces, policies, newest_version, or current_version by @Uk1288 [#497](https://github.com/hashicorp/go-tfe/pull/497) in 1.7.0, but this was a mistake and the field is added in v1.8.0

# v1.7.0

## Enhancements

* Adds new run creation attributes: `allow-empty-apply`, `terraform-version`, `plan-only` by @sebasslash [#482](https://github.com/hashicorp/go-tfe/pull/447)
* Adds additional Task Stage and Run Statuses for Pre-plan run tasks by @glennsarti [#469](https://github.com/hashicorp/go-tfe/pull/469)
* Adds `stage` field to the create and update methods for Workspace Run Tasks by @glennsarti [#469](https://github.com/hashicorp/go-tfe/pull/469)
* Adds `ResourcesProcessed`, `StateVersion`, `TerraformVersion`, `Modules`, `Providers`, and `Resources` fields to the State Version struct by @laurenolivia [#484](https://github.com/hashicorp/go-tfe/pull/484)
* Add `Include` param field to `PolicySetListOptions` to allow policy list to include related resource data such as workspaces, policies, newest_version, or current_version by @Uk1288 [#497](https://github.com/hashicorp/go-tfe/pull/497)
* Allow FileTriggersEnabled to be set to false when Git tags are present by @mjyocca @hashimoon  [#468] (https://github.com/hashicorp/go-tfe/pull/468)

# v1.6.0

## Enhancements
* Remove beta messaging for Run Tasks by @glennsarti [#447](https://github.com/hashicorp/go-tfe/pull/447)
* Adds `Description` field to the `RunTask` object by @glennsarti [#447](https://github.com/hashicorp/go-tfe/pull/447)
* Add `Name` field to `OAuthClient` by @barrettclark [#466](https://github.com/hashicorp/go-tfe/pull/466)
* Add support for creating both public and private `RegistryModule` with no VCS connection by @Uk1288 [#460](https://github.com/hashicorp/go-tfe/pull/460)
* Add `ConfigurationSourceAdo` configuration source option by @mjyocca [#467](https://github.com/hashicorp/go-tfe/pull/467)
* [beta] state version outputs may now include a detailed-type attribute in a future API release by @brandonc [#479](https://github.com/hashicorp/go-tfe/pull/429)

# v1.5.0

## Enhancements
* [beta] Add support for triggering Workspace runs through matching Git tags [#434](https://github.com/hashicorp/go-tfe/pull/434)
* Add `Query` param field to `AgentPoolListOptions` to allow searching based on agent pool name, by @JarrettSpiker [#417](https://github.com/hashicorp/go-tfe/pull/417)
* Add organization scope and allowed workspaces field for scope agents by @Netra2104 [#453](https://github.com/hashicorp/go-tfe/pull/453)
* Adds `Namespace` and `RegistryName` fields to `RegistryModuleID` to allow reading of Public Registry Modules by @Uk1288 [#464](https://github.com/hashicorp/go-tfe/pull/464)

## Bug fixes
* Fixed JSON mapping for Configuration Versions failing to properly set the `speculative` property [#459](https://github.com/hashicorp/go-tfe/pull/459)

# v1.4.0

## Enhancements
* Adds `RetryServerErrors` field to the `Config` object by @sebasslash [#439](https://github.com/hashicorp/go-tfe/pull/439)
* Adds support for the GPG Keys API by @sebasslash [#429](https://github.com/hashicorp/go-tfe/pull/429)
* Adds support for new `WorkspaceLimit` Admin setting for organizations [#425](https://github.com/hashicorp/go-tfe/pull/425)
* Adds support for new `ExcludeTags` workspace list filter field by @Uk1288 [#438](https://github.com/hashicorp/go-tfe/pull/438)
* [beta] Adds additional filter fields to `RunListOptions` by @mjyocca [#424](https://github.com/hashicorp/go-tfe/pull/424)
* [beta] Renames the optional StateVersion field `ExtState` to `JSONStateOutputs` and changes the purpose and type by @annawinkler [#444](https://github.com/hashicorp/go-tfe/pull/444) and @brandoncroft [#452](https://github.com/hashicorp/go-tfe/pull/452)

# v1.3.0

## Enhancements
* Adds support for Microsoft Teams notification configuration by @JarrettSpiker [#398](https://github.com/hashicorp/go-tfe/pull/389)
* Add support for Audit Trail API by @sebasslash [#407](https://github.com/hashicorp/go-tfe/pull/407)
* Adds Private Registry Provider, Provider Version, and Provider Platform APIs support by @joekarl and @annawinkler [#313](https://github.com/hashicorp/go-tfe/pull/313)
* Adds List Registry Modules endpoint by @chroju [#385](https://github.com/hashicorp/go-tfe/pull/385)
* Adds `WebhookURL` field to `VCSRepo` struct by @kgns [#413](https://github.com/hashicorp/go-tfe/pull/413)
* Adds `Category` field to `VariableUpdateOptions` struct by @jtyr [#397](https://github.com/hashicorp/go-tfe/pull/397)
* Adds `TriggerPatterns` to `Workspace` by @matejrisek [#400](https://github.com/hashicorp/go-tfe/pull/400)
* [beta] Adds `ExtState` field to `StateVersionCreateOptions` by @brandonc [#416](https://github.com/hashicorp/go-tfe/pull/416)

# v1.2.0

## Enhancements
* Adds support for reading current state version outputs to StateVersionOutputs, which can be useful for reading outputs when users don't have the necessary permissions to read the entire state by @brandonc [#370](https://github.com/hashicorp/go-tfe/pull/370)
* Adds Variable Set methods for `ApplyToWorkspaces` and `RemoveFromWorkspaces` by @byronwolfman [#375](https://github.com/hashicorp/go-tfe/pull/375)
* Adds `Names` query param field to `TeamListOptions` by @sebasslash [#393](https://github.com/hashicorp/go-tfe/pull/393)
* Adds `Emails` query param field to `OrganizationMembershipListOptions` by @sebasslash [#393](https://github.com/hashicorp/go-tfe/pull/393)
* Adds Run Tasks API support by @glennsarti [#381](https://github.com/hashicorp/go-tfe/pull/381), [#382](https://github.com/hashicorp/go-tfe/pull/382) and [#383](https://github.com/hashicorp/go-tfe/pull/383)


## Bug fixes
* Fixes ignored comment when performing apply, discard, cancel, and force-cancel run actions [#388](https://github.com/hashicorp/go-tfe/pull/388)

# v1.1.0

## Enhancements

* Add Variable Set API support by @rexredinger [#305](https://github.com/hashicorp/go-tfe/pull/305)
* Add Comments API support by @alex-ikse [#355](https://github.com/hashicorp/go-tfe/pull/355)
* Add beta support for SSOTeamID to `Team`, `TeamCreateOptions`, `TeamUpdateOptions` by @xlgmokha [#364](https://github.com/hashicorp/go-tfe/pull/364)

# v1.0.0

## Breaking Changes
* Renamed methods named Generate to Create for `AgentTokens`, `OrganizationTokens`, `TeamTokens`, `UserTokens` by @sebasslash [#327](https://github.com/hashicorp/go-tfe/pull/327)
* Methods that express an action on a relationship have been prefixed with a verb, e.g `Current()` is now `ReadCurrent()` by @sebasslash [#327](https://github.com/hashicorp/go-tfe/pull/327)
* All list option structs are now pointers @uturunku1 [#309](https://github.com/hashicorp/go-tfe/pull/309)
* All errors have been refactored into constants in `errors.go` @uturunku1 [#310](https://github.com/hashicorp/go-tfe/pull/310)
* The `ID` field in Create/Update option structs has been renamed to `Type` in accordance with the JSON:API spec by @omarismail, @uturunku1 [#190](https://github.com/hashicorp/go-tfe/pull/190), [#323](https://github.com/hashicorp/go-tfe/pull/323), [#332](https://github.com/hashicorp/go-tfe/pull/332)
* Nested URL params (consisting of an organization, module and provider name) used to identify a `RegistryModule` have been refactored into a struct `RegistryModuleID` by @sebasslash [#337](https://github.com/hashicorp/go-tfe/pull/337)


## Enhancements
* Added missing include fields for `AdminRuns`, `AgentPools`, `ConfigurationVersions`, `OAuthClients`, `Organizations`, `PolicyChecks`, `PolicySets`, `Policies` and `RunTriggers` by @uturunku1 [#339](https://github.com/hashicorp/go-tfe/pull/339)
* Cleanup documentation and improve consistency by @uturunku1 [#331](https://github.com/hashicorp/go-tfe/pull/331)
* Add more linters to our CI pipeline by @sebasslash [#326](https://github.com/hashicorp/go-tfe/pull/326)
* Resolve `TFE_HOSTNAME` as fallback for `TFE_ADDRESS` by @sebasslash [#340](https://github.com/hashicorp/go-tfe/pull/326)
* Adds a `fetching` status to `RunStatus` and adds the `Archive` method to the ConfigurationVersions interface by @mpminardi [#338](https://github.com/hashicorp/go-tfe/pull/338)
* Added a `Download` method to the `ConfigurationVersions` interface by @tylerwolf [#358](https://github.com/hashicorp/go-tfe/pull/358)
* API Coverage documentation by @laurenolivia [#334](https://github.com/hashicorp/go-tfe/pull/334)

## Bug Fixes
* Fixed invalid memory address error when `AdminSMTPSettingsUpdateOptions.Auth` field is empty and accessed by @uturunku1 [#335](https://github.com/hashicorp/go-tfe/pull/335)
