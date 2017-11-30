const nodesPageCommands = {
  test(json) {
    this.selectOption('input[type=radio]#external');
    this.expect.element('@etcdCount').to.not.be.present;
    this.expect.element('@externalEtcdAddress').to.be.present;

    this.selectOption('input[type=radio]#provisioned');
    this.expect.element('@etcdCount').to.be.present;
    this.expect.element('@externalEtcdAddress').to.not.be.present;

    this
      .setField('[id=aws_controllers--number]', json['aws_controllers-numberOfInstances'])
      .selectOption(`[id=aws_controllers--instance] [value="${json['aws_controllers-instanceType']}"]`)
      .setField('[id=aws_controllers--storage-size]', json['aws_controllers-storageSizeInGiB'])
      .setField('[id=aws_workers--number]', json['aws_workers-numberOfInstances'])
      .selectOption(`[id=aws_workers--instance] [value="${json['aws_workers-instanceType']}"]`)
      .setField('[id=aws_workers--storage-size]', json['aws_workers-storageSizeInGiB'])
      .setField('@etcdCount', json['aws_etcds-numberOfInstances'])
      .selectOption(`[id=aws_etcds--instance] [value="${json['aws_etcds-instanceType']}"]`)
      .setField('[id=aws_etcds--storage-size]', json['aws_etcds-storageSizeInGiB']);
  },
};

module.exports = {
  commands: [nodesPageCommands],
  elements: {
    etcdCount: '[id=aws_etcds--number]',
    externalEtcdAddress: 'input#externalETCDClient[type=text]',
  },
};
