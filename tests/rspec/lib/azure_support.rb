# frozen_string_literal: true

require 'azure/storage/common'
require 'azure/storage/blob'

# Shared support code for Azure-based operations
#
module AzureSupport
  LOCATIONS = %w[canadacentral westcentralus westeurope southeastasia japaneast].freeze

  def self.random_location_unless_defined
    ENV['TF_VAR_tectonic_azure_location'] || LOCATIONS.sample
  end

  def self.set_ssh_key_path
    dir_home = ENV['HOME']
    ssh_pub_key_path = "#{dir_home}/.ssh/id_rsa.pub"
    ENV['TF_VAR_tectonic_azure_ssh_key'] = ssh_pub_key_path
    ssh_pub_key_path
  end

  def self.collect_azure_vm_console_logs(storage_name, api_key)
    output = {}
    common_client = Azure::Storage::Common::Client.create(storage_account_name: storage_name,
                                                          storage_access_key: api_key)

    blob_client = Azure::Storage::Blob::BlobService.new(client: common_client)

    blob_containers = blob_client.list_containers

    blob_containers.each do |blob_container|
      blobs = blob_client.list_blobs(blob_container.name)
      blobs.each do |blob|
        if blob.name.end_with?('.log')
          _, content = blob_client.get_blob(blob_container.name, blob.name)
          output.merge!(blob_container.name => content)
        end
      end
    end
    output
  end
end
