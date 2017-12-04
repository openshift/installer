# frozen_string_literal: true

require 'jenkins'
require 'securerandom'

# PasswordGenerator contains helper functions to generate test passwords and also the bcrypted hash
module PasswordGenerator
  def self.generate_password
    SecureRandom.urlsafe_base64
  end
end
