import React from 'react';
import { connect } from 'react-redux';
import { Set as ImmutableSet } from 'immutable';
import _ from 'lodash';

import * as awsActions from '../aws-actions';
import { AWS_HOSTED_ZONE_ID, CLUSTER_SUBDOMAIN } from '../cluster-config';
import { Alert } from './alert';
import { TectonicGA } from '../tectonic-ga';
import { toExtraData } from '../utils';

const stateToProps = ({aws, clusterConfig}) => {
  const hostedZoneID = clusterConfig[AWS_HOSTED_ZONE_ID];
  const isPrivate = _.get(clusterConfig, toExtraData(AWS_HOSTED_ZONE_ID) + '.privateZones.' + hostedZoneID);
  return {
    hostedZoneID,
    isPrivate,
    domainInfo: aws.domainInfo.value,
    domain: _.get(clusterConfig, toExtraData(AWS_HOSTED_ZONE_ID) + '.zoneToName.' + hostedZoneID),
    clusterSubdomain: clusterConfig[CLUSTER_SUBDOMAIN],
  };
};

const dispatchToProps = dispatch => ({
  getDomainInfo: body => dispatch(awsActions.getDomainInfo(body)),
});

class DomainInfo extends React.Component {
  getDomainInfo (props) {
    const { domain, hostedZoneID } = props;
    if (!hostedZoneID || !domain) {
      return;
    }
    this.props.getDomainInfo({id: hostedZoneID, name: domain});
  }

  componentDidMount () {
    this.getDomainInfo(this.props);
  }

  componentWillReceiveProps (nextProps) {
    if (nextProps.hostedZoneID === this.props.hostedZoneID && nextProps.domain === this.props.domain) {
      return;
    }
    this.getDomainInfo(nextProps);
  }

  render () {
    const warnings = [];
    const {clusterSubdomain, domainInfo, domain, isPrivate} = this.props;
    const {soaTTL, soaValue, registered, awsNS, publicNS} = domainInfo;

    if (registered && registered.startsWith('AVAILABLE')) {
      warnings.push(<Alert severity="error" key="registered">{domain} is not registered!</Alert>);
    }

    if (awsNS && awsNS.length === 0) {
      warnings.push(<Alert severity="error" key="noAWSNS">{domain} has no NS records in Route53.</Alert>);
    }

    if (!isPrivate) {
      if (publicNS && publicNS.length === 0) {
        warnings.push(<Alert severity="error" key="noPublicNS">{domain} has no NS records.</Alert>);
      }

      if (!ImmutableSet(publicNS).subtract(awsNS).isEmpty()) {
        const preStyle = {marginLeft: 15, marginTop: 0};
        warnings.push(<Alert severity="error" key="wrongNS">
          Public NS records for {domain} do not match Route 53. You should update your domain registrar's NS records to match.
          <br /><br />
          Public:
          <pre style={preStyle}>{publicNS.join('\n')}</pre>
          <br />
          Route 53:
          <pre style={preStyle}>{awsNS.join('\n')}</pre>
        </Alert>);
      }
    }

    const minimumTTL = soaValue && _.toInteger(_.last(soaValue.split(' ')));
    if (soaTTL > 300 && minimumTTL > 300) {
      warnings.push(<Alert key="soa">
        <b>{domain}'s SOA TTL is {soaTTL} seconds and its SOA minimum TTL is {minimumTTL} seconds.</b>&nbsp;
        The SOA record TTL and minimum TTL values determine how long to cache NXDOMAIN responses for. Installation cannot complete until {clusterSubdomain}-k8s.{domain} resolves.&nbsp;
        {/* eslint-disable react/jsx-no-target-blank */}
        <a href="https://coreos.com/tectonic/docs/latest/install/aws/troubleshooting.html#route53-dns-resolution" onClick={() => TectonicGA.sendDocsEvent('aws-tf')} rel="noopener" target="_blank">Read more here</a>.
        {/* eslint-enable react/jsx-no-target-blank */}
      </Alert>);
    }

    if (warnings.length) {
      return <div>{warnings}</div>;
    }

    return null;
  }
}

export const AWS_DomainValidation = connect(stateToProps, dispatchToProps)(DomainInfo);
