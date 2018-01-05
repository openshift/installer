# frozen_string_literal: true

require 'selenium-webdriver'
require 'headless'

# WebdriverHelpers contains helper functions to create the browser to use in the UI tests.
module WebdriverHelpers
  def self.start_webdriver
    @headless = Headless.new
    @headless.start
    options = Selenium::WebDriver::Chrome::Options.new
    options.add_argument('--ignore-certificate-errors')
    options.add_argument('--disable-popup-blocking')
    options.add_argument('--disable-translate')
    options.add_argument('--disable-gpu')
    options.add_argument('--no-sandbox')
    # options.add_argument('--headless')
    driver = Selenium::WebDriver.for :chrome, options: options
    target_size = Selenium::WebDriver::Dimension.new(1920, 1080)
    driver.manage.window.size = target_size
    driver.manage.timeouts.implicit_wait = 20 # seconds
    driver.manage.timeouts.page_load = 60 # seconds
    driver
  end

  def self.stop_webdriver(driver)
    driver.quit
    @headless.destroy
  end
end
