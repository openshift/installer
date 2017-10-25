const consoleLoginPageCommands = {
  test(json) {
    this
      .setField('@email', 'abc')
      .expectValidationErrorContains('Invalid email address')
      .setField('@email', json.adminEmail)
      .expectNoValidationError();

    this.setField('@password', 'password');
    this.setField('@confirmPassword', 'abc');
    this.expect.element('@alertError').text.to.contain('Passwords do not match');

    this.setField('@password', 'abc');
    this.setField('@confirmPassword', 'password');
    this.expect.element('@alertError').text.to.contain('Passwords do not match');

    this.setField('@password', 'password');
    this.setField('@confirmPassword', 'password');
    this.expect.element('@alertError').to.not.be.present;
  },
};

module.exports = {
  commands: [consoleLoginPageCommands],
  elements: {
    alertError: '.alert-error',
    email: 'input#adminEmail',
    password: 'input#adminPassword',
    confirmPassword: 'input#adminPassword2',
  },
};
