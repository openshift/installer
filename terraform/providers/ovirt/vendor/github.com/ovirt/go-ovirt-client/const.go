package ovirtclient

// DefaultBlankTemplateID returns the ID for the factory-default blank template. This should not be used
// as the template may be deleted from the oVirt engine. Instead, use the API call to find the blank template.
const DefaultBlankTemplateID TemplateID = "00000000-0000-0000-0000-000000000000"

// MinDiskSizeOVirt defines the minimum size of 1M for disks in oVirt. Smaller disks can be created, but they
// lead to bugs in oVirt when creating disks from templates and changing the format.
const MinDiskSizeOVirt uint64 = 1048576
