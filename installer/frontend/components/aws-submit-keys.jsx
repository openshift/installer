import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { validate } from '../validate';
import { Connect, Selector } from './ui';
import { Field, Form } from '../form';

import * as awsActions from '../aws-actions';
import { AWS_SSH, AWS_REGION_FORM, AWS_REGION } from '../cluster-config';

const awsSshForm = new Form('AWSSSHForm', [
  new Field(AWS_SSH, {
    default: '',
    validator: validate.nonEmpty,
    dependencies: [AWS_REGION_FORM],
    getExtraStuff: (dispatch, isNow) => dispatch(awsActions.getSsh(null, null, isNow)).then(options => ({options: _.sortBy(options, 'label')})),
  })], {
  validator: (data, cc) => {
    const key = data[AWS_SSH];
    const options = _.get(cc, ['extra', AWS_SSH, 'options']);
    if (options && key && !_.some(options, o => o.value === key)) {
      return `SSH key ${key} does not exist in this region.`;
    }
  },
}
);

const Title = connect(
  ({clusterConfig}) => ({region: clusterConfig[AWS_REGION]})
)(
  ({region}) => <h4>SSH Keys in {region}</h4>
);

export const AWS_SubmitKeys = () => <div>
  <div className="row form-group">
    <div className="col-xs-12">
      {/* eslint-disable react/jsx-no-target-blank */}
      <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#ssh-key" rel="noopener" target="_blank">Generate a new key</a> if you don't have an existing one in this region.
      {/* eslint-enable react/jsx-no-target-blank */}
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-12">
      <Title />
      <Connect field={AWS_SSH}>
        <Selector refreshBtn={true} disabledValue="Please select a SSH Key Pair from this region." />
      </Connect>
      <awsSshForm.Errors />
    </div>
  </div>
</div>;

AWS_SubmitKeys.canNavigateForward = awsSshForm.canNavigateForward;
