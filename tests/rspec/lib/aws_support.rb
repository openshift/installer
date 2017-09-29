# frozen_string_literal: true

require 'aws-sdk'

# Shared support code for AWS-based operations
#
module AwsSupport
  def self.sorted_auto_scaling_instances(aws_autoscaling_group_id, aws_region)
    aws = Aws::AutoScaling::Client.new(region: aws_region)
    resp = aws.describe_auto_scaling_groups(auto_scaling_group_names: [
                                              aws_autoscaling_group_id.to_s
                                            ])
    resp.auto_scaling_groups[0].instances.map(&:instance_id).sort
  end

  def self.preferred_instance_ip_address(instance_id, aws_region)
    aws = Aws::EC2::Client.new(region: aws_region)
    resp = aws.describe_instances(instance_ids: [instance_id.to_s])
    ssh_master_ip = if resp.reservations[0].instances[0].network_interfaces[0].association.nil?
                      resp.reservations[0].instances[0].network_interfaces[0].private_ip_address
                    else
                      resp.reservations[0].instances[0].network_interfaces[0].association.public_ip
                    end
    ssh_master_ip
  end

  def self.create_aws_key_pairs(aws_region)
    client = Aws::EC2::Client.new(region: aws_region)
    key_pair_name = NameGenerator.generate_short_name

    dir_home = `echo ${HOME}`.chomp
    ssh_pub_key = File.read("#{dir_home}/.ssh/id_rsa.pub")

    resp = client.import_key_pair(dry_run: false,
                                  key_name: key_pair_name,
                                  public_key_material: ssh_pub_key)
    raise 'Error to import AWS key pair' unless resp.successful?
    resp.key_name
  end

  def self.delete_aws_key_pairs(key_pair_name, aws_region)
    client = Aws::EC2::Client.new(region: aws_region)

    resp = client.delete_key_pair(key_name: key_pair_name,
                                  dry_run: false)
    raise "Error to delete AWS key pair - #{key_pair_name}" unless resp.successful?
  end
end
