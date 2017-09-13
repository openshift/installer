exports.command = function() {
  this.expect.element('.wiz-error-message:not(.hidden)').to.not.be.present;
  return this;
};
