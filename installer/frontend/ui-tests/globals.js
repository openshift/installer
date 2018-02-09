const chromedriver = require('chromedriver');
const path = require('path');

module.exports = {
  default: {
    waitForConditionTimeout: 10000,
  },

  before (done) {
    chromedriver.start();
    done();
  },

  after (done) {
    chromedriver.stop();
    done();
  },

  tmpUploadPath: path.join(path.resolve(__dirname), 'tmp-upload-test'),
};
