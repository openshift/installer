// Nightwatch's setValue() command actually appends to any existing value, so use this helper to clear any existing
// value first (see github.com/nightwatchjs/nightwatch/issues/4).
exports.command = function(selector, value) {
  this.expect.element(selector).to.be.visible;
  this.clearValue(selector)
    .setValue(selector, value)
    .expect.element(selector).to.have.value.that.equals(value);
  return this;
};
