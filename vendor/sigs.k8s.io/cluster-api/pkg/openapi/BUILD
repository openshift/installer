load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "openapi_generated.go",
    ],
    importmap = "installer/vendor/sigs.k8s.io/cluster-api/pkg/openapi",
    importpath = "sigs.k8s.io/cluster-api/pkg/openapi",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/go-openapi/spec:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//vendor/k8s.io/kube-openapi/pkg/common:go_default_library",
    ],
)
