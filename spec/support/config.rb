require "securerandom"
require "yaml"

class Config
  BASE_PATH = File.join("..", "..", "examples", "tectonic.aws.yaml")

  attr_reader :config, :uuid

  def initialize(options = {})
    @uuid = SecureRandom.hex(8)

    @config = defaults.merge(options)
  end

  def init(*_args)
    write

    Cluster.new(path)
  end

  private

  def defaults
    example.merge("name" => uuid,
                  "baseDomain" => ENV.fetch("BASE_DOMAIN") { "tectonic-ci.de" },
                  "aws" => { "region" => ENV.fetch("AWS_REGION") })
  end

  def example
    YAML.load_file(File.expand_path(BASE_PATH, __dir__))
  end

  def path
    "examples/#{uuid}.yaml"
  end

  def write
    File.open(path, "w") { |f| f.write(config.to_yaml) }
  end
end
