const chromedriver = require('chromedriver');

module.exports = {
  before (done) {
    chromedriver.start();
    done();
  },

  after (done) {
    chromedriver.stop();
    done();
  },
};
