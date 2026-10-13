/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Manages the collection of labels of a cluster.
resource Labels {
	// Retrieves the list of labels.
	method List {
		// Index of the requested page, where one corresponds to the first page.
		in out Page Integer = 1

		// Number of items contained in the returned page.
		in out Size Integer = 100

		// Total number of items of the collection.
		out Total Integer

		// Retrieved list of labels.
		out Items []Label
	}

	// Adds a new label to the cluster.
	method Add {
		// Description of the label.
		in out Body Label
	}

	// Reference to the service that manages an specific label.
	locator Label {
		target Label
		variable ID
	}
}
