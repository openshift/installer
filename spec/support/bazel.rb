require "open3"

module Bazel
  class Tarball
    def initialize(path)
      @path = path
    end

    def untar(dir)
      stdout, stderr, status = Open3.capture3("tar -zxf #{@path} -C #{dir}")
      return if status.success?

      STDOUT.puts(stdout)
      STDERR.puts(stderr)

      exit(1)
    end
  end

  module_function

  def build(target)
    case target
    when :tarball
      build_tarball
    else
      raise "target not supported"
    end
  end

  def build_tarball
    stdout, stderr, status = Open3.capture3("bazel build #{target}")
    unless status.success?
      STDOUT.puts(stdout)
      STDERR.puts(stderr)

      exit(1)
    end

    Tarball.new("bazel-bin/tectonic-dev.tar.gz")
  end
end
