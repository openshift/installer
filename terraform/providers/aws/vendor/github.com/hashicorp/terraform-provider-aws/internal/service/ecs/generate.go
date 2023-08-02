//go:generate go run ../../generate/listpages/main.go -ListOps=DescribeCapacityProviders
//go:generate go run ../../generate/tagresource/main.go
//go:generate go run ../../generate/tags/main.go -GetTag -ListTags -ServiceTagsSlice -UpdateTags -CreateTags -ParentNotFoundErrCode=InvalidParameterException "-ParentNotFoundErrMsg=The specified cluster is inactive. Specify an active cluster and try again."
// ONLY generate directives and package declaration! Do not add anything else to this file.

package ecs
