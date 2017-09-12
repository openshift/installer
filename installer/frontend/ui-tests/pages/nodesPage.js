const nodesPageCommands = {
  test() {
    this.selectOption('input[type=radio]#provisioned');
    this.expect.element('@provisionedEtcdCount').to.be.present;
    this.expect.element('@externalEtcdAddress').to.not.be.present;

    this.selectOption('input[type=radio]#external');
    this.expect.element('@provisionedEtcdCount').to.not.be.present;
    this.expect.element('@externalEtcdAddress').to.be.present;

    this.selectOption('input[type=radio]#selfHosted');
    this.expect.element('@provisionedEtcdCount').to.not.be.present;
    this.expect.element('@externalEtcdAddress').to.not.be.present;
  },
};

module.exports = {
  commands: [nodesPageCommands],
  elements: {
    externalEtcdAddress: 'input#externalETCDClient[type=text]',
    provisionedEtcdCount: 'input[id="etcd--number"][type=number][min="1"]',
  },
};
