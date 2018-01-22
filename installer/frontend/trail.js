import _ from 'lodash';
import React from 'react';
import { Map as ImmutableMap, Set as ImmutableSet } from 'immutable';

import { commitPhases } from './actions';
import { PLATFORM_TYPE, DRY_RUN } from './cluster-config';

import { CertificateAuthority } from './components/certificate-authority';
import { ClusterInfo } from './components/cluster-info';
import { ClusterType } from './components/cluster-type';
import { DryRun } from './components/dry-run';
import { SubmitDefinition } from './components/submit-definition';
import { Success } from './components/success';
import { TF_PowerOn } from './components/tf-poweron';
import { Users } from './components/users';

import { BM_Controllers, BM_Workers } from './components/bm-nodeforms';
import { BM_Credentials } from './components/bm-credentials';
import { BM_Hostname } from './components/bm-hostname';
import { BM_Matchbox } from './components/bm-matchbox';
import { KubernetesCIDRs } from './components/k8s-cidrs';
import { BM_SSHKeys } from './components/bm-sshkeys';

import { AWS_CloudCredentials } from './components/aws-cloud-credentials';
import { AWS_ClusterInfo } from './components/aws-cluster-info';
import { AWS_Nodes, Etcd } from './components/nodes';
import { AWS_SubmitKeys } from './components/aws-submit-keys';
import { AWS_VPC } from './components/aws-vpc';

import {
  AWS_TF,
  BARE_METAL_TF,
  isSupported,
} from './platforms';

// Common pages
const certificateAuthorityPage = {
  path: '/define/certificate-authority',
  component: CertificateAuthority,
  title: 'Certificate Authority',
};

const clusterTypePage = {
  path: '/define/cluster-type',
  component: ClusterType,
  title: 'Platform',
};

const dryRunPage = {
  path: '/define/advanced',
  component: DryRun,
  title: 'Download Assets',
};

const etcdPage = {
  path: '/define/etcd',
  component: Etcd,
  title: 'etcd Connection',
};

const submitDefinitionPage = {
  path: '/define/submit',
  component: SubmitDefinition,
  title: 'Submit',
  hidePager: true,
};

const successPage = {
  path: '/boot/complete',
  component: Success,
  title: 'Installation Complete',
  hidePager: true,
};

const TFPowerOnPage = {
  path: '/boot/tf/poweron',
  component: TF_PowerOn,
  title: 'Start Installation',
};

const usersPage = {
  path: '/define/users',
  component: Users,
  title: 'Console Login',
};

// Baremetal pages
const bmClusterInfoPage = {
  path: '/define/cluster-info',
  component: ClusterInfo,
  title: 'Cluster Info',
};

const bmControllersPage = {
  path: '/define/controllers',
  component: BM_Controllers,
  title: 'Define Masters',
};

const bmCredentialsPage = {
  path: '/define/credentials',
  component: BM_Credentials,
  title: 'Matchbox Credentials',
};

const bmHostnamePage = {
  path: '/define/hostname',
  component: BM_Hostname,
  title: 'Cluster DNS',
};

const bmMatchboxPage = {
  path: '/define/matchbox',
  component: BM_Matchbox,
  title: 'Matchbox Address',
};

const bmNetworkConfigPage = {
  path: '/define/network-config',
  component: KubernetesCIDRs,
  title: 'Network Configuration',
};

const bmSshKeysPage = {
  path: '/define/ssh-keys',
  component: BM_SSHKeys,
  title: 'SSH Key',
};

const bmWorkersPage = {
  path: '/define/workers',
  component: BM_Workers,
  title: 'Define Workers',
};

// AWS pages
const awsClusterInfoPage = {
  path: '/define/aws/cluster-info',
  component: AWS_ClusterInfo,
  title: 'Cluster Info',
};

const awsCloudCredentialsPage = {
  path: '/define/aws/cloud-credentials',
  component: AWS_CloudCredentials,
  title: 'AWS Credentials',
};

const awsDefineNodesPage = {
  path: '/define/aws/nodes',
  component: AWS_Nodes,
  title: 'Define Nodes',
};

const awsSubmitKeysPage = {
  path: '/define/aws/keys',
  component: AWS_SubmitKeys,
  title: 'SSH Key',
};

const awsVPCPage = {
  path: '/define/aws/vpc',
  component: AWS_VPC,
  title: 'Networking',
};

// This component is never visible!
const loadingPage = {
  path: '/',
  component: React.createElement('div'),
  title: 'Tectonic',
};

export const trailSections = new Map([
  ['loading', [
    loadingPage,
  ]],
  ['choose', [
    clusterTypePage,
  ]],
  ['defineBaremetal', [
    bmClusterInfoPage,
    bmHostnamePage,
    certificateAuthorityPage,
    bmMatchboxPage,
    bmCredentialsPage,
    bmControllersPage,
    bmWorkersPage,
    bmNetworkConfigPage,
    etcdPage,
    bmSshKeysPage,
    usersPage,
    submitDefinitionPage,
  ]],
  ['defineAWS', [
    awsCloudCredentialsPage,
    awsClusterInfoPage,
    certificateAuthorityPage,
    awsSubmitKeysPage,
    awsDefineNodesPage,
    awsVPCPage,
    usersPage,
    submitDefinitionPage,
  ]],
  ['bootBaremetalTF', [
    TFPowerOnPage,
    successPage,
  ]],
  ['bootAWSTF', [
    TFPowerOnPage,
    successPage,
  ]],
  ['bootDryRun', [
    dryRunPage,
  ]],
]);

// Lets us do 'trailSections.defineAWS' instead of using trailSections.get('defineAWS')
trailSections.forEach((v, k) => {
  trailSections[k] = v;
});

// A Trail is an immutable representation of the navigation options available to a user.
class Trail {
  constructor (sections, whitelist) {
    this.sections = sections;
    const sectionPages = this.sections.reduce((ls, l) => ls.concat(l), []);
    sectionPages.forEach(p => {
      if (!p.component) {
        throw new Error(`${p.title} has no component`);
      }
    });
    this.pages = !whitelist ? sectionPages : sectionPages.filter(p => whitelist.includes(p));
    this.ixByPage = this.pages.reduce((m, v, ix) => m.set(v, ix), ImmutableMap());
    this.pageByPath = this.pages.reduce((m, v) => m.set(v.path, v), ImmutableMap());
  }

  // Returns the previous page in the trail if that page exists
  previousFrom (page) {
    const myIx = this.ixByPage.get(page);
    return this.pages[myIx - 1];
  }

  // Returns the next page in the trail if that exists
  nextFrom (page) {
    const myIx = this.ixByPage.get(page);
    return this.pages[myIx + 1];
  }
}

const doingStuff = ImmutableSet([commitPhases.REQUESTED, commitPhases.WAITING]);

const platformToSection = {
  [AWS_TF]: {
    choose: new Trail([trailSections.choose, trailSections.defineAWS]),
    define: new Trail([trailSections.defineAWS], [submitDefinitionPage]),
    boot: new Trail([trailSections.bootAWSTF]),
  },
  [BARE_METAL_TF]: {
    choose: new Trail([trailSections.choose, trailSections.defineBaremetal]),
    define: new Trail([trailSections.defineBaremetal], [submitDefinitionPage]),
    boot: new Trail([trailSections.bootBaremetalTF]),
  },
};

export const trail = ({cluster, clusterConfig, commitState}) => {
  const platform = clusterConfig[PLATFORM_TYPE];

  if (platform === '') {
    return new Trail([trailSections.loading]);
  }
  if (!isSupported(platform)) {
    return new Trail([trailSections.choose]);
  }

  const {phase} = commitState;
  const submitted = cluster.ready || (phase === commitPhases.SUCCEEDED);
  if (submitted) {
    // If we detected a dry run in progress when the app started, then clusterConfig will not be populated, so also
    // check for a Terraform "show" (dry run) action
    if (clusterConfig[DRY_RUN] || _.get(cluster, 'status.terraform.action') === 'show') {
      return new Trail([trailSections.bootDryRun]);
    }
    return platformToSection[platform].boot;
  }
  if (doingStuff.has(phase)) {
    return platformToSection[platform].define;
  }

  return platformToSection[platform].choose;
};
