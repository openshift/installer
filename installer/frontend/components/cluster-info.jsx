import _ from 'lodash';
import React from 'react';
import jwt_decode from 'jwt-decode';

import { validate } from '../validate';
import { CLUSTER_NAME, PULL_SECRET, TECTONIC_LICENSE, LICENSING } from '../cluster-config';
import fields from '../fields';
import { Field, Form } from '../form';
import { readFile } from '../readfile';

import { Alert } from './alert';
import { Connect, Input } from './ui';

// eslint-disable-next-line react/jsx-no-target-blank
const accountLink = <a href="https://account.coreos.com" rel="noopener" target="_blank">account.coreos.com</a>;

const licenseField = new Field(TECTONIC_LICENSE, {
  default: '',
  validator: token => {
    const err = validate.nonEmpty(token);
    if (err) {
      return err;
    }
    try {
      const decoded = jwt_decode(token, {header: false});
      if (!decoded.license) {
        return 'Error parsing license';
      }
    } catch (unused) {
      return 'Error parsing license';
    }
    return;
  },
});

const pullSecretField = new Field(PULL_SECRET, {
  default: '',
  // eslint-disable-next-line react/display-name
  validator: secret => {
    const err = validate.nonEmpty(secret);
    if (err) {
      return err;
    }
    try {
      JSON.parse(secret);
    } catch (unused) {
      return 'Error parsing pull secret';
    }
    return;
  },
});

const form = new Form(LICENSING, [fields[CLUSTER_NAME], licenseField, pullSecretField]);

const FileInput = ({id, onValue}) => {
  const upload = e => {
    readFile(e.target.files.item(0))
      .then(onValue)
      .catch(msg => console.error(msg));

    // Reset value so that onChange fires if you pick the same file again.
    e.target.value = null;
  };
  return <input type="file" id={id} onChange={upload} style={{display: 'none'}} />;
};

const FileUpload = ({buttonTitle, description, ErrorHelp, field, id, onValue, value}) => {
  const invalid = field.validator(value);
  if (invalid) {
    return <div style={{marginTop: 8}}>
      <i className="fa fa-ban wiz-error-fg"></i>&nbsp;&nbsp;Invalid {description}
      <Alert noIcon severity="error">
        <b>{invalid}</b>
        {ErrorHelp}
        <label className="btn btn-flat btn-warning">
          {buttonTitle}
          <FileInput id={id} onValue={onValue} />
        </label>
      </Alert>
    </div>;
  }
  return <div style={{marginTop: 8}}>
    <i className="fa fa-check-circle wiz-success-fg"></i>&nbsp;&nbsp;Valid {description}
    <label style={{fontSize: 14, margin: '0 0 0 15px'}}>
      <a>Edit</a>
      <FileInput id={id} onValue={onValue} />
    </label>
    <p className="text-muted">Your "{_.startCase(description)}" from {accountLink} has been included.</p>
  </div>;
};

const License = () => <Connect field={TECTONIC_LICENSE}>
  <FileUpload
    buttonTitle={'Upload "coreos-license.txt"'}
    description="CoreOS license"
    ErrorHelp={<p>Please make sure you upload the "raw format" license from {accountLink}.</p>}
    field={licenseField}
  />
</Connect>;

const Secret = () => <Connect field={PULL_SECRET}>
  <FileUpload
    buttonTitle={'Upload "config.json"'}
    description="pull secret"
    ErrorHelp={<p>Please make sure you upload the pull secret from {accountLink} in a valid JSON format.</p>}
    field={pullSecretField}
  />
</Connect>;

export const ClusterInfo = () => <div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={CLUSTER_NAME}>Cluster Name</label>
    </div>
    <div className="col-xs-9">
      <Connect field={CLUSTER_NAME}>
        <Input placeholder="production" autoFocus={true} />
      </Connect>
      <p className="text-muted">Give this cluster a name that will help you identify it.</p>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={TECTONIC_LICENSE}>CoreOS License</label>
    </div>
    <div className="col-xs-9">
      <License />
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={PULL_SECRET}>Pull Secret</label>
    </div>
    <div className="col-xs-9">
      <Secret />
    </div>
  </div>
</div>;

ClusterInfo.canNavigateForward = form.canNavigateForward;
