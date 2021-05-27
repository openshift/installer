// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type QualifiedName struct {
	namespace   string // namespace. does not include leading '/'.  may be "" (i.e. default namespace)
	packageName string // package.  may be "".  does not include leading/trailing '/'
	entity      string // entity.  should not be ""
	EntityName  string // pkg+entity
}

// Imported code from openwhisk cli https://github.com/apache/incubator-openwhisk/tree/26146368f1dd07f817062e662db64c73a8d486d6/tools/cli/go-whisk-cli/commands
///////////////////////////
// QualifiedName Methods //
///////////////////////////

//  GetFullQualifiedName() returns a full qualified name in proper string format
//      from qualifiedName with proper syntax.
//  Example: /namespace/[package/]entity
func (qualifiedName *QualifiedName) GetFullQualifiedName() string {
	output := []string{}

	if len(qualifiedName.GetNamespace()) > 0 {
		output = append(output, "/", qualifiedName.GetNamespace(), "/")
	}
	if len(qualifiedName.GetPackageName()) > 0 {
		output = append(output, qualifiedName.GetPackageName(), "/")
	}
	output = append(output, qualifiedName.GetEntity())

	return strings.Join(output, "")
}

//  GetPackageName() returns the package name from qualifiedName without a
//      leading '/'
func (qualifiedName *QualifiedName) GetPackageName() string {
	return qualifiedName.packageName
}

//  GetEntityName() returns the entity name ([package/]entity) of qualifiedName
//      without a leading '/'
func (qualifiedName *QualifiedName) GetEntityName() string {
	return qualifiedName.EntityName
}

//  GetEntity() returns the name of entity in qualifiedName without a leading '/'
func (qualifiedName *QualifiedName) GetEntity() string {
	return qualifiedName.entity
}

//  GetNamespace() returns the name of the namespace in qualifiedName without
//      a leading '/'
func (qualifiedName *QualifiedName) GetNamespace() string {
	return qualifiedName.namespace
}

//  NewQualifiedName(name) initializes and constructs a (possibly fully qualified)
//      QualifiedName struct.
//
//      NOTE: If the given qualified name is None, then this is a default qualified
//          name and it is resolved from properties.
//      NOTE: If the namespace is missing from the qualified name, the namespace
//          is also resolved from the property file.
//
//  Examples:
//      foo => qualifiedName {namespace: "_", entityName: foo}
//      pkg/foo => qualifiedName {namespace: "_", entityName: pkg/foo}
//      /ns/foo => qualifiedName {namespace: ns, entityName: foo}
//      /ns/pkg/foo => qualifiedName {namespace: ns, entityName: pkg/foo}
func NewQualifiedName(name string) (*QualifiedName, error) {
	qualifiedName := new(QualifiedName)

	// If name has a preceding delimiter (/), or if it has two delimiters with a
	// leading non-empty string, then it contains a namespace. Otherwise the name
	// does not specify a namespace, so default the namespace to the namespace
	// value set in the properties file; if that is not set, use "_"
	name = addLeadSlash(name)
	parts := strings.Split(name, "/")
	if strings.HasPrefix(name, "/") {
		qualifiedName.namespace = parts[1]

		if len(parts) < 2 || len(parts) > 4 {
			return qualifiedName, qualifiedNameNotSpecifiedErr()
		}

		for i := 1; i < len(parts); i++ {
			if len(parts[i]) == 0 || parts[i] == "." {
				return qualifiedName, qualifiedNameNotSpecifiedErr()
			}
		}

		qualifiedName.EntityName = strings.Join(parts[2:], "/")
		if len(parts) == 4 {
			qualifiedName.packageName = parts[2]
		}
		qualifiedName.entity = parts[len(parts)-1]
	} else {
		if len(name) == 0 || name == "." {
			return qualifiedName, qualifiedNameNotSpecifiedErr()
		}

		qualifiedName.entity = parts[len(parts)-1]
		if len(parts) == 2 {
			qualifiedName.packageName = parts[0]
		}
		qualifiedName.EntityName = name
		qualifiedName.namespace = getNamespaceFromProp()
	}

	return qualifiedName, nil
}

/////////////////////
// Error Functions //
/////////////////////

//  qualifiedNameNotSpecifiedErr() returns generic whisk error for
//      invalid qualified names detected while building a new
//      QualifiedName struct.
func qualifiedNameNotSpecifiedErr() error {
	return errors.New("A valid qualified name must be specified.")
}

//  NewQualifiedNameError(entityName, err) returns specific whisk error
//      for invalid qualified names.
func NewQualifiedNameError(entityName string, err error) error {
	errorMsg := fmt.Sprintf("%s is not a alid qualified name %s", entityName, err)
	return errors.New(errorMsg)
}

///////////////////////////
// Helper/Misc Functions //
///////////////////////////

//  addLeadSlash(name) returns a (possibly fully qualified) resource name,
//      inserting a leading '/' if it is of 3 parts (namespace/package/action)
//      and lacking the leading '/'.
func addLeadSlash(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) == 3 && parts[0] != "" {
		name = "/" + name
	}
	return name
}

//  getNamespaceFromProp() returns a namespace from Properties if one exists,
//      else defaults to returning "_"
func getNamespaceFromProp() string {
	namespace := os.Getenv("FUNCTION_NAMESPACE")
	return namespace
}

//  getQualifiedName(name, namespace) returns a fully qualified name given a
//      (possibly fully qualified) resource name and optional namespace.
//
//  Examples:
//      (foo, None) => /_/foo
//      (pkg/foo, None) => /_/pkg/foo
//      (foo, ns) => /ns/foo
//      (/ns/pkg/foo, None) => /ns/pkg/foo
//      (/ns/pkg/foo, otherns) => /ns/pkg/foo
func getQualifiedName(name string, namespace string) string {
	name = addLeadSlash(name)
	if strings.HasPrefix(name, "/") {
		return name
	} else if strings.HasPrefix(namespace, "/") {
		return fmt.Sprintf("%s/%s", namespace, name)
	} else {
		if len(namespace) == 0 {
			namespace = getNamespaceFromProp()
		}
		return fmt.Sprintf("/%s/%s", namespace, name)
	}
}
