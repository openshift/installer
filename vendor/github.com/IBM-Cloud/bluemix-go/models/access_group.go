package models

// AccessGroup represents the access group of IAM UUM
type AccessGroup struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type AccessGroupV2 struct {
	AccessGroup
	AccountID        string `json:"account_id,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	CreatedByID      string `json:"created_by_id,omitempty"`
	LastModifiedAt   string `json:"last_modified_at,omitempty"`
	LastModifiedByID string `json:"last_modified_by_id,omitempty"`
}
