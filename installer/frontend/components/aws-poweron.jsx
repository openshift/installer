import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { SaveAssets } from './save-assets';
import { WaitingLi } from './ui';
import { SSHInstructions } from './ssh-instructions';
import { AWS_DomainValidation } from './aws-domain-validation';

const componentDescriptions = {
  cloudFormation: 'Starting AWS CloudFormation',
  kubelet: 'Starting Kubelet',
  tectonicConsole: 'Starting Tectonic components',
};

const cfCategories = {
  vpc: 'VPC',
  securityGroups: 'Security Groups',
  networking: 'Networking',
  iam: 'IAM',
  ec2: 'EC2 Instances',
  autoscale: 'Autoscale Groups',
  loadBalancer: 'Load Balancer',
  other: 'Other VPC Resources',
  dns: 'DNS',
};

const cfStatusMapping = {
  'AWS::AutoScaling::AutoScalingGroup': cfCategories.autoscale,
  'AWS::AutoScaling::LaunchConfiguration': cfCategories.autoscale,
  'AWS::CloudWatch::Alarm': cfCategories.other,
  'AWS::EC2::EIP': cfCategories.networking,
  'AWS::EC2::Instance': cfCategories.ec2,
  'AWS::EC2::InternetGateway': cfCategories.networking,
  'AWS::EC2::Route': cfCategories.networking,
  'AWS::EC2::RouteTable': cfCategories.networking,
  'AWS::EC2::SecurityGroup': cfCategories.securityGroups,
  'AWS::EC2::SecurityGroupIngress': cfCategories.securityGroups,
  'AWS::EC2::Subnet': cfCategories.networking,
  'AWS::EC2::SubnetRouteTableAssociation': cfCategories.networking,
  'AWS::EC2::VPC': cfCategories.vpc,
  'AWS::EC2::VPCGatewayAttachment': cfCategories.vpc,
  'AWS::ElasticLoadBalancing::LoadBalancer': cfCategories.loadBalancer,
  'AWS::IAM::InstanceProfile': cfCategories.iam,
  'AWS::IAM::Role': cfCategories.iam,
  'AWS::Route53::RecordSet': cfCategories.dns,
};

const toCategory = rt => {
  if (cfStatusMapping[rt]) {
    return cfStatusMapping[rt];
  }
  return cfCategories.other;
};

const stateToProps = ({cluster, commitState, clusterConfig}) => ({
  error: cluster.error,
  status: cluster.status,
  blob: commitState.response,
  region: clusterConfig.awsRegion,
});

const ERROR_STATUSES = [
  'CREATE_FAILED',
  'DELETE_COMPLETE',
  'DELETE_IN_PROGRESS',
  'DELETED',
];

const SUCCESS = 'CREATE_COMPLETE';

const CloudFormationAnchor = ({arn, region, children}) => {
  if (!arn || !region) {
    return <span>{children}</span>;
  }
  const href = `https://${region}.console.aws.amazon.com/cloudformation/home?region=${region}#/stack/detail?stackId=${arn.replace(/\//g, '%2F')}`;
  return <a href={href} target="_blank">{children}</a>;
};


export const AWS_PowerOn = connect(stateToProps)(
class PowerOn extends React.Component {
  constructor (props) {
    super(props);
    this.state = {remoteAddr: null};
  }

  componentWillReceiveProps ({status}) {
    if (!status) {
      return;
    }
    if (!status.cloudFormation) {
      return;
    }
    if (!status.kubelet || !status.kubelet.ready) {
      return;
    }
    if (_.includes(status.kubelet.addrs, this.state.remoteAddr)) {
      // Only change remoteAddr previously-checked remoteAddr is no longer in DNS.
      return;
    }
    // Cloud formation & kublet are ready. OK to tell user to run scp/ssh commands.
    this.setState({
      remoteAddr: status.kubelet.remoteAddr,
    });
  }

  render () {
    const {error, status, blob} = this.props;
    const substeps = {};
    if (status) {
      if (status.cloudFormation && !status.cloudFormation.ready) {
        const cfStatus = {};
        _.each(cfCategories, (category) => {
          cfStatus[category] = [];
        });
        status.cloudFormation.resources.forEach(r => {
          const category = toCategory(r.resourceType);
          cfStatus[category].push(r);
        });
        substeps.cloudFormation = _.map(cfStatus, (resources, category) => {
          const done = _.size(resources) > 0 && _.every(resources, r => r.resourceStatus === SUCCESS);
          const errResource = _.find(resources, r => _.includes(ERROR_STATUSES, r.resourceStatus));
          let errMsg;
          if (errResource) {
            errMsg = 'Failure';
            if (errResource.resourceStatusReason) {
              errMsg = `Failure: ${errResource.resourceStatusReason}`;
            }
          }
          return (
            <WaitingLi done={done} error={!!errResource} key={category} substep={true}>
              <span title={category}>{category} {errMsg}</span>
            </WaitingLi>
            );
        });
      } else if (status.kubelet && !status.kubelet.ready) {
        substeps.kubelet = [];
        let msg = `Resolving ${status.kubelet.instance}`;
        const dnsReady = (status.kubelet.message || '').search('no such host') === -1;
        substeps.kubelet.push(
          <WaitingLi done={dnsReady} key="dns" substep={true}>
            <span title={msg}>{msg}</span>
          </WaitingLi>
        );
        msg = `Waiting for kubelet to be ready`;
        substeps.kubelet.push(
          <WaitingLi done={status.kubelet.ready} key="port" substep={true}>
            <span title={msg}>{msg}</span>
          </WaitingLi>
        );
        substeps.kubelet.push(<AWS_DomainValidation key="domain" />);
      }
    }

    const statusList = (['cloudFormation', 'kubelet', 'tectonicConsole'].map(key => {
      const value = status ? status[key] : {ready: false, instance: ''};
      const err = value.error || _.includes(ERROR_STATUSES, value.message);
      return (
        <WaitingLi done={value.ready} error={err} key={key}>
          <span title={value.instance}>
            {componentDescriptions[key]}
          </span>
          {err && <span>: {value.message}</span> }
          { !value.ready && substeps[key] &&
            <ul>{substeps[key]}</ul>
          }
        </WaitingLi>
      );
    }));

    const cloudFormationId = status && status.cloudFormation && status.cloudFormation.id;

    return (
      <div className="row">
        <div className="col-sm-12">
          <div className="form-group">
            Kubernetes is starting up. We're commiting your cluster details.
            Grab some tea and sit tight. This process can take up to 20 minutes.
            Status updates will appear below or you can check in the <CloudFormationAnchor arn={cloudFormationId} region={this.props.region}>AWS Console</CloudFormationAnchor>.
          </div>
          <div className="form-group">
            <div className="wiz-launch__progress">
              <ul className="wiz-launch-progress">
                {statusList}
              </ul>
            </div>
          </div>
          <div className="from-group">
            <SaveAssets blob={blob} />
          </div>
          { this.state.remoteAddr && <SSHInstructions controllerIP={this.state.remoteAddr} /> }
          { error && <div className='wiz-error-message'>{error.toString()}</div> }
          <br />
        </div>
      </div>
    );
  }
});

AWS_PowerOn.canNavigateForward = ({cluster}) => {
  if (!cluster.ready) {
    return false;
  }

  return _.every(cluster.status, v => v.ready);
};
