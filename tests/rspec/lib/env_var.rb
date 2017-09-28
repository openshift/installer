# frozen_string_literal: true

# EnvVars contains helper functions for system environment variables
module EnvVar
  def self.set?(vars)
    vars.all? { |var| ENV.key?(var) }
  end

  def self.contains_any?(vars)
    vars.any? { |var| ENV.key?(var) }
  end
end
