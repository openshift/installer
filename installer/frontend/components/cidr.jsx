import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { CIDR, Connect, Deselect } from './ui';
import { validate } from '../validate';
import { DESELECTED_FIELDS } from '../cluster-config.js';

export const cidrSize = cidr => {
  if (validate.CIDR(cidr)) {
    return undefined;
  }
  const [, bits] = cidr.split('/');

  // JavaScript's bit shifting only works on signed 32bit ints so <<31 would be negative :(
  return Math.pow(2, 32 - parseInt(bits, 10));
};

const CIDRTooltip = connect(
  ({clusterConfig}, {field}) => ({clusterConfig, value: _.get(clusterConfig, field)})
)(({value}) => {
  const addresses = cidrSize(value);
  if (!_.isInteger(addresses)) {
    return null;
  }
  return <div className="tooltip">{addresses.toLocaleString('en', {useGrouping: true})} IP address{addresses > 1 && 'es'}</div>;
});

export const CIDRRow = ({field, name, disabled, placeholder, autoFocus, selectable, fieldName, validator}) => {
  fieldName = fieldName || field;
  return <div className="row form-group">
    <div className="col-xs-4">
      {selectable && <Deselect field={fieldName} />}
      <label htmlFor={(selectable ? `${DESELECTED_FIELDS}.` : '') + fieldName}>{name}</label>
    </div>
    <div className="col-xs-5">
      <div className="withtooltip">
        <Connect field={field}>
          <CIDR autoFocus={autoFocus} disabled={disabled} id={field} placeholder={placeholder} validator={validator} />
        </Connect>
        <CIDRTooltip field={field} />
      </div>
    </div>
  </div>;
};
