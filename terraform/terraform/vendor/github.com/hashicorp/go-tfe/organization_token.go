package tfe

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ OrganizationTokens = (*organizationTokens)(nil)

// OrganizationTokens describes all the organization token related methods
// that the Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/cloud/api/organization-tokens.html
type OrganizationTokens interface {
	// Create a new organization token, replacing any existing token.
	Create(ctx context.Context, organization string) (*OrganizationToken, error)

	// Read an organization token.
	Read(ctx context.Context, organization string) (*OrganizationToken, error)

	// Delete an organization token.
	Delete(ctx context.Context, organization string) error
}

// organizationTokens implements OrganizationTokens.
type organizationTokens struct {
	client *Client
}

// OrganizationToken represents a Terraform Enterprise organization token.
type OrganizationToken struct {
	ID          string    `jsonapi:"primary,authentication-tokens"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`
	LastUsedAt  time.Time `jsonapi:"attr,last-used-at,iso8601"`
	Token       string    `jsonapi:"attr,token"`
}

// Create a new organization token, replacing any existing token.
func (s *organizationTokens) Create(ctx context.Context, organization string) (*OrganizationToken, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.QueryEscape(organization))
	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	ot := &OrganizationToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Read an organization token.
func (s *organizationTokens) Read(ctx context.Context, organization string) (*OrganizationToken, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.QueryEscape(organization))
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	ot := &OrganizationToken{}
	err = req.Do(ctx, ot)
	if err != nil {
		return nil, err
	}

	return ot, err
}

// Delete an organization token.
func (s *organizationTokens) Delete(ctx context.Context, organization string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/authentication-token", url.QueryEscape(organization))
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}
