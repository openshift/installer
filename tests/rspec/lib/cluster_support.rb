# frozen_string_literal: true

require 'fileutils'

def save_docker_logs(ip, cluster_name)
  # Get the Container IDs
  command = "docker ps -a --format '{{.ID}} {{.Names}}'"
  stdout, stderr, exitcode = ssh_exec(ip, command)
  if exitcode != 0
    puts "failed to execut docker ps on #{ip} (exitcode #{exitcode})"
    puts1 "Standard error: \n#{stderr}"
    return
  end

  containers_id = stdout.split("\n")

  # For each container get the docker logs and save in a file
  containers_id.each do |container|
    id, image_name = container.split(' ')
    command_log = "docker logs #{id}"
    puts "Trying to get the docker logs: #{command_log} - image: #{image_name}"
    stdout, stderr, exitcode = ssh_exec(ip, command_log)
    output = ''
    output << "Docker Logs of #{image_name} on #{ip} (exitcode #{exitcode})\n"
    output << "Standard output: \n#{stdout}"
    output << "\nStandard error: \n#{stderr}"
    output << "\nEnd of docker logs of #{image_name} container on #{ip}"

    save_to_file(cluster_name, 'docker', ip, image_name, output)
  end
rescue => e
  puts "failed to retrieve docker logs on ip #{ip} with: #{e}"
end

def print_service_logs(ip, service, cluster_name)
  command = "journalctl --no-pager -u '#{service}'"
  begin
    stdout, stderr, exitcode = ssh_exec(ip, command)
    output = ''
    output << "Journal of #{service} service on #{ip} (exitcode #{exitcode})\n"
    output << "Standard output: \n#{stdout}"
    output << "\nStandard error: \n#{stderr}"
    output << "\nEnd of journal of #{service} service on #{ip}"
    puts output

    save_to_file(cluster_name, 'journal', ip, service, output)
  rescue => e
    puts "Cannot retrieve logs of service #{service} - failed to ssh exec on ip #{ip} with: #{e}"
  end
end

def save_to_file(cluster_name, service_type, ip, service, output_to_save)
  logs_path = "../../build/#{cluster_name}/logs/#{service_type}_logs_#{cluster_name}"
  FileUtils.mkdir_p(service_logs_path)
  save_to_file = File.open("#{logs_path}/#{ip}_#{service}.log", 'w+')
  save_to_file << output_to_save
  save_to_file.close
end
