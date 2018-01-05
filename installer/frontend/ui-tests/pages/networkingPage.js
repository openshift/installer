const networkConfigurationPage = require('./networkConfigurationPage');

const pageCommands = {
  testCidrInputs(json) {
    this.setField('#serviceCIDR', json.serviceCIDR);

    this.setField('#podCIDR', '10.2.0.0/15');
    this.expectValidationErrorContains('AWS subnets must be between /16 and /28');
    this.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/16');
    this.expectNoValidationError();
    this.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/22');
    this.expectNoValidationError();
    this.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    this.expect.element('@k8sCIDRsWarningTitle').text.to.equal('Pod range mostly assigned');

    this.setField('#podCIDR', '10.2.0.0/23');
    this.expectNoValidationError();
    this.expect.element('@k8sCIDRsErrorTitle').text.to.equal('Pod range too small');
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/29');
    this.expectValidationErrorContains('AWS subnets must be between /16 and /28');
    this.expect.element('@k8sCIDRsErrorTitle').text.to.equal('Pod range too small');
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;

    networkConfigurationPage.testDockerBridgeValidation(this);

    this.setField('#podCIDR', json.podCIDR);
    this.expectNoValidationError();
    this.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;
  },

  test(json) {
    this.expect.element('@vpcOptionNewPublic').to.be.selected;
    this.expect.element('#awsVpcId').to.not.be.present;

    this.selectOption(`#awsHostedZoneId option[value=${json.awsHostedZoneId}]`);
    this.selectOption('#awsSplitDNS option[value=off]');

    // If a AWS VPC CIDR is specified, do a full test of the advanced networking options
    const isAdvanced = !!json.awsVpcCIDR || !!json.awsControllerSubnets || !!json.awsWorkerSubnets;

    if (isAdvanced) {
      this.click('@advanced');
      this.expect.element('#awsVpcCIDR').to.be.visible;
      this.expect.element('@masterSubnet1a').to.be.visible;
      this.expect.element('@masterSubnet1c').to.be.visible;
      this.expect.element('@workerSubnet1a').to.be.visible;
      this.expect.element('@workerSubnet1c').to.be.visible;

      // Subnet CIDR outside of VPC CIDR
      this.setField('#awsVpcCIDR', '0.0.0.0/20');
      this.expect.element('@alertError').text.to.contain('vpcCIDR (0.0.0.0/20) does not contain instanceCIDR');
      this.setField('#awsVpcCIDR', '10.0.0.0/16');
      this.expect.element('@alertError').to.not.be.present;

      // Overlapping public subnets
      this.setField('@masterSubnet1c', '10.0.0.0/19');
      this.expect.element('@alertError').text.to.contain('CIDR of subnet 1 (10.0.0.0/19) overlaps with CIDR of subnet 0');
      this.setField('@masterSubnet1c', '10.0.32.0/20');
      this.expect.element('@alertError').to.not.be.present;

      // Overlapping private subnets
      this.setField('@workerSubnet1c', '10.0.64.0/19');
      this.expect.element('@alertError').text.to.contain('CIDR of subnet 1 (10.0.64.0/19) overlaps with CIDR of subnet 0');
      this.setField('@workerSubnet1c', '10.0.96.0/20');
      this.expect.element('@alertError').to.not.be.present;

      // Pod CIDR overlaps with AWS VPC CIDR
      this.setField('#podCIDR', '10.0.0.0/16');
      this.expect.element('@alertError').text.to.contain('IP Network "10.0.0.0/16" has been given twice');

      // Pod CIDR overlaps with Service CIDR
      this.setField('#podCIDR', '10.3.0.0/16');
      this.expect.element('@alertError').text.to.contain('IP Network "10.3.0.0/16" has been given twice');

      this.setField('#podCIDR', '10.2.0.0/16');
      this.expect.element('@alertError').to.not.be.present;

      this.testCidrInputs(json);
    }

    this.selectOption('@vpcOptionExistingPublic');
    this.expect.element('#awsVpcId').to.be.visible;
    if (isAdvanced) {
      this.testCidrInputs(json);
    }

    this.selectOption('@vpcOptionExistingPrivate');
    this.expect.element('#awsVpcId').to.be.visible;
    if (isAdvanced) {
      this.testCidrInputs(json);
    }

    this.selectOption('@vpcOptionNewPublic');
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    advanced: {
      selector: '//*[text()[contains(.,"Advanced Settings")]]',
      locateStrategy: 'xpath',
    },
    alertError: '.alert-error',
    k8sCIDRsErrorTitle: '#k8sCIDRs .alert-error b',
    k8sCIDRsWarningTitle: '#k8sCIDRs .alert-info b',
    masterSubnet1a: '[id="awsControllerSubnets.us-west-1a"]',
    masterSubnet1c: '[id="awsControllerSubnets.us-west-1c"]',
    vpcOptionExistingPrivate: '.wiz-radio-group:nth-child(3) input[type=radio]',
    vpcOptionExistingPublic: '.wiz-radio-group:nth-child(2) input[type=radio]',
    vpcOptionNewPublic: '.wiz-radio-group:nth-child(1) input[type=radio]',
    workerSubnet1a: '[id="awsWorkerSubnets.us-west-1a"]',
    workerSubnet1c: '[id="awsWorkerSubnets.us-west-1c"]',
  },
};
