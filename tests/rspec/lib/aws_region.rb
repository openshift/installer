def export_random_region_if_not_defined
  regions = %w[us-east-1 us-east-2 us-west-1 us-west-2]

  return if ENV.key?('TF_VAR_tectonic_aws_region')

  ENV['TF_VAR_tectonic_aws_region'] = regions.sample
end
