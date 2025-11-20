package operatorroles

import (
	"fmt"
	"os"
	"strings"

	common "github.com/openshift-online/ocm-common/pkg/aws/validations"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/openshift/rosa/pkg/aws"
	awscb "github.com/openshift/rosa/pkg/aws/commandbuilder"
	"github.com/openshift/rosa/pkg/aws/tags"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/output"
	"github.com/openshift/rosa/pkg/roles"
	"github.com/openshift/rosa/pkg/rosa"
)

type operatorRolesInput struct {
	prefix              string
	permissionsBoundary string
	cluster             *cmv1.Cluster
	accountRoleVersion  string
	policies            map[string]*cmv1.AWSSTSPolicy
	defaultVersion      string
	credRequests        map[string]*cmv1.STSOperator
	managedPolicies     bool
	hostedCPPolicies    bool
	isHcpSharedVpc      bool
	route53RoleArn      string
	vpcEndpointRoleArn  string
}

func handleOperatorRoleCreationByClusterKey(r *rosa.Runtime, env string,
	permissionsBoundary string, mode string,
	policies map[string]*cmv1.AWSSTSPolicy,
	defaultPolicyVersion string, isHcpSharedVpc bool) error {
	clusterKey := r.GetClusterKey()
	cluster := r.FetchCluster()
	route53RoleArn := args.sharedVpcRoleArn
	vpcEndpointRoleArn := args.vpcEndpointRoleArn
	if cluster.AWS().STS().RoleARN() == "" {
		r.Reporter.Errorf("Cluster '%s' is not an STS cluster.", clusterKey)
		os.Exit(1)
	}

	// Check to see if IAM operator roles have already created
	missingRoles, err := validateOperatorRoles(r, cluster)
	if err != nil {
		if strings.Contains(err.Error(), "AccessDenied") {
			r.Reporter.Debugf("Failed to verify if operator roles exist: '%v'", err)
		} else {
			r.Reporter.Errorf("Failed to verify if operator roles exist: '%v'", err)
			os.Exit(1)
		}
	}

	hostedCPPolicies := aws.IsHostedCPManagedPolicies(cluster)

	operatorRolePolicyPrefix, err := aws.GetOperatorRolePolicyPrefixFromCluster(cluster, r.AWSClient)
	if err != nil {
		r.Reporter.Errorf("%s", err)
	}

	credRequests, err := r.OCMClient.GetCredRequests(cluster.Hypershift().Enabled())
	if err != nil {
		r.Reporter.Errorf("Error getting operator credential request from OCM '%v'", err)
		os.Exit(1)
	}

	managedPolicies := cluster.AWS().STS().ManagedPolicies()
	if args.forcePolicyCreation && managedPolicies {
		r.Reporter.Warnf("Forcing creation of policies only works for unmanaged policies")
		os.Exit(1)
	}

	switch mode {
	case interactive.ModeAuto:

		if len(missingRoles) == 0 {
			if ocm.IsOidcConfigReusable(cluster) {
				err := validateOperatorRolesMatchOidcProvider(r, cluster)
				if err != nil {
					return err
				}
			}

			if !args.forcePolicyCreation {
				r.Reporter.Infof("Operator Roles already exists")
				return nil
			}
		}
		roleName, err := aws.GetInstallerAccountRoleName(cluster)
		if err != nil {
			r.Reporter.Errorf("Expected parsing role account role '%s': '%v'", cluster.AWS().STS().RoleARN(), err)
			os.Exit(1)
		}

		path, err := aws.GetPathFromAccountRole(cluster, aws.AccountRoles[aws.InstallerAccountRole].Name)
		if err != nil {
			r.Reporter.Errorf("Expected a valid path for '%s': '%v'", cluster.AWS().STS().RoleARN(), err)
			os.Exit(1)
		}
		if path != "" && !output.HasFlag() && r.Reporter.IsTerminal() {
			r.Reporter.Infof("ARN path '%s' detected in installer role '%s'. "+
				"This ARN path will be used for subsequent created operator roles and policies.",
				path, cluster.AWS().STS().RoleARN())
		}
		var accountRoleVersion string

		if !output.HasFlag() || r.Reporter.IsTerminal() {
			r.Reporter.Infof("Creating roles using '%s'", r.Creator.ARN)
		}
		accountRoleVersion, err = r.AWSClient.GetAccountRoleVersion(roleName)
		if err != nil {
			r.Reporter.Errorf("Error getting account role version '%v'", err)
			os.Exit(1)
		}
		err = createRoles(r, operatorRolesInput{
			prefix:              operatorRolePolicyPrefix,
			permissionsBoundary: permissionsBoundary,
			cluster:             cluster,
			accountRoleVersion:  accountRoleVersion,
			policies:            policies,
			defaultVersion:      defaultPolicyVersion,
			credRequests:        credRequests,
			managedPolicies:     managedPolicies,
			hostedCPPolicies:    hostedCPPolicies,
			isHcpSharedVpc:      isHcpSharedVpc,
			route53RoleArn:      route53RoleArn,
			vpcEndpointRoleArn:  vpcEndpointRoleArn,
		})
		if err != nil {
			r.Reporter.Errorf("There was an error creating the operator roles: '%v'", err)
			isThrottle := "false"
			if strings.Contains(err.Error(), "Throttling") {
				isThrottle = helper.True
			}
			r.OCMClient.LogEvent("ROSACreateOperatorRolesModeAuto", map[string]string{
				ocm.ClusterID:  clusterKey,
				ocm.Response:   ocm.Failure,
				ocm.IsThrottle: isThrottle,
			})
			os.Exit(1)
		}
		r.OCMClient.LogEvent("ROSACreateOperatorRolesModeAuto", map[string]string{
			ocm.ClusterID: clusterKey,
			ocm.Response:  ocm.Success,
		})
	case interactive.ModeManual:
		commands, err := buildCommands(r, env, operatorRolePolicyPrefix, permissionsBoundary, defaultPolicyVersion,
			cluster, policies, credRequests, managedPolicies, hostedCPPolicies, route53RoleArn, vpcEndpointRoleArn)
		if err != nil {
			r.Reporter.Errorf("There was an error building the list of resources: '%v'", err)
			os.Exit(1)
			r.OCMClient.LogEvent("ROSACreateOperatorRolesModeManual", map[string]string{
				ocm.ClusterID: clusterKey,
				ocm.Response:  ocm.Failure,
			})
		}
		if r.Reporter.IsTerminal() {
			r.Reporter.Infof("All policy files saved to the current directory")
			r.Reporter.Infof("Run the following commands to create the operator roles:\n")
		}
		r.OCMClient.LogEvent("ROSACreateOperatorRolesModeManual", map[string]string{
			ocm.ClusterID: clusterKey,
		})
		fmt.Println(commands)

	default:
		r.Reporter.Errorf("Invalid mode. Allowed values are '%s'", interactive.Modes)
		os.Exit(1)
	}
	return nil
}

func createRoles(r *rosa.Runtime, createInput operatorRolesInput) error {
	sharedVpcRoleArn := createInput.cluster.AWS().PrivateHostedZoneRoleARN()
	isSharedVpc := sharedVpcRoleArn != ""

	for credrequest, operator := range createInput.credRequests {
		ver := createInput.cluster.Version()
		if ver != nil && operator.MinVersion() != "" {
			isSupported, err := ocm.CheckSupportedVersion(ocm.GetVersionMinor(ver.ID()), operator.MinVersion())
			if err != nil {
				r.Reporter.Errorf("Error validating operator role '%s' version %s", operator.Name(), err)
				os.Exit(1)
			}
			if !isSupported {
				continue
			}
		}
		roleName, _ := aws.FindOperatorRoleNameBySTSOperator(createInput.cluster, operator)
		if roleName == "" {
			return fmt.Errorf("Failed to find operator IAM role")
		}

		path, err := aws.GetPathFromAccountRole(createInput.cluster, aws.AccountRoles[aws.InstallerAccountRole].Name)
		if err != nil {
			return err
		}

		var policyArn string
		var policyArns []string
		filename := aws.GetOperatorPolicyKey(credrequest, createInput.hostedCPPolicies, isSharedVpc)
		if createInput.managedPolicies {
			policyArn, err = aws.GetManagedPolicyARN(createInput.policies, filename)
			if err != nil {
				return err
			}
			if createInput.isHcpSharedVpc {
				if credrequest == aws.IngressOperatorCloudCredentialsRoleType {
					sharedVpcPolicyArn, err := getHcpSharedVpcPolicy(r, sharedVpcRoleArn, createInput.defaultVersion)
					if err != nil {
						return err
					}
					policyArns = append(policyArns, sharedVpcPolicyArn)
				} else if credrequest == aws.ControlPlaneCloudCredentialsRoleType {
					for _, arn := range []string{createInput.vpcEndpointRoleArn, sharedVpcRoleArn} {
						sharedVpcPolicyArn, err := getHcpSharedVpcPolicy(r, arn, createInput.defaultVersion)
						if err != nil {
							return err
						}
						policyArns = append(policyArns, sharedVpcPolicyArn)
					}
				}
			}
		} else {
			policyArn = aws.GetOperatorPolicyARN(r.Creator.Partition, r.Creator.AccountID, createInput.prefix,
				operator.Namespace(), operator.Name(), path)
			policyDetails := aws.GetPolicyDetails(createInput.policies, filename)

			if isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType {
				err = validateIngressOperatorPolicyOverride(r, policyArn, sharedVpcRoleArn, createInput.prefix)
				if err != nil {
					return err
				}

				policyDetails = aws.InterpolatePolicyDocument(r.Creator.Partition, policyDetails, map[string]string{
					"shared_vpc_role_arn": sharedVpcRoleArn,
				})
			}

			operatorPolicyTags := map[string]string{
				common.OpenShiftVersion: createInput.accountRoleVersion,
				tags.RolePrefix:         createInput.prefix,
				tags.RedHatManaged:      helper.True,
				tags.OperatorNamespace:  operator.Namespace(),
				tags.OperatorName:       operator.Name(),
			}

			if args.forcePolicyCreation || (isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType) {
				policyArn, err = r.AWSClient.ForceEnsurePolicy(policyArn, policyDetails,
					createInput.defaultVersion, operatorPolicyTags, path)
				if err != nil {
					return err
				}
			} else {
				policyArn, err = r.AWSClient.EnsurePolicy(policyArn, policyDetails,
					createInput.defaultVersion, operatorPolicyTags, path)
				if err != nil {
					return err
				}
			}
		}
		policyArns = append(policyArns, policyArn)

		policyDetails := aws.GetPolicyDetails(createInput.policies, "operator_iam_role_policy")
		policy, err := aws.GenerateOperatorRolePolicyDoc(r.Creator.Partition, createInput.cluster,
			r.Creator.AccountID, operator, policyDetails)
		if err != nil {
			return err
		}

		r.Reporter.Debugf("Creating role '%s'", roleName)
		tagsList := map[string]string{
			tags.OperatorNamespace: operator.Namespace(),
			tags.OperatorName:      operator.Name(),
			tags.RedHatManaged:     helper.True,
		}
		if !ocm.IsOidcConfigReusable(createInput.cluster) {
			tagsList[tags.ClusterID] = createInput.cluster.ID()
		}
		if createInput.managedPolicies {
			tagsList[common.ManagedPolicies] = helper.True
		}
		if createInput.hostedCPPolicies {
			tagsList[tags.HypershiftPolicies] = helper.True
		}

		roleARN, err := r.AWSClient.EnsureRole(r.Reporter, roleName, policy, createInput.permissionsBoundary,
			createInput.accountRoleVersion, tagsList, path, createInput.managedPolicies)
		if err != nil {
			return err
		}
		if !output.HasFlag() || r.Reporter.IsTerminal() {
			r.Reporter.Infof("Created role '%s' with ARN '%s'", roleName, roleARN)
		}

		for _, arn := range policyArns {
			r.Reporter.Debugf("Attaching permission policy '%s' to role '%s'", arn, roleName)
			err = r.AWSClient.AttachRolePolicy(r.Reporter, roleName, arn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func buildCommands(r *rosa.Runtime, env string,
	prefix string, permissionsBoundary string, defaultPolicyVersion string, cluster *cmv1.Cluster,
	policies map[string]*cmv1.AWSSTSPolicy, credRequests map[string]*cmv1.STSOperator,
	managedPolicies bool, hostedCPPolicies bool, route53RoleArn string, vpcEndpointRoleArn string) (string, error) {
	sharedVpcRoleArn := cluster.AWS().PrivateHostedZoneRoleARN()
	isSharedVpc := sharedVpcRoleArn != ""
	var policyDetails = make(map[string]roles.ManualSharedVpcPolicyDetails)

	if !managedPolicies {
		err := aws.GenerateOperatorRolePolicyFiles(r.Reporter, policies, credRequests, sharedVpcRoleArn, r.Creator.Partition)
		if err != nil {
			r.Reporter.Errorf("There was an error generating the policy files: %s", err)
			os.Exit(1)
		}
	}

	commands := []string{}

	for credrequest, operator := range credRequests {
		ver := cluster.Version()
		if ver != nil && operator.MinVersion() != "" {
			isSupported, err := ocm.CheckSupportedVersion(ocm.GetVersionMinor(ver.ID()), operator.MinVersion())
			if err != nil {
				r.Reporter.Errorf("Error validating operator role '%s' version %s", operator.Name(), err)
				os.Exit(1)
			}
			if !isSupported {
				continue
			}
		}
		roleName, _ := aws.FindOperatorRoleNameBySTSOperator(cluster, operator)
		path, err := aws.GetPathFromAccountRole(cluster, aws.AccountRoles[aws.InstallerAccountRole].Name)
		if err != nil {
			return "", err
		}

		var policyARN string
		if managedPolicies {
			policyARN, err = aws.GetManagedPolicyARN(policies, aws.GetOperatorPolicyKey(
				credrequest, hostedCPPolicies, isSharedVpc))
			if err != nil {
				return "", err
			}
		} else {
			policyARN = computePolicyARN(*r.Creator, prefix, operator.Namespace(), operator.Name(), path)
			name := aws.GetOperatorPolicyName(prefix, operator.Namespace(), operator.Name())
			iamTags := map[string]string{
				common.OpenShiftVersion: defaultPolicyVersion,
				tags.RolePrefix:         prefix,
				tags.OperatorNamespace:  operator.Namespace(),
				tags.OperatorName:       operator.Name(),
				tags.RedHatManaged:      helper.True,
			}
			operatorPolicyKey := aws.GetOperatorPolicyKey(credrequest, hostedCPPolicies, isSharedVpc)
			fileName := fmt.Sprintf("file://%s.json", operatorPolicyKey)
			_, err = r.AWSClient.IsPolicyExists(policyARN)
			if err != nil {
				createPolicy := awscb.NewIAMCommandBuilder().
					SetCommand(awscb.CreatePolicy).
					AddParam(awscb.PolicyName, name).
					AddParam(awscb.PolicyDocument, fileName).
					AddTags(iamTags).
					AddParam(awscb.Path, path).
					Build()
				commands = append(commands, createPolicy)
			} else if isSharedVpc && credrequest == aws.IngressOperatorCloudCredentialsRoleType {
				err := validateIngressOperatorPolicyOverride(r, policyARN, sharedVpcRoleArn, prefix)
				if err != nil {
					return "", err
				}

				createPolicyVersion := awscb.NewIAMCommandBuilder().
					SetCommand(awscb.CreatePolicyVersion).
					AddParam(awscb.PolicyArn, policyARN).
					AddParam(awscb.PolicyDocument, fileName).
					AddParamNoValue(awscb.SetAsDefault).
					Build()
				commands = append(commands, createPolicyVersion)
			}
		}

		policyDetail := aws.GetPolicyDetails(policies, "operator_iam_role_policy")
		policy, err := aws.GenerateOperatorRolePolicyDoc(r.Creator.Partition, cluster,
			r.Creator.AccountID, operator, policyDetail)
		if err != nil {
			return "", err
		}

		filename := fmt.Sprintf("operator_%s_policy", credrequest)
		filename = aws.GetFormattedFileName(filename)
		r.Reporter.Debugf("Saving '%s' to the current directory", filename)
		err = helper.SaveDocument(policy, filename)
		if err != nil {
			return "", err
		}
		iamTags := map[string]string{
			tags.OperatorNamespace: operator.Namespace(),
			tags.OperatorName:      operator.Name(),
			tags.RedHatManaged:     helper.True,
		}
		if !ocm.IsOidcConfigReusable(cluster) {
			iamTags[tags.ClusterID] = cluster.ID()
		}
		if managedPolicies {
			iamTags[common.ManagedPolicies] = helper.True
		}
		if hostedCPPolicies {
			iamTags[tags.HypershiftPolicies] = helper.True
		}
		createRole := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.CreateRole).
			AddParam(awscb.RoleName, roleName).
			AddParam(awscb.AssumeRolePolicyDocument, fmt.Sprintf("file://%s", filename)).
			AddParam(awscb.PermissionsBoundary, permissionsBoundary).
			AddTags(iamTags).
			AddParam(awscb.Path, path).
			Build()

		attachRolePolicy := awscb.NewIAMCommandBuilder().
			SetCommand(awscb.AttachRolePolicy).
			AddParam(awscb.RoleName, roleName).
			AddParam(awscb.PolicyArn, policyARN).
			Build()

		var attachSharedVpcRolePolicy string
		var policyCommands []string

		if isSharedVpc { // HCP Shared VPC policy attachment

			// Precreate HCP shared VPC policies for less memory usage + time to execute
			// Shared VPC role arn (route53)
			if _, ok := policyDetails[aws.IngressOperatorCloudCredentialsRoleType]; !ok {
				exists, createPolicyCommand, policyName, err := roles.GetHcpSharedVpcPolicyDetails(r, sharedVpcRoleArn)
				if err != nil {
					return "", err
				}

				sharedVpcRolePath, err := aws.GetPathFromARN(sharedVpcRoleArn)
				if err != nil {
					return "", err
				}

				policyDetails[aws.IngressOperatorCloudCredentialsRoleType] = roles.ManualSharedVpcPolicyDetails{
					Command:       createPolicyCommand,
					Name:          policyName,
					AlreadyExists: exists,
					Path:          sharedVpcRolePath,
				}
			}
			// VPC endpoint role arn
			if _, ok := policyDetails[aws.ControlPlaneCloudCredentialsRoleType]; !ok {

				exists, createPolicyCommand, policyName, err := roles.GetHcpSharedVpcPolicyDetails(r, vpcEndpointRoleArn)
				if err != nil {
					return "", err
				}

				vpcEndpointRolePath, err := aws.GetPathFromARN(vpcEndpointRoleArn)
				if err != nil {
					return "", err
				}

				policyDetails[aws.ControlPlaneCloudCredentialsRoleType] = roles.ManualSharedVpcPolicyDetails{
					Command:       createPolicyCommand,
					Name:          policyName,
					AlreadyExists: exists,
					Path:          vpcEndpointRolePath,
				}
			}

			var policies []string

			// Attach HCP shared VPC policies
			switch credrequest {
			case aws.IngressOperatorCloudCredentialsRoleType:
				if details, ok := policyDetails[credrequest]; ok {
					policies = append(policies, policyDetails[credrequest].Name)
					if !policyDetails[credrequest].AlreadyExists { // Skip creation if already exists
						policyCommands = append(policyCommands, policyDetails[credrequest].Command)
						// Allow only one creation command for this policy to be printed
						details.AlreadyExists = true
						policyDetails[credrequest] = details
					}
				}
			case aws.ControlPlaneCloudCredentialsRoleType:
				for i, details := range policyDetails {
					policies = append(policies, details.Name)
					if !details.AlreadyExists {
						policyCommands = append(policyCommands, details.Command)
						// Allow only one creation command for this policy to be printed
						details.AlreadyExists = true
						policyDetails[i] = details
					}
				}
			}

			// Attach policies to roles
			for _, policy := range policies {
				details, err := roles.GetPolicyDetailsByName(policyDetails, policy)
				if err != nil {
					return "", err
				}
				arn := aws.GetPolicyArn(r.Creator.Partition, r.Creator.AccountID, policy, details.Path)

				attachSharedVpcRolePolicy = awscb.NewIAMCommandBuilder().
					SetCommand(awscb.AttachRolePolicy).
					AddParam(awscb.RoleName, roleName).
					AddParam(awscb.PolicyArn, arn).
					Build()
				policyCommands = append(policyCommands, attachSharedVpcRolePolicy)
			}
		}
		commands = append(commands, createRole, attachRolePolicy)
		commands = append(commands, policyCommands...)

	}
	return awscb.JoinCommands(commands), nil
}

func validateOperatorRoles(r *rosa.Runtime, cluster *cmv1.Cluster) ([]string, error) {
	var missingRoles []string
	operatorIAMRoles := cluster.AWS().STS().OperatorIAMRoles()
	if len(operatorIAMRoles) == 0 {
		return missingRoles, fmt.Errorf("No Operator IAM roles found for cluster %s", cluster.Name())
	}
	for _, operatorIAMRole := range operatorIAMRoles {
		roleARN := operatorIAMRole.RoleARN()
		roleName, err := aws.GetResourceIdFromARN(roleARN)
		if err != nil {
			return missingRoles, err
		}
		exists, _, err := r.AWSClient.CheckRoleExists(roleName)
		if err != nil {
			return missingRoles, err
		}
		if !exists {
			missingRoles = append(missingRoles, roleName)
		}
	}
	return missingRoles, nil
}

func validateOperatorRolesMatchOidcProvider(r *rosa.Runtime, cluster *cmv1.Cluster) error {
	operatorRolesList, err := convertV1OperatorIAMRoleIntoOcmOperatorIamRole(
		cluster.AWS().STS().OperatorIAMRoles())
	if err != nil {
		return err
	}
	expectedPath, err := aws.GetPathFromARN(cluster.AWS().STS().RoleARN())
	if err != nil {
		return err
	}
	return ocm.ValidateOperatorRolesMatchOidcProvider(r.Reporter, r.AWSClient,
		operatorRolesList, cluster.AWS().STS().OidcConfig().IssuerUrl(),
		ocm.GetVersionMinor(cluster.Version().RawID()), expectedPath, cluster.AWS().STS().ManagedPolicies(), false)
}
