def generate_script(command):
  return """
set -ex
{command}
""".format(command=command)

def _impl(ctx):
  script = generate_script(ctx.attr.command)

  # Write the file, it is executed by 'bazel test'.
  ctx.actions.write(
      output=ctx.outputs.executable,
      content=script
  )

  # To ensure the files needed by the script are available, we put them in
  # the runfiles.
  runfiles = ctx.runfiles(files=ctx.files.deps)
  return [DefaultInfo(runfiles=runfiles)]

gen_test = rule(
    implementation=_impl,
    attrs={
        "command": attr.string(),
        "deps": attr.label_list(allow_files=True),
    },
  test=True,
)
