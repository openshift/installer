package containerv2

/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Schematics
 * (C) Copyright IBM Corp. 2017 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

/*******************************************************************************
 * A file for openshift related utility functions, like getting kube
 * config
 ******************************************************************************/

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	yaml "github.com/ghodss/yaml"

	"github.com/IBM-Cloud/bluemix-go/client"
	bxhttp "github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	"github.com/IBM-Cloud/bluemix-go/trace"
)

const (
	// IAMHTTPtimeout -
	IAMHTTPtimeout            = 10 * time.Second
	VirtualPrivateEndpoint    = "vpe"
	PrivateServiceEndpoint    = "private"
	VirtualPrivateEndpointDNS = ".vpe.private"
	PrivateEndpointDNS        = ".private"
)

// Frame -
type Frame uintptr

// StackTrace -
type StackTrace []Frame
type stackTracer interface {
	StackTrace() StackTrace
}

type openShiftUser struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
	} `json:"metadata"`
	Identities []string `json:"identities"`
	Groups     []string `json:"groups"`
}

type authEndpoints struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	ServerURL             string `json:"server_endpoint,omitempty"`
}

// PanicCatch - Catch panic and give error
func PanicCatch(r interface{}) error {
	if r != nil {
		var e error
		switch x := r.(type) {
		case string:
			e = errors.New(x)
		case error:
			e = x
		default:
			e = errors.New("Unknown panic")
		}
		fmt.Printf("Panic error %v", e)
		if err, ok := e.(stackTracer); ok {
			fmt.Printf("Panic stack trace %v", err.StackTrace())
		} else {
			debug.PrintStack()
		}
		return e
	}
	return nil
}

// NormalizeName -
func NormalizeName(name string) (string, error) {
	name = strings.ToLower(name)
	reg, err := regexp.Compile("[^A-Za-z0-9:]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(name, "-"), nil
}

// logInAndFillOCToken will update kubeConfig with an Openshift token, if one is not there
func (r *clusters) FetchOCTokenForKubeConfig(kubecfg []byte, cMeta *ClusterInfo, skipSSLVerification bool, endpointType string) (kubecfgEdited []byte, host string, rerr error) {
	// TODO: this is not a a standard manner to login ... using propriatary OC cli reverse engineering
	defer func() {
		err := PanicCatch(recover())
		if err != nil {
			rerr = fmt.Errorf("could not login to openshift account %s", err)
		}
	}()

	var cfg map[string]interface{}
	err := yaml.Unmarshal(kubecfg, &cfg)
	if err != nil {
		return kubecfg, "", err
	}
	var token, passcode string
	if r.client.Config.BluemixAPIKey == "" {
		trace.Logger.Println("Creating user passcode to login for getting oc token")

		// Retry to cover rate limiting on passcode endpoint in particular
		for try := 1; try <= 3; try++ {
			passcode, err = r.client.TokenRefresher.GetPasscode()

			if err == nil {
				break
			}

			if err != nil && try == 3 {
				return kubecfg, "", err
			}

			time.Sleep(1 * time.Second)
		}
	}

	// honor the endpointType parameter if the current parameter is different
	switch endpointType {
	case PrivateServiceEndpoint:
		if !strings.Contains(cMeta.ServerURL, PrivateEndpointDNS) || strings.Contains(cMeta.ServerURL, VirtualPrivateEndpointDNS) {
			// Could be changed to private only if the cluster's private service endpoint is enabled and public is enabled
			if cMeta.ServiceEndpoints.PrivateServiceEndpointEnabled && cMeta.ServiceEndpoints.PrivateServiceEndpointURL != "" && !cMeta.ServiceEndpoints.PublicServiceEndpointEnabled {
				// As this is Openshift, we need to use the URL with the signed certificate (-e) (the right URL is not available in getCluster response)
				urlParts := strings.Split(cMeta.ServiceEndpoints.PrivateServiceEndpointURL, ".")
				cMeta.ServerURL = urlParts[0] + "-e." + strings.Join(urlParts[1:], ".")
			} else {
				trace.Logger.Println("Ignore endpoint parameter and use default ServerURL - currently unsupported scenario")
			}
		}
	case VirtualPrivateEndpoint:
		if !strings.Contains(cMeta.ServerURL, VirtualPrivateEndpointDNS) && !cMeta.ServiceEndpoints.PublicServiceEndpointEnabled {
			if cMeta.VirtualPrivateEndpointURL != "" {
				cMeta.ServerURL = cMeta.VirtualPrivateEndpointURL
			} else {
				return kubecfg, "", fmt.Errorf("virtual private endpoint is not supported by the cluster")
			}
		}
	}

	authEP, err := func(meta *ClusterInfo) (*authEndpoints, error) {
		request := rest.GetRequest(meta.ServerURL + "/.well-known/oauth-authorization-server")
		var auth authEndpoints

		// Create new REST client - reusing modified existing client instances could lead to race conditions
		restClient := &rest.Client{}
		resp, err := restClient.Do(request, &auth, nil)

		if err != nil {
			return &auth, err
		}
		defer resp.Body.Close()
		if resp.StatusCode > 299 {
			msg, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("bad status code [%d] returned when fetching Cluster authentication endpoints: %s", resp.StatusCode, msg)
		}
		if endpointType != "" {
			auth.AuthorizationEndpoint, err = reconfigureAuthorizationEndpoint(auth.AuthorizationEndpoint, endpointType, meta)
			if err != nil {
				return &auth, err
			}
		}
		auth.ServerURL = meta.ServerURL
		return &auth, nil
	}(cMeta)

	if err != nil {
		return kubecfg, "", err
	}

	trace.Logger.Println("Got authentication endpoints for getting oc token")
	token, uname, err := r.openShiftAuthorizePasscode(authEP, passcode, cMeta.IsStagingSatelliteCluster())

	if err != nil {
		return kubecfg, "", err
	}

	trace.Logger.Println("Got the token and user ", uname)
	clusterName, _ := NormalizeName(authEP.ServerURL[len("https://"):len(authEP.ServerURL)]) //TODO deal with http
	ccontext := "default/" + clusterName + "/" + uname
	uname = uname + "/" + clusterName
	clusters := cfg["clusters"].([]interface{})
	newCluster := map[string]interface{}{"name": clusterName, "cluster": map[string]interface{}{"server": authEP.ServerURL}}
	if skipSSLVerification {
		newCluster["cluster"].(map[string]interface{})["insecure-skip-tls-verify"] = true
	}
	clusters = append(clusters, newCluster)
	cfg["clusters"] = clusters

	contexts := cfg["contexts"].([]interface{})
	newContext := map[string]interface{}{"name": ccontext, "context": map[string]interface{}{"cluster": clusterName, "namespace": "default", "user": uname}}
	contexts = append(contexts, newContext)
	cfg["contexts"] = contexts

	users := cfg["users"].([]interface{})
	newUser := map[string]interface{}{"name": uname, "user": map[string]interface{}{"token": token}}
	users = append(users, newUser)
	cfg["users"] = users

	cfg["current-context"] = ccontext

	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return kubecfg, "", err
	}
	kubecfg = bytes
	return kubecfg, cMeta.ServerURL, nil
}

// Never redirect. Let caller handle. This is an http.Client callback method (CheckRedirect)
func neverRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (r *clusters) openShiftAuthorizePasscode(authEP *authEndpoints, passcode string, skipSSLVerification bool) (string, string, error) {
	var request *rest.Request
	authString := "passcode:" + passcode
	if r.client.Config.BluemixAPIKey != "" {
		apikey := r.client.Config.BluemixAPIKey
		authString = "apikey:" + apikey
	}
	request = rest.GetRequest(authEP.AuthorizationEndpoint+"?response_type=token&client_id=openshift-challenging-client").
		Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))
	// Creating a new client instance (instead of tempering with existing one) to avoid race conditions
	copyConfig := r.client.Config.Copy()
	copyConfig.SSLDisable = skipSSLVerification
	copyConfig.HTTPClient = bxhttp.NewHTTPClient(copyConfig)
	copyConfig.HTTPClient.CheckRedirect = neverRedirect

	client := client.New(copyConfig, r.client.ServiceName, r.client.TokenRefresher)

	var respInterface interface{}
	var resp *http.Response
	var err error
	for try := 1; try <= 3; try++ {
		// bmxerror.NewRequestFailure("ServerErrorResponse", string(raw), resp.StatusCode)
		resp, err = client.SendRequest(request, respInterface)
		if err != nil {
			if resp.StatusCode != 302 {
				return "", "", err
			}
		}
		defer resp.Body.Close()
		if resp.StatusCode > 399 {
			if try >= 3 {
				msg, _ := io.ReadAll(resp.Body)
				return "", "", fmt.Errorf("bad status code [%d] returned when openshift login: %s", resp.StatusCode, string(msg))
			}
			time.Sleep(200 * time.Millisecond)
		} else {
			break
		}
	}

	loc, err := resp.Location()
	if err != nil {
		return "", "", err
	}
	val, err := url.ParseQuery(loc.Fragment)
	if err != nil {
		return "", "", err
	}
	token := val.Get("access_token")
	trace.Logger.Println("Getting username after getting the token")
	name, err := r.getOpenShiftUser(authEP, token)
	if err != nil {
		return "", "", err
	}
	return token, name, nil
}

func (r *clusters) getOpenShiftUser(authEP *authEndpoints, token string) (string, error) {
	request := rest.GetRequest(authEP.ServerURL+"/apis/user.openshift.io/v1/users/~").
		Set("Authorization", "Bearer "+token)

	var user openShiftUser
	resp, err := r.client.SendRequest(request, &user)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status code [%d] returned when fetching OpenShift user Details: %s", resp.StatusCode, string(msg))
	}

	return user.Metadata.Name, nil
}

// honor endpointType for OauthServer if the current parameter is different
func reconfigureAuthorizationEndpoint(originalAuthEndpoint string, endpointType string, clusterInfo *ClusterInfo) (string, error) {
	urlDefault, err := url.ParseRequestURI(originalAuthEndpoint)
	if err != nil || urlDefault.Host == "" {
		return "", fmt.Errorf("could not parse original auth endpoint raw url: %s, error: %v", originalAuthEndpoint, err)
	}
	switch endpointType {
	case PrivateServiceEndpoint:
		if (!strings.Contains(originalAuthEndpoint, PrivateEndpointDNS) || strings.Contains(originalAuthEndpoint, VirtualPrivateEndpointDNS)) &&
			!clusterInfo.ServiceEndpoints.PublicServiceEndpointEnabled &&
			clusterInfo.ServiceEndpoints.PrivateServiceEndpointEnabled {
			urlPrivate, err := url.ParseRequestURI(clusterInfo.ServiceEndpoints.PrivateServiceEndpointURL)
			if err != nil || urlPrivate.Host == "" {
				return "", fmt.Errorf("could not parse private service endpoint raw url, cluster may not support it: %s, error: %v", clusterInfo.ServiceEndpoints.PrivateServiceEndpointURL, err)
			}
			// As this is Openshift, we need to use the URL with the signed certificate (the right URL is not available in getCluster response)
			hostNameParts := strings.Split(urlPrivate.Hostname(), ".")
			hostName := hostNameParts[0] + "-e." + strings.Join(hostNameParts[1:], ".")

			u := url.URL{
				Scheme: urlDefault.Scheme,
				Host:   hostName + ":" + urlDefault.Port(),
				Path:   urlDefault.Path,
			}
			return u.String(), nil
		} else {
			trace.Logger.Println("Ignore endpoint parameter and use default OauthServerURL - currently unsupported scenario")
		}
	case VirtualPrivateEndpoint:
		if !strings.Contains(originalAuthEndpoint, VirtualPrivateEndpointDNS) && !clusterInfo.ServiceEndpoints.PublicServiceEndpointEnabled {
			urlVPE, err := url.ParseRequestURI(clusterInfo.VirtualPrivateEndpointURL)
			if err != nil || urlVPE.Host == "" {
				return "", fmt.Errorf("could not parse virtual private endpoint raw url, cluster may not support it: %s, error: %v", clusterInfo.VirtualPrivateEndpointURL, err)
			}
			u := url.URL{
				Scheme: urlDefault.Scheme,
				Host:   urlVPE.Hostname() + ":" + urlDefault.Port(),
				Path:   urlDefault.Path,
			}
			return u.String(), nil
		}
	}
	return originalAuthEndpoint, nil
}
