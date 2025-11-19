/*
Copyright (c) 2023 Red Hat, Inc.

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

package oidcconfig

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/openshift-online/ocm-common/pkg/rosa/oidcconfigs"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"
	"github.com/zgalor/weberr"

	"github.com/openshift/rosa/cmd/create/oidcprovider"
	"github.com/openshift/rosa/pkg/arguments"
	"github.com/openshift/rosa/pkg/aws"
	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
	. "github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/interactive/confirm"
	interactiveRoles "github.com/openshift/rosa/pkg/interactive/roles"
	"github.com/openshift/rosa/pkg/output"
	"github.com/openshift/rosa/pkg/rosa"
)

var args struct {
	region           string
	rawFiles         bool
	userPrefix       string
	managed          bool
	installerRoleArn string
}

var Cmd = &cobra.Command{
	Use:     "oidc-config",
	Aliases: []string{"oidcconfig"},
	Short:   "Create OIDC config compliant with OIDC protocol.",
	Long: "Create OIDC config in a S3 bucket for the " +
		"client AWS account and populates it to be compliant with OIDC protocol. " +
		"It also creates a Secret in Secrets Manager containing the private key.",
	Example: `  # Create OIDC config
	rosa create oidc-config`,
	Run:  run,
	Args: cobra.NoArgs,
}

const (
	maxLengthUserPrefix = 15

	rawFilesFlag   = "raw-files"
	userPrefixFlag = "prefix"
	managedFlag    = "managed"
)

func init() {
	flags := Cmd.Flags()

	flags.BoolVar(
		&args.rawFiles,
		rawFilesFlag,
		false,
		"Creates OIDC config documents (Private RSA key, Discovery document, JSON Web Key Set) "+
			"and saves locally for the client to create the configuration.",
	)

	flags.StringVar(
		&args.userPrefix,
		userPrefixFlag,
		"",
		"Prefix for the OIDC configuration, secret and provider.",
	)

	flags.BoolVar(
		&args.managed,
		managedFlag,
		true,
		"Indicates whether it is a Red Hat managed or unmanaged (Customer hosted) OIDC Configuration.",
	)

	// normalizing installer role argument to support deprecated flag
	flags.SetNormalizeFunc(arguments.NormalizeFlags)
	flags.StringVar(
		&args.installerRoleArn,
		InstallerRoleArnFlag,
		"",
		"STS Role ARN with get secrets permission.",
	)

	interactive.AddModeFlag(Cmd)

	confirm.AddFlag(flags)
	interactive.AddFlag(flags)
	arguments.AddRegionFlag(flags)
	output.AddFlag(Cmd)
}

func checkInteractiveModeNeeded(cmd *cobra.Command) {
	modeNotChanged := !cmd.Flags().Changed("mode")
	if modeNotChanged && !cmd.Flags().Changed(rawFilesFlag) {
		interactive.Enable()
		return
	}
	oidcConfigKindNotSet := !cmd.Flags().Changed(managedFlag)
	if oidcConfigKindNotSet && !confirm.Yes() {
		interactive.Enable()
		return
	}
	modeIsAuto := cmd.Flag("mode").Value.String() == interactive.ModeAuto
	installerRoleArnNotSet := (!cmd.Flags().Changed(InstallerRoleArnFlag) || args.installerRoleArn == "") &&
		!confirm.Yes()
	if !args.managed && (modeNotChanged || (modeIsAuto && installerRoleArnNotSet)) {
		interactive.Enable()
		return
	}
}

func run(cmd *cobra.Command, _ []string) {
	r := rosa.NewRuntime().WithAWS().WithOCM()
	defer r.Cleanup()

	mode, err := interactive.GetMode()
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}

	// Get AWS region
	region, err := aws.GetRegion(arguments.GetRegion())
	if err != nil {
		r.Reporter.Errorf("Error getting region: %v", err)
		os.Exit(1)
	}
	args.region = region

	checkInteractiveModeNeeded(cmd)

	if interactive.Enabled() && !cmd.Flags().Changed(managedFlag) {
		args.managed = confirm.Prompt(true, "Would you like to create a Managed (Red Hat hosted) OIDC Configuration")
	}

	if args.rawFiles && mode != "" {
		r.Reporter.Warnf("--%s param is not supported alongside --mode param.", rawFilesFlag)
		os.Exit(1)
	}

	if args.rawFiles && args.installerRoleArn != "" {
		r.Reporter.Warnf("--%s param is not supported alongside --%s param", rawFilesFlag, InstallerRoleArnFlag)
		os.Exit(1)
	}

	if args.rawFiles && args.managed {
		r.Reporter.Warnf("--%s param is not supported alongside --%s param", rawFilesFlag, managedFlag)
		os.Exit(1)
	}

	if !args.rawFiles && interactive.Enabled() && !cmd.Flags().Changed("mode") {
		question := "OIDC Config creation mode"
		if args.managed {
			r.Reporter.Warnf("For a managed OIDC Config only auto mode is supported. " +
				"However, you may choose the provider creation mode")
			question = "OIDC Provider creation mode"
		}
		mode, err = interactive.GetOptionMode(cmd, mode, question)
		if err != nil {
			r.Reporter.Errorf("Expected a valid %s: %s", question, err)
			os.Exit(1)
		}
	}

	if output.HasFlag() && mode != "" && mode != interactive.ModeAuto {
		r.Reporter.Warnf("--output param is not supported outside auto mode.")
		os.Exit(1)
	}

	if args.managed && args.userPrefix != "" {
		r.Reporter.Warnf("--%s param is not supported for managed OIDC config", userPrefixFlag)
		os.Exit(1)
	}

	if args.managed && args.installerRoleArn != "" {
		r.Reporter.Warnf("--%s param is not supported for managed OIDC config", InstallerRoleArnFlag)
		os.Exit(1)
	}

	if !args.managed {
		if !args.rawFiles {
			if !output.HasFlag() && r.Reporter.IsTerminal() {
				r.Reporter.Infof("This command will create a S3 bucket populating it with documents " +
					"to be compliant with OIDC protocol. It will also create a Secret in Secrets Manager containing the private key")
			}
			if mode == interactive.ModeAuto && (interactive.Enabled() || (confirm.Yes() && args.installerRoleArn == "")) {
				args.installerRoleArn = interactiveRoles.
					GetInstallerRoleArn(
						r,
						cmd,
						args.installerRoleArn,
						MinorVersionForGetSecret,
						r.AWSClient.FindRoleARNs,
					)
			}
			if interactive.Enabled() {
				prefix, err := interactive.GetString(interactive.Input{
					Question:   "Prefix for OIDC",
					Help:       cmd.Flags().Lookup(userPrefixFlag).Usage,
					Default:    args.userPrefix,
					Validators: []interactive.Validator{interactive.MaxLength(maxLengthUserPrefix)},
				})
				if err != nil {
					r.Reporter.Errorf("Expected a valid prefix for the configuration: %s", err)
					os.Exit(1)
				}
				args.userPrefix = prefix
			}
			roleName, _ := aws.GetResourceIdFromARN(args.installerRoleArn)
			if roleName != "" {
				if !output.HasFlag() && r.Reporter.IsTerminal() && mode == interactive.ModeAuto {
					r.Reporter.Infof("Using %s for the installer role", args.installerRoleArn)
				}
				err := aws.ARNValidator(args.installerRoleArn)
				if err != nil {
					r.Reporter.Errorf("Expected a valid ARN: %s", err)
					os.Exit(1)
				}
				roleExists, _, err := r.AWSClient.CheckRoleExists(roleName)
				if err != nil {
					r.Reporter.Errorf(
						"There was a problem checking if role '%s' exists: %v",
						args.installerRoleArn,
						err,
					)
					os.Exit(1)
				}
				if !roleExists {
					r.Reporter.Errorf("Role '%s' does not exist", args.installerRoleArn)
					os.Exit(1)
				}
				isValid, err := r.AWSClient.ValidateAccountRoleVersionCompatibility(
					roleName, aws.InstallerAccountRole, MinorVersionForGetSecret)
				if err != nil {
					r.Reporter.Errorf("There was a problem listing role tags: %v", err)
					os.Exit(1)
				}
				if !isValid {
					r.Reporter.Errorf(
						"Role '%s' is not of minimum version '%s'",
						args.installerRoleArn,
						MinorVersionForGetSecret,
					)
					os.Exit(1)
				}
			}
		}

		args.userPrefix = strings.Trim(args.userPrefix, " \t")

		if len([]rune(args.userPrefix)) > maxLengthUserPrefix {
			r.Reporter.Errorf("Expected a valid prefix for the configuration: "+
				"length of prefix is limited to %d characters", maxLengthUserPrefix)
			os.Exit(1)
		}
	}

	oidcConfigInput := oidcconfigs.OidcConfigInput{}
	if !args.managed {
		oidcConfigInput, err = oidcconfigs.BuildOidcConfigInput(args.userPrefix, args.region)
		if err != nil {
			r.Reporter.Errorf("%s", err)
			os.Exit(1)
		}
	}

	oidcConfigStrategy, err := getOidcConfigStrategy(mode, &oidcConfigInput)
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}
	oidcConfigId := oidcConfigStrategy.execute(r)
	if !args.rawFiles {
		arguments.DisableRegionDeprecationWarning = true // disable region deprecation warning
		providerArgs := []string{"", mode, oidcConfigInput.IssuerUrl}
		if oidcConfigId != "" {
			providerArgs = append(providerArgs, "--oidc-config-id", oidcConfigId)
			err = oidcprovider.Cmd.Flags().Set("oidc-config-id", oidcConfigId)
			if err != nil {
				r.Reporter.Errorf("Unable to attempt creation of OIDC provider; oidc config ID"+
					" not found / not created successfully: %s", err)
				os.Exit(1)
			}
		} else {
			r.Reporter.Infof("To create the OIDC provider, please run 'rosa create oidc-provider' with the ID " +
				"of the OIDC config or cluster you want to associate it with.")
			os.Exit(0)
		}
		oidcprovider.Cmd.Run(oidcprovider.Cmd, providerArgs)
		arguments.DisableRegionDeprecationWarning = false // enable region deprecation again
	}
}

func CreateOIDCConfig(r *rosa.Runtime, managed bool, userPrefix, region string) (string, error) {
	// userPrefix, region are used only for unmanaged
	oidcConfigInput, err := oidcconfigs.BuildOidcConfigInput(userPrefix, region)
	if err != nil {
		return "", nil
	}
	if managed {
		strategy := CreateManagedOidcConfigAutoStrategy{oidcConfigInput: &oidcConfigInput}
		return strategy.executeNoExit(r)
	}

	strategy := CreateUnmanagedOidcConfigAutoStrategy{oidcConfig: &oidcConfigInput}
	return strategy.executeNoExit(r)
}

type CreateOidcConfigStrategy interface {
	execute(r *rosa.Runtime) string
}

type CreateUnmanagedOidcConfigRawStrategy struct {
	oidcConfig *oidcconfigs.OidcConfigInput
}

func (s *CreateUnmanagedOidcConfigRawStrategy) execute(r *rosa.Runtime) string {
	bucketName := s.oidcConfig.BucketName
	discoveryDocument := s.oidcConfig.DiscoveryDocument
	jwks := s.oidcConfig.Jwks
	privateKey := s.oidcConfig.PrivateKey
	privateKeyFilename := s.oidcConfig.PrivateKeyFilename
	err := helper.SaveDocument(string(privateKey), privateKeyFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving private key to a file: %s", err)
		os.Exit(1)
	}
	discoveryDocumentFilename := fmt.Sprintf("discovery-document-%s.json", bucketName)
	err = helper.SaveDocument(discoveryDocument, discoveryDocumentFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving discovery document to a file: %s", err)
		os.Exit(1)
	}
	jwksFilename := fmt.Sprintf("jwks-%s.json", bucketName)
	err = helper.SaveDocument(string(jwks[:]), jwksFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving JSON Web Key Set to a file: %s", err)
		os.Exit(1)
	}
	if !output.HasFlag() && r.Reporter.IsTerminal() {
		r.Reporter.Infof(
			"Please refer to documentation to use generated files to create an OIDC compliant configuration.",
		)
	}
	return ""
}

type CreateUnmanagedOidcConfigAutoStrategy struct {
	oidcConfig *oidcconfigs.OidcConfigInput
}

const (
	discoveryDocumentKey = ".well-known/openid-configuration"
	jwksKey              = "keys.json"
)

func (s *CreateUnmanagedOidcConfigAutoStrategy) executeNoExit(r *rosa.Runtime) (string, error) {
	bucketUrl := s.oidcConfig.IssuerUrl
	bucketName := s.oidcConfig.BucketName
	discoveryDocument := s.oidcConfig.DiscoveryDocument
	jwks := s.oidcConfig.Jwks
	privateKey := s.oidcConfig.PrivateKey
	privateKeySecretName := s.oidcConfig.PrivateKeySecretName
	installerRoleArn := args.installerRoleArn
	err := r.AWSClient.CreateS3Bucket(bucketName, args.region)
	if err != nil {
		return "", fmt.Errorf("There was a problem creating S3 bucket '%s': %s", bucketName, err)
	}
	err = r.AWSClient.PutPublicReadObjectInS3Bucket(bucketName, strings.NewReader(discoveryDocument), discoveryDocumentKey)
	if err != nil {
		return "", fmt.Errorf("There was a problem populating discovery "+
			"document to S3 bucket '%s': %s", bucketName, err)
	}
	err = r.AWSClient.PutPublicReadObjectInS3Bucket(bucketName, bytes.NewReader(jwks), jwksKey)
	if err != nil {
		return "", fmt.Errorf("There was a problem populating JWKS "+
			"to S3 bucket '%s': %s", bucketName, err)
	}
	secretARN, err := r.AWSClient.CreateSecretInSecretsManager(privateKeySecretName, string(privateKey[:]))
	if err != nil {
		return "", fmt.Errorf("There was a problem saving private key to secrets manager: %s", err)
	}
	oidcConfig, err := v1.NewOidcConfig().
		Managed(false).
		SecretArn(secretARN).
		IssuerUrl(bucketUrl).
		InstallerRoleArn(installerRoleArn).
		Build()
	if err == nil {
		oidcConfig, err = r.OCMClient.CreateOidcConfig(oidcConfig)
	}
	if err != nil {
		return "", fmt.Errorf("There was a problem building your unmanaged OIDC Configuration: %v.\n"+
			"Please refer to documentation and try again through:\n"+
			"\trosa register oidc-config --issuer-url %s --secret-arn %s --role-arn %s",
			err, bucketUrl, secretARN, installerRoleArn)
	}
	return oidcConfig.ID(), nil
}

func (s *CreateUnmanagedOidcConfigAutoStrategy) execute(r *rosa.Runtime) string {
	bucketUrl := s.oidcConfig.IssuerUrl
	bucketName := s.oidcConfig.BucketName
	discoveryDocument := s.oidcConfig.DiscoveryDocument
	jwks := s.oidcConfig.Jwks
	privateKey := s.oidcConfig.PrivateKey
	privateKeySecretName := s.oidcConfig.PrivateKeySecretName
	installerRoleArn := args.installerRoleArn
	var spin *spinner.Spinner
	if !output.HasFlag() && r.Reporter.IsTerminal() {
		spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		r.Reporter.Infof("Setting up unmanaged OIDC configuration '%s'", bucketName)
	}
	if spin != nil {
		spin.Start()
	}
	err := r.AWSClient.CreateS3Bucket(bucketName, args.region)
	if err != nil {
		r.Reporter.Errorf("There was a problem creating S3 bucket '%s': %s", bucketName, err)
		os.Exit(1)
	}
	err = r.AWSClient.PutPublicReadObjectInS3Bucket(
		bucketName, strings.NewReader(discoveryDocument), discoveryDocumentKey)
	if err != nil {
		r.Reporter.Errorf("There was a problem populating discovery "+
			"document to S3 bucket '%s': %s", bucketName, err)
		os.Exit(1)
	}
	err = r.AWSClient.PutPublicReadObjectInS3Bucket(bucketName, bytes.NewReader(jwks), jwksKey)
	if err != nil {
		if spin != nil {
			spin.Stop()
		}
		r.Reporter.Errorf("There was a problem populating JWKS "+
			"to S3 bucket '%s': %s", bucketName, err)
		os.Exit(1)
	}
	secretARN, err := r.AWSClient.CreateSecretInSecretsManager(privateKeySecretName, string(privateKey[:]))
	if err != nil {
		r.Reporter.Errorf("There was a problem saving private key to secrets manager: %s", err)
		os.Exit(1)
	}
	oidcConfig, err := v1.NewOidcConfig().
		Managed(false).
		SecretArn(secretARN).
		IssuerUrl(bucketUrl).
		InstallerRoleArn(installerRoleArn).
		Build()
	if err == nil {
		oidcConfig, err = r.OCMClient.CreateOidcConfig(oidcConfig)
	}
	if err != nil {
		if spin != nil {
			spin.Stop()
		}
		r.Reporter.Errorf("There was a problem building your unmanaged OIDC Configuration %v.\n"+
			"Please refer to documentation and try again through:\n"+
			"\trosa register oidc-config --issuer-url %s --secret-arn %s --role-arn %s",
			err, bucketUrl, secretARN, installerRoleArn)
		os.Exit(1)
	}
	if output.HasFlag() {
		err = output.Print(oidcConfig)
		if err != nil {
			r.Reporter.Errorf("%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if r.Reporter.IsTerminal() {
		if spin != nil {
			spin.Stop()
		}
		output := fmt.Sprintf(InformOperatorRolesOutput, oidcConfig.ID())
		r.Reporter.Infof(output)
	}
	return oidcConfig.ID()
}

type CreateUnmanagedOidcConfigManualStrategy struct {
	oidcConfig *oidcconfigs.OidcConfigInput
}

func (s *CreateUnmanagedOidcConfigManualStrategy) execute(r *rosa.Runtime) string {
	commands := []string{}
	bucketName := s.oidcConfig.BucketName
	discoveryDocument := s.oidcConfig.DiscoveryDocument
	jwks := s.oidcConfig.Jwks
	privateKey := s.oidcConfig.PrivateKey
	privateKeyFilename := s.oidcConfig.PrivateKeyFilename
	privateKeySecretName := s.oidcConfig.PrivateKeySecretName
	err := helper.SaveDocument(string(privateKey), privateKeyFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving private key to a file: %s", err)
		os.Exit(1)
	}
	createBucketConfig := ""
	if args.region != aws.DefaultRegion {
		createBucketConfig = fmt.Sprintf("LocationConstraint=%s", args.region)
	}
	createS3BucketCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.CreateBucket).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.CreateBucketConfiguration, createBucketConfig).
		AddParam(awscb.Region, args.region).
		Build()
	commands = append(commands, createS3BucketCommand)

	putBucketTaggingCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.PutBucketTagging).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.Tagging, fmt.Sprintf("'TagSet=[{Key=%s,Value=%s}]'", tags.RedHatManaged, tags.True)).
		Build()
	commands = append(commands, putBucketTaggingCommand)

	PutPublicAccessBlockCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.PutPublicAccessBlock).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.PublicAccessBlockConfiguration,
			"BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=false,RestrictPublicBuckets=false").
		Build()
	commands = append(commands, PutPublicAccessBlockCommand)

	readOnlyPolicyFilename := fmt.Sprintf("readOnlyPolicy-%s.json", bucketName)
	err = helper.SaveDocument(fmt.Sprintf(aws.ReadOnlyAnonUserPolicyTemplate, bucketName), readOnlyPolicyFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving bucket policy document to a file: %s", err)
		os.Exit(1)
	}
	putBucketBucketPolicyCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.PutBucketPolicy).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.Policy, fmt.Sprintf("file://%s", readOnlyPolicyFilename)).
		Build()
	commands = append(commands, putBucketBucketPolicyCommand)
	commands = append(commands, fmt.Sprintf("rm %s", readOnlyPolicyFilename))

	discoveryDocumentFilename := fmt.Sprintf("discovery-document-%s.json", bucketName)
	err = helper.SaveDocument(discoveryDocument, discoveryDocumentFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving discovery document to a file: %s", err)
		os.Exit(1)
	}
	putDiscoveryDocumentCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.PutObject).
		AddParam(awscb.Body, fmt.Sprintf("./%s", discoveryDocumentFilename)).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.Key, discoveryDocumentKey).
		AddParam(awscb.Tagging, fmt.Sprintf("'%s=%s'", tags.RedHatManaged, tags.True)).
		Build()
	commands = append(commands, putDiscoveryDocumentCommand)
	commands = append(commands, fmt.Sprintf("rm %s", discoveryDocumentFilename))
	jwksFilename := fmt.Sprintf("jwks-%s.json", bucketName)
	err = helper.SaveDocument(string(jwks[:]), jwksFilename)
	if err != nil {
		r.Reporter.Errorf("There was a problem saving JSON Web Key Set to a file: %s", err)
		os.Exit(1)
	}
	putJwksCommand := awscb.NewS3ApiCommandBuilder().
		SetCommand(awscb.PutObject).
		AddParam(awscb.Body, fmt.Sprintf("./%s", jwksFilename)).
		AddParam(awscb.Bucket, bucketName).
		AddParam(awscb.Key, jwksKey).
		AddParam(awscb.Tagging, fmt.Sprintf("'%s=%s'", tags.RedHatManaged, tags.True)).
		Build()
	commands = append(commands, putJwksCommand)
	commands = append(commands, fmt.Sprintf("rm %s", jwksFilename))
	createSecretCommand := awscb.NewSecretsManagerCommandBuilder().
		SetCommand(awscb.CreateSecret).
		AddParam(awscb.Name, privateKeySecretName).
		AddParam(awscb.SecretString, fmt.Sprintf("file://%s", privateKeyFilename)).
		AddParam(awscb.Description, fmt.Sprintf("\"Secret for %s\"", bucketName)).
		AddParam(awscb.Region, args.region).
		AddTags(map[string]string{
			tags.RedHatManaged: "true",
		}).
		Build()
	commands = append(commands, createSecretCommand)
	commands = append(commands, fmt.Sprintf("rm %s", privateKeyFilename))
	fmt.Println(awscb.JoinCommands(commands))
	if r.Reporter.IsTerminal() {
		r.Reporter.Infof("Please run commands above to generate OIDC compliant configuration in your AWS account. " +
			"To register this OIDC Configuration, please run the following command:\n" +
			"rosa register oidc-config\n" +
			"For more information please refer to the documentation")
	}
	return ""
}

type CreateManagedOidcConfigAutoStrategy struct {
	oidcConfigInput *oidcconfigs.OidcConfigInput
}

func (s *CreateManagedOidcConfigAutoStrategy) execute(r *rosa.Runtime) string {
	var spin *spinner.Spinner
	if !output.HasFlag() && r.Reporter.IsTerminal() {
		spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		r.Reporter.Infof("Setting up managed OIDC configuration")
	}
	if spin != nil {
		spin.Start()
	}
	oidcConfig, err := v1.NewOidcConfig().Managed(true).Build()
	if err != nil {
		r.Reporter.Errorf("There was a problem building the managed OIDC Configuration: %v", err)
		os.Exit(1)
	}
	oidcConfig, err = r.OCMClient.CreateOidcConfig(oidcConfig)
	if err != nil {
		if spin != nil {
			spin.Stop()
		}
		r.Reporter.Errorf("There was a problem registering your managed OIDC Configuration: %v", err)
		os.Exit(1)
	}
	s.oidcConfigInput.IssuerUrl = oidcConfig.IssuerUrl()
	if output.HasFlag() {
		err = output.Print(oidcConfig)
		if err != nil {
			r.Reporter.Errorf("%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if r.Reporter.IsTerminal() {
		if spin != nil {
			spin.Stop()
		}
		output := fmt.Sprintf(InformOperatorRolesOutput, oidcConfig.ID())
		r.Reporter.Infof(output)
	}
	return oidcConfig.ID()
}

func (s *CreateManagedOidcConfigAutoStrategy) executeNoExit(r *rosa.Runtime) (string, error) {
	oidcConfig, err := v1.NewOidcConfig().Managed(true).Build()
	if err != nil {
		return "", fmt.Errorf("There was a problem building the managed OIDC Configuration: %v", err)
	}
	oidcConfig, err = r.OCMClient.CreateOidcConfig(oidcConfig)
	if err != nil {
		return "", fmt.Errorf("There was a problem building the managed OIDC Configuration: %v", err)
	}
	s.oidcConfigInput.IssuerUrl = oidcConfig.IssuerUrl()
	return oidcConfig.ID(), nil
}

func getOidcConfigStrategy(mode string, input *oidcconfigs.OidcConfigInput) (CreateOidcConfigStrategy, error) {
	if args.rawFiles {
		return &CreateUnmanagedOidcConfigRawStrategy{oidcConfig: input}, nil
	}
	if args.managed {
		return &CreateManagedOidcConfigAutoStrategy{oidcConfigInput: input}, nil
	}
	switch mode {
	case interactive.ModeAuto:
		return &CreateUnmanagedOidcConfigAutoStrategy{oidcConfig: input}, nil
	case interactive.ModeManual:
		return &CreateUnmanagedOidcConfigManualStrategy{oidcConfig: input}, nil
	default:
		return nil, weberr.Errorf("Invalid mode. Allowed values are %s", interactive.Modes)
	}
}
