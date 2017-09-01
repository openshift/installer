const platformPageCommands = {
  test(platformEl) {
    this.expect.element('select#platformType').to.be.visible.before(60000);
    return this.click(platformEl)
      .expect.element(platformEl).to.be.selected;
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
