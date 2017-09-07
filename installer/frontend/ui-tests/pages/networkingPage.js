const networkingPageCommands = {
  test(json) {
    return this
      .selectOption('#awsHostedZoneId option[value=Z1ILIMNSJGTMO2]')
      .selectOption('#awsSplitDNS option[value=off]')
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
  },
};
