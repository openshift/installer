require 'aws-sdk-iam'

# frozen_string_literal: true

# The AWSIAM contains helper functions to interact with AWS IAM
module AWSIAM
  def self.assume_role
    return if ENV.key?('AWS_SESSION_TOKEN')

    role_name = ENV['TECTONIC_INSTALLER_ROLE']
    if role_name.to_s.empty?
      raise 'TECTONIC_INSTALLER_ROLE environment variable not set'
    end

    client = Aws::IAM::Client.new(region: @aws_region)

    resp = client.get_role(role_name: role_name)
    credentials = request_credentials(resp.role.arn)

    export_env_variables(credentials)
  end

  def self.request_credentials(role_arn)
    sts = Aws::STS::Client.new(access_key_id: ENV['AWS_ACCESS_KEY_ID'],
                               secret_access_key: ENV['AWS_SECRET_ACCESS_KEY'],
                               region: @aws_region)
    role_credentials = Aws::AssumeRoleCredentials.new(
      client: sts,
      role_arn: role_arn,
      role_session_name: 'temp'
    )
    role_credentials
  rescue StandardError => e
    raise "Error creating Role credentials: #{e.message}"
  end

  def self.export_env_variables(credential)
    ENV['AWS_ACCESS_KEY_ID'] = credential.credentials.access_key_id
    ENV['AWS_SECRET_ACCESS_KEY'] = credential.credentials.secret_access_key
    ENV['AWS_SESSION_TOKEN'] = credential.credentials.session_token
  end
end
