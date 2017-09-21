import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import {
  AWS_ETCDS,
  EXTERNAL_ETCD_CLIENT,
  ETCD_OPTION,
  ETCD_OPTIONS,
  PLATFORM_TYPE,
} from '../cluster-config';
import { validate } from '../validate';
import { Connect, Input, Radio } from './ui';
import { DefineNode } from './aws-define-nodes';
import { makeNodeForm } from './make-node-form';
import { Field, Form } from '../form';
import { AWS_TF } from '../platforms';

const EtcdForm = 'EtcdForm';
const one2Nine = validate.int({min: 1, max: 9});

const fields = [
  new Field(ETCD_OPTION, {
    default: ETCD_OPTIONS.PROVISIONED,
  }),
  makeNodeForm(AWS_ETCDS, value => one2Nine(value) || validate.isOdd(value), {
    dependencies: [ETCD_OPTION],
    ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.PROVISIONED,
  }),
  new Field(EXTERNAL_ETCD_CLIENT, {
    default: '',
    validator: validate.hostPort,
    dependencies: [ETCD_OPTION],
    ignoreWhen: cc => cc[ETCD_OPTION] !== ETCD_OPTIONS.EXTERNAL,
  }),
];

const form = new Form(EtcdForm, fields, {
  validator: value => {
    const etcd = value[ETCD_OPTION];
    if (!_.values(ETCD_OPTIONS).includes(etcd)) {
      return 'Please select an option.';
    }
  },
});

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
              I have my own v3 etcd cluster
            </label>
            <p className="text-muted wiz-help-text">Your Tectonic cluster will be configured to use an external etcd, which you specify.</p>
          </div>
        </div>
        <form.Errors />
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
                autoFocus
                className="wiz-inline-field wiz-inline-field--protocol"
                prefix={<span className="input__prefix--protocol">http://</span>}
                placeholder="etcd.example.com:2379" />
            </Connect>
            <p className="text-muted">Hostname and port of etcd client endpoint</p>
          </div>
        </div>
      </div>
    }
    {isAWS && etcdOption === ETCD_OPTIONS.PROVISIONED && <hr />}
    {isAWS && etcdOption === ETCD_OPTIONS.PROVISIONED &&
      <div className="row form-group col-xs-12">
        <DefineNode type={AWS_ETCDS} name="etcd" withoutTitle={true} max={9} />
      </div>
    }
  </div>
);

Etcd.canNavigateForward = form.canNavigateForward;
