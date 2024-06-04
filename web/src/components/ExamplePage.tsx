import * as React from 'react';
import axios from 'axios';
import  { useEffect } from 'react';

import './example.css';
// import design system element definitions,
// which auto-register their tagnames once executed
import '@rhds/elements/rh-accordion/rh-accordion.js';
import '@rhds/elements/rh-tabs/rh-tabs.js';
import { Page, PageSection } from '@patternfly/react-core';

const ExamplePage: React.FC = () => {
  useEffect(() => {
    fetchPods();
  }, []);
  const fetchPods = () => {
    axios.get('https://cnf-certsuite-plugin.cnf-certsuite-plugin-ns.svc.cluster.local:9443/api/pods')
      .then(response => {
        const data = response.data;
        const podList = document.getElementById('podListAcc');
        const logContainer = document.getElementById('logContainer');
        const backButton = document.getElementById('backButton');

        if (podList) podList.innerHTML = ''; // Clear existing list
        if (logContainer) logContainer.innerHTML = ''; // Clear log container
        if (backButton) backButton.style.display = 'none'; // Hide back button

        data.forEach((pod: { name: string, containers: string[] }) => {
          console.log(pod)
          const listItem = document.createElement('rh-accordion-header');
          listItem.innerHTML = `<h2> ${pod.name} </h2>`;
          listItem.id = `pod-${pod.name}`;
          podList?.appendChild(listItem);
          displayContainers(pod.name, pod.containers);
        });
      })
      .catch(error => console.error('Error fetching pods:', error));
  };

  const displayContainers = (podName: string, containers: string[]) => {
    const podList = document.getElementById(`pod-${podName}`);
    if (!podList) return;
    const containerList = document.createElement('rh-accordion-panel');
    containers.forEach(container => {
      const accordion = document.createElement('rh-accordion');
      accordion.id = `acc-${podName}`;
      const listItem = document.createElement('rh-accordion-header');
      listItem.id = `con-${podName}-${container}`;
      listItem.innerHTML = `<h2> ${container} </h2>`;
      accordion.appendChild(listItem);
      containerList.appendChild(accordion);
      fetchPodLogs(podName, container);
    });

    podList.insertAdjacentElement('afterend', containerList);
  };

  const fetchPodLogs = (podName: string, containerName: string) => {
    axios.get(`https://cnf-certsuite-plugin.cnf-certsuite-plugin-ns.svc.cluster.local:9443/api/logs/${podName}/${containerName}`)
      .then(response => {
        const data = response.data;
        const logContainer = document.createElement('rh-accordion-panel');
        const podList = document.getElementById(`con-${podName}-${containerName}`);
        if (podList) {
          logContainer.innerHTML = `<pre>${data}</pre>`;
          podList.insertAdjacentElement('afterend', logContainer);
        }
      })
      .catch(error => console.error(`Error fetching logs for pod ${podName}:`, error));
  };

  return (
    <div>
      <head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>CNF Certification Test</title>
        <link rel="shortcut icon" type="image/svg+xml" sizes="any" href="https://ux.redhat.com/assets/logo-red-hat.svg" />
        <link rel="stylesheet" href="https://ux.redhat.com/assets/packages/@rhds/elements/elements/rh-table/rh-table-lightdom.css" />
        <link rel="stylesheet" href="https://ux.redhat.com/assets/packages/@rhds/elements/elements/rh-footer/rh-footer-lightdom.css" />
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
        <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.0.13/css/all.css"
          integrity="sha384-DNOHZ68U8hZfKXOrtjWvjxusGo9WQnrNx2sqG0tfsghAvtVlRW3tvkXWZh58N9jp" crossOrigin="anonymous" />
        <script type="importmap">
          {`{
            "imports": {
              "@rhds/elements/": "https://ga.jspm.io/npm:@rhds/elements@1.2.0/elements/",
              "@rhds/elements/lib/": "https://ga.jspm.io/npm:@rhds/elements@1.2.0/elements/lib/",
              "@patternfly/elements/": "https://ga.jspm.io/npm:@patternfly/elements@2.4.0/"
            },
            "scopes": {
              "https://ga.jspm.io/": {
                "@lit/reactive-element": "https://ga.jspm.io/npm:@lit/reactive-element@1.6.3/reactive-element.js",
                "@lit/reactive-element/decorators/": "https://ga.jspm.io/npm:@lit/reactive-element@1.6.3/decorators/",
                "@patternfly/elements/": "https://ga.jspm.io/npm:@patternfly/elements@2.4.0/",
                "@patternfly/pfe-core": "https://ga.jspm.io/npm:@patternfly/pfe-core@2.4.1/core.js",
                "@patternfly/pfe-core/": "https://ga.jspm.io/npm:@patternfly/pfe-core@2.4.1/",
                "@rhds/tokens/media.js": "https://ga.jspm.io/npm:@rhds/tokens@1.1.2/js/media.js",
                "lit": "https://ga.jspm.io/npm:lit@2.8.0/index.js",
                "lit-element/lit-element.js": "https://ga.jspm.io/npm:lit-element@3.3.3/lit-element.js",
                "lit-html": "https://ga.jspm.io/npm:lit-html@2.8.0/lit-html.js",
                "lit-html/": "https://ga.jspm.io/npm:lit-html@2.8.0/",
                "lit/": "https://ga.jspm.io/npm:lit@2.8.0/",
                "tslib": "https://ga.jspm.io/npm:tslib@2.6.2/tslib.es6.mjs"
              },
              "https://ga.jspm.io/npm:@patternfly/elements@2.4.0/": {
                "lit": "https://ga.jspm.io/npm:lit@2.6.1/index.js",
                "lit/": "https://ga.jspm.io/npm:lit@2.6.1/"
              }
            }
          }`}
        </script>
      </head>
      
      <Page>
        <PageSection variant="light">
          <div style={{ "--main-opacity": "0" } as React.CSSProperties}>
            <h1>List of Pods</h1>
            <rh-tabs>
              <rh-tab id="backButton" slot="tab" onClick={fetchPods}>Pod List</rh-tab>
              <rh-tab-panel>
                <rh-accordion id="podListAcc">
                  <rh-accordion-header>
                    <h4>Item One</h4>
                  </rh-accordion-header>
                  <rh-accordion-panel>
                    <rh-accordion-header>
                      <h4>Item One</h4>
                    </rh-accordion-header>
                    <rh-accordion-panel>
                      <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
                    </rh-accordion-panel>
                  </rh-accordion-panel>
                </rh-accordion>
              </rh-tab-panel>
            </rh-tabs>
            <ul id="podList"></ul>
            <div id="logContainer"></div>
          </div>
        </PageSection>
      </Page>
    </div>
  );
}
export default ExamplePage;
