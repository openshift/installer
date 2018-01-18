import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { Input, Connect, FieldRowList } from './ui';
import { AWS_TAGS, CLUSTER_NAME } from '../cluster-config';
import { FieldList, Form } from '../form';

const validators = {
  AWSTag: value => {
    if (value.startsWith('aws:')) {
      return 'Tags may not start with "aws:"';
    }
    return;
  },

  AWSTagKey: value => {
    if (!value) {
      return;
    }
    if (value.length > 255) {
      return 'Tags must be <= 255 characters';
    }
    return validators.AWSTag(value);
  },

  AWSTagValue: value => {
    if (!value) {
      return;
    }
    if (value.length > 127) {
      return 'Tags must be <= 127 characters';
    }
    return validators.AWSTag(value);
  },

  AWSTagUniqueKeys: tags => {
    const REQUIRED_MSG = 'Both fields are required';
    const keyCounts = _.countBy(tags, 'key');

    return _.map(tags, tag => {
      const error = {};
      if (tag.key && !tag.value) {
        error.value = REQUIRED_MSG;
      }
      if (!tag.key && tag.value) {
        error.key = REQUIRED_MSG;
      } else if (keyCounts[tag.key] > 1) {
        error.key = 'Tag keys must be unique';
      }
      return error;
    });
  },
};

const rowFields = {
  key: {
    default: '',
    validator: validators.AWSTagKey,
  },
  value: {
    default: '',
    validator: validators.AWSTagValue,
  },
};

const tagsFields = new FieldList(AWS_TAGS, rowFields, {
  validator: validators.AWSTagUniqueKeys,
});

export const tagsForm = new Form('AWS_TAGS_FORM', [tagsFields]);

const Tag = ({autoFocus, placeholder, row}) => <div>
  <div className="col-xs-5" style={{paddingRight: 0}}>
    <Connect field={row.key}>
      <Input autoFocus={autoFocus} placeholder="e.g. Name" />
    </Connect>
  </div>
  <div className="col-xs-6" style={{paddingRight: 0}}>
    <Connect field={row.value}>
      <Input placeholder={placeholder} />
    </Connect>
  </div>
</div>;

export const AWS_Tags = connect(
  ({clusterConfig}) => ({placeholder: `e.g. ${clusterConfig[CLUSTER_NAME] || 'myclustername'}`})
)(({placeholder}) => <div>
  <div className="row">
    <div className="col-xs-5">
      <label className="text-muted cos-thin-label">KEY</label>
    </div>
    <div className="col-xs-6">
      <label className="text-muted cos-thin-label">VALUE</label>
    </div>
  </div>

  <FieldRowList id={AWS_TAGS} placeholder={placeholder} Row={Tag} rowFields={rowFields} />
</div>);
