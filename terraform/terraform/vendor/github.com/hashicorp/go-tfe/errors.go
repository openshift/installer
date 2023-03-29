package tfe

import (
	"errors"
)

// Generic errors applicable to all resources.
var (
	// ErrUnauthorized is returned when receiving a 401.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrResourceNotFound is returned when receiving a 404.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrMissingDirectory is returned when the path does not have an existing directory.
	ErrMissingDirectory = errors.New("path needs to be an existing directory")

	// ErrNamespaceNotAuthorized is returned when a user attempts to perform an action
	// on a namespace (organization) they do not have access to.
	ErrNamespaceNotAuthorized = errors.New("namespace not authorized")
)

// Options/fields that cannot be defined
var (
	ErrUnsupportedOperations = errors.New("operations is deprecated and cannot be specified when execution mode is used")

	ErrUnsupportedPrivateKey = errors.New("private Key can only be present with Azure DevOps Server service provider")

	ErrUnsupportedBothTagsRegexAndFileTriggersEnabled = errors.New(`"TagsRegex" cannot be populated when "FileTriggersEnabled" is true`)

	ErrUnsupportedBothTagsRegexAndTriggerPatterns = errors.New(`"TagsRegex" and "TriggerPrefixes" cannot be populated at the same time`)

	ErrUnsupportedBothTagsRegexAndTriggerPrefixes = errors.New(`"TagsRegex" and "TriggerPatterns" cannot be populated at the same time`)

	ErrUnsupportedRunTriggerType = errors.New(`"RunTriggerType" must be "inbound" when requesting "include" query params`)

	ErrUnsupportedBothTriggerPatternsAndPrefixes = errors.New(`"TriggerPatterns" and "TriggerPrefixes" cannot be populated at the same time`)

	ErrUnsupportedBothNamespaceAndPrivateRegistryName = errors.New(`"Namespace" cannot be populated when "RegistryName" is "private"`)
)

// Library errors that usually indicate a bug in the implementation of go-tfe
var (
	ErrItemsMustBeSlice = errors.New(`model field "Items" must be a slice`) // ErrItemsMustBeSlice is returned when an API response attribute called Items is not a slice

	ErrInvalidRequestBody = errors.New("go-tfe bug: DELETE/PATCH/POST body must be nil, ptr, or ptr slice") // ErrInvalidRequestBody is returned when a request body for DELETE/PATCH/POST is not a reference type

	ErrInvalidStructFormat = errors.New("go-tfe bug: struct can't use both json and jsonapi attributes") // ErrInvalidStructFormat is returned when a mix of json and jsonapi tagged fields are used in the same struct
)

// Resource Errors
var (
	ErrWorkspaceLocked = errors.New("workspace already locked") // ErrWorkspaceLocked is returned when trying to lock a
	// locked workspace.

	ErrWorkspaceNotLocked = errors.New("workspace already unlocked") // ErrWorkspaceNotLocked is returned when trying to unlock
	// a unlocked workspace.

	ErrWorkspaceLockedByRun = errors.New("unable to unlock workspace locked by run") // ErrWorkspaceLockedByRun is returned when trying to unlock a
	// workspace locked by a run
)

// Invalid values for resources/struct fields
var (
	ErrInvalidWorkspaceID = errors.New("invalid value for workspace ID")

	ErrInvalidWorkspaceValue = errors.New("invalid value for workspace")

	ErrInvalidTerraformVersionID = errors.New("invalid value for terraform version ID")

	ErrInvalidTerraformVersionType = errors.New("invalid type for terraform version. Please use 'terraform-version'")

	ErrInvalidConfigVersionID = errors.New("invalid value for configuration version ID")

	ErrInvalidCostEstimateID = errors.New("invalid value for cost estimate ID")

	ErrInvalidSMTPAuth = errors.New("invalid smtp auth type")

	ErrInvalidAgentPoolID = errors.New("invalid value for agent pool ID")

	ErrInvalidAgentTokenID = errors.New("invalid value for agent token ID")

	ErrInvalidRunID = errors.New("invalid value for run ID")

	ErrInvalidRunTaskCategory = errors.New(`category must be "task"`)

	ErrInvalidRunTaskID = errors.New("invalid value for run task ID")

	ErrInvalidRunTaskURL = errors.New("invalid url for run task URL")

	ErrInvalidWorkspaceRunTaskID = errors.New("invalid value for workspace run task ID")

	ErrInvalidWorkspaceRunTaskType = errors.New(`invalid value for type, please use "workspace-tasks"`)

	ErrInvalidTaskResultID = errors.New("invalid value for task result ID")

	ErrInvalidTaskStageID = errors.New("invalid value for task stage ID")

	ErrInvalidApplyID = errors.New("invalid value for apply ID")

	ErrInvalidOrg = errors.New("invalid value for organization")

	ErrInvalidName = errors.New("invalid value for name")

	ErrInvalidNotificationConfigID = errors.New("invalid value for notification configuration ID")

	ErrInvalidMembership = errors.New("invalid value for membership")

	ErrInvalidMembershipIDs = errors.New("invalid value for organization membership ids")

	ErrInvalidOauthClientID = errors.New("invalid value for OAuth client ID")

	ErrInvalidOauthTokenID = errors.New("invalid value for OAuth token ID")

	ErrInvalidPolicySetID = errors.New("invalid value for policy set ID")

	ErrInvalidPolicyCheckID = errors.New("invalid value for policy check ID")

	ErrInvalidTag = errors.New("invalid tag id")

	ErrInvalidPlanExportID = errors.New("invalid value for plan export ID")

	ErrInvalidPlanID = errors.New("invalid value for plan ID")

	ErrInvalidParamID = errors.New("invalid value for parameter ID")

	ErrInvalidPolicyID = errors.New("invalid value for policy ID")

	ErrInvalidProvider = errors.New("invalid value for provider")

	ErrInvalidVersion = errors.New("invalid value for version")

	ErrInvalidRunTriggerID = errors.New("invalid value for run trigger ID")

	ErrInvalidRunTriggerType = errors.New(`invalid value or no value for RunTriggerType. It must be either "inbound" or "outbound"`)

	ErrInvalidIncludeValue = errors.New(`invalid value for "include" field`)

	ErrInvalidSHHKeyID = errors.New("invalid value for SSH key ID")

	ErrInvalidStateVerID = errors.New("invalid value for state version ID")

	ErrInvalidOutputID = errors.New("invalid value for state version output ID")

	ErrInvalidAccessTeamID = errors.New("invalid value for team access ID")

	ErrInvalidTeamID = errors.New("invalid value for team ID")

	ErrInvalidUsernames = errors.New("invalid value for usernames")

	ErrInvalidUserID = errors.New("invalid value for user ID")

	ErrInvalidUserValue = errors.New("invalid value for user")

	ErrInvalidTokenID = errors.New("invalid value for token ID")

	ErrInvalidCategory = errors.New("category must be policy-set")

	ErrInvalidPolicies = errors.New("must provide at least one policy")

	ErrInvalidVariableID = errors.New("invalid value for variable ID")

	ErrInvalidNotificationTrigger = errors.New("invalid value for notification trigger")

	ErrInvalidVariableSetID = errors.New("invalid variable set ID")

	ErrInvalidCommentID = errors.New("invalid value for comment ID")

	ErrInvalidCommentBody = errors.New("invalid value for comment body")

	ErrInvalidNamespace = errors.New("invalid value for namespace")

	ErrInvalidKeyID = errors.New("invalid value for key-id")

	ErrInvalidOS = errors.New("invalid value for OS")

	ErrInvalidArch = errors.New("invalid value for arch")

	ErrInvalidAgentID = errors.New("invalid value for Agent ID")

	ErrInvalidRegistryName = errors.New(`invalid value for registry-name. It must be either "private" or "public"`)
)

// Missing values for required field/option
var (
	ErrRequiredAccess = errors.New("access is required")

	ErrRequiredAgentPoolID = errors.New("'agent' execution mode requires an agent pool ID to be specified")

	ErrRequiredAgentMode = errors.New("specifying an agent pool ID requires 'agent' execution mode")

	ErrRequiredCategory = errors.New("category is required")

	ErrRequiredDestinationType = errors.New("destination type is required")

	ErrRequiredDataType = errors.New("data type is required")

	ErrRequiredKey = errors.New("key is required")

	ErrRequiredName = errors.New("name is required")

	ErrRequiredEnabled = errors.New("enabled is required")

	ErrRequiredEnforce = errors.New("enforce is required")

	ErrRequiredEnforcementPath = errors.New("enforcement path is required")

	ErrRequiredEnforcementMode = errors.New("enforcement mode is required")

	ErrRequiredEmail = errors.New("email is required")

	ErrRequiredM5 = errors.New("MD5 is required")

	ErrRequiredURL = errors.New("url is required")

	ErrRequiredAPIURL = errors.New("API URL is required")

	ErrRequiredHTTPURL = errors.New("HTTP URL is required")

	ErrRequiredServiceProvider = errors.New("service provider is required")

	ErrRequiredProvider = errors.New("provider is required")

	ErrRequiredOauthToken = errors.New("OAuth token is required")

	ErrRequiredOauthTokenID = errors.New("oauth token ID is required")

	ErrRequiredTestNumber = errors.New("TestNumber is required")

	ErrMissingTagIdentifier = errors.New("must specify at least one tag by ID or name")

	ErrAgentTokenDescription = errors.New("agent token description can't be blank")

	ErrRequiredTagID = errors.New("you must specify at least one tag id to remove")

	ErrRequiredTagWorkspaceID = errors.New("you must specify at least one workspace to add tag to")

	ErrRequiredWorkspace = errors.New("workspace is required")

	ErrRequiredWorkspaceID = errors.New("workspace ID is required")

	ErrWorkspacesRequired = errors.New("workspaces is required")

	ErrWorkspaceMinLimit = errors.New("must provide at least one workspace")

	ErrRequiredPlan = errors.New("plan is required")

	ErrRequiredPolicies = errors.New("policies is required")

	ErrRequiredVersion = errors.New("version is required")

	ErrRequiredVCSRepo = errors.New("vcs repo is required")

	ErrRequiredIdentifier = errors.New("identifier is required")

	ErrRequiredDisplayIdentifier = errors.New("display identifier is required")

	ErrRequiredSha = errors.New("sha is required")

	ErrRequiredSourceable = errors.New("sourceable is required")

	ErrRequiredValue = errors.New("value is required")

	ErrRequiredOrg = errors.New("organization is required")

	ErrRequiredTeam = errors.New("team is required")

	ErrRequiredStateVerListOps = errors.New("StateVersionListOptions is required")

	ErrRequiredTeamAccessListOps = errors.New("TeamAccessListOptions is required")

	ErrRequiredRunTriggerListOps = errors.New("RunTriggerListOptions is required")

	ErrRequiredTFVerCreateOps = errors.New("version, URL and sha is required for AdminTerraformVersionCreateOptions")

	ErrRequiredSerial = errors.New("serial is required")

	ErrRequiredState = errors.New("state is required")

	ErrRequiredSHHKeyID = errors.New("SSH key ID is required")

	ErrRequiredOnlyOneField = errors.New("only one of usernames or organization membership ids can be provided")

	ErrRequiredUsernameOrMembershipIds = errors.New("usernames or organization membership ids are required")

	ErrRequiredGlobalFlag = errors.New("global flag is required")

	ErrRequiredWorkspacesList = errors.New("no workspaces list provided")

	ErrCommentBody = errors.New("comment body is required")

	ErrEmptyTeamName = errors.New("team name can not be empty")

	ErrInvalidEmail = errors.New("email is invalid")

	ErrRequiredPrivateRegistry = errors.New("only private registry is allowed")

	ErrRequiredOS = errors.New("OS is required")

	ErrRequiredArch = errors.New("arch is required")

	ErrRequiredShasum = errors.New("shasum is required")

	ErrRequiredFilename = errors.New("filename is required")

	ErrInvalidAsciiArmor = errors.New("ascii armor is invalid")

	ErrRequiredNamespace = errors.New("namespace is required for public registry")

	ErrTerraformVersionValidForPlanOnly = errors.New("setting terraform-version is only valid when plan-only is set to true")
)
