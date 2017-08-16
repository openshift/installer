const networkingPageCommands = {
  provideNetworkingDetails() {
    return this
      .waitForElementPresent('@domain', 10000)
      .click('@domain')
      .click('@nextStep');
  },
};

module.exports = {
  url: '',
  commands: [networkingPageCommands],
  elements: {
    domain: {
      selector: 'option[value=Z1ILIMNSJGTMO2]',
    },
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
  },
};
