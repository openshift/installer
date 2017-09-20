import React from 'react';
import { connect } from 'react-redux';

import { configActionTypes } from '../actions';
import { CA_TYPE, CA_CERTIFICATE, CA_PRIVATE_KEY } from '../cluster-config';
import { Field, Form } from '../form';
import { validate } from '../validate';
import { CertArea, Connect, PrivateKeyArea } from './ui';

const form = new Form('CERTIFICATE_AUTHORITY', [
  new Field(CA_CERTIFICATE, {default: '', validator: validate.certificate}),
  new Field(CA_PRIVATE_KEY, {default: '', validator: validate.privateKey}),
]);

export const CertificateAuthority = connect(
  ({clusterConfig}) => {
    return {
      caType: clusterConfig[CA_TYPE],
    };
  },

  (dispatch) => {
    return {
      setCAType: (value) => {
        dispatch({
          type: configActionTypes.SET,
          payload: {[CA_TYPE]: value},
        });
      },
    };
  }
)(({caType, setCAType}) => {
  // TODO: (ggreer) use checkbox from ui
  return (
    <div>
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
                <input
                  type="radio"
                  name="certificateAuthority"
                  defaultChecked={caType === 'self-signed'}
                  onChange={() => setCAType('self-signed')} />
                Generate a CA certificate and key for me.
              </label>&nbsp;(default)
              <p className="text-muted wiz-help-text">Component certificates will not be trusted in web browsers without additional configuration.</p>
            </div>
          </div>
          <div className="wiz-radio-group">
            <div className="radio wiz-radio-group__radio">
              <label>
                <input
                  type="radio"
                  name="certificateAuthority"
                  defaultChecked={caType === 'owned'}
                  onChange={() => setCAType('owned')} />
                I'll provide a CA certificate and key in PEM format.
              </label>
              <p className="text-muted wiz-help-text">Your CA will be used to issue certificates for cluster components.</p>
            </div>
            <div className="wiz-radio-group__body">
              {
                caType === 'owned' && <div>
                  <div className="row form-group">
                    <div className="col-xs-12">
                      <Connect field={CA_CERTIFICATE}>
                        <CertArea autoFocus="true" uploadButtonLabel="Upload CA Certificate" />
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
      </div>
    </div>
  );
});

CertificateAuthority.canNavigateForward = state => state.clusterConfig[CA_TYPE] === 'self-signed' || form.canNavigateForward(state);
