# frozen_string_literal: true

require 'timeout'

# K8sConformanceTest represents the Kubernetes upstream conformance tests
class K8sConformanceTest
  attr_reader :kubeconfig_path

  def initialize(kubeconfig_path)
    @kubeconfig_path = kubeconfig_path
  end

  def run
    ::Timeout.timeout(90 * 60) do # 1 1/2 hour
      succeeded = system("docker run -v #{@kubeconfig_path}:/kubeconfig #{ENV['KUBE_CONFORMANCE_IMAGE']}")
      raise 'Running k8s conformance tests failed' unless succeeded
    end
  end
end
