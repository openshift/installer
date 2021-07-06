package models

type AccessGroupMember struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type AccessGroupMemberV2 struct {
	ID          string `json:"iam_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Href        string `json:"href,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	CreatedByID string `json:"created_by_id,omitempty"`
}
