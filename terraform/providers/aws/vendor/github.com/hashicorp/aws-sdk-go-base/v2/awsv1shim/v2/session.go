package awsv1shim

import ( // nosemgrep: no-sdkv2-imports-in-awsv1shim
	"context"
	"fmt"
	"log"
	"os"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/awsconfig"
	"github.com/hashicorp/aws-sdk-go-base/v2/internal/constants"
)

// getSessionOptions attempts to return valid AWS Go SDK session authentication
// options based on pre-existing credential provider, configured profile, or
// fallback to automatically a determined session via the AWS Go SDK.
func getSessionOptions(awsC *awsv2.Config, c *awsbase.Config) (*session.Options, error) {
	useFIPSEndpoint, _, err := awsconfig.ResolveUseFIPSEndpoint(context.Background(), awsC.ConfigSources)
	if err != nil {
		return nil, fmt.Errorf("error resolving configuration: %w", err)
	}

	useDualStackEndpoint, _, err := awsconfig.ResolveUseDualStackEndpoint(context.Background(), awsC.ConfigSources)
	if err != nil {
		return nil, fmt.Errorf("error resolving configuration: %w", err)
	}

	httpClient, err := defaultHttpClient(c)
	if err != nil {
		return nil, err
	}

	options := &session.Options{
		Config: aws.Config{
			Credentials:          newV2Credentials(awsC.Credentials),
			HTTPClient:           httpClient,
			LogLevel:             aws.LogLevel(aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors),
			Logger:               debugLogger{},
			MaxRetries:           aws.Int(0),
			Region:               aws.String(awsC.Region),
			UseFIPSEndpoint:      convertFIPSEndpointState(useFIPSEndpoint),
			UseDualStackEndpoint: convertDualStackEndpointState(useDualStackEndpoint),
		},
	}

	if c.CustomCABundle != "" {
		reader, err := c.CustomCABundleReader()
		if err != nil {
			return nil, err
		}
		options.CustomCABundle = reader
	}

	return options, nil
}

// GetSession attempts to return valid AWS Go SDK session.
func GetSession(awsC *awsv2.Config, c *awsbase.Config) (*session.Session, error) {
	options, err := getSessionOptions(awsC, c)
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(*options)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, "NoCredentialProviders") {
			return nil, c.NewNoValidCredentialSourcesError(err)
		}
		return nil, fmt.Errorf("Error creating AWS session: %w", err)
	}

	// Set retries after resolving credentials to prevent retries during resolution
	if retryer := awsC.Retryer(); retryer != nil {
		sess = sess.Copy(&aws.Config{MaxRetries: aws.Int(retryer.MaxAttempts())})
	}

	SetSessionUserAgent(sess, c.APNInfo, c.UserAgent)

	// Add custom input from ENV to the User-Agent request header
	// Reference: https://github.com/terraform-providers/terraform-provider-aws/issues/9149
	if v := os.Getenv(constants.AppendUserAgentEnvVar); v != "" {
		log.Printf("[DEBUG] Using additional User-Agent Info: %s", v)
		sess.Handlers.Build.PushBack(request.MakeAddToUserAgentFreeFormHandler(v))
	}

	// Generally, we want to configure a lower retry theshold for networking issues
	// as the session retry threshold is very high by default and can mask permanent
	// networking failures, such as a non-existent service endpoint.
	// MaxRetries will override this logic if it has a lower retry threshold.
	// NOTE: This logic can be fooled by other request errors raising the retry count
	//       before any networking error occurs
	sess.Handlers.Retry.PushBack(func(r *request.Request) {
		if r.RetryCount < constants.MaxNetworkRetryCount {
			return
		}
		// RequestError: send request failed
		// caused by: Post https://FQDN/: dial tcp: lookup FQDN: no such host
		if tfawserr.ErrMessageAndOrigErrContain(r.Error, request.ErrCodeRequestError, "send request failed", "no such host") {
			log.Printf("[WARN] Disabling retries after next request due to networking issue")
			r.Retryable = aws.Bool(false)
		}
		// RequestError: send request failed
		// caused by: Post https://FQDN/: dial tcp IPADDRESS:443: connect: connection refused
		if tfawserr.ErrMessageAndOrigErrContain(r.Error, request.ErrCodeRequestError, "send request failed", "connection refused") {
			log.Printf("[WARN] Disabling retries after next request due to networking issue")
			r.Retryable = aws.Bool(false)
		}
	})

	return sess, nil
}

func convertFIPSEndpointState(value awsv2.FIPSEndpointState) endpoints.FIPSEndpointState {
	switch value {
	case awsv2.FIPSEndpointStateEnabled:
		return endpoints.FIPSEndpointStateEnabled
	case awsv2.FIPSEndpointStateDisabled:
		return endpoints.FIPSEndpointStateDisabled
	default:
		return endpoints.FIPSEndpointStateUnset
	}
}

func convertDualStackEndpointState(value awsv2.DualStackEndpointState) endpoints.DualStackEndpointState {
	switch value {
	case awsv2.DualStackEndpointStateEnabled:
		return endpoints.DualStackEndpointStateEnabled
	case awsv2.DualStackEndpointStateDisabled:
		return endpoints.DualStackEndpointStateDisabled
	default:
		return endpoints.DualStackEndpointStateUnset
	}
}

func SetSessionUserAgent(sess *session.Session, apnInfo *awsbase.APNInfo, userAgentProducts awsbase.UserAgentProducts) {
	// AWS SDK Go automatically adds a User-Agent product to HTTP requests,
	// which contains helpful information about the SDK version and runtime.
	// The configuration of additional User-Agent header products should take
	// precedence over that product. Since the AWS SDK Go request package
	// functions only append, we must PushFront on the build handlers instead
	// of PushBack.
	if apnInfo != nil {
		sess.Handlers.Build.PushFront(
			request.MakeAddToUserAgentFreeFormHandler(apnInfo.BuildUserAgentString()),
		)
	}

	if len(userAgentProducts) > 0 {
		sess.Handlers.Build.PushBack(request.MakeAddToUserAgentFreeFormHandler(userAgentProducts.BuildUserAgentString()))
	}
}
