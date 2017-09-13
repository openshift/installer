const networkingPageCommands = {
  testCidrInputs(json) {
    this.setField('#serviceCIDR', json.tectonic_service_cidr);

    this.setField('#podCIDR', '10.2.0.0/15');
    this.expectValidationErrorContains('AWS subnets must be between /16 and /28');
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/16');
    this.expectNoValidationError();
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/21');
    this.expectNoValidationError();
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').text.to.equal('Pod Range Mostly Assigned');

    this.setField('#podCIDR', '10.2.0.0/22');
    this.expectNoValidationError();
    this.expect.element('@alertErrorTitle').text.to.equal('Pod Range Too Small');
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/29');
    this.expectValidationErrorContains('AWS subnets must be between /16 and /28');
    this.expect.element('@alertErrorTitle').text.to.equal('Pod Range Too Small');
    this.expect.element('@alertWarningTitle').to.not.be.present;

    this.setField('#podCIDR', json.tectonic_cluster_cidr);
    this.expectNoValidationError();
    this.expect.element('@alertErrorTitle').to.not.be.present;
    this.expect.element('@alertWarningTitle').to.not.be.present;
  },

  test(json) {
    this.expect.element('@vpcOptionNewPublic').to.be.selected;

    this.selectOption('#awsHostedZoneId option[value=Z1ILIMNSJGTMO2]');
    this.selectOption('#awsSplitDNS option[value=off]');

    this.click('@advanced');
    this.expect.element('#awsVpcCIDR').to.be.visible;
    this.expect.element('[id="awsControllerSubnets.us-west-1a"]').to.be.visible;
    this.expect.element('[id="awsControllerSubnets.us-west-1c"]').to.be.visible;
    this.expect.element('[id="awsWorkerSubnets.us-west-1a"]').to.be.visible;
    this.expect.element('[id="awsWorkerSubnets.us-west-1a"]').to.be.visible;

    this.testCidrInputs(json);

    this.selectOption('@vpcOptionExistingPublic');
    this.testCidrInputs(json);

    this.selectOption('@vpcOptionExistingPrivate');
    this.testCidrInputs(json);

    this.selectOption('@vpcOptionNewPublic');

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
    vpcOptionNewPublic: '.wiz-radio-group:nth-child(1) input[type=radio]',
    vpcOptionExistingPublic: '.wiz-radio-group:nth-child(2) input[type=radio]',
    vpcOptionExistingPrivate: '.wiz-radio-group:nth-child(3) input[type=radio]',
  },
};
