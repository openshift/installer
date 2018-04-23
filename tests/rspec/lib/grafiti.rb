# frozen_string_literal: true

require 'yaml'
require 'json'

# Grafiti contains helper functions to use the http://github.com/coreos/grafiti tool
#
class Grafiti
  attr_reader :build_path
  attr_reader :tmp_dir
  attr_reader :config_file_path
  attr_reader :tag_file_path
  attr_reader :terraform_internal_file
  attr_reader :aws_region

  def initialize(build_path, region)
    @aws_region = region
    @build_path = build_path
    @tmp_dir = `mktemp -d -p #{@build_path}`.chomp
    @config_file_path = File.join(@tmp_dir, 'config.toml')
    @tag_file_path = File.join(@tmp_dir, 'tag.json')
    @terraform_internal_file = File.join(build_path, 'internal.yaml')
    write_config_file
    write_tag_file
  end

  def clean
    cmd = 'grafiti'\
          " --config #{@config_file_path}"\
          ' --ignore-errors'\
          ' delete'\
          ' --all-deps'\
          " --delete-file #{@tag_file_path}"

    succeded = system({ 'AWS_REGION' => @aws_region }, cmd)
    raise 'failed to run grafiti delete' unless succeded
  end

  def cluster_id
    config = YAML.load_file(terraform_internal_file)
    config['clusterId'] || ''
  end

  def write_config_file
    IO.write(@config_file_path, 'maxNumRequestRetries = 11')
  end

  def write_tag_file
    tags = { 'TagFilters' => [{ 'Key' => 'tectonicClusterID', 'Values' => [cluster_id] }] }
    IO.write(@tag_file_path, tags.to_json)
  end
end
