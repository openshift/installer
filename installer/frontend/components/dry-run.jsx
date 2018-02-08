import React from 'react';

import { DocsA } from './ui';

export const DryRun = () => <div className="row">
  <div className="col-xs-12">
    <div className="form-group">
      Your cluster assets have been created. You can download these <DocsA path="/admin/assets-zip.html">assets</DocsA> and customize underlying infrastructure as needed.
      Note: changes to Kubernetes manifests or Tectonic components run in the cluster are not supported.&nbsp; <DocsA path="/install/aws/manual-boot.html">Read more here.</DocsA>
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
