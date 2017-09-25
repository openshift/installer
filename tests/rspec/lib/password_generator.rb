# frozen_string_literal: true

require 'jenkins'
require 'securerandom'
require 'bcrypt'

# PasswordGenerator contains helper functions to generate test passwords and also the bcrypted hash
module PasswordGenerator
  def self.generate_password
    SecureRandom.urlsafe_base64
  end

  def self.generate_hash(plain_password)
    BCrypt::Password.create(plain_password)
  end
end
