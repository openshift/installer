const platformPageCommands = {
  test(platformEl) {
    this.expect.element('select#platformType').to.be.visible.before(60000);
    return this.selectOption(platformEl);
  },
};

module.exports = {
  url: () => `${this.api.launchUrl}/define/cluster-type`,
  commands: [platformPageCommands],
  elements: {
    awsPlatform: 'option[value="aws-tf"]',
    bareMetalPlatform: 'option[value="bare-metal-tf"]',
  },
};
