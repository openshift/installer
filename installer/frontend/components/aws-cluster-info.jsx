import { connect } from 'react-redux';
import React from 'react';

import { Input, Connect } from './ui';
import { TectonicLicense, licenseForm } from './tectonic-license';
import { AWS_Tags, tagsFields } from './aws-tags';
import { Form } from '../form';
import fields from '../fields';
import { AWS_CLUSTER_INFO, CLUSTER_NAME } from '../cluster-config';

const clusterInfoForm = new Form(AWS_CLUSTER_INFO, [
  licenseForm,
  tagsFields,
  fields[CLUSTER_NAME],
]);

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

  <TectonicLicense />

  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="tags">AWS Tags</label>
    </div>
    <div className="col-xs-9">
      <TagsWithPlaceholder />
    </div>
  </div>

  <clusterInfoForm.Errors />
</div>;

AWS_ClusterInfo.canNavigateForward = clusterInfoForm.canNavigateForward;
