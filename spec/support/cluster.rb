require "open3"
require "yaml"

class Cluster
  Error = Class.new(StandardError)

  attr_reader :cluster_name

  def initialize(config_path, *args)
    config = YAML.load_file(config_path)
    @cluster_name = config.fetch("name")

    run(:init, "--config=#{config_path}", *args)
  end

  def destroy(*args)
    run_cluster_command(:destroy, *args)
  end

  def install(*args)
    run_cluster_command(:install, *args)

    return unless block_given?

    begin
      yield
    ensure
      destroy
    end
  end

  private

  def run(*args)
    cmd = [File.join("installer", "tectonic"), *args].join(" ")
    Open3.popen2e(cmd) do |_, stdout_stderr, wait_thr|
      buffer = ""

      Thread.new do
        stdout_stderr.each { |line| buffer << line }
      end

      raise Error, buffer unless wait_thr.value.success?
    end
  end

  def run_cluster_command(action, *args)
    run(action, "--dir=#{cluster_name}", *args)
  end
end
