// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/functions"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/apache/openwhisk-client-go/whisk"
)

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.functions.cloud.ibm.com"

//FunctionClient ...
func FunctionClient(c *bluemix.Config) (*whisk.Client, error) {
	baseEndpoint := getBaseURL(c.Region)
	u, err := url.Parse(fmt.Sprintf("%s/api", baseEndpoint))
	if err != nil {
		return nil, err
	}

	functionsClient, err := whisk.NewClient(http.DefaultClient, &whisk.Config{
		Host:    u.Host,
		Version: "v1",
	})

	return functionsClient, err
}

//getBaseURL ..
func getBaseURL(region string) string {
	baseEndpoint := fmt.Sprintf(DefaultServiceURL)
	if region != "us-south" {
		baseEndpoint = fmt.Sprintf("https://%s.functions.cloud.ibm.com", region)
	}

	return baseEndpoint
}

/*
 *
 * Configure a HTTP client using the OpenWhisk properties (i.e. host, auth, iamtoken)
 * Only cf-based namespaces needs auth key value.
 * iam-based namespace don't have an auth key and needs only iam token for authorization.
 *
 */
func setupOpenWhiskClientConfig(namespace string, sess *bxsession.Session, functionNamespace functions.FunctionServiceAPI) (*whisk.Client, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.functions.cloud.ibm.com/api", sess.Config.Region))
	wskClient, _ := whisk.NewClient(http.DefaultClient, &whisk.Config{
		Host:    u.Host,
		Version: "v1",
	})

	nsList, err := functionNamespace.Namespaces().GetNamespaces()
	if err != nil {
		return nil, err
	}

	var validNamespace bool
	var isCFNamespace bool
	allNamespaces := []string{}
	for _, n := range nsList.Namespaces {
		allNamespaces = append(allNamespaces, n.GetName())
		if n.GetName() == namespace || n.GetID() == namespace {
			if os.Getenv("TF_LOG") != "" {
				whisk.SetDebug(true)
			}
			if n.IsCf() {
				isCFNamespace = true
				break
			}
			validNamespace = true
			// Configure whisk properties to handle iam-based/iam-migrated  namespaces.
			if n.IsIamEnabled() {
				additionalHeaders := make(http.Header)

				err := refreshToken(sess)
				if err != nil {
					for count := sess.Config.MaxRetries; *count >= 0; *count-- {
						if err == nil || !isRetryable(err) {
							break
						}
						err = refreshToken(sess)
					}
					if err != nil {
						return nil, err
					}

				}
				additionalHeaders.Add("Authorization", sess.Config.IAMAccessToken)
				additionalHeaders.Add("X-Namespace-Id", n.GetID())

				wskClient.Config.Namespace = n.GetID()
				wskClient.Config.AdditionalHeaders = additionalHeaders
				return wskClient, nil
			}
		}
	}

	// Configure whisk properties to handle cf-based namespaces.
	if isCFNamespace {
		if sess.Config.UAAAccessToken == "" && sess.Config.UAARefreshToken == "" {
			return nil, fmt.Errorf("Couldn't retrieve auth key for IBM Cloud Function")
		}
		err := validateNamespace(namespace)
		if err != nil {
			return nil, err
		}

		nsList, err := functionNamespace.Namespaces().GetCloudFoundaryNamespaces()
		if err != nil {
			return nil, err
		}

		for _, n := range nsList.Namespaces {
			if n.GetName() == namespace {
				wskClient.Config.Namespace = n.GetName()
				wskClient.Config.AuthToken = fmt.Sprintf("%s:%s", n.GetUUID(), n.GetKey())
				return wskClient, nil
			}
		}
	}

	if !validNamespace {
		return nil, fmt.Errorf("Namespace '%s' is not in the list of entitled namespaces. Available namespaces are %s", namespace, allNamespaces)
	}

	return nil, fmt.Errorf("Failed to create whisk config object for namespace '%s'", namespace)
}
