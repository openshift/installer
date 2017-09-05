# frozen_string_literal: true

require 'env_var'

# Jenkins contains helper functions to interact with Jenkins
module Jenkins
  def self.environment?
    EnvVar.set?(%w[JENKINS_HOME])
  end
end
