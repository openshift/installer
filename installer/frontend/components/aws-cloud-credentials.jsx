import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { compose, validate } from '../validate';
import { configActions } from '../actions';
import { Input, Password, Select, RadioBoolean, Connect } from './ui';
import { Alert } from './alert';

import { getRegions, getDefaultSubnets } from '../aws-actions';
import { Field, Form } from '../form';
import { TectonicGA } from '../tectonic-ga';

import {
  AWS_ACCESS_KEY_ID,
  AWS_CONTROLLER_SUBNETS,
  AWS_WORKER_SUBNETS,
  AWS_CONTROLLER_SUBNET_IDS,
  AWS_WORKER_SUBNET_IDS,
  AWS_REGION,
  AWS_REGION_FORM,
  AWS_SECRET_ACCESS_KEY,
  AWS_SESSION_TOKEN,
  STS_ENABLED,
} from '../cluster-config';

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
  'sa-east-1': 'São Paulo',
  'us-east-1': 'Northern Virginia',
  'us-east-2': 'Ohio',
  'us-gov-west-1': 'AWS GovCloud',
  'us-west-1': 'Northern California',
  'us-west-2': 'Oregon',
};

const { batchSetIn } = configActions;
const AWS_CREDS = 'AWSCreds';

const awsCredsForm = new Form(AWS_CREDS, [
  new Field(STS_ENABLED, {default: false}),
  new Field(AWS_ACCESS_KEY_ID, {
    default: '',
    validator: compose(validate.nonEmpty, (v) => {
      if (!v.startsWith('AK')) {
        return 'AWS access key IDs always begin with \'AK\'.';
      }
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
    default: '',
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
    default: '',
    validator: validate.nonEmpty,
    dependencies: [STS_ENABLED],
    ignoreWhen: cc => !cc[STS_ENABLED],
  }),
], {
  asyncValidator: (dispatch, getState, data, ignored, isNow) => {
    const creds = {
      AccessKeyID: data[AWS_ACCESS_KEY_ID],
      SecretAccessKey: data[AWS_SECRET_ACCESS_KEY],
      SessionToken: data[AWS_SESSION_TOKEN] || '',
      // you must send a region to get a region!!!
      Region: data[AWS_REGION] || 'us-east-1',
    };

    // This returns the list of regions, which we don't want to show as an error message.
    return dispatch(getRegions(null, creds, isNow)).then(() => {});
  },
});

const selectRegionForm = new Form(AWS_REGION_FORM, [
  awsCredsForm,
  new Field(AWS_REGION, {
    default: '',
    validator: validate.nonEmpty,
    dependencies: [AWS_CREDS],
  }),
], {
  asyncValidator: (dispatch, getState, data, oldData, isNow) => {
    const oldAws = oldData[AWS_REGION];
    const newAws = data[AWS_REGION];
    if (!oldAws || newAws === oldAws) {
      return Promise.resolve();
    }
    batchSetIn(dispatch, [
      [AWS_CONTROLLER_SUBNETS, {}],
      [AWS_CONTROLLER_SUBNET_IDS, {}],
      [AWS_WORKER_SUBNETS, {}],
      [AWS_WORKER_SUBNET_IDS, {}],
    ]);

    // TODO: (ggreer) move this to getextradata in another form
    return dispatch(getDefaultSubnets(null, null, isNow));
  },
});

const OPT_GROUPS = { ap: 'Asia Pacific', ca: 'Canada', cn: 'China', eu: 'European Union', sa: 'South America', us: 'United States'};

const stateToProps = ({aws, clusterConfig, serverFacts}) => {
  const regionSelections = {
    inFly: aws.availableRegions.inFly,
    value: [],
    error: aws.availableRegions.error,
  };

  // Calculate intersection of AWS regions: those w/coreos images & those the user has access to
  const coreosRegionsValues = _.map(serverFacts.awsRegions, v => v.value);
  if (!aws.availableRegions.inFly && !_.isEmpty(aws.availableRegions.value) && !_.isEmpty(coreosRegionsValues)) {
    const intersection = _.filter(aws.availableRegions.value, v => _.includes(coreosRegionsValues, v));
    // Format for use with <AsyncSelect/>
    regionSelections.value = intersection.sort().map(value => {
      const optgroupKey = value.split('-')[0];
      return {
        value,
        label: `${REGION_NAMES[value]} (${value})`,
        optgroup: OPT_GROUPS[optgroupKey] || optgroupKey,
      };
    });
  }

  return {
    stsEnabled: clusterConfig[STS_ENABLED],
    regionSelections,
  };
};

const awsCreds = <div>
  <div className="row form-group">
    <div className="col-xs-4">
      <label htmlFor="accessKeyId">Access Key ID</label>
    </div>
    <div className="col-xs-8">
      <Connect field={AWS_ACCESS_KEY_ID}>
        <Input id="accessKeyId" autoFocus="true" placeholder="AKxxxxxxxxxxxxxxxxxx" />
      </Connect>
    </div>
  </div>
  <div className="row form-group">
    <div className="col-xs-4">
      <label htmlFor="secretAccessKey">Secret Access Key</label>
    </div>
    <div className="col-xs-8">
      <Connect field={AWS_SECRET_ACCESS_KEY}>
        <Password id="secretAccessKey"/>
      </Connect>
    </div>
  </div>
</div>;

export const AWS_CloudCredentials = connect(stateToProps)(
({regionSelections, stsEnabled}) =>
  <div>
    <div className="row form-group">
      <div className="col-xs-12">
        Enter your Amazon Web Services (AWS) credentials to create and configure the required resources.
        It is strongly suggested that you create a <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#privileges" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} target="_blank">limited access role</a> for Tectonic's communication with your cloud provider.
      </div>
    </div>

    <div className="row form-group">
      <div className="col-xs-12">
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={STS_ENABLED}>
                <RadioBoolean inverted={true} name="stsEnabled" />
              </Connect>
              Use a normal access key
            </label>&nbsp;(default)
            <p className="text-muted">
              Go to the <a href="https://console.aws.amazon.com/iam/home#/users" target="_blank">AWS console user section</a>, select your user name, and the Security Credentials tab.
            </p>
          </div>
          <div className="wiz-radio-group__body">
            { !stsEnabled && awsCreds }
          </div>
        </div>
        <div className="wiz-radio-group">
          <div className="radio wiz-radio-group__radio">
            <label>
              <Connect field={STS_ENABLED}>
                <RadioBoolean name="stsEnabled" />
              </Connect>
              Use a temporary session token
            </label>
          </div>
          <div className="wiz-radio-group__body">
            { stsEnabled && <div>
              { awsCreds }
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
              <Alert>
                Temporary session tokens have a maximum lifetime of one hour. You must complete the Tectonic Installer before the token expires.
              </Alert>
            </div>}
          </div>
        </div>
      </div>
    </div>
    <hr />
    <div className="row form-group">
      <div className="col-xs-4">
        <label htmlFor="awsRegion">Region</label>
      </div>
      <div className="col-xs-8">
        <Connect field={AWS_REGION}>
          <Select id="awsRegion" availableValues={regionSelections} disabled={regionSelections.inFly}>
            <option value="" disabled>Please select region</option>
          </Select>
        </Connect>
      </div>
    </div>

    <div className="row form-group">
      <div className="col-xs-12">
        <awsCredsForm.Errors />
      </div>
    </div>
  </div>
);

AWS_CloudCredentials.canNavigateForward = selectRegionForm.canNavigateForward;
