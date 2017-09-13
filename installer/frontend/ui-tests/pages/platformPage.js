const wizard = require('../utils/wizard');

const platformPageCommands = {
  test(platformEl) {
    this.expect.element('select#platformType').to.be.visible.before(60000);

    this.selectOption('@awsGUI');
    this.expect.element(wizard.nextStep).to.be.present;
    this.selectOption('@awsAdvanced');
    this.expect.element(wizard.nextStep).to.not.be.present;
    this.selectOption('@azureAdvanced');
    this.expect.element(wizard.nextStep).to.not.be.present;
    this.selectOption('@metalGUI');
    this.expect.element(wizard.nextStep).to.be.present;
    this.selectOption('@metalAdvanced');
    this.expect.element(wizard.nextStep).to.not.be.present;
    this.selectOption('@openstackAdvanced');
    this.expect.element(wizard.nextStep).to.not.be.present;

    this.selectOption(platformEl);
    this.expect.element(wizard.nextStep).to.be.present;
  },
};

module.exports = {
  commands: [platformPageCommands],
  elements: {
    awsAdvanced: 'option[value="aws"]',
    awsGUI: 'option[value="aws-tf"]',
    azureAdvanced: 'option[value="azure"]',
    metalAdvanced: 'option[value="bare-metal"]',
    metalGUI: 'option[value="bare-metal-tf"]',
    openstackAdvanced: 'option[value="openstack"]',
  },
};
