# Contributing to terraform-provider-ovirt

Hi and thank you for wanting to contribute to this Terraform provider! This guide will take you through the most important steps of writing code for this library and getting it merged.

## Before you begin

It can be tempting to quickly add a new function to create something in Terraform. However, Terraform is not like Ansible, it isn't just about creating things. In Terraform, you will need to implement the full lifecycle. Think about what happens if a certain parameter of a resource changes? Can you update the resource? Do you have to re-create it? What happens if someone manually destroys the resource on the oVirt Engine and Terraform doesn't know about it?

Or, most importantly, what happens if you need to send two API calls for one resource, but the second one fails? This is why Terraform resources should match API calls as close as possible. Avoid creating composite resources that require sending more than one API call.

If you think about all these, your Terraform resource will be robust. If you don't, you'll see random errors happen.

## Important design consideration

The general rule for this library is: **one API call = one resource**.

Why? Because Terraform does a pretty good job at state management. This saves you from a lot of trouble.

Think of this: you want to create a VM and then resize its disk. What happens if you successfully create the VM, but then fail on the resize? If the two API calls are separate resources Terraform will handle it for you. If you do it in one resource you will have to delete the VM so Terraform can try the whole process again.

Hence, if you can, please try and create separate resources for separate API calls.

## Using go-ovirt-client

This provider is based on the [go-ovirt-client](https://github.com/ovirt/go-ovirt-client) library, a hand-written overlay for the Go oVirt SDK. This library provides many functions we rely on, most importantly mocking the oVirt Engine so we don't have to run one for testing.

You may run into a situation where you don't have the necessary API calls you need to implement a Terraform resource. In this case you must first get your API call into that library. Don't worry, there's a [contributing guide there too](https://github.com/oVirt/go-ovirt-client/blob/main/CONTRIBUTING.md).

Once your change to go-ovirt-client has been merged, you can start developing against it in this Terraform provider by running:

```
go get github.com/ovirt/go-ovirt-client/v2@<your commit hash>
```

## Creating a resource

**👉 Tip:** Even if you don't want to create a new resource, this section is worth reading through.

### Creating a schema

Before you even begin writing actual code, you will need to decide on the schema of your provider. This typically looks like this:

```go
package ovirt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var diskSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	// More schema here
}
```

There are two types of fields: the computed ones and the non-computed ones.

**Computed fields** should be used for attributes that are read-only, such as identifiers automatically assigned by oVirt, or statuses automatically managed by oVirt.

**Non-computed fields** are ones where the user needs to provide the value. You can still update them, but the initial value should in all cases be provided by the user, or the field should have a default value.

When it comes to non-computed fields you should also decide on the update strategy: can you update a resource in-place without deleting and re-creating it? For example, you may be able to change the VM's name without destroying it, but not the template ID it's based on. If your field cannot be updated, you should set the `ForceNew` field to `true`. If you have at least one field which is not computed (`Computed=true`) and `ForceNew` is also not set, you will need to provide an update function.

When writing the schema you should make sure to provide ample description and validation so that users can reasonably write their Terraform code without tripping over low level errors. The [validation.go](ovirt/validation.go) file already contains a number of validators you can add to your schema.

### Adding your resource

Next, we need to declare the resource with the schema we created. We create the resource on the provider struct like this:

```go
func (p *provider) vmResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.vmCreate,
		ReadContext:   p.vmRead,
		UpdateContext: p.vmUpdate,
		DeleteContext: p.vmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.vmImport,
		},
		Schema:      vmSchema,
		Description: "The ovirt_vm resource creates a virtual machine in oVirt.",
	}
}
```

Each of the functions mentioned here (`vmCreate`, `vmRead`, `vmUpdate`, `vmDelete`, and `vmImport`) will need to be implemented here. If the resource doesn't have any fields that can be updated, you can leave the `vmUpdate` function empty. 

Next, you will need to add the resource in [provider.go](ovirt/provider.go);

```go
func (p *provider) provider() *schema.Provider {
	return &schema.Provider{
		Schema:               providerSchema,
		ConfigureContextFunc: p.configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"ovirt_vm": p.vmResource(),
			// More resources here.
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Data sources here
        },
	}
}
```

### Writing the create function

The create function is responsible for creating the resource the first time. The function signature looks like this:

```go
func (p *provider) vmCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	// Code here
}
```

It accepts three parameters:

1. The context. You should pass this context to any go-ovirt-client functions you call using `ovirtclient.ContextStrategy()` as the last parameter.
2. The data record. This is where you can get your parameters from. You will also need to update this data set once your resource has been created. At the very least, you will need to set the `id` field on it so that Terraform knows which ID it belongs to.
3. The unused provider interface. We don't use this as we access the go-ovirt-client over the `p` receiver.

This function returns a list of diagnostics. If there is a diagnostic with the type `diag.Error`, the VM creation will return with an error.

Since you will need to update the `data` record after the resource is done, and this update will need to be done in the update as well, you should create a function like this:

```go
func vmResourceUpdate(vm ovirtclient.VMData, data *schema.ResourceData) diag.Diagnostics {
    diags := diag.Diagnostics{}
    data.SetId(vm.ID())
    diags = setResourceField(data, "cluster_id", vm.ClusterID(), diags)
    //...
    return diags
}
```

### Writing the read and update function

The signature of the read and update functions look exactly the same as the create function. The difference is, that `update` should take the parameters from `data` and update the resource denoted in `id`. Both read and update should then update `data` with the current state of the resource. (This is what you need the `vmResourceUpdate` helper function described above.)

It is worth noting, that in both cases you should explicitly check if the resource has been deleted and set the ID to `""` if that is the case. For read:

```go
vm, err := p.client.GetVM(id, ovirtclient.ContextStrategy(ctx))
if isNotFound(err) {
    data.SetId("")
	// This is fine, return no error
    return nil
}
if err != nil {
    // Handle other errors
}
```

For update:

```go
vm, err := p.client.GetVM(id, ovirtclient.ContextStrategy(ctx))
if isNotFound(err) {
    data.SetId("")
	// Continue processing errors below
}
if err != nil {
    // Handle error and return diagnostics.
}
```

### Writing the delete function

The delete function does exactly what the name says: it takes the ID and possibly other fields from `data` and deletes the resource, then sets the ID to `""`.

### Writing the import function

The import function is a tricky one: the signature is exactly the same as before, but the `data` parameter will contain only a single ID, nothing else. This ID is not necessarily the resource ID, it is whatever the user entered.

You can use this to your advantage when needing multiple parameters on import, for example by splitting the ID by a slash (`/`).

You must then use the provided information to get the current state of the resource and update the `data` records as before.

### Creating blocks (avoid if possible)

There is a special case when you want to create resource blocks, such as this:

```terraform
resource "ovirt_foo" "bar" {
  some_block {
    other_property = "baz"
  }
}
```

This is very tricky to program and should generally be avoided. However, if you need such a block you can define it in the schema as follows:

```go
var fooSchema = map[string]*schema.Schema{
    "some_block": {
        Type:     schema.TypeSet,
        Optional: true,
        MaxItems: 1,
        ForceNew: true,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "other_property": {
                    Type:     schema.TypeString,
                    Optional: true,
                    ForceNew: true,
                },
            },
        },
    },
}
```

There are several limitations with this approach:

1. You cannot define a validation function.
2. You cannot define defaults.
3. You need to handle keys and values manually (see below).

Now, on how to handle these cases. The `some_block` attribute will be a set (or a list). This means that you can have either 0 or 1 entries. (If you remove the `MaxItems` it can have more.) You need to handle both cases.

```go
if someBlockSet, ok := data.GetOk("some_block"); ok {
    someBlockList := someBlockSet.(*schema.Set).List()
    if len(someBlockList) == 1 {
        someBlockEntries := someBlockList[0].(map[string]interface{})
        otherProperty := ""
        if otherPropertyContents, ok := someBlockEntries["other_property"]; ok {
            otherProperty = otherPropertyContents.(string)
        }
        // Use otherProperty here
    }
}
```

Now, this is part 1, and as you can see it's already pretty complicated. Now comes part 2: reading the resource. Here you must make sure that the output of the read produces the exact same output. For example, the oVirt Engine may set a default, but you must ignore that default if the user didn't provide the `some_block` block.

```go
ovirtEngineSomeEntry := ovirtClient.GetSomeEntry()
if rawSomeBlock, ok := data.GetOk("some_block"); ok {
    someBlockList := rawSomeBlock.(*schema.Set).List()
    if len(someBlockList) == 1 {
        // The user provided input.
        someBlockEntry := someBlockList[0].(map[string]interface{})
        // Get the original value
        otherProperty := osEntry["other_property"]
        if otherProperty != ovirtEngineSomeEntry {
            // The engine returned a different value from the user input, set the value.
            data.Set("some_block", []map[string]interface{}{{
                "other_property": ovirtEngineSomeEntry,
            }})
        }
    }
} else if ovirtEngineSomeEntry != "defaultValue" {
    // The user didn't provide input, but the oVirt Engine returned a non-default value.
    data.Set("some_block", []map[string]interface{}{{
        "other_property": ovirtEngineSomeEntry,
    }})
}
```

Did we mention you may want to avoid blocks whenever possible?

## Writing tests

So far so good, you have a resource that works in theory. In practice Terraform can be a tricky beast to deal with though, so you should always write a test for your resource. We exclusively rely on the mocks provided by go-ovirt-client for this functionality, otherwise this provider would be a headache to test.

Additionally, all examples in the [examples](examples) directory are executed automatically against a live engine if you provide the `OVIRT_URL`, `OVIRT_USERNAME`, and `OVIRT_PASSWORD` environment variables.

In order to write a test you must create the appropriate test file and add your test:

```go
func TestVMResource(t *testing.T) {
    
}
```

This is a regular Go test. Next, we will initialize the provider and the test helper:

```go
func TestVMResource(t *testing.T) {
	p := newProvider(ovirtclientlog.NewTestLogger(t))
	
}
```

The provider has multiple functions: first, you can obtain a go-ovirt-sdk client to run API calls for setup/teardown:

```go
client := p.getTestHelper().GetClient()
```

Second, you can use the test helper to get a variety of IDs for testing:

```go
clusterID := p.getTestHelper().GetClusterID()
```

Now that we have this sorted out, let's set up the Terraform tests:

```go
resource.UnitTest(t, resource.TestCase{
    ProviderFactories: p.getProviderFactories(),
        Steps: []resource.TestStep{

        }
    })
```

Here you can add your test steps. Each unit test has a number of options, we'll list the more important ones here:

- **Config**: This is the Terraform config to apply on this step.
- **Destroy**: Set to true to destroy instead of apply.
- **ImportState**: Set to true to import instead of apply. You must set the `ImportStateIdFunc` option.
- **ResourceName**: Contains the name of the resource in the `Config` that the test is meant for. This is especially important for import tests.
- **ImportStateIdFunc**: This function will be run to determine the ID to import. Use this function to create resources to tests against dynamically.
- **Check**: You can add a test function here to verify that the apply/destroy/import was completed successfully. You will have access to the Terraform state here for verification.

**⚠️ Important!** Your `Config` field must include the Terraform `provider {}` section with the `mock = true` option!

When you're done, run `go test -v ./...` to run the tests.

## Generating documentation

Now that your resource works, tests are done, the only thing left to do is generate the documentation. Go ahead and run `go generate`.

## Submitting your PR

From here it's simple: push to your fork and submit a PR on GitHub. Follow the description there and we'll review your change in short order.
