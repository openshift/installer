package cfschema

const (
	HandlerTypeCreate = "create"
	HandlerTypeDelete = "delete"
	HandlerTypeList   = "list"
	HandlerTypeRead   = "read"
	HandlerTypeUpdate = "update"
)

type Handler struct {
	Permissions      []string `json:"permissions,omitempty"`
	TimeoutInMinutes int      `json:"timeoutInMinutes,omitempty"`
}
