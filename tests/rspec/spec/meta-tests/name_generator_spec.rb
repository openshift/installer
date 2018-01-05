# frozen_string_literal: true

require 'name_generator'

describe NameGenerator do
  before(:all) do
    @curent_env = ENV.clone
  end

  after(:all) do
    ENV.clear
    ENV = @curent_env.clone
  end

  it 'prefix should not contain more than 10 chars' do
    ENV['BUILD_ID'] = '1'
    ENV['BRANCH_NAME'] = 'master'
    prefix = 'my-long-string-with-prefix'

    expect(NameGenerator.jenkins_env_name(prefix).split('-')[0]).to eq('mylonstrwi')
  end
end
