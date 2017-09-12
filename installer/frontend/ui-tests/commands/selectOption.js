// Click a select option
exports.command = function(selector) {
  this.expect.element(selector).to.be.visible;
  this.click(selector);
  this.expect.element(selector).to.be.selected;
  return this;
};
