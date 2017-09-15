const wizard = require('../utils/wizard');

exports.command = function(text) {
  this.expect.element('.wiz-error-message:not(.hidden)').text.contains(text);

  // Wizard next button should be disabled when there is a validation error
  this.expect.element(wizard.nextStep).to.have.attribute('class').which.contains('disabled');

  return this;
};
