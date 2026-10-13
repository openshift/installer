package authentication

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

const (
	DeviceAuthURL = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/auth/device"
)

type DeviceAuthConfig struct {
	conf               *oauth2.Config
	verifierOpt        oauth2.AuthCodeOption
	DeviceAuthResponse *oauth2.DeviceAuthResponse
	ClientID           string
}

// Step 1:
// Initiates device code flow and returns the device auth config.
// After running, use your DeviceAuthConfig to display the user code and verification URI
//
//	fmt.Printf("To continue login, navigate to %v and enter code %v\n", deviceAuthResp.VerificationURI, deviceAuthResp.UserCode)
//	fmt.Printf("Checking status every %v seconds...\n", deviceAuthResp.Interval)
func (d *DeviceAuthConfig) InitiateDeviceAuth(ctx context.Context) (*DeviceAuthConfig, error) {
	d.conf = &oauth2.Config{
		ClientID:     d.ClientID,
		ClientSecret: "",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			DeviceAuthURL: DeviceAuthURL,
			TokenURL:      DefaultTokenURL,
		},
	}

	// Verifiers and Challenges are required for device auth
	verifier := oauth2.GenerateVerifier()
	verifierOpt := oauth2.VerifierOption(verifier)
	challenge := oauth2.S256ChallengeOption(verifier)

	// Get device code
	deviceAuthResp, err := d.conf.DeviceAuth(ctx, challenge, verifierOpt)
	if err != nil {
		return d, fmt.Errorf("failed to get device code: %v", err)
	}

	d.DeviceAuthResponse = deviceAuthResp
	d.verifierOpt = verifierOpt

	return d, nil
}

// Step 2:
// Initiates polling for token exchange and returns a refresh token
func (d *DeviceAuthConfig) PollForTokenExchange(ctx context.Context) (string, error) {
	if d.DeviceAuthResponse == nil || d.verifierOpt == nil {
		return "", fmt.Errorf("required config is nil, please run InitiateDeviceAuth first")
	}
	// Wait for the user to enter the code, polls at interval specified in deviceAuthResp.Interval
	token, err := d.conf.DeviceAccessToken(ctx, d.DeviceAuthResponse, d.verifierOpt)
	if err != nil {
		return "", fmt.Errorf("error exchanging for token: %v", err)
	}

	return token.RefreshToken, nil
}
