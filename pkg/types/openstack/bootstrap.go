package openstack

// ResolveBootstrapFlavor returns the effective flavor to use for the bootstrap
// machine. If the Platform has an explicit BootstrapFlavor set, that value is
// returned. Otherwise, the provided controlPlaneFlavor is returned as a
// fallback. A nil receiver is handled gracefully and returns controlPlaneFlavor.
func (p *Platform) ResolveBootstrapFlavor(controlPlaneFlavor string) string {
	if p != nil && p.BootstrapFlavor != "" {
		return p.BootstrapFlavor
	}
	return controlPlaneFlavor
}
