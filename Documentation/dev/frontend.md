# Frontend Architecture
## Tools
The frontend is a single React app written in ES6 and compiled with Babel. eslint is used to enforce some coding style rules.

## Directory Structure
* `/installer/api/`: Go code for the backend.
* `/installer/assets/frontend/`: Assets for the frontend such as the CSS and the compiled JavaScript.
* `/installer/bin/`: Compiled backend binary and generated assets for clusters created with the GUI.
* `/installer/frontend/`: All frontend JavaScript code.
* `/installer/frontend/ui-tests/`: Frontend Nightwatch end-to-end tests.
* `/installer/frontend/__tests__/`: Frontend Jest unit tests.
* `/installer/frontend/__tests__/examples/`: Input files for both the unit tests and end-to-end tests (kept in the same directory because some input files are used by both the unit and end-to-end tests).

## Testing
The frontend has both Jest unit tests and Nightwatch end-to-end tests, but we make more use of end-to-end tests with the unit tests currently just testing some `Form` / `Field` functionality.

Some end-to-end tests output .tfvars files that match those fed into smoke tests. For example, test input `/installer/frontend/__tests__/examples/aws.progress` generates .tfvars matching `/tests/smoke/aws/vars/aws.tfvars.json`. However, a few tfvars are missing from one file or the other for different reasons, so those are ignored by the GUI test.

## Wizard Navigation
We use `Trail`s to define the sequence of pages for the wizard and the React component for each page. React Router is used to map URLs to those page components.

Each page component has a `canNavigateForward()` function that determines whether the page’s Next button is enabled.

The app automatically jumps back to the first page before the current page that has a disabled Next button. This should only happen if the user changes the URL manually. Similarly, the sidebar has links to each page in the current `Trail`, but they are only enabled up to and including the first page that has a disabled Next button.

## State / Redux
As the user progresses through the wizard, field values and related data are all stored in Redux state under the following keys.
* `aws`: For tracking the status and response of AWS related API calls.
* `cluster`: Tracks the status of the cluster during installation.
* `clusterConfig`: The actual field values that will be used to generate the .tfvars file (often abbreviated to `cc` in the code).
* `clusterConfig.error`: All field validation errors, including “form level” validation errors.
* `clusterConfig.inFly`: Flags indicating which fields / forms are currently being validated (field `validator()` functions can be `async`).
* `clusterConfig.extra`: Data fetched by `getExtraStuff()` functions.
* `clusterConfig.extraError`: Errors while executing `getExtraStuff()` functions.
* `clusterConfig.extraInFly`: Flags indicating which `getExtraStuff()` calls are in progress.
* `commitState`: Tracks the status of the `/terraform/apply` API call, which actually initiates the Terraform apply action.
* `dirty`: Tracks which fields have been changed. Validation errors are only shown for fields marked as `dirty`.
* `eventErrors`: Another place where errors are stored. Almost unused and should probably be refactored away.
* `serverFacts`: Stores additional data loaded from the backend when the app launches.

We use `sessionStorage` instead of `localStorage` to persist the installer state, so all field values are lost when the tab is closed. The `sessionStorage` data is updated periodically and also when the `beforeunload` event fires, so that no field values are lost when refreshing the page.

The Start Over link provides a way to clear the app state without closing the tab. First `sessionStorage` is cleared and then the app is reloaded to clear the Redux state. This does not cancel any cluster install or destroy actions that may be in progress.

Redux actions (`actions.js`, `aws-actions.js`) are currently defined inconsistently with some using action creator functions and some not and also with inconsistent naming and parameter ordering. It would be good to refactor, e.g. by consistently using action creators.

## App Initialization
When the installer GUI starts, code in `app.jsx` goes through the initialization steps below.
1. Restore state from `sessionStorage`, if it exists.
1. Router initialization.
1. Query `/tectonic/facts` to load additional data required by the frontend.
1. Validate all forms and fields so that validation errors are immediately displayed for any invalid values that were loaded from `sessionStorage`.
1. Query `/tectonic/status` to see if an install is already in progress, and if there is, jump to the Start Installation page to display the installation progress. This can happen even if there is no data in `sessionStorage`, so the app needs to handle displaying Installation progress without relying on the existence of any field values in Redux.

## Forms / Fields
Forms are created using the `Form` and `Field` classes, which both take the following options.
* `ignoreWhen`
  * Function that returns a bool indicating whether the field’s validation should be ignored.
* `validator`
  * Function to test whether this form / field is valid. `validate.js` contains various useful validation functions that can be used here and a `compose()` helper for combining multiple validator functions.
  * Fields are valid if their `ignoreWhen()` returns `true` or if their `validator()` returns `true`.
  * Forms are valid if their `ignoreWhen()` returns `true` or if both their own `validator()` and all of their fields' `validator()`s return `true`.
* `dependencies`
  * Other forms and fields that this form / field depends on. Used to delay calling `getExtraStuff()` until we have any necessary values from earlier wizard steps.
* `getExtraStuff`
  * Function to fetch additional data for the form / field. For example, a list of options for a dropdown that depends on an earlier field and therefore cannot be fetched when the app first loads. `getExtraStuff()` is not called until any `dependencies` have passed validation.

## cluster-config.js
Converts the Redux `clusterConfig` into actual Terraform variables for each platform.

Also defines constants for string IDs. In general, if a string ID is used more than once in the code, it should be converted to a constant here.

There is also a `DEFAULT_CLUSTER_CONFIG` object, which was used to specify default values for all fields. Now, defaults are specified by passing an option to `new Field()` instead, but several bits of code still rely on `DEFAULT_CLUSTER_CONFIG`, so they will need to be changed before this can be removed completely.

## Input Components
Various input field components are defined in `components/ui.jsx`, including
* `CIDR`
  * A text input for inputting IP address ranges in CIDR notation.
* `Connect`
  * Wraps input field components in order to automatically inject props that provide the field’s Redux data.
* `Deselect` / `DeselectField`
  * For inputs than can optionally be deselected. `Deselect` is the checkbox that controls whether the field is deselected. `DeselectField` wraps the field that should be deselected.
* `Select`
  * A dropdown input.
* `Selector`
  * `Select` with an associated Redux `extraData` entry and optionally with a refresh button that triggers reloading the dropdown options.
* `AsyncSelect`
  * Similar to `Selector`, so merging these two components would be a good area for refactoring.
* `FileArea`
  * A combined file input and textarea for either uploading a file or entering the text directly.
* `ConnectedFieldList`
  * For expandable lists of fields (e.g. key value pairs).

## Interaction with Backend
All interaction with the backend is via a JSON API using `fetch()` promises. Most API calls are made from Redux actions which then call `dispatch()` with the response data.

## Install Screen
Once the wizard is complete, `/terraform/apply` is called to send the results to the backend. After that, `/tectonic/status` is called periodically to monitor the installation progress. The `/tectonic/status` response includes the Terraform log output, which is then shown in the GUI. The log output is also used to generate a progress bar by searching the logs to estimate how much progress has been made.

Once the Terraform step is complete, the Tectonic install step begins. This step has several sub-steps, with the current status of each also being indicated by the `/tectonic/status` response.
