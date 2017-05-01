import React from 'react';

import { WithClusterConfig } from './ui';
import { validate } from '../validate';
import { ADMIN_EMAIL, ADMIN_PASSWORD } from '../cluster-config';

import { Input, Password } from './ui';

export const Users = () => {
  return (
    <div>
      <div className="form-group">
        These credentials will be used to log in to your Tectonic Console. Additional identity services and user management can be configured in the console.
      </div>
      <div className="row form-group">
        <div className="col-sm-3">
          <label htmlFor={ADMIN_EMAIL}>Email Address</label>
        </div>
        <div className="col-sm-9">
          <WithClusterConfig field={ADMIN_EMAIL} validator={validate.email}>
            <Input id={ADMIN_EMAIL} placeholder="admin@example.com" />
          </WithClusterConfig>
        </div>
      </div>
      <div className="row form-group">
        <div className="col-sm-3">
          <label htmlFor={ADMIN_PASSWORD}>Password</label>
        </div>
        <div className="col-sm-9">
          <WithClusterConfig field={ADMIN_PASSWORD} validator={validate.nonEmpty}>
            <Password id={ADMIN_PASSWORD} />
          </WithClusterConfig>
        </div>
      </div>
    </div>
  );
};

Users.canNavigateForward = ({clusterConfig}) => {
  return !validate.email(clusterConfig[ADMIN_EMAIL]) &&
         !validate.nonEmpty(clusterConfig[ADMIN_PASSWORD]);
};
