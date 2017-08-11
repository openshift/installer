import React from 'react';

import { Deselect, Input, WithClusterConfig } from './ui';
import { validate } from '../validate';
import { DESELECTED_FIELDS } from '../cluster-config.js';
import { WithTooltip } from './tooltip';

const generateTooltipText = ({value}) => {
  if (validate.CIDR(value)) {
    return;
  }
  const [, bits] = value.split('/');
  // javascript's bit shifting only works on signed 32bit ints so <<31
  // would be negative :(
  const addresses = Math.pow(2, 32 - parseInt(bits, 10));
  return `${addresses} IP address${addresses === 1 ? '' : 'es'}`;
};

export const CIDR = ({field, name, disabled, placeholder, autoFocus, validator, selectable, fieldName}) => {
  fieldName = fieldName || field;
  return <div className="row form-group">
    <div className="col-xs-3">
      {selectable && <Deselect field={fieldName} />}
      <label htmlFor={(selectable ? `${DESELECTED_FIELDS}.` : '') + fieldName}>{name}</label>
    </div>
    <div className="col-xs-5 withtooltip">
      <WithClusterConfig field={field} validator={validator || validate.CIDR}>
        <WithTooltip generateText={generateTooltipText}>
          <Input placeholder={placeholder} id={field} disabled={disabled} autoFocus={autoFocus} />
        </WithTooltip>
      </WithClusterConfig>
    </div>
  </div>;
};
