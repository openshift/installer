const chromedriver = require('chromedriver');

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
};
