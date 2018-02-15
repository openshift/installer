import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { CIDR, Connect, Deselect, localeNum } from './ui';
import { cidrEnd, cidrSize, cidrStart } from '../cidr';

export const CIDRTooltip = ({cidr}) => {
  const size = cidrSize(cidr);
  if (!_.isInteger(size)) {
    return null;
  }
  return <div className="tooltip">
    {localeNum(size)} IP address{size > 1 && 'es'} ({size === 1 ? cidrStart(cidr) : `${cidrStart(cidr)} to ${cidrEnd(cidr)}`})
  </div>;
};

export const CIDRRow = connect(
  ({clusterConfig}, {field}) => ({cidr: _.get(clusterConfig, field)})
)(({autoFocus, cidr, deselectId, disabled, field, placeholder, name, validator}) => <div className="row form-group">
  <div className="col-xs-4">
    {deselectId ? <Deselect field={deselectId} label={name} /> : <label htmlFor={field}>{name}</label>}
  </div>
  <div className="col-xs-5">
    <div className="withtooltip">
      <Connect field={field}>
        <CIDR autoFocus={autoFocus} disabled={disabled} id={field} placeholder={placeholder} validator={validator} />
      </Connect>
      <CIDRTooltip cidr={cidr} />
    </div>
  </div>
</div>);
