import { isReleaseVersion } from './utils';

const send = (obj) => {
  if (!isReleaseVersion()) {
    return;
  }

  try {
    const ga = window[window.GoogleAnalyticsObject || 'ga'];
    if (typeof ga !== 'function') {
      throw new Error('ga is not a function!');
    }
    if (obj.type === 'pageview') {
      ga('TectonicInstaller.send', 'pageview', obj.page);
      ga('CoreOS.send', 'pageview', obj.page);
    } else if (obj.type === 'event') {
      const {category, action, label, value} = obj;
      ga('TectonicInstaller.send', 'event', category, action, label, value);
      ga('CoreOS.send', 'event', category, action, label, value);
    }
  } catch (err) {
    console.error(`Failed to send GA event: ${err.message}`);
  }
};

export const TectonicGA = {
  initialize: () => {

    // https://developers.google.com/analytics/devguides/collection/analyticsjs/
    /* eslint-disable */
    (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
      (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
      m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
    })(window,document,'script','https://www.google-analytics.com/analytics.js','ga')
    /* eslint-enable */

    const ga = window.ga;

    ga('create', 'UA-42684979-10', 'none', 'TectonicInstaller');
    ga('create', 'UA-42684979-6', 'none', 'CoreOS');
  },

  sendPageView: (page) => {
    send({ type: 'pageview', page});
  },

  sendError: (message, stack = '') => {
    send({
      type: 'event',
      category: 'installerError',
      label: `${GIT_TAG} ${message} Stack: ${stack}`,
    });
  },

  sendEvent: (category, action, label = '', platform = '') => {
    send({
      type: 'event',
      category, action,
      label: `${platform}-${GIT_TAG} ${label}`,
    });
  },

  sendDocsEvent: (platform = '') => {
    send({
      type: 'event',
      category: 'Installer Docs Link',
      action: 'click',
      label: `${platform}-${GIT_TAG} User clicks on documentation link`,
    });
  },
};
