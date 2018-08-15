package installerassets

func init() {
	Rebuilders["machines/master-count"] = PlatformOverrideRebuilder(
		"machines/master-count",
		ConstantDefault([]byte("3")),
	)
}
