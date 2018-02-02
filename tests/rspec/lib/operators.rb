# frozen_string_literal: true

# Operators contains helper functions to test creation of the CoreOS operators
module Operators
  BOOTKUBE_OPERATOR_NAMES = [
    'tectonic-network-operator'
  ].freeze

  OPERATOR_NAMES = [
    'kube-core-operator',
    'kubernetes-addon-operator',
    'tectonic-channel-operator',
    'tectonic-cluo-operator',
    'tectonic-prometheus-operator',
    'tectonic-utility-operator'
  ].freeze

  def self.manifests_generated?(manifest_path)
    BOOTKUBE_OPERATOR_NAMES.each do |operator_name|
      file_path = File.join(
        manifest_path, 'manifests', "#{operator_name}.yaml"
      )
      next if File.exist?(file_path)

      raise "could not find manifest for #{operator_name}"
    end

    OPERATOR_NAMES.each do |operator_name|
      file_path = File.join(
        manifest_path, 'tectonic/updater/operators', "#{operator_name}.yaml"
      )
      next if File.exist?(file_path)

      raise "could not find manifest for #{operator_name}"
    end
  end
end
