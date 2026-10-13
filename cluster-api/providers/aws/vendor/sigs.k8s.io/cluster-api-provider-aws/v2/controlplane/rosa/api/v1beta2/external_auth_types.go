package v1beta2

// ExternalAuthProvider is an external OIDC identity provider that can issue tokens for this cluster
type ExternalAuthProvider struct {
	// Name of the OIDC provider
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
	// Issuer describes attributes of the OIDC token issuer
	//
	// +kubebuilder:validation:Required
	// +required
	Issuer TokenIssuer `json:"issuer"`

	// OIDCClients contains configuration for the platform's clients that
	// need to request tokens from the issuer
	//
	// +listType=map
	// +listMapKey=componentNamespace
	// +listMapKey=componentName
	// +kubebuilder:validation:MaxItems=20
	// +optional
	OIDCClients []OIDCClientConfig `json:"oidcClients,omitempty"`

	// ClaimMappings describes rules on how to transform information from an
	// ID token into a cluster identity
	// +optional
	ClaimMappings *TokenClaimMappings `json:"claimMappings,omitempty"`

	// ClaimValidationRules are rules that are applied to validate token claims to authenticate users.
	//
	// +listType=atomic
	ClaimValidationRules []TokenClaimValidationRule `json:"claimValidationRules,omitempty"`
}

// TokenAudience is the audience that the token was issued for.
//
// +kubebuilder:validation:MinLength=1
type TokenAudience string

// TokenIssuer describes attributes of the OIDC token issuer
type TokenIssuer struct {
	// URL is the serving URL of the token issuer.
	// Must use the https:// scheme.
	//
	// +kubebuilder:validation:Pattern=`^https:\/\/[^\s]`
	// +kubebuilder:validation:Required
	// +required
	URL string `json:"issuerURL"`

	// Audiences is an array of audiences that the token was issued for.
	// Valid tokens must include at least one of these values in their
	// "aud" claim.
	// Must be set to exactly one value.
	//
	// +listType=set
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +required
	Audiences []TokenAudience `json:"audiences"`

	// CertificateAuthority is a reference to a config map in the
	// configuration namespace. The .data of the configMap must contain
	// the "ca-bundle.crt" key.
	// If unset, system trust is used instead.
	CertificateAuthority *LocalObjectReference `json:"issuerCertificateAuthority,omitempty"`
}

// OIDCClientConfig contains configuration for the platform's client that
// need to request tokens from the issuer.
type OIDCClientConfig struct {
	// ComponentName is the name of the component that is supposed to consume this
	// client configuration
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Required
	// +required
	ComponentName string `json:"componentName"`

	// ComponentNamespace is the namespace of the component that is supposed to consume this
	// client configuration
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Required
	// +required
	ComponentNamespace string `json:"componentNamespace"`

	// ClientID is the identifier of the OIDC client from the OIDC provider
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +required
	ClientID string `json:"clientID"`

	// ClientSecret refers to a secret that
	// contains the client secret in the `clientSecret` key of the `.data` field
	ClientSecret LocalObjectReference `json:"clientSecret"`

	// ExtraScopes is an optional set of scopes to request tokens with.
	//
	// +listType=set
	// +optional
	ExtraScopes []string `json:"extraScopes,omitempty"`
}

// TokenClaimMappings describes rules on how to transform information from an
// ID token into a cluster identity.
type TokenClaimMappings struct {
	// Username is a name of the claim that should be used to construct
	// usernames for the cluster identity.
	//
	// Default value: "sub"
	// +optional
	Username *UsernameClaimMapping `json:"username,omitempty"`

	// Groups is a name of the claim that should be used to construct
	// groups for the cluster identity.
	// The referenced claim must use array of strings values.
	// +optional
	Groups *PrefixedClaimMapping `json:"groups,omitempty"`
}

// PrefixedClaimMapping defines claims with a prefix.
type PrefixedClaimMapping struct {
	// Claim is a JWT token claim to be used in the mapping
	//
	// +kubebuilder:validation:Required
	// +required
	Claim string `json:"claim"`

	// Prefix is a string to prefix the value from the token in the result of the
	// claim mapping.
	//
	// By default, no prefixing occurs.
	//
	// Example: if `prefix` is set to "myoidc:"" and the `claim` in JWT contains
	// an array of strings "a", "b" and  "c", the mapping will result in an
	// array of string "myoidc:a", "myoidc:b" and "myoidc:c".
	Prefix string `json:"prefix,omitempty"`
}

// UsernameClaimMapping defines the claim that should be used to construct usernames for the cluster identity.
//
// +kubebuilder:validation:XValidation:rule="self.prefixPolicy == 'Prefix' ? has(self.prefix) : !has(self.prefix)",message="prefix must be set if prefixPolicy is 'Prefix', but must remain unset otherwise"
type UsernameClaimMapping struct {
	// Claim is a JWT token claim to be used in the mapping
	//
	// +kubebuilder:validation:Required
	// +required
	Claim string `json:"claim"`

	// PrefixPolicy specifies how a prefix should apply.
	//
	// By default, claims other than `email` will be prefixed with the issuer URL to
	// prevent naming clashes with other plugins.
	//
	// Set to "NoPrefix" to disable prefixing.
	//
	// Example:
	//     (1) `prefix` is set to "myoidc:" and `claim` is set to "username".
	//         If the JWT claim `username` contains value `userA`, the resulting
	//         mapped value will be "myoidc:userA".
	//     (2) `prefix` is set to "myoidc:" and `claim` is set to "email". If the
	//         JWT `email` claim contains value "userA@myoidc.tld", the resulting
	//         mapped value will be "myoidc:userA@myoidc.tld".
	//     (3) `prefix` is unset, `issuerURL` is set to `https://myoidc.tld`,
	//         the JWT claims include "username":"userA" and "email":"userA@myoidc.tld",
	//         and `claim` is set to:
	//         (a) "username": the mapped value will be "https://myoidc.tld#userA"
	//         (b) "email": the mapped value will be "userA@myoidc.tld"
	//
	// +kubebuilder:validation:Enum={"", "NoPrefix", "Prefix"}
	// +optional
	PrefixPolicy UsernamePrefixPolicy `json:"prefixPolicy,omitempty"`

	// Prefix is prepended to claim to prevent clashes with existing names.
	//
	// +kubebuilder:validation:MinLength=1
	// +optional
	Prefix *string `json:"prefix,omitempty"`
}

// UsernamePrefixPolicy specifies how a prefix should apply.
type UsernamePrefixPolicy string

const (
	// NoOpinion let's the cluster assign prefixes.  If the username claim is email, there is no prefix
	// If the username claim is anything else, it is prefixed by the issuerURL
	NoOpinion UsernamePrefixPolicy = ""

	// NoPrefix means the username claim value will not have any  prefix
	NoPrefix UsernamePrefixPolicy = "NoPrefix"

	// Prefix means the prefix value must be specified.  It cannot be empty
	Prefix UsernamePrefixPolicy = "Prefix"
)

// TokenValidationRuleType defines the type of the validation rule.
type TokenValidationRuleType string

const (
	// TokenValidationRuleTypeRequiredClaim defines the type for RequiredClaim.
	TokenValidationRuleTypeRequiredClaim TokenValidationRuleType = "RequiredClaim"
)

// TokenClaimValidationRule validates token claims to authenticate users.
type TokenClaimValidationRule struct {
	// Type sets the type of the validation rule
	//
	// +kubebuilder:validation:Enum={"RequiredClaim"}
	// +kubebuilder:default="RequiredClaim"
	Type TokenValidationRuleType `json:"type"`

	// RequiredClaim allows configuring a required claim name and its expected value
	// +kubebuilder:validation:Required
	RequiredClaim TokenRequiredClaim `json:"requiredClaim"`
}

// TokenRequiredClaim allows configuring a required claim name and its expected value.
type TokenRequiredClaim struct {
	// Claim is a name of a required claim. Only claims with string values are
	// supported.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +required
	Claim string `json:"claim"`

	// RequiredValue is the required value for the claim.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +required
	RequiredValue string `json:"requiredValue"`
}

// LocalObjectReference references an object in the same namespace.
type LocalObjectReference struct {
	// Name is the metadata.name of the referenced object.
	//
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}
