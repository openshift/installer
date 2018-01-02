import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import {
  AWS_CONTROLLERS,
  AWS_ETCDS,
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
import { toError } from '../utils';
import { compose, validate } from '../validate';
import { makeNodeForm, toKey } from './make-node-form';
import { Connect, Input, NumberInput, Radio, Select } from './ui';

const Row = ({label, htmlFor, children}) => <div className="row form-group">
  <div className="col-xs-4">
    <label htmlFor={htmlFor}>{label}</label>
  </div>
  <div className="col-xs-8">
    {children}
  </div>
</div>;

const IOPs = connect(
  ({clusterConfig}, {fieldName}) => ({type: clusterConfig[toKey(fieldName, STORAGE_TYPE)]})
)(
  ({type, fieldName}) => type !== 'io1' ? null : <Row htmlFor={`${fieldName}--storage-iops`} label="Storage Speed">
    <Connect field={toKey(fieldName, STORAGE_IOPS)}>
      <NumberInput id={`${fieldName}--storage-iops`} className="wiz-super-short-input" suffix="&nbsp;&nbsp;IOPS" />
    </Connect>
  </Row>
);

const IamRoles = connect(
  ({clusterConfig}) => ({roles: _.get(clusterConfig, ['extra', IAM_ROLE], [])})
)(
  ({roles, type}) => <Row htmlFor={`${type}--iam-role`} label="IAM Role">
    <Connect field={toKey(type, IAM_ROLE)}>
      <Select id={`${type}--iam-role`}>
        <option value={IAM_ROLE_CREATE_OPTION}>Create an IAM role for me (default)</option>
        {_.isArray(roles) && roles.map(r => <option value={r} key={r}>{r}</option>)}
      </Select>
    </Connect>
    {!_.isArray(roles) && <div className="wiz-error-message">Could not load IAM role list</div>}
  </Row>
);

const Errors = connect(
  ({clusterConfig}, {type}) => ({error: _.get(clusterConfig, toError(type))})
)(props => props.error ? <div className="wiz-error-message">{props.error}</div> : <span />);

export const DefineNode = ({type, max, withIamRole = true}) => <div>
  {withIamRole && <IamRoles type={type} />}
  <Row htmlFor={`${type}--number`} label="Instances">
    <Connect field={toKey(type, NUMBER_OF_INSTANCES)}>
      <NumberInput className="wiz-super-short-input" id={`${type}--number`} min="1" max={max} />
    </Connect>
  </Row>
  <Row htmlFor={`${type}--instance`} label="Instance Type">
    <Connect field={toKey(type, INSTANCE_TYPE)}>
      <Select id={`${type}--instance`}>
        <option value="" disabled>Please select AWS EC2 instance type</option>
        {AWS_INSTANCE_TYPES.map(({value, label}) => <option value={value} key={value}>{label}</option>)}
      </Select>
    </Connect>
    {type === 'aws_etcds' && <p className="text-muted wiz-help-text">
      {/* eslint-disable react/jsx-no-target-blank */}
      Read the <a href="https://coreos.com/etcd/docs/latest/op-guide/hardware.html" rel="noopener" target="_blank">etcd recommended hardware</a> guide for best performance.
      {/* eslint-enable react/jsx-no-target-blank */}
    </p>}
  </Row>
  <Row htmlFor={`${type}--storage-size`} label="Storage Size">
    <Connect field={toKey(type, STORAGE_SIZE_IN_GIB)}>
      <NumberInput id={`${type}--storage-size`} className="wiz-super-short-input" suffix="&nbsp;&nbsp;GiB" />
    </Connect>
  </Row>
  <Row htmlFor={`${type}--storage-type`} label="Storage Type">
    <Connect field={toKey(type, STORAGE_TYPE)}>
      <Select id={`${type}--storage-type`}>
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

const MAX_MASTERS = 10;
const MAX_WORKERS = 1000;

const etcdForm = new Form('etcdForm', [
  new Field(ETCD_OPTION, {default: ETCD_OPTIONS.PROVISIONED}),
  makeNodeForm(AWS_ETCDS, compose(validate.int({min: 1, max: 9}), validate.isOdd), false, {
    dependencies: [ETCD_OPTION],
    ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.PROVISIONED,
  }),
  new Field(EXTERNAL_ETCD_CLIENT, {
    default: '',
    validator: validate.url,
    dependencies: [ETCD_OPTION],
    ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.EXTERNAL,
  }),
], {
  validator: value => {
    const etcd = value[ETCD_OPTION];
    if (!_.values(ETCD_OPTIONS).includes(etcd)) {
      return 'Please select an option.';
    }
  },
});

const awsNodesForm = new Form('awsNodesForm', [
  makeNodeForm(AWS_CONTROLLERS, validate.int({min: 1, max: MAX_MASTERS})),
  makeNodeForm(AWS_WORKERS, validate.int({min: 1, max: MAX_WORKERS})),
  etcdForm,
]);

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
        <etcdForm.Errors />
      </div>
    </div>
    {etcdOption === ETCD_OPTIONS.EXTERNAL && <hr />}
    {etcdOption === ETCD_OPTIONS.EXTERNAL &&
      <div className="form-group">
        <div className="row">
          <div className="col-xs-3">
            <label htmlFor={EXTERNAL_ETCD_CLIENT}>Client Address</label>
          </div>
          <div className="col-xs-8">
            <Connect field={EXTERNAL_ETCD_CLIENT}>
              <Input id={EXTERNAL_ETCD_CLIENT}
                autoFocus={true}
                className="wiz-inline-field wiz-inline-field--suffix"
                suffix={<span className="input__suffix">:2379</span>}
                placeholder="https://etcd.example.com" />
            </Connect>
            <p className="text-muted">Address of etcd client endpoint</p>
          </div>
        </div>
      </div>
    }
    {isAWS && etcdOption === ETCD_OPTIONS.PROVISIONED && <hr />}
    {isAWS && etcdOption === ETCD_OPTIONS.PROVISIONED &&
      <div className="row form-group col-xs-12">
        <DefineNode type={AWS_ETCDS} max={9} withIamRole={false} />
      </div>
    }
  </div>
);

export const AWS_Nodes = () => <div>
  <h3>Master Nodes</h3>
  <br />
  <DefineNode type={AWS_CONTROLLERS} max={MAX_MASTERS} />
  <hr />
  <h3>Worker Nodes</h3>
  <br />
  <DefineNode type={AWS_WORKERS} max={MAX_WORKERS} />
  <hr />
  <h3>etcd Nodes</h3>
  <br />
  <Etcd />
</div>;

Etcd.canNavigateForward = etcdForm.canNavigateForward;
AWS_Nodes.canNavigateForward = awsNodesForm.canNavigateForward;
