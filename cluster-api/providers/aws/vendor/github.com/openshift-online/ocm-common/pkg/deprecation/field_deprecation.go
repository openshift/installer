package deprecation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

type fieldDeprecationsContextKey string

const fieldDeprecationsKey fieldDeprecationsContextKey = "fieldDeprecations"

// FieldDeprecations is stored in the request context and contains field deprecation information
type FieldDeprecations struct {
	messages map[string]map[string]string
}

func (f *FieldDeprecations) Add(field, message string, sunsetDate time.Time) error {
	now := time.Now().UTC()
	if now.After(sunsetDate) {
		return errors.New(message)
	}

	f.messages[field] = map[string]string{
		"details":    message,
		"sunsetDate": sunsetDate.Format(time.RFC3339),
	}
	return nil
}

func (f *FieldDeprecations) ToJSON() ([]byte, error) {
	output := make(map[string]string)
	for field, message := range f.messages {
		output[field] = message["details"]
	}
	return json.Marshal(output)
}

func (f *FieldDeprecations) IsEmpty() bool {
	return len(f.messages) == 0
}

func NewFieldDeprecations() FieldDeprecations {
	return FieldDeprecations{messages: make(map[string]map[string]string)}
}

func WithFieldDeprecations(ctx context.Context) context.Context {
	return context.WithValue(ctx, fieldDeprecationsKey, NewFieldDeprecations())
}

func GetFieldDeprecations(ctx context.Context) FieldDeprecations {
	fieldDeprecations, ok := ctx.Value(fieldDeprecationsKey).(FieldDeprecations)
	if !ok {
		return NewFieldDeprecations()
	}
	return fieldDeprecations
}

// FieldDeprecationResponseWriter is a wrapping for the response writer that sets field deprecation headers
type FieldDeprecationResponseWriter struct {
	http.ResponseWriter
	Request *http.Request
}

func (w *FieldDeprecationResponseWriter) WriteHeader(statusCode int) {
	w.setFieldDeprecationHeaders()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *FieldDeprecationResponseWriter) Write(data []byte) (int, error) {
	w.setFieldDeprecationHeaders()
	return w.ResponseWriter.Write(data)
}

func (w *FieldDeprecationResponseWriter) setFieldDeprecationHeaders() {
	deprecatedFields := GetFieldDeprecations(w.Request.Context())
	if !deprecatedFields.IsEmpty() {
		deprecatedFieldsJSON, _ := deprecatedFields.ToJSON()
		w.ResponseWriter.Header().Set(consts.OcmFieldDeprecation, string(deprecatedFieldsJSON))
	}
}
