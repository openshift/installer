require 'aws-sdk-iam'

# frozen_string_literal: true

# The AWSIAM contains helper functions to interact with AWS IAM
module AWSIAM
  def self.assume_role(aws_region)
    role_name = ENV['TECTONIC_INSTALLER_ROLE']
    if role_name.to_s.empty?
      raise 'TECTONIC_INSTALLER_ROLE environment variable not set'
    end

    client = Aws::IAM::Client.new(region: aws_region)

    resp = client.get_role(role_name: role_name)
    role_credentials = request_credentials(aws_region, resp.role.arn) # return the role_credentials
    role_credentials
  end

  def self.request_credentials(aws_region, role_arn)
    sts = Aws::STS::Client.new(region: aws_region)
    role_credentials = Aws::AssumeRoleCredentials.new(
      duration_seconds: 3600,
      client: sts,
      role_arn: role_arn,
      role_session_name: 'temp'
    )
    role_credentials
  rescue StandardError => e
    raise "Error creating Role credentials: #{e.message}"
  end
end
