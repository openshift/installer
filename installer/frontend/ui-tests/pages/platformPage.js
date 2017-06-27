const awsPlatformPageCommands = {
  selectPlatform() {
    return this
      .waitForElementVisible('select#platformType', 100000)
      .click('@awsPlatform')
      .click('@nextStep');
  },
};

const bareMetalplatformPageCommands = {
  selectBareMetalPlatform() {
    return this
      .waitForElementVisible('select#platformType', 100000)
      .click('@bareMetalPlatform')
      .click('@nextStep');
  },
};

module.exports = {

  url: () => {
    return this.api.launchUrl + '/define/cluster-type';
  },

  commands: [awsPlatformPageCommands,bareMetalplatformPageCommands],
  elements: {
    awsPlatform: {
      selector: 'option[value="aws-tf"]',
    },
    bareMetalPlatform: {
      selector: 'option[value="bare-metal-tf"]',
    },
    nextStep: {
      selector:'//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
