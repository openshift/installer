# frozen_string_literal: true

require 'timeout'

# K8sConformanceTest represents the Kubernetes upstream conformance tests
class K8sConformanceTest
  attr_reader :kubeconfig_path, :vpn_tunnel

  def initialize(kubeconfig_path, vpn_tunnel)
    @kubeconfig_path = kubeconfig_path
    @vpn_tunnel = vpn_tunnel
  end

  def run
    ::Timeout.timeout(90 * 60) do # 1 1/2 hour
      image = ENV['KUBE_CONFORMANCE_IMAGE']
      succeeded = system("docker run -v #{@kubeconfig_path}:/kubeconfig #{network_config} #{image}")
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
