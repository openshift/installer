exports.command = function(text) {
  this.expect.element('.wiz-error-message:not(.hidden)').text.contains(text);
  return this;
};
