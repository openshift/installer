# InstallConfig Explain

The installer allows it's users to see all the configuration available using the `explain` subcommand. The command works almost like the
the `oc explain` subcommand.

## Generating the documentation

Like `oc explain` and Custom Resource Definitions, the installer also generates an internal only Custom Resource Definition for the `InstallConfig`.

```sh
go generate ./pkg/types/installconfig.go
```

The code generation uses the upstream project kubebuilder, which provides tools like controller-tools to [generate][kubebuilder-generate-crd] Custom Resource Definitions.

## Descriptions of fields and various types

The generator uses the Godoc comments on the fields and types to create the descriptions. Therefore, making sure that everything used by the `InstallConfig` definition has user friendly comments will automatically transfer to user friendly explain output.

## Kubebuilder markers

The generator allows use of various markers for fields and types to provide information about the valid inputs. Various available markers are defined [here][kubebuilder-validation-markers]

## Additional guidelines on godoc comments

All definitions in installer codebase already follow and enforce [effective-go] guidelines, but for better explain output for the types these additional guidelines should also be followed

1. No external types are allowed to be used in the `InstallConfig`. All the types used should be defined in the installer repository itself. The only exception is the `TypeMeta` and `ObjectMeta`.

2. Fields should always have `json` tags that follow [mixedCaps][go-mixed-caps] and should always start with lowercase.

3. The comments on the fields must use the `json` tag to reference the field.

4. Optional fields must have the `+optional` marker defined. Optional fields must also define the default value. If the values are static, `+kubebuilder:validation=Default` marker should be used, otherwise the comment must clearly define how the default value will be calculated.

5. Only string based enum types are allowed. Also such enum type must use `+kubebuilder:validation:Enum` marker.

[go-effective]: https://golang.org/doc/effective_go.html
[go-mixed-caps]: https://golang.org/doc/effective_go.html#mixed-caps
[kubebuilder-generate-crd]: https://book.kubebuilder.io/reference/generating-crd.html
[kubebuilder-validation-markers]: https://book.kubebuilder.io/reference/markers/crd-validation.html
