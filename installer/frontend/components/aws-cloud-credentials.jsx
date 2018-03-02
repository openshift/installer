import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { compose, validate } from '../validate';
import { LoaderInline } from './loader';
import { A, DocsA, ExternalLinkIcon, Input, Password, Select, RadioBoolean, Connect } from './ui';
import { Alert } from './alert';

import { getRegions } from '../aws-actions';
import { Field, Form } from '../form';

import {
  AWS_ACCESS_KEY_ID,
  AWS_CREDS,
  AWS_REGION,
  AWS_REGION_FORM,
  AWS_SECRET_ACCESS_KEY,
  AWS_SESSION_TOKEN,
  STS_ENABLED,
} from '../cluster-config';

const awsCredsForm = new Form(AWS_CREDS, [
  new Field(STS_ENABLED, {default: false}),
  new Field(AWS_ACCESS_KEY_ID, {
    validator: compose(validate.nonEmpty, (v) => {
      if (v.indexOf('@') >= 0) {
        return 'AWS access key IDs are not email addresses.';
      }
      if (v.length < 20) {
        return 'AWS key IDs are at least 20 characters.';
      }
      if (v.trim() !== v) {
        return 'AWS key IDs cannot start or end with whitespace.';
      }
    }),
  }),
  new Field(AWS_SECRET_ACCESS_KEY, {
    validator: compose(validate.nonEmpty, (v) => {
      if (v.length < 40) {
        return 'AWS secrets are at least 40 characters.';
      }
      if (v.trim() !== v) {
        return 'AWS secrets cannot start or end with whitespace.';
      }
    }),
  }),
  new Field(AWS_SESSION_TOKEN, {
    validator: validate.nonEmpty,
    dependencies: [STS_ENABLED],
    ignoreWhen: cc => !cc[STS_ENABLED],
  }),
]);

const selectRegionForm = new Form(AWS_REGION_FORM, [
  awsCredsForm,
  new Field(AWS_REGION, {
    validator: validate.nonEmpty,
    dependencies: [AWS_CREDS],
    getExtraStuff: (dispatch, isNow) => dispatch(getRegions(isNow)),
  }),
]);

const REGION_NAMES = {
  'ap-northeast-1': 'Tokyo',
  'ap-northeast-2': 'Seoul',
  'ap-south-1': 'Mumbai',
  'ap-southeast-1': 'Singapore',
  'ap-southeast-2': 'Sydney',
  'ca-central-1': 'Canada',
  'cn-north-1': 'Beijing',
  'eu-central-1': 'Frankfurt',
  'eu-west-1': 'Ireland',
  'eu-west-2': 'London',
  'eu-west-3': 'Paris',
  'sa-east-1': 'São Paulo',
  'us-east-1': 'Northern Virginia',
  'us-east-2': 'Ohio',
  'us-gov-west-1': 'AWS GovCloud',
  'us-west-1': 'Northern California',
  'us-west-2': 'Oregon',
};

const OPT_GROUPS = {
  ap: 'Asia Pacific',
  ca: 'Canada',
  cn: 'China',
  eu: 'European Union',
  sa: 'South America',
  us: 'United States',
};

const stateToProps = ({aws, serverFacts}) => {
  const {error, inFly} = aws.availableRegions;

  // Calculate intersection of AWS regions: those w/coreos images & those the user has access to
  let options = [];
  if (!inFly && !_.isEmpty(aws.availableRegions.value) && !_.isEmpty(serverFacts.awsRegions)) {
    const intersection = _.intersection(aws.availableRegions.value, serverFacts.awsRegions);
    // Format for use with <Select />
    options = intersection.sort().map(value => {
      const optgroupKey = value.split('-')[0];
      return {
        value,
        label: REGION_NAMES[value] ? `${REGION_NAMES[value]} (${value})` : value,
        optgroup: OPT_GROUPS[optgroupKey] || optgroupKey,
      };
    });
  }

  return {error: _.get(error, 'message') || error, inFly, options};
};

const Region = connect(stateToProps)(({error, inFly, options}) => {
  if (error && !inFly) {
    return <div className="row form-group">
      <div className="col-xs-12">
        <Alert severity="error">{error}</Alert>
      </div>
    </div>;
  }
  return <div className="row form-group">
    <div className="col-xs-4">
      <label htmlFor="awsRegion">Region</label>
    </div>
    <div className="col-xs-8">
      {inFly ? <LoaderInline /> : <Connect field={AWS_REGION}>
        <Select id="awsRegion" options={options}>
          <option value="" disabled>Please select region</option>
        </Select>
      </Connect>}
    </div>
  </div>;
});

const awsCreds = <div>
  <div className="row form-group">
    <div className="col-xs-4">
      <label htmlFor="accessKeyId">Access Key ID</label>
    </div>
    <div className="col-xs-8">
      <Connect field={AWS_ACCESS_KEY_ID}>
        <Input id="accessKeyId" autoFocus={true} placeholder="AKxxxxxxxxxxxxxxxxxx" />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-4">
      <label htmlFor="secretAccessKey">Secret Access Key</label>
    </div>
    <div className="col-xs-8">
      <Connect field={AWS_SECRET_ACCESS_KEY}>
        <Password id="secretAccessKey" />
      </Connect>
    </div>
  </div>
</div>;

export const AWS_CloudCredentials = connect(
  ({clusterConfig}) => ({stsEnabled: clusterConfig[STS_ENABLED]})
)(({stsEnabled}) => <div>
  <div className="row form-group">
    <div className="col-xs-12">
      Enter your Amazon Web Services (AWS) credentials to create and configure the required resources. It is strongly suggested that you create a <DocsA path="/install/aws/requirements.html#privileges">limited access role</DocsA> for Tectonic's communication with your cloud provider.
    </div>
  </div>

  <div className="row form-group">
    <div className="col-xs-12">
      <div className="wiz-radio-group">
        <div className="radio wiz-radio-group__radio">
          <label>
            <Connect field={STS_ENABLED}>
              <RadioBoolean inverted={true} name="stsEnabled" id="stsEnabledFalse" />
            </Connect>
            Use a normal access key
          </label>&nbsp;(default)
          <p className="text-muted">
            Go to the <A href="https://console.aws.amazon.com/iam/home#/users">AWS console user section<ExternalLinkIcon /></A>, select your user name, and the Security Credentials tab.
          </p>
        </div>
        <div className="wiz-radio-group__body">
          {!stsEnabled && awsCreds}
        </div>
      </div>
      <div className="wiz-radio-group">
        <div className="radio wiz-radio-group__radio">
          <label>
            <Connect field={STS_ENABLED}>
              <RadioBoolean name="stsEnabled" id="stsEnabledTrue" />
            </Connect>
            Use a temporary session token
          </label>
        </div>
        <div className="wiz-radio-group__body">
          {stsEnabled && <div>
            {awsCreds}
            <div className="row form-group">
              <div className="col-xs-4">
                <label htmlFor={AWS_SESSION_TOKEN}>Session Token</label>
              </div>
              <div className="col-xs-8">
                <Connect field={AWS_SESSION_TOKEN}>
                  <Input id={AWS_SESSION_TOKEN} />
                </Connect>
              </div>
            </div>
            <Alert>Temporary session tokens have a maximum lifetime of one hour. You must complete the Tectonic Installer before the token expires.</Alert>
          </div>}
        </div>
      </div>
    </div>
  </div>
  <hr />
  <Region />
</div>);

AWS_CloudCredentials.canNavigateForward = selectRegionForm.canNavigateForward;
