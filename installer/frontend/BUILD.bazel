package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_nodejs//:defs.bzl", "nodejs_binary")

genrule(
    name = "yarn-build",
    message = "Installing JS dependencies...",
    srcs = [":node_modules"],
    outs = ["installer/assets/frontend/scripts/app-bundle.js"],
    cmd = "$(location @yarn//:yarn) run build",
    tools = [
        "@yarn//:yarn",
        ":browserify",
    ],
)

nodejs_binary(
    name = "browserify",
    entry_point = "browserify/.bin/browserify",
)
