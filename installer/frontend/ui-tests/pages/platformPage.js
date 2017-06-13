const platformPageCommands = {
  selectPlatform() {
    return this
      .waitForElementVisible('select#platformType', 100000)
      .click('@awsPlatform')
      .click('@nextStep');
  },
};

module.exports = {

  url: () => {
    return this.api.launchUrl + '/define/cluster-type';
  },

  commands: [platformPageCommands],
  elements: {
    awsPlatform: {
      selector: 'option[value="aws-tf"]',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
