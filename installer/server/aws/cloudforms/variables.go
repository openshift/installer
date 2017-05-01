package cloudforms

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/route53"
	fuze "github.com/coreos/container-linux-config-transpiler/config"
	templating "github.com/coreos/container-linux-config-transpiler/config/templating"
	"github.com/coreos/coreos-cloudinit/config/validate"
)

// stackVars contain variables for rendering a cloud formation stack template.
// TODO(dghubble) This type contains sooo much more than just variables. Clean
// this up.
type stackVars struct {
	*Config
	UserDataEtcd       string
	UserDataController string
	UserDataWorker     string
	// TODO(dghubble): Remove output StackBody from the variables struct
	StackBody string
	StackURL  string
}

// newStackVars returns a new StackVars.
func newStackVars(c *Config) *stackVars {
	return &stackVars{
		Config: c,
	}
}

// encodeSecretAssets encrypts the given SecretAssets and adds the compressed
// ciphertexts to the StackConfig. Requests are made to the AWS key management
// service to perform the encryption.
func (s *stackVars) encodeSecretAssets(sess *session.Session, assets *SecretAssets) error {
	kmsSvc := kms.New(sess)

	compactAssets, err := assets.compact(s.Config, kmsSvc)
	if err != nil {
		return err
	}

	s.EncodedSecrets = compactAssets
	return nil
}

// render renders the user-data for different kinds of nodes used in the stack.
func (s *stackVars) render() error {
	controller, err := renderTemplate(s.ControllerTemplate, s.Config)
	if err != nil {
		return fmt.Errorf("failed to render controller template %v", err)
	}
	controllerIgn, err := convertToIgnition(controller)
	if err != nil {
		return err
	}
	s.UserDataController = base64.StdEncoding.EncodeToString(controllerIgn)

	worker, err := renderTemplate(s.WorkerTemplate, s.Config)
	if err != nil {
		return fmt.Errorf("failed to render worker template %v", err)
	}
	workerIgn, err := convertToIgnition(worker)
	if err != nil {
		return err
	}
	s.UserDataWorker = base64.StdEncoding.EncodeToString(workerIgn)

	etcd, err := renderTemplate(s.EtcdTemplate, s.Config)
	if err != nil {
		return fmt.Errorf("failed to render etcd template %v", err)
	}
	etcdIgn, err := convertToIgnition(etcd)
	if err != nil {
		return err
	}
	s.UserDataEtcd = base64.StdEncoding.EncodeToString(etcdIgn)

	body, err := renderTemplate(s.StackTemplate, s)
	if err != nil {
		return err
	}

	// validate the JSON document
	if err = validateJSON(body); err != nil {
		return err
	}

	// minify JSON
	var buf bytes.Buffer
	if err := json.Compact(&buf, body); err != nil {
		return err
	}
	s.StackBody = buf.String()

	return nil
}

func (s *stackVars) validateUserData() error {
	errors := []string{}

	for _, userData := range []struct {
		Name    string
		Content string
	}{
		{
			Content: s.UserDataWorker,
			Name:    "UserDataWorker",
		},
		{
			Content: s.UserDataController,
			Name:    "UserDataController",
		},
	} {
		report, err := validate.Validate([]byte(userData.Content))

		if err != nil {
			errors = append(
				errors,
				fmt.Sprintf("cloud-config %s could not be parsed: %v",
					userData.Name,
					err,
				),
			)
			continue
		}

		for _, entry := range report.Entries() {
			errors = append(errors, fmt.Sprintf("%s: %+v", userData.Name, entry))
		}
	}

	if len(errors) > 0 {
		reportString := strings.Join(errors, "\n")
		return fmt.Errorf("cloud-config validation errors:\n%s\n", reportString)
	}

	return nil
}

// Begin dynamic checks against AWS API
func (c *stackVars) validateAll(session *session.Session) error {
	r53Svc := route53.New(session)
	ec2Svc := ec2.New(session)

	if _, err := c.validateStack(session); err != nil {
		return err
	}

	if c.VPCID == "" {
		if err := ValidateSubnets(c.VPCCIDR, c.ControllerSubnets); err != nil {
			return err
		}

		if err := ValidateSubnets(c.VPCCIDR, c.WorkerSubnets); err != nil {
			return err
		}

		if err := ValidateKubernetesCIDRs(c.VPCCIDR, c.PodCIDR, c.ServiceCIDR); err != nil {
			return err
		}
	} else {
		if err := CheckSubnetsAgainstExistingVPC(session, c.VPCID, c.ControllerSubnets, c.WorkerSubnets); err != nil {
			return err
		}

		if err := CheckKubernetesCIDRs(session, c.VPCID, c.PodCIDR, c.ServiceCIDR); err != nil {
			return err
		}
	}

	if err := c.validateDNSConfig(r53Svc); err != nil {
		return err
	}

	if err := c.validateKeyPair(ec2Svc); err != nil {
		return err
	}

	if err := c.validateControllerRootVolume(ec2Svc); err != nil {
		return err
	}

	if err := c.validateWorkerRootVolume(ec2Svc); err != nil {
		return err
	}

	return nil
}

func (c *stackVars) validateStack(session *session.Session) (string, error) {
	if c.StackURL == "" {
		return "", errors.New("must upload stack before validating it")
	}

	validateInput := cloudformation.ValidateTemplateInput{
		TemplateURL: aws.String(c.StackURL),
	}

	cfSvc := cloudformation.New(session)
	validationReport, err := cfSvc.ValidateTemplate(&validateInput)
	if err != nil {
		return "", maybeAwsErr(err)
	}

	return validationReport.String(), nil
}

func (c *stackVars) validateKeyPair(ec2Svc ec2Service) error {
	_, err := ec2Svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		KeyNames: []*string{aws.String(c.KeyName)},
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "InvalidKeyPair.NotFound" {
				return fmt.Errorf("Key %s does not exist.", c.KeyName)
			}
		}
		return maybeAwsErr(err)
	}
	return nil
}

type r53Service interface {
	ListHostedZonesByName(*route53.ListHostedZonesByNameInput) (*route53.ListHostedZonesByNameOutput, error)
	ListResourceRecordSets(*route53.ListResourceRecordSetsInput) (*route53.ListResourceRecordSetsOutput, error)
	GetHostedZone(*route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error)
}

func (c *stackVars) validateDNSConfig(r53 r53Service) error {
	hzOut, err := r53.GetHostedZone(&route53.GetHostedZoneInput{Id: aws.String(c.HostedZoneID)})
	if err != nil {
		return fmt.Errorf("error getting hosted zone %s: %v", c.HostedZoneID, err)
	}

	if !isSubdomain(c.ControllerDomain, aws.StringValue(hzOut.HostedZone.Name)) {
		return fmt.Errorf("controllerDomain %s is not a sub-domain of hosted-zone %s", c.ControllerDomain, aws.StringValue(hzOut.HostedZone.Name))
	}

	recordSetsResp, err := r53.ListResourceRecordSets(
		&route53.ListResourceRecordSetsInput{
			HostedZoneId: hzOut.HostedZone.Id,
		},
	)
	if err != nil {
		return fmt.Errorf("error listing record sets for hosted zone id = %s: %v", c.HostedZoneID, err)
	}

	if len(recordSetsResp.ResourceRecordSets) > 0 {
		for _, recordSet := range recordSetsResp.ResourceRecordSets {
			if *recordSet.Name == withTrailingDot(c.ControllerDomain) {
				return fmt.Errorf(
					"RecordSet for \"%s\" already exists in Hosted Zone \"%s\"",
					c.ControllerDomain,
					aws.StringValue(hzOut.HostedZone.Name),
				)
			}
		}
	}

	return nil
}

func (c *stackVars) validateControllerRootVolume(ec2Svc ec2Service) error {

	// Send a dry-run request to validate the controller root volume parameters
	controllerRootVolume := &ec2.CreateVolumeInput{
		DryRun:           aws.Bool(true),
		AvailabilityZone: aws.String(c.availabilityZones()[0]),
		Iops:             aws.Int64(int64(c.ControllerRootVolumeIOPS)),
		Size:             aws.Int64(int64(c.ControllerRootVolumeSize)),
		VolumeType:       aws.String(c.ControllerRootVolumeType),
	}

	if _, err := ec2Svc.CreateVolume(controllerRootVolume); err != nil {
		if operr, ok := err.(awserr.Error); ok && operr.Code() != "DryRunOperation" {
			return fmt.Errorf("create volume dry-run request failed: %v", err)
		}
	}

	return nil
}

func (c *stackVars) validateWorkerRootVolume(ec2Svc ec2Service) error {

	// Send a dry-run request to validate the worker root volume parameters
	workerRootVolume := &ec2.CreateVolumeInput{
		DryRun:           aws.Bool(true),
		AvailabilityZone: aws.String(c.availabilityZones()[0]),
		Iops:             aws.Int64(int64(c.WorkerRootVolumeIOPS)),
		Size:             aws.Int64(int64(c.WorkerRootVolumeSize)),
		VolumeType:       aws.String(c.WorkerRootVolumeType),
	}

	if _, err := ec2Svc.CreateVolume(workerRootVolume); err != nil {
		operr, ok := err.(awserr.Error)

		if !ok || (ok && operr.Code() != "DryRunOperation") {
			return fmt.Errorf("create volume dry-run request failed: %v", err)
		}
	}

	return nil
}

// stackLocations returns an unique S3 bucket and object name for the cloud
// formation stack.
func (c *stackVars) stackLocation(session *session.Session) (string, string) {
	return uniqueS3Bucket(session, c.Config.HostedZoneName), fmt.Sprintf("%s/cloudforms.template", c.ClusterName)
}

func (c *stackVars) upload(session *session.Session) error {
	if c.StackBody == "" {
		return errors.New("must render before uploading")
	}

	bucket, name := c.stackLocation(session)
	stackURL, err := uploadS3(session, bucket, name, []byte(c.StackBody))
	if err != nil {
		return err
	}
	c.StackURL = stackURL

	return nil
}

func (c *stackVars) remove(session *session.Session) error {
	bucket, name := c.stackLocation(session)
	return deleteS3(session, bucket, name)
}

// convertToIgnition parses a container linux config and converts it into
// a machine readable Ignition config JSON bytes.
func convertToIgnition(config []byte) ([]byte, error) {
	// parse bytes into a fuze / config-transpiler Config
	tc, report := fuze.Parse(config)
	if report.IsFatal() {
		return nil, fmt.Errorf("error parsing config: %s", report.String())
	}
	// convert into an Ignition config
	ign, report := fuze.ConvertAs2_0(tc, templating.PlatformEC2)
	if report.IsFatal() {
		return nil, fmt.Errorf("error converting into Ignition: %s", report.String())
	}
	return json.Marshal(ign)
}

// renderTemplate executes the given template.Template.
func renderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

func isSubdomain(sub, parent string) bool {
	sub, parent = withTrailingDot(sub), withTrailingDot(parent)
	subParts, parentParts := strings.Split(sub, "."), strings.Split(parent, ".")

	if len(parentParts) > len(subParts) {
		return false
	}

	subSuffixes := subParts[len(subParts)-len(parentParts):]

	if len(subSuffixes) != len(parentParts) {
		return false
	}
	for i := range subSuffixes {
		if subSuffixes[i] != parentParts[i] {
			return false
		}
	}
	return true
}
