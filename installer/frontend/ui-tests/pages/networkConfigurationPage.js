const networkConfigurationPageCommands = {
  enterCIDRs(json) {
    return this
      .setField('#podCIDR', json.tectonic_cluster_cidr, true)
      .setField('#serviceCIDR', json.tectonic_service_cidr, true);
  },
};

module.exports = {
  commands: [networkConfigurationPageCommands],
  elements: {},
};
