module.exports = {
  url: '',
  elements: {
    nextStep: {
      selector: '//*[text()[contains(.,"Next Step")]]',
      locateStrategy: 'xpath',
    },
    manuallyBoot: {
      selector: '//a[contains(text(), "Manually boot")]',
      locateStrategy: 'xpath',
    },
  },
};
