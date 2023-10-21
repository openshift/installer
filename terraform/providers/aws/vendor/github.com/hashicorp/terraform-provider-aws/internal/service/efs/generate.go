//go:generate go run ../../generate/listpages/main.go -ListOps=DescribeMountTargets -InputPaginator=Marker -OutputPaginator=NextMarker
//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=DescribeTags -ListTagsInIDElem=FileSystemId -ServiceTagsSlice -TagInIDElem=ResourceId -UpdateTags
// ONLY generate directives and package declaration! Do not add anything else to this file.

package efs
