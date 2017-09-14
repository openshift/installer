# frozen_string_literal: true

# Shared support code for Azure-based operations
#
module AzureSupport
  LOCATIONS = %w[eastus eastus2 centralus westus westus2].freeze

  def self.random_location_unless_defined
    ENV['TF_VAR_tectonic_azure_location'] || LOCATIONS.sample
  end
end
