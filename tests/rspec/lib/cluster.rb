# frozen_string_literal: true

require 'cluster_support'
require 'fileutils'
require 'jenkins'
require 'kubectl_helpers'
require 'name_generator'
require 'password_generator'
require 'securerandom'
require 'ssh'
require 'tfstate_file'
require 'tfvars_file'
require 'timeout'
require 'with_retries'
require 'open3'

# Cluster represents a k8s cluster
class Cluster
  attr_reader :tfvars_file, :kubeconfig, :manifest_path, :build_path,
              :tectonic_admin_email, :tectonic_admin_password, :tfstate_file

  def initialize(tfvars_file)
    @tfvars_file = tfvars_file

    # Enable local testers to specify a static cluster name
    # S3 buckets can only handle lower case names
    @name = ENV['CLUSTER'] || NameGenerator.generate(@tfvars_file.prefix)
    @tectonic_admin_email = ENV['TF_VAR_tectonic_admin_email'] || NameGenerator.generate_fake_email
    @tectonic_admin_password = ENV['TF_VAR_tectonic_admin_password'] || PasswordGenerator.generate_password

    @build_path = File.join(File.realpath('../../'), "build/#{@name}")
    @manifest_path = File.join(@build_path, 'generated')
    @kubeconfig = File.join(manifest_path, 'auth/kubeconfig')
    @tfstate_file = TFStateFile.new(@build_path)

    check_prerequisites
  end

  def plan(terraform_options = nil)
    env = env_variables
    env['TF_PLAN_OPTIONS'] = terraform_options unless terraform_options.nil?
    stdout, stderr, exit_status = Open3.capture3(env, 'make -C ../.. plan')
    [stdout, stderr, exit_status]
  end

  def start
    apply
    wait_til_ready
  end

  def update_cluster
    start
  end

  def stop
    if ENV.key?('TECTONIC_TESTS_DONT_CLEAN_UP')
      puts "*** Cleanup inhibiting flag set. Stopping here. ***\n"
      puts '*** Your email/password to use in the tectonic console is:'\
           "#{@tectonic_admin_email} / #{@tectonic_admin_password} ***\n"
      return
    end
    destroy
  end

  def check_prerequisites
    return if license_and_pull_secret_defined?
    raise 'Tectonic license and pull secret are not defined as environment'\
          'variables.'
  end

  def env_variables
    {
      'CLUSTER' => @name,
      'TF_VAR_tectonic_cluster_name' => @name,
      'TF_VAR_tectonic_admin_email' => @tectonic_admin_email,
      'TF_VAR_tectonic_admin_password' => @tectonic_admin_password
    }
  end

  def tf_var(v)
    tf_value "var.#{v}"
  end

  def tf_value(v)
    Dir.chdir(@build_path) do
      `echo '#{v}' | terraform console ../../platforms/#{env_variables['PLATFORM']}`.chomp
    end
  end

  def secret_files(namespace, secret)
    cmd = "get secret -n #{namespace} #{secret} -o go-template "\
          "\'--template={{range $key, $value := .data}}{{$key}}\n{{end}}\'"
    KubeCTL.run(@kubeconfig, cmd).split("\n")
  end

  def api_ip_addresses
    nodes = KubeCTL.run(
      @kubeconfig,
      'get node -l=node-role.kubernetes.io/master '\
      '-o jsonpath=\'{range .items[*]}'\
      '{@.metadata.name}{"\t"}{@.status.addresses[?(@.type=="ExternalIP")].address}'\
      '{"\n"}{end}\''
    )

    nodes = nodes.split("\n").map { |node| node.split("\t") }.to_h

    api_pods = KubeCTL.run(
      @kubeconfig,
      'get pod -n kube-system -l k8s-app=kube-apiserver '\
      '-o \'jsonpath={range .items[*]}'\
      '{@.metadata.name}{"\t"}{@.spec.nodeName}'\
      '{"\n"}{end}\''
    )

    api_pods
      .split("\n")
      .map { |pod| pod.split("\t") }
      .map { |pod| [pod[0], nodes[pod[1]]] }.to_h
  end

  def forensic(events = true)
    save_kubernetes_events(@kubeconfig, @name) if events

    master_ip_addresses.each do |master_ip|
      save_docker_logs(master_ip, @name)

      ['bootkube', 'tectonic', 'kubelet', 'k8s-node-bootstrap'].each do |service|
        print_service_logs(master_ip, service, @name)
      end
    end

    worker_ip_addresses.each do |worker_ip|
      save_docker_logs(worker_ip, @name, master_ip_address)

      ['kubelet'].each do |service|
        print_service_logs(worker_ip, service, @name, master_ip_address)
      end
    end

    etcd_ip_addresses.each do |etcd_ip|
      ['etcd-member'].each do |service|
        print_service_logs(etcd_ip, service, @name, master_ip_address)
      end
    end
  end

  private

  def license_and_pull_secret_defined?
    license_path = 'TF_VAR_tectonic_license_path'
    pull_secret_path = 'TF_VAR_tectonic_pull_secret_path'

    EnvVar.set?([license_path, pull_secret_path])
  end

  def apply
    ::Timeout.timeout(30 * 60) do # 30 minutes
      3.times do |idx|
        env = env_variables
        env['TF_LOG'] = 'TRACE' if idx.positive?
        env['TF_APPLY_OPTIONS'] = '-no-color'
        env['TF_INIT_OPTIONS'] = '-no-color'
        return true if system(env, "bash -co pipefail 'make -C ../.. apply |
                                    tee ../../build/#{@name}/terraform-apply.log'")
      end
    end
  rescue Timeout::Error
    forensic(false)
    raise 'Applying cluster failed'
  end

  def destroy
    ::Timeout.timeout(30 * 60) do # 30 minutes
      3.times do |idx|
        env = env_variables
        env['TF_LOG'] = 'TRACE' if idx.positive?
        env['TF_DESTROY_OPTIONS'] = '-no-color'
        env['TF_INIT_OPTIONS'] = '-no-color'
        return true if system(env, "bash -co pipefail 'make -C ../.. destroy |
                                    tee ../../build/#{@name}/terraform-destroy.log'")
      end
    end

    recover_from_failed_destroy
    raise 'Destroying cluster failed'
  rescue => e
    recover_from_failed_destroy
    raise e
  end

  def recover_from_failed_destroy() end

  def wait_til_ready
    wait_for_bootstrapping

    from = Time.now
    loop do
      begin
        KubeCTL.run(@kubeconfig, 'cluster-info')
        break
      rescue KubeCTL::KubeCTLCmdError
        elapsed = Time.now - from
        raise 'kubectl cluster-info never returned with successful error code' if elapsed > 1200 # 20 mins timeout
        sleep 10
      end
    end

    wait_nodes_ready
  end

  def wait_nodes_ready
    from = Time.now
    loop do
      puts 'Waiting for nodes become in ready state after an update'
      Retriable.with_retries(KubeCTL::KubeCTLCmdError, limit: 5, sleep: 10) do
        nodes = describe_nodes
        nodes_ready = Array.new(@tfvars_file.node_count, false)
        nodes['items'].each_with_index do |item, index|
          item['status']['conditions'].each do |condition|
            if condition['type'] == 'Ready' && condition['status'] == 'True'
              nodes_ready[index] = true
            end
          end
        end
        if nodes_ready.uniq.length == 1 && nodes_ready.uniq.include?(true)
          puts '**All nodes are Ready!**'
          return true
        end
        puts "One or more nodes are not ready yet or missing nodes. Waiting...\n" \
             "# of returned nodes #{nodes['items'].size}. Expected #{@tfvars_file.node_count}"
        elapsed = Time.now - from
        raise 'waiting for all nodes to become ready timed out' if elapsed > 1200 # 20 mins timeout
        sleep 20
      end
    end
  end

  def wait_for_bootstrapping
    ips = master_ip_addresses
    wait_for_service('bootkube', ips)
    wait_for_service('tectonic', ips)
    puts 'HOORAY! The cluster is up'
  end

  def wait_for_service(service, ips)
    from = Time.now

    180.times do # 180 * 10 = 1800 seconds = 30 minutes
      return if service_finished_bootstrapping?(ips, service)

      elapsed = Time.now - from
      if (elapsed.round % 5).zero?
        puts "Waiting for bootstrapping of #{service} service to complete..."
        puts "Checked master nodes: #{ips}"
      end
      sleep 10
    end

    puts 'Trying to collecting the logs...'
    forensic(false) # Call forensic to collect logs when service timeout
    raise "timeout waiting for #{service} service to bootstrap on any of: #{ips}"
  end

  def service_finished_bootstrapping?(ips, service)
    command = "test -e /opt/tectonic/init_#{service}.done"

    ips.each do |ip|
      finished = 1
      begin
        _, _, finished = ssh_exec(ip, command)
      rescue => e
        puts "failed to ssh exec on ip #{ip} with: #{e}"
      end

      if finished.zero?
        puts "#{service} service finished successfully on ip #{ip}"
        return true
      end
    end

    false
  end

  def describe_nodes
    KubeCTL.run_and_parse(@kubeconfig, 'get nodes')
  end
end
