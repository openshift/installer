# Tectonic Installer GUI Tests

### Set the following environment variables:

```
export AWS_ACCESS_KEY_ID = <awsAccessKey>
export AWS_SECRET_ACCESS_KEY = <awsSecretKey>
export TECTONIC_LICENSE = <coreOSLicense>
export PULL_SECRET = <pullSecret>
```

### Here is the list of the make targets to build and test the tectonic installer:

`make build` : builds the tectonic installer component

`make launch-installer` : launches the tectonic installer gui

`make launch-installer-guitests` : runs the tectonic installer gui tests

`make gui-tests-cleanup` : removes all test artifacts


## Running tests

### Chromedriver

`nightwatch.json` should include:

```json
{
  "selenium": {
    "start_process": false
  },
  "test_settings": {
    "default": {
      "selenium_port": 9515,
      "selenium_host": "localhost",
      "default_path_prefix": "",
      "desiredCapabilities": {
        "browserName": "chrome",
        "chromeOptions": {
          "args": [
            "--disable-gpu",
            "--no-sandbox"
          ]
        }
      },
      "launch_url": "http://localhost:4444"
    }
  }
}
```

To execute the tests in headless mode (no GUI shown), you'll need Google Chrome v60 or later. Add `"--headless"` to chromeOptions's args.

```json
"chromeOptions": {
  "args": [
    "--disable-gpu",
    "--no-sandbox",
    "--headless"
  ]
}
```


### Selenium

Selenium requires Java. To execute the tests on browser using Selenium, `nightwatch.json` should include:

```json
{
  "selenium":{
    "start_process":true,
    "host":"127.0.0.1",
    "port":4445,
    "server_path":"./bin/selenium-server-standalone-3.3.1.jar",
    "cli_args":{
       "webdriver.chrome.driver":"./bin/chromedriver"
    }
  },
  "test_settings":{
    "default":{
      "launch_url": "http://localhost:4444"
      "selenium_host":"127.0.0.1",
      "selenium_port":4445,
      "desiredCapabilities":{
        "browserName":"chrome",
        "javascriptEnabled":true,
        "acceptSslCerts":true
      }
    }
  }
}
```
