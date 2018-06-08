require "tmpdir"

Dir[File.expand_path("support/**/*.rb", __dir__)].each { |f| require f }

RSpec.configure do |config| # rubocop:disable Metrics/BlockLength
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups

  config.filter_run_when_matching :focus

  config.example_status_persistence_file_path = "spec/examples.txt"

  config.disable_monkey_patching!

  config.warnings = true

  config.default_formatter = "doc" if config.files_to_run.one?

  config.profile_examples = 10

  config.order = :random

  Kernel.srand config.seed

  # rubocop:disable Style/GlobalVars
  config.before(:suite) do
    tarball = Bazel.build(:tarball)
    $tmp_dir = Dir.mktmpdir
    tarball.untar($tmp_dir)
  end

  config.around(:example) do |example|
    Dir.chdir(File.join($tmp_dir, "tectonic-dev")) { example.run }
  end

  config.after(:suite) do
    FileUtils.remove_entry_secure($tmp_dir)
  end
  # rubocop:enable Style/GlobalVars
end
