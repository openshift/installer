const fs = require('fs');

exports.command = function (selector, text) {
  this.expect.element(selector).to.be.present;

  // eslint-disable-next-line no-sync
  fs.writeFileSync(this.globals.tmpUploadPath, text);
  return this.setValue(selector, this.globals.tmpUploadPath);
};
