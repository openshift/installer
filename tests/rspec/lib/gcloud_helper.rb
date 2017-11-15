# frozen_string_literal: true

# GcloudHelper contains helper functions to interact with the command line tool
# gcloud
class GcloudHelper
  attr_accessor :authenticated

  def initialize
    @authenticated = true
    auth_command = 'gcloud auth activate-service-account --key-file="${GOOGLE_APPLICATION_CREDENTIALS}"'
    project_command = 'gcloud config set project "${GOOGLE_PROJECT}"'
    auth = system(auth_command) && system(project_command)
    raise 'Problem login with gcloud' unless auth
  end

  def run(args)
    if @authenticated
      out = `gcloud #{args}`
      raise GcloudCmdError if $CHILD_STATUS.exitstatus != 0
      return out
    end
    raise 'You need to be authenticated for running gcloud'
  end

  # Gcloud is raised whenever the shell command 'gcloud' fails
  class GcloudCmdError < StandardError
    def initialize(msg = 'failed to call gcloud')
      super
    end
  end
end
