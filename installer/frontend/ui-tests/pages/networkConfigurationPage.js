const pageCommands = {
  test(json) {
    this.setField('#serviceCIDR', json.serviceCIDR);

    this.setField('#podCIDR', '10.2.0.0/21');
    this.expectNoValidationError();
    this.expect.element('@alertError').to.not.be.present;
    this.expect.element('@alertWarning').to.not.be.present;

    this.setField('#podCIDR', '10.2.0.0/22');
    this.expectNoValidationError();
    this.expect.element('@alertError').to.not.be.present;
    this.expect.element('@alertWarning').text.to.contain('Pod range mostly assigned');

    this.setField('#podCIDR', '10.2.0.0/23');
    this.expectNoValidationError();
    this.expect.element('@alertError').text.to.contain('Pod range too small');
    this.expect.element('@alertWarning').to.not.be.present;

    this.setField('#podCIDR', json.podCIDR);
    this.expectNoValidationError();
    this.expect.element('@alertError').to.not.be.present;
    this.expect.element('@alertWarning').to.not.be.present;
  },
};

module.exports = {
  commands: [pageCommands],
  elements: {
    alertWarning: '.alert-info',
    alertError: '.alert-error',
  },
};
