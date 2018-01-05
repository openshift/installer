# frozen_string_literal: true

require 'aws-sdk-autoscaling'
require 'aws-sdk-ec2'
require 'with_retries'

# Shared support code for AWS-based operations
#
module AwsSupport
  def self.sorted_auto_scaling_instances(aws_autoscaling_group_id, aws_region)
    aws = Aws::AutoScaling::Client.new(region: aws_region)
    resp = ''
    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 5, sleep: 10) do
      resp = aws.describe_auto_scaling_groups(auto_scaling_group_names: [
                                                aws_autoscaling_group_id.to_s
                                              ])
    end
    resp.auto_scaling_groups[0].instances.map(&:instance_id).sort
  end

  def self.preferred_instance_ip_address(instance_id, aws_region)
    aws = Aws::EC2::Client.new(region: aws_region)
    resp = ''
    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 5, sleep: 10) do
      resp = aws.describe_instances(instance_ids: [instance_id.to_s])
    end
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

    resp = ''
    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 5, sleep: 10) do
      resp = client.import_key_pair(dry_run: false,
                                    key_name: key_pair_name,
                                    public_key_material: ssh_pub_key)
    end
    resp.key_name
  end

  def self.delete_aws_key_pairs(key_pair_name, aws_region)
    client = Aws::EC2::Client.new(region: aws_region)

    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 3, sleep: 10) do
      client.delete_key_pair(key_name: key_pair_name,
                             dry_run: false)
    end
  end
end
