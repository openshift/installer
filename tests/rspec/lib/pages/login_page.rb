# frozen_string_literal: true

# filename: login_page.rb

require_relative 'base_page'
require 'json'

# Login class to deal with tectonic login page
class Login < BasePage
  include RSpec::Matchers

  USERNAME_INPUT          = { id: 'login' }.freeze
  PASSWORD_INPUT          = { id: 'password' }.freeze
  SUBMIT_INPUT            = { id: 'submit-login' }.freeze
  CLUSTER_STATUS_LABEL    = { id: 'resource-title' }.freeze
  LOGIN_FAIL              = { id: 'login-error' }.freeze
  ADMIN_SIDE_BAR          = { css: '#sidebar > div.navigation-container > div:nth-child(6) > div' }.freeze
  LOGOUT                  = {
    css: '#sidebar > div.navigation-container > div:nth-child(6) > ul > li:nth-child(2) > a'
  }.freeze

  def initialize(driver)
    super
  end

  def login_page(console_url)
    check_console_health(console_url)
    visit(console_url)
    expect(displayed?(USERNAME_INPUT)).to be_truthy
  end

  def logout
    click(ADMIN_SIDE_BAR)
    wait_for { displayed?(LOGOUT) }
    click(LOGOUT)
    expect(displayed?(USERNAME_INPUT)).to be_truthy
  end

  def with(username, password)
    type username, USERNAME_INPUT
    type password, PASSWORD_INPUT
    submit SUBMIT_INPUT
  end

  def success_login?
    wait_for { displayed?(CLUSTER_STATUS_LABEL) }
    text_of(CLUSTER_STATUS_LABEL).include? 'Cluster Status'
  end

  def fail_to_login?
    wait_for { displayed?(LOGIN_FAIL) }
    displayed?(LOGIN_FAIL)
  end

  def check_console_health(console_url)
    from = Time.now
    loop do
      status_json = nil
      begin
        status = `curl -k #{console_url}/health`
        status_json = JSON.parse(status)
        elapsed = Time.now - from
      rescue JSON::ParserError => e
        puts 'Not able to parse the /health result. waiting...'
        sleep 2
        raise "Console was not ready. Not able to get a response from /health. Error #{e}" if elapsed > 1200
        retry
      end
      break if status_json['status'].eql? 'ok'
      puts 'Waiting for Console to be ready...' if (elapsed.round % 5).zero?
      raise "Console was not ready. Response from /health = #{status_json}" if elapsed > 1200 # 20 mins timeout
      sleep 2
    end
  end
end
