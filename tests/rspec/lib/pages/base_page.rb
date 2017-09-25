# frozen_string_literal: true

# filename: base_page.rb
require 'selenium-webdriver'

# Base class to provide basic element interaction
class BasePage
  def initialize(driver)
    @driver = driver
  end

  def visit(url_path)
    @driver.get url_path
  end

  def current_url
    @driver.current_url
  end

  def refresh
    @driver.navigate.refresh
  end

  def find(locator)
    retries ||= 0
    wait_for { @driver.find_element locator }
  rescue
    retry if (retries += 1) < 3
  end

  def finds(locator)
    @driver.find_elements locator
  end

  def type(text, locator)
    find(locator).send_keys text
  end

  def clear_text(locator)
    find(locator).clear
  end

  def submit(locator)
    find(locator).submit
  end

  def click(locator)
    retries ||= 0
    find(locator).click
  rescue
    retry if (retries += 1) < 3
  end

  def text_of(locator)
    find(locator).text
  end

  def text_include(locator, message)
    retries ||= 0
    find(locator).text.include?(message)
  rescue
    retry if (retries += 1) < 3
  end

  def get_attibute(locator, attribute)
    find(locator).attribute(attribute)
  end

  def scroll_to(locator)
    retries ||= 0
    find(locator).location_once_scrolled_into_view
  rescue
    retry if (retries += 1) < 10
  end

  def wait_for(seconds = 15)
    Selenium::WebDriver::Wait.new(timeout: seconds).until { yield }
  end

  def wait_for_message(locator, message, seconds = 15)
    Selenium::WebDriver::Wait.new(timeout: seconds).until do
      text_include(locator, message)
    end
  end

  def displayed?(locator)
    retries ||= 0
    find(locator).displayed?
  rescue Selenium::WebDriver::Error::NoSuchElementError
    false
  rescue
    retry if (retries += 1) < 10
  end

  def enabled?(locator)
    retries ||= 0
    find(locator).enabled?
  rescue
    retry if (retries += 1) < 10
  end
end
