const {testExternalEtcd} = require('./etcdConnectionPage');

const nodesPageCommands = {
  test (json) {
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

    testInstanceCount('@mastersCount', json['aws_controllers-numberOfInstances'], 100);
    testInstanceCount('@workersCount', json['aws_workers-numberOfInstances'], 1000);
    if (json['aws_etcds-numberOfInstances']) {
      testInstanceCount('@etcdCount', json['aws_etcds-numberOfInstances'], 9);
    }

    this
      .selectOption(`[id=aws_controllers--instance] [value="${json['aws_controllers-instanceType']}"]`)
      .selectOption(`[id=aws_workers--instance] [value="${json['aws_workers-instanceType']}"]`)
      .expectNoValidationError();

    const etcdInstanceType = json['aws_etcds-instanceType'];
    if (etcdInstanceType) {
      this
        .selectOption(`[id=aws_etcds--instance] [value="${etcdInstanceType}"]`)
        .expectNoValidationError();
    }

    const testVolumeSize = (field, value) => {
      const min = 30;
      this
        .setField(field, min - 1)
        .expectValidationErrorContains(`Cannot be less than ${min}`)
        .setField(field, min)
        .expectNoValidationError();

      if (value) {
        this
          .setField(field, value)
          .expectNoValidationError();
      }
    };

    testVolumeSize('[id=aws_controllers--storage-size]', json['aws_controllers-storageSizeInGiB']);
    testVolumeSize('[id=aws_workers--storage-size]', json['aws_workers-storageSizeInGiB']);
    testVolumeSize('[id=aws_etcds--storage-size]', json['aws_etcds-storageSizeInGiB']);

    const mastersStorageType = json['aws_controllers-storageType'];
    if (mastersStorageType) {
      this.selectOption(`[id=aws_controllers--storage-type] [value="${mastersStorageType}"]`);
    }
    const workersStorageType = json['aws_workers-storageType'];
    if (workersStorageType) {
      this.selectOption(`[id=aws_workers--storage-type] [value="${workersStorageType}"]`);
    }
    this.expectNoValidationError();

    const mastersIOPS = json['aws_controllers-storageIOPS'];
    if (mastersIOPS) {
      this.setField('[id=aws_controllers--storage-iops]', mastersIOPS);
    }
    const workersIOPS = json['aws_workers-storageIOPS'];
    if (workersIOPS) {
      this.setField('[id=aws_workers--storage-iops]', workersIOPS);
    }
    this.expectNoValidationError();

    testExternalEtcd(this);

    if (json.etcdOption === 'external') {
      this
        .selectOption('#external')
        .setField('#externalETCDClient', json.externalETCDClient)
        .expectNoValidationError();
    }
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
