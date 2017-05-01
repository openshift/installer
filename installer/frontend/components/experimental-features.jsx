import React from 'react';
import { connect } from 'react-redux';
import { CheckBox, WithClusterConfig } from './ui';

import { UPDATER_ENABLED } from '../cluster-config';
import { Alert } from './alert';
import { TectonicGA } from '../tectonic-ga';

export const ExperimentalFeatures = connect(
  state => ({updater: state.clusterConfig[UPDATER_ENABLED]})
)(props => <div className="row form-group">
  <div className="col-xs-3">
    <label htmlFor="tectonicOperator">Experimental Features</label>
  </div>
  <div className="col-xs-9">
    <WithClusterConfig field={UPDATER_ENABLED}>
      <CheckBox id="tectonicOperator" suffix={
        <label htmlFor="tectonicOperator" className="text-muted wiz-help-text">
          Install experimental <a href="https://coreos.com/blog/introducing-operators.html" onClick={TectonicGA.sendDocsEvent} target="_blank">Tectonic Operators</a>.
        </label>
      } />
    </WithClusterConfig>
    { props.updater &&
      <Alert noIcon>
        Only use on clusters that can be easily replaced.
      </Alert>
    }
  </div>
</div>);
