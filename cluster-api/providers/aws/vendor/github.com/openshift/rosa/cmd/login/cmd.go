/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package login

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/golang-jwt/jwt/v4"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/authentication"
	"github.com/openshift-online/ocm-sdk-go/authentication/securestore"
	"github.com/spf13/cobra"
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/cmd/logout"
	"github.com/openshift/rosa/pkg/arguments"
	"github.com/openshift/rosa/pkg/config"
	"github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/fedramp"
	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/output"
	"github.com/openshift/rosa/pkg/properties"
	"github.com/openshift/rosa/pkg/reporter"
	"github.com/openshift/rosa/pkg/rosa"
)

// #nosec G101
var uiTokenPage string = "https://console.redhat.com/openshift/token/rosa"

const oauthClientId = "ocm-cli"

var reAttempt bool

var env string

var args struct {
	tokenURL      string
	clientID      string
	clientSecret  string
	scopes        []string
	env           string
	token         string
	insecure      bool
	useAuthCode   bool
	useDeviceCode bool
	rhRegion      string
}

var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your Red Hat account",
	Long: fmt.Sprintf("Log in to your Red Hat account, saving the credentials to the configuration file or OS Keyring.\n"+
		"The supported mechanism is by using a token, which can be obtained at: %s\n\n"+
		"The application looks for the token in the following order, stopping when it finds it:\n"+
		fmt.Sprintf("\t1. OS Keyring via Environment variable (%s)\n", properties.KeyringEnvKey)+
		"\t2. Command-line flags\n"+
		"\t3. Environment variable (ROSA_TOKEN)\n"+
		"\t4. Environment variable (OCM_TOKEN)\n"+
		"\t5. Configuration file\n"+
		"\t6. Command-line prompt\n", uiTokenPage),
	Example: fmt.Sprintf(`  # Login to the OpenShift API with an existing token generated from %s
  rosa login --token=$OFFLINE_ACCESS_TOKEN`, uiTokenPage),
	Run:  run,
	Args: cobra.NoArgs,
}

func init() {
	flags := Cmd.Flags()
	flags.StringVar(
		&args.tokenURL,
		"token-url",
		"",
		fmt.Sprintf(
			"OpenID token URL. The default value is '%s'.",
			sdk.DefaultTokenURL,
		),
	)
	flags.StringVar(
		&args.clientID,
		"client-id",
		"",
		fmt.Sprintf(
			"OpenID client identifier. The default value is '%s'.",
			sdk.DefaultClientID,
		),
	)
	flags.StringVar(
		&args.clientSecret,
		"client-secret",
		"",
		"OpenID client secret.",
	)
	flags.StringSliceVar(
		&args.scopes,
		"scope",
		sdk.DefaultScopes,
		"OpenID scope. If this option is used it will completely replace the default "+
			"scopes. Can be repeated multiple times to specify multiple scopes.",
	)
	flags.SetNormalizeFunc(arguments.NormalizeFlags)
	flags.StringVar(
		&args.env,
		arguments.NewEnvFlag,
		"",
		"Environment of the API gateway. The value can be the complete URL or an alias. "+
			"The valid aliases are 'production', 'staging' and 'integration'.",
	)
	flags.MarkHidden(arguments.NewEnvFlag)
	flags.StringVarP(
		&args.token,
		"token",
		"t",
		"",
		fmt.Sprintf("Access or refresh token generated from %s.", uiTokenPage),
	)
	flags.BoolVar(
		&args.insecure,
		"insecure",
		false,
		"Enables insecure communication with the server. This disables verification of TLS "+
			"certificates and host names.",
	)
	flags.BoolVar(
		&args.useAuthCode,
		"use-auth-code",
		false,
		"Login using OAuth Authorization Code. This should be used for most cases where a "+
			"browser is available. See --use-device-code for remote hosts and containers.",
	)
	flags.BoolVar(
		&args.useDeviceCode,
		"use-device-code",
		false,
		"Login using OAuth Device Code. "+
			"This should only be used for remote hosts and containers where browsers are "+
			"not available. See --use-auth-code for all other scenarios.",
	)
	flags.StringVar(
		&args.rhRegion,
		"rh-region",
		"",
		"OCM data sovereignty region identifier. --env will be used to initiate a service discovery "+
			"request to find the region URL matching the provided identifier. Use `rosa list rh-regions` "+
			"to see available regions.",
	)
	flags.MarkHidden("rh-region")
	arguments.AddRegionFlag(flags)
	fedramp.AddFlag(flags)
}

func run(cmd *cobra.Command, argv []string) {
	r := rosa.NewRuntime()
	defer r.Cleanup()
	err := runWithRuntime(r, cmd, argv)
	if err != nil {
		r.Reporter.Errorf(err.Error())
		os.Exit(1)
	}
}

func runWithRuntime(r *rosa.Runtime, cmd *cobra.Command, argv []string) error {
	ctx := cmd.Context()
	var spin *spinner.Spinner
	if r.Reporter.IsTerminal() && !output.HasFlag() {
		spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	}

	// Check mandatory options:
	env = args.env

	// Fail fast if config is keyring managed and invalid
	if keyring, ok := config.IsKeyringManaged(); ok {
		err := securestore.ValidateBackend(keyring)
		if err != nil {
			return fmt.Errorf("Error validating keyring: %v", err)
		}
	}

	// Confirm that token is not passed with auth code flags
	if (args.useAuthCode || args.useDeviceCode) && args.token != "" {
		r.Reporter.Errorf("Token cannot be passed with '--use-auth-code' or '--use-device-code' commands")
		os.Exit(1)
	}

	// FedRAMP does not support oauth code flow login yet
	if fedramp.HasFlag(cmd) && (args.useAuthCode || args.useDeviceCode) {
		r.Reporter.Errorf("This login method is currently not supported with FedRAMP")
		os.Exit(1)
	}

	if args.useAuthCode {
		r.Reporter.Infof("You will now be redirected to Red Hat SSO login")

		if spin != nil {
			spin.Start()
		}
		// Short wait for a less jarring experience
		time.Sleep(2 * time.Second)
		if spin != nil {
			spin.Stop()
		}
		token, err := authentication.InitiateAuthCode(oauthClientId)
		if err != nil {
			return fmt.Errorf("An error occurred while retrieving the token: %v", err)
		}
		args.token = token
		args.clientID = oauthClientId
		r.Reporter.Infof("Token received successfully")
	} else if args.useDeviceCode {
		deviceAuthConfig := &authentication.DeviceAuthConfig{
			ClientID: oauthClientId,
		}
		_, err := deviceAuthConfig.InitiateDeviceAuth(ctx)
		if err != nil {
			return fmt.Errorf("An error occurred while initiating device auth: %v", err)
		}
		deviceAuthResp := deviceAuthConfig.DeviceAuthResponse

		r.Reporter.Infof("To login, navigate to %v on another device and enter code %v",
			deviceAuthResp.VerificationURI, deviceAuthResp.UserCode)
		r.Reporter.Infof("Checking status every %v seconds...", deviceAuthResp.Interval)
		token, err := deviceAuthConfig.PollForTokenExchange(ctx)
		if err != nil {
			return fmt.Errorf("An error occurred while polling for token exchange: %v", err)
		}
		args.token = token
		args.clientID = oauthClientId
	}

	// Load the configuration file:
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load config file: %v", err)
	}
	if cfg == nil || config.IsNotValid(cfg) {
		cfg = new(config.Config)
	}

	token := args.token

	// Determine if we should be using the FedRAMP environment:
	err = CheckAndLogIntoFedramp(fedramp.HasFlag(cmd), fedramp.HasAdminFlag(cmd), cfg, token, r)
	if err != nil {
		r.Reporter.Errorf("%s", err.Error())
		os.Exit(1)
	}

	haveReqs := token != "" || (args.clientID != "" && args.clientSecret != "")

	// Verify environment variables:
	if !haveReqs && !reAttempt && !fedramp.Enabled() {
		token = os.Getenv(constants.RosaToken)
		if token == "" {
			token = os.Getenv(constants.OcmToken)
		}
		haveReqs = token != ""
	}

	// Verify configuration file:
	if !haveReqs {
		armed, err := cfg.Armed()
		if err != nil {
			return fmt.Errorf("Failed to verify configuration: %v", err)
		}
		haveReqs = armed
	}

	// Prompt the user for token:
	if !haveReqs {
		fmt.Println("To login to your Red Hat account, get an offline access token at", uiTokenPage)
		token, err = interactive.GetPassword(interactive.Input{
			Question: "Copy the token and paste it here",
			Required: true,
		})
		if err != nil {
			return fmt.Errorf("Failed to parse token: %v", err)
		}
		haveReqs = token != ""
	}

	if !haveReqs {
		return fmt.Errorf("Failed to login to OCM. See 'rosa login --help' for information.")
	}

	// Red Hat SSO does not issue encrypted refresh tokens, but AWS Cognito does. If the token
	// is encrypted we can safely assume that the user is trying to use the FedRAMP environment.
	if config.IsEncryptedToken(token) {
		fedramp.Enable()
	}

	gatewayURL, err := ocm.ResolveGatewayUrl(env, cfg)
	if err != nil {
		r.Reporter.Errorf("Failed to resolve gateway URL: %v", err)
		os.Exit(1)
	}

	var ok bool
	tokenURL := sdk.DefaultTokenURL
	if args.tokenURL != "" {
		tokenURL = args.tokenURL
	}
	clientID := sdk.DefaultClientID
	if args.clientID != "" {
		clientID = args.clientID
	}

	if strings.ToLower(env) == ocm.ProductionAlias {
		r.Reporter.Warnf("\"prod\" provided as the environment, aliasing environment to \"production\"")
		env = ocm.Production
	}

	// Override configuration details for FedRAMP:
	if fedramp.Enabled() {
		clientID = fedramp.ClientID
		if args.clientID != "" {
			clientID = args.clientID
		}
		if fedramp.HasAdminFlag(cmd) {
			if !fedramp.IsValidEnv(env) {
				_ = r.Reporter.Errorf("%s is an invalid environment name, please use one of: ",
					strings.Join(ocm.ValidOCMUrlAliases(), ", "))
				os.Exit(1)
			}
			gatewayURL, ok = fedramp.AdminURLAliases[env]
			if !ok {
				gatewayURL = env
			}
			tokenURL, ok = fedramp.AdminTokenURLs[env]
			if !ok {
				tokenURL = args.tokenURL
			}
			clientID, ok = fedramp.AdminClientIDs[env]
			if !ok {
				clientID = args.clientID
			}
		} else {
			gatewayURL, ok = fedramp.URLAliases[env]
			if !ok {
				gatewayURL = env
			}
			tokenURL, ok = fedramp.TokenURLs[env]
			if !ok {
				tokenURL = args.tokenURL
			}
		}
	}

	// If an --rh-region is provided, gatewayURL is resolved from ResolveGatewayUrl() and then used to initiate
	// service discovery for the environment gatewayURL is a part of, but it (and
	// ultimately the cfg.URL) is then updated to the URL of the matching --rh-region:
	//   1. resolve the gatewayURL as above
	//   2. fetch a well-known file from sdk.GetRhRegion
	//   3. update the gatewayURL to the region URL matching args.rhRegion
	//
	// So `--env=https://api.stage.openshift.com --rh-region=ap-southeast-1` might result in
	// gatewayURL/cfg.URL being mutated to "https://api.aws.ap-southeast-1.stage.openshift.com"
	//
	// See ocm-sdk-go/rh_region.go for full details on how service discovery works.
	if args.rhRegion != "" {
		regValue, err := sdk.GetRhRegion(gatewayURL, args.rhRegion)
		if err != nil {
			r.Reporter.Errorf("Can't find region: %v", err)
			os.Exit(1)
		}
		gatewayURL = fmt.Sprintf("https://%s", regValue.URL)
	}

	// Update the configuration with the values given in the command line:
	cfg.TokenURL = tokenURL
	cfg.ClientID = clientID
	cfg.ClientSecret = args.clientSecret
	cfg.Scopes = args.scopes
	cfg.URL = gatewayURL
	cfg.Insecure = args.insecure
	cfg.FedRAMP = fedramp.Enabled()

	if token != "" {
		if config.IsEncryptedToken(token) {
			cfg.AccessToken = ""
			cfg.RefreshToken = token
		} else {
			// If a token has been provided parse it:
			jwtToken, err := config.ParseToken(token)
			if err != nil {
				return fmt.Errorf("Failed to parse token: %v", err)
			}

			// Put the token in the place of the configuration that corresponds to its type:
			typ, err := tokenType(jwtToken)
			if err != nil {
				return fmt.Errorf("Failed to extract type from 'typ' claim of token: %v", err)
			}
			switch typ {
			case "Bearer", "":
				cfg.AccessToken = token
				cfg.RefreshToken = ""
			case "Refresh", "Offline":
				cfg.AccessToken = ""
				cfg.RefreshToken = token
			default:
				return fmt.Errorf("Don't know how to handle token type '%s' in token", typ)
			}
		}
	}

	// Create a connection and get the token to verify that the crendentials are correct:
	r.OCMClient, err = ocm.NewClient().
		Config(cfg).
		Logger(r.Logger).
		Build()
	if err != nil {
		if strings.Contains(err.Error(), "token needs to be updated") && !reAttempt {
			reattemptLogin(cmd, argv)
			return nil
		} else {
			return fmt.Errorf("Failed to create OCM connection: %v", err)
		}
	}

	accessToken, refreshToken, err := r.OCMClient.GetConnectionTokens()
	if err != nil {
		return fmt.Errorf(
			"Failed to get token. Your session might be expired: %v\nGet a new offline access token at %s",
			err, uiTokenPage)
	}
	reAttempt = false
	// Save the configuration:
	cfg.AccessToken = accessToken
	cfg.RefreshToken = refreshToken
	err = config.Save(cfg)
	if err != nil {
		return fmt.Errorf("Failed to save config file: %v", err)
	}

	username, err := cfg.GetData("preferred_username")
	if err != nil {
		username, err = cfg.GetData("username")
		if err != nil {
			return fmt.Errorf("Failed to get username: %v", err)
		}
	}

	r.Reporter.Infof("Logged in as '%s' on '%s'", username, cfg.URL)
	r.OCMClient.LogEvent("ROSALoginSuccess", map[string]string{
		ocm.Response: ocm.Success,
		ocm.Username: username,
		ocm.URL:      cfg.URL,
	})

	if args.useAuthCode || args.useDeviceCode {
		ssoURL, err := url.Parse(cfg.TokenURL)
		if err != nil {
			return fmt.Errorf("can't parse token url '%s': %v", args.tokenURL, err)
		}
		ssoHost := ssoURL.Scheme + "://" + ssoURL.Hostname()

		r.Reporter.Infof("To switch accounts, logout from %s and run `rosa logout` "+
			"before attempting to login again", ssoHost)
	}

	return nil
}

func reattemptLogin(cmd *cobra.Command, argv []string) {
	logout.Cmd.Run(cmd, argv)
	reAttempt = true
	run(cmd, argv)
}

// tokenType extracts the value of the `typ` claim. It returns the value as a string, or the empty
// string if there is no such claim.
func tokenType(jwtToken *jwt.Token) (typ string, err error) {
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("Expected map claims but got %T", claims)
		return
	}
	claim, ok := claims["typ"]
	if !ok {
		return
	}
	value, ok := claim.(string)
	if !ok {
		err = fmt.Errorf("Expected string 'typ' but got %T", claim)
		return
	}
	typ = value
	return
}

func Call(cmd *cobra.Command, argv []string, reporter reporter.Logger) error {
	loginFlags := []string{"token-url", "client-id", "client-secret", "scope", arguments.NewEnvFlag, "token", "insecure"}
	hasLoginFlags := false
	// Check if the user set login flags
	for _, loginFlag := range loginFlags {
		if cmd.Flags().Changed(loginFlag) {
			hasLoginFlags = true
			break
		}
	}
	if hasLoginFlags {
		// Always force login if user sets login flags
		run(cmd, argv)
		return nil
	}

	// Verify if user is already logged in:
	isLoggedIn := false
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Failed to load config file: %v", err)
	}
	if cfg != nil && !config.IsNotValid(cfg) {
		// Check that credentials in the config file are valid
		isLoggedIn, err = cfg.Armed()
		if err != nil {
			return fmt.Errorf("Failed to determine if user is logged in: %v", err)
		}
	}

	if isLoggedIn {
		username, err := cfg.GetData("preferred_username")
		if err != nil {
			username, err = cfg.GetData("username")
			if err != nil {
				return fmt.Errorf("Failed to get username: %v", err)
			}
		}

		if reporter.IsTerminal() {
			reporter.Infof("Logged in as '%s' on '%s'", username, cfg.URL)
		}
		return nil
	}

	run(cmd, argv)
	return nil
}

func CheckAndLogIntoFedramp(hasFlag, hasAdminFlag bool, cfg *config.Config, token string,
	runtime *rosa.Runtime) error {
	if hasFlag ||
		(cfg.FedRAMP && token == "") ||
		fedramp.IsGovRegion(arguments.GetRegion()) ||
		config.IsEncryptedToken(token) {
		// Display error to user if they attempt to log into govcloud without a region specified (fixes OCM-5718)
		if !fedramp.IsGovRegion(arguments.GetRegion()) {
			return errors.Errorf("When logging into the FedRAMP environment, a recognized us-gov region needs " +
				"to be specified. Example: --region us-gov-west-1")
		}

		fedramp.Enable()
		// Always default to prod
		if env == sdk.DefaultURL || env == "" {
			env = "production"
		}
		if hasAdminFlag {
			uiTokenPage = fedramp.AdminLoginURLs[env]
		} else {
			uiTokenPage = fedramp.LoginURLs[env]
		}
	} else {
		fedramp.Disable()
	}
	return nil
}
