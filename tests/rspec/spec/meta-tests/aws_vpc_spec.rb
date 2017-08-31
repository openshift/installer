require 'aws_vpc'
require 'aws_iam'

describe AWSVPC do
  before(:all) do
    AWSIAM.assume_role if Jenkins.environment?
    @vpc = described_class.new('test-vpc')
  end

  it '#initialize generated a password' do
    expect(@vpc.ovpn_password.nil?).to be_falsey
    expect(@vpc.ovpn_password.empty?).to be_falsey
  end

  context '#create' do
    before(:all) do
      @vpc.create
    end

    after(:all) do
      @vpc.destroy
    end

    it 'sets the vpn_url' do
      expect(@vpc.vpn_url.nil?).to be_falsey
      expect(@vpc.vpn_url.empty?).to be_falsey
    end

    it 'sets the vpn_config' do
      expect(@vpc.vpn_connection.vpn_conf.nil?).to be_falsey
      expect(@vpc.vpn_connection.vpn_conf.empty?).to be_falsey
    end
  end
end
