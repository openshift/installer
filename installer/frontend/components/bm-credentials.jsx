import React from 'react';
import {
  BM_MATCHBOX_CA,
  BM_MATCHBOX_CLIENT_CERT,
  BM_MATCHBOX_CLIENT_KEY,
} from '../cluster-config';
import { Field, Form } from '../form';
import { validate } from '../validate';
import { CertArea, Connect, PrivateKeyArea } from './ui';

const form = new Form('BM_MATCHBOX_CREDENTIALS', [
  new Field(BM_MATCHBOX_CA, {default: '', validator: validate.certificate}),
  new Field(BM_MATCHBOX_CLIENT_CERT, {default: '', validator: validate.certificate}),
  new Field(BM_MATCHBOX_CLIENT_KEY, {default: '', validator: validate.privateKey}),
]);

export const BM_Credentials = () => <div>
  <div className="form-group">
    Credentials were generated during the matchbox installation. Provide the certificates and keys here:
  </div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={BM_MATCHBOX_CA}>CA Certificate</label>
    </div>
    <div className="col-xs-9">
      <Connect field={BM_MATCHBOX_CA}>
        <CertArea autoFocus="true" />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={BM_MATCHBOX_CLIENT_CERT}>Client Certificate</label>
    </div>
    <div className="col-xs-9">
      <Connect field={BM_MATCHBOX_CLIENT_CERT}>
        <CertArea />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={BM_MATCHBOX_CLIENT_KEY}>Client Key</label>
    </div>
    <div className="col-xs-9">
      <Connect field={BM_MATCHBOX_CLIENT_KEY}>
        <PrivateKeyArea />
      </Connect>
    </div>
  </div>
</div>;

BM_Credentials.canNavigateForward = form.canNavigateForward;
