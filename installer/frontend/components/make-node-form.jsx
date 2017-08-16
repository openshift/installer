import _ from 'lodash';

import {
  AWS_CONTROLLERS,
  AWS_WORKERS,
  AWS_ETCDS,
  ENTITLEMENTS,
  INSTANCE_TYPE,
  NUMBER_OF_INSTANCES,
  STORAGE_IOPS,
  STORAGE_SIZE_IN_GIB,
  STORAGE_TYPE,
} from '../cluster-config';

import { Field, Form } from '../form';
import { AWS_INSTANCE_TYPES } from '../facts';
import { validate } from '../validate';

const toKey = (name, field) => `${name}-${field}`;

export const makeNodeForm = (name, instanceValidator = validate.int({min: 1, max: 999}), opts) => {
  const storageType = toKey(name, STORAGE_TYPE);

  // all fields must have a unique name!
  const fields = [
    new Field(toKey(name, NUMBER_OF_INSTANCES), {
      validator: instanceValidator,
      default: 3,
      name: NUMBER_OF_INSTANCES,
    }),
    new Field(toKey(name, INSTANCE_TYPE), {
      validator: validate.nonEmpty,
      default: 't2.medium',
      name: INSTANCE_TYPE,
    }),
    new Field(toKey(name, STORAGE_SIZE_IN_GIB), {
      validator: validate.int({min: 1, max: 15999}),
      default: 30,
      name: STORAGE_SIZE_IN_GIB,
    }),
    new Field(storageType, {
      validator: validate.nonEmpty,
      default: 'gp2',
      name: STORAGE_TYPE,
    }),
    new Field(toKey(name, STORAGE_IOPS), {
      validator: validate.int({min: 100, max: 20000}),
      default: 1000,
      name: STORAGE_IOPS,
      dependencies: [storageType],
      ignoreWhen: cc => cc[storageType] !== 'io1',
    }),
  ];
  const validator = (data, clusterConfig) => {
    const type = data[STORAGE_TYPE];
    const size = data[STORAGE_SIZE_IN_GIB];
    const ops = data[STORAGE_IOPS];
    if (type === 'io1' && ops > size * 50) {

      return `IOPS can't be larger than ${size * 50} (50 IOPS/GiB)`;
    }

    // ETCD resources do not count against entitlements
    //  because they are a temporary hack until hosted
    if (name === AWS_ETCDS) {
      return;
    }

    const entitlements = clusterConfig[ENTITLEMENTS];

    if (!entitlements) {
      return;
    }

    if (entitlements.bypass) {
      return;
    }

    const controllerNumber = clusterConfig[toKey(AWS_CONTROLLERS, NUMBER_OF_INSTANCES)];
    const controllerInstanceType = clusterConfig[toKey(AWS_CONTROLLERS, INSTANCE_TYPE)];
    const workerNumber = clusterConfig[toKey(AWS_WORKERS, NUMBER_OF_INSTANCES)];
    const workerInstanceType = clusterConfig[toKey(AWS_WORKERS, INSTANCE_TYPE)];

    const msg = 'Your Tectonic license limits you to a total of';
    const nodeCount = controllerNumber + workerNumber;
    if (entitlements.nodeCount && nodeCount > entitlements.nodeCount) {
      return `${msg} ${entitlements.nodeCount} nodes. You have ${nodeCount}.`;
    }

    const controllerType = _.find(AWS_INSTANCE_TYPES, v => v.value === controllerInstanceType);
    const workerType = _.find(AWS_INSTANCE_TYPES, v => v.value === workerInstanceType);
    let vcpus = controllerNumber * controllerType.vcpus;
    vcpus += controllerNumber * workerType.vcpus;
    if (entitlements.vCPUsCount && vcpus > entitlements.vCPUsCount) {
      return `${msg} ${entitlements.vCPUsCount} vCPUs. You have ${vcpus}.`;
    }
  };
  return new Form(name, fields, _.defaults({validator}, opts));
};
