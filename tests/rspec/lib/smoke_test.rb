# frozen_string_literal: true

require 'timeout'

# SmokeTest contains helper functions to operate the smoke tests written in
# golang
module SmokeTest
  def self.run(cluster)
    ::Timeout.timeout(30 * 60) do # 30 minutes
      os = Gem::Platform.local.os
      succeeded = system(
        env_variables(cluster),
        File.join(
          File.dirname(ENV['RELEASE_TARBALL_PATH']), "tests/smoke/#{os}_amd64_stripped/smoke"
        ) + ' -test.v -test.parallel=1 --cluster'
      )
      raise 'SmokeTests failed' unless succeeded
    end
  end

  def self.env_variables(cluster)
    env = if cluster.env_variables['PLATFORM'].include?('aws')
            {
              'SMOKE_KUBECONFIG' => cluster.kubeconfig,
              'SMOKE_NODE_COUNT' => cluster.config_file.node_count.to_s,
              'SMOKE_MANIFEST_PATHS' => cluster.manifest_path,
              'SMOKE_NETWORKING' => cluster.config_file.networking.to_s
            }
          else
            {
              'SMOKE_KUBECONFIG' => cluster.kubeconfig,
              'SMOKE_NODE_COUNT' => cluster.tfvars_file.node_count.to_s,
              'SMOKE_MANIFEST_PATHS' => cluster.manifest_path,
              'SMOKE_NETWORKING' => cluster.tfvars_file.networking.to_s
            }
          end
    env
  end
end
