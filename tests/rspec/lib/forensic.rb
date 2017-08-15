require 'ssh'

# Forensic contains helper functions to run the cluster forensic scripts
module Forensic
  def self.run(cluster)
    return if cluster.class.name != 'AWSCluster'
    check_prerequisites

    env_variables = cluster.env_variables
    env_variables['AWS_REGION'] = cluster.tfvars_file.region

    succeeded = system(
      env_variables,
      './../smoke/aws/cluster-foreach.sh ./../smoke/forensics.sh'
    )
    raise 'Forensic script failed' unless succeeded
  end
end
