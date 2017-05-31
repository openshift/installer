import React from 'react';
import { Dropdown } from './ui';

export class Header extends React.Component {
  render() {
    const hasNewVersion = GIT_TAG !== GIT_LATEST_TAG;
    const productDdItems = {
      'Tectonic - Kubernetes': 'http://coreos.com/tectonic/',
      'Quay - Registry': 'http://coreos.com/quay-enterprise',
      'Container Linux Support': 'http://coreos.com/products/container-linux-subscription/',
      'Training': 'http://coreos.com/training/',
    };

    const openSourceDdItems = {
      'Open Source Docs': 'http://coreos.com/docs/',
      'Kubernetes': 'http://coreos.com/kubernetes',
      'Operators': 'http://coreos.com/operators/',
      'Container Linux': 'http://coreos.com/os/docs/latest',
      'rkt': 'http://coreos.com/rkt/',
      'etcd': 'http://coreos.com/etcd/',
      'Clair': 'http://coreos.com/clair/',
      'flannel': 'http://coreos.com/flannel/',
      'Ignition': 'http://coreos.com/ignition/',
      'Matchbox': 'http://coreos.com/matchbox/',
      '90+ more on Github': 'https://github.com/coreos/',
    };

    return <div className="co-navbar">
      <div className="co-navbar__header">
        <a href="/" className="co-navbar__logo-link">
          <img className="co-navbar__logo" src="/frontend/lib/mochi/img/logo/coreos/logo.svg" />
        </a>
      </div>
      <div>
        <ul className="co-navbar-nav">
          <li>
            <Dropdown items={productDdItems} header="Product"/>
          </li>
          <li>
            <Dropdown items={openSourceDdItems} header="Open Source"/>
          </li>
          <li className="tectonic-dropdown-menu-title">
            <a href="http://coreos.com/tectonic/docs/latest/" target="_blank" className="tectonic-dropdown-menu-title__link">Tectonic Docs</a>
          </li>
        </ul>
        <div className="co-navbar--right">
          <ul className="co-navbar-nav">
            {hasNewVersion && <li className="co-navbar-nav-item__version">
              <span className="co-navbar-nav-item__version--new">
                New installer version: <a href="http://coreos.com/tectonic/releases/" target="_blank">Release notes {GIT_LATEST_TAG}</a>
              </span>
            </li>}
            <li className="co-navbar-nav-item__version">
              <span>Version: {GIT_TAG}</span>
            </li>
          </ul>
        </div>
      </div>
    </div>;
  }
}

