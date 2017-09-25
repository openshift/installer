# frozen_string_literal: true

# filename: login_page.rb

require_relative 'base_page'

# Login class to deal with tectonic login page
class Login < BasePage
  include RSpec::Matchers

  USERNAME_INPUT          = { id: 'login' }.freeze
  PASSWORD_INPUT          = { id: 'password' }.freeze
  SUBMIT_INPUT            = { css: 'body > div.dex-container > div > form > button' }.freeze
  CLUSTER_STATUS_LABEL    = {
    css: '#content > div > div.co-p-cluster__body > div.row.co-m-nav-title > div > h1 > span'
  }.freeze
  LOGIN_FAIL              = { css: 'body > div.dex-container > div > form > div.dex-error-box' }.freeze
  ADMIN_SIDE_BAR          = { css: '#sidebar > div > div:nth-child(7) > div' }.freeze
  LOGOUT                  = { css: '#sidebar > div > div:nth-child(7) > ul > li:nth-child(2) > a' }.freeze

  def initialize(driver)
    super
  end

  def login_page(url)
    visit url
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
    text_of(LOGIN_FAIL).include? 'Invalid username and password.'
  end
end
