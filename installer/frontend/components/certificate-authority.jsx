import React from 'react';
import { connect } from 'react-redux';

import { CA_CERTIFICATE, CA_PRIVATE_KEY, CA_TYPE, CA_TYPES } from '../cluster-config';
import { Field, Form } from '../form';
import { validate } from '../validate';
import { CertArea, Connect, PrivateKeyArea, Radio } from './ui';

const form = new Form('CERTIFICATE_AUTHORITY', [
  new Field(CA_TYPE, {default: CA_TYPES.SELF_SIGNED}),
  new Field(CA_CERTIFICATE, {
    default: '',
    ignoreWhen: cc => cc[CA_TYPE] !== CA_TYPES.OWNED,
  }),
  new Field(CA_PRIVATE_KEY, {
    default: '',
    ignoreWhen: cc => cc[CA_TYPE] !== CA_TYPES.OWNED,
  }),
], {
  validator: (data, cc) => cc[CA_TYPE] === CA_TYPES.SELF_SIGNED ||
    validate.certificate(cc[CA_CERTIFICATE]) ||
    validate.privateKey(cc[CA_PRIVATE_KEY]),
});

export const CertificateAuthority = connect(
  ({clusterConfig}) => ({caType: clusterConfig[CA_TYPE]})
)(({caType}) => <div>
  <div className="row form-group">
    <div className="col-xs-12">
      A certificate authority (CA) is needed so we can generate certificates for cluster components.
    </div>
  </div>

  <div className="row form-group">
    <div className="col-xs-12">
      <div className="wiz-radio-group">
        <div className="radio wiz-radio-group__radio">
          <label>
            <Connect field={CA_TYPE}>
              <Radio name="certificateAuthority" id="certificateAuthoritySelfSigned" value={CA_TYPES.SELF_SIGNED} />
            </Connect>
            Generate a CA certificate and key for me.
          </label>&nbsp;(default)
          <p className="text-muted wiz-help-text">Component certificates will not be trusted in web browsers without additional configuration.</p>
        </div>
      </div>

      <div className="wiz-radio-group">
        <div className="radio wiz-radio-group__radio">
          <label>
            <Connect field={CA_TYPE}>
              <Radio name="certificateAuthority" id="certificateAuthorityOwned" value={CA_TYPES.OWNED} />
            </Connect>
            I'll provide a CA certificate and key in PEM format.
          </label>
          <p className="text-muted wiz-help-text">Your CA will be used to issue certificates for cluster components.</p>
        </div>
        {caType === CA_TYPES.OWNED && <div className="wiz-radio-group__body">
          <div className="row form-group">
            <div className="col-xs-12">
              <Connect field={CA_CERTIFICATE}>
                <CertArea autoFocus={true} uploadButtonLabel="Upload CA Certificate" />
              </Connect>
            </div>
          </div>
          <div className="row form-group">
            <div className="col-xs-12">
              <Connect field={CA_PRIVATE_KEY}>
                <PrivateKeyArea uploadButtonLabel="Upload CA Private Key" />
              </Connect>
            </div>
          </div>
        </div>
        }
      </div>
    </div>
  </div>
</div>);

CertificateAuthority.canNavigateForward = form.canNavigateForward;
