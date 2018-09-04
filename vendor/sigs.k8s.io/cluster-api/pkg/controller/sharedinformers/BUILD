load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "informers.go",
        "zz_generated.api.register.go",
    ],
    importmap = "installer/vendor/sigs.k8s.io/cluster-api/pkg/controller/sharedinformers",
    importpath = "sigs.k8s.io/cluster-api/pkg/controller/sharedinformers",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/kubernetes-incubator/apiserver-builder/pkg/controller:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset:go_default_library",
        "//vendor/sigs.k8s.io/cluster-api/pkg/client/informers_generated/externalversions:go_default_library",
    ],
)
