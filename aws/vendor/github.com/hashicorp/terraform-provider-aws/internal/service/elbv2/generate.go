//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=DescribeTags -ListTagsInIDElem=ResourceArns -ListTagsInIDNeedSlice=yes -ListTagsOutTagsElem=TagDescriptions[0].Tags -ServiceTagsSlice -TagOp=AddTags -TagInIDElem=ResourceArns -TagInIDNeedSlice=yes -UntagOp=RemoveTags -UpdateTags -CreateTags
// ONLY generate directives and package declaration! Do not add anything else to this file.

package elbv2
