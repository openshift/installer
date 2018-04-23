workspace(name = "installer")

terrafom_version = "0.11.7"

supported_platforms = [
    "linux",
    "darwin",
]

# We need a feature from recent master to get rid of the platform specific part of the go_binary path.
# https://github.com/bazelbuild/rules_go/pull/1393

git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    commit = "361556b0d27318d1b8fa42c91a4baa4ab5ea1c58"
)

# http_archive(
#     name = "io_bazel_rules_go",
#     url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.1/rules_go-0.10.1.tar.gz",
#     sha256 = "4b14d8dd31c6dbaf3ff871adcd03f28c3274e42abc855cb8fb4d01233c0154dc",
# )

http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.10.1/bazel-gazelle-0.10.1.tar.gz",
    sha256 = "d03625db67e9fb0905bbd206fa97e32ae9da894fe234a493e7517fd25faec914",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

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
