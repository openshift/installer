require 'json'
require 'English'

# KubeCTL contains helper functions to interact with the command line tool
# kubectl
module KubeCTL
  def self.run(kubeconfig, args)
    out = `KUBECONFIG=#{kubeconfig} kubectl #{args}`
    raise KubeCTLCmdError if $CHILD_STATUS.exitstatus != 0
    out
  end

  def self.run_and_parse(kubeconfig, args)
    out = run(kubeconfig, args + ' -ojson')
    JSON.parse(out)
  end

  # KubeCTLCmdError is raised whenever the shell command 'kubectl' fails
  class KubeCTLCmdError < StandardError
    def initialize(msg = 'failed to call kubectl')
      super
    end
  end
end
