# frozen_string_literal: true

require 'timeout'

# K8sConformanceTest represents the Kubernetes upstream conformance tests
class K8sConformanceTest
  attr_reader :kubeconfig_path, :vpn_tunnel

  def initialize(kubeconfig_path, vpn_tunnel, platform)
    @kubeconfig_path = kubeconfig_path
    @vpn_tunnel = vpn_tunnel
    @platform = platform
  end

  def run
    ::Timeout.timeout(2 * 60 * 60) do # 2 hour
      image = ENV['KUBE_CONFORMANCE_IMAGE']
      command = if @platform.include?('metal')
                  "sudo rkt run --volume kubecfg,kind=host,readOnly=false,source=#{@kubeconfig_path} \
                  --mount volume=kubecfg,target=/kubeconfig #{network_config} --dns=host \
                  --insecure-options=image #{image}"
                else
                  "docker run -v #{@kubeconfig_path}:/kubeconfig #{network_config} #{image}"
                end

      succeeded = system(command)
      raise 'Running k8s conformance tests failed' unless succeeded
    end
  end

  private

  # When the cluster is only reachable via a VPN connection, the
  # kube-conformance container has to share the same linux network namespace
  # like the current container to be able to use the same VPN tunnel.
  def network_config
    return '--net=host' unless @vpn_tunnel

    hostname = `hostname`.chomp
    "--net=container:#{hostname}"
  end
end
