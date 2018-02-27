# frozen_string_literal: true

# TFStateFile represents a Terraform state file
class TFStateFile
  def initialize(build_path, tfstate_file_name = 'terraform.tfstate')
    @build_path = build_path
    @tfstate_file_name = tfstate_file_name
  end

  def value(address, wanted_key)
    file_exists?

    Dir.chdir(@build_path) do
      state = `terraform state show -state=#{@tfstate_file_name} #{address}`.chomp.split("\n")
      state.each do |value|
        key, value = value.split('=')
        key = key.strip.chomp
        value = value.strip.chomp
        return value if key == wanted_key
      end
    end

    msg = "could not find value for key \"#{wanted_key}\" in tfstate file #{@build_path}"
    raise TFStateFileValueForKeyDoesNotExist, msg
  end

  def output(module_name, wanted_key)
    Dir.chdir(@build_path) do
      out = `terraform output -module=#{module_name} -state=#{@tfstate_file_name} -json`.chomp
      out = JSON.parse(out)
      return out[wanted_key]['value']
    end
  end

  def file_exists?
    tfstate_file = File.join(@build_path, @tfstate_file_name)
    raise "tfstate file #{tfstate_file} does not exist" unless File.exist? tfstate_file
  end
end

class TFStateFileValueForKeyDoesNotExist < StandardError; end
