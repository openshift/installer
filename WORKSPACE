workspace(name = "tectonic_installer")

terrafom_version = "0.10.8"

provider_matchbox_version = "0.2.2"

supported_platforms = [
    "linux",
    "darwin",
]

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "91fca9cf860a1476abdc185a5f675b641b60d3acf0596679a27b580af60bf19c",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.0/rules_go-0.7.0.tar.gz",
)

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "0ee9f6d9a34994a338b374ecb59df814dd5ba2952ad57fe27ddf44ef858a2c09",
    type = "tar.gz",
    strip_prefix = "rules_nodejs-0.2.2",
    url = "https://codeload.github.com/bazelbuild/rules_nodejs/tar.gz/0.2.2",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

node_repositories(package_json = ["//installer/frontend:package.json"])

# Runtime binary dependencies follow.
# These will be fetched and included in the build output verbatim.
#
[new_http_archive(
    name = "terraform_runtime_%s" % platform,
    build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
    type = "zip",
    url = "https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_amd64.zip" % (terrafom_version, terrafom_version, platform),
) for platform in supported_platforms]

[new_http_archive(
    name = "terraform_provider_matchbox_%s" % platform,
    build_file_content = """exports_files(
["terraform-provider-matchbox"],
visibility = ["//visibility:public"]
)""",
    strip_prefix = "terraform-provider-matchbox-v%s-%s-amd64/" % (provider_matchbox_version, platform),
    url = "https://github.com/coreos/terraform-provider-matchbox/releases/download/v%s/terraform-provider-matchbox-v%s-%s-amd64.tar.gz" % (provider_matchbox_version, provider_matchbox_version, platform),
) for platform in supported_platforms]
