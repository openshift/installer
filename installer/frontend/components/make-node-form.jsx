import _ from 'lodash';

import {
  INSTANCE_TYPE,
  NUMBER_OF_INSTANCES,
  STORAGE_IOPS,
  STORAGE_SIZE_IN_GIB,
  STORAGE_TYPE,
} from '../cluster-config';

import { Field, Form } from '../form';
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

  const validator = (data) => {
    const type = data[STORAGE_TYPE];
    const size = data[STORAGE_SIZE_IN_GIB];
    const ops = data[STORAGE_IOPS];

    if (type === 'io1' && ops > size * 50) {
      return `IOPS can't be larger than ${size * 50} (50 IOPS/GiB)`;
    }
  };

  return new Form(name, fields, _.defaults({validator}, opts));
};
