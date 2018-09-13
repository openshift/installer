load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importmap = "installer/vendor/sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions/cluster",
    importpath = "sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions/cluster",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions/cluster/v1alpha1:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions/internalinterfaces:go_default_library",
    ],
)
