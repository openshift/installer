import React from 'react';

import { TectonicGA } from '../tectonic-ga';

export const DryRun = () => <div className="row">
  <div className="col-xs-12">
    <div className="form-group">
      {/* eslint-disable react/jsx-no-target-blank */}
      Your cluster assets have been created. You can download these <a href="https://coreos.com/tectonic/docs/latest/admin/assets-zip.html" rel="noopener" target="_blank">assets</a> and customize underlying infrastructure as needed.
      Note: changes to Kubernetes manifests or Tectonic components run in the cluster are not supported.
      &nbsp;<a href="https://coreos.com/tectonic/docs/latest/install/aws/manual-boot.html"
        onClick={TectonicGA.sendDocsEvent} rel="noopener" target="_blank">
        Read more here.&nbsp;&nbsp;<i className="fa fa-external-link" />
      </a>
      {/* eslint-enable react/jsx-no-target-blank */}
    </div>
    <div className="from-group">
      <div className="wiz-giant-button-container">
        <a href="/terraform/assets" download>
          <button className="btn btn-primary wiz-giant-button">
            <i className="fa fa-download"></i>&nbsp;&nbsp;Download assets
          </button>
        </a>
      </div>
    </div>
  </div>
</div>;

DryRun.canNavigateForward = () => false;
