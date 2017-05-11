import _ from 'lodash';
import React from 'react';
import { Map as ImmutableMap, Set as ImmutableSet } from 'immutable';

import { commitPhases } from './actions';
import { PLATFORM_TYPE, DRY_RUN } from './cluster-config';

import { CertificateAuthority } from './components/certificate-authority';
import { ClusterType } from './components/cluster-type';
import { DryRun } from './components/dry-run';
import { Etcd } from './components/etcd';
import { SubmitDefinition } from './components/submit-definition';
import { Success } from './components/success';
import { TF_PowerOn } from './components/tf-poweron';
import { Users } from './components/users';

import { BM_ClusterInfo } from './components/bm-cluster-info';
import { BM_Controllers, BM_Workers } from './components/bm-nodeforms';
import { BM_Credentials } from './components/bm-credentials';
import { BM_Hostname } from './components/bm-hostname';
import { BM_Matchbox } from './components/bm-matchbox';
import { BM_NetworkConfig } from './components/bm-network-config';
import { BM_SSHKeys } from './components/bm-sshkeys';

import { AWS_CloudCredentials } from './components/aws-cloud-credentials';
import { AWS_ClusterInfo } from './components/aws-cluster-info';
import { AWS_DefineNodes } from './components/aws-define-nodes';
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
  hideSave: true,
  showRestore: true,
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
  hideSave: true,
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
  component: BM_ClusterInfo,
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
  component: BM_NetworkConfig,
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
  component: AWS_DefineNodes,
  title: 'Define Nodes',
};

const awsSubmitKeysPage = {
  path: '/define/aws/keys',
  component: AWS_SubmitKeys,
  title: 'Submit Keys',
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

export const sections = new Map([
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
    bmNetworkConfigPage,
    bmControllersPage,
    bmWorkersPage,
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

// Lets us do 'sections.defineAWS' instead of using sections.get('defineAWS')
sections.forEach((v, k) => {
  sections[k] = v;
});

// A Trail is an immutable representation of the navigation options available to a user.
// Navigation involves
//    - moving from one page to the next page
//    - moving from one page to some previous pages
//    - presenting a (possibly disabled) list of all pages
export class Trail {
  constructor(trailSections, whitelist, opts={}) {
    this.canReset = opts.canReset;
    this.sections = trailSections;
    const sectionPages = this.sections.reduce((ls, l) => ls.concat(l), []);
    sectionPages.forEach(p => {
      if (!p.component) {
        throw new Error(`${p.title} has no component`);
      }
    });
    this._pages = !whitelist ? sectionPages : sectionPages.filter(p => whitelist.includes(p));
    this.ixByPage = this._pages.reduce((m, v, ix) => m.set(v, ix), ImmutableMap());
    this.pageByPath = this._pages.reduce((m, v) => m.set(v.path, v), ImmutableMap());
  }

  sectionFor(page) {
    return _.find(this.sections, s => s.includes(page));
  }

  navigable(page) {
    return this.ixByPage.has(page);
  }

  // Given a path, return a "better" path for this trail. Will not navigate to that path.
  fixPath(path, state) {
    const page = this.pageByPath.get(path);
    if (!page) {
      return this._pages[0].path;
    }

    // If the NEXT button on a previous page is disabled, you
    // shouldn't be able to see this page. Show the first page, before
    // or including this one, that won't allow forward navigation.
    const pred = this._pages.find(p => {
      return p === page || !(p.component.canNavigateForward ? p.component.canNavigateForward(state) : true);
    });

    return pred.path;
  }

  // Returns -1, 0, or 1, if a is before, the same as, or after b in the trail
  cmp(a, b) {
    return this.ixByPage.get(a) - this.ixByPage.get(b);
  }

  canNavigateForward(pageA, pageB, state) {
    const a = this.ixByPage.get(pageA);
    const b = this.ixByPage.get(pageB);

    if (!_.isFinite(a) || !_.isFinite(b)) {
      return false;
    }

    const start = Math.min(a, b);
    const end = Math.max(a, b);

    for (let i = start; i < end; i++) {
      const component = this._pages[i].component;
      if (!component.canNavigateForward) {
        continue;
      }
      if (!component.canNavigateForward(state)) {
        return false;
      }
    }
    return true;
  }

   // Returns the previous page in the trail if that page exists
  previousFrom(page) {
    const myIx = this.ixByPage.get(page);
    return this._pages[myIx - 1];
  }

  // Returns the next page in the trail if that exists
  nextFrom(page) {
    const myIx = this.ixByPage.get(page);
    return this._pages[myIx + 1];
  }
}

const doingStuff = ImmutableSet([commitPhases.REQUESTED, commitPhases.WAITING]);

const platformToSection = {
  [AWS_TF]: {
    choose: new Trail([sections.choose, sections.defineAWS]),
    define: new Trail([sections.defineAWS], [submitDefinitionPage]),
    boot: new Trail([sections.bootAWSTF]),
  },
  [BARE_METAL_TF]: {
    choose: new Trail([sections.choose, sections.defineBaremetal]),
    define: new Trail([sections.defineBaremetal], [submitDefinitionPage]),
    boot: new Trail([sections.bootBaremetalTF]),
  },
};

export const trail = ({cluster, clusterConfig, commitState}) => {
  let platform = clusterConfig[PLATFORM_TYPE];
  const { ready } = cluster;

  if (platform === '') {
    return new Trail([sections.loading]);
  }
  if (!isSupported(platform)) {
    return new Trail([sections.choose]);
  }

  const { phase } = commitState;
  const submitted = ready || (phase === commitPhases.SUCCEEDED);
  if (submitted) {
    if (clusterConfig[DRY_RUN]) {
      return new Trail([sections.bootDryRun], null, {canReset: true});
    }
    return platformToSection[platform].boot;
  }
  if (doingStuff.has(phase)) {
    return platformToSection[platform].define;
  }

  return platformToSection[platform].choose;
};

export const getAllRoutes = () => {
  // No components have the same path, so this is safe.
  // If a user guesses an invalid URL, they could get in a weird state. Oh well.
  let routes = [];
  _.each(sections, v => {
    _.each(v, o => {
      routes.push(o);
    });
  });
  return routes;
};
