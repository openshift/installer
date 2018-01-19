workspace(name = "tectonic_installer")

terrafom_version = "0.10.8"

provider_matchbox_version = "0.2.2"

supported_platforms = [
    "linux",
    "darwin",
]

# Latest working commit for cross compilation [0] until fix [1] lands.
# [0] https://github.com/bazelbuild/rules_go/issues/1240#issuecomment-357789209
# [1] https://github.com/bazelbuild/rules_go/pull/1248
git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    commit = "3f38260eda98d23e9142bb905caede5912508770"
)

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "7550c6f7904f602b69c2a69d92f7c739db87479336554c7a31d0649519ec414e",
    type = "tar.gz",
    strip_prefix = "rules_nodejs-0.3.1",
    url = "https://codeload.github.com/bazelbuild/rules_nodejs/tar.gz/0.3.1",
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "4952295aa35241082eefbb53decd7d4dd4e67a1f52655d708a1da942e3f38975",
    url = "https://github.com/bazelbuild/bazel-gazelle/archive/eaa1e87d2a3ca716780ca6650ef5b9b9663b8773.zip",
    strip_prefix = "bazel-gazelle-eaa1e87d2a3ca716780ca6650ef5b9b9663b8773",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

node_repositories(package_json = ["//installer/frontend:package.json"])

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

[new_http_archive(
    name = "terraform_provider_matchbox_%s" % platform,
    build_file_content = """exports_files(
["terraform-provider-matchbox"],
visibility = ["//visibility:public"]
)""",
    strip_prefix = "terraform-provider-matchbox-v%s-%s-amd64/" % (provider_matchbox_version, platform),
    url = "https://github.com/coreos/terraform-provider-matchbox/releases/download/v%s/terraform-provider-matchbox-v%s-%s-amd64.tar.gz" % (provider_matchbox_version, provider_matchbox_version, platform),
) for platform in supported_platforms]

# Import Go dependencies.
go_repository(
    name = "com_github_hashicorp_hcl",
    importpath = "github.com/hashicorp/hcl",
    commit = "23c074d0eceb2b8a5bfdbb271ab780cde70f05a8",
)

go_repository(
    name = "com_github_segmentio_terraform-docs",
    importpath = "github.com/segmentio/terraform-docs",
    tag = "v0.3.0",
)

go_repository(
    name = "com_github_tj_docopt",
    importpath = "github.com/tj/docopt",
    tag = "v1.0.0",
)
