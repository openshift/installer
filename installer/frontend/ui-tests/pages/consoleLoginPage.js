const consoleLoginPageCommands = {
  test(json) {
    return this
      .setValue('@email', json.tectonic_admin_email)
      .setValue('@password', 'password')
      .setValue('@confirmPassword', 'password');
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
