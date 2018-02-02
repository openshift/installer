# frozen_string_literal: true

require 'aws-sdk-autoscaling'
require 'aws-sdk-ec2'
require 'aws_iam'
require 'with_retries'

# Shared support code for AWS-based operations
#
module AwsSupport
  def self.sorted_auto_scaling_instances(aws_autoscaling_group_id, aws_region, role_credentials = nil)
    options = define_options(aws_region, role_credentials)

    aws = Aws::AutoScaling::Client.new(options)
    resp = ''
    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 5, sleep: 10) do
      resp = aws.describe_auto_scaling_groups(auto_scaling_group_names: [
                                                aws_autoscaling_group_id.to_s
                                              ])
    end
    resp.auto_scaling_groups[0].instances.map(&:instance_id).sort
  end

  def self.instance_ip_address(instance_id, aws_region, role_credentials = nil)
    options = define_options(aws_region, role_credentials)
    aws = Aws::EC2::Client.new(options)
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

  def self.create_aws_key_pairs(aws_region, role_credentials = nil)
    options = define_options(aws_region, role_credentials)

    client = Aws::EC2::Client.new(options)
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

  def self.delete_aws_key_pairs(key_pair_name, aws_region, role_credentials = nil)
    options = define_options(aws_region, role_credentials)
    client = Aws::EC2::Client.new(options)

    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 3, sleep: 10) do
      client.delete_key_pair(key_name: key_pair_name,
                             dry_run: false)
    end
  end

  def self.collect_ec2_console_logs(instance_id, aws_region, role_credentials = nil)
    options = define_options(aws_region, role_credentials)
    client = Aws::EC2::Client.new(options)

    resp = ''
    Retriable.with_retries(Aws::EC2::Errors::RequestLimitExceeded, limit: 3, sleep: 10) do
      resp = client.get_console_output(instance_id: instance_id.to_s,
                                       dry_run: false)
    end
    resp.output
  end

  def self.check_expiration_role(aws_region, role_credentials)
    if Time.now.utc > role_credentials.expiration
      puts 'Need to renew the AWS tokens. Assumed role expired.'
      role_credentials = AWSIAM.assume_role(aws_region)
    end
    role_credentials
  end

  def self.define_options(aws_region, role_credentials)
    options = { region: aws_region }
    unless role_credentials.nil?
      role_credentials = check_expiration_role(aws_region, role_credentials)
      options[:credentials] = role_credentials
    end
    options
  end
end
