import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { Deselect, Input, WithClusterConfig } from './ui';
import { validate } from '../validate';
import { DESELECTED_FIELDS } from '../cluster-config.js';

const CIDRTooltip = connect(
  ({clusterConfig}, {field}) => ({clusterConfig: clusterConfig, value: _.get(clusterConfig, field)})
)(({value}) => {
  if (validate.CIDR(value)) {
    return null;
  }
  const [, bits] = value.split('/');
  // javascript's bit shifting only works on signed 32bit ints so <<31
  // would be negative :(
  const addresses = Math.pow(2, 32 - parseInt(bits, 10));
  return <div className="tooltip">{addresses} IP address{addresses > 1 && 'es'}</div>;
});

export const CIDR = ({field, name, disabled, placeholder, autoFocus, validator, selectable, fieldName}) => {
  fieldName = fieldName || field;
  return <div className="row form-group">
    <div className="col-xs-3">
      {selectable && <Deselect field={fieldName} />}
      <label htmlFor={(selectable ? `${DESELECTED_FIELDS}.` : '') + fieldName}>{name}</label>
    </div>
    <div className="col-xs-5">
      <div className="withtooltip">
        <WithClusterConfig field={field} validator={validator || validate.CIDR}>
          <Input placeholder={placeholder} id={field} disabled={disabled} autoFocus={autoFocus} />
        </WithClusterConfig>
        <CIDRTooltip field={field} />
      </div>
    </div>
  </div>;
};
