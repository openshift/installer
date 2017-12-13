# frozen_string_literal: true

require 'fileutils'

def save_docker_logs(destination_ip, cluster_name, via_host_ip = nil)
  # Get the Container IDs
  command = "docker ps -a --format '{{.ID}} {{.Names}}'"
  stdout, stderr, exitcode = ssh_exec(destination_ip, command, via_host_ip)
  if exitcode != 0
    puts "failed to execut docker ps on #{destination_ip} (exitcode #{exitcode})"
    puts "Standard error: \n#{stderr}"
    return
  end

  containers_id = stdout.split("\n")

  # For each container get the docker logs and save in a file
  containers_id.each do |container|
    id, image_name = container.split(' ')
    command_log = "docker logs #{id}"
    puts "Trying to get the docker logs: #{command_log} - image: #{image_name}"
    stdout, stderr, exitcode = ssh_exec(destination_ip, command_log, via_host_ip)
    output = ''
    output += "Docker Logs of #{image_name} on #{destination_ip} (exitcode #{exitcode})\n"
    output += "Standard output: \n#{stdout}"
    output += "\nStandard error: \n#{stderr}"
    output += "\nEnd of docker logs of #{image_name} container on #{destination_ip}"

    save_to_file(cluster_name, 'docker', destination_ip, image_name, output)
  end
rescue => e
  puts "failed to retrieve docker logs on ip #{destination_ip} with: #{e}"
end

def print_service_logs(destination_ip, service, cluster_name, via_host_ip = nil)
  command = "journalctl --no-pager -u '#{service}'"
  begin
    stdout, stderr, exitcode = ssh_exec(destination_ip, command, via_host_ip)
    output = ''
    output += "Journal of #{service} service on #{destination_ip} (exitcode #{exitcode})\n"
    output += "Standard output: \n#{stdout}"
    output += "\nStandard error: \n#{stderr}"
    output += "\nEnd of journal of #{service} service on #{destination_ip}"
    puts output

    save_to_file(cluster_name, 'systemd', destination_ip, service, output)
  rescue => e
    puts "Cannot retrieve logs of service #{service} - failed to ssh exec on ip #{destination_ip} with: #{e}"
  end
end

def save_kubernetes_events(kubeconfig, cluster_name)
  logs = KubeCTL.run(kubeconfig, 'get ev --all-namespaces -o wide')

  save_to_file(cluster_name, 'kubernetes_events', cluster_name, 'kubernetes_ev_all_namespaces', logs)
end

def save_to_file(cluster_name, service_type, ip, service, output_to_save)
  logs_path = "../../build/#{cluster_name}/logs/#{ip}/#{service_type}"
  FileUtils.mkdir_p(logs_path)
  save_to_file = File.open("#{logs_path}/ip=#{ip},cluster=#{cluster_name},service=#{service}.log", 'w+')
  save_to_file << output_to_save
  save_to_file.close
end
