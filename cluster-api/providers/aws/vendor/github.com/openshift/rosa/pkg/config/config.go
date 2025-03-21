/*
Copyright (c) 2022 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the types and functions used to manage the configuration of the command line
// client.

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/glog"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/authentication/securestore"

	"github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/debug"
	"github.com/openshift/rosa/pkg/properties"
)

var (
	UpsertConfigToKeyring   = securestore.UpsertConfigToKeyring
	RemoveConfigFromKeyring = securestore.RemoveConfigFromKeyring
	GetConfigFromKeyring    = securestore.GetConfigFromKeyring
)

// Config is the type used to store the configuration of the client.
type Config struct {
	AccessToken  string   `json:"access_token,omitempty" doc:"Bearer access token."`
	ClientID     string   `json:"client_id,omitempty" doc:"OpenID client identifier."`
	ClientSecret string   `json:"client_secret,omitempty" doc:"OpenID client secret."`
	Insecure     bool     `json:"insecure,omitempty" doc:"Enables insecure communication with the server."`
	RefreshToken string   `json:"refresh_token,omitempty" doc:"Offline or refresh token."`
	Scopes       []string `json:"scopes,omitempty" doc:"OpenID scope."`
	TokenURL     string   `json:"token_url,omitempty" doc:"OpenID token URL."`
	URL          string   `json:"url,omitempty" doc:"URL of the API gateway."`
	UserAgent    string   `json:"user_agent,omitempty" doc:"OCM client UserAgent. Default value is used if not set."`
	Version      string   `json:"version,omitempty" doc:"OCM client version. Default value is used if not set."`
	FedRAMP      bool     `json:"fedramp,omitempty" doc:"Indicates FedRAMP."`
}

var DisallowedSetConfigProperties = []string{"scopes"}

func ConfigPropertiesNamesAndDocs() ([]string, []string) {
	configType := reflect.ValueOf(Config{}).Type()
	names := make([]string, configType.NumField())
	docs := make([]string, configType.NumField())
	for i := 0; i < configType.NumField(); i++ {
		tag := configType.Field(i).Tag
		propName := strings.Split(tag.Get("json"), ",")[0]
		names[i] = propName
		propDoc := tag.Get("doc")
		docs[i] = propDoc
	}
	return names, docs
}

func ConfigVarDocs() []string {
	names, docs := ConfigPropertiesNamesAndDocs()
	fieldHelps := make([]string, len(names))
	for i := 0; i < len(names); i++ {
		fieldHelps[i] = fmt.Sprintf("\t%-15s%s", names[i], docs[i])
	}
	return fieldHelps
}

func GetAllConfigProperties() []string {
	properties, _ := ConfigPropertiesNamesAndDocs()
	return properties
}

func GetAllowedConfigProperties() []string {
	allProperties, _ := ConfigPropertiesNamesAndDocs()
	var allowedProperties []string
	for _, prop := range allProperties {
		if !slices.Contains(DisallowedSetConfigProperties, prop) {
			allowedProperties = append(allowedProperties, prop)
		}
	}
	return allowedProperties
}

// Loads the configuration from the OS keyring if requested, load from the configuration file if not
func Load() (cfg *Config, err error) {
	if keyring, ok := IsKeyringManaged(); ok {
		return loadFromOS(keyring)
	}

	return loadFromFile()
}

// Loads the configuration from the OS keyring. If the configuration doesn't exist
// it will return an empty configuration object.
func loadFromOS(keyring string) (cfg *Config, err error) {
	cfg = &Config{}

	data, err := GetConfigFromKeyring(keyring)
	if err != nil {
		return nil, fmt.Errorf("can't load config from OS keyring [%s]: %v", keyring, err)
	}
	// No config found, return
	if len(data) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		// Treat the config as empty if it can't be unmarshalled, it is invalid
		return nil, nil
	}
	return cfg, nil
}

// Loads the configuration from the configuration file. If the configuration file doesn't exist
// it will return an empty configuration object.
func loadFromFile() (cfg *Config, err error) {
	file, err := Location()
	if err != nil {
		return
	}
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		cfg = nil
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("Failed to check if config file '%s' exists: %v", file, err)
		return
	}
	// #nosec G304
	data, err := os.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("Failed to read config file '%s': %v", file, err)
		return
	}
	cfg = new(Config)
	err = json.Unmarshal(data, cfg)
	if err != nil {
		err = fmt.Errorf("Failed to parse config file '%s': %v", file, err)
		return
	}
	return
}

// Save saves the given configuration to the configuration file.
func Save(cfg *Config) error {
	file, err := Location()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshal config: %v", err)
	}

	if keyring, ok := IsKeyringManaged(); ok {
		err := UpsertConfigToKeyring(keyring, data)
		if err != nil {
			return fmt.Errorf("can't save config to OS keyring [%s]: %v", keyring, err)
		}
		return nil
	}

	dir := filepath.Dir(file)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("Failed to create directory %s: %v", dir, err)
	}
	err = os.WriteFile(file, data, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write file '%s': %v", file, err)
	}
	return nil
}

// Remove removes the configuration file.
func Remove() error {
	if keyring, ok := IsKeyringManaged(); ok {
		err := RemoveConfigFromKeyring(keyring)
		if err != nil {
			return fmt.Errorf("can't remove configuration from keyring [%s]: %w", keyring, err)
		}
		return nil
	}

	file, err := Location()
	if err != nil {
		return err
	}
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		return nil
	}
	err = os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

// Location returns the location of the configuration file. If a configuration file
// already exists in the HOME directory, it uses that, otherwise it prefers to
// use the XDG config directory.
func Location() (path string, err error) {
	// Use env variable
	if ocmconfig := os.Getenv(constants.OcmConfig); ocmconfig != "" {
		return ocmconfig, nil
	}

	// Determine home directory to use for the legacy file path
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path = filepath.Join(home, ".ocm.json")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// Determine standard config directory
		configDir, err := os.UserConfigDir()
		if err != nil {
			return path, err
		}

		// Use standard config directory
		path = filepath.Join(configDir, "/ocm/ocm.json")
	}

	return path, nil
}

func (c *Config) GetData(key string) (value string, err error) {
	if c.AccessToken == "" {
		return
	}

	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(c.AccessToken, jwt.MapClaims{})
	if err != nil {
		err = fmt.Errorf("Failed to parse token: %v", err)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = fmt.Errorf("Expected map claims but got %T", claims)
		return
	}
	claim, ok := claims[key]
	if !ok {
		err = fmt.Errorf("Token does not contain the '%s' claim", key)
		return
	}
	value, ok = claim.(string)
	if !ok {
		err = fmt.Errorf("Expected string '%s' but got %T", key, claim)
		return
	}

	return
}

// Armed checks if the configuration contains either credentials or tokens that haven't expired, so
// that it can be used to perform authenticated requests.
func (c *Config) Armed() (armed bool, err error) {
	if c.ClientID != "" && c.ClientSecret != "" {
		armed = true
		return
	}
	now := time.Now()
	if c.AccessToken != "" {
		var expires bool
		var left time.Duration
		var accessToken *jwt.Token
		accessToken, err = ParseToken(c.AccessToken)
		if err != nil {
			err = fmt.Errorf("Failed to parse token: %v", err)
			return
		}
		expires, left, err = getTokenExpiry(accessToken, now)
		if err != nil {
			return
		}
		if !expires || left > 5*time.Second {
			armed = true
			return
		}
	}
	if c.RefreshToken != "" {
		if IsEncryptedToken(c.RefreshToken) {
			// We have no way of knowing an encrypted token expiration, so
			// we assume it's valid and let the access token request fail.
			armed = true
			return
		}
		var expires bool
		var left time.Duration
		var refreshToken *jwt.Token
		refreshToken, err = ParseToken(c.RefreshToken)
		if err != nil {
			err = fmt.Errorf("Failed to parse token: %v", err)
			return
		}
		expires, left, err = getTokenExpiry(refreshToken, now)
		if err != nil {
			return
		}
		if !expires || left > 10*time.Second {
			armed = true
			return
		}
	}
	return
}

// Connection creates a connection using this configuration.
func (c *Config) Connection() (connection *sdk.Connection, err error) {
	// Create the logger:
	level := glog.Level(1)
	if debug.Enabled() {
		level = glog.Level(0)
	}
	logger, err := sdk.NewGlogLoggerBuilder().
		DebugV(level).
		InfoV(level).
		WarnV(level).
		Build()
	if err != nil {
		return
	}

	// Prepare the builder for the connection adding only the properties that have explicit
	// values in the configuration, so that default values won't be overridden:
	builder := sdk.NewConnectionBuilder()
	builder.Logger(logger)
	if c.TokenURL != "" {
		builder.TokenURL(c.TokenURL)
	}
	if c.ClientID != "" || c.ClientSecret != "" {
		builder.Client(c.ClientID, c.ClientSecret)
	}
	if c.Scopes != nil {
		builder.Scopes(c.Scopes...)
	}
	if c.URL != "" {
		builder.URL(c.URL)
	}
	tokens := make([]string, 0, 2)
	if c.AccessToken != "" {
		tokens = append(tokens, c.AccessToken)
	}
	if c.RefreshToken != "" {
		tokens = append(tokens, c.RefreshToken)
	}
	if len(tokens) > 0 {
		builder.Tokens(tokens...)
	}
	builder.Insecure(c.Insecure)

	// Create the connection:
	connection, err = builder.Build()
	if err != nil {
		return
	}

	return
}

func PersistTokens(cfg *Config, accessToken string, refreshToken string) error {
	var err error
	activeCfg := cfg

	if activeCfg == nil {
		// Load the configuration if none is provided
		activeCfg, err = Load()
		if err != nil {
			return err
		}
	}
	activeCfg.AccessToken = accessToken
	activeCfg.RefreshToken = refreshToken
	return Save(activeCfg)
}

// IsKeyringManaged returns the keyring name and a boolean indicating if the config is managed by the keyring.
func IsKeyringManaged() (keyring string, ok bool) {
	keyring = os.Getenv(properties.KeyringEnvKey)
	return keyring, keyring != ""
}

// Returns the available keyrings on the current host
func GetKeyrings() []string {
	backends := securestore.AvailableBackends()
	if len(backends) == 0 {
		return []string{"no available backends"}
	}
	return backends
}
