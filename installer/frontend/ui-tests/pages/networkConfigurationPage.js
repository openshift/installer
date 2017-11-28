const testDockerBridgeValidation = page => {
  page.setField('#podCIDR', '172.17.0.0/16');
  page.expectValidationErrorContains('Overlaps with default Docker Bridge subnet (172.17.0.0/16)');
  page.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
  page.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;

  ['172.18.0.0', '172.31.255.255', '192.168.0.0', '192.168.15.255'].forEach(ip => {
    page.setField('#podCIDR', `${ip}/16`);
    page.expectNoValidationError();
    page.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    page.expect.element('@k8sCIDRsWarningTitle').text.to.equal('Pod range may conflict with Docker Bridge subnet');
  });

  ['172.32.0.0', '192.168.16.0'].forEach(ip => {
    page.setField('#podCIDR', `${ip}/16`);
    page.expectNoValidationError();
    page.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    page.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;
  });
};

const pageCommands = {
  test(json) {
    this.setField('#serviceCIDR', json.serviceCIDR);

    this.setField('#podCIDR', '10.2.0.0/21');
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

    testDockerBridgeValidation(this);

    this.setField('#podCIDR', json.podCIDR);
    this.expectNoValidationError();
    this.expect.element('@k8sCIDRsErrorTitle').to.not.be.present;
    this.expect.element('@k8sCIDRsWarningTitle').to.not.be.present;
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    k8sCIDRsErrorTitle: '#k8sCIDRs .alert-error b',
    k8sCIDRsWarningTitle: '#k8sCIDRs .alert-info b',
  },
  testDockerBridgeValidation,
};
