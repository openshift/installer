package common

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// FencingCredential stores the information about a baremetal host's management controller.
type FencingCredential struct {
	HostName string `json:"hostName,omitempty" validate:"required,uniqueField"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address" validate:"required,uniqueField"`
}

// ValidateUniqueAndRequiredFields validated unique fields are indeed unique and that required fields exist on a generic element.
func ValidateUniqueAndRequiredFields[T any](elements []T, fldPath *field.Path, filter validator.FilterFunc, fieldName string) field.ErrorList {
	errs := field.ErrorList{}

	values := make(map[string]map[interface{}]struct{})

	// Initialize a new validator and register a custom validation rule for the tag `uniqueField`.
	validate := validator.New()
	if err := validate.RegisterValidation("uniqueField", func(fl validator.FieldLevel) bool {
		valueFound := false
		fieldName := fl.Parent().Type().Name() + "." + fl.FieldName()
		fieldValue := fl.Field().Interface()

		if fl.Field().Type().Comparable() {
			if _, present := values[fieldName]; !present {
				values[fieldName] = make(map[interface{}]struct{})
			}

			fieldValues := values[fieldName]
			if _, valueFound = fieldValues[fieldValue]; !valueFound {
				fieldValues[fieldValue] = struct{}{}
			}
		} else {
			panic(fmt.Sprintf("Cannot apply validation rule 'uniqueField' on field %s", fl.FieldName()))
		}

		return !valueFound
	}); err != nil {
		fmt.Printf("unexpected error registering validation: %v\n", err)
	}

	// Apply validations and translate errors.

	fldPath = fldPath.Child(fieldName)

	for idx, element := range elements {
		err := validate.StructFiltered(element, filter)
		if err != nil {
			elementType := reflect.TypeOf(elements).Elem().Elem().Name()
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				for _, fieldErr := range validationErrs {
					childName := fldPath.Index(idx).Child(fieldErr.Namespace()[len(elementType)+1:])
					switch fieldErr.Tag() {
					case "required":
						errs = append(errs, field.Required(childName, "missing "+fieldErr.Field()))
					case "uniqueField":
						errs = append(errs, field.Duplicate(childName, fieldErr.Value()))
					}
				}
			}
		}
	}
	return errs
}

// ValidateTwoFencingCredentials in case fencing credentials exists validates there are exactly 2.
func ValidateTwoFencingCredentials(fencingCredentials []*FencingCredential, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}
	fencingCredentialsLength := len(fencingCredentials)
	if fencingCredentialsLength > 0 && fencingCredentialsLength != 2 {
		errs = append(errs, field.Forbidden(fldPath, fmt.Sprintf("there should be exactly two fencingCredentials to support the two node cluster, instead %d fencingCredentials were found", fencingCredentialsLength)))
	}

	return errs
}
