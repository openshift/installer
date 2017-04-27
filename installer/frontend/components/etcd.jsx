import React from 'react';
import { connect } from 'react-redux';

import {
  AWS_ETCDS,
  EXTERNAL_ETCD_CLIENT,
  EXTERNAL_ETCD_ENABLED,
  PLATFORM_TYPE,
} from '../cluster-config';
import { validate } from '../validate';
import { Connect, Input, RadioBoolean } from './ui';
import { DefineNode } from './aws-define-nodes';
import { makeNodeForm } from './make-node-form';
import { Field, Form } from '../form';
import { AWS, AWS_TF } from '../platforms';

const EtcdForm = 'EtcdForm';
const one2Nine = validate.int({min: 1, max: 9});

const fields = [
  new Field(EXTERNAL_ETCD_ENABLED, {
    default: false,
  }),
  makeNodeForm(AWS_ETCDS, value => one2Nine(value) || validate.isOdd(value), {
    dependencies: [EXTERNAL_ETCD_ENABLED],
    ignoreWhen: cc => cc[EXTERNAL_ETCD_ENABLED],
  }),
  new Field(EXTERNAL_ETCD_CLIENT, {
    default: '',
    validator: validate.hostPort,
    dependencies: [EXTERNAL_ETCD_ENABLED],
    ignoreWhen: cc => !cc[EXTERNAL_ETCD_ENABLED],
  }),
];

const form = new Form(EtcdForm, fields);

export const Etcd = connect(({clusterConfig}) => ({
  externalETCDEnabled: clusterConfig[EXTERNAL_ETCD_ENABLED],
  isAWS: clusterConfig[PLATFORM_TYPE] === AWS || clusterConfig[PLATFORM_TYPE] === AWS_TF,
}))(
class ExternalETCD extends React.Component {
  render () {
    const {externalETCDEnabled, isAWS} = this.props;
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
                    <Connect field={EXTERNAL_ETCD_ENABLED}>
                      <RadioBoolean inverted={true} name="externalETCDEnabled" />
                    </Connect>
                    Launch etcd for me
                  </label>&nbsp;(default)
                  <p className="text-muted wiz-help-text">The installer will automatically launch and configure etcd inside your Tectonic cluster.</p>
                </div>
              </div>
              <div className="wiz-radio-group">
                <div className="radio wiz-radio-group__radio">
                  <label>
                    <Connect field={EXTERNAL_ETCD_ENABLED}>
                      <RadioBoolean name="externalETCDEnabled" />
                    </Connect>
                    I have my own v3 etcd cluster
                  </label>
                  <p className="text-muted wiz-help-text">Your Tectonic cluster will be configured to to use an external etcd, which you specify.</p>
                </div>
              </div>
            </div>
          </div>
          { (isAWS || externalETCDEnabled) && <hr /> }
          {
            externalETCDEnabled &&
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
            !externalETCDEnabled && isAWS &&
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
