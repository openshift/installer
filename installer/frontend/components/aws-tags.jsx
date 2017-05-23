import _ from 'lodash';
import React from 'react';

import { Input, Connect } from './ui';
import { AWS_TAGS } from '../cluster-config';
import { FieldList } from '../form';

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
    const keys = _.map(tags, t => t.key);
    const errors = {};
    let i = 1;
    for (let name1 of keys) {
      for (let name2 of keys.slice(i)) {
        if (name1 === name2) {
          errors[i] = {key: 'Tag keys must be unique'};
        }
        i += 1;
      }
    }
    _.each(tags, (tag, index) => {
      if (tag.value ? !tag.key : tag.key) {
        errors[index] = errors[index] || {};
        errors[index].value = 'Both fields are required';
      }
    });
    return errors;
  },
};

export const tagsFields = new FieldList(AWS_TAGS, {
  fields: {
    key: {
      default: '',
      validator: validators.AWSTagKey,
    },
    value: {
      default: '',
      validator: validators.AWSTagValue,
    },
  },
  validator: validators.AWSTagUniqueKeys,
});

const Tag = ({row, remove, placeholder}) =>
  <div className="row" style={{padding: '0 0 20px 0'}}>
    <div className="col-xs-5" style={{paddingRight: 0}}>
      <Connect field={row.key}>
        <Input placeholder="e.g. Name" blurry />
      </Connect>
    </div>
    <div className="col-xs-6" style={{paddingRight: 0}}>
      <Connect field={row.value}>
        <Input placeholder={placeholder} blurry />
      </Connect>
    </div>
    <div className="col-xs-1">
      <i className="fa fa-minus-circle list-add-or-subtract pull-right" onClick={remove}></i>
    </div>
  </div>;

export const AWS_Tags = props => {
  return <div>
    <div className="row">
      <div className="col-xs-5">
        <label className="text-muted cos-thin-label">KEY</label>
      </div>
      <div className="col-xs-6">
        <label className="text-muted cos-thin-label">VALUE</label>
      </div>
    </div>

    <tagsFields.Map>
      <Tag {...props} />
    </tagsFields.Map>

    <div className="row">
      <div className="col-xs-3">
        <span className="wiz-link" onClick={tagsFields.addOnClick}>
          <i className="fa fa-plus-circle list-add wiz-link"></i>&nbsp; Add More
        </span>
      </div>
    </div>
  </div>;
};
