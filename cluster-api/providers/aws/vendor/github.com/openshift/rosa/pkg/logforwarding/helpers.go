package logforwarding

import (
	"os"

	"gopkg.in/yaml.v3"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"
)

// FlagName contains the common log forwarding config command flag name
const FlagName = "log-fwd-config"
const LogFwdConfigHelpMessage = "A path to a log forwarding config file. This should be a YAML file with the" +
	" following structure:\n\n" +
	"cloudwatch:\n" +
	"  cloudwatch_log_role_arn: \"role_arn_here\"\n" +
	"  cloudwatch_log_group_name: \"group_name_here\"\n" +
	"  applications: [\"example_app_1\", \"example_app_2\"]\n" +
	"  groups: [\"group-name\", \"group_name-2\"]\n" +
	"s3:\n" +
	"  s3_config_bucket_name: \"bucket_name_here\"\n" +
	"  s3_config_bucket_prefix: \"bucket_prefix_here\"\n" +
	"  applications: [\"example_app_1\", \"example_app_2\"]\n" +
	"  groups: [\"group-name\"]"

// S3LogForwarderConfig represents the log forward config for S3
type S3LogForwarderConfig struct {
	Applications         []string `yaml:"applications,omitempty"`
	GroupsLogVersions    []string `yaml:"groups,omitempty"`
	S3ConfigBucketName   string   `yaml:"s3_config_bucket_name,omitempty"`
	S3ConfigBucketPrefix string   `yaml:"s3_config_bucket_prefix,omitempty"`
}

// CloudWatchLogForwarderConfig represents the log forward config for CloudWatch
type CloudWatchLogForwarderConfig struct {
	Applications           []string `yaml:"applications,omitempty"`
	GroupsLogVersions      []string `yaml:"groups,omitempty"`
	CloudWatchLogRoleArn   string   `yaml:"cloudwatch_log_role_arn,omitempty"`
	CloudWatchLogGroupName string   `yaml:"cloudwatch_log_group_name,omitempty"`
}

type LogForwarderYaml struct {
	S3         *S3LogForwarderConfig         `yaml:"s3"`
	CloudWatch *CloudWatchLogForwarderConfig `yaml:"cloudwatch"`
}

func ConstructPodGroupsHelpMessage(options []*cmv1.LogForwarderGroupVersions) (s string) {
	s = ""
	for _, option := range options {
		apps := ""
		for i, application := range option.Versions()[len(option.Versions())-1].Applications() {
			if i != 0 {
				apps += ","
			}
			apps += application
		}
		s = s + option.Name() + ": " + apps + "\n"
	}
	return
}

func ConstructPodGroupsInteractiveOptions(options []*cmv1.LogForwarderGroupVersions) (l []string) {
	for _, option := range options {
		l = append(l, option.Name())
	}
	return
}

func BindCloudWatchLogForwarder(input *CloudWatchLogForwarderConfig) *cmv1.LogForwarderBuilder {
	cloudWatchBuilder := cmv1.NewLogForwarderCloudWatchConfig()
	cloudWatchBuilder.LogDistributionRoleArn(input.CloudWatchLogRoleArn)
	cloudWatchBuilder.LogGroupName(input.CloudWatchLogGroupName)
	outputBuilder := cmv1.NewLogForwarder().Cloudwatch(cloudWatchBuilder)
	outputBuilder.Applications(input.Applications...)
	if len(input.GroupsLogVersions) > 0 {
		logForwarderGroups := make([]*cmv1.LogForwarderGroupBuilder, 0)
		for _, group := range input.GroupsLogVersions {
			logForwarderGroups = append(logForwarderGroups, cmv1.NewLogForwarderGroup().ID(group))
		}
		outputBuilder.Groups(logForwarderGroups...)
	}
	return outputBuilder
}

func BindS3LogForwarder(input *S3LogForwarderConfig) *cmv1.LogForwarderBuilder {
	s3Builder := cmv1.NewLogForwarderS3Config()
	s3Builder.BucketName(input.S3ConfigBucketName)
	s3Builder.BucketPrefix(input.S3ConfigBucketPrefix)
	outputBuilder := cmv1.NewLogForwarder().S3(s3Builder)
	outputBuilder.Applications(input.Applications...)
	if len(input.GroupsLogVersions) > 0 {
		logForwarderGroups := make([]*cmv1.LogForwarderGroupBuilder, 0)
		for _, group := range input.GroupsLogVersions {
			logForwarderGroups = append(logForwarderGroups, cmv1.NewLogForwarderGroup().ID(group))
		}
		outputBuilder.Groups(logForwarderGroups...)
	}
	return outputBuilder
}

func UnmarshalLogForwarderConfigYaml(yamlFile string) (*LogForwarderYaml, error) {
	fileContents, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, errors.UserWrapf(err, "error reading log-fwd-config YAML file '%s'", yamlFile)
	}
	tempFwdConfigObject := &LogForwarderYaml{}
	err = yaml.Unmarshal(fileContents, &tempFwdConfigObject)
	if err != nil {
		return nil, errors.UserWrapf(err, "error parsing log forwarder config YAML file '%s'", yamlFile)
	}

	return tempFwdConfigObject, nil
}
