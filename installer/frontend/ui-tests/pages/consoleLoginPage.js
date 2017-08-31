const consoleLoginPageCommands = {
  test(json) {
    return this
      .setField('@email', json.tectonic_admin_email)
      .setField('@password', 'password')
      .setField('@confirmPassword', 'password');
  },
};

module.exports = {
  commands: [consoleLoginPageCommands],
  elements: {
    email: 'input#adminEmail',
    password: 'input#adminPassword',
    confirmPassword: 'input#adminPassword2',
  },
};
