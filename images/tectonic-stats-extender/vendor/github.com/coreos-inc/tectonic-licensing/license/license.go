package license

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"
)

const (
	LicenseSchemaVersion2       string = "v2"
	LicenseSchemaVersionCurrent string = LicenseSchemaVersion2
)

var (
	ErrLicenseExpired = errors.New("invalid license: license expired")
	ErrInvalidLicense = errors.New("invalid license: license signature invalid")
)

type License struct {
	SchemaVersion  string    `json:"schemaVersion"`
	Version        string    `json:"version"`
	AccountID      string    `json:"accountID"`
	AccountSecret  string    `json:"accountSecret"`
	CreationDate   time.Time `json:"creationDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	// Subscriptions is a map of SubscriptionID -> Subscription
	Subscriptions map[string]Subscription `json:"subscriptions"`
}

func (l *License) Expired(now time.Time) bool {
	return now.After(l.ExpirationDate)
}

type Subscription struct {
	// PlanName is the internal, private plan name for a subscription. Used
	// internally, as a unique, possibly changing identifier
	PlanName string `json:"planName"`
	// PublicPlanName is the user facing plan name. Ex: "Basic - Monthly"
	PublicPlanName string `json:"publicPlanName"`
	PlanID         string `json:"planID"`
	// ProductName is the internal, private plan name for a subscription. Used
	// internally, as a unique, possibly changing identifier
	ProductName string `json:"productName"`
	// PublicPlanName is the user facing product name. Ex: "Quay Enterprise"
	PublicProductName string `json:"publicProductName"`
	ProductID         string `json:"productID"`
	// ServiceStart is when the subscription's service started.
	ServiceStart time.Time `json:"serviceStart"`
	// ServiceEnd is when the subscription's service is set to end. If a
	// subscription is recurring, this value may change over time as the
	// license is updated/replaced.
	ServiceEnd time.Time `json:"serviceEnd"`
	// TrialEnd is when the subscription is scheduled to end
	TrialEnd *time.Time `json:"trialEnd,omitempty"`
	// InTrial indicates if the subscription was actively in a trial when
	// the license was obtained.
	InTrial bool `json:"inTrial"`
	// TrialOnly indicates if this subscription is an evaluation only
	// subscription, meaning it will end immediately when the subscription's
	// TrialEnd is passed.
	TrialOnly bool `json:"trialOnly"`
	// Duration is the length-measures which constitute the rate plan's billing
	// period/interval. Combined with the duration period to get the
	// subscription's billing frequency.
	Duration int32 `json:"duration"`
	// DurationPeriod possible values are "minutes", "days", "months", "years".
	// It measures the magnitude of the rate plan's period.
	DurationPeriod string `json:"durationPeriod"`
	// Entitlements is a mapping of entitlement name (ex:
	// "software.quay.builders") to the value of the entitlement "ex: 25".
	// This is meant to be used by software to control features available to
	// the end user.
	Entitlements map[string]int64 `json:"entitlements"`
}

func (sub Subscription) Inactive(now time.Time) bool {
	return now.Before(sub.ServiceStart)
}

func (sub Subscription) Active(now time.Time) bool {
	return !sub.Inactive(now) && !sub.Expired(now)
}

func (sub Subscription) Expired(now time.Time) bool {
	return now.After(sub.ServiceEnd)
}

func New(accountID string, creationDate, expirationDate time.Time, subscriptions map[string]Subscription) *License {
	return &License{
		SchemaVersion:  LicenseSchemaVersionCurrent,
		AccountID:      accountID,
		CreationDate:   creationDate.UTC(),
		ExpirationDate: expirationDate.UTC(),
		Subscriptions:  subscriptions,
	}
}

type licenseClaims struct {
	SchemaVersion  string `json:"schemaVersion"`
	Version        string `json:"version"`
	CreationDate   string `json:"creationDate"`
	ExpirationDate string `json:"expirationDate"`
	License        string `json:"license"`
}

// Decode takes a public key, and a license. If the license signature can
// be validated using the provided public key, the license is returned.
func Decode(publicKey *rsa.PublicKey, licenseContents string) (*License, error) {
	jws, err := jose.ParseSigned(licenseContents)
	if err != nil {
		return nil, fmt.Errorf("error parsing license as JWT: %s", err)
	}
	payload, err := jws.Verify(publicKey)
	if err != nil {
		return nil, ErrInvalidLicense
	}

	var claims licenseClaims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, err
	}

	err = validateLicenseSchemaVersion(claims)
	if err != nil {
		return nil, err
	}

	var l License
	err = json.Unmarshal([]byte(claims.License), &l)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// Encode takes an rsa.PrivateKey, and a License. It will attempt to
// marshal the license into JSON, and then creates a JWT containing the entire
// license and a few other pieces of metadata as claims in the JWT. This JWT
// is then signed using the provided private key, and serialized as a string to
// be returned.
func Encode(privateKey *rsa.PrivateKey, l *License) (string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, nil)
	if err != nil {
		return "", err
	}

	rawLicenseJSON, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	claims := licenseClaims{
		SchemaVersion:  string(l.SchemaVersion),
		Version:        l.Version,
		CreationDate:   l.CreationDate.UTC().Format(time.RFC3339),
		ExpirationDate: l.ExpirationDate.UTC().Format(time.RFC3339),
		License:        string(rawLicenseJSON),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	jws, err := signer.Sign(claimsJSON)
	if err != nil {
		return "", err
	}

	return jws.CompactSerialize()
}

func validateLicenseSchemaVersion(claims licenseClaims) error {
	// For the v1 license, everything was a claim, and it's schema version
	// was in the "versions" claim with a value of "1".
	// For the v2 license, most of everything is within the license claim, but
	// specific values are also claims, such as "schemaVersion" which has a
	// value of "v2" currently.

	// both licenses have a claim called "version", so when checking v1, check
	// that it doesn't have the schemaVersion field, which exists only in the
	// v2 license format license format
	if claims.SchemaVersion == "" && claims.Version == "1" {
		return errors.New("invalid license: license is too old, please upgrade")
	}

	if claims.SchemaVersion != LicenseSchemaVersionCurrent {
		return fmt.Errorf("invalid license: cannot handle license schema version '%s', expected '%s'", claims.SchemaVersion, LicenseSchemaVersionCurrent)
	}

	return nil
}
