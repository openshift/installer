import React from 'react';

import { validate } from '../validate';
import { ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_PASSWORD2, CREDS } from '../cluster-config';
import { Field, Form } from '../form';
import { Input, Password, Connect } from './ui';

const credsForm = new Form(CREDS, [
  new Field(ADMIN_EMAIL, {default: '', validator: validate.email}),
  new Field(ADMIN_PASSWORD, {default: '', validator: validate.nonEmpty}),
  new Field(ADMIN_PASSWORD2, {default: '', validator: validate.nonEmpty}),
], {
  validator: (data, cc) => {
    if (cc[ADMIN_PASSWORD] && cc[ADMIN_PASSWORD2] && cc[ADMIN_PASSWORD] !== cc[ADMIN_PASSWORD2]) {
      return 'Passwords do not match.';
    }
  },
});

export const Users = () => <div>
  <div className="form-group">
    These credentials will be used to log in to your Tectonic Console. Additional identity services and user management can be configured in the console.
  </div>
  <div className="row form-group">
    <div className="col-sm-4">
      <label htmlFor={ADMIN_EMAIL}>Email Address</label>
    </div>
    <div className="col-sm-8">
      <Connect field={ADMIN_EMAIL}>
        <Input autoFocus={true} type="email" placeholder="admin@example.com" />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-sm-4">
      <label htmlFor={ADMIN_PASSWORD}>Password</label>
    </div>
    <div className="col-sm-8">
      <Connect field={ADMIN_PASSWORD}>
        <Password />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-sm-4">
      <label htmlFor={ADMIN_PASSWORD2}>Confirm Password</label>
    </div>
    <div className="col-sm-8">
      <Connect field={ADMIN_PASSWORD2}>
        <Password />
      </Connect>
    </div>
  </div>
  <credsForm.Errors />
</div>;

Users.canNavigateForward = credsForm.canNavigateForward;
