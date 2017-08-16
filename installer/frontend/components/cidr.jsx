import React from 'react';

import { Deselect, Input, WithClusterConfig } from './ui';
import { validate } from '../validate';
import { DESELECTED_FIELDS } from '../cluster-config.js';

export const CIDR = ({field, name, disabled, placeholder, autoFocus, validator, selectable, fieldName}) => {
  fieldName = fieldName || field;
  return <div className="row form-group">
    <div className="col-xs-3">
      {selectable && <Deselect field={fieldName} />}
      <label htmlFor={(selectable ? `${DESELECTED_FIELDS}.` : '') + fieldName}>{name}</label>
    </div>
    <div className="col-xs-5">
      <WithClusterConfig field={field} validator={validator || validate.CIDR}>
        <Input placeholder={placeholder} id={field} disabled={disabled} autoFocus={autoFocus} />
      </WithClusterConfig>
    </div>
  </div>;
};
