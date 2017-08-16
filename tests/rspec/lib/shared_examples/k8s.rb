require 'smoke_test'
require 'aws_cluster'
require 'forensic'

RSpec.shared_examples 'withCluster' do |tf_vars_path|
  before(:all) do
    prefix = File.basename(tf_vars_path).split('.').first
    raise 'could not extract prefix from tfvars file name' if prefix == ''

    @cluster = AWSCluster.new(prefix, tf_vars_path)
    @cluster.start
  end

  after(:each) do |example|
    Forensic.run(@cluster) if example.exception
  end

  after(:all) do
    @cluster.stop
  end

  it 'succeeds with the golang test suit' do
    expect { SmokeTest.run(@cluster) }.to_not raise_error
  end
end
