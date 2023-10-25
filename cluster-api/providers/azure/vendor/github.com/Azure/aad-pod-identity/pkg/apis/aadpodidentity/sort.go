package aadpodidentity

type AzureIdentityBindings []AzureIdentityBinding

func (a AzureIdentityBindings) Len() int {
	return len(a)
}

func (a AzureIdentityBindings) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AzureIdentityBindings) Less(i, j int) bool {
	if a[i].Namespace == a[j].Namespace {
		return a[i].Name < a[j].Name
	}
	return a[i].Namespace < a[j].Namespace
}
