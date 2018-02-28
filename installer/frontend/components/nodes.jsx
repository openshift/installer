import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { getIamRoles } from '../aws-actions';
import {
  AWS_CONTROLLERS,
  AWS_ETCDS,
  AWS_REGION_FORM,
  AWS_WORKERS,
  ETCD_OPTION,
  ETCD_OPTIONS,
  EXTERNAL_ETCD_CLIENT,
  IAM_ROLE,
  IAM_ROLE_CREATE_OPTION,
  INSTANCE_TYPE,
  NUMBER_OF_INSTANCES,
  PLATFORM_TYPE,
  STORAGE_IOPS,
  STORAGE_SIZE_IN_GIB,
  STORAGE_TYPE,
} from '../cluster-config';

import { AWS_INSTANCE_TYPES } from '../facts';
import { Field, Form } from '../form';
import { AWS_TF } from '../platforms';
import { compose, validate } from '../validate';
import { A, Connect, Input, NumberInput, Radio, Select } from './ui';

const toId = (name, field) => `${name}-${field}`;

// Use this single dummy form / field to trigger loading the IAM roles list. Then IAM role fields can set this as their
// dependency, which avoids triggering a separate API request for each field.
new Form('DUMMY_NODE_FORM', [
  new Field(IAM_ROLE, {
    default: 'DUMMY_VALUE',
    name: IAM_ROLE,
    dependencies: [AWS_REGION_FORM],
    getExtraStuff: (dispatch, isNow) => dispatch(getIamRoles(null, null, isNow)),
  }),
]);

const makeNodeForm = (name, instanceValidator, withIamRole = true, opts) => {
  const storageTypeId = toId(name, STORAGE_TYPE);

  // all fields must have a unique name!
  const fields = [
    new Field(toId(name, NUMBER_OF_INSTANCES), {
      validator: instanceValidator,
      default: 3,
      name: NUMBER_OF_INSTANCES,
    }),
    new Field(toId(name, INSTANCE_TYPE), {
      validator: validate.nonEmpty,
      default: 't2.medium',
      name: INSTANCE_TYPE,
    }),
    new Field(toId(name, STORAGE_SIZE_IN_GIB), {
      validator: validate.int({min: 30, max: 15999}),
      default: 30,
      name: STORAGE_SIZE_IN_GIB,
    }),
    new Field(storageTypeId, {
      validator: validate.nonEmpty,
      default: 'gp2',
      name: STORAGE_TYPE,
    }),
    new Field(toId(name, STORAGE_IOPS), {
      validator: validate.int({min: 100, max: 20000}),
      default: 1000,
      name: STORAGE_IOPS,
      dependencies: [storageTypeId],
      ignoreWhen: cc => cc[storageTypeId] !== 'io1',
    }),
  ];

  if (withIamRole) {
    fields.unshift(new Field(toId(name, IAM_ROLE), {
      default: IAM_ROLE_CREATE_OPTION,
      name: IAM_ROLE,
      dependencies: [IAM_ROLE],
    }));
  }

  const validator = (data) => {
    const type = data[STORAGE_TYPE];
    const maxIops = 50 * data[STORAGE_SIZE_IN_GIB];
    const ops = data[STORAGE_IOPS];

    if (type === 'io1' && ops > maxIops) {
      return `IOPS can't be larger than ${maxIops} (50 IOPS/GiB)`;
    }
  };

  return new Form(name, fields, _.defaults({validator}, opts));
};

const Row = ({label, htmlFor, children}) => <div className="row form-group">
  <div className="col-xs-4">
    <label htmlFor={htmlFor}>{label}</label>
  </div>
  <div className="col-xs-8">
    {children}
  </div>
</div>;

const IOPs = connect(
  ({clusterConfig}, {fieldName}) => ({type: clusterConfig[toId(fieldName, STORAGE_TYPE)]})
)(
  ({type, fieldName}) => type !== 'io1' ? null : <Row htmlFor={`${fieldName}--storage-iops`} label="Storage Speed">
    <Connect field={toId(fieldName, STORAGE_IOPS)}>
      <NumberInput id={`${fieldName}--storage-iops`} className="wiz-super-short-input" suffix="&nbsp;&nbsp;IOPS" />
    </Connect>
  </Row>
);

const IamRoles = connect(
  ({clusterConfig}) => ({roles: _.get(clusterConfig, ['extra', IAM_ROLE], [])})
)(
  ({roles, type}) => <Row htmlFor={`${type}--iam-role`} label="IAM Role">
    <Connect field={toId(type, IAM_ROLE)}>
      <Select id={`${type}--iam-role`}>
        <option value={IAM_ROLE_CREATE_OPTION}>Create an IAM role for me (default)</option>
        {_.isArray(roles) && roles.map(r => <option value={r} key={r}>{r}</option>)}
      </Select>
    </Connect>
    {!_.isArray(roles) && <div className="wiz-error-message">Could not load IAM role list</div>}
  </Row>
);

const DefineNode = ({type, max, withIamRole = true}) => <div>
  {withIamRole && <IamRoles type={type} />}
  <Row htmlFor={`${type}--number`} label="Instances">
    <Connect field={toId(type, NUMBER_OF_INSTANCES)}>
      <NumberInput className="wiz-super-short-input" id={`${type}--number`} min="1" max={max} />
    </Connect>
  </Row>
  <Row htmlFor={`${type}--instance`} label="Instance Type">
    <Connect field={toId(type, INSTANCE_TYPE)}>
      <Select id={`${type}--instance`}>
        <option value="" disabled>Please select AWS EC2 instance type</option>
        {AWS_INSTANCE_TYPES.map(({value, label}) => <option value={value} key={value}>{label}</option>)}
      </Select>
    </Connect>
    {type === 'aws_etcds' && <p className="text-muted wiz-help-text">
      Read the <A href="https://coreos.com/etcd/docs/latest/op-guide/hardware.html" rel="noopener">etcd recommended hardware</A> guide for best performance.
    </p>}
  </Row>
  <Row htmlFor={`${type}--storage-size`} label="Storage Size">
    <Connect field={toId(type, STORAGE_SIZE_IN_GIB)}>
      <NumberInput id={`${type}--storage-size`} className="wiz-super-short-input" suffix="&nbsp;&nbsp;GiB" />
    </Connect>
  </Row>
  <Row htmlFor={`${type}--storage-type`} label="Storage Type">
    <Connect field={toId(type, STORAGE_TYPE)}>
      <Select id={`${type}--storage-type`}>
        <option value="" disabled>Please select storage type</option>
        <option value="gp2" key="gp2">General Purpose SSD (GP2)</option>
        <option value="io1" key="io1">Provisioned IOPS SSD (IO1)</option>
        <option value="standard" key="standard">Magnetic</option>
      </Select>
    </Connect>
  </Row>

  <IOPs fieldName={type} />
</div>;

export const MAX_MASTERS = 100;
export const MAX_WORKERS = 1000;

const etcdNodeForm = makeNodeForm(AWS_ETCDS, compose(validate.int({min: 1, max: 9}), validate.isOdd), false, {
  dependencies: [ETCD_OPTION],
  ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.PROVISIONED,
});
const etcdForm = new Form('etcdForm', [
  new Field(ETCD_OPTION, {default: ETCD_OPTIONS.PROVISIONED}),
  etcdNodeForm,
  new Field(EXTERNAL_ETCD_CLIENT, {
    default: '',
    validator: compose(validate.nonEmpty, validate.host),
    dependencies: [ETCD_OPTION],
    ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.EXTERNAL,
  }),
]);

const mastersForm = makeNodeForm(AWS_CONTROLLERS, validate.int({min: 1, max: MAX_MASTERS}));
const workersForm = makeNodeForm(AWS_WORKERS, validate.int({min: 1, max: MAX_WORKERS}));
const awsNodesForm = new Form('awsNodesForm', [mastersForm, workersForm, etcdForm]);

export const Etcd = connect(({clusterConfig}) => ({
  etcdOption: clusterConfig[ETCD_OPTION],
  isAWS: clusterConfig[PLATFORM_TYPE] === AWS_TF,
}))(
  ({etcdOption, isAWS}) => <div>
    <div className="row form-group">
      <div className="col-xs-12">
        etcd is the key-value store used by Kubernetes for cluster coordination and state management.
      </div>
    </div>

    <div className="row form-group">
      <div className="col-xs-12">
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={ETCD_OPTION}>
                <Radio name={ETCD_OPTION} value={ETCD_OPTIONS.PROVISIONED} id={ETCD_OPTIONS.PROVISIONED} />
              </Connect>
              {isAWS && <span>Provision AWS etcd cluster</span>}
              {!isAWS && <span>Provision etcd cluster directly on controller nodes</span>}
            </label>&nbsp;(default)
            <p className="text-muted wiz-help-text">
              {isAWS && <span>Create EC2 instances to run an etcd cluster.</span>}
              {!isAWS && <span>Run etcd directly on controller nodes.</span>}
            </p>
          </div>
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={ETCD_OPTION}>
                <Radio name={ETCD_OPTION} value={ETCD_OPTIONS.EXTERNAL} id={ETCD_OPTIONS.EXTERNAL} />
              </Connect>
              I have my own etcd v3 cluster
            </label>
            <p className="text-muted wiz-help-text">Your Tectonic cluster will be configured to use an external etcd, which you specify.</p>
          </div>
        </div>
      </div>
    </div>
    {etcdOption === ETCD_OPTIONS.EXTERNAL && <div>
      <hr />
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={EXTERNAL_ETCD_CLIENT}>Client Address</label>
        </div>
        <div className="col-xs-9">
          <Connect field={EXTERNAL_ETCD_CLIENT}>
            <Input id={EXTERNAL_ETCD_CLIENT}
              autoFocus={true}
              className="wiz-inline-field wiz-inline-field--prefix wiz-inline-field--suffix"
              prefix={<span className="input__prefix">https://</span>}
              suffix={<span className="input__suffix">:2379</span>}
              placeholder="etcd.example.com" />
          </Connect>
          <p className="text-muted">Address of etcd client endpoint</p>
        </div>
      </div>
    </div>
    }
    {isAWS && etcdOption === ETCD_OPTIONS.PROVISIONED && <div>
      <hr />
      <div className="row form-group col-xs-12">
        <DefineNode type={AWS_ETCDS} max={9} withIamRole={false} />
        <etcdNodeForm.Errors />
      </div>
    </div>
    }
  </div>
);

export const AWS_Nodes = () => <div>
  <h3>Master Nodes</h3>
  <br />
  <DefineNode type={AWS_CONTROLLERS} max={MAX_MASTERS} />
  <mastersForm.Errors />
  <hr />
  <h3>Worker Nodes</h3>
  <br />
  <DefineNode type={AWS_WORKERS} max={MAX_WORKERS} />
  <workersForm.Errors />
  <hr />
  <h3>etcd Nodes</h3>
  <br />
  <Etcd />
</div>;

Etcd.canNavigateForward = etcdForm.canNavigateForward;
AWS_Nodes.canNavigateForward = awsNodesForm.canNavigateForward;
