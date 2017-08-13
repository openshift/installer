package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"

	paws "github.com/coreos/tectonic-installer/installer/pkg/aws"
)

type labelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// awsDescribeRegionsHandler returns the list of AWS regions.
func awsDescribeRegionsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}
	ec2svc := ec2.New(awsSession)

	resp, err := ec2svc.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return fromAWSErr(err)
	}

	regions := make([]string, len(resp.Regions))
	for i, region := range resp.Regions {
		regions[i] = aws.StringValue(region.RegionName)
	}

	return writeJSONResponse(w, req, http.StatusOK, regions)
}

// awsDefaultSubnetsHandler turns the list of default public/private subnets.
func awsDefaultSubnetsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}

	input := struct {
		VpcCIDR string `json:"vpcCIDR"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}

	publicSubnets, privateSubnets, err := paws.GetDefaultSubnets(awsSession, input.VpcCIDR)
	if err != nil {
		return newBadRequestError("Could not get default subnets: %s", err)
	}

	response := struct {
		Public  []paws.VPCSubnet `json:"public"`
		Private []paws.VPCSubnet `json:"private"`
	}{publicSubnets, privateSubnets}

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// awsValidateSubnets validates that the given VPC and Subnet inputs are
// coherent, either statically in the case of a new VPC, or dynamically against
// an existing VPC.
func awsValidateSubnetsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	input := struct {
		VpcCIDR        string           `json:"vpcCIDR"`
		PodCIDR        string           `json:"podCIDR"`
		ServiceCIDR    string           `json:"serviceCIDR"`
		PublicSubnets  []paws.VPCSubnet `json:"publicSubnets"`
		PrivateSubnets []paws.VPCSubnet `json:"privateSubnets"`
		ExistingVPCID  string           `json:"awsVpcId"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}

	type response struct {
		Message string `json:"message"`
		Valid   bool   `json:"valid"`
	}

	if input.ExistingVPCID == "" {
		// Statically validating a new VPC.
		// # Validate private subnets.
		if err := paws.ValidateSubnets(input.VpcCIDR, input.PublicSubnets); err != nil {
			return writeJSONResponse(w, req, http.StatusOK, response{err.Error(), false})
		}
		// # Validate public subnets.
		if err := paws.ValidateSubnets(input.VpcCIDR, input.PrivateSubnets); err != nil {
			return writeJSONResponse(w, req, http.StatusOK, response{err.Error(), false})
		}
		// # Validate Kubernetes CIDRs.
		if err := paws.ValidateKubernetesCIDRs(input.VpcCIDR, input.PodCIDR, input.ServiceCIDR); err != nil {
			return writeJSONResponse(w, req, http.StatusOK, response{err.Error(), false})
		}
		return writeJSONResponse(w, req, http.StatusOK, response{"", true})
	}

	// Dynamically check against an existing VPC.
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}
	// # Check all subnets against the existing VPC.
	err = paws.CheckSubnetsAgainstExistingVPC(awsSession, input.ExistingVPCID, input.PublicSubnets, input.PrivateSubnets)
	if err != nil {
		return writeJSONResponse(w, req, http.StatusOK, response{err.Error(), false})
	}
	// # Check the Kubernetes CIDRs.
	err = paws.CheckKubernetesCIDRs(awsSession, input.ExistingVPCID, input.PodCIDR, input.ServiceCIDR)
	if err != nil {
		return writeJSONResponse(w, req, http.StatusOK, response{err.Error(), false})
	}
	return writeJSONResponse(w, req, http.StatusOK, response{"", true})
}

// awsGetVPCsHandler responds with the list of AWS VPC instances. An AWS
// Session is read from the context.
func awsGetVPCsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}
	ec2svc := ec2.New(awsSession)

	vpcs, err := ec2svc.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		return fromAWSErr(err)
	}

	response := []labelValue{}

	for _, vpc := range vpcs.Vpcs {
		if vpc.VpcId == nil {
			continue
		}

		label := aws.StringValue(vpc.VpcId)
		for _, tag := range vpc.Tags {
			if aws.StringValue(tag.Key) == "Name" {
				label = fmt.Sprintf("%s - %s", aws.StringValue(tag.Value), label)
				break
			}
		}

		response = append(response, labelValue{label, aws.StringValue(vpc.VpcId)})
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].Label >= response[j].Label
	})

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// awsGetVPCsSubnetsHandler returns the list of public/private subnets present
// in the given VPC.
func awsGetVPCsSubnetsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}

	input := struct {
		VpcID string `json:"vpcID"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}
	if input.VpcID == "" {
		return newBadRequestError("missing required VPC ID")
	}

	publicSubnets, privateSubnets, err := paws.GetVPCSubnets(awsSession, input.VpcID)
	if err != nil {
		return fromAWSErr(err)
	}

	response := struct {
		Public  []paws.VPCSubnet `json:"public"`
		Private []paws.VPCSubnet `json:"private"`
	}{publicSubnets, privateSubnets}

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// awsGetKeyPairsHandler returns the list of EC2 key pairs.
func awsGetKeyPairsHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}
	ec2svc := ec2.New(awsSession)

	resp, err := ec2svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return fromAWSErr(err)
	}

	response := []labelValue{}

	for _, keypair := range resp.KeyPairs {
		response = append(response, labelValue{
			aws.StringValue(keypair.KeyName),
			aws.StringValue(keypair.KeyName),
		})
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].Label >= response[j].Label
	})

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// awsGetZonesHandler returns the list of Route53 Hosted Zones.
func awsGetZonesHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {

	addHostedZones := func(allZones *[]*route53.HostedZone, marker *string) error {

		awsSession, err := awsSessionFromRequest(req)
		if err != nil {
			return err
		}
		route53svc := route53.New(awsSession)

		requestInput := route53.ListHostedZonesInput{}
		if *marker != "" {
			requestInput = route53.ListHostedZonesInput{Marker: marker}
		}

		resp, err := route53svc.ListHostedZones(&requestInput)
		if err != nil {
			return err
		}

		*allZones = append(*allZones, resp.HostedZones...)

		if *resp.IsTruncated == true {
			*marker = *resp.NextMarker
		} else {
			*marker = "complete"
		}

		return nil
	}

	var allZones []*route53.HostedZone
	var marker string

	for marker != "complete" {
		err := addHostedZones(&allZones, &marker)
		if err != nil {
			return fromAWSErr(err)
		}
	}

	response := []labelValue{}

	for _, key := range allZones {
		// Strip trailing dot off domain names & add "(private)" to private zones
		label := strings.TrimSuffix(aws.StringValue(key.Name), ".")
		if aws.BoolValue(key.Config.PrivateZone) {
			label = fmt.Sprintf("%s (private)", label)
		}
		response = append(response, labelValue{label, aws.StringValue(key.Id)})
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].Label >= response[j].Label
	})
	return writeJSONResponse(w, req, http.StatusOK, response)
}

// awsGetDomainInfoHandler returns the SOA record for a given zone.
func awsGetDomainInfoHandler(w http.ResponseWriter, req *http.Request, _ *Context) error {
	req.Header.Set("Tectonic-Region", "us-east-1")
	awsSession, err := awsSessionFromRequest(req)
	if err != nil {
		return fromAWSErr(err)
	}

	input := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}

	response := struct {
		SoaTTL     int64    `json:"soaTTL"`
		SoaValue   string   `json:"soaValue"`
		Registered string   `json:"registered"`
		AWSNS      []string `json:"awsNS"`
		PublicNS   []string `json:"publicNS"`
	}{}

	var wg sync.WaitGroup
	done := make(chan struct{})
	errch := make(chan error)
	wg.Add(4)
	go func() {
		wg.Wait()
		close(done)
	}()

	// Read the SOA Value/TTL.
	go func() {
		defer wg.Done()

		records, err := paws.GetHostedZoneRecords(awsSession, input.ID, input.Name, "SOA", 1)
		if err != nil {
			errch <- err
			return
		}

		if len(records) == 1 {
			response.SoaTTL = aws.Int64Value(records[0].TTL)
			response.SoaValue = aws.StringValue(records[0].ResourceRecords[0].Value)
		}
	}()

	// Read the defined nameservers.
	go func() {
		defer wg.Done()

		records, err := paws.GetHostedZoneRecords(awsSession, input.ID, input.Name, "NS", 10)
		if err != nil {
			errch <- err
			return
		}

		if len(records) > 0 {
			domains := make([]string, len(records[0].ResourceRecords))
			for i, record := range records[0].ResourceRecords {
				domains[i] = aws.StringValue(record.Value)
			}
			sort.Strings(domains)
			response.AWSNS = domains
		}
	}()

	// Verify that the domain is registered.
	go func() {
		defer wg.Done()

		split := strings.Split(input.Name, ".")
		if len(split) < 2 {
			response.Registered = route53domains.DomainAvailabilityDontKnow
			return
		}
		domainName := strings.Join(split[len(split)-2:], ".")

		availability, err := route53domains.New(awsSession).CheckDomainAvailability(
			&route53domains.CheckDomainAvailabilityInput{
				DomainName: aws.String(domainName),
			},
		)
		if err != nil {
			errch <- err
			return
		}
		response.Registered = aws.StringValue(availability.Availability)
	}()

	// Get the nameservers we see being set for the domain.
	go func() {
		defer wg.Done()

		nsRecords, err := net.LookupNS(input.Name)
		if err != nil {
			errch <- err
			return
		}

		domains := make([]string, len(nsRecords))
		for i, ns := range nsRecords {
			domains[i] = ns.Host
		}
		sort.Strings(domains)
		response.PublicNS = domains
	}()

	select {
	case <-done:
		return writeJSONResponse(w, req, http.StatusOK, response)
	case <-errch:
		return fromAWSErr(err)
	case <-time.After(10 * time.Second):
		return newError(http.StatusRequestTimeout, "Time-out while querying domain informations")
	}
}

// awsSessionFromRequest creates an AWS Session from credentials in the POST Body.
func awsSessionFromRequest(req *http.Request) (*session.Session, error) {
	credentials := credentials.NewStaticCredentials(
		req.Header.Get("Tectonic-AccessKeyID"),
		req.Header.Get("Tectonic-SecretAccessKey"),
		req.Header.Get("Tectonic-SessionToken"),
	)

	awsConfig := aws.NewConfig().
		WithCredentials(credentials).
		WithRegion(req.Header.Get("Tectonic-Region")).
		WithCredentialsChainVerboseErrors(true)

	return session.NewSession(awsConfig)
}

// fromAWSError returns the HTTP status code as well as the error message from
// an awserr.Error instance.
func fromAWSErr(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			return newError(reqErr.StatusCode(), awsErr.Message())
		}
		return newInternalServerError(awsErr.Message())
	}
	return newInternalServerError("AWS API Error: %v", err)
}
