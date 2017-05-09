import React from 'react';

export const Alert = props => <div className={`alert alert-${props.severity || 'info'}`}
  style={{padding: 15, marginBottom: 10, marginTop: 10}}>
  { !props.noIcon && <span><i className="fa fa-fw fa-exclamation-triangle"></i>&nbsp;</span> }
  {props.children}
</div>;
