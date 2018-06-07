RSpec.describe "AWS Basic" do
  let(:config) { Config.new }

  # TODO: This is the framework of a test, which needs custom assertions.
  it "test" do
    cluster = config.init
    cluster.install do
      cmd = system("echo SUCCESS")
      expect(cmd).to equal(true)
    end
  end
end
