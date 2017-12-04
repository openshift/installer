# frozen_string_literal: true

require 'tfvars_file'

TFVARS_FILE_PATH = '../smoke/aws/vars/aws.tfvars.json'

describe TFVarsFile do
  subject { described_class.new(TFVARS_FILE_PATH) }

  context 'with wrong file path' do
    subject { described_class.new('wrong-file-path') }
    it 'throws an exception' do
      expect { subject }.to raise_error(RuntimeError)
    end
  end

  it '#path returns the file path' do
    expect(subject.path).to be(TFVARS_FILE_PATH)
  end

  it '#node_count returns correct #' do
    expect(subject.node_count).to eq(4)
  end

  it '#networking? returns empty string if not set' do
    expect(subject.networking).to eq('')
  end

  it 'returns raw variable values via reflection' do
    expect(subject.tectonic_aws_worker_ec2_type).to eq('m4.large')
  end
end
