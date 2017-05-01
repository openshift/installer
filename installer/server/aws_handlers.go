package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/coreos/tectonic-installer/installer/server/aws/cloudforms"
	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

type listItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ListItems are a slice of listItems
type ListItems []listItem

func (slice ListItems) Len() int {
	return len(slice)
}

func (slice ListItems) Less(i, j int) bool {
	return slice[i].Label < slice[j].Label
}

func (slice ListItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type awsCredentials struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Region          string `json:"region"`
}

// Characters that aren't allowed in AWS KMS aliases
var coerceKMSAlias = regexp.MustCompile("[^a-zA-Z0-9/_-]+")

// toAwsAppError returns an AWS-specific error along with AWS's error message
func toAwsAppError(err error) *ctxh.AppError {
	// Generic AWS Error with Code, Message, and original error (if any)
	if awsErr, ok := err.(awserr.Error); ok {
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			return ctxh.NewAppError(err, awsErr.Message(), reqErr.StatusCode())
		}
		return ctxh.NewAppError(err, awsErr.Message(), http.StatusInternalServerError)
	}
	return ctxh.NewAppError(err, fmt.Sprintf("AWS API Error: %v", err), http.StatusInternalServerError)
}

// ec2FromRequest takes an http request and returns an ec2 session
func ec2FromRequest(req *http.Request) (*ec2.EC2, *ctxh.AppError) {
	sess, err := awsSessionFromRequest(req)
	if err != nil {
		return nil, ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
	}
	return ec2.New(sess), nil
}

// awsGetDomainInfoHandler returns the SOA record for a given zone.
func awsGetDomainInfoHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		input := struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{}

		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshall json", http.StatusBadRequest)
		}
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		AccessKeyID := req.Header.Get("Tectonic-AccessKeyID")
		SecretAccessKey := req.Header.Get("Tectonic-SecretAccessKey")
		SessionToken := req.Header.Get("Tectonic-SessionToken")
		// route53domains doesn't work in all regions. use us-east-1
		// ignore error as usEast1Sess is only used by route53DomainsSvc
		usEast1Sess, _ := getAWSSession(AccessKeyID, SecretAccessKey, SessionToken, "us-east-1")

		domainInfo := struct {
			Errors struct {
				SOA        string `json:"soa"`
				Registered string `json:"registered"`
				AWSNS      string `json:"awsNS"`
				PublicNS   string `json:"publicNS"`
			} `json:"errors"`
			SoaTTL     int64    `json:"soaTTL"`
			SoaValue   string   `json:"soaValue"`
			Registered string   `json:"registered"`
			AWSNS      []string `json:"awsNS"`
			PublicNS   []string `json:"publicNS"`
		}{}

		route53Svc := route53.New(sess)

		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			defer wg.Done()
			r53input := route53.ListResourceRecordSetsInput{
				HostedZoneId:    aws.String(input.ID),
				StartRecordName: aws.String(input.Name),
				StartRecordType: aws.String("SOA"),
				MaxItems:        aws.String("1"),
			}
			resp, err := route53Svc.ListResourceRecordSets(&r53input)
			if err != nil {
				domainInfo.Errors.SOA = err.Error()
				return
			}
			if len(resp.ResourceRecordSets) == 1 {
				domainInfo.SoaTTL = aws.Int64Value(resp.ResourceRecordSets[0].TTL)
				domainInfo.SoaValue = aws.StringValue(resp.ResourceRecordSets[0].ResourceRecords[0].Value)
			}
		}()
		go func() {
			defer wg.Done()
			r53input := route53.ListResourceRecordSetsInput{
				HostedZoneId:    aws.String(input.ID),
				StartRecordName: aws.String(input.Name),
				StartRecordType: aws.String("NS"),
				MaxItems:        aws.String("10"),
			}
			resp, err := route53Svc.ListResourceRecordSets(&r53input)
			if err != nil {
				domainInfo.Errors.AWSNS = err.Error()
				return
			}

			if len(resp.ResourceRecordSets) > 0 {
				domains := make([]string, len(resp.ResourceRecordSets[0].ResourceRecords))
				for i, record := range resp.ResourceRecordSets[0].ResourceRecords {
					domains[i] = aws.StringValue(record.Value)
				}
				sort.Strings(domains)
				domainInfo.AWSNS = domains
			}
		}()

		go func() {
			defer wg.Done()
			split := strings.Split(input.Name, ".")
			if usEast1Sess == nil || len(split) < 2 {
				domainInfo.Registered = route53domains.DomainAvailabilityDontKnow
				return
			}
			domainName := strings.Join(split[len(split)-2:], ".")

			route53DomainsSvc := route53domains.New(usEast1Sess)
			availableResp, err := route53DomainsSvc.CheckDomainAvailability(&route53domains.CheckDomainAvailabilityInput{
				DomainName: aws.String(domainName),
			})
			if err != nil {
				domainInfo.Errors.Registered = err.Error()
				return
			}
			domainInfo.Registered = aws.StringValue(availableResp.Availability)
		}()

		go func() {
			defer wg.Done()

			nsRecords, err := net.LookupNS(input.Name)
			if err != nil {
				domainInfo.Errors.PublicNS = err.Error()
				return
			}

			domains := make([]string, len(nsRecords))
			for i, ns := range nsRecords {
				domains[i] = ns.Host
			}
			sort.Strings(domains)
			domainInfo.PublicNS = domains
		}()

		wg.Wait()

		writeJSONData(w, domainInfo)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

func awsGetZonesHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		route53Svc := route53.New(sess)
		resp, err := route53Svc.ListHostedZones(&route53.ListHostedZonesInput{})
		if err != nil {
			return toAwsAppError(err)
		}

		keys := []listItem{}

		for _, key := range resp.HostedZones {
			// Strip trailing dot off domain names & add "(private)" to private zones
			label := strings.TrimSuffix(aws.StringValue(key.Name), ".")
			if aws.BoolValue(key.Config.PrivateZone) {
				label = fmt.Sprintf("%s (private)", label)
			}
			keys = append(keys, listItem{
				Label: label,
				Value: aws.StringValue(key.Id),
			})
		}

		writeJSONData(w, keys)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsGetKeyPairsHandler responds with the list of AWS keypairs. An AWS Session
// is read from the context.
func awsGetKeyPairsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		ec2Svc, appErr := ec2FromRequest(req)
		if appErr != nil {
			return appErr
		}

		resp, err := ec2Svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
		if err != nil {
			return toAwsAppError(err)
		}

		keys := make(ListItems, len(resp.KeyPairs))
		for i, keypair := range resp.KeyPairs {
			keyName := aws.StringValue(keypair.KeyName)
			keys[i] = listItem{
				Label: keyName,
				Value: keyName,
			}
		}
		sort.Sort(keys)
		writeJSONData(w, keys)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsDescribeRegionsHandler responds with the list of AWS regions.
func awsDescribeRegionsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		ec2Svc, appErr := ec2FromRequest(req)
		if appErr != nil {
			return appErr
		}

		resp, err := ec2Svc.DescribeRegions(&ec2.DescribeRegionsInput{})
		if err != nil {
			return toAwsAppError(err)
		}

		regions := make([]string, len(resp.Regions))
		for i, region := range resp.Regions {
			regions[i] = aws.StringValue(region.RegionName)
		}

		writeJSONData(w, regions)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

func formatKeyLabel(alias string, keyID string) string {
	alias = strings.TrimPrefix(alias, "alias/")
	if alias == "" {
		return keyID
	}
	return fmt.Sprintf("%s - %s", alias, keyID)
}

// awsCreateKMSHandler creates a new KMS key.
func awsCreateKMSHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}
		kmsSvc := kms.New(sess)

		// Create a new key.
		keyOutput, err := kmsSvc.CreateKey(&kms.CreateKeyInput{
			Description: aws.String("Tectonic installer"),
		})
		if err != nil {
			return toAwsAppError(err)
		}

		// Generate intelligible name for the key.
		// The name's format is <username>-tectonic/<random-number>.
		// The random-number is used so multiple keys can be created by the same
		// user.
		num, err := rand.Int(rand.Reader, new(big.Int).SetInt64(10000))
		if err != nil {
			return ctxh.NewAppError(err, "could not create bigInt", http.StatusInternalServerError)
		}
		username, err := awsGetUsername(sess)
		if err != nil {
			return ctxh.NewAppError(err, "could not retrieve username", http.StatusInternalServerError)
		}
		username = coerceKMSAlias.ReplaceAllString(username, "")
		alias := fmt.Sprintf("%v-tectonic/%v", username, num)

		// Create an alias for the key using the name we just generated.
		_, err = kmsSvc.CreateAlias(&kms.CreateAliasInput{
			AliasName:   aws.String("alias/" + alias),
			TargetKeyId: keyOutput.KeyMetadata.KeyId,
		})
		if err != nil {
			// Failed to create the alias, print the error and set the alias name to
			// an empty string so it doesn't appear in the key label.
			log.Errorf("failed to create alias: %v", err)
			alias = ""
		}

		writeJSONData(w, listItem{
			Label: formatKeyLabel(alias, *keyOutput.KeyMetadata.KeyId),
			Value: *keyOutput.KeyMetadata.Arn,
		})
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsGetKMSHandler responds with the list of AWS KMS keys.
func awsGetKMSHandler() ctxh.ContextHandler {
	toLabel := func(aliases []*kms.AliasListEntry, keyID string) string {
		for _, alias := range aliases {
			if alias.TargetKeyId == nil {
				continue
			}
			if *alias.TargetKeyId == keyID {
				return formatKeyLabel(*alias.AliasName, keyID)
			}
		}
		return keyID
	}

	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		kmsSvc := kms.New(sess)

		output := struct {
			keys       *kms.ListKeysOutput
			aliases    *kms.ListAliasesOutput
			keysErr    error
			aliasesErr error
		}{}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			output.keys, output.keysErr = kmsSvc.ListKeys(&kms.ListKeysInput{})
		}()
		go func() {
			defer wg.Done()
			output.aliases, output.aliasesErr = kmsSvc.ListAliases(&kms.ListAliasesInput{})
		}()

		// Wait for health checks to get responses or timeout
		wg.Wait()

		if output.aliasesErr != nil {
			return toAwsAppError(output.aliasesErr)
		}
		if output.keysErr != nil {
			return toAwsAppError(output.keysErr)
		}

		keys := make(ListItems, len(output.keys.Keys))
		for i, key := range output.keys.Keys {
			keys[i] = listItem{
				Label: toLabel(output.aliases.Aliases, aws.StringValue(key.KeyId)),
				Value: aws.StringValue(key.KeyArn),
			}
		}
		sort.Sort(keys)
		writeJSONData(w, keys)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsGetVPCsHandler responds with the list of AWS VPC instances. An AWS
// Session is read from the context.
func awsGetVPCsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		ec2Svc := ec2.New(sess)
		resp, err := ec2Svc.DescribeVpcs(&ec2.DescribeVpcsInput{})
		if err != nil {
			return ctxh.NewAppError(err, "could not describe VPCs", http.StatusInternalServerError)
		}

		vpcs := make(ListItems, 0)
		for _, vpc := range resp.Vpcs {
			if vpc.VpcId == nil {
				continue
			}

			name := ""
			for _, tag := range vpc.Tags {
				if *tag.Key == "Name" {
					name = aws.StringValue(tag.Value)
					break
				}
			}

			vpcID := aws.StringValue(vpc.VpcId)
			var label string
			if name == "" {
				label = vpcID
			} else {
				label = fmt.Sprintf("%s - %s", name, vpcID)
			}

			vpcs = append(vpcs, listItem{
				Label: label,
				Value: vpcID,
			})
		}

		sort.Sort(vpcs)
		writeJSONData(w, vpcs)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

func awsGetVPCsSubnetsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		input := struct {
			VpcID string `json:"vpcID"`
		}{}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal VPC ID", http.StatusBadRequest)
		}
		if len(input.VpcID) == 0 {
			return ctxh.NewAppError(err, "need VPC ID", http.StatusBadRequest)
		}

		publicSubnets, privateSubnets, err := cloudforms.GetVPCSubnets(sess, input.VpcID)
		if err != nil {
			return ctxh.NewAppError(err, "could not get net slices", http.StatusBadRequest)
		}
		response := struct {
			Public  []cloudforms.VPCSubnet `json:"public"`
			Private []cloudforms.VPCSubnet `json:"private"`
		}{publicSubnets, privateSubnets}

		writeJSONData(w, response)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsDefaultSubnetsHandler responds with the default public/private subnets.
func awsDefaultSubnetsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		sess, err := awsSessionFromRequest(req)
		if err != nil {
			return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
		}

		vpcCIDR := struct {
			VpcCIDR string `json:"vpcCIDR"`
		}{}

		if err := json.NewDecoder(req.Body).Decode(&vpcCIDR); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal VPC CIDR", http.StatusBadRequest)
		}

		publicSubnets, privateSubnets, err := cloudforms.GetDefaultSubnets(sess, vpcCIDR.VpcCIDR)
		if err != nil {
			return ctxh.NewAppError(err, "could not get net slices", http.StatusBadRequest)
		}
		response := struct {
			Public  []cloudforms.VPCSubnet `json:"public"`
			Private []cloudforms.VPCSubnet `json:"private"`
		}{publicSubnets, privateSubnets}

		writeJSONData(w, response)
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsValidateSubnets checks VPC and Subnet choices.
func awsValidateSubnetsHandler() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		input := struct {
			VpcCIDR        string                 `json:"vpcCIDR"`
			PodCIDR        string                 `json:"podCIDR"`
			ServiceCIDR    string                 `json:"serviceCIDR"`
			PublicSubnets  []cloudforms.VPCSubnet `json:"publicSubnets"`
			PrivateSubnets []cloudforms.VPCSubnet `json:"privateSubnets"`
			ExistingVPCID  string                 `json:"awsVpcId"`
		}{}

		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal subnets", http.StatusBadRequest)
		}

		type Validation struct {
			Message string `json:"message"`
			Valid   bool   `json:"valid"`
		}

		if input.ExistingVPCID == "" {
			// new VPC will be created
			if err := cloudforms.ValidateSubnets(input.VpcCIDR, input.PublicSubnets); err != nil {
				writeJSONData(w, Validation{err.Error(), false})
				return nil
			}
			if err := cloudforms.ValidateSubnets(input.VpcCIDR, input.PrivateSubnets); err != nil {
				writeJSONData(w, Validation{err.Error(), false})
				return nil
			}
			if err := cloudforms.ValidateKubernetesCIDRs(input.VpcCIDR, input.PodCIDR, input.ServiceCIDR); err != nil {
				writeJSONData(w, Validation{err.Error(), false})
				return nil
			}
		} else {
			// existing VPC will be used, check against it
			sess, err := awsSessionFromRequest(req)
			if err != nil {
				return ctxh.NewAppError(err, "could not create AWS session", http.StatusInternalServerError)
			}
			err = cloudforms.CheckSubnetsAgainstExistingVPC(sess, input.ExistingVPCID, input.PublicSubnets, input.PrivateSubnets)
			if err != nil {
				writeJSONData(w, Validation{err.Error(), false})
				return nil
			}
			err = cloudforms.CheckKubernetesCIDRs(sess, input.ExistingVPCID, input.PodCIDR, input.ServiceCIDR)
			if err != nil {
				writeJSONData(w, Validation{err.Error(), false})
				return nil
			}
		}

		writeJSONData(w, Validation{"", true})
		return nil
	}
	return requireHTTPMethod("POST", ctxh.ContextHandlerFuncWithError(fn))
}

// awsSessionFromRequest creates an AWS Session from credentials in the POST Body.
func awsSessionFromRequest(req *http.Request) (*session.Session, error) {
	AccessKeyID := req.Header.Get("Tectonic-AccessKeyID")
	SecretAccessKey := req.Header.Get("Tectonic-SecretAccessKey")
	SessionToken := req.Header.Get("Tectonic-SessionToken")
	Region := req.Header.Get("Tectonic-Region")
	return getAWSSession(AccessKeyID, SecretAccessKey, SessionToken, Region)
}

// awsGetUsername returns the username of the session.
func awsGetUsername(session *session.Session) (string, error) {
	// Regardless of the credentials used by the session being a static
	// AccessKeyID / SecretAccessKey or a temporary set using an STS Token,
	// GetCallerIdentity always works, requires no permission and returns the ARN.
	//
	// The last part of the ARN contains the username, the role session name
	// when STS is used, or the role name when STS is used and there was no role
	// session name provided.
	stsOutput, err := sts.New(session).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		// Typically happens if the credentials are expired.
		return "", err
	}

	arnS := strings.Split(*stsOutput.Arn, "/")
	if len(arnS) == 0 {
		// Should never happen, the ARN should always contain at least one '/'.
		return "", fmt.Errorf("misformated ARN: %q", *stsOutput.Arn)
	}

	return arnS[len(arnS)-1], nil
}
