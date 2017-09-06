const networkingPageCommands = {
  test(json) {
    return this
      .waitForElementPresent('@domain')
      .click('@domain')
      .click('@advanced')
      .setField('#podCIDR', json.tectonic_cluster_cidr)
      .setField('#serviceCIDR', json.tectonic_service_cidr);
  },
};

module.exports = {
  commands: [networkingPageCommands],
  elements: {
    advanced: {
      selector: '//*[text()[contains(.,"Advanced Settings")]]',
      locateStrategy: 'xpath',
    },
    domain: 'option[value=Z1ILIMNSJGTMO2]',
  },
};
