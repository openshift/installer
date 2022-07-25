# aws-cloudformation-resource-schema-sdk-go

This package provides [AWS CloudFormation Resource Schema](https://github.com/aws-cloudformation/aws-cloudformation-resource-schema/) functionality in Go, including the validation of schema documents, parsing of schema documents into native Go types, and offering methods for interacting with these schemas.

_NOTE: There is a separate [AWS CloudFormation resource specification](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-resource-specification.html), which is different than what is being described or handled in this package._

To browse the documentation before it is published on https://pkg.go.dev, you can browse it locally via [`godoc`](https://pkg.go.dev/golang.org/x/tools/cmd/godoc):

```shell
godoc -http=":6060" &
open "http://localhost:6060/pkg/github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go"
```

## Example Usage

Adding Go Module dependency to your project:

```shell
go get github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go
```

Adding the import, using an import alias to simplify usage:

```go
import {
  # ... other imports ...
  cfschema "github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go"
}
```

Loading the meta-schema:

```go
metaSchema, err := cfschema.NewMetaJsonSchemaPath("provider.definition.schema.v1.json")
```

Quickly validating a resource schema file path against the meta-schema:

```go
err := metaSchema.ValidateResourcePath("aws-logs-loggroup.json")
```

Loading a resource schema for further processing:

```go
resourceSchema, err := cfschema.NewResourceJsonSchemaPath("aws-logs-loggroup.json")
```

Validating a loaded resource schema against the meta-schema:

```go
err := metaSchema.ValidateResourceJsonSchema(resourceSchema)
```

Validating a configuration against a loaded resource schema:

```go
err := resourceSchema.ValidateConfigurationDocument("{...}")
```

Parsing the resource schema into Go:

```go
resource, err := resourceSchema.Resource()
```

Expanding a resource schema to replace JSON Pointer references:

```go
err := resource.Expand()
```

## Further Reading

### CloudFormation Resource Providers Schema

The specification for CloudFormation Resource Types is based on the [CloudFormation Resource Providers Meta-Schema](https://github.com/aws-cloudformation/aws-cloudformation-resource-schema/blob/master/src/main/resources/schema/provider.definition.schema.v1.json), which defines all the valid fields and structures for a resource schema. Additional information about creating these schemas files can be found in the [Modeling resource types for use in AWS CloudFormation](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/resource-type-model.html) documentation.

Conceptually, the naming, typing, and some validation of attributes in the CloudFormation Schema are a flattened set of properties and re-usable definitions. Any nesting or re-usability is defined through JSON Pointer references. Additional concepts such as read-only, write-only, create-only attributes are implemented at the resource level and reference (potentially nested) attributes using JSON Pointers. The [`Initech::TPS::Report` example resource schema](https://github.com/aws-cloudformation/aws-cloudformation-resource-schema/blob/master/src/main/resources/examples/resource/initech.tps.report.v1.json) can provide high level insight into the structure of these files.

### JSON Pointers

CloudFormation Resource Providers Schemas make extensive use of JavaScript Object Notation (JSON) Pointers as described in [RFC 6901](https://tools.ietf.org/html/rfc6901).

### JSON Schema

[JSON Schema](http://json-schema.org/) is a vocabulary that allows you to annotate and validate JSON documents. This is the specification on which CloudFormation Schemas are built. Understanding the core concepts and high level implementation details of this specification will provide a much clearer picture into the details of this Go package.

Some helpful resources for learning JSON Schema include:

- [Understanding JSON Schema](https://json-schema.org/understanding-json-schema/)
- [Getting Started Step-By-Step](https://json-schema.org/learn/getting-started-step-by-step)
