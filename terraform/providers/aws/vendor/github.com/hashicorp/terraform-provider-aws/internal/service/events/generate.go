//go:generate go run ../../generate/listpages/main.go -ListOps=ListEventBuses,ListRules,ListTargetsByRule
//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsInIDElem=ResourceARN -ServiceTagsSlice -TagInIDElem=ResourceARN -UpdateTags -CreateTags
// ONLY generate directives and package declaration! Do not add anything else to this file.

package events
