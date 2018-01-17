import React from 'react';

import { TectonicGA } from '../tectonic-ga';
import { A } from './ui';

export const DryRun = () => <div className="row">
  <div className="col-xs-12">
    <div className="form-group">
      Your cluster assets have been created. You can download these <A href="https://coreos.com/tectonic/docs/latest/admin/assets-zip.html" rel="noopener">assets</A> and customize underlying infrastructure as needed.
      Note: changes to Kubernetes manifests or Tectonic components run in the cluster are not supported.&nbsp; <A href="https://coreos.com/tectonic/docs/latest/install/aws/manual-boot.html" onClick={TectonicGA.sendDocsEvent} rel="noopener">Read more here.&nbsp;&nbsp;<i className="fa fa-external-link" /></A>
    </div>
    <div className="from-group">
      <div className="wiz-giant-button-container">
        <a href="/terraform/assets" download="assets.zip">
          <button className="btn btn-primary wiz-giant-button">
            <i className="fa fa-download"></i>&nbsp;&nbsp;Download Assets
          </button>
        </a>
      </div>
    </div>
  </div>
</div>;

DryRun.canNavigateForward = () => false;
DryRun.canReset = () => true;
