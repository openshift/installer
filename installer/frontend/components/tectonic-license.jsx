import React from 'react';
import jwt_decode from 'jwt-decode';

import { validate } from '../validate';
import { FileArea, Connect } from './ui';
import { PULL_SECRET, TECTONIC_LICENSE, LICENSING } from '../cluster-config';
import { Field, Form } from '../form';

const TECTONIC_LICENSE_PLACEHOLDER = `Raw formatted license:

06jd9rqTAr1DZWs/ssB0k128C1nfcq0v4yqL4PpDLXg...`;

const PULL_SECRET_PLACEHOLDER = `{
  "auths": {
    "quay.io": {
      "auth": "Y29yZW9zK3RlY190ZXN0aW5nOmZha2VwYXNzd29yZAo=",
      "email": "user@example.com"
    }
  }
}...`;

const TectonicLicenseField = new Field(TECTONIC_LICENSE, {
  default: '',
  validator: token => {
    const err = validate.nonEmpty(token);
    if (err) {
      return err;
    }
    try {
      const decoded = jwt_decode(token, {header: false});
      if (!decoded.license) {
        return 'Error parsing license.';
      }
    } catch (unused) {
      return 'Error parsing license. Please make sure you pasted the "raw format" license from account.coreos.com.';
    }
    return;
  },
});

const PullSecretField = new Field(PULL_SECRET, {
  default: '',
  validator: secret => {
    const err = validate.nonEmpty(secret);
    if (err) {
      return err;
    }
    try {
      JSON.parse(secret);
    } catch (unused) {
      return 'Pull secret must be valid JSON.';
    }
    return;
  },
});

export const licenseForm = new Form(LICENSING, [TectonicLicenseField, PullSecretField]);

export const TectonicLicense = () =>
  <div>
    <div className="row form-group">
      <div className="col-xs-3">
        <label htmlFor="tectonicLicense">CoreOS License</label>
      </div>
      <div className="col-xs-9">
        <Connect field={TECTONIC_LICENSE}>
          <FileArea id="tectonicLicense" placeholder={TECTONIC_LICENSE_PLACEHOLDER} uploadButtonLabel='Upload "tectonic-license.txt"' />
        </Connect>
        {/* eslint-disable react/jsx-no-target-blank */}
        <p className="text-muted">Input "CoreOS License" from <a href="https://account.coreos.com" rel="noopener" target="_blank">account.coreos.com</a>.</p>
        {/* eslint-enable react/jsx-no-target-blank */}
      </div>
    </div>

    <div className="row form-group">
      <div className="col-xs-3">
        <label htmlFor={PULL_SECRET}>Pull Secret</label>
      </div>
      <div className="col-xs-9">
        <Connect field={PULL_SECRET}>
          <FileArea id={PULL_SECRET} placeholder={PULL_SECRET_PLACEHOLDER} uploadButtonLabel='Upload "config.json"' />
        </Connect>
        {/* eslint-disable react/jsx-no-target-blank */}
        <p className="text-muted">Input "Pull Secret" from <a href="https://account.coreos.com" rel="noopener" target="_blank">account.coreos.com</a>.</p>
        {/* eslint-enable react/jsx-no-target-blank */}
      </div>
    </div>

    <licenseForm.Errors />

  </div>;
