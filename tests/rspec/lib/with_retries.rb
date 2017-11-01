# frozen_string_literal: true

# inspiration: https://gist.github.com/suciuvlad/4078129
# This will catch any exception and retry twice (three tries total):
#   with_retries { ... }
#
# This will catch any exception and retry four times (five tries total):
#   with_retries(:limit => 5) { ... }
#
# This will catch a specific exception and retry once (two tries total):
#   with_retries(Some::Error, :limit => 2) { ... }
#
# You can also sleep in between tries. This is helpful if you're hoping
# that some external service recovers from its issues.
#   with_retries(Service::Error, :sleep => 1) { ... }
#
module Retriable
  def self.with_retries(*args, **keyword_args)
    exceptions = args
    options = {}

    options[:limit] =  keyword_args[:limit] ||= 3
    options[:sleep] =  keyword_args[:sleep] ||= 0
    exceptions = [Exception] if exceptions.empty?

    retried = 0
    begin
      yield
    rescue *exceptions => e
      raise e if retried + 1 > options[:limit]
      retried += 1
      sleep_time = options[:sleep] * retried
      sleep sleep_time
      puts "Error #{e}; retrying in #{sleep_time} seconds"
      retry
    end
  end
end
