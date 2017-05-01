import React from 'react';
import { connect } from 'react-redux';

import { configActionTypes, sequenceActionTypes } from '../actions';
import { BARE_METAL_TF } from '../platforms';
import { validate } from '../validate';
import { PLATFORM_TYPE, SSH_AUTHORIZED_KEYS } from '../cluster-config';

import { FileArea } from './ui';

const keyPlaceholder = `Paste your SSH public key here. It's often found in ~/.ssh/id_rsa.pub.

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDxL
...
you@example.com
`;

export const BM_SSHKeys = connect(
  ({clusterConfig, sequence}) => {
    let keys = clusterConfig[SSH_AUTHORIZED_KEYS];
    const platformType = clusterConfig[PLATFORM_TYPE];
    if (keys.length > 1 && platformType === BARE_METAL_TF) {
      // BARE_METAL_TF only supports 1 ssh key :(
      keys = clusterConfig[SSH_AUTHORIZED_KEYS][0];
    }
    return {
      keys,
      platformType,
      sequence,
    };
  },
  (dispatch) => {
    return {
      handleKey: (key, id, index) => {
        dispatch({
          type: configActionTypes.SET_SSH_AUTHORIZED_KEYS,
          payload: {
            index: index,
            value: {id, key},
          },
        });
      },
      addKey: (index, sequence) => {
        dispatch({
          type: sequenceActionTypes.INCREMENT,
        });
        dispatch({
          type: configActionTypes.SET_SSH_AUTHORIZED_KEYS,
          payload: {
            index: index,
            value: {
              // sequenceActionTypes.INCREMENT was dispatched, but we still have the old value
              id: `ssh-key-${sequence + 1}`,
              key: '',
            },
          },
        });
      },
      removeKey: (index) => {
        dispatch({
          type: configActionTypes.REMOVE_SSH_AUTHORIZED_KEYS,
          payload: index,
        });
      },
    };
  }
)(({addKey, handleKey, removeKey, keys, platformType, sequence}) => {
  const fields = keys.map(({key, id}, i) => {
    return (
      <div className="row form-group" key={id}>
        <div className="col-xs-3">
          <label htmlFor="ssh-key">Public Key</label>
        </div>
        <div className="col-xs-9">
          {
            keys.length > 1 &&
            <div className="fa fa-times-circle wiz-teeny-close-button wiz-teeny-close-button--ssh"
                 onClick={() => removeKey(i)}></div>
          }
          <div className="wiz-ssh-key-container">
            <FileArea
                className="wiz-ssh-key-container__input"
                id={id}
                value={key}
                data-index={i}
                invalid={validate.SSHKey(key)}
                onValue={v => handleKey(v, id, i)}
                placeholder={keyPlaceholder}
                autoFocus={i > 0} />
            {
              i + 1 === keys.length && platformType !== BARE_METAL_TF &&
              <p>
                <a onClick={() => addKey(keys.length, sequence)}>
                  <span className="fa fa-plus"></span> Add another public key
                </a>
              </p>
            }
          </div>
        </div>
      </div>
    );
  });

  return (
    <div>
      <div className="form-group">Access to the nodes is intended for use by admins.
        End users can run applications using the API, CLI and Tectonic
        Console and typically don’t need SSH access.</div>
      <div className="form-group">The public keys below will be added to all machines in this cluster.</div>
      {fields}
    </div>
  );
});

BM_SSHKeys.canNavigateForward = ({clusterConfig}) => {
  const ks = clusterConfig[SSH_AUTHORIZED_KEYS].map(k => k.key);
  return ks.length && ks.every(v => !validate.SSHKey(v));
};
