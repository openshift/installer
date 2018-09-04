load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "zz_generated.api.register.go",
    ],
    importmap = "installer/vendor/sigs.k8s.io/cluster-api/pkg/apis/cluster/install",
    importpath = "sigs.k8s.io/cluster-api/pkg/apis/cluster/install",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apimachinery/pkg/apimachinery/announced:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apimachinery/registered:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/apis:go_default_library",
    ],
)
