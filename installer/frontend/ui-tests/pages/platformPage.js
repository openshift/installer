const platformPageCommands = {
  test(platformEl) {
    this.expect.element('select#platformType').to.be.visible.before(60000);
    this.selectOption(platformEl);
  },
};

module.exports = {
  commands: [platformPageCommands],
  elements: {
    awsPlatform: 'option[value="aws-tf"]',
    bareMetalPlatform: 'option[value="bare-metal-tf"]',
  },
};
