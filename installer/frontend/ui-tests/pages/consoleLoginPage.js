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
    email: {
      selector: 'input#adminEmail',
    },
    password: {
      selector: 'input#adminPassword',
    },
    confirmPassword: {
      selector: 'input#adminPassword2',
    },
  },
};
