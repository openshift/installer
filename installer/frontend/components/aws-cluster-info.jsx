import { connect } from 'react-redux';
import React from 'react';

import { Input, Select, Connect } from './ui';
import { TectonicLicense, licenseForm } from './tectonic-license';
import { ExperimentalFeatures } from './experimental-features';
import { AWS_Tags, tagsFields } from './aws-tags';
import { Alert } from './alert';
import { Form, Field } from '../form';
import fields from '../fields';
import { AWS_CLUSTER_INFO, CHANNEL_TO_USE, CLUSTER_NAME } from '../cluster-config';

const clusterInfoForm = new Form(AWS_CLUSTER_INFO, [
  licenseForm,
  tagsFields,
  fields[CLUSTER_NAME],
  new Field(CHANNEL_TO_USE, {
    default: 'stable',
    validator: v => ['stable', 'beta', 'alpha'].includes(v) ? '' : 'unknown channel',
  }),
]);

const ChannelWarning = connect(({clusterConfig}) => ({channel: clusterConfig[CHANNEL_TO_USE]}))(
  ({channel}) => (channel !== 'alpha' ? null : <Alert>
    The Alpha channel can be unstable. It should only be used for testing or development.
  </Alert>)
);

const TagsWithPlaceholder = connect(({clusterConfig}) => ({clusterName: clusterConfig[CLUSTER_NAME] || 'myclustername'}))(({clusterName}) => <AWS_Tags placeholder={`e.g. ${clusterName}`} />);

export const AWS_ClusterInfo = () => <div>
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor={CLUSTER_NAME}>Name</label>
    </div>
    <div className="col-xs-6">
      <Connect field={CLUSTER_NAME}>
        <Input placeholder="production" autoFocus="true" />
      </Connect>
      <p className="text-muted" style={{marginBottom: 0}}>
        This name is used in the Tectonic Console to identify this cluster.
      </p>
    </div>
  </div>

  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="coreOsChannel">Container Linux Update Channel</label>
    </div>
    <div className="col-xs-6">
      <Connect field={CHANNEL_TO_USE}>
        <Select id="coreOsChannel">
          <option value="" disabled>Select a CoreOS Container Linux channel</option>
          <option value="stable">Stable</option>
          <option value="beta">Beta</option>
          <option value="alpha">Alpha</option>
        </Select>
      </Connect>
      <ChannelWarning/>
    </div>
  </div>

  <ExperimentalFeatures />

  <TectonicLicense />

  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="tags">Tags</label>
    </div>
    <div className="col-xs-9">
      <TagsWithPlaceholder />
    </div>
  </div>

  <clusterInfoForm.Errors />
</div>;

AWS_ClusterInfo.canNavigateForward = clusterInfoForm.canNavigateForward;
