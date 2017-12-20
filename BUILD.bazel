package(
    default_visibility = ["//visibility:public"],
)

template_files = glob([
    "modules/**/*",
    "platforms/**/*",
])

config_setting(
    name = "darwin",
    values = {"cpu": "darwin"},
    visibility = ["//visibility:public"],
)

config_setting(
    name = "linux",
    values = {"cpu": "k8"}, # don't ask...
    visibility = ["//visibility:public"],
)

genrule(
    name = "terraform_runtime",
    output_to_bindir = 1,
    srcs = select({
        "//:linux": ["@terraform_runtime_linux//:terraform"],
        "//:darwin": ["@terraform_runtime_darwin//:terraform"],
    }),
    outs = ["bin/terraform"],
    cmd = "cp $(<) $(@)",
)

genrule(
    name = "provider_matchbox",
    output_to_bindir = 1,
    srcs =  select({
        "//:linux": ["@terraform_provider_matchbox_linux//:terraform-provider-matchbox"],
        "//:darwin": ["@terraform_provider_matchbox_darwin//:terraform-provider-matchbox"],
    }),
    outs = ["bin/terraform-provider-matchbox"],
    cmd = "cp $(<) $(@)",
)

genrule(
    name = "templates",
    message = "Copying templates...",
    output_to_bindir = 1,
    srcs = template_files,
    outs = ["templates/%s" % f for f in template_files],
    cmd = '\n'.join([
        "for tf_file in $(SRCS); do",
        "target=\"$(@D)/templates/$$(dirname $${tf_file})\"",
        "mkdir -p $${target}",
        "cp $${tf_file} $${target}",
        "done"
    ]),
)

load("@io_bazel_rules_go//go:def.bzl", "go_prefix")

go_prefix("github.com/coreos/tectonic-installer")

alias(
    name = "smoke_tests",
    actual = "//tests/smoke:smoke",
)

alias(
    name = "backend",
    actual = "//installer/cmd/installer:installer",
)

# This absolutely must exist in the top package
# Moving it into installer/frontend currently won't work
#
filegroup(
    name = "node_modules", 
    srcs = glob(["node_modules/**/*"])
)
