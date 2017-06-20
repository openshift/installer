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
import { AWS_TF, BARE_METAL_TF } from '../platforms';

const EtcdForm = 'EtcdForm';
const one2Nine = validate.int({min: 1, max: 9});

const fields = [
  new Field(ETCD_OPTION, {
    default: ETCD_OPTIONS.SELF_HOSTED,
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
  validator: (value, clusterConfig) => {
    if (clusterConfig[PLATFORM_TYPE] === BARE_METAL_TF && value === ETCD_OPTIONS.PROVISIONED) {
      return 'Please select an option.';
    }
    if (!(value in ETCD_OPTIONS)) {
      return 'Please select an option.';
    }
  },
});

export const Etcd = connect(({clusterConfig}) => ({
  etcdOption: clusterConfig[ETCD_OPTION],
  isAWS: clusterConfig[PLATFORM_TYPE] === AWS_TF,
}))(
class ExternalETCD extends React.Component {
  render () {
    const {etcdOption, isAWS} = this.props;

    return (
        <div>
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
                      <Radio name={ETCD_OPTION} value={ETCD_OPTIONS.SELF_HOSTED} />
                    </Connect>
                    Create self-hosted etcd cluster (alpha)
                  </label>&nbsp;(default)
                  <p className="text-muted wiz-help-text">The installer will automatically launch and configure etcd inside your Tectonic cluster.</p>
                </div>
              </div>
              <div className="wiz-radio-group">
                <div className="radio wiz-radio-group__radio">
                  <label>
                    <Connect field={ETCD_OPTION}>
                      <Radio name={ETCD_OPTION} value={ETCD_OPTIONS.EXTERNAL} />
                    </Connect>
                    I have my own v3 etcd cluster
                  </label>
                  <p className="text-muted wiz-help-text">Your Tectonic cluster will be configured to to use an external etcd, which you specify.</p>
                </div>
              </div>
              { isAWS && <div className="radio wiz-radio-group__radio">
                <label>
                  <Connect field={ETCD_OPTION}>
                    <Radio name={ETCD_OPTION} value={ETCD_OPTIONS.PROVISIONED} />
                  </Connect>
                  Provision AWS etcd cluster
                </label>
                <p className="text-muted wiz-help-text">Create EC2 instances to run an etcd cluster.</p>
              </div> }
            </div>
          </div>
          { (etcdOption !== ETCD_OPTIONS.SELF_HOSTED) && <hr /> }
          {
            etcdOption === ETCD_OPTIONS.EXTERNAL &&
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
                           placeholder="etcd.example.com:2379"/>
                  </Connect>
                  <p className="text-muted">Hostname and port of etcd client endpoint</p>
                </div>
              </div>
            </div>
          }
          {
            etcdOption === ETCD_OPTIONS.PROVISIONED &&
              <div className="row form-group col-xs-12">
                <DefineNode type={AWS_ETCDS} name="etcd" withoutTitle={true} max={9} />
              </div>
          }
          <form.Errors />
        </div>
    );
  }
});

Etcd.canNavigateForward = form.canNavigateForward;
