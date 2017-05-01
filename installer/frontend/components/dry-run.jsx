import React from 'react';

import { connect } from 'react-redux';
import { SaveAssets } from './save-assets';
import { isTerraform } from '../platforms';
import { PLATFORM_TYPE } from '../cluster-config';

const stateToProps = ({commitState, clusterConfig}) => ({
  blob: commitState.response,
  terraForm: isTerraform(clusterConfig[PLATFORM_TYPE]),
});

export const DryRun = connect(stateToProps)(
({blob, terraForm}) => <div className="row">
  <div className="col-sm-12">
    <div className="form-group">
      Your cluster assets have been created. You can download these assets and customize them as needed.
    </div>
    <div className="from-group">
      <div className="wiz-giant-button-container">
      { terraForm
        ? <a href="/terraform/assets" download>
            <button className="btn btn-primary wiz-giant-button">
              <i className="fa fa-download"></i>&nbsp;&nbsp;Download assets
            </button>
          </a>
        : <SaveAssets blob={blob} />
      }
      </div>
    </div>
  </div>
</div>);

DryRun.canNavigateForward = () => false;
