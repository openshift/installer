const networkingPageCommands = {
  test(json) {
    this.click('@advanced');
    this.expect.element('#awsVpcCIDR').to.be.visible;
    this.expect.element('[id="awsControllerSubnets.us-west-1a"]').to.be.visible;
    this.expect.element('[id="awsControllerSubnets.us-west-1c"]').to.be.visible;
    this.expect.element('[id="awsWorkerSubnets.us-west-1a"]').to.be.visible;
    this.expect.element('[id="awsWorkerSubnets.us-west-1a"]').to.be.visible;
    this
      .selectOption('#awsHostedZoneId option[value=Z1ILIMNSJGTMO2]')
      .selectOption('#awsSplitDNS option[value=off]')
      .setField('#serviceCIDR', json.tectonic_service_cidr);

    this.setField('#podCIDR', '10.2.0.0/15');
    this.expect.element('@validationError').text.to.equal('AWS subnets must be between /16 and /28.');
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/16');
    this.expect.element('@validationError').to.not.be.present;
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/21');
    this.expect.element('@validationError').to.not.be.present;
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').text.to.equal('Pod Range Mostly Assigned');

    this.setField('#podCIDR', '10.2.0.0/22');
    this.expect.element('@validationError').to.not.be.present;
    this.expect.element('@alertErrorTitle').text.to.equal('Pod Range Too Small');
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/29');
    this.expect.element('@validationError').text.to.equal('AWS subnets must be between /16 and /28.');
    this.expect.element('@alertErrorTitle').text.to.equal('Pod Range Too Small');
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', json.tectonic_cluster_cidr);
    this.expect.element('@validationError').to.not.be.present;
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;

    return this;
  },
};

module.exports = {
  commands: [networkingPageCommands],
  elements: {
    advanced: {
      selector: '//*[text()[contains(.,"Advanced Settings")]]',
      locateStrategy: 'xpath',
    },
    alertWarningTitle: '.alert-info b',
    alertErrorTitle: '.alert-error b',
    validationError: '.wiz-error-message:not(.hidden)',
  },
};
