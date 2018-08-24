workspace(name = "installer")

terrafom_version = "0.11.7"

supported_platforms = [
    "linux",
    "darwin",
]

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.1/rules_go-0.12.1.tar.gz",
    sha256 = "8b68d0630d63d95dacc0016c3bb4b76154fe34fca93efd65d1c366de3fcb4294",
)

http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.12.0/bazel-gazelle-0.12.0.tar.gz",
    sha256 = "ddedc7aaeb61f2654d7d7d4fd7940052ea992ccdb031b8f9797ed143ac7e8d43",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# Runtime binary dependencies follow.
# These will be fetched and included in the build output verbatim.
#
[new_http_archive(
    name = "terraform_runtime_%s" % platform,
    build_file_content = """exports_files(["terraform"], visibility = ["//visibility:public"])""",
    type = "zip",
    url = "https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_amd64.zip" % (terrafom_version, terrafom_version, platform),
) for platform in supported_platforms]
