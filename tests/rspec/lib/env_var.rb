# EnvVars contains helper functions for system environment variables
module EnvVar
  def self.set?(vars)
    vars.all? { |var| ENV.key?(var) }
  end
end
