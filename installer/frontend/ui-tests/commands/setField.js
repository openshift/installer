// Nightwatch's setValue() command actually appends to any existing value, so use this helper to clear any existing
// value first (see github.com/nightwatchjs/nightwatch/issues/4).
exports.command = function(selector, value) {
  this.expect.element(selector).to.be.visible;

  // Hack: Prefix the value with an End key code, followed by lots of backspace characters to clear any existing value
  // first. We do this because Nightwatch's clearValue() function sometimes fails to actually clear the field.
  // Also add a tab character to the end to cause the input to blur after setting its value. This is necessary for
  // testing validation errors.
  this.setValue(selector, `${'\uE010\b'.repeat(50)}${value}\t`);

  this.expect.element(selector).to.have.value.that.equals(value);

  return this;
};
