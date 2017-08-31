# Grafiti contains helper functions to use the http://github.com/coreos/grafiti
# tool
class Grafiti
  attr_reader :build_path
  attr_reader :tmp_dir
  attr_reader :config_file_path
  attr_reader :tag_file_path
  attr_reader :terraform_log_path
  attr_reader :aws_region

  def initialize(build_path, region)
    @aws_region = region
    @build_path = build_path
    @tmp_dir = `mktemp -d -p #{@build_path}`.chomp
    @config_file_path = File.join(@tmp_dir, 'config.toml')
    @tag_file_path = File.join(@tmp_dir, 'tag.json')
    @terraform_log_path = File.join(build_path, 'terraform.tfstate')
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
    cmd = 'grep -m 1 -oP'\
          ' \'^ *\"tags\.tectonicClusterID\": \"\K[0-9a-z-]*(?=\",$)\''\
          " #{terraform_log_path}"
    `#{cmd}`.chomp
  end

  def write_config_file
    IO.write(@config_file_path, 'maxNumRequestRetries = 11')
  end

  def write_tag_file
    content = '{"TagFilters":['\
              "  {\"Key\":\"tectonicClusterID\",\"Values\":[\"#{cluster_id}\""\
              ']}]}'
    IO.write(@tag_file_path, content)
  end
end
