import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { CIDR, Connect, Deselect } from './ui';
import { cidrEnd, cidrSize, cidrStart } from '../cidr';
import { DESELECTED_FIELDS } from '../cluster-config';

const CIDRTooltip = connect(
  ({clusterConfig}, {field}) => ({cidr: _.get(clusterConfig, field)})
)(({cidr}) => {
  const size = cidrSize(cidr);
  if (!_.isInteger(size)) {
    return null;
  }
  return <div className="tooltip">
    {size.toLocaleString('en', {useGrouping: true})} IP address{size > 1 && 'es'} ({size === 1 ? cidrStart(cidr) : `${cidrStart(cidr)} to ${cidrEnd(cidr)}`})
  </div>;
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
