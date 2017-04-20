import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

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
import { toError, toAsyncError } from '../utils';
import { AWS_INSTANCE_TYPES } from '../facts';
import { validate } from '../validate';
import { NumberInput, Connect, Select } from './ui';

const toKey = (name, field) => `${name}-${field}`;

const Row = ({label, htmlFor, children}) => <div className="row form-group">
  <div className="col-xs-4">
    <label htmlFor={htmlFor}>{label}</label>
  </div>
  <div className="col-xs-8">
    {children}
  </div>
</div>;

export const makeNodeForm = (name, instanceValidator=validate.int({min: 1, max: 999}), opts) => {
  const storageType = toKey(name, STORAGE_TYPE);

  // all fields must have a unique name!
  const fields = [
    new Field(toKey(name, NUMBER_OF_INSTANCES), {
      validator: instanceValidator,
      default: 1,
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

const IOPs = connect(
  ({clusterConfig}, {fieldName}) => ({type: clusterConfig[toKey(fieldName, STORAGE_TYPE)]})
)(
  ({type, fieldName}) => type !== 'io1' ? null : <Row htmlFor={`${fieldName}--storage-iops`} label="Storage Speed">
    <Connect field={toKey(fieldName, STORAGE_IOPS)}>
      <NumberInput id={`${fieldName}--storage-iops`} className="wiz-super-short-input" suffix="IOPS" />
    </Connect>
  </Row>
);

const Errors = connect(
  ({clusterConfig}, {type}) => ({
    error: _.get(clusterConfig, toError(type)) || _.get(clusterConfig, toAsyncError(type)),
  })
)(props => props.error ? <div className="wiz-error-message">{props.error}</div> : <span/>);

export const DefineNode = ({name, type, disabled, withoutTitle, max}) =>
  <div>
    { !withoutTitle &&
      <div>
        <h3>{name}</h3>
        <br />
      </div>
    }
    <Row htmlFor={`${name}--number`} label="Instances">
      <Connect field={toKey(type, NUMBER_OF_INSTANCES)}>
        <NumberInput
          id={`${name}--number`}
          className="wiz-super-short-input"
          disabled={disabled}
          min="1"
          max={max || 1000} />
      </Connect>
    </Row>
    <Row htmlFor={`${name}--instance`} label="Instance Type">
      <Connect field={toKey(type, INSTANCE_TYPE)}>
        <Select id={`${name}--instance`}>
          <option value="" disabled>Please select AWS EC2 instance type</option>
          {
            AWS_INSTANCE_TYPES.map(({value, label}, ix) => {
              return <option value={value} key={`${value}-${ix}`}>{label}</option>;
            })
          }
        </Select>
      </Connect>
    </Row>
    <Row htmlFor={`${name}--storage-size`} label="Storage Size">
      <Connect field={toKey(type, STORAGE_SIZE_IN_GIB)}>
        <NumberInput id={`${name}--storage-size`} className="wiz-super-short-input" suffix="GiB" />
        </Connect>
    </Row>
    <Row htmlFor={`${name}--storage-type`} label="Storage Type">
      <Connect field={toKey(type, STORAGE_TYPE)}>
        <Select id={`${name}--storage-type`}>
          <option value="" disabled>Please select storage type</option>
          <option value="gp2" key="gp2">General Purpose SSD (GP2)</option>
          <option value="io1" key="io1">Provisioned IOPS SSD (IO1)</option>
          <option value="standard" key="standard">Magnetic</option>
        </Select>
      </Connect>
    </Row>

    <IOPs fieldName={type} />

    <Errors type={type} />
  </div>;

const fields = [
  makeNodeForm(AWS_CONTROLLERS),
  makeNodeForm(AWS_WORKERS),
];

const DefineNodesForm = 'DefineNodesForm';
const form = new Form(DefineNodesForm, fields, {
  // TODO: add after we have an ENTITLEMENTS field...
  // dependencies: [ENTITLEMENTS],
});

export const AWS_DefineNodes = () =>
  <div>
    <DefineNode type={AWS_CONTROLLERS} name="Masters" max={10} />
    <hr/>
    <DefineNode type={AWS_WORKERS} name="Workers" />
    <form.Errors />
  </div>;

AWS_DefineNodes.canNavigateForward = form.canNavigateForward;
