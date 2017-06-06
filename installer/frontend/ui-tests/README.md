# Tectonic Installer GUI Tests:

### Set the following environment variables :

`export AWS_ACCESS_KEY_ID = <awsAccessKey> `

`export AWS_SECRET_ACCESS_KEY = <awsSecretKey> `

`export TECTONIC_LICENSE = <coreOSLicense> `

`export PULL_SECRET = <pullSecret> `

### Here is the list of the make targets to build and test the tectonic installer:

`make build` : builds the tectonic installer component

`make launch-installer` : launches the tectonic installer gui

`make launch-installer-guitests` : runs the tectonic installer gui tests

`make gui-tests-cleanup` : removes all test artifacts


#### To execute the tests on browser, `package.json` should include:

```
   "selenium":{
      "start_process":true,
      "host":"127.0.0.1",
      "port":4444,
      "server_path":"./bin/selenium-server-standalone-3.3.1.jar",
      "cli_args":{
         "webdriver.chrome.driver":"./bin/chromedriver"
      }
   },
   "test_settings":{
      "default":{
         "launch_url":"127.0.0.1:8080",
         "selenium_host":"127.0.0.1",
         "selenium_port":4444,
         "desiredCapabilities":{
            "browserName":"chrome",
            "javascriptEnabled":true,
            "acceptSslCerts":true
         }
      }
   }

```

#### To execute the tests in headless mode: [Download Google Chrome Canary](https://www.google.com/chrome/browser/canary.html)  and `package.json` should include:

```
   "selenium":{
      "start_process":false
   },
   "test_settings":{
      "default":{
         "launch_url":"127.0.0.1:8080",
         "selenium_host":"127.0.0.1",
         "selenium_port":9515,
         "desiredCapabilities":{
            "browserName":"chrome",
            "chromeOptions":{
               "args":[
                  "--headless",
                  "--disable-gpu",
                  "--no-sandbox"
               ],
               "binary":"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary"
            }
         }
      }
   }
```
