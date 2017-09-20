# frozen_string_literal: true

# Shared support code for Azure-based operations
#
module AzureSupport
  LOCATIONS = %w[eastus westus northcentralus southcentralus].freeze

  def self.random_location_unless_defined
    ENV['TF_VAR_tectonic_azure_location'] || LOCATIONS.sample
  end
end
