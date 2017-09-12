const networkConfigurationPageCommands = {
  enterCIDRs(json) {
    this
      .setField('#podCIDR', json.tectonic_cluster_cidr)
      .setField('#serviceCIDR', json.tectonic_service_cidr);
  },
};

module.exports = {
  commands: [networkConfigurationPageCommands],
  elements: {},
};
