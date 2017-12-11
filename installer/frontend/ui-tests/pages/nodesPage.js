const nodesPageCommands = {
  test(json) {
    this.selectOption('input[type=radio]#external');
    this.expect.element('@etcdCount').to.not.be.present;
    this.expect.element('@externalEtcdAddress').to.be.present;

    this.selectOption('input[type=radio]#provisioned');
    this.expect.element('@etcdCount').to.be.present;
    this.expect.element('@externalEtcdAddress').to.not.be.present;

    const testInstanceCount = (field, value, max) => {
      this
        // Browser won't let us set this to zero because of the min="0" attribute, so use -1 instead
        .setField(field, -1)
        .expectValidationErrorContains('Cannot be less than 1')
        .setField(field, 1)
        .expectNoValidationError()
        .setField(field, max)
        .expectNoValidationError()
        .setField(field, max + 1)
        .expectValidationErrorContains(`Cannot be greater than ${max}`)
        .setField(field, value);
    };

    testInstanceCount('@mastersCount', json['aws_controllers-numberOfInstances'], 10);
    testInstanceCount('@workersCount', json['aws_workers-numberOfInstances'], 1000);
    testInstanceCount('@etcdCount', json['aws_etcds-numberOfInstances'], 9);

    this
      .selectOption(`[id=aws_controllers--instance] [value="${json['aws_controllers-instanceType']}"]`)
      .setField('[id=aws_controllers--storage-size]', json['aws_controllers-storageSizeInGiB'])
      .selectOption(`[id=aws_workers--instance] [value="${json['aws_workers-instanceType']}"]`)
      .setField('[id=aws_workers--storage-size]', json['aws_workers-storageSizeInGiB'])
      .selectOption(`[id=aws_etcds--instance] [value="${json['aws_etcds-instanceType']}"]`)
      .setField('[id=aws_etcds--storage-size]', json['aws_etcds-storageSizeInGiB']);
  },
};

module.exports = {
  commands: [nodesPageCommands],
  elements: {
    etcdCount: '[id=aws_etcds--number]',
    externalEtcdAddress: 'input#externalETCDClient[type=text]',
    mastersCount: '[id=aws_controllers--number]',
    workersCount: '[id=aws_workers--number]',
  },
};
