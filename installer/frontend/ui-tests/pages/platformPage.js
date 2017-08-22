const platformPageCommands = {
  selectPlatform(platformEl) {
    this.waitForElementVisible('select#platformType', 100000)
      .expect.element(platformEl).to.be.present;
    this.click(platformEl)
      .expect.element(platformEl).to.be.selected;
    return this.click('@nextStep');
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
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
