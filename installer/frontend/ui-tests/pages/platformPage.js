const platformPageCommands = {
  test(platformEl) {
    this.waitForElementVisible('select#platformType', 100000)
      .expect.element(platformEl).to.be.present;
    return this.click(platformEl)
      .expect.element(platformEl).to.be.selected;
  },
};

module.exports = {
  url: () => `${this.api.launchUrl}/define/cluster-type`,
  commands: [platformPageCommands],
  elements: {
    awsPlatform: {
      selector: 'option[value="aws-tf"]',
    },
    bareMetalPlatform: {
      selector: 'option[value="bare-metal-tf"]',
    },
  },
};
