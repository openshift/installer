import React from 'react';

import { validate, compose } from '../validate';
import { Connect, Input } from './ui';

import { Form, Field } from '../form';
import { CONTROLLER_DOMAIN, BM_TECTONIC_DOMAIN } from '../cluster-config';

const validator = compose(validate.nonEmpty, validate.domainName);

const hostNamesForm = new Form('BM_Hostname', [
  new Field(CONTROLLER_DOMAIN, { validator, default: '' }),
  new Field(BM_TECTONIC_DOMAIN, { validator, default: '' }),
]);

export const BM_Hostname = () =>
  <div>
    <div className="form-group">
      Enter a DNS name which resolves to any master node. This is the name we'll use when configuring components.
    </div>
    <div className="row form-group">
      <div className="col-xs-3">
        <label htmlFor={CONTROLLER_DOMAIN}>Master DNS</label>
      </div>
      <div className="col-xs-9">
        <Connect field={CONTROLLER_DOMAIN}>
          <Input autoFocus={true} id={CONTROLLER_DOMAIN} placeholder="masters.example.com" />
        </Connect>
      </div>
    </div>
    <div className="form-group">Enter a DNS name which resolves to any worker node(s) for Ingress.</div>
    <div className="row form-group">
      <div className="col-xs-3">
        <label htmlFor={BM_TECTONIC_DOMAIN}>Tectonic DNS</label>
      </div>
      <div className="col-xs-9">
        <Connect field={BM_TECTONIC_DOMAIN}>
          <Input placeholder="tectonic.example.com" />
        </Connect>
      </div>
    </div>
  </div>;

BM_Hostname.canNavigateForward = hostNamesForm.canNavigateForward;
